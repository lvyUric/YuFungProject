import React, { useRef, useState } from 'react';
import {
  PageContainer,
  ProTable,
  ModalForm,
  ProForm,
  ProFormText,
  ProFormSelect,
  ProFormDigit,
  ProFormTextArea,
} from '@ant-design/pro-components';
import {
  Button,
  message,
  Popconfirm,
  Tag,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
} from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';

import {
  getSystemConfigList,
  createSystemConfig,
  getSystemConfig,
  updateSystemConfig,
  deleteSystemConfig,
  configTypeOptions,
  statusOptions,
  type SystemConfigInfo,
  type SystemConfigListParams,
  type SystemConfigCreateRequest,
  type SystemConfigUpdateRequest,
} from '@/services/system-config';

const SystemConfigManagement: React.FC = () => {
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [updateModalVisible, setUpdateModalVisible] = useState(false);
  const [currentRow, setCurrentRow] = useState<SystemConfigInfo>();
  const actionRef = useRef<ActionType>(null);

  // 创建系统配置
  const handleAdd = async (fields: SystemConfigCreateRequest) => {
    const hide = message.loading('正在创建...');
    try {
      const response = await createSystemConfig(fields);
      hide();
      if (response.code === 200) {
        message.success('创建成功');
        setCreateModalVisible(false);
        actionRef.current?.reload();
        return true;
      } else {
        message.error(response.message || '创建失败');
        return false;
      }
    } catch (error) {
      hide();
      message.error('创建失败');
      return false;
    }
  };

  // 更新系统配置
  const handleUpdate = async (fields: SystemConfigUpdateRequest) => {
    const hide = message.loading('正在更新...');
    try {
      const response = await updateSystemConfig(currentRow!.config_id, fields);
      hide();
      if (response.code === 200) {
        message.success('更新成功');
        setUpdateModalVisible(false);
        setCurrentRow(undefined);
        actionRef.current?.reload();
        return true;
      } else {
        message.error(response.message || '更新失败');
        return false;
      }
    } catch (error) {
      hide();
      message.error('更新失败');
      return false;
    }
  };

  // 删除系统配置
  const handleDelete = async (record: SystemConfigInfo) => {
    const hide = message.loading('正在删除...');
    try {
      const response = await deleteSystemConfig(record.config_id);
      hide();
      if (response.code === 200) {
        message.success('删除成功');
        actionRef.current?.reload();
      } else {
        message.error(response.message || '删除失败');
      }
    } catch (error) {
      hide();
      message.error('删除失败');
    }
  };

  // 表格列定义
  const columns: ProColumns<SystemConfigInfo>[] = [
    {
      title: '配置类型',
      dataIndex: 'config_type',
      width: 120,
      valueType: 'select',
      valueEnum: configTypeOptions.reduce((acc, item) => {
        acc[item.value] = { text: item.label };
        return acc;
      }, {} as Record<string, { text: string }>),
    },
    {
      title: '配置键',
      dataIndex: 'config_key',
      width: 150,
      copyable: true,
    },
    {
      title: '配置值',
      dataIndex: 'config_value',
      width: 150,
      ellipsis: true,
    },
    {
      title: '显示名称',
      dataIndex: 'display_name',
      width: 150,
      ellipsis: true,
    },
    {
      title: '排序',
      dataIndex: 'sort_order',
      width: 80,
      search: false,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 100,
      valueType: 'select',
      valueEnum: statusOptions.reduce((acc, item) => {
        acc[item.value] = { text: item.label };
        return acc;
      }, {} as Record<string, { text: string }>),
      render: (_, record) => (
        <Tag color={record.status === 'enable' ? 'green' : 'red'}>
          {record.status_text}
        </Tag>
      ),
    },
    {
      title: '备注',
      dataIndex: 'remark',
      width: 200,
      search: false,
      ellipsis: true,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 150,
      search: false,
      valueType: 'dateTime',
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      width: 150,
      fixed: 'right',
      render: (_, record) => [
        <Button
          key="edit"
          type="link"
          size="small"
          icon={<EditOutlined />}
          onClick={() => {
            setCurrentRow(record);
            setUpdateModalVisible(true);
          }}
        >
          编辑
        </Button>,
        <Popconfirm
          key="delete"
          title="确定要删除这条记录吗？"
          onConfirm={() => handleDelete(record)}
        >
          <Button
            type="link"
            size="small"
            danger
            icon={<DeleteOutlined />}
          >
            删除
          </Button>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer
      header={{
        title: '系统配置管理',
        breadcrumb: {
          items: [
            { path: '/system', title: '系统管理' },
            { title: '系统配置' },
          ],
        },
      }}
    >
      <ProTable<SystemConfigInfo, SystemConfigListParams>
        headerTitle="配置列表"
        actionRef={actionRef}
        rowKey="config_id"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button
            type="primary"
            key="primary"
            icon={<PlusOutlined />}
            onClick={() => setCreateModalVisible(true)}
          >
            新建配置
          </Button>,
        ]}
        request={async (params) => {
          try {
            const response = await getSystemConfigList({
              ...params,
              page: params.current,
              page_size: params.pageSize,
            });
            return {
              data: response.data?.list || [],
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
        scroll={{ x: 1200 }}
        pagination={{
          pageSize: 20,
          showQuickJumper: true,
          showSizeChanger: true,
          showTotal: (total, range) =>
            `第 ${range[0]}-${range[1]} 条/总共 ${total} 条`,
        }}
      />

      {/* 创建配置模态框 */}
      <ModalForm
        title="新建系统配置"
        width={600}
        open={createModalVisible}
        onOpenChange={setCreateModalVisible}
        onFinish={handleAdd}
        modalProps={{
          destroyOnClose: true,
        }}
      >
        <ProForm.Group>
          <ProFormSelect
            width="md"
            name="config_type"
            label="配置类型"
            options={configTypeOptions}
            rules={[{ required: true, message: '请选择配置类型' }]}
          />
          <ProFormText
            width="md"
            name="config_key"
            label="配置键"
            rules={[{ required: true, message: '请输入配置键' }]}
          />
        </ProForm.Group>
        <ProForm.Group>
          <ProFormText
            width="md"
            name="config_value"
            label="配置值"
            rules={[{ required: true, message: '请输入配置值' }]}
          />
          <ProFormText
            width="md"
            name="display_name"
            label="显示名称"
            rules={[{ required: true, message: '请输入显示名称' }]}
          />
        </ProForm.Group>
        <ProForm.Group>
          <ProFormDigit
            width="md"
            name="sort_order"
            label="排序"
            min={0}
            initialValue={0}
          />
          <ProFormSelect
            width="md"
            name="status"
            label="状态"
            options={statusOptions}
            initialValue="enable"
            rules={[{ required: true, message: '请选择状态' }]}
          />
        </ProForm.Group>
        <ProFormTextArea
          name="remark"
          label="备注说明"
          placeholder="请输入备注信息"
        />
      </ModalForm>

      {/* 编辑配置模态框 */}
      <ModalForm
        title="编辑系统配置"
        width={600}
        open={updateModalVisible}
        onOpenChange={setUpdateModalVisible}
        onFinish={handleUpdate}
        initialValues={currentRow}
        modalProps={{
          destroyOnClose: true,
        }}
      >
        <ProForm.Group>
          <ProFormText
            width="md"
            name="config_type"
            label="配置类型"
            disabled
          />
          <ProFormText
            width="md"
            name="config_key"
            label="配置键"
            disabled
          />
        </ProForm.Group>
        <ProForm.Group>
          <ProFormText
            width="md"
            name="config_value"
            label="配置值"
            rules={[{ required: true, message: '请输入配置值' }]}
          />
          <ProFormText
            width="md"
            name="display_name"
            label="显示名称"
            rules={[{ required: true, message: '请输入显示名称' }]}
          />
        </ProForm.Group>
        <ProForm.Group>
          <ProFormDigit
            width="md"
            name="sort_order"
            label="排序"
            min={0}
          />
          <ProFormSelect
            width="md"
            name="status"
            label="状态"
            options={statusOptions}
            rules={[{ required: true, message: '请选择状态' }]}
          />
        </ProForm.Group>
        <ProFormTextArea
          name="remark"
          label="备注说明"
          placeholder="请输入备注信息"
        />
      </ModalForm>
    </PageContainer>
  );
};

export default SystemConfigManagement; 