import React, { useState, useEffect } from 'react';
import { Card, Timeline, Tag, Spin, Button, Empty, Tooltip, Space, Typography, Divider } from 'antd';
import { 
  ClockCircleOutlined, 
  UserOutlined, 
  EditOutlined, 
  PlusOutlined, 
  DeleteOutlined,
  MoreOutlined,
  HistoryOutlined,
  GlobalOutlined,
  InfoCircleOutlined
} from '@ant-design/icons';
import { request } from '@umijs/max';
import styles from './index.less';

const { Text, Paragraph } = Typography;

// 临时类型定义
interface ChangeDetail {
  field_name: string;
  field_label: string;
  old_value: any;
  new_value: any;
  old_value_text: string;
  new_value_text: string;
}

interface ChangeRecordItem {
  id: string;
  change_id: string;
  table_name: string;
  record_id: string;
  user_id: string;
  username: string;
  company_id: string;
  change_type: 'insert' | 'update' | 'delete';
  old_values?: Record<string, any>;
  new_values?: Record<string, any>;
  changed_fields: string[];
  change_time: string;
  change_reason?: string;
  ip_address?: string;
  user_agent?: string;
  change_time_formatted: string;
  change_details: ChangeDetail[];
}

interface ChangeRecordParams {
  days?: number;
  page?: number;
  page_size?: number;
}

interface ChangeRecordResponse {
  code: number;
  message: string;
  data: {
    records: ChangeRecordItem[];
    total: number;
    page: number;
    page_size: number;
    has_more: boolean;
  };
}

// 临时API调用函数
const getPolicyChangeRecords = async (
  policyId: string,
  params?: ChangeRecordParams,
): Promise<ChangeRecordResponse> => {
  return request(`/api/policies/${policyId}/change-records`, {
    method: 'GET',
    params: {
      days: 10,
      page: 1,
      page_size: 10,
      ...params,
    },
  });
};

interface ChangeRecordProps {
  recordId: string;
  tableName?: string;
  showTitle?: boolean;
  maxHeight?: number;
}

