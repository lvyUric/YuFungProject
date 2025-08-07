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
} from 'antd';
import { InboxOutlined, DownloadOutlined, CheckCircleOutlined, ExclamationCircleOutlined } from '@ant-design/icons';
import type { UploadProps } from 'antd';
import { FormattedMessage, useIntl } from '@umijs/max';
import { downloadCompanyTemplate, previewCompanyImport, importCompany } from '@/services/ant-design-pro/company';

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
  const [previewData, setPreviewData] = useState<API.CompanyImportResponse | null>(null);
  const [importResult, setImportResult] = useState<API.CompanyImportResponse | null>(null);
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
      const response = await downloadCompanyTemplate(format);
      
      // 创建下载链接
      const blob = new Blob([response], {
        type: format === 'xlsx' 
          ? 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
          : 'text/csv',
      });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `company_template.${format}`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      message.success(intl.formatMessage({
        id: 'pages.companyList.import.templateDownloadSuccess',
        defaultMessage: '模板下载成功',
      }));
    } catch (error) {
      message.error(intl.formatMessage({
        id: 'pages.companyList.import.templateDownloadError',
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
          id: 'pages.companyList.import.invalidFileFormat',
          defaultMessage: '只支持 Excel(.xlsx,.xls) 和 CSV(.csv) 格式文件',
        }));
        return false;
      }

      const isLt10M = file.size / 1024 / 1024 < 10;
      if (!isLt10M) {
        message.error(intl.formatMessage({
          id: 'pages.companyList.import.fileSizeLimit',
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
        id: 'pages.companyList.import.selectFile',
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

      const response = await previewCompanyImport(formData);
      
      if (response.code === 200 && response.data) {
        setPreviewData(response.data);
        setCurrentStep(1);
      } else {
        message.error(response.message || intl.formatMessage({
          id: 'pages.companyList.import.previewError',
          defaultMessage: '预览失败',
        }));
      }
    } catch (error) {
      message.error(intl.formatMessage({
        id: 'pages.companyList.import.previewError',
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

      const response = await importCompany(formData);
      
      if (response.code === 200 && response.data) {
        setImportResult(response.data);
        setCurrentStep(2);
        
        if (response.data.error_count === 0) {
          message.success(intl.formatMessage({
            id: 'pages.companyList.import.success',
            defaultMessage: '导入成功',
          }));
          onFinish?.(true);
        } else {
          message.warning(intl.formatMessage({
            id: 'pages.companyList.import.partialSuccess',
            defaultMessage: '部分数据导入成功，请查看错误详情',
          }));
        }
      } else {
        message.error(response.message || intl.formatMessage({
          id: 'pages.companyList.import.error',
          defaultMessage: '导入失败',
        }));
      }
    } catch (error) {
      message.error(intl.formatMessage({
        id: 'pages.companyList.import.error',
        defaultMessage: '导入失败',
      }));
    } finally {
      setLoading(false);
    }
  };

  // 预览表格列配置
  const previewColumns = [
    {
      title: '公司名称',
      dataIndex: 'company_name',
      key: 'company_name',
      width: 150,
    },
    {
      title: '公司代码',
      dataIndex: 'company_code',
      key: 'company_code',
      width: 120,
    },
    {
      title: '联系人',
      dataIndex: 'contact_person',
      key: 'contact_person',
      width: 100,
    },
    {
      title: '联系电话',
      dataIndex: 'mobile',
      key: 'mobile',
      width: 130,
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
      width: 200,
    },
    {
      title: '省/直辖市',
      dataIndex: 'address_cn_province',
      key: 'address_cn_province',
      width: 100,
    },
    {
      title: '市',
      dataIndex: 'address_cn_city',
      key: 'address_cn_city',
      width: 80,
    },
    {
      title: 'Broker Code',
      dataIndex: 'broker_code',
      key: 'broker_code',
      width: 120,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
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
              message={intl.formatMessage({
                id: 'pages.companyList.import.uploadTip',
                defaultMessage: '请先下载模板，按照模板格式填写数据后上传',
              })}
              type="info"
              showIcon
              style={{ marginBottom: 16 }}
            />
            
            <Space direction="vertical" style={{ width: '100%' }}>
              <Card title="下载模板" size="small">
                <Space>
                  <Button
                    icon={<DownloadOutlined />}
                    onClick={() => handleDownloadTemplate('xlsx')}
                    loading={loading}
                  >
                    下载 Excel 模板
                  </Button>
                  <Button
                    icon={<DownloadOutlined />}
                    onClick={() => handleDownloadTemplate('csv')}
                    loading={loading}
                  >
                    下载 CSV 模板
                  </Button>
                </Space>
              </Card>

              <Card title="上传文件" size="small">
                <Dragger {...uploadProps}>
                  <p className="ant-upload-drag-icon">
                    <InboxOutlined />
                  </p>
                  <p className="ant-upload-text">
                    {intl.formatMessage({
                      id: 'pages.companyList.import.uploadText',
                      defaultMessage: '点击或拖拽文件到此区域上传',
                    })}
                  </p>
                  <p className="ant-upload-hint">
                    {intl.formatMessage({
                      id: 'pages.companyList.import.uploadHint',
                      defaultMessage: '支持 Excel(.xlsx,.xls) 和 CSV(.csv) 格式，文件大小不超过 10MB',
                    })}
                  </p>
                </Dragger>
              </Card>

              <Card title="导入选项" size="small">
                <Space direction="vertical">
                  <Checkbox
                    checked={skipHeader}
                    onChange={(e) => setSkipHeader(e.target.checked)}
                  >
                    跳过表头行
                  </Checkbox>
                  <Checkbox
                    checked={updateExisting}
                    onChange={(e) => setUpdateExisting(e.target.checked)}
                  >
                    更新已存在的公司（根据公司名称匹配）
                  </Checkbox>
                </Space>
              </Card>
            </Space>
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
        id: 'pages.companyList.import.title',
        defaultMessage: '导入公司数据',
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