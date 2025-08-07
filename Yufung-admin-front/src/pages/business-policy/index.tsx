import React, { useState, useRef, useEffect } from 'react';
import {
  PageContainer,
  ProTable,
  ModalForm,
  ProForm,
  ProFormText,
  ProFormSelect,
  ProFormDigit,
  ProFormDatePicker,
  ProFormSwitch,
  ProFormTextArea,
} from '@ant-design/pro-components';
import {
  Button,
  message,
  Modal,
  Tag,
  Space,
  Popconfirm,
  Card,
  Row,
  Col,
  Select,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  ExportOutlined,
  ImportOutlined,
  DownloadOutlined,
  BarChartOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { useNavigate } from '@umijs/max';
import type { ActionType, ProColumns } from '@ant-design/pro-components';

import {
  getPolicyList,
  getPolicyDetail,
  createPolicy,
  updatePolicy,
  deletePolicy,
  getPolicyStatistics,
  batchUpdatePolicyStatus,
  getSystemConfigOptions,
  getCompanyOptions,
  type PolicyInfo,
  type PolicyListParams,
  type PolicyCreateRequest,
  type PolicyUpdateRequest,
  type PolicyStatistics,
  type SystemConfigOption,
} from '@/services/policy';

// 引入新的导入导出组件
import PolicyStepForm from './components/PolicyStepForm';
import ImportModal from './components/ImportModal';
import ExportModal from './components/ExportModal';

const PolicyManagement: React.FC = () => {
  const navigate = useNavigate();
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [stepFormVisible, setStepFormVisible] = useState(false);
  const [editStepFormVisible, setEditStepFormVisible] = useState(false);
  const [importModalVisible, setImportModalVisible] = useState(false);
  const [exportModalVisible, setExportModalVisible] = useState(false);
  const [selectedRows, setSelectedRows] = useState<PolicyInfo[]>([]);
  const [currentRow, setCurrentRow] = useState<PolicyInfo>();
  const [statistics, setStatistics] = useState<PolicyStatistics>();
  const [showStatistics, setShowStatistics] = useState(true);
  const [editingId, setEditingId] = useState<string>('');
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const actionRef = useRef<ActionType>(null);

  // 系统配置选项状态
  const [hkManagerOptions, setHkManagerOptions] = useState<SystemConfigOption[]>([]);
  const [referralBranchOptions, setReferralBranchOptions] = useState<SystemConfigOption[]>([]);
  const [partnerOptions, setPartnerOptions] = useState<SystemConfigOption[]>([]);
  const [companyOptions, setCompanyOptions] = useState<any[]>([]);

  // 加载系统配置选项
  useEffect(() => {
    const loadConfigOptions = async () => {
      try {
        // 加载港分客户经理选项
        const hkManagerRes = await getSystemConfigOptions('hk_manager');
        if (hkManagerRes.code === 200 && Array.isArray(hkManagerRes.data)) {
          setHkManagerOptions(hkManagerRes.data || []);
        } else {
          setHkManagerOptions([]);
        }

        // 加载转介分行选项
        const referralBranchRes = await getSystemConfigOptions('referral_branch');
        if (referralBranchRes.code === 200 && Array.isArray(referralBranchRes.data)) {
          setReferralBranchOptions(referralBranchRes.data || []);
        } else {
          setReferralBranchOptions([]);
        }

        // 加载合作伙伴选项
        const partnerRes = await getSystemConfigOptions('partner');
        if (partnerRes.code === 200 && Array.isArray(partnerRes.data)) {
          setPartnerOptions(partnerRes.data || []);
        } else {
          setPartnerOptions([]);
        }

        // 加载公司选项
        const companyRes = await getCompanyOptions();
        if (companyRes.code === 200 && Array.isArray(companyRes.data)) {
          setCompanyOptions(companyRes.data || []);
        } else {
          setCompanyOptions([]);
        }
      } catch (error) {
        console.error('加载配置选项失败:', error);
      }
    };

    loadConfigOptions();
  }, []);

  // 币种选项
  const currencyOptions = [
    { label: '美元 (USD)', value: 'USD' },
    { label: '港币 (HKD)', value: 'HKD' },
    { label: '人民币 (CNY)', value: 'CNY' },
  ];

  // 缴费方式选项
  const paymentMethodOptions = [
    { label: '期缴', value: '期缴' },
    { label: '趸缴', value: '趸缴' },
    { label: '预缴', value: '预缴' },
  ];

  // 加载统计数据
  useEffect(() => {
    const loadStatistics = async () => {
      try {
        const response = await getPolicyStatistics();
        if (response.code === 200) {
          setStatistics(response.data);
        }
      } catch (error) {
        console.error('加载统计数据失败:', error);
      }
    };

    loadStatistics();
  }, []);

  // 创建保单
  const handleAdd = async (fields: PolicyCreateRequest) => {
    const hide = message.loading('正在添加');
    try {
      await createPolicy(fields);
      hide();
      message.success('添加成功');
      return true;
    } catch (error) {
      hide();
      message.error('添加失败请重试！');
      return false;
    }
  };

  // 更新保单
  const handleUpdate = async (fields: PolicyUpdateRequest) => {
    const hide = message.loading('正在配置');
    try {
      await updatePolicy(currentRow?.policy_id || '', fields);
      hide();
      message.success('配置成功');
      return true;
    } catch (error) {
      hide();
      message.error('配置失败请重试！');
      return false;
    }
  };

  // 删除保单
  const handleRemove = async (selectedRows: PolicyInfo[]) => {
    const hide = message.loading('正在删除');
    if (!selectedRows) return true;
    try {
      await Promise.all(
        selectedRows.map((row) => deletePolicy(row.policy_id))
      );
      hide();
      message.success('删除成功，即将刷新');
      return true;
    } catch (error) {
      hide();
      message.error('删除失败，请重试');
      return false;
    }
  };

  // 批量删除
  const handleBatchDelete = async () => {
    if (selectedRows.length === 0) {
      message.warning('请选择要删除的记录');
      return;
    }

    Modal.confirm({
      title: '确认删除',
      content: `确定要删除选中的 ${selectedRows.length} 条记录吗？`,
      onOk: async () => {
        const success = await handleRemove(selectedRows);
        if (success) {
          setSelectedRows([]);
          actionRef.current?.reloadAndRest?.();
        }
      },
    });
  };

  // 批量更新状态
  const handleBatchUpdateStatus = async (field: string, value: boolean) => {
    if (selectedRows.length === 0) {
      message.warning('请选择要更新的记录');
      return;
    }

    Modal.confirm({
      title: '确认更新',
      content: `确定要更新选中的 ${selectedRows.length} 条记录的状态吗？`,
      onOk: async () => {
        try {
          const hide = message.loading('正在更新...');
          await batchUpdatePolicyStatus({
            policy_ids: selectedRows.map(row => row.policy_id),
            [field]: value,
          });
          hide();
          message.success('批量更新成功');
          setSelectedRows([]);
          actionRef.current?.reload();
        } catch (error) {
          message.error('批量更新失败');
        }
      },
    });
  };

  // 分步表单成功回调
  const handleStepFormSuccess = () => {
    setStepFormVisible(false);
    if (actionRef.current) {
      actionRef.current.reload();
    }
  };

  // 分步表单成功回调（编辑）
  const handleEditStepFormSuccess = () => {
    setEditStepFormVisible(false);
    setCurrentRow(undefined);
    if (actionRef.current) {
      actionRef.current.reload();
    }
  };

  // 编辑表单关闭回调
  const handleEditStepFormCancel = (visible: boolean) => {
    setEditStepFormVisible(visible);
    if (!visible) {
      setCurrentRow(undefined);
      setEditingId('');
    }
  };

  // 导出数据
  const handleExport = () => {
    setExportModalVisible(true);
  };

  // 批量导入
  const handleImport = () => {
    setImportModalVisible(true);
  };

  // 导入完成后的回调
  const handleImportFinish = (success: boolean) => {
    if (success && actionRef.current) {
      actionRef.current.reload();
    }
    // 无论成功还是失败，都关闭导入模态框
    setImportModalVisible(false);
  };

  // 编辑保单 - 先获取详情数据
  const handleEdit = async (record: PolicyInfo) => {
    try {
      setEditingId(record.policy_id);
      console.log('开始编辑保单，记录ID:', record.policy_id);
      
      const response = await getPolicyDetail(record.policy_id);
      console.log('获取保单详情响应:', response);
      
      if (response.code === 200 && response.data) {
        console.log('设置当前行数据:', response.data);
        setCurrentRow(response.data);
        setEditStepFormVisible(true);
      } else {
        message.error(response.message || '获取保单详情失败');
      }
    } catch (error) {
      console.error('获取保单详情失败:', error);
      message.error('获取保单详情失败，请重试');
    } finally {
      setEditingId('');
    }
  };

  // 查看详情
  const handleViewDetail = (record: PolicyInfo) => {
    navigate(`/business/policy/detail/${record.policy_id}`);
  };

  // 表格列定义 - 按照模板字段顺序
  const columns: ProColumns<PolicyInfo>[] = [
    {
      title: '序号',
      dataIndex: 'serial_number',
      hideInSearch: true,
      width: 60,
      fixed: 'left',
    },
    {
      title: '账户号',
      dataIndex: 'account_number',
      copyable: true,
      width: 120,
      fixed: 'left',
    },
    {
      title: '客户号',
      dataIndex: 'customer_number',
      copyable: true,
      width: 100,
    },
    {
      title: '客户中文名',
      dataIndex: 'customer_name_cn',
      width: 120,
    },
    {
      title: '客户英文名',
      dataIndex: 'customer_name_en',
      width: 120,
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: '投保单号',
      dataIndex: 'proposal_number',
      copyable: true,
      width: 150,
    },
    {
      title: '保单币种',
      dataIndex: 'policy_currency',
      valueType: 'select',
      valueEnum: {
        USD: { text: 'USD' },
        HKD: { text: 'HKD' },
        CNY: { text: 'CNY' },
      },
      width: 100,
    },
    {
      title: '合作伙伴',
      dataIndex: 'partner',
      valueType: 'select',
      request: async () => {
        return partnerOptions.map(option => ({
          label: option.display_name,
          value: option.config_value,
        }));
      },
      width: 120,
      ellipsis: true,
    },
    {
      title: '转介编号',
      dataIndex: 'referral_code',
      copyable: true,
      width: 120,
      hideInSearch: true,
    },
    {
      title: '港分客户经理',
      dataIndex: 'hk_manager',
      valueType: 'select',
      request: async () => {
        return hkManagerOptions.map(option => ({
          label: option.display_name,
          value: option.config_value,
        }));
      },
      width: 120,
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '转介理财经理',
      dataIndex: 'referral_pm',
      width: 120,
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '转介分行',
      dataIndex: 'referral_branch',
      valueType: 'select',
      request: async () => {
        return referralBranchOptions.map(option => ({
          label: option.display_name,
          value: option.config_value,
        }));
      },
      width: 120,
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '转介支行',
      dataIndex: 'referral_sub_branch',
      width: 120,
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '转介日期',
      dataIndex: 'referral_date',
      valueType: 'date',
      width: 120,
      hideInSearch: true,
    },
    {
      title: '是否退保',
      dataIndex: 'is_surrendered',
      valueType: 'select',
      valueEnum: {
        true: { text: '是', status: 'Error' },
        false: { text: '否', status: 'Success' },
      },
      render: (_, record) => (
        <Tag color={record.is_surrendered ? 'red' : 'green'}>
          {record.is_surrendered ? '是' : '否'}
        </Tag>
      ),
      width: 100,
    },
    {
      title: '缴费日期',
      dataIndex: 'payment_date',
      valueType: 'date',
      width: 120,
      hideInSearch: true,
    },
    {
      title: '生效日期',
      dataIndex: 'effective_date',
      valueType: 'date',
      width: 120,
      hideInSearch: true,
    },
    {
      title: '缴费方式',
      dataIndex: 'payment_method',
      valueType: 'select',
      valueEnum: {
        '期缴': { text: '期缴' },
        '趸缴': { text: '趸缴' },
        '预缴': { text: '预缴' },
      },
      width: 120,
      hideInSearch: true,
    },
    {
      title: '缴费年期',
      dataIndex: 'payment_years',
      width: 100,
      hideInSearch: true,
      render: (_, record) => record.payment_years ? `${record.payment_years}年` : '-',
    },
    {
      title: '期缴期数',
      dataIndex: 'payment_periods',
      width: 100,
      hideInSearch: true,
      render: (_, record) => record.payment_periods ? `${record.payment_periods}期` : '-',
    },
    {
      title: '实际缴纳保费',
      dataIndex: 'actual_premium',
      valueType: 'money',
      hideInSearch: true,
      width: 130,
      render: (_, record) => record.actual_premium?.toLocaleString() || '-',
    },
    {
      title: 'AUM',
      dataIndex: 'aum',
      valueType: 'money',
      hideInSearch: true,
      width: 120,
      render: (_, record) => record.aum?.toLocaleString() || '-',
    },
    {
      title: '是否已过冷静期',
      dataIndex: 'past_cooling_period',
      valueType: 'select',
      valueEnum: {
        true: { text: '是', status: 'Success' },
        false: { text: '否', status: 'Warning' },
      },
      render: (_, record) => (
        <Tag color={record.past_cooling_period ? 'green' : 'orange'}>
          {record.past_cooling_period ? '是' : '否'}
        </Tag>
      ),
      width: 130,
    },
    {
      title: '是否支付佣金',
      dataIndex: 'is_paid_commission',
      valueType: 'select',
      valueEnum: {
        true: { text: '是', status: 'Success' },
        false: { text: '否', status: 'Warning' },
      },
      render: (_, record) => (
        <Tag color={record.is_paid_commission ? 'green' : 'orange'}>
          {record.is_paid_commission ? '是' : '否'}
        </Tag>
      ),
      width: 120,
    },
    {
      title: '转介费率',
      dataIndex: 'referral_rate',
      width: 100,
      hideInSearch: true,
      render: (_, record) => record.referral_rate ? `${record.referral_rate}%` : '-',
    },
    {
      title: '汇率',
      dataIndex: 'exchange_rate',
      width: 100,
      hideInSearch: true,
      render: (_, record) => record.exchange_rate ? record.exchange_rate.toFixed(4) : '-',
    },
    {
      title: '预计转介费',
      dataIndex: 'expected_fee',
      valueType: 'money',
      width: 120,
      hideInSearch: true,
      render: (_, record) => record.expected_fee?.toLocaleString() || '-',
    },
    {
      title: '支付日期',
      dataIndex: 'payment_pay_date',
      valueType: 'date',
      width: 120,
      hideInSearch: true,
    },
    {
      title: '是否员工',
      dataIndex: 'is_employee',
      valueType: 'select',
      valueEnum: {
        true: { text: '是', status: 'Processing' },
        false: { text: '否', status: 'Default' },
      },
      render: (_, record) => (
        <Tag color={record.is_employee ? 'purple' : 'default'}>
          {record.is_employee ? '是' : '否'}
        </Tag>
      ),
      width: 100,
      hideInSearch: true,
    },
    {
      title: '承保公司',
      dataIndex: 'insurance_company',
      valueType: 'select',
      request: async () => {
        return companyOptions.map(option => ({
          label: option.company_name,
          value: option.company_name,
        }));
      },
      width: 120,
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '保险产品名称',
      dataIndex: 'product_name',
      width: 150,
      ellipsis: true,
    },
    {
      title: '产品类型',
      dataIndex: 'product_type',
      width: 120,
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '备注说明',
      dataIndex: 'remark',
      width: 150,
      ellipsis: true,
      hideInSearch: true,
    },
    {
      title: '操作',
      dataIndex: 'option',
      valueType: 'option',
      width: 260,
      fixed: 'right',
      render: (_, record) => [
        <Button
          key="detail"
          type="link"
          size="small"
          icon={<EyeOutlined />}
          onClick={() => handleViewDetail(record)}
        >
          详情
        </Button>,
        <Button
          key="editable"
          type="link"
          size="small"
          icon={<EditOutlined />}
          loading={editingId === record.policy_id}
          onClick={() => handleEdit(record)}
        >
          编辑
        </Button>,
        <Popconfirm
          key="delete"
          title="确认删除这条记录吗？"
          onConfirm={async () => {
            await handleRemove([record]);
            actionRef.current?.reload();
          }}
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
    <PageContainer>
      {/* 统计卡片 */}
      {showStatistics && statistics && (
        <div style={{ marginBottom: 16 }}>
          <div style={{ display: 'flex', gap: 16 }}>
            <Card size="small" style={{ flex: 1 }}>
              <div style={{ textAlign: 'center' }}>
                <div style={{ fontSize: 24, fontWeight: 'bold', color: '#1890ff' }}>
                  {statistics.total_policies || 0}
                </div>
                <div>总保单数</div>
              </div>
            </Card>
            <Card size="small" style={{ flex: 1 }}>
              <div style={{ textAlign: 'center' }}>
                <div style={{ fontSize: 24, fontWeight: 'bold', color: '#52c41a' }}>
                  {(statistics.total_premium / 10000).toFixed(2)}
                </div>
                <div>总保费(万)</div>
              </div>
            </Card>
            <Card size="small" style={{ flex: 1 }}>
              <div style={{ textAlign: 'center' }}>
                <div style={{ fontSize: 24, fontWeight: 'bold', color: '#f5222d' }}>
                  {(statistics.total_aum / 10000).toFixed(2)}
                </div>
                <div>总AUM(万)</div>
              </div>
            </Card>
            <Card size="small" style={{ flex: 1 }}>
              <div style={{ textAlign: 'center' }}>
                <div style={{ fontSize: 24, fontWeight: 'bold', color: '#722ed1' }}>
                  {(statistics.total_expected_fee / 10000).toFixed(2)}
                </div>
                <div>预计转介费(万)</div>
              </div>
            </Card>
            <Card size="small" style={{ flex: 1 }}>
              <div style={{ textAlign: 'center' }}>
                <div style={{ fontSize: 24, fontWeight: 'bold', color: '#fa8c16' }}>
                  {statistics.surrendered_count || 0}
                </div>
                <div>退保保单</div>
              </div>
            </Card>
            <Card size="small" style={{ flex: 1 }}>
              <div style={{ textAlign: 'center' }}>
                <div style={{ fontSize: 24, fontWeight: 'bold', color: '#13c2c2' }}>
                  {statistics.paid_commission_count || 0}
                </div>
                <div>已付佣金</div>
              </div>
            </Card>
          </div>
        </div>
      )}

      {/* 保单表格 */}
      <ProTable<PolicyInfo, PolicyListParams>
        headerTitle="保单管理"
        actionRef={actionRef}
        rowKey="id"
        search={{
          labelWidth: 120,
        }}
        // 添加列设置功能
        columnsState={{
          persistenceKey: 'policy-management-table',
          persistenceType: 'localStorage',
          defaultValue: {
            // 默认隐藏一些不常用的列
            customer_name_en: { show: false },
            referral_code: { show: false },
            referral_pm: { show: false },
            referral_sub_branch: { show: false },
            referral_date: { show: false },
            payment_date: { show: false },
            effective_date: { show: false },
            payment_method: { show: false },
            payment_years: { show: false },
            payment_periods: { show: false },
            referral_rate: { show: false },
            exchange_rate: { show: false },
            payment_pay_date: { show: false },
            is_employee: { show: false },
            product_type: { show: false },
            remark: { show: false },
          },
        }}
        toolBarRender={() => [
          <Button
            type="primary"
            key="primary"
            onClick={() => {
              setStepFormVisible(true);
            }}
          >
            <PlusOutlined /> 新建保单
          </Button>,
          <Button
            key="import"
            onClick={handleImport}
          >
            <ImportOutlined /> 批量导入
          </Button>,
          <Button
            key="export"
            onClick={handleExport}
          >
            <ExportOutlined /> 导出数据
          </Button>,
          <Button
            key="statistics"
            onClick={() => setShowStatistics(!showStatistics)}
          >
            <BarChartOutlined /> {showStatistics ? '隐藏' : '显示'}统计
          </Button>,
        ]}
        request={async (params, sort, filter) => {
          const queryParams: PolicyListParams = {
            page: params.current || 1,
            page_size: params.pageSize || 20,
            account_number: params.account_number,
            customer_number: params.customer_number,
            customer_name_cn: params.customer_name_cn,
            proposal_number: params.proposal_number,
            policy_currency: params.policy_currency,
            partner: params.partner,
            product_name: params.product_name,
            is_surrendered: params.is_surrendered,
            past_cooling_period: params.past_cooling_period,
            is_paid_commission: params.is_paid_commission,
            is_employee: params.is_employee,
          };

          try {
            const response = await getPolicyList(queryParams);
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
        tableAlertRender={({ selectedRowKeys, onCleanSelected }) => (
          <Space size={24}>
            <span>
              已选择 {selectedRowKeys.length} 项
              <Button type="link" size="small" onClick={onCleanSelected}>
                取消选择
              </Button>
            </span>
          </Space>
        )}
        tableAlertOptionRender={() => (
          <Space size={16}>
            <Button
              size="small"
              onClick={() => handleBatchUpdateStatus('is_surrendered', true)}
            >
              批量标记退保
            </Button>
            <Button
              size="small"
              onClick={() => handleBatchUpdateStatus('past_cooling_period', true)}
            >
              批量标记过冷静期
            </Button>
            <Button
              size="small"
              onClick={() => handleBatchUpdateStatus('is_paid_commission', true)}
            >
              批量标记已付佣金
            </Button>
            <Button size="small" danger onClick={handleBatchDelete}>
              批量删除
            </Button>
          </Space>
        )}
        scroll={{ x: 3500 }}
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

      {/* 创建表单 */}
      <ModalForm
        title="新建保单"
        width="400px"
        open={createModalVisible}
        onOpenChange={setCreateModalVisible}
        onFinish={async (value) => {
          const success = await handleAdd(value as PolicyCreateRequest);
          if (success) {
            setCreateModalVisible(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
          <ProFormText
          rules={[
            {
              required: true,
              message: '账户号为必填项',
            },
          ]}
            label="账户号"
          name="account_number"
          />
          <ProFormText
          rules={[
            {
              required: true,
              message: '客户姓名为必填项',
            },
          ]}
          label="客户姓名"
            name="customer_name_cn"
          />
          <ProFormText
          rules={[
            {
              required: true,
              message: '投保单号为必填项',
            },
          ]}
            label="投保单号"
          name="proposal_number"
          />
          <ProFormSelect
          rules={[
            {
              required: true,
              message: '保单币种为必填项',
            },
          ]}
            label="保单币种"
          name="policy_currency"
            options={currencyOptions}
        />
      </ModalForm>

      {/* 新的导入模态框 */}
      <ImportModal
        open={importModalVisible}
        onOpenChange={setImportModalVisible}
        onFinish={handleImportFinish}
      />

      {/* 新的导出模态框 */}
      <ExportModal
        open={exportModalVisible}
        onOpenChange={setExportModalVisible}
        selectedRows={selectedRows}
      />

      {/* 分步表单模态框（新增） */}
      <PolicyStepForm
        visible={stepFormVisible}
        onVisibleChange={setStepFormVisible}
        onSuccess={handleStepFormSuccess}
      />

      {/* 分步表单模态框（编辑） */}
      <PolicyStepForm
        visible={editStepFormVisible}
        onVisibleChange={handleEditStepFormCancel}
        onSuccess={handleEditStepFormSuccess}
        initialValues={currentRow}
      />
    </PageContainer>
  );
};

export default PolicyManagement; 