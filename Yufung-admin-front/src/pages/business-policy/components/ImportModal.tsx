import React, { useState } from 'react';
import {
  Modal,
  Upload,
  Button,
  Steps,
  Table,
  message,
  Alert,
  Space,
  Checkbox,
  Typography,
  Card,
  Divider,
  Progress,
  Form,
} from 'antd';
import { InboxOutlined, DownloadOutlined, CheckCircleOutlined, ExclamationCircleOutlined, UploadOutlined } from '@ant-design/icons';
import type { UploadProps } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { 
  downloadPolicyTemplate, 
  previewPolicyImport, 
  importPoliciesFromFile,
  type PolicyImportResponse,
  type PolicyImportError 
} from '@/services/policy';

const { Dragger } = Upload;
const { Step } = Steps;
const { Text, Title } = Typography;

interface ImportModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onFinish?: (success: boolean) => void;
}

const ImportModal: React.FC<ImportModalProps> = ({ open, onOpenChange, onFinish }) => {
  const intl = useIntl();
  const [currentStep, setCurrentStep] = useState(0);
  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [previewData, setPreviewData] = useState<PolicyImportResponse | null>(null);
  const [importResult, setImportResult] = useState<PolicyImportResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [skipHeader, setSkipHeader] = useState(true);
  const [updateExisting, setUpdateExisting] = useState(false);

  const handleClose = () => {
    setCurrentStep(0);
    setUploadFile(null);
    setPreviewData(null);
    setImportResult(null);
    setSkipHeader(true);
    setUpdateExisting(false);
    onOpenChange(false);
  };

  // 下载模板
  const handleDownloadTemplate = async (format: 'xlsx' | 'csv') => {
    try {
      setLoading(true);
      const response = await downloadPolicyTemplate(format);
      
      // 创建下载链接
      const blob = new Blob([response], {
        type: format === 'xlsx' 
          ? 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
          : 'text/csv',
      });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `policy_template.${format}`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      message.success(intl.formatMessage({
        id: 'pages.policyList.import.templateDownloadSuccess',
        defaultMessage: '模板下载成功',
      }));
    } catch (error) {
      message.error(intl.formatMessage({
        id: 'pages.policyList.import.templateDownloadError',
        defaultMessage: '模板下载失败',
      }));
    } finally {
      setLoading(false);
    }
  };

  // 文件上传配置
  const uploadProps: UploadProps = {
    name: 'file',
    multiple: false,
    accept: '.xlsx,.xls,.csv',
    beforeUpload: (file) => {
      const isValidFormat = file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
        file.type === 'application/vnd.ms-excel' ||
        file.type === 'text/csv' ||
        file.name.endsWith('.xlsx') ||
        file.name.endsWith('.xls') ||
        file.name.endsWith('.csv');
      
      if (!isValidFormat) {
        message.error(intl.formatMessage({
          id: 'pages.policyList.import.invalidFileFormat',
          defaultMessage: '只支持 Excel(.xlsx,.xls) 和 CSV(.csv) 格式文件',
        }));
        return false;
      }

      const isLt10M = file.size / 1024 / 1024 < 10;
      if (!isLt10M) {
        message.error(intl.formatMessage({
          id: 'pages.policyList.import.fileSizeLimit',
          defaultMessage: '文件大小不能超过 10MB',
        }));
        return false;
      }

      setUploadFile(file);
      return false; // 阻止自动上传
    },
    fileList: uploadFile ? [{ uid: '1', name: uploadFile.name, status: 'done' }] : [],
    onRemove: () => {
      setUploadFile(null);
    },
  };

  // 预览数据
  const handlePreview = async () => {
    if (!uploadFile) {
      message.error(intl.formatMessage({
        id: 'pages.policyList.import.selectFile',
        defaultMessage: '请先选择文件',
      }));
      return;
    }

    try {
      setLoading(true);
      const formData = new FormData();
      formData.append('file', uploadFile);
      formData.append('skip_header', skipHeader.toString());
      formData.append('update_existing', updateExisting.toString());

      const response = await previewPolicyImport(formData);
      
      if (response.code === 200 && response.data) {
        setPreviewData(response.data);
        setCurrentStep(1);
      } else {
        message.error(response.message || intl.formatMessage({
          id: 'pages.policyList.import.previewError',
          defaultMessage: '预览失败',
        }));
      }
    } catch (error) {
      message.error(intl.formatMessage({
        id: 'pages.policyList.import.previewError',
        defaultMessage: '预览失败',
      }));
    } finally {
      setLoading(false);
    }
  };

  // 确认导入
  const handleImport = async () => {
    if (!uploadFile) return;

    try {
      setLoading(true);
      const formData = new FormData();
      formData.append('file', uploadFile);
      formData.append('skip_header', skipHeader.toString());
      formData.append('update_existing', updateExisting.toString());

      const response = await importPoliciesFromFile(formData);
      
      if (response.code === 200 && response.data) {
        setImportResult(response.data);
        setCurrentStep(2);
        
        if (response.data.error_count === 0) {
          message.success(intl.formatMessage({
            id: 'pages.policyList.import.success',
            defaultMessage: '导入成功',
          }));
          onFinish?.(true);
        } else if (response.data.success_count > 0) {
          message.warning(intl.formatMessage({
            id: 'pages.policyList.import.partialSuccess',
            defaultMessage: '部分数据导入成功，请查看错误详情',
          }));
          // 即使有部分错误，只要成功导入了数据，也刷新页面
          onFinish?.(true);
        } else {
          message.error(intl.formatMessage({
            id: 'pages.policyList.import.allFailed',
            defaultMessage: '所有数据导入失败',
          }));
          onFinish?.(false);
        }
      } else {
        message.error(response.message || intl.formatMessage({
          id: 'pages.policyList.import.error',
          defaultMessage: '导入失败',
        }));
        onFinish?.(false);
      }
    } catch (error) {
      message.error(intl.formatMessage({
        id: 'pages.policyList.import.error',
        defaultMessage: '导入失败',
      }));
    } finally {
      setLoading(false);
    }
  };

  // 预览表格列配置 - 按照模板字段顺序
  const previewColumns = [
    {
      title: '序号',
      dataIndex: 'serial_number',
      key: 'serial_number',
      width: 60,
    },
    {
      title: '账户号',
      dataIndex: 'account_number',
      key: 'account_number',
      width: 120,
    },
    {
      title: '客户号',
      dataIndex: 'customer_number',
      key: 'customer_number',
      width: 100,
    },
    {
      title: '客户中文名',
      dataIndex: 'customer_name_cn',
      key: 'customer_name_cn',
      width: 120,
    },
    {
      title: '客户英文名',
      dataIndex: 'customer_name_en',
      key: 'customer_name_en',
      width: 120,
    },
    {
      title: '投保单号',
      dataIndex: 'proposal_number',
      key: 'proposal_number',
      width: 150,
    },
    {
      title: '保单币种',
      dataIndex: 'policy_currency',
      key: 'policy_currency',
      width: 100,
    },
    {
      title: '合作伙伴',
      dataIndex: 'partner',
      key: 'partner',
      width: 120,
    },
    {
      title: '转介编号',
      dataIndex: 'referral_code',
      key: 'referral_code',
      width: 120,
    },
    {
      title: '港分客户经理',
      dataIndex: 'hk_manager',
      key: 'hk_manager',
      width: 120,
    },
    {
      title: '转介理财经理',
      dataIndex: 'referral_pm',
      key: 'referral_pm',
      width: 120,
    },
    {
      title: '转介分行',
      dataIndex: 'referral_branch',
      key: 'referral_branch',
      width: 120,
    },
    {
      title: '转介支行',
      dataIndex: 'referral_sub_branch',
      key: 'referral_sub_branch',
      width: 120,
    },
    {
      title: '转介日期',
      dataIndex: 'referral_date',
      key: 'referral_date',
      width: 120,
      render: (value: string) => value ? new Date(value).toLocaleDateString() : '-',
    },
    {
      title: '是否退保',
      dataIndex: 'is_surrendered',
      key: 'is_surrendered',
      width: 100,
      render: (value: boolean) => value ? '是' : '否',
    },
    {
      title: '缴费日期',
      dataIndex: 'payment_date',
      key: 'payment_date',
      width: 120,
      render: (value: string) => value ? new Date(value).toLocaleDateString() : '-',
    },
    {
      title: '生效日期',
      dataIndex: 'effective_date',
      key: 'effective_date',
      width: 120,
      render: (value: string) => value ? new Date(value).toLocaleDateString() : '-',
    },
    {
      title: '缴费方式',
      dataIndex: 'payment_method',
      key: 'payment_method',
      width: 120,
    },
    {
      title: '缴费年期',
      dataIndex: 'payment_years',
      key: 'payment_years',
      width: 100,
      render: (value: number) => value ? `${value}年` : '-',
    },
    {
      title: '期缴期数',
      dataIndex: 'payment_periods',
      key: 'payment_periods',
      width: 100,
      render: (value: number) => value ? `${value}期` : '-',
    },
    {
      title: '实际缴纳保费',
      dataIndex: 'actual_premium',
      key: 'actual_premium',
      width: 130,
      render: (value: number) => value?.toLocaleString() || '-',
    },
    {
      title: 'AUM',
      dataIndex: 'aum',
      key: 'aum',
      width: 120,
      render: (value: number) => value?.toLocaleString() || '-',
    },
    {
      title: '是否已过冷静期',
      dataIndex: 'past_cooling_period',
      key: 'past_cooling_period',
      width: 130,
      render: (value: boolean) => value ? '是' : '否',
    },
    {
      title: '是否支付佣金',
      dataIndex: 'is_paid_commission',
      key: 'is_paid_commission',
      width: 120,
      render: (value: boolean) => value ? '是' : '否',
    },
    {
      title: '转介费率',
      dataIndex: 'referral_rate',
      key: 'referral_rate',
      width: 100,
      render: (value: number) => value ? `${value}%` : '-',
    },
    {
      title: '汇率',
      dataIndex: 'exchange_rate',
      key: 'exchange_rate',
      width: 100,
      render: (value: number) => value ? value.toFixed(4) : '-',
    },
    {
      title: '预计转介费',
      dataIndex: 'expected_fee',
      key: 'expected_fee',
      width: 120,
      render: (value: number) => value?.toLocaleString() || '-',
    },
    {
      title: '支付日期',
      dataIndex: 'payment_pay_date',
      key: 'payment_pay_date',
      width: 120,
      render: (value: string) => value ? new Date(value).toLocaleDateString() : '-',
    },
    {
      title: '是否员工',
      dataIndex: 'is_employee',
      key: 'is_employee',
      width: 100,
      render: (value: boolean) => value ? '是' : '否',
    },
    {
      title: '承保公司',
      dataIndex: 'insurance_company',
      key: 'insurance_company',
      width: 120,
    },
    {
      title: '保险产品名称',
      dataIndex: 'product_name',
      key: 'product_name',
      width: 150,
    },
    {
      title: '产品类型',
      dataIndex: 'product_type',
      key: 'product_type',
      width: 120,
    },
    {
      title: '备注说明',
      dataIndex: 'remark',
      key: 'remark',
      width: 150,
      ellipsis: true,
    },
  ];

  // 错误表格列配置
  const errorColumns = [
    {
      title: '行号',
      dataIndex: 'row',
      key: 'row',
      width: 80,
    },
    {
      title: '错误信息',
      dataIndex: 'errors',
      key: 'errors',
      render: (errors: string[]) => errors.join('; '),
    },
    {
      title: '数据',
      dataIndex: 'data',
      key: 'data',
      render: (data: any) => JSON.stringify(data),
      ellipsis: true,
    },
  ];

  const renderStepContent = () => {
    switch (currentStep) {
      case 0:
        return (
          <div>
            <Alert
              message="导入说明"
              description="请上传Excel或CSV文件。注意：只有投保单号为必填字段，其他字段可选。如果数据列数不足，系统会自动补充空列。"
              type="info"
              showIcon
              style={{ marginBottom: 16 }}
            />
            
            <Form layout="vertical">
              <Form.Item label="上传文件">
                <Upload {...uploadProps}>
                  <Button icon={<UploadOutlined />}>选择文件</Button>
                </Upload>
                <div style={{ marginTop: 8, color: '#666', fontSize: '12px' }}>
                  支持格式：Excel(.xlsx, .xls) 和 CSV(.csv)，文件大小不超过10MB
                </div>
              </Form.Item>

              <Form.Item>
                <Space>
                  <Checkbox checked={skipHeader} onChange={(e) => setSkipHeader(e.target.checked)}>
                    跳过表头行
                  </Checkbox>
                  <Checkbox checked={updateExisting} onChange={(e) => setUpdateExisting(e.target.checked)}>
                    更新已存在的记录
                  </Checkbox>
                </Space>
              </Form.Item>

              <Form.Item>
                <Space>
                  <Button type="primary" onClick={handlePreview} loading={loading}>
                    预览数据
                  </Button>
                  <Button onClick={handleDownloadTemplate.bind(null, 'xlsx')}>
                    下载Excel模板
                  </Button>
                  <Button onClick={handleDownloadTemplate.bind(null, 'csv')}>
                    下载CSV模板
                  </Button>
                </Space>
              </Form.Item>
            </Form>
          </div>
        );

      case 1:
        return (
          <div>
            <Alert
              message={`预览成功！共 ${previewData?.total_count || 0} 条数据，其中 ${previewData?.error_count || 0} 条有错误`}
              type={previewData?.error_count === 0 ? 'success' : 'warning'}
              showIcon
              style={{ marginBottom: 16 }}
            />

            {previewData?.preview && previewData.preview.length > 0 && (
              <Card title="数据预览（前10条）" size="small" style={{ marginBottom: 16 }}>
                <Table
                  columns={previewColumns}
                  dataSource={previewData.preview.slice(0, 10)}
                  size="small"
                  pagination={false}
                  scroll={{ x: 800 }}
                  rowKey={(record, index) => `preview_${index}`}
                />
              </Card>
            )}

            {previewData?.errors && previewData.errors.length > 0 && (
              <Card title="错误详情" size="small">
                <Table
                  columns={errorColumns}
                  dataSource={previewData.errors}
                  size="small"
                  pagination={{ pageSize: 5 }}
                  rowKey="row"
                />
              </Card>
            )}
          </div>
        );

      case 2:
        return (
          <div>
            <div style={{ textAlign: 'center', marginBottom: 24 }}>
              {importResult?.error_count === 0 ? (
                <CheckCircleOutlined style={{ fontSize: 48, color: '#52c41a' }} />
              ) : (
                <ExclamationCircleOutlined style={{ fontSize: 48, color: '#faad14' }} />
              )}
              <Title level={4} style={{ marginTop: 16 }}>
                {importResult?.error_count === 0 ? '导入成功' : '部分导入成功'}
              </Title>
            </div>

            <Card size="small">
              <Space direction="vertical" style={{ width: '100%' }}>
                <div>
                  <Text strong>导入统计：</Text>
                </div>
                <div>
                  <Text>总计：{importResult?.total_count || 0} 条</Text>
                </div>
                <div>
                  <Text type="success">成功：{importResult?.success_count || 0} 条</Text>
                </div>
                <div>
                  <Text type="danger">失败：{importResult?.error_count || 0} 条</Text>
                </div>
                
                {importResult?.total_count && (
                  <Progress
                    percent={Math.round(((importResult.success_count || 0) / importResult.total_count) * 100)}
                    status={importResult.error_count === 0 ? 'success' : 'active'}
                  />
                )}
              </Space>
            </Card>

            {importResult?.errors && importResult.errors.length > 0 && (
              <Card title="错误详情" size="small" style={{ marginTop: 16 }}>
                <Table
                  columns={errorColumns}
                  dataSource={importResult.errors}
                  size="small"
                  pagination={{ pageSize: 5 }}
                  rowKey="row"
                />
              </Card>
            )}
          </div>
        );

      default:
        return null;
    }
  };

  const renderFooter = () => {
    const buttons = [];

    if (currentStep === 0) {
      buttons.push(
        <Button key="cancel" onClick={handleClose}>
          取消
        </Button>
      );
      buttons.push(
        <Button
          key="preview"
          type="primary"
          onClick={handlePreview}
          loading={loading}
          disabled={!uploadFile}
        >
          预览数据
        </Button>
      );
    } else if (currentStep === 1) {
      buttons.push(
        <Button key="back" onClick={() => setCurrentStep(0)}>
          上一步
        </Button>
      );
      buttons.push(
        <Button
          key="import"
          type="primary"
          onClick={handleImport}
          loading={loading}
          disabled={previewData?.total_count === 0}
        >
          确认导入
        </Button>
      );
    } else if (currentStep === 2) {
      buttons.push(
        <Button key="finish" type="primary" onClick={handleClose}>
          完成
        </Button>
      );
    }

    return buttons;
  };

  return (
    <Modal
      title={intl.formatMessage({
        id: 'pages.policyList.import.title',
        defaultMessage: '导入保单数据',
      })}
      open={open}
      onCancel={handleClose}
      footer={renderFooter()}
      width={800}
      destroyOnClose
    >
      <Steps current={currentStep} style={{ marginBottom: 24 }}>
        <Step title="上传文件" />
        <Step title="预览数据" />
        <Step title="导入结果" />
      </Steps>

      {renderStepContent()}
    </Modal>
  );
};

export default ImportModal; 