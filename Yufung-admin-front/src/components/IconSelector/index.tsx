import React, { useState } from 'react';
import { Select, Input, Space, Popover, Button, Row, Col } from 'antd';
import {
  // 通用图标
  HomeOutlined,
  DashboardOutlined,
  SettingOutlined,
  UserOutlined,
  TeamOutlined,
  MenuOutlined,
  BankOutlined,
  FileTextOutlined,
  ContactsOutlined,
  ShoppingCartOutlined,
  AppstoreOutlined,
  DatabaseOutlined,
  UnorderedListOutlined,
  BuildOutlined,
  // 系统管理图标
  SafetyOutlined,
  KeyOutlined,
  AuditOutlined,
  MonitorOutlined,
  CloudServerOutlined,
  ToolOutlined,
  // 业务图标
  ShopOutlined,
  DollarOutlined,
  CreditCardOutlined,
  LineChartOutlined,
  BarChartOutlined,
  PieChartOutlined,
  // 其他常用图标
  FolderOutlined,
  FileOutlined,
  ControlOutlined,
  SearchOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  ReloadOutlined,
} from '@ant-design/icons';

const { Option } = Select;

// 图标映射表
const iconMap = {
  // 通用图标
  'home': HomeOutlined,
  'dashboard': DashboardOutlined,
  'setting': SettingOutlined,
  'user': UserOutlined,
  'team': TeamOutlined,
  'menu': MenuOutlined,
  'bank': BankOutlined,
  'file-text': FileTextOutlined,
  'contacts': ContactsOutlined,
  'shopping-cart': ShoppingCartOutlined,
  'appstore': AppstoreOutlined,
  'database': DatabaseOutlined,
  'unordered-list': UnorderedListOutlined,
  // 添加数据库中使用的图标
  'company': BankOutlined, // 公司管理使用银行图标
  'policy': FileTextOutlined, // 保单管理使用文档图标
  'system': SettingOutlined, // 系统管理使用设置图标
  // 系统管理图标
  'safety': SafetyOutlined,
  'key': KeyOutlined,
  'audit': AuditOutlined,
  'monitor': MonitorOutlined,
  'cloud-server': CloudServerOutlined,
  'tool': ToolOutlined,
  // 业务图标
  'shop': ShopOutlined,
  'dollar': DollarOutlined,
  'credit-card': CreditCardOutlined,
  'line-chart': LineChartOutlined,
  'bar-chart': BarChartOutlined,
  'pie-chart': PieChartOutlined,
  // 其他常用图标
  'folder': FolderOutlined,
  'file': FileOutlined,
  'control': ControlOutlined,
  'search': SearchOutlined,
  'plus': PlusOutlined,
  'edit': EditOutlined,
  'delete': DeleteOutlined,
  'eye': EyeOutlined,
  'reload': ReloadOutlined,
};

// 图标分类
const iconCategories = {
  '通用图标': ['home', 'dashboard', 'setting', 'user', 'team', 'menu', 'bank', 'appstore', 'database', 'company', 'system'],
  '系统管理': ['safety', 'key', 'audit', 'monitor', 'cloud-server', 'tool', 'folder', 'file', 'control', 'unordered-list'],
  '业务管理': ['shop', 'dollar', 'credit-card', 'line-chart', 'bar-chart', 'pie-chart', 'file-text', 'contacts', 'shopping-cart', 'policy'],
  '操作图标': ['search', 'plus', 'edit', 'delete', 'eye', 'reload'],
};

interface IconSelectorProps {
  value?: string;
  onChange?: (value: string) => void;
  placeholder?: string;
}

const IconSelector: React.FC<IconSelectorProps> = ({ value, onChange, placeholder = "请选择图标" }) => {
  const [searchText, setSearchText] = useState('');
  const [visible, setVisible] = useState(false);

  // 渲染图标
  const renderIcon = (iconName: string) => {
    const IconComponent = iconMap[iconName as keyof typeof iconMap];
    return IconComponent ? <IconComponent /> : null;
  };

  // 过滤图标
  const getFilteredIcons = () => {
    if (!searchText) return iconCategories;
    
    const filtered: { [key: string]: string[] } = {};
    Object.entries(iconCategories).forEach(([category, icons]) => {
      const matchedIcons = icons.filter(iconName => 
        iconName.toLowerCase().includes(searchText.toLowerCase())
      );
      if (matchedIcons.length > 0) {
        filtered[category] = matchedIcons;
      }
    });
    return filtered;
  };

  // 图标选择器内容
  const iconSelectorContent = (
    <div style={{ width: 320, maxHeight: 400, overflow: 'auto' }}>
      <Input
        placeholder="搜索图标"
        value={searchText}
        onChange={(e) => setSearchText(e.target.value)}
        style={{ marginBottom: 12 }}
      />
      {Object.entries(getFilteredIcons()).map(([category, icons]) => (
        <div key={category} style={{ marginBottom: 16 }}>
          <div style={{ fontWeight: 'bold', marginBottom: 8, color: '#666' }}>
            {category}
          </div>
          <Row gutter={[8, 8]}>
            {icons.map(iconName => (
              <Col span={6} key={iconName}>
                <Button
                  size="small"
                  style={{
                    width: '100%',
                    height: 36,
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    padding: '4px 8px',
                    backgroundColor: value === iconName ? '#1890ff' : undefined,
                    color: value === iconName ? 'white' : undefined,
                  }}
                  onClick={() => {
                    onChange?.(iconName);
                    setVisible(false);
                  }}
                >
                  {renderIcon(iconName)}
                  <span style={{ fontSize: 10, marginTop: 2 }}>{iconName}</span>
                </Button>
              </Col>
            ))}
          </Row>
        </div>
      ))}
    </div>
  );

  return (
    <Space.Compact style={{ width: '100%' }}>
      <Input
        value={value}
        onChange={(e) => onChange?.(e.target.value)}
        placeholder={placeholder}
        prefix={value ? renderIcon(value) : null}
        style={{ flex: 1 }}
      />
      <Popover
        content={iconSelectorContent}
        title="选择图标"
        trigger="click"
        open={visible}
        onOpenChange={setVisible}
        placement="bottomRight"
      >
        <Button>选择</Button>
      </Popover>
    </Space.Compact>
  );
};

export default IconSelector;

// 导出图标渲染函数供其他组件使用
export const renderMenuIcon = (iconName: string) => {
  if (!iconName) return null;
  const IconComponent = iconMap[iconName as keyof typeof iconMap];
  return IconComponent ? <IconComponent /> : null;
}; 