import {
  LockOutlined,
  MailOutlined,
  PhoneOutlined,
  UserOutlined,
} from '@ant-design/icons';
import {
  ProFormText,
} from '@ant-design/pro-components';
import {
  FormattedMessage,
  Helmet,
  SelectLang,
  useIntl,
  history,
} from '@umijs/max';
import { Alert, App, Button, Form } from 'antd';
import { createStyles } from 'antd-style';
import React, { useState } from 'react';
import { Footer } from '@/components';
import { register } from '@/services/ant-design-pro/api';
import Settings from '../../../../config/defaultSettings';

const useStyles = createStyles(({ token }) => {
  return {
    lang: {
      width: 42,
      height: 42,
      lineHeight: '42px',
      position: 'fixed',
      right: 16,
      top: 16,
      borderRadius: token.borderRadius,
      backgroundColor: 'rgba(255, 255, 255, 0.9)',
      backdropFilter: 'blur(10px)',
      zIndex: 1000,
      ':hover': {
        backgroundColor: 'rgba(255, 255, 255, 1)',
      },
    },
    container: {
      display: 'flex',
      flexDirection: 'column',
      minHeight: '100vh',
      width: '100vw',
      position: 'relative',
      backgroundImage:
        "url('https://dev.admin.lifebee.tech/imgs/instance--assets.svg')",
      backgroundSize: 'cover',
      backgroundPosition: 'center',
      backgroundRepeat: 'no-repeat',
      backgroundAttachment: 'fixed',
      '&::before': {
        content: '""',
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        background: 'linear-gradient(135deg, rgba(0, 0, 0, 0.4) 0%, rgba(0, 0, 0, 0.2) 100%)',
        zIndex: 1,
      },
    },
    content: {
      position: 'relative',
      zIndex: 2,
      flex: '1',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '20px',
      minHeight: '100vh',
    },
    logo: {
      width: '100px',
      height: '100px',
      objectFit: 'contain',
      filter: 'drop-shadow(0 4px 8px rgba(0, 0, 0, 0.15))',
      transition: 'transform 0.3s ease',
      '&:hover': {
        transform: 'scale(1.05)',
      },
    },
    formContainer: {
      width: '100%',
      maxWidth: '520px',
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      backdropFilter: 'blur(20px)',
      borderRadius: '24px',
      boxShadow: '0 20px 60px rgba(0, 0, 0, 0.15), 0 8px 32px rgba(0, 0, 0, 0.1)',
      border: '1px solid rgba(255, 255, 255, 0.3)',
      padding: '48px',
      position: 'relative',
      overflow: 'hidden',
      '&::before': {
        content: '""',
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        height: '4px',
        background: 'linear-gradient(90deg, #1890ff, #722ed1, #eb2f96)',
        borderRadius: '24px 24px 0 0',
      },
    },
    title: {
      fontSize: '32px',
      fontWeight: 600,
      color: '#1a1a1a',
      marginBottom: '12px',
      textAlign: 'center',
    },
    subtitle: {
      fontSize: '16px',
      color: '#666',
      textAlign: 'center',
      marginBottom: '40px',
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
      <div className={styles.content}>
        <div className={styles.formContainer}>
          <div style={{ textAlign: 'center', marginBottom: '40px' }}>
            <img alt="logo" src="/icons/昱丰logo鸟.png" className={styles.logo} />
            <div className={styles.title}>
              <FormattedMessage
                id="pages.register.title"
                defaultMessage="用户注册"
              />
            </div>
            <div className={styles.subtitle}>
              <FormattedMessage
                id="pages.register.subtitle"
                defaultMessage="创建您的Yufung Admin账户"
              />
            </div>
          </div>
          
          <Form
            layout="vertical"
            onFinish={handleSubmit}
            style={{ maxWidth: '100%' }}
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

            <Button
              type="primary"
              htmlType="submit"
              size="large"
              style={{
                width: '100%',
                height: '48px',
                fontSize: '16px',
                fontWeight: 500,
                borderRadius: '8px',
                marginBottom: '16px',
                marginTop: '16px',
              }}
            >
              <FormattedMessage
                id="pages.register.submit"
                defaultMessage="注册"
              />
            </Button>

            <div style={{ marginTop: '24px', textAlign: 'center' }}>
              <FormattedMessage
                id="pages.register.backToLogin"
                defaultMessage="已有账号？立即登录"
              />
              <a
                onClick={() => {
                  history.push('/user/login');
                }}
                style={{ marginLeft: '8px' }}
              >
                <FormattedMessage
                  id="pages.register.loginLink"
                  defaultMessage="登录"
                />
              </a>
            </div>
          </Form>
        </div>
      </div>
      <Footer />
    </div>
  );
};

export default Register; 