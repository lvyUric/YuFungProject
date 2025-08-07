import React, { useState, useEffect } from 'react';
import { Card, List, Avatar, Tag, Spin, Empty } from 'antd';
import { UserOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { getRecentActivityLogs, ActivityLog, OPERATION_TYPE_LABELS, RESULT_STATUS_LABELS } from '@/services/activityLog';
import { formatDateTime, getRelativeTime } from '@/utils/date';

interface ActivityLogListProps {
  limit?: number;
}

const ActivityLogList: React.FC<ActivityLogListProps> = ({ limit = 5 }) => {
  const [loading, setLoading] = useState(false);
  const [logs, setLogs] = useState<ActivityLog[]>([]);

  // 获取最近的活动记录
  const fetchRecentLogs = async () => {
    setLoading(true);
    try {
      const data = await getRecentActivityLogs(limit);
      setLogs(data);
    } catch (error) {
      console.error('获取最近活动记录失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 获取操作类型标签颜色
  const getOperationTypeColor = (type: string) => {
    switch (type) {
      case 'create':
        return 'green';
      case 'delete':
        return 'red';
      case 'update':
        return 'blue';
      case 'view':
        return 'default';
      case 'export':
        return 'orange';
      case 'import':
        return 'purple';
      case 'login':
        return 'cyan';
      case 'logout':
        return 'magenta';
      default:
        return 'default';
    }
  };

  // 获取状态标签颜色
  const getStatusColor = (status: string) => {
    return status === 'success' ? 'green' : 'red';
  };

  // 渲染操作描述
  const renderOperationDesc = (log: ActivityLog) => {
    const { operation_desc, target_name, target_id } = log;
    
    if (target_name) {
      return (
        <div>
          <span>{operation_desc}</span>
          {target_id && (
            <div style={{ fontSize: '12px', color: '#999', marginTop: '2px' }}>
              ID: {target_id}
            </div>
          )}
        </div>
      );
    }
    
    return operation_desc;
  };

  // 渲染列表项
  const renderListItem = (log: ActivityLog) => (
    <List.Item>
      <List.Item.Meta
        avatar={
          <Avatar 
            icon={<UserOutlined />} 
            style={{ 
              backgroundColor: getOperationTypeColor(log.operation_type) === 'green' ? '#52c41a' : 
                             getOperationTypeColor(log.operation_type) === 'red' ? '#ff4d4f' : '#1890ff' 
            }}
          />
        }
        title={
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <span style={{ fontWeight: 'bold' }}>{log.username}</span>
            <div>
              <Tag color={getOperationTypeColor(log.operation_type)}>
                {OPERATION_TYPE_LABELS[log.operation_type as keyof typeof OPERATION_TYPE_LABELS] || log.operation_type}
              </Tag>
              <Tag color={getStatusColor(log.result_status)}>
                {RESULT_STATUS_LABELS[log.result_status as keyof typeof RESULT_STATUS_LABELS] || log.result_status}
              </Tag>
            </div>
          </div>
        }
        description={
          <div>
            <div style={{ marginBottom: '4px' }}>
              {renderOperationDesc(log)}
            </div>
            <div style={{ fontSize: '12px', color: '#999', display: 'flex', alignItems: 'center', gap: '16px' }}>
              <span>
                <ClockCircleOutlined style={{ marginRight: '4px' }} />
                {formatDateTime(log.operation_time)}
              </span>
              <span>公司: {log.company_name}</span>
              {log.ip_address && <span>IP: {log.ip_address}</span>}
              {log.execution_time > 0 && <span>耗时: {log.execution_time}ms</span>}
            </div>
            <div style={{ fontSize: '12px', color: '#999', marginTop: '2px' }}>
              {getRelativeTime(log.operation_time)}
            </div>
          </div>
        }
      />
    </List.Item>
  );

  // 初始化加载
  useEffect(() => {
    fetchRecentLogs();
  }, [limit]);

  return (
    <Card 
      title="最近活动记录" 
      size="small"
      style={{ height: '100%' }}
      bodyStyle={{ padding: '12px' }}
    >
      <Spin spinning={loading}>
        {logs.length > 0 ? (
          <List
            dataSource={logs}
            renderItem={renderListItem}
            size="small"
            style={{ maxHeight: '400px', overflow: 'auto' }}
          />
        ) : (
          <Empty 
            description="暂无活动记录" 
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            style={{ padding: '20px 0' }}
          />
        )}
      </Spin>
    </Card>
  );
};

export default ActivityLogList; 