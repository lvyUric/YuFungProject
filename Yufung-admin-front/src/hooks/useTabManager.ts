import { useState, useEffect, useCallback } from 'react';
import { history, useLocation } from '@umijs/max';
import type { TabItem } from '@/components/TabManager';

// 标签页存储的key
const TAB_STORAGE_KEY = 'yufung_tabs';

// 默认标签页
const DEFAULT_TABS: TabItem[] = [
  {
    key: 'welcome',
    label: '首页',
    path: '/welcome',
    closable: false,
  },
];

export const useTabManager = () => {
  const location = useLocation();
  const [tabs, setTabs] = useState<TabItem[]>(DEFAULT_TABS);
  const [activeKey, setActiveKey] = useState<string>('welcome');

  // 从localStorage加载标签页
  useEffect(() => {
    const savedTabs = localStorage.getItem(TAB_STORAGE_KEY);
    if (savedTabs) {
      try {
        const parsedTabs = JSON.parse(savedTabs);
        if (Array.isArray(parsedTabs) && parsedTabs.length > 0) {
          setTabs(parsedTabs);
          // 设置当前活跃的标签页
          const currentTab = parsedTabs.find(tab => tab.path === location.pathname);
          if (currentTab) {
            setActiveKey(currentTab.key);
          } else {
            setActiveKey(parsedTabs[0].key);
          }
        }
      } catch (error) {
        console.error('Failed to parse saved tabs:', error);
      }
    }
  }, []);

  // 保存标签页到localStorage
  const saveTabs = useCallback((newTabs: TabItem[]) => {
    localStorage.setItem(TAB_STORAGE_KEY, JSON.stringify(newTabs));
  }, []);

  // 添加标签页
  const addTab = useCallback((tab: TabItem) => {
    setTabs(prevTabs => {
      // 检查是否已存在相同路径的标签页
      const existingTab = prevTabs.find(t => t.path === tab.path);
      if (existingTab) {
        setActiveKey(existingTab.key);
        return prevTabs;
      }
      
      const newTabs = [...prevTabs, tab];
      saveTabs(newTabs);
      setActiveKey(tab.key);
      return newTabs;
    });
  }, [saveTabs]);

  // 移除标签页
  const removeTab = useCallback((targetKey: string) => {
    setTabs(prevTabs => {
      const targetIndex = prevTabs.findIndex(tab => tab.key === targetKey);
      if (targetIndex === -1) return prevTabs;

      const newTabs = prevTabs.filter(tab => tab.key !== targetKey);
      
      // 如果关闭的是当前活跃标签页，需要切换到其他标签页
      if (targetKey === activeKey) {
        let newActiveKey = activeKey;
        
        if (newTabs.length === 0) {
          // 如果没有标签页了，添加默认标签页
          newTabs.push(DEFAULT_TABS[0]);
          newActiveKey = DEFAULT_TABS[0].key;
        } else if (targetIndex === prevTabs.length - 1) {
          // 如果关闭的是最后一个标签页，切换到前一个
          newActiveKey = newTabs[newTabs.length - 1].key;
        } else {
          // 否则切换到下一个标签页
          newActiveKey = newTabs[targetIndex].key;
        }
        
        setActiveKey(newActiveKey);
        
        // 跳转到新活跃标签页的路径
        const newActiveTab = newTabs.find(tab => tab.key === newActiveKey);
        if (newActiveTab && newActiveTab.path !== location.pathname) {
          history.push(newActiveTab.path);
        }
      }
      
      saveTabs(newTabs);
      return newTabs;
    });
  }, [activeKey, location.pathname, saveTabs]);

  // 切换标签页
  const changeTab = useCallback((key: string) => {
    setActiveKey(key);
    const targetTab = tabs.find(tab => tab.key === key);
    if (targetTab && targetTab.path !== location.pathname) {
      history.push(targetTab.path);
    }
  }, [tabs, location.pathname]);

  // 根据路径添加或激活标签页
  const addOrActivateTab = useCallback((path: string, label: string, icon?: React.ReactNode) => {
    const existingTab = tabs.find(tab => tab.path === path);
    
    if (existingTab) {
      // 如果标签页已存在，激活它
      setActiveKey(existingTab.key);
      if (existingTab.path !== location.pathname) {
        history.push(existingTab.path);
      }
    } else {
      // 如果标签页不存在，添加新标签页
      const newTab: TabItem = {
        key: `tab_${Date.now()}`,
        label,
        path,
        closable: true,
        icon,
      };
      addTab(newTab);
    }
  }, [tabs, location.pathname, addTab]);

  // 监听路由变化，自动添加标签页
  useEffect(() => {
    // 跳过登录页面
    if (location.pathname === '/user/login') {
      return;
    }

    // 检查当前路径是否已有对应的标签页
    const existingTab = tabs.find(tab => tab.path === location.pathname);
    if (!existingTab && location.pathname !== '/welcome') {
      // 根据路径生成标签页标题
      const pathSegments = location.pathname.split('/').filter(Boolean);
      let title = '新页面';
      
      // 简单的路径到标题映射
      const pathTitleMap: Record<string, string> = {
        'user-management': '用户管理',
        'role-management': '角色管理',
        'menu-management': '菜单管理',
        'company-list': '公司管理',
        'business-policy': '保单管理',
        'business-customer': '客户管理',
        'system-config': '系统配置',
        'table-list': '列表页面',
      };
      
      if (pathSegments.length > 0) {
        const lastSegment = pathSegments[pathSegments.length - 1];
        title = pathTitleMap[lastSegment] || lastSegment;
      }
      
      addOrActivateTab(location.pathname, title);
    } else if (existingTab) {
      setActiveKey(existingTab.key);
    }
  }, [location.pathname, tabs, addOrActivateTab]);

  return {
    tabs,
    activeKey,
    addTab,
    removeTab,
    changeTab,
    addOrActivateTab,
  };
}; 