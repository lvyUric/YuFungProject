import { LockOutlined } from '@ant-design/icons';
import { ProForm, ProFormText } from '@ant-design/pro-components';
import { PageContainer } from '@ant-design/pro-components';
import { FormattedMessage, useIntl } from '@umijs/max';
import { Card, App } from 'antd';
import React from 'react';
import { changePassword } from '@/services/ant-design-pro/api';

const ChangePassword: React.FC = () => {
  const { message } = App.useApp();
  const intl = useIntl();

  const handleSubmit = async (values: API.ChangePasswordParams & { confirmPassword: string }) => {
    try {
      const response = await changePassword({
        old_password: values.old_password,
        new_password: values.new_password,
      });

      if (response.success) {
        const successMessage = intl.formatMessage({
          id: 'pages.changePassword.success',
          defaultMessage: '密码修改成功！',
        });
        message.success(successMessage);
        return true;
      }

      message.error(response.message || '密码修改失败');
      return false;
    } catch (error: any) {
      const defaultFailureMessage = intl.formatMessage({
        id: 'pages.changePassword.failure',
        defaultMessage: '密码修改失败，请重试！',
      });

      let errorMessage = defaultFailureMessage;
      if (error?.response?.data?.message) {
        errorMessage = error.response.data.message;
      } else if (error?.message) {
        errorMessage = error.message;
      }

      message.error(errorMessage);
      return false;
    }
  };

  return (
    <PageContainer
      title={intl.formatMessage({
        id: 'pages.changePassword.title',
        defaultMessage: '修改密码',
      })}
    >
      <Card>
        <ProForm
          style={{ maxWidth: 400, margin: '0 auto' }}
          onFinish={async (values: API.ChangePasswordParams & { confirmPassword: string }) => {
            const success = await handleSubmit(values);
            if (success) {
              // 重置表单
              return true;
            }
            return false;
          }}
          submitter={{
            searchConfig: {
              submitText: intl.formatMessage({
                id: 'pages.changePassword.submit',
                defaultMessage: '修改密码',
              }),
            },
          }}
        >
          <ProFormText.Password
            name="old_password"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined />,
            }}
            placeholder={intl.formatMessage({
              id: 'pages.changePassword.oldPassword.placeholder',
              defaultMessage: '请输入当前密码',
            })}
            rules={[
              {
                required: true,
                message: (
                  <FormattedMessage
                    id="pages.changePassword.oldPassword.required"
                    defaultMessage="请输入当前密码！"
                  />
                ),
              },
            ]}
          />

          <ProFormText.Password
            name="new_password"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined />,
            }}
            placeholder={intl.formatMessage({
              id: 'pages.changePassword.newPassword.placeholder',
              defaultMessage: '请输入新密码',
            })}
            rules={[
              {
                required: true,
                message: (
                  <FormattedMessage
                    id="pages.changePassword.newPassword.required"
                    defaultMessage="请输入新密码！"
                  />
                ),
              },
              {
                min: 8,
                message: (
                  <FormattedMessage
                    id="pages.changePassword.newPassword.length"
                    defaultMessage="密码长度至少8个字符！"
                  />
                ),
              },
            ]}
          />

          <ProFormText.Password
            name="confirmPassword"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined />,
            }}
            placeholder={intl.formatMessage({
              id: 'pages.changePassword.confirmPassword.placeholder',
              defaultMessage: '请确认新密码',
            })}
            rules={[
              {
                required: true,
                message: (
                  <FormattedMessage
                    id="pages.changePassword.confirmPassword.required"
                    defaultMessage="请确认新密码！"
                  />
                ),
              },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('new_password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(
                    new Error(
                      intl.formatMessage({
                        id: 'pages.changePassword.confirmPassword.mismatch',
                        defaultMessage: '两次输入的密码不一致！',
                      })
                    )
                  );
                },
              }),
            ]}
          />
        </ProForm>
      </Card>
    </PageContainer>
  );
};

export default ChangePassword; 