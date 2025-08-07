import React from 'react';
import {
  ModalForm,
  ProFormText,
  ProFormTextArea,
  ProFormDigit,
  ProFormDatePicker,
  ProFormSelect,
  ProFormGroup,
  ProCard,
} from '@ant-design/pro-components';
import { FormattedMessage, useIntl } from '@umijs/max';
import { message } from 'antd';
import { updateCompany } from '@/services/ant-design-pro/company';
import dayjs from 'dayjs';

export type UpdateFormProps = {
  modalOpen: boolean;
  values: API.CompanyInfo;
  onOpenChange: (visible: boolean) => void;
  onFinish: (value: boolean) => Promise<void>;
};

const UpdateForm: React.FC<UpdateFormProps> = ({
  modalOpen,
  values,
  onOpenChange,
  onFinish,
}) => {
  const intl = useIntl();

  const handleSubmit = async (formValues: API.UpdateCompanyRequest) => {
    try {
      if (!values.id) {
        message.error('公司ID不能为空');
        return false;
      }

      const updateData = {
        ...formValues,
      };

      // 如果选择了新的日期，格式化日期
      if (formValues.valid_start_date) {
        updateData.valid_start_date = dayjs(formValues.valid_start_date).format('YYYY-MM-DD');
      }
      if (formValues.valid_end_date) {
        updateData.valid_end_date = dayjs(formValues.valid_end_date).format('YYYY-MM-DD');
      }

      const response = await updateCompany(values.id, updateData);

      if (response.code === 200) {
        message.success(
          intl.formatMessage({
            id: 'pages.companyList.updateSuccess',
            defaultMessage: '更新成功',
          }),
        );
        await onFinish(true);
        return true;
      }

      message.error(response.message || '更新失败');
      return false;
    } catch (error: any) {
      message.error(
        error?.response?.data?.message ||
          intl.formatMessage({
            id: 'pages.companyList.updateError',
            defaultMessage: '更新失败，请重试',
          }),
      );
      return false;
    }
  };

  return (
    <ModalForm
      title={intl.formatMessage({
        id: 'pages.companyList.updateForm.title',
        defaultMessage: '编辑保险公司',
      })}
      width="1000px"
      open={modalOpen}
      onOpenChange={onOpenChange}
      onFinish={async (value) => {
        return await handleSubmit(value as API.UpdateCompanyRequest);
      }}
      modalProps={{
        destroyOnClose: true,
      }}
      layout="horizontal"
      grid
      initialValues={{
        ...values,
        valid_start_date: values.valid_start_date ? dayjs(values.valid_start_date) : undefined,
        valid_end_date: values.valid_end_date ? dayjs(values.valid_end_date) : undefined,
      }}
    >
      {/* 基本信息 */}
      <ProCard title="基本信息" bordered style={{ marginBottom: 16 }}>
        <ProFormGroup>
          <ProFormText
            name="company_name"
            label="公司名称（显示）"
            placeholder="请输入公司名称"
            colProps={{ span: 12 }}
            rules={[
              { required: true, message: '请输入公司名称！' },
              { min: 2, max: 100, message: '公司名称长度应为2-100个字符！' },
            ]}
          />
          <ProFormText
            name="company_code"
            label="保险公司代码"
            placeholder="请输入内部公司代码"
            colProps={{ span: 12 }}
            rules={[{ max: 50, message: '代码长度不能超过50个字符！' }]}
          />
        </ProFormGroup>
      </ProCard>

      {/* 负责人信息 */}
      <ProCard title="负责人信息" bordered style={{ marginBottom: 16 }}>
        <ProFormGroup>
          <ProFormText
            name="contact_person"
            label="联络人"
            placeholder="请输入联络人姓名"
            colProps={{ span: 12 }}
            rules={[{ max: 100, message: '姓名长度不能超过100个字符！' }]}
          />
        </ProFormGroup>
      </ProCard>

      {/* 联系方式 */}
      <ProCard title="联系方式" bordered style={{ marginBottom: 16 }}>
        <ProFormGroup>
          <ProFormText
            name="tel_no"
            label="Tel No.（固定电话）"
            placeholder="请输入固定电话"
            colProps={{ span: 8 }}
          />
          <ProFormText
            name="mobile"
            label="Mobile（移动电话）"
            placeholder="请输入移动电话"
            colProps={{ span: 8 }}
          />
          <ProFormText
            name="email"
            label="Email（邮件）"
            placeholder="请输入邮箱地址"
            colProps={{ span: 8 }}
            rules={[
              { required: true, message: '请输入邮箱地址！' },
              { type: 'email', message: '邮箱格式不正确！' },
            ]}
          />
        </ProFormGroup>
      </ProCard>

      {/* 中文地址 */}
      <ProCard title="Address（中文）" bordered style={{ marginBottom: 16 }}>
        <ProFormGroup>
          <ProFormText
            name="address_cn_province"
            label="省/自治区/直辖市"
            placeholder="请输入省/自治区/直辖市"
            colProps={{ span: 6 }}
            rules={[{ max: 50, message: '长度不能超过50个字符！' }]}
          />
          <ProFormText
            name="address_cn_city"
            label="市"
            placeholder="请输入城市"
            colProps={{ span: 6 }}
            rules={[{ max: 50, message: '长度不能超过50个字符！' }]}
          />
          <ProFormText
            name="address_cn_district"
            label="县/区"
            placeholder="请输入县/区"
            colProps={{ span: 6 }}
            rules={[{ max: 50, message: '长度不能超过50个字符！' }]}
          />
          <ProFormText
            name="address_cn_detail"
            label="详细地址"
            placeholder="请输入详细地址"
            colProps={{ span: 6 }}
            rules={[{ max: 200, message: '长度不能超过200个字符！' }]}
          />
        </ProFormGroup>
      </ProCard>

      {/* 英文地址 */}
      <ProCard title="Address（英文）" bordered style={{ marginBottom: 16 }}>
        <ProFormGroup>
          <ProFormText
            name="address_en_province"
            label="Province/State"
            placeholder="Enter province/state"
            colProps={{ span: 6 }}
            rules={[{ max: 50, message: '长度不能超过50个字符！' }]}
          />
          <ProFormText
            name="address_en_city"
            label="City"
            placeholder="Enter city"
            colProps={{ span: 6 }}
            rules={[{ max: 50, message: '长度不能超过50个字符！' }]}
          />
          <ProFormText
            name="address_en_district"
            label="District"
            placeholder="Enter district"
            colProps={{ span: 6 }}
            rules={[{ max: 50, message: '长度不能超过50个字符！' }]}
          />
          <ProFormText
            name="address_en_detail"
            label="Detailed Address"
            placeholder="Enter detailed address"
            colProps={{ span: 6 }}
            rules={[{ max: 200, message: '长度不能超过200个字符！' }]}
          />
        </ProFormGroup>
      </ProCard>

      {/* 业务信息 */}
      <ProCard title="业务信息" bordered style={{ marginBottom: 16 }}>
        <ProFormGroup>
          <ProFormText
            name="broker_code"
            label="Broker Code"
            placeholder="请输入经纪人代码"
            colProps={{ span: 8 }}
            rules={[{ max: 50, message: '长度不能超过50个字符！' }]}
          />
          <ProFormText
            name="link"
            label="Link"
            placeholder="请输入相关链接"
            colProps={{ span: 8 }}
            rules={[{ type: 'url', message: '请输入有效的URL！' }]}
          />
          <ProFormText
            name="username"
            label="Username"
            placeholder="请输入用户名"
            colProps={{ span: 8 }}
            rules={[
              { min: 3, max: 50, message: '用户名长度应为3-50个字符！' },
            ]}
          />
        </ProFormGroup>
        <ProFormGroup>
          <ProFormText.Password
            name="password"
            label="Password（留空则不修改）"
            placeholder="请输入新密码"
            colProps={{ span: 12 }}
            rules={[{ min: 8, message: '密码长度至少8个字符！' }]}
          />
        </ProFormGroup>
      </ProCard>

      {/* 系统字段 */}
      <ProCard title="系统设置" bordered style={{ marginBottom: 16 }}>
        <ProFormGroup>
          <ProFormDatePicker
            name="valid_start_date"
            label="有效期开始"
            placeholder="请选择有效期开始日期"
            colProps={{ span: 6 }}
            rules={[{ required: true, message: '请选择有效期开始日期！' }]}
          />
          <ProFormDatePicker
            name="valid_end_date"
            label="有效期结束"
            placeholder="请选择有效期结束日期"
            colProps={{ span: 6 }}
            rules={[{ required: true, message: '请选择有效期结束日期！' }]}
          />
          <ProFormDigit
            name="user_quota"
            label="用户配额"
            placeholder="请输入用户配额"
            colProps={{ span: 6 }}
            min={1}
            max={10000}
            rules={[{ required: true, message: '请输入用户配额！' }]}
          />
          <ProFormSelect
            name="status"
            label="状态"
            colProps={{ span: 6 }}
            options={[
              {
                label: '有效',
                value: 'active',
              },
              {
                label: '停用',
                value: 'inactive',
              },
            ]}
            placeholder="请选择状态"
            rules={[{ required: true, message: '请选择状态！' }]}
          />
        </ProFormGroup>
      </ProCard>

      {/* 扩展信息 */}
      <ProCard title="扩展信息" bordered style={{ marginBottom: 16 }}>
        <ProFormTextArea
          name="remark"
          label="备注信息"
          placeholder="请输入备注信息"
          rules={[{ max: 500, message: '长度不能超过500个字符！' }]}
        />
      </ProCard>
    </ModalForm>
  );
};

export default UpdateForm; 