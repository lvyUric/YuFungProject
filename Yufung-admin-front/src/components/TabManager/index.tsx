import React, { useState, useEffect, useCallback } from 'react';
import { Tabs, Button, Dropdown, Menu } from 'antd';
import { CloseOutlined, ReloadOutlined, MoreOutlined } from '@ant-design/icons';
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

  // 将TabItem转换为Tabs组件的items格式
  useEffect(() => {
    const items = tabs.map((tab) => ({
      key: tab.key,
      label: (
        <div className="tab-label">
          {tab.icon && <span className="tab-icon">{tab.icon}</span>}
          <span className="tab-text">{tab.label}</span>
          {tab.closable !== false && (
            <CloseOutlined
              className="tab-close-icon"
              onClick={(e) => {
                e.stopPropagation();
                onTabRemove(tab.key);
              }}
            />
          )}
        </div>
      ),
      children: null, // 标签页内容由路由渲染
    }));
    setTabItems(items);
  }, [tabs, onTabRemove]);

  // 标签页切换处理
  const handleTabChange = (key: string) => {
    const targetTab = tabs.find(tab => tab.key === key);
    if (targetTab) {
      onTabChange(key);
      // 如果路径不同，进行路由跳转
      if (targetTab.path !== location.pathname) {
        history.push(targetTab.path);
      }
    }
  };

  // 右键菜单处理
  const handleContextMenu = useCallback((e: React.MouseEvent, tabKey: string) => {
    e.preventDefault();
    const targetTab = tabs.find(tab => tab.key === tabKey);
    if (!targetTab) return;

    const menu = (
      <Menu>
        <Menu.Item
          key="refresh"
          icon={<ReloadOutlined />}
          onClick={() => {
            // 刷新当前页面
            window.location.reload();
          }}
        >
          刷新页面
        </Menu.Item>
        {tabs.length > 1 && (
          <Menu.Item
            key="close"
            icon={<CloseOutlined />}
            onClick={() => onTabRemove(tabKey)}
          >
            关闭标签页
          </Menu.Item>
        )}
        <Menu.Item
          key="closeOthers"
          icon={<CloseOutlined />}
          onClick={() => {
            // 关闭其他标签页
            tabs.forEach(tab => {
              if (tab.key !== tabKey && tab.closable !== false) {
                onTabRemove(tab.key);
              }
            });
          }}
        >
          关闭其他标签页
        </Menu.Item>
        <Menu.Item
          key="closeAll"
          icon={<CloseOutlined />}
          onClick={() => {
            // 关闭所有标签页
            tabs.forEach(tab => {
              if (tab.closable !== false) {
                onTabRemove(tab.key);
              }
            });
          }}
        >
          关闭所有标签页
        </Menu.Item>
      </Menu>
    );

    // 这里可以使用Dropdown组件显示右键菜单
    // 由于antd的Dropdown需要特定的触发方式，这里简化处理
  }, [tabs, onTabRemove]);

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
        onTabClick={(key, e) => {
          // 处理右键点击
          if ('button' in e && e.button === 2) {
            handleContextMenu(e as React.MouseEvent, key);
          }
        }}
      />
    </div>
  );
};

export default TabManager; 