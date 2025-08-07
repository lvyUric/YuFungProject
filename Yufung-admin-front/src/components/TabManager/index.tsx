import React, { useState, useEffect } from 'react';
import { Tabs } from 'antd';
import { history, useLocation } from '@umijs/max';
import type { TabsProps } from 'antd';
import './index.less';

export interface TabItem {
  key: string;
  label: string;
  path: string;
  closable?: boolean;
  icon?: React.ReactNode;
}

interface TabManagerProps {
  tabs: TabItem[];
  activeKey: string;
  onTabChange: (activeKey: string) => void;
  onTabRemove: (targetKey: string) => void;
  onTabAdd?: () => void;
}

const TabManager: React.FC<TabManagerProps> = ({
  tabs,
  activeKey,
  onTabChange,
  onTabRemove,
  onTabAdd,
}) => {
  const location = useLocation();
  const [tabItems, setTabItems] = useState<TabsProps['items']>([]);

  useEffect(() => {
    const items = tabs.map((tab) => ({
      key: tab.key,
      label: (
        <div className="tab-label">
          {tab.icon && <span className="tab-icon">{tab.icon}</span>}
          <span className="tab-text">{tab.label}</span>
        </div>
      ),
      children: null,
      closable: tab.closable !== false,
    }));
    setTabItems(items);
  }, [tabs]);

  const handleTabChange = (key: string) => {
    const targetTab = tabs.find(tab => tab.key === key);
    if (targetTab) {
      onTabChange(key);
      if (targetTab.path !== location.pathname) {
        history.push(targetTab.path);
      }
    }
  };

  return (
    <div className="tab-manager">
      <Tabs
        type="editable-card"
        activeKey={activeKey}
        onChange={handleTabChange}
        onEdit={(targetKey, action) => {
          if (action === 'add' && onTabAdd) {
            onTabAdd();
          } else if (action === 'remove' && typeof targetKey === 'string') {
            onTabRemove(targetKey);
          }
        }}
        items={tabItems}
        hideAdd={!onTabAdd}
        className="custom-tabs"
      />
    </div>
  );
};

export default TabManager; 