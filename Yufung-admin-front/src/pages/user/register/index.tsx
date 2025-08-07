import {
  LockOutlined,
  MailOutlined,
  PhoneOutlined,
  UserOutlined,
} from '@ant-design/icons';
import {
  LoginForm,
  ProFormText,
} from '@ant-design/pro-components';
import {
  FormattedMessage,
  Helmet,
  SelectLang,
  useIntl,
  history,
} from '@umijs/max';
import { Alert, App } from 'antd';
import { createStyles } from 'antd-style';
import React, { useState } from 'react';
import { Footer } from '@/components';
import { register } from '@/services/ant-design-pro/api';
import Settings from '../../../../config/defaultSettings';

const useStyles = createStyles(({ token }) => {
  return {
    action: {
      marginLeft: '8px',
      color: 'rgba(0, 0, 0, 0.2)',
      fontSize: '24px',
      verticalAlign: 'middle',
      cursor: 'pointer',
      transition: 'color 0.3s',
      '&:hover': {
        color: token.colorPrimaryActive,
      },
    },
    lang: {
      width: 42,
      height: 42,
      lineHeight: '42px',
      position: 'fixed',
      right: 16,
      borderRadius: token.borderRadius,
      ':hover': {
        backgroundColor: token.colorBgTextHover,
      },
    },
    container: {
      display: 'flex',
      flexDirection: 'column',
      height: '100vh',
      overflow: 'auto',
      backgroundImage:
        "url('https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/V-_oS6r-i7wAAAAAAAAAAAAAFl94AQBr')",
      backgroundSize: '100% 100%',
    },
  };
});

const Lang = () => {
  const { styles } = useStyles();

  return (
    <div className={styles.lang} data-lang>
      {SelectLang && <SelectLang />}
    </div>
  );
};

const RegisterMessage: React.FC<{
  content: string;
}> = ({ content }) => {
  return (
    <Alert
      style={{
        marginBottom: 24,
      }}
      message={content}
      type="error"
      showIcon
    />
  );
};

