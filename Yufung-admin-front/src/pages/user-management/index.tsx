import React, { useState, useRef } from 'react';
import { PageContainer } from '@ant-design/pro-components';
import type { ProColumns, ActionType } from '@ant-design/pro-components';
import { ProTable } from '@ant-design/pro-components';
import { Button, message, Modal, Tag, Dropdown, Space } from 'antd';
import type { MenuProps } from 'antd';
import { PlusOutlined, DownOutlined, ExportOutlined } from '@ant-design/icons';
import { getUserList, deleteUser, batchUpdateUserStatus, exportUsers, downloadUserTemplate, previewUserImport, importUsers } from '@/services/ant-design-pro/userApi';
import UserForm from './components/UserForm';
import UserImportModal from './components/UserImportModal';
import type { UserItem, UserListParams } from '@/services/ant-design-pro/userApi';

const UserManagement: React.FC = () => {
  const [createModalOpen, setCreateModalOpen] = useState<boolean>(false);
  const [updateModalOpen, setUpdateModalOpen] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<UserItem>();
  const [selectedRowsState, setSelectedRows] = useState<UserItem[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const actionRef = useRef<ActionType>(null);

  // 导入相关状态
  const [importModalOpen, setImportModalOpen] = useState<boolean>(false);
  const [importLoading, setImportLoading] = useState<boolean>(false);

  // 下载模板
  const handleDownloadTemplate = async () => {
    try {
      const blob = await downloadUserTemplate('xlsx');
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `用户导入模板_${new Date().toISOString().slice(0, 10)}.xlsx`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
      message.success('模板下载成功');
    } catch (error) {
      message.error('模板下载失败');
    }
  };

  // 处理文件导入
  const handleImport = async (file: File, options: { skipHeader: boolean; updateExisting: boolean }) => {
    setImportLoading(true);
    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('skip_header', options.skipHeader.toString());
      formData.append('update_existing', options.updateExisting.toString());

      const response = await importUsers(formData);
      if (response.code === 200) {
        message.success(`导入完成！成功：${response.data.success_count}，失败：${response.data.error_count}`);
        if (response.data.error_count > 0) {
          console.log('导入错误详情：', response.data.errors);
        }
        actionRef.current?.reload();
        setImportModalOpen(false);
      } else {
        message.error(response.message || '导入失败');
      }
    } catch (error) {
      message.error('导入失败');
    } finally {
      setImportLoading(false);
    }
  };

  // 表格列定义
  const columns: ProColumns<UserItem>[] = [
    {
      title: '用户ID',
      dataIndex: 'user_id',
      width: 150,
      ellipsis: true,
      search: false,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      width: 120,
      ellipsis: true,
    },
    {
      title: '显示名称',
      dataIndex: 'display_name',
      width: 120,
      ellipsis: true,
    },
    {
      title: '所属公司',
      dataIndex: 'company_id',
      width: 150,
      ellipsis: true,
      valueType: 'select',
      request: async () => {
        // TODO: 获取公司列表
        return [];
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      valueEnum: {
        active: {
          text: '激活',
          status: 'Success',
        },
        inactive: {
          text: '禁用',
          status: 'Error',
        },
        locked: {
          text: '锁定',
          status: 'Warning',
        },
      },
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      width: 200,
      ellipsis: true,
      search: false,
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      width: 150,
      ellipsis: true,
      search: false,
    },
    {
      title: '最后登录',
      dataIndex: 'last_login',
      width: 180,
      valueType: 'dateTime',
      search: false,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 180,
      valueType: 'dateTime',
      search: false,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      width: 200,
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setCurrentRow(record);
            setUpdateModalOpen(true);
          }}
        >
          编辑
        </a>,
        <a
          key="delete"
          onClick={() => handleDelete(record)}
          style={{ color: '#ff4d4f' }}
        >
          删除
        </a>,
        <a
          key="disable"
          onClick={() => handleQuickDisable(record)}
        >
          {record.status === 'active' ? '停用' : '启用'}
        </a>,
      ],
    },
  ];

  // 删除用户
  const handleDelete = (record: UserItem) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除用户 "${record.display_name}" 吗？`,
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        try {
          await deleteUser(record.user_id);
          message.success('删除成功');
          actionRef.current?.reload();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  // 快捷停用/启用
  const handleQuickDisable = async (record: UserItem) => {
    try {
      const newStatus = record.status === 'active' ? 'inactive' : 'active';
      await batchUpdateUserStatus([record.user_id], newStatus);
      message.success(`${newStatus === 'active' ? '启用' : '停用'}成功`);
      actionRef.current?.reload();
    } catch (error) {
      message.error('操作失败');
    }
  };

  // 批量操作菜单
  const batchMenuItems: MenuProps['items'] = [
    {
      key: 'active',
      label: '批量启用',
      onClick: () => handleBatchStatus('active'),
    },
    {
      key: 'inactive',
      label: '批量停用',
      onClick: () => handleBatchStatus('inactive'),
    },
    {
      key: 'delete',
      label: '批量删除',
      danger: true,
      onClick: () => handleBatchDelete(),
    },
  ];

  // 批量更新状态
  const handleBatchStatus = async (status: string) => {
    if (selectedRowsState.length === 0) {
      message.warning('请选择要操作的用户');
      return;
    }

    try {
      const userIds = selectedRowsState.map(row => row.user_id);
      await batchUpdateUserStatus(userIds, status);
      message.success(`批量${status === 'active' ? '启用' : '停用'}成功`);
      setSelectedRows([]);
      actionRef.current?.reload();
    } catch (error) {
      message.error('批量操作失败');
    }
  };

  // 批量删除
  const handleBatchDelete = () => {
    if (selectedRowsState.length === 0) {
      message.warning('请选择要删除的用户');
      return;
    }

    Modal.confirm({
      title: '确认批量删除',
      content: `确定要删除选中的 ${selectedRowsState.length} 个用户吗？`,
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        try {
          await Promise.all(selectedRowsState.map(row => deleteUser(row.user_id)));
          message.success('批量删除成功');
          setSelectedRows([]);
          actionRef.current?.reload();
        } catch (error) {
          message.error('批量删除失败');
        }
      },
    });
  };

  // 导出用户
  const handleExport = async () => {
    try {
      const blob = await exportUsers({});
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `用户列表_${new Date().toISOString().slice(0, 10)}.xlsx`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
      message.success('导出成功');
    } catch (error) {
      message.error('导出失败');
    }
  };

  return (
    <PageContainer>
      <ProTable<UserItem, UserListParams>
        headerTitle="用户列表"
        actionRef={actionRef}
        rowKey="user_id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button
            type="primary"
            key="primary"
            onClick={() => {
              setCreateModalOpen(true);
            }}
          >
            <PlusOutlined /> 新建用户
          </Button>,
          <Button
            key="import"
            onClick={() => setImportModalOpen(true)}
          >
            导入用户
          </Button>,
          <Button
            key="template"
            onClick={handleDownloadTemplate}
          >
            下载模板
          </Button>,
          <Button
            key="export"
            onClick={handleExport}
          >
            <ExportOutlined /> 导出
          </Button>,
        ]}
        request={async (params, sort, filter) => {
          try {
            const response = await getUserList({
              page: params.current || 1,
              page_size: params.pageSize || 20,
              username: params.username,
              display_name: params.display_name,
              company_id: params.company_id,
              status: params.status,
            });

            return {
              data: response.data?.users || [],
              success: response.code === 200,
              total: response.data?.total || 0,
            };
          } catch (error) {
            return {
              data: [],
              success: false,
              total: 0,
            };
          }
        }}
        columns={columns}
        rowSelection={{
          onChange: (_, selectedRows) => {
            setSelectedRows(selectedRows);
          },
        }}
        onChange={(pagination, filters, sorter) => {
          // 处理分页变化
          console.log('分页变化:', pagination);
          setCurrentPage(pagination.current || 1);
          setPageSize(pagination.pageSize || 20);
        }}
        pagination={{
          current: currentPage,
          pageSize: pageSize,
          showQuickJumper: true,
          showSizeChanger: true,
          pageSizeOptions: [10, 20, 50, 100],
          showTotal: (total, range) =>
            `第 ${range[0]}-${range[1]} 条/总共 ${total} 条`,
        }}
      />
      
      {selectedRowsState?.length > 0 && (
        <div
          style={{
            position: 'fixed',
            bottom: 24,
            left: '50%',
            transform: 'translateX(-50%)',
            zIndex: 1000,
            background: '#fff',
            padding: '12px 24px',
            borderRadius: 8,
            boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
          }}
        >
          <Space>
            <span>已选择 {selectedRowsState.length} 项</span>
            <Dropdown menu={{ items: batchMenuItems }} placement="topLeft">
              <Button>
                批量操作 <DownOutlined />
              </Button>
            </Dropdown>
            <Button onClick={() => setSelectedRows([])}>取消选择</Button>
          </Space>
        </div>
      )}

      {/* 新建用户模态框 */}
      <UserForm
        open={createModalOpen}
        onCancel={() => {
          setCreateModalOpen(false);
          setCurrentRow(undefined);
        }}
        onFinish={async (value: any) => {
          setCreateModalOpen(false);
          setCurrentRow(undefined);
          actionRef.current?.reload();
        }}
      />

      {/* 编辑用户模态框 */}
      <UserForm
        open={updateModalOpen}
        onCancel={() => {
          setUpdateModalOpen(false);
          setCurrentRow(undefined);
        }}
        onFinish={async (value: any) => {
          setUpdateModalOpen(false);
          setCurrentRow(undefined);
          actionRef.current?.reload();
        }}
        initialValues={currentRow}
      />

      {/* 导入用户模态框 */}
      <UserImportModal
        open={importModalOpen}
        onCancel={() => setImportModalOpen(false)}
        onImport={handleImport}
        loading={importLoading}
      />
    </PageContainer>
  );
};

export default UserManagement; 