import React, { useState, useRef } from 'react';
import {
  PageContainer,
  ProTable,
} from '@ant-design/pro-components';
import { Button, Space, Tag, Tooltip, message, Modal } from 'antd';
import { DeleteOutlined, ReloadOutlined } from '@ant-design/icons';
import { getActivityLogs, deleteActivityLogsByCompany, ActivityLog, ActivityLogQuery, OPERATION_TYPE_LABELS, RESULT_STATUS_LABELS } from '@/services/activityLog';
import { formatDateTime } from '@/utils/date';
import type { ActionType, ProColumns } from '@ant-design/pro-components';

const ActivityLogPage: React.FC = () => {
  const actionRef = useRef<ActionType>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);

  // 删除公司活动记录
  const handleDeleteCompanyLogs = async (companyId: string) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除该公司(${companyId})的所有活动记录吗？此操作不可恢复。`,
      onOk: async () => {
        try {
          await deleteActivityLogsByCompany(companyId);
          message.success('删除成功');
          actionRef.current?.reload();
        } catch (error) {
          message.error('删除失败');
          console.error('删除失败:', error);
        }
      },
    });
  };

  // 表格列定义
  const columns: ProColumns<ActivityLog>[] = [
    {
      title: '用户',
      dataIndex: 'username',
      width: 120,
      fixed: 'left',
      render: (_, record) => (
        <div>
          <div>{record.username}</div>
          <div style={{ fontSize: '12px', color: '#999' }}>{record.user_id}</div>
        </div>
      ),
    },
    {
      title: '公司',
      dataIndex: 'company_name',
      width: 150,
      render: (_, record) => (
        <div>
          <div>{record.company_name}</div>
          <div style={{ fontSize: '12px', color: '#999' }}>{record.company_id}</div>
        </div>
      ),
    },
    {
      title: '操作类型',
      dataIndex: 'operation_type',
      width: 100,
      valueType: 'select',
      valueEnum: {
        login: { text: '登录', status: 'Processing' },
        logout: { text: '退出', status: 'Default' },
        create: { text: '创建', status: 'Success' },
        update: { text: '更新', status: 'Warning' },
        delete: { text: '删除', status: 'Error' },
        view: { text: '查看', status: 'Processing' },
        export: { text: '导出', status: 'Processing' },
        import: { text: '导入', status: 'Processing' },
      },
      render: (_, record) => {
        const label = OPERATION_TYPE_LABELS[record.operation_type as keyof typeof OPERATION_TYPE_LABELS] || record.operation_type;
        const color = record.operation_type === 'create' ? 'green' : 
                     record.operation_type === 'delete' ? 'red' : 
                     record.operation_type === 'update' ? 'orange' : 'blue';
        return <Tag color={color}>{label}</Tag>;
      },
    },
    {
      title: '模块',
      dataIndex: 'module_name',
      width: 120,
      valueType: 'select',
      valueEnum: {
        '用户管理': { text: '用户管理' },
        '角色管理': { text: '角色管理' },
        '菜单管理': { text: '菜单管理' },
        '公司管理': { text: '公司管理' },
        '保单管理': { text: '保单管理' },
        '客户管理': { text: '客户管理' },
        '系统管理': { text: '系统管理' },
        '认证授权': { text: '认证授权' },
      },
    },
    {
      title: '操作描述',
      dataIndex: 'operation_desc',
      ellipsis: true,
      width: 200,
      render: (text) => (
        <Tooltip title={text}>
          <span>{text}</span>
        </Tooltip>
      ),
    },
    {
      title: '目标',
      dataIndex: 'target_name',
      width: 120,
      render: (_, record) => (
        <div>
          {record.target_name && (
            <>
              <div>{record.target_name}</div>
              {record.target_id && (
                <div style={{ fontSize: '12px', color: '#999' }}>{record.target_id}</div>
              )}
            </>
          )}
        </div>
      ),
    },
    {
      title: 'IP地址',
      dataIndex: 'ip_address',
      width: 120,
      copyable: true,
    },
    {
      title: '执行时间',
      dataIndex: 'execution_time',
      width: 100,
      render: (text) => `${text}ms`,
    },
    {
      title: '状态',
      dataIndex: 'result_status',
      width: 80,
      valueType: 'select',
      valueEnum: {
        success: { text: '成功', status: 'Success' },
        failure: { text: '失败', status: 'Error' },
      },
      render: (_, record) => {
        const label = RESULT_STATUS_LABELS[record.result_status as keyof typeof RESULT_STATUS_LABELS] || record.result_status;
        const color = record.result_status === 'success' ? 'green' : 'red';
        return <Tag color={color}>{label}</Tag>;
      },
    },
    {
      title: '操作时间',
      dataIndex: 'operation_time',
      width: 160,
      valueType: 'dateTimeRange',
      render: (text: any) => text ? formatDateTime(text) : '-',
      sorter: true,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      width: 100,
      fixed: 'right',
      render: (_, record) => [
        <Tooltip key="delete" title="删除该公司所有记录">
          <Button
            type="link"
            size="small"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDeleteCompanyLogs(record.company_id)}
          >
            删除
          </Button>
        </Tooltip>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<ActivityLog, ActivityLogQuery>
        headerTitle="系统活动记录"
        actionRef={actionRef}
        rowKey={(record) => record.log_id || record.id || record.user_id + record.operation_time}
        search={{
          labelWidth: 120,
          defaultCollapsed: false,
        }}
        toolBarRender={() => [
          <Button
            key="refresh"
            icon={<ReloadOutlined />}
            onClick={() => actionRef.current?.reload()}
          >
            刷新
          </Button>,
        ]}
        request={async (params: any, sort: any, filter: any) => {
          const queryParams: ActivityLogQuery = {
            page: params.current || 1,
            page_size: params.pageSize || 20,
            user_id: params.user_id,
            operation_type: params.operation_type,
            module_name: params.module_name,
            result_status: params.result_status,
          };

          // 处理时间范围
          if (params.operation_time && Array.isArray(params.operation_time) && params.operation_time.length === 2) {
            queryParams.start_time = params.operation_time[0] as string;
            queryParams.end_time = params.operation_time[1] as string;
          }

          // 处理排序
          if (sort && Object.keys(sort).length > 0) {
            const sortKey = Object.keys(sort)[0];
            const sortOrder = sort[sortKey];
            queryParams.sort_by = sortKey;
            queryParams.sort_order = sortOrder === 'ascend' ? 'asc' : 'desc';
          }

          try {
            const response = await getActivityLogs(queryParams);
            // 检查响应结构
            if (response.code === 200 && response.data) {
              return {
                data: response.data.list || [],
                success: true,
                total: response.data.total || 0,
              };
            } else {
              message.error(response.message || '获取活动记录失败');
              return {
                data: [],
                success: false,
                total: 0,
              };
            }
          } catch (error) {
            message.error('获取活动记录失败');
            return {
              data: [],
              success: false,
              total: 0,
            };
          }
        }}
        columns={columns}
        onChange={(pagination, filters, sorter) => {
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
        scroll={{ x: 1400 }}
        options={{
          density: true,
          fullScreen: true,
          reload: true,
          setting: true,
        }}
        columnsState={{
          persistenceKey: 'activity-log-table',
          persistenceType: 'localStorage',
        }}
      />
    </PageContainer>
  );
};

export default ActivityLogPage; 