import React, { useState } from 'react';
import {
  Modal,
  Form,
  Radio,
  Select,
  Button,
  message,
  Space,
  Card,
  Typography,
  Divider,
  Alert,
} from 'antd';
import { DownloadOutlined, FileExcelOutlined, FileTextOutlined } from '@ant-design/icons';
import { FormattedMessage, useIntl } from '@umijs/max';
import { exportCompany } from '@/services/ant-design-pro/company';

const { Option } = Select;
const { Text } = Typography;

interface ExportModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  selectedRows?: API.CompanyInfo[];
  currentFilters?: {
    status?: string;
    keyword?: string;
  };
}

const ExportModal: React.FC<ExportModalProps> = ({ 
  open, 
  onOpenChange, 
  selectedRows = [],
  currentFilters = {} 
}) => {
  const intl = useIntl();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  const handleClose = () => {
    form.resetFields();
    onOpenChange(false);
  };

  const handleExport = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      const exportRequest: API.CompanyExportRequest = {
        export_type: values.export_type,
        format: values.format,
        template: values.template,
      };

      // 根据导出类型设置不同的参数
      if (values.export_type === 'selected') {
        if (selectedRows.length === 0) {
          message.error(intl.formatMessage({
            id: 'pages.companyList.export.noSelectedRows',
            defaultMessage: '请先选择要导出的数据',
          }));
          return;
        }
        exportRequest.ids = selectedRows.map(row => row.id).filter(Boolean) as string[];
      } else if (values.export_type === 'filtered') {
        exportRequest.status = currentFilters.status;
        exportRequest.keyword = currentFilters.keyword;
      }

      const response = await exportCompany(exportRequest);

      if (response.code === 200 && response.data) {
        // 如果返回文件URL，直接下载
        if (response.data.file_url) {
          const link = document.createElement('a');
          link.href = response.data.file_url;
          link.download = response.data.file_name || `companies_export_${Date.now()}.${values.format}`;
          document.body.appendChild(link);
          link.click();
          document.body.removeChild(link);
        }
        
        message.success(intl.formatMessage({
          id: 'pages.companyList.export.success',
          defaultMessage: '导出成功',
        }));
        
        handleClose();
      } else {
        message.error(response.message || intl.formatMessage({
          id: 'pages.companyList.export.error',
          defaultMessage: '导出失败',
        }));
      }
    } catch (error) {
      message.error(intl.formatMessage({
        id: 'pages.companyList.export.error',
        defaultMessage: '导出失败',
      }));
    } finally {
      setLoading(false);
    }
  };

  const getExportTypeOptions = () => {
    const options = [
      {
        value: 'all',
        label: '导出全部数据',
        description: '导出系统中的所有公司数据',
      },
    ];

    if (selectedRows.length > 0) {
      options.push({
        value: 'selected',
        label: `导出选中数据（${selectedRows.length}条）`,
        description: '导出当前选中的公司数据',
      });
    }

    if (currentFilters.status || currentFilters.keyword) {
      options.push({
        value: 'filtered',
        label: '导出筛选结果',
        description: '导出当前筛选条件下的公司数据',
      });
    }

    return options;
  };

  const formatOptions = [
    {
      value: 'xlsx',
      label: 'Excel 格式 (.xlsx)',
      icon: <FileExcelOutlined style={{ color: '#107c41' }} />,
      description: '推荐格式，支持丰富的数据格式和样式',
    },
    {
      value: 'csv',
      label: 'CSV 格式 (.csv)',
      icon: <FileTextOutlined style={{ color: '#666' }} />,
      description: '通用格式，可用于各种数据处理工具',
    },
  ];

  return (
    <Modal
      title={intl.formatMessage({
        id: 'pages.companyList.export.title',
        defaultMessage: '导出公司数据',
      })}
      open={open}
      onCancel={handleClose}
      footer={[
        <Button key="cancel" onClick={handleClose}>
          取消
        </Button>,
        <Button
          key="export"
          type="primary"
          icon={<DownloadOutlined />}
          loading={loading}
          onClick={handleExport}
        >
          开始导出
        </Button>,
      ]}
      width={600}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={{
          export_type: selectedRows.length > 0 ? 'selected' : 'all',
          format: 'xlsx',
          template: false,
        }}
      >
        <Alert
          message="导出说明"
          description="请选择导出范围和格式，系统将根据您的权限导出相应的数据。"
          type="info"
          showIcon
          style={{ marginBottom: 24 }}
        />

        <Form.Item
          name="export_type"
          label="导出范围"
          rules={[{ required: true, message: '请选择导出范围' }]}
        >
          <Radio.Group>
            <Space direction="vertical" style={{ width: '100%' }}>
              {getExportTypeOptions().map(option => (
                <Card key={option.value} size="small" style={{ width: '100%' }}>
                  <Radio value={option.value}>
                    <div>
                      <Text strong>{option.label}</Text>
                      <div>
                        <Text type="secondary" style={{ fontSize: '12px' }}>
                          {option.description}
                        </Text>
                      </div>
                    </div>
                  </Radio>
                </Card>
              ))}
            </Space>
          </Radio.Group>
        </Form.Item>

        <Divider />

        <Form.Item
          name="format"
          label="导出格式"
          rules={[{ required: true, message: '请选择导出格式' }]}
        >
          <Radio.Group>
            <Space direction="vertical" style={{ width: '100%' }}>
              {formatOptions.map(option => (
                <Card key={option.value} size="small" style={{ width: '100%' }}>
                  <Radio value={option.value}>
                    <Space>
                      {option.icon}
                      <div>
                        <Text strong>{option.label}</Text>
                        <div>
                          <Text type="secondary" style={{ fontSize: '12px' }}>
                            {option.description}
                          </Text>
                        </div>
                      </div>
                    </Space>
                  </Radio>
                </Card>
              ))}
            </Space>
          </Radio.Group>
        </Form.Item>

        <Divider />

        <Form.Item
          name="template"
          label="导出选项"
        >
          <Radio.Group>
            <Space direction="vertical">
              <Radio value={false}>
                导出数据文件
                <div>
                  <Text type="secondary" style={{ fontSize: '12px' }}>
                    导出包含实际数据的文件
                  </Text>
                </div>
              </Radio>
              <Radio value={true}>
                导出模板文件
                <div>
                  <Text type="secondary" style={{ fontSize: '12px' }}>
                    导出空白模板，用于批量导入时的格式参考
                  </Text>
                </div>
              </Radio>
            </Space>
          </Radio.Group>
        </Form.Item>

        {/* 当前筛选条件显示 */}
        {(currentFilters.status || currentFilters.keyword) && (
          <Card title="当前筛选条件" size="small" style={{ marginTop: 16 }}>
            <Space direction="vertical" size="small">
              {currentFilters.status && (
                <div>
                  <Text strong>状态：</Text>
                  <Text>{currentFilters.status}</Text>
                </div>
              )}
              {currentFilters.keyword && (
                <div>
                  <Text strong>关键词：</Text>
                  <Text>{currentFilters.keyword}</Text>
                </div>
              )}
            </Space>
          </Card>
        )}

        {/* 选中数据信息 */}
        {selectedRows.length > 0 && (
          <Card title="选中数据" size="small" style={{ marginTop: 16 }}>
            <Text>已选中 {selectedRows.length} 条公司数据</Text>
            <div style={{ marginTop: 8 }}>
              <Text type="secondary" style={{ fontSize: '12px' }}>
                选中的公司：{selectedRows.slice(0, 3).map(row => row.company_name).join('、')}
                {selectedRows.length > 3 && '...'}
              </Text>
            </div>
          </Card>
        )}
      </Form>
    </Modal>
  );
};

export default ExportModal; 