const Register: React.FC = () => {
  const [registerState, setRegisterState] = useState<{
    status?: string;
    message?: string;
  }>({});
  const { styles } = useStyles();
  const { message } = App.useApp();
  const intl = useIntl();

  const handleSubmit = async (values: API.RegisterParams) => {
    try {
      const response = await register(values);
      
      if (response.success) {
        const successMessage = intl.formatMessage({
          id: 'pages.register.success',
          defaultMessage: '注册成功！即将跳转到登录页面',
        });
        message.success(successMessage);
        
        // 延迟跳转到登录页面
        setTimeout(() => {
          history.push('/user/login');
        }, 2000);
        return;
      }

      setRegisterState({
        status: 'error',
        message: response.message || '注册失败'
      });
    } catch (error: any) {
      const defaultRegisterFailureMessage = intl.formatMessage({
        id: 'pages.register.failure',
        defaultMessage: '注册失败，请重试！',
      });
      
      let errorMessage = defaultRegisterFailureMessage;
      if (error?.response?.data?.message) {
        errorMessage = error.response.data.message;
      } else if (error?.message) {
        errorMessage = error.message;
      }
      
      message.error(errorMessage);
      setRegisterState({
        status: 'error',
        message: errorMessage
      });
    }
  };

  const { status } = registerState;

  return (
    <div className={styles.container}>
      <Helmet>
        <title>
          {intl.formatMessage({
            id: 'menu.register',
            defaultMessage: '注册页',
          })}
          {Settings.title && ` - ${Settings.title}`}
        </title>
      </Helmet>
      <Lang />
      <div
        style={{
          flex: '1',
          padding: '32px 0',
        }}
      >
        <LoginForm
          contentStyle={{
            minWidth: 280,
            maxWidth: '75vw',
          }}
          logo={<img alt="logo" src="/logo.svg" />}
          title="Yufung Admin"
          subTitle={intl.formatMessage({
            id: 'pages.register.subtitle',
            defaultMessage: '用户注册',
          })}
          onFinish={async (values) => {
            await handleSubmit(values as API.RegisterParams);
          }}
          submitter={{
            searchConfig: {
              submitText: intl.formatMessage({
                id: 'pages.register.submit',
                defaultMessage: '注册',
              }),
            },
          }}
        >
          {status === 'error' && (
            <RegisterMessage
              content={registerState.message || intl.formatMessage({
                id: 'pages.register.errorMessage',
                defaultMessage: '注册失败，请检查输入信息',
              })}
            />
          )}

          <ProFormText
            name="username"
            fieldProps={{
              size: 'large',
              prefix: <UserOutlined />,
            }}
            placeholder={intl.formatMessage({
              id: 'pages.register.username.placeholder',
              defaultMessage: '请输入用户名',
            })}
            rules={[
              {
                required: true,
                message: (
                  <FormattedMessage
                    id="pages.register.username.required"
                    defaultMessage="请输入用户名！"
                  />
                ),
              },
              {
                min: 3,
                max: 50,
                message: (
                  <FormattedMessage
                    id="pages.register.username.length"
                    defaultMessage="用户名长度应为3-50个字符！"
                  />
                ),
              },
            ]}
          />

          <ProFormText
            name="display_name"
            fieldProps={{
              size: 'large',
              prefix: <UserOutlined />,
            }}
            placeholder={intl.formatMessage({
              id: 'pages.register.displayName.placeholder',
              defaultMessage: '请输入显示名称',
            })}
            rules={[
              {
                required: true,
                message: (
                  <FormattedMessage
                    id="pages.register.displayName.required"
                    defaultMessage="请输入显示名称！"
                  />
                ),
              },
              {
                min: 2,
                max: 100,
                message: (
                  <FormattedMessage
                    id="pages.register.displayName.length"
                    defaultMessage="显示名称长度应为2-100个字符！"
                  />
                ),
              },
            ]}
          />

          <ProFormText.Password
            name="password"
            fieldProps={{
              size: 'large',
              prefix: <LockOutlined />,
            }}
            placeholder={intl.formatMessage({
              id: 'pages.register.password.placeholder',
              defaultMessage: '请输入密码',
            })}
            rules={[
              {
                required: true,
                message: (
                  <FormattedMessage
                    id="pages.register.password.required"
                    defaultMessage="请输入密码！"
                  />
                ),
              },
              {
                min: 8,
                message: (
                  <FormattedMessage
                    id="pages.register.password.length"
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
              id: 'pages.register.confirmPassword.placeholder',
              defaultMessage: '请确认密码',
            })}
            rules={[
              {
                required: true,
                message: (
                  <FormattedMessage
                    id="pages.register.confirmPassword.required"
                    defaultMessage="请确认密码！"
                  />
                ),
              },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error(intl.formatMessage({
                    id: 'pages.register.confirmPassword.mismatch',
                    defaultMessage: '两次输入的密码不一致！',
                  })));
                },
              }),
            ]}
          />

          <ProFormText
            name="email"
            fieldProps={{
              size: 'large',
              prefix: <MailOutlined />,
            }}
            placeholder={intl.formatMessage({
              id: 'pages.register.email.placeholder',
              defaultMessage: '请输入邮箱（可选）',
            })}
            rules={[
              {
                type: 'email',
                message: (
                  <FormattedMessage
                    id="pages.register.email.invalid"
                    defaultMessage="邮箱格式不正确！"
                  />
                ),
              },
            ]}
          />

          <ProFormText
            name="phone"
            fieldProps={{
              size: 'large',
              prefix: <PhoneOutlined />,
            }}
            placeholder={intl.formatMessage({
              id: 'pages.register.phone.placeholder',
              defaultMessage: '请输入手机号（可选）',
            })}
          />

          <div style={{ marginBottom: 24, textAlign: 'center' }}>
            <a
              onClick={() => {
                history.push('/user/login');
              }}
            >
              <FormattedMessage
                id="pages.register.backToLogin"
                defaultMessage="已有账号？立即登录"
              />
            </a>
          </div>
        </LoginForm>
      </div>
      <Footer />
    </div>
  );
};

export default Register; 