import React, { useEffect } from 'react';
import { Modal, Form, Input, Select, Switch, message } from 'antd';
import type { UserItem, UserCreateRequest, UserUpdateRequest } from '@/services/ant-design-pro/userApi';
import { createUser, updateUser } from '@/services/ant-design-pro/userApi';

interface UserFormProps {
  open: boolean;
  onCancel: () => void;
  onFinish: (values: any) => Promise<void>;
  initialValues?: UserItem;
}

const UserForm: React.FC<UserFormProps> = ({
  open,
  onCancel,
  onFinish,
  initialValues,
}) => {
  const [form] = Form.useForm();
  const isEdit = !!initialValues;

  useEffect(() => {
    if (open) {
      if (initialValues) {
        // 编辑模式，设置初始值
        form.setFieldsValue({
          username: initialValues.username,
          display_name: initialValues.display_name,
          company_id: initialValues.company_id,
          role_ids: initialValues.role_ids,
          status: initialValues.status,
          email: initialValues.email,
          phone: initialValues.phone,
          remark: initialValues.remark,
        });
      } else {
        // 新建模式，重置表单
        form.resetFields();
      }
    }
  }, [open, initialValues, form]);

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      
      if (isEdit) {
        // 编辑用户
        const updateData: UserUpdateRequest = {
          display_name: values.display_name,
          role_ids: values.role_ids,
          email: values.email,
          phone: values.phone,
          remark: values.remark,
          status: values.status,
        };
        
        await updateUser(initialValues!.user_id, updateData);
        message.success('更新用户成功');
      } else {
        // 新建用户
        const createData: UserCreateRequest = {
          username: values.username,
          display_name: values.display_name,
          password: values.password,
          company_id: values.company_id,
          role_ids: values.role_ids,
          email: values.email,
          phone: values.phone,
          remark: values.remark,
        };
        
        await createUser(createData);
        message.success('创建用户成功');
      }
      
      await onFinish(values);
    } catch (error) {
      console.error('提交失败:', error);
      message.error(isEdit ? '更新用户失败' : '创建用户失败');
    }
  };

  return (
    <Modal
      title={isEdit ? '编辑用户' : '新建用户'}
      open={open}
      onCancel={onCancel}
      onOk={handleSubmit}
      width={600}
      destroyOnClose
    >
      <Form
        form={form}
        layout="vertical"
        preserve={false}
      >
        <Form.Item
          name="username"
          label="用户名"
          rules={[
            { required: true, message: '请输入用户名' },
            { min: 3, max: 50, message: '用户名长度应在3-50字符之间' },
          ]}
        >
          <Input 
            placeholder="请输入用户名" 
            disabled={isEdit} // 编辑时用户名不可修改
          />
        </Form.Item>

        <Form.Item
          name="display_name"
          label="显示名称"
          rules={[
            { required: true, message: '请输入显示名称' },
            { min: 2, max: 100, message: '显示名称长度应在2-100字符之间' },
          ]}
        >
          <Input placeholder="请输入显示名称" />
        </Form.Item>

        {!isEdit && (
          <Form.Item
            name="password"
            label="密码"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 8, message: '密码长度至少8位' },
            ]}
          >
            <Input.Password placeholder="请输入密码" />
          </Form.Item>
        )}

        <Form.Item
          name="company_id"
          label="所属公司"
          rules={[{ required: true, message: '请选择所属公司' }]}
        >
          <Select
            placeholder="请选择所属公司"
            disabled={isEdit} // 编辑时公司不可修改
            options={[
              // TODO: 从API获取公司列表
              { label: '中国平安保险经纪有限公司', value: 'CMP1735967088DA82E1D9' },
              { label: '太平洋保险经纪有限公司', value: 'CMP17359670884B2F8C3A' },
              { label: '阳光保险经纪有限公司', value: 'CMP1735967088E5C4A6B8' },
            ]}
          />
        </Form.Item>

        <Form.Item
          name="role_ids"
          label="用户角色"
          rules={[{ required: true, message: '请选择用户角色' }]}
        >
          <Select
            mode="multiple"
            placeholder="请选择用户角色"
            options={[
              // TODO: 从API获取角色列表
              { label: '平台管理员', value: 'ADMIN' },
              { label: '超级管理员', value: 'SUPER_ADMIN' },
              { label: '普通用户', value: 'USER' },
              { label: '只读用户', value: 'READONLY' },
            ]}
          />
        </Form.Item>

        {isEdit && (
          <Form.Item
            name="status"
            label="用户状态"
          >
            <Select
              options={[
                { label: '激活', value: 'active' },
                { label: '禁用', value: 'inactive' },
                { label: '锁定', value: 'locked' },
              ]}
            />
          </Form.Item>
        )}

        <Form.Item
          name="email"
          label="邮箱地址"
          rules={[
            { type: 'email', message: '请输入有效的邮箱地址' },
          ]}
        >
          <Input placeholder="请输入邮箱地址" />
        </Form.Item>

        <Form.Item
          name="phone"
          label="手机号码"
        >
          <Input placeholder="请输入手机号码" />
        </Form.Item>

        <Form.Item
          name="remark"
          label="备注信息"
        >
          <Input.TextArea 
            rows={3} 
            placeholder="请输入备注信息" 
            maxLength={500}
            showCount
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default UserForm; 