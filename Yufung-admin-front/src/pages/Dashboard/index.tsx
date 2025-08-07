import React from 'react';
import { Card, Row, Col, Statistic, Button } from 'antd';
import { UserOutlined, TeamOutlined, FileTextOutlined, SettingOutlined } from '@ant-design/icons';
import ActivityLogList from './components/ActivityLogList';

const Dashboard: React.FC = () => {
  return (
    <div style={{ padding: '24px' }}>
      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总用户数"
              value={1128}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="公司数量"
              value={93}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="保单总数"
              value={11280}
              prefix={<FileTextOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="系统配置"
              value={15}
              prefix={<SettingOutlined />}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 活动记录和快捷操作 */}
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={16}>
          <Card title="系统概览" extra={<Button type="link">查看更多</Button>}>
            <div style={{ height: '300px', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#999' }}>
              这里可以放置系统概览图表
            </div>
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <ActivityLogList limit={8} />
        </Col>
      </Row>

      {/* 快捷操作 */}
      <Row gutter={[16, 16]} style={{ marginTop: '24px' }}>
        <Col xs={24}>
          <Card title="快捷操作">
            <Row gutter={[16, 16]}>
              <Col xs={12} sm={6}>
                <Button type="primary" block icon={<UserOutlined />}>
                  用户管理
                </Button>
              </Col>
              <Col xs={12} sm={6}>
                <Button type="default" block icon={<TeamOutlined />}>
                  公司管理
                </Button>
              </Col>
              <Col xs={12} sm={6}>
                <Button type="default" block icon={<FileTextOutlined />}>
                  保单管理
                </Button>
              </Col>
              <Col xs={12} sm={6}>
                <Button type="default" block icon={<SettingOutlined />}>
                  系统设置
                </Button>
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard; 