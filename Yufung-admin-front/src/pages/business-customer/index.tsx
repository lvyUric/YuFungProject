import React from 'react';
import { PageContainer } from '@ant-design/pro-components';
import { Card, Typography } from 'antd';

const { Title, Paragraph } = Typography;

const CustomerManagement: React.FC = () => {
  return (
    <PageContainer
      header={{
        title: '客户管理',
        breadcrumb: {
          items: [
            {
              path: '/business',
              title: '业务管理',
            },
            {
              title: '客户管理',
            },
          ],
        },
      }}
    >
      <Card>
        <div style={{ textAlign: 'center', padding: '40px 0' }}>
          <Title level={3}>👥 客户管理页面</Title>
          <Paragraph>
            这个页面正在开发中...
          </Paragraph>
          <Paragraph type="secondary">
            将包含客户信息管理、客户分类、客户关系维护等功能
          </Paragraph>
        </div>
      </Card>
    </PageContainer>
  );
};

export default CustomerManagement; 