import React, { useState, useEffect } from 'react';
import { 
  PageContainer, 
  ProCard,
} from '@ant-design/pro-components';
import { 
  Card, 
  Button, 
  message, 
  Spin, 
  Space, 
  Divider, 
  Tag,
  Row,
  Col,
  Typography,
  Tabs,
  Avatar,
  Statistic,
  Badge,
  Descriptions,
} from 'antd';
import { 
  EditOutlined, 
  ArrowLeftOutlined,
  HistoryOutlined,
  FileTextOutlined,
  UserOutlined,
  BankOutlined,
  DollarOutlined,
  CalendarOutlined,
  SafetyOutlined,
  TeamOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons';
import { useParams, useNavigate } from '@umijs/max';
import { getPolicyDetail, type PolicyInfo } from '@/services/policy';
import ChangeRecord from '@/components/ChangeRecord';
import PolicyStepForm from '../components/PolicyStepForm';
import styles from './index.less';

const { Text, Title } = Typography;

const PolicyDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [policyData, setPolicyData] = useState<PolicyInfo | null>(null);
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [activeTab, setActiveTab] = useState('info');

  // 加载保单详情
  const loadPolicyDetail = async () => {
    if (!id) return;
    
    setLoading(true);
    try {
      const response = await getPolicyDetail(id);
      if (response.code === 200 && response.data) {
        setPolicyData(response.data);
      } else {
        message.error(response.message || '获取保单详情失败');
      }
    } catch (error) {
      console.error('Failed to load policy detail:', error);
      message.error('获取保单详情失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPolicyDetail();
  }, [id]);

  // 返回列表页
  const handleBack = () => {
    navigate('/business/policy');
  };

  // 编辑成功回调
  const handleEditSuccess = () => {
    setEditModalVisible(false);
    loadPolicyDetail(); // 重新加载数据
    message.success('保单更新成功');
  };

  // 格式化金额显示
  const formatCurrency = (amount: number, currency: string) => {
    if (!amount) return '-';
    return `${currency} ${amount.toLocaleString()}`;
  };

  // 格式化日期
  const formatDate = (dateStr?: string) => {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleDateString('zh-CN');
  };

  // 格式化布尔值
  const formatBoolean = (value: boolean, trueText = '是', falseText = '否') => {
    return value ? 
      <Badge status="success" text={trueText} /> : 
      <Badge status="default" text={falseText} />;
  };

  // 获取币种颜色
  const getCurrencyColor = (currency: string) => {
    const colors: Record<string, string> = {
      USD: 'green',
      HKD: 'blue', 
      CNY: 'red',
    };
    return colors[currency] || 'default';
  };

  if (loading) {
    return (
      <PageContainer>
        <div className={styles.loadingContainer}>
          <Spin size="large" />
          <div style={{ marginTop: 16 }}>
            <Text type="secondary">加载保单详情中...</Text>
          </div>
        </div>
      </PageContainer>
    );
  }

  if (!policyData) {
    return (
      <PageContainer>
        <Card className={styles.emptyCard}>
          <div className={styles.emptyContainer}>
            <InfoCircleOutlined className={styles.emptyIcon} />
            <Text type="secondary">保单不存在或已被删除</Text>
            <div style={{ marginTop: 16 }}>
              <Button type="primary" onClick={handleBack}>
                返回列表
              </Button>
            </div>
          </div>
        </Card>
      </PageContainer>
    );
  }

  return (
    <PageContainer
      title={
        <div className={styles.pageHeader}>
          <Space>
            <Button 
              type="text" 
              icon={<ArrowLeftOutlined />} 
              onClick={handleBack}
              className={styles.backButton}
            >
              返回
            </Button>
            <Divider type="vertical" />
            <div className={styles.titleSection}>
              <Avatar 
                size={40} 
                icon={<SafetyOutlined />} 
                className={styles.titleAvatar}
              />
              <div>
                <Title level={4} style={{ margin: 0 }}>
                  保单详情
                </Title>
                <Text type="secondary">#{policyData.serial_number}</Text>
              </div>
            </div>
          </Space>
        </div>
      }
      extra={[
        <Button
          key="edit"
          type="primary"
          icon={<EditOutlined />}
          onClick={() => setEditModalVisible(true)}
          className={styles.editButton}
        >
          编辑保单
        </Button>
      ]}
      tabList={[
        {
          tab: (
            <span className={styles.tabItem}>
              <FileTextOutlined />
              保单信息
            </span>
          ),
          key: 'info',
        },
        {
          tab: (
            <span className={styles.tabItem}>
              <HistoryOutlined />
              变更记录
            </span>
          ),
          key: 'changes',
        },
      ]}
      tabActiveKey={activeTab}
      onTabChange={setActiveTab}
      className={styles.pageContainer}
    >
      {activeTab === 'info' && (
        <div className={styles.contentContainer}>
          {/* 概览卡片 */}
          <Card className={styles.overviewCard}>
            <Row gutter={24}>
              <Col span={6}>
                <Statistic
                  title="实际缴纳保费"
                  value={policyData.actual_premium}
                  precision={2}
                  prefix={policyData.policy_currency}
                  valueStyle={{ color: '#3f8600' }}
                />
              </Col>
              <Col span={6}>
                <Statistic
                  title="AUM"
                  value={policyData.aum}
                  precision={2}
                  prefix={policyData.policy_currency}
                  valueStyle={{ color: '#1890ff' }}
                />
              </Col>
              <Col span={6}>
                <div className={styles.statusItem}>
                  <Text type="secondary">退保状态</Text>
                  <div style={{ marginTop: 4 }}>
                    {formatBoolean(policyData.is_surrendered, '已退保', '正常')}
                  </div>
                </div>
              </Col>
              <Col span={6}>
                <div className={styles.statusItem}>
                  <Text type="secondary">过冷静期</Text>
                  <div style={{ marginTop: 4 }}>
                    {formatBoolean(policyData.past_cooling_period)}
                  </div>
                </div>
              </Col>
            </Row>
          </Card>

          <Row gutter={24}>
            <Col span={12}>
              {/* 基本信息卡片 */}
              <Card 
                title={
                  <Space>
                    <UserOutlined className={styles.cardIcon} />
                    基本信息
                  </Space>
                }
                className={styles.infoCard}
              >
                <Descriptions column={1} size="small">
                  <Descriptions.Item label="序号">
                    <Tag color="blue">#{policyData.serial_number}</Tag>
                  </Descriptions.Item>
                  <Descriptions.Item label="账户号">
                    <Text code>{policyData.account_number || '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="客户号">
                    <Text strong>{policyData.customer_number}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="客户中文名">
                    <Text strong className={styles.customerName}>
                      {policyData.customer_name_cn}
                    </Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="客户英文名">
                    <Text>{policyData.customer_name_en || '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="投保单号">
                    <Text code>{policyData.proposal_number}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="保单币种">
                    <Tag color={getCurrencyColor(policyData.policy_currency)} icon={<DollarOutlined />}>
                      {policyData.policy_currency}
                    </Tag>
                  </Descriptions.Item>
                </Descriptions>
              </Card>

              {/* 缴费信息卡片 */}
              <Card 
                title={
                  <Space>
                    <DollarOutlined className={styles.cardIcon} />
                    缴费信息
                  </Space>
                }
                className={styles.infoCard}
                style={{ marginTop: 16 }}
              >
                <Descriptions column={1} size="small">
                  <Descriptions.Item label="缴费日期">
                    <Space>
                      <CalendarOutlined />
                      <Text>{formatDate(policyData.payment_date)}</Text>
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label="生效日期">
                    <Space>
                      <CalendarOutlined />
                      <Text>{formatDate(policyData.effective_date)}</Text>
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label="缴费方式">
                    <Tag color="processing">{policyData.payment_method}</Tag>
                  </Descriptions.Item>
                  <Descriptions.Item label="缴费年期">
                    <Text>{policyData.payment_years ? `${policyData.payment_years}年` : '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="期缴期数">
                    <Text>{policyData.payment_periods ? `${policyData.payment_periods}期` : '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="实际缴纳保费">
                    <Text strong className={styles.premiumAmount}>
                      {formatCurrency(policyData.actual_premium, policyData.policy_currency)}
                    </Text>
                  </Descriptions.Item>
                </Descriptions>
              </Card>
            </Col>

            <Col span={12}>
              {/* 转介信息卡片 */}
              <Card 
                title={
                  <Space>
                    <TeamOutlined className={styles.cardIcon} />
                    转介信息
                  </Space>
                }
                className={styles.infoCard}
              >
                <Descriptions column={1} size="small">
                  <Descriptions.Item label="合作伙伴">
                    <Text>{policyData.partner || '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="转介编号">
                    <Text code>{policyData.referral_code || '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="港分客户经理">
                    <Text>{policyData.hk_manager || '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="转介理财经理">
                    <Text>{policyData.referral_pm || '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="转介分行">
                    <Text>{policyData.referral_branch || '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="转介支行">
                    <Text>{policyData.referral_sub_branch || '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="转介日期">
                    <Space>
                      <CalendarOutlined />
                      <Text>{formatDate(policyData.referral_date)}</Text>
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label="转介比例">
                    <Text>{policyData.referral_rate ? `${policyData.referral_rate}%` : '未填写'}</Text>
                  </Descriptions.Item>
                </Descriptions>
              </Card>

              {/* 产品信息卡片 */}
              <Card 
                title={
                  <Space>
                    <BankOutlined className={styles.cardIcon} />
                    产品信息
                  </Space>
                }
                className={styles.infoCard}
                style={{ marginTop: 16 }}
              >
                <Descriptions column={1} size="small">
                  <Descriptions.Item label="保险公司">
                    <Text strong>{policyData.insurance_company || '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="产品名称">
                    <Text>{policyData.product_name || '未填写'}</Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="产品类型">
                    <Tag>{policyData.product_type || '未分类'}</Tag>
                  </Descriptions.Item>
                  <Descriptions.Item label="AUM">
                    <Text strong className={styles.aumAmount}>
                      {formatCurrency(policyData.aum, policyData.policy_currency)}
                    </Text>
                  </Descriptions.Item>
                  <Descriptions.Item label="已付佣金">
                    {formatBoolean(policyData.is_paid_commission)}
                  </Descriptions.Item>
                  <Descriptions.Item label="是否员工">
                    {formatBoolean(policyData.is_employee)}
                  </Descriptions.Item>
                </Descriptions>
              </Card>
            </Col>
          </Row>

          {/* 其他信息卡片 */}
          <Card 
            title={
              <Space>
                <InfoCircleOutlined className={styles.cardIcon} />
                其他信息
              </Space>
            }
            className={styles.infoCard}
            style={{ marginTop: 16 }}
          >
            <Row gutter={24}>
              <Col span={12}>
                <Descriptions column={1} size="small">
                  <Descriptions.Item label="创建时间">
                    <Space>
                      <CalendarOutlined />
                      <Text>{formatDate(policyData.created_at)}</Text>
                    </Space>
                  </Descriptions.Item>
                  <Descriptions.Item label="更新时间">
                    <Space>
                      <CalendarOutlined />
                      <Text>{formatDate(policyData.updated_at)}</Text>
                    </Space>
                  </Descriptions.Item>
                </Descriptions>
              </Col>
              {policyData.remark && (
                <Col span={12}>
                  <Descriptions column={1} size="small">
                    <Descriptions.Item label="备注">
                      <div className={styles.remarkContent}>
                        {policyData.remark}
                      </div>
                    </Descriptions.Item>
                  </Descriptions>
                </Col>
              )}
            </Row>
          </Card>
        </div>
      )}

      {activeTab === 'changes' && (
        <ChangeRecord
          recordId={policyData.policy_id}
          tableName="policies"
          showTitle={false}
          maxHeight={700}
        />
      )}

      {/* 编辑表单 */}
      <PolicyStepForm
        visible={editModalVisible}
        onVisibleChange={setEditModalVisible}
        onSuccess={handleEditSuccess}
        initialValues={policyData}
      />
    </PageContainer>
  );
};

export default PolicyDetail; 