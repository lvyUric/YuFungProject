import React, { useRef, useState, useEffect } from 'react';
import type { ActionType, ProColumns, ProDescriptionsItemProps } from '@ant-design/pro-components';
import {
  FooterToolbar,
  PageContainer,
  ProDescriptions,
  ProTable,
} from '@ant-design/pro-components';
import { FormattedMessage, useIntl } from '@umijs/max';
import {
  Button,
  Drawer,
  message,
  Modal,
  Tag,
  Dropdown,
  Space,
  Tooltip,
} from 'antd';
import {
  PlusOutlined,
  ExclamationCircleOutlined,
  ImportOutlined,
  ExportOutlined,
  DownOutlined,
  EyeOutlined,
  EyeInvisibleOutlined,
} from '@ant-design/icons';
import {
  getCompanyList,
  deleteCompany,
} from '@/services/ant-design-pro/company';
import CreateForm from './components/CreateForm';
import UpdateForm from './components/UpdateForm';
import ImportModal from './components/ImportModal';
import ExportModal from './components/ExportModal';

const { confirm } = Modal;

const CompanyList: React.FC = () => {
  const actionRef = useRef<ActionType>(null);
  const [showDetail, setShowDetail] = useState<boolean>(false);
  const [currentRow, setCurrentRow] = useState<API.CompanyInfo>();
  const [selectedRowsState, setSelectedRows] = useState<API.CompanyInfo[]>([]);
  const [createModalOpen, setCreateModalOpen] = useState<boolean>(false);
  const [updateModalOpen, setUpdateModalOpen] = useState<boolean>(false);
  const [importModalOpen, setImportModalOpen] = useState<boolean>(false);
  const [exportModalOpen, setExportModalOpen] = useState<boolean>(false);
  const [currentFilters, setCurrentFilters] = useState<{ status?: string; keyword?: string }>({});

  // 响应式断点状态
  const [isMobile, setIsMobile] = useState(false);
  const [isTablet, setIsTablet] = useState(false);

  useEffect(() => {
    const handleResize = () => {
      const width = window.innerWidth;
      setIsMobile(width < 768);
      setIsTablet(width >= 768 && width < 1200);
      
      // 根据屏幕大小自动调整列显示
      if (width < 768) {
        // 移动端：只显示核心列
        // Mobile responsive handling can be done through ProTable's built-in column settings
      } else if (width < 1200) {
        // 平板端：显示重要列  
        // Tablet responsive handling can be done through ProTable's built-in column settings
      }
      // 桌面端保持用户自定义设置
    };

    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  const intl = useIntl();

  /**
   * 删除节点
   */
  const handleRemove = async (selectedRows: API.CompanyInfo[]) => {
    const hide = message.loading(
      intl.formatMessage({
        id: 'pages.companyList.processing',
        defaultMessage: '正在删除...',
      }),
    );
    if (!selectedRows) return true;
    try {
      await Promise.all(
        selectedRows.map((row) => {
          if (row.id) {
            return deleteCompany(row.id);
          }
          return Promise.resolve();
        })
      );
      hide();
      message.success(
        intl.formatMessage({
          id: 'pages.companyList.deleteSuccess',
          defaultMessage: '删除成功，即将刷新',
        }),
      );
      return true;
    } catch (error) {
      hide();
      message.error(
        intl.formatMessage({
          id: 'pages.companyList.deleteError',
          defaultMessage: '删除失败，请重试',
        }),
      );
      return false;
    }
  };

  /**
   * 删除确认
   */
  const showDeleteConfirm = (record: API.CompanyInfo) => {
    confirm({
      title: intl.formatMessage({
        id: 'pages.companyList.deleteConfirm.title',
        defaultMessage: '确认删除',
      }),
      icon: <ExclamationCircleOutlined />,
      content: intl.formatMessage({
        id: 'pages.companyList.deleteConfirm.content',
        defaultMessage: '确定要删除这个公司吗？删除后不可恢复。',
      }),
      okText: intl.formatMessage({
        id: 'pages.companyList.deleteConfirm.ok',
        defaultMessage: '确定',
      }),
      cancelText: intl.formatMessage({
        id: 'pages.companyList.deleteConfirm.cancel',
        defaultMessage: '取消',
      }),
      onOk: async () => {
        const success = await handleRemove([record]);
        if (success) {
          actionRef.current?.reloadAndRest?.();
        }
      },
    });
  };

  /**
   * 批量删除确认
   */
  const showBatchDeleteConfirm = () => {
    confirm({
      title: intl.formatMessage({
        id: 'pages.companyList.batchDeleteConfirm.title',
        defaultMessage: '确认批量删除',
      }),
      icon: <ExclamationCircleOutlined />,
      content: intl.formatMessage(
        {
          id: 'pages.companyList.batchDeleteConfirm.content',
          defaultMessage: '确定要删除选中的 {count} 个公司吗？删除后不可恢复。',
        },
        { count: selectedRowsState.length },
      ),
      okText: intl.formatMessage({
        id: 'pages.companyList.batchDeleteConfirm.ok',
        defaultMessage: '确定',
      }),
      cancelText: intl.formatMessage({
        id: 'pages.companyList.batchDeleteConfirm.cancel',
        defaultMessage: '取消',
      }),
      onOk: async () => {
        const success = await handleRemove(selectedRowsState);
        if (success) {
          setSelectedRows([]);
          actionRef.current?.reloadAndRest?.();
        }
      },
    });
  };

  // 导入导出菜单
  const importExportMenuItems = [
    {
      key: 'import',
      label: '导入数据',
      icon: <ImportOutlined />,
      onClick: () => setImportModalOpen(true),
    },
    {
      key: 'export',
      label: '导出数据',
      icon: <ExportOutlined />,
      onClick: () => setExportModalOpen(true),
    },
  ];

  const columns: ProColumns<API.CompanyInfo>[] = [
    {
      title: <FormattedMessage id="pages.companyList.companyName" defaultMessage="公司名称" />,
      dataIndex: 'company_name',
      key: 'company_name',
      fixed: 'left',
      width: isMobile ? 120 : 150,
      ellipsis: true,
      tip: intl.formatMessage({
        id: 'pages.companyList.companyName.tip',
        defaultMessage: '公司名称是唯一的 key',
      }),
      render: (dom, entity) => {
        return (
          <Tooltip title={dom}>
            <a
              onClick={() => {
                setCurrentRow(entity);
                setShowDetail(true);
              }}
              style={{ 
                display: 'block',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap'
              }}
            >
              {dom}
            </a>
          </Tooltip>
        );
      },
    },
    {
      title: <FormattedMessage id="pages.companyList.companyCode" defaultMessage="公司代码" />,
      dataIndex: 'company_code',
      sorter: true,
      hideInForm: true,
      width: isMobile ? 80 : 120,
      ellipsis: true,
    },
    {
      title: <FormattedMessage id="pages.companyList.contactPerson" defaultMessage="联络人" />,
      dataIndex: 'contact_person',
      valueType: 'text',
      hideInForm: true,
      width: isMobile ? 70 : 100,
      ellipsis: true,
    },
    {
      title: <FormattedMessage id="pages.companyList.mobile" defaultMessage="联系电话" />,
      dataIndex: 'mobile',
      hideInForm: true,
      width: isMobile ? 100 : 130,
      ellipsis: true,
      render: (_, record) => {
        const phone = record.mobile || record.tel_no || record.contact_phone;
        return phone || '-';
      },
    },
    {
      title: <FormattedMessage id="pages.companyList.email" defaultMessage="邮箱" />,
      dataIndex: 'email',
      valueType: 'text',
      hideInForm: true,
      width: isMobile ? 120 : 180,
      ellipsis: true,
      render: (_, record) => {
        return record.email ? (
          <a href={`mailto:${record.email}`}>{record.email}</a>
        ) : '-';
      },
    },
    {
      title: '中文地址',
      dataIndex: 'address_cn_detail',
      valueType: 'text',
      hideInForm: true,
      hideInSearch: true,
      ellipsis: true,
      width: 200,
      render: (_, record) => {
        const fullAddress = record.address_cn_detail || record.address;
        
        return fullAddress ? (
          <Tooltip title={fullAddress}>
            <span>{fullAddress}</span>
          </Tooltip>
        ) : '-';
      },
    },
    {
      title: 'Broker Code',
      dataIndex: 'broker_code',
      hideInForm: true,
      hideInSearch: true,
      width: isMobile ? 80 : 120,
      ellipsis: true,
    },
    {
      title: <FormattedMessage id="pages.companyList.validPeriod" defaultMessage="有效期" />,
      dataIndex: 'valid_start_date',
      valueType: 'dateRange',
      width: isMobile ? 140 : 180,
      render: (_, record) => {
        const startDate = record.valid_start_date;
        const endDate = record.valid_end_date;
        
        if (!startDate || !endDate) return '-';
        
        const start = new Date(startDate).toLocaleDateString('zh-CN');
        const end = new Date(endDate).toLocaleDateString('zh-CN');
        
        return `${start} ~ ${end}`;
      },
      hideInForm: true,
      hideInSearch: true,
    },
    {
      title: <FormattedMessage id="pages.companyList.userQuota" defaultMessage="用户配额" />,
      dataIndex: 'user_quota',
      valueType: 'digit',
      width: isMobile ? 80 : 100,
      render: (_, record) => (
        <span>
          {record.current_user_count || 0} / {record.user_quota || 0}
        </span>
      ),
      hideInForm: true,
      hideInSearch: true,
    },
    {
      title: <FormattedMessage id="pages.companyList.status" defaultMessage="状态" />,
      dataIndex: 'status',
      hideInForm: true,
      width: isMobile ? 60 : 80,
      valueEnum: {
        active: {
          text: <FormattedMessage id="pages.companyList.status.active" defaultMessage="有效" />,
          status: 'Success',
        },
        inactive: {
          text: <FormattedMessage id="pages.companyList.status.inactive" defaultMessage="停用" />,
          status: 'Default',
        },
        expired: {
          text: <FormattedMessage id="pages.companyList.status.expired" defaultMessage="过期" />,
          status: 'Error',
        },
      },
    },
    {
      title: '提交人',
      dataIndex: 'submitted_by',
      hideInForm: true,
      width: isMobile ? 70 : 100,
      ellipsis: true,
      search: false,
    },
    {
      title: <FormattedMessage id="pages.companyList.createdAt" defaultMessage="创建时间" />,
      dataIndex: 'created_at',
      valueType: 'dateTime',
      width: isMobile ? 120 : 150,
      search: false,
      render: (_, record) => {
        return record.created_at ? new Date(record.created_at).toLocaleString('zh-CN') : '-';
      },
      hideInForm: true,
    },
    {
      title: <FormattedMessage id="pages.companyList.option" defaultMessage="操作" />,
      dataIndex: 'option',
      valueType: 'option',
      fixed: 'right',
      width: isMobile ? 80 : 120,
      render: (_, record) => [
        <a
          key="edit"
          onClick={() => {
            setCurrentRow(record);
            setUpdateModalOpen(true);
          }}
        >
          <FormattedMessage id="pages.companyList.edit" defaultMessage="编辑" />
        </a>,
        <a
          key="delete"
          style={{ color: 'red' }}
          onClick={() => showDeleteConfirm(record)}
        >
          <FormattedMessage id="pages.companyList.delete" defaultMessage="删除" />
        </a>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.CompanyInfo, API.CompanyQueryParams>
        headerTitle={intl.formatMessage({
          id: 'pages.companyList.title',
          defaultMessage: '公司列表',
        })}
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        // 添加表格滚动配置
        scroll={{ 
          x: 'max-content', // 根据内容自动调整宽度
          y: 'calc(100vh - 350px)' // 设置垂直滚动高度
        }}
        // 优化表格尺寸
        size={isMobile ? 'small' : 'middle'}
        // 添加列设置工具
        toolBarRender={() => [
          <Button
            type="primary"
            key="primary"
            onClick={() => {
              setCreateModalOpen(true);
            }}
          >
            <PlusOutlined /> <FormattedMessage id="pages.companyList.new" defaultMessage="新建" />
          </Button>,
          <Dropdown
            key="importExport"
            menu={{
              items: importExportMenuItems,
            }}
            placement="bottomLeft"
          >
            <Button>
              <Space>
                数据管理
                <DownOutlined />
              </Space>
            </Button>
          </Dropdown>,
        ]}
        request={async (params) => {
          // 保存当前筛选条件
          setCurrentFilters({
            status: params.status,
            keyword: params.keyword,
          });

          const response = await getCompanyList({
            page: params.current || 1,
            page_size: params.pageSize || 20,
            status: params.status,
            keyword: params.keyword,
          });

          if (response.code === 200 && response.data) {
            return {
              data: response.data.companies || [],
              success: true,
              total: response.data.total || 0,
            };
          }

          return {
            data: [],
            success: false,
            total: 0,
          };
        }}
        columns={columns}
        rowSelection={{
          onChange: (_, selectedRows) => {
            setSelectedRows(selectedRows);
          },
        }}
        pagination={{
          defaultPageSize: 20,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total, range) =>
            intl.formatMessage(
              {
                id: 'pages.companyList.pagination.total',
                defaultMessage: '共 {total} 条记录，第 {start}-{end} 条',
              },
              {
                total,
                start: range[0],
                end: range[1],
              },
            ),
        }}
      />
      {selectedRowsState?.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              <FormattedMessage
                id="pages.companyList.chosen"
                defaultMessage="已选择"
              />{' '}
              <a style={{ fontWeight: 600 }}>{selectedRowsState.length}</a>{' '}
              <FormattedMessage id="pages.companyList.item" defaultMessage="项" />
            </div>
          }
        >
          <Button
            type="primary"
            danger
            onClick={showBatchDeleteConfirm}
          >
            <FormattedMessage
              id="pages.companyList.batchDeletion"
              defaultMessage="批量删除"
            />
          </Button>
        </FooterToolbar>
      )}

      <CreateForm
        onOpenChange={setCreateModalOpen}
        modalOpen={createModalOpen}
        onFinish={async (success: boolean) => {
          if (success) {
            setCreateModalOpen(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      />

      {updateModalOpen && currentRow && (
        <UpdateForm
          onOpenChange={setUpdateModalOpen}
          modalOpen={updateModalOpen}
          values={currentRow}
          onFinish={async (success: boolean) => {
            if (success) {
              setUpdateModalOpen(false);
              setCurrentRow(undefined);
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
        />
      )}

      <ImportModal
        open={importModalOpen}
        onOpenChange={setImportModalOpen}
        onFinish={async (success: boolean) => {
          if (success) {
            setImportModalOpen(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      />

      <ExportModal
        open={exportModalOpen}
        onOpenChange={setExportModalOpen}
        selectedRows={selectedRowsState}
        currentFilters={currentFilters}
      />

      <Drawer
        width={600}
        open={showDetail}
        onClose={() => {
          setCurrentRow(undefined);
          setShowDetail(false);
        }}
        closable={false}
      >
        {currentRow?.company_name && (
          <ProDescriptions<API.CompanyInfo>
            column={2}
            title={currentRow?.company_name}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.company_name,
            }}
            columns={columns as ProDescriptionsItemProps<API.CompanyInfo>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default CompanyList; 