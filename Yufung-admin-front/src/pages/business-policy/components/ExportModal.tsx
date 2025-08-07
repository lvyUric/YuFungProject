import React, { useState } from 'react';
import {
  Modal,
  Button,
  Radio,
  Space,
  message,
  Alert,
  Card,
  Typography,
  Progress,
} from 'antd';
import { DownloadOutlined, ExportOutlined } from '@ant-design/icons';
import { FormattedMessage, useIntl } from '@umijs/max';
import { exportPoliciesToFile } from '@/services/policy';

const { Text, Title } = Typography;

interface ExportModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  selectedRows?: any[];
}

const ExportModal: React.FC<ExportModalProps> = ({ open, onOpenChange, selectedRows = [] }) => {
  const intl = useIntl();
  const [exportFormat, setExportFormat] = useState<'xlsx' | 'csv'>('xlsx');
  const [exportType, setExportType] = useState<'all' | 'selected'>('all');
  const [loading, setLoading] = useState(false);
  const [progress, setProgress] = useState(0);

  const handleClose = () => {
    setExportFormat('xlsx');
    setExportType('all');
    setProgress(0);
    onOpenChange(false);
  };

  const handleExport = async () => {
    try {
      setLoading(true);
      setProgress(0);

      // 模拟进度更新
      const progressInterval = setInterval(() => {
        setProgress(prev => {
          if (prev >= 90) {
            clearInterval(progressInterval);
            return prev;
          }
          return prev + 10;
        });
      }, 200);

      const policyIds = exportType === 'selected' && selectedRows.length > 0
        ? selectedRows.map((record) => record.policy_id)
        : undefined;

      const response = await exportPoliciesToFile({
        policy_ids: policyIds,
        export_type: exportFormat,
        format: exportFormat,
      });

      clearInterval(progressInterval);
      setProgress(100);

      // 创建下载链接
      const blob = new Blob([response], {
        type: exportFormat === 'xlsx'
          ? 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
          : 'text/csv',
      });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
      const exportScope = exportType === 'selected' ? '选中' : '全部';
      link.download = `保单导出_${exportScope}_${timestamp}.${exportFormat}`;
      
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);

      message.success(intl.formatMessage({
        id: 'pages.policyList.export.success',
        defaultMessage: '导出成功',
      }));

      setTimeout(() => {
        handleClose();
      }, 1000);

    } catch (error) {
      message.error(intl.formatMessage({
        id: 'pages.policyList.export.error',
        defaultMessage: '导出失败',
      }));
    } finally {
      setLoading(false);
    }
  };

  const getExportDescription = () => {
    if (exportType === 'selected') {
      return selectedRows.length > 0 
        ? `将导出选中的 ${selectedRows.length} 条保单记录`
        : '请先选择要导出的保单记录';
    }
    return '将导出所有保单记录';
  };

  return (
    <Modal
      title={intl.formatMessage({
        id: 'pages.policyList.export.title',
        defaultMessage: '导出保单数据',
      })}
      open={open}
      onCancel={handleClose}
      footer={[
        <Button key="cancel" onClick={handleClose} disabled={loading}>
          取消
        </Button>,
        <Button
          key="export"
          type="primary"
          icon={<ExportOutlined />}
          onClick={handleExport}
          loading={loading}
          disabled={exportType === 'selected' && selectedRows.length === 0}
        >
          {loading ? '导出中...' : '开始导出'}
        </Button>,
      ]}
      width={500}
      destroyOnClose
    >
      <Space direction="vertical" style={{ width: '100%' }} size="large">
        <Alert
          message="导出说明"
          description="支持导出为 Excel 或 CSV 格式，可选择导出全部数据或仅导出选中的记录"
          type="info"
          showIcon
        />

        <Card title="导出范围" size="small">
          <Radio.Group
            value={exportType}
            onChange={(e) => setExportType(e.target.value)}
          >
            <Space direction="vertical">
              <Radio value="all">导出全部保单</Radio>
              <Radio value="selected" disabled={selectedRows.length === 0}>
                导出选中保单 {selectedRows.length > 0 && `(${selectedRows.length} 条)`}
              </Radio>
            </Space>
          </Radio.Group>
          <div style={{ marginTop: 8, color: '#666' }}>
            <Text type="secondary">{getExportDescription()}</Text>
          </div>
        </Card>

        <Card title="文件格式" size="small">
          <Radio.Group
            value={exportFormat}
            onChange={(e) => setExportFormat(e.target.value)}
          >
            <Space direction="vertical">
              <Radio value="xlsx">
                <Space>
                  Excel (.xlsx)
                  <Text type="secondary">- 推荐，支持更丰富的格式</Text>
                </Space>
              </Radio>
              <Radio value="csv">
                <Space>
                  CSV (.csv)
                  <Text type="secondary">- 通用格式，兼容性更好</Text>
                </Space>
              </Radio>
            </Space>
          </Radio.Group>
        </Card>

        {loading && (
          <Card title="导出进度" size="small">
            <Progress
              percent={progress}
              status={progress === 100 ? 'success' : 'active'}
              strokeColor={{
                '0%': '#108ee9',
                '100%': '#87d068',
              }}
            />
            <Text type="secondary">
              {progress < 100 ? '正在处理数据...' : '导出完成！'}
            </Text>
          </Card>
        )}
      </Space>
    </Modal>
  );
};

export default ExportModal; 