const ChangeRecord: React.FC<ChangeRecordProps> = ({
  recordId,
  tableName = 'policies',
  showTitle = true,
  maxHeight = 600,
}) => {
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [records, setRecords] = useState<ChangeRecordItem[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(false);
  const pageSize = 10;

  // 加载变更记录
  const loadRecords = async (currentPage: number = 1, append: boolean = false) => {
    if (currentPage === 1) {
      setLoading(true);
    } else {
      setLoadingMore(true);
    }

    try {
      console.log('Loading change records for:', recordId, { page: currentPage, pageSize });
      const response = await getPolicyChangeRecords(recordId, {
        days: 30, // 默认查询30天内的记录
        page: currentPage,
        page_size: pageSize,
      });

      console.log('Change records response:', response);

      // 修正响应判断逻辑 - 使用code字段
      if (response.code === 200) {
        const newRecords = response.data.records || [];
        if (append) {
          setRecords(prev => [...prev, ...newRecords]);
        } else {
          setRecords(newRecords);
        }
        setTotal(response.data.total || 0);
        setHasMore(response.data.has_more || false);
      } else {
        console.error('API response error:', response.message);
      }
    } catch (error) {
      console.error('Failed to load change records:', error);
    } finally {
      setLoading(false);
      setLoadingMore(false);
    }
  };

  // 加载更多
  const handleLoadMore = () => {
    const nextPage = page + 1;
    setPage(nextPage);
    loadRecords(nextPage, true);
  };

  useEffect(() => {
    console.log('ChangeRecord component mounted with recordId:', recordId);
    if (recordId) {
      setPage(1);
      loadRecords(1, false);
    }
  }, [recordId]);

  // 获取操作类型图标和颜色
  const getChangeTypeConfig = (changeType: string) => {
    switch (changeType) {
      case 'insert':
        return {
          icon: <PlusOutlined />,
          color: 'success',
          text: '新增记录',
          bgColor: '#f6ffed',
          borderColor: '#b7eb8f',
        };
      case 'update':
        return {
          icon: <EditOutlined />,
          color: 'processing',
          text: '修改记录',
          bgColor: '#e6f7ff',
          borderColor: '#91d5ff',
        };
      case 'delete':
        return {
          icon: <DeleteOutlined />,
          color: 'error',
          text: '删除记录',
          bgColor: '#fff2f0',
          borderColor: '#ffadd2',
        };
      default:
        return {
          icon: <EditOutlined />,
          color: 'default',
          text: '变更记录',
          bgColor: '#fafafa',
          borderColor: '#d9d9d9',
        };
    }
  };

  // 格式化变更详情
  const renderChangeDetails = (details: any[]) => {
    if (!details || details.length === 0) {
      return (
        <div className={styles.emptyDetails}>
          <InfoCircleOutlined style={{ marginRight: 8, color: '#bfbfbf' }} />
          <Text type="secondary">无详细变更信息</Text>
        </div>
      );
    }

    return (
      <div className={styles.changeDetails}>
        {details.map((detail, index) => (
          <div key={index} className={styles.changeItem}>
            <Text strong className={styles.fieldLabel}>
              {detail.field_label}
            </Text>
            <div className={styles.changeValue}>
              {detail.old_value_text && (
                <>
                  <Text delete type="secondary" className={styles.oldValue}>
                    {detail.old_value_text}
                  </Text>
                  <Text type="secondary" className={styles.arrow}> → </Text>
                </>
              )}
              <Text mark className={styles.newValue}>
                {detail.new_value_text}
              </Text>
            </div>
          </div>
        ))}
      </div>
    );
  };

  // 渲染时间线项目
  const renderTimelineItem = (record: ChangeRecordItem) => {
    const typeConfig = getChangeTypeConfig(record.change_type);
    
    return {
      dot: (
        <div className={styles.timelineDot}>
          {React.cloneElement(typeConfig.icon, {
            style: { color: 'white' }
          })}
        </div>
      ),
      children: (
        <Card 
          size="small" 
          className={styles.changeCard}
          bodyStyle={{ padding: '16px 20px' }}
        >
          <div className={styles.changeHeader}>
            <Space size={[12, 8]} wrap>
              <Tag 
                color={typeConfig.color} 
                icon={typeConfig.icon}
                className={styles.typeTag}
              >
                {typeConfig.text}
              </Tag>
              
              <Tooltip title={`用户ID: ${record.user_id}`}>
                <Space size={4} className={styles.userInfo}>
                  <UserOutlined style={{ color: '#667eea' }} />
                  <Text strong>{record.username}</Text>
                </Space>
              </Tooltip>
              
              <Tooltip title={`完整时间: ${record.change_time_formatted || record.change_time}`}>
                <Space size={4} className={styles.timeInfo}>
                  <ClockCircleOutlined style={{ color: '#8c8c8c' }} />
                  <Text type="secondary" className={styles.relativeTime}>
                    {formatRelativeTime(record.change_time)}
                  </Text>
                </Space>
              </Tooltip>

              {record.ip_address && (
                <Tooltip title={`IP地址: ${record.ip_address}`}>
                  <Space size={4} className={styles.ipInfo}>
                    <GlobalOutlined style={{ color: '#8c8c8c' }} />
                    <Text type="secondary" className={styles.ipText}>
                      {record.ip_address}
                    </Text>
                  </Space>
                </Tooltip>
              )}
            </Space>
          </div>

          <div className={styles.changeContent}>
            {renderChangeDetails(record.change_details)}
          </div>

          {record.change_reason && (
            <div className={styles.changeReason}>
              <Space size={4}>
                <InfoCircleOutlined style={{ color: '#667eea' }} />
                <Text type="secondary">变更原因：</Text>
              </Space>
              <Text className={styles.reasonText}>{record.change_reason}</Text>
            </div>
          )}
        </Card>
      ),
    };
  };

  // 格式化相对时间
  const formatRelativeTime = (timeStr: string) => {
    const time = new Date(timeStr);
    const now = new Date();
    const diff = now.getTime() - time.getTime();
    
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);
    
    if (minutes < 1) return '刚刚';
    if (minutes < 60) return `${minutes}分钟前`;
    if (hours < 24) return `${hours}小时前`;
    if (days < 7) return `${days}天前`;
    
    return time.toLocaleDateString('zh-CN');
  };

  if (loading) {
    return (
      <Card className={styles.changeRecordCard}>
        <div className={styles.loadingWrapper}>
          <Spin size="large" />
          <div style={{ marginTop: 16 }}>
            <Text type="secondary">加载变更记录中...</Text>
          </div>
        </div>
      </Card>
    );
  }

  console.log('Rendering ChangeRecord, records length:', records.length);

  return (
    <Card 
      title={showTitle ? (
        <Space>
          <HistoryOutlined style={{ color: '#667eea' }} />
          <span>变更记录</span>
          {total > 0 && (
            <Tag color="processing" style={{ marginLeft: 8 }}>
              共 {total} 条记录
            </Tag>
          )}
        </Space>
      ) : undefined}
      size="small"
      className={styles.changeRecordCard}
      bodyStyle={{ padding: showTitle ? '24px' : '0' }}
    >
      {records.length === 0 ? (
        <div className={styles.emptyWrapper}>
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={
              <div>
                <Text type="secondary" style={{ fontSize: 16, marginBottom: 8, display: 'block' }}>
                  暂无变更记录
                </Text>
                <Text type="secondary" style={{ fontSize: 12 }}>
                  该保单暂时没有变更历史
                </Text>
              </div>
            }
            style={{ padding: '60px 0' }}
          />
          <Divider dashed />
          <div style={{ textAlign: 'center', marginTop: 16 }}>
            <Space direction="vertical" size={4}>
              <Text type="secondary" style={{ fontSize: 12 }}>
                <InfoCircleOutlined style={{ marginRight: 4 }} />
                记录ID: {recordId}
              </Text>
              <Text type="secondary" style={{ fontSize: 12 }}>
                调试信息: 请检查API是否正常返回数据
              </Text>
            </Space>
          </div>
        </div>
      ) : (
        <div style={{ maxHeight, overflowY: 'auto' }} className={styles.timelineWrapper}>
          <Timeline 
            mode="left"
            className={styles.changeTimeline}
            items={records.map(renderTimelineItem)}
          />
          
          {hasMore && (
            <div className={styles.loadMoreContainer}>
              <Button
                type="link"
                icon={<MoreOutlined />}
                loading={loadingMore}
                onClick={handleLoadMore}
                className={styles.loadMoreButton}
              >
                {loadingMore ? '加载中...' : `查看更多 (还有 ${total - records.length} 条)`}
              </Button>
            </div>
          )}
          
          {!hasMore && records.length > 0 && (
            <div className={styles.noMoreContainer}>
              <Divider dashed />
              <Text type="secondary">
                <HistoryOutlined style={{ marginRight: 4 }} />
                已显示全部 {total} 条记录
              </Text>
            </div>
          )}
        </div>
      )}
    </Card>
  );
};

export default ChangeRecord; 