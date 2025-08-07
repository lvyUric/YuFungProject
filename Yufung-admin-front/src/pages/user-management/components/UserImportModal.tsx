import React, { useState } from 'react';
import { Modal, Upload, Form, Checkbox, message } from 'antd';
import { InboxOutlined } from '@ant-design/icons';
import type { UploadProps } from 'antd';

const { Dragger } = Upload;

interface UserImportModalProps {
  open: boolean;
  onCancel: () => void;
  onImport: (file: File, options: { skipHeader: boolean; updateExisting: boolean }) => Promise<void>;
  loading: boolean;
}

const UserImportModal: React.FC<UserImportModalProps> = ({
  open,
  onCancel,
  onImport,
  loading,
}) => {
  const [form] = Form.useForm();
  const [fileList, setFileList] = useState<any[]>([]);

  const uploadProps: UploadProps = {
    accept: '.xlsx,.xls,.csv',
    multiple: false,
    fileList,
    beforeUpload: (file) => {
      setFileList([file]);
      return false; // 阻止自动上传
    },
    onRemove: () => {
      setFileList([]);
    },
  };

  const handleOk = async () => {
    if (fileList.length === 0) {
      message.error('请选择要导入的文件');
      return;
    }

    try {
      const values = await form.validateFields();
      await onImport(fileList[0], {
        skipHeader: values.skipHeader || false,
        updateExisting: values.updateExisting || false,
      });
    } catch (error) {
      // 表单验证失败
    }
  };

  const handleCancel = () => {
    setFileList([]);
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title="导入用户"
      open={open}
      onCancel={handleCancel}
      onOk={handleOk}
      confirmLoading={loading}
      width={600}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          label="选择文件"
          required
        >
          <Dragger {...uploadProps}>
            <p className="ant-upload-drag-icon">
              <InboxOutlined />
            </p>
            <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
            <p className="ant-upload-hint">
              支持 .xlsx、.xls、.csv 格式文件
            </p>
          </Dragger>
        </Form.Item>

        <Form.Item
          name="skipHeader"
          valuePropName="checked"
        >
          <Checkbox>跳过表头行</Checkbox>
        </Form.Item>

        <Form.Item
          name="updateExisting"
          valuePropName="checked"
        >
          <Checkbox>更新已存在的用户</Checkbox>
        </Form.Item>
      </Form>

      <div style={{ marginTop: 16, padding: 16, backgroundColor: '#f5f5f5', borderRadius: 4 }}>
        <h4>导入说明：</h4>
        <ul style={{ marginLeft: 20, marginBottom: 0 }}>
          <li>文件格式：Excel (.xlsx, .xls) 或 CSV (.csv)</li>
          <li>列顺序：用户名、显示名称、密码、所属公司ID、角色ID(多个用逗号分隔)、邮箱地址、手机号码、备注信息</li>
          <li>必填字段：用户名、显示名称、密码、所属公司ID</li>
          <li>建议先下载模板文件，按照格式填写数据</li>
        </ul>
      </div>
    </Modal>
  );
};

export default UserImportModal; 