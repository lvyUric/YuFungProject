import { LinkOutlined } from '@ant-design/icons';
import type { Settings as LayoutSettings } from '@ant-design/pro-components';
import { SettingDrawer } from '@ant-design/pro-components';
import type { RequestConfig, RunTimeLayoutConfig } from '@umijs/max';
import { history, Link } from '@umijs/max';
import React from 'react';
import {
  AvatarDropdown,
  AvatarName,
  Footer,
  Question,
  SelectLang,
} from '@/components';
import { currentUser as queryCurrentUser } from '@/services/ant-design-pro/api';
import { getUserMenus, type UserMenuResponse } from '@/services/menu';
import { renderMenuIcon } from '@/components/IconSelector';
import TabManager from '@/components/TabManager';
import { useTabManager } from '@/hooks/useTabManager';
import defaultSettings from '../config/defaultSettings';
import { errorConfig } from './requestErrorConfig';
import '@ant-design/v5-patch-for-react-19';

const isDev = process.env.NODE_ENV === 'development';
const isDevOrTest = isDev || process.env.CI;
console.log('Current Environment:', {
  NODE_ENV: process.env.NODE_ENV,
  isDev,
  isDevOrTest
});

const loginPath = '/user/login';

// 将用户菜单转换为ProLayout菜单格式，启用国际化功能
const transformUserMenus = (userMenus: UserMenuResponse[] = []): any[] => {
  return userMenus.map((menu) => {
    // 根据菜单名称创建国际化key
    const createLocaleKey = (menuName: string, routePath: string) => {
      // 根据路由路径或菜单名称创建国际化key
      if (routePath) {
        const pathSegments = routePath.split('/').filter(Boolean);
        if (pathSegments.length >= 2) {
          return `menu.${pathSegments.join('.')}`;
        }
      }
      
      // 备用映射规则
      const nameMapping: Record<string, string> = {
        '业务管理': 'menu.business',
        '保单管理': 'menu.business.policy',
        '客户管理': 'menu.business.customer',
        '系统管理': 'menu.system',
        '用户管理': 'menu.system.user',
        '角色管理': 'menu.system.role',
        '菜单管理': 'menu.system.menu',
        '公司管理': 'menu.system.company',
        '仪表板': 'menu.dashboard',
        '欢迎': 'menu.welcome',
        '首页': 'menu.home',
      };
      
      return nameMapping[menuName] || `menu.${menuName.toLowerCase().replace(/\s+/g, '-')}`;
    };

    const localeKey = createLocaleKey(menu.menu_name, menu.route_path);
    console.log(`Menu: ${menu.menu_name} -> Locale: ${localeKey}`);

    const transformedMenu: any = {
      key: menu.menu_id, // 使用menu_id作为唯一标识
      name: menu.menu_name, // 保留原始名称作为备用
      path: menu.route_path,
      icon: menu.icon ? renderMenuIcon(menu.icon) : undefined, // 使用图标渲染函数
      locale: localeKey, // 添加国际化key
    };

    // 如果有子菜单，递归转换
    if (menu.children && menu.children.length > 0) {
      transformedMenu.routes = transformUserMenus(menu.children);
    }

    return transformedMenu;
  });
};

/**
 * @see https://umijs.org/docs/api/runtime-config#getinitialstate
 * */
export async function getInitialState(): Promise<{
  settings?: Partial<LayoutSettings>;
  currentUser?: API.CurrentUser;
  loading?: boolean;
  fetchUserInfo?: () => Promise<API.CurrentUser | undefined>;
  userMenus?: UserMenuResponse[];
}> {
  const fetchUserInfo = async () => {
    try {
      // 检查是否有token
      const token = localStorage.getItem('access_token');
      if (!token) {
        history.push(loginPath);
        return undefined;
      }

      const msg = await queryCurrentUser({
        skipErrorHandler: true,
      });
      
      // 适配新的响应格式
      if (msg.code === 200 && msg.data) {
        return msg.data;
      }
      return undefined;
    } catch (error) {
      history.push(loginPath);
    }
    return undefined;
  };

  const fetchUserMenus = async () => {
    try {
      const token = localStorage.getItem('access_token');
      if (!token) {
        return [];
      }

      const response = await getUserMenus();
      if (response.code === 200 && response.data) {
        return response.data;
      }
      return [];
    } catch (error) {
      console.error('Failed to fetch user menus:', error);
      return [];
    }
  };

  // 如果不是登录页面，执行
  if (history.location.pathname !== loginPath) {
    const currentUser = await fetchUserInfo();
    const userMenus = await fetchUserMenus();
    return {
      fetchUserInfo,
      currentUser,
      userMenus,
      settings: defaultSettings as Partial<LayoutSettings>,
    };
  }
  return {
    fetchUserInfo,
    settings: defaultSettings as Partial<LayoutSettings>,
  };
}

// 标签页管理组件
const TabManagerWrapper: React.FC = () => {
  const { tabs, activeKey, changeTab, removeTab } = useTabManager();
  
  return (
    <TabManager
      tabs={tabs}
      activeKey={activeKey}
      onTabChange={changeTab}
      onTabRemove={removeTab}
    />
  );
};

export const layout: RunTimeLayoutConfig = ({
  initialState,
  setInitialState,
}) => {
  // 获取用户菜单并转换格式
  const userMenus = transformUserMenus(initialState?.userMenus);

  return {
    actionsRender: () => [
      <Question key="doc" />,
      <SelectLang key="SelectLang" />,
    ],
    avatarProps: {
      src: initialState?.currentUser?.avatar,
      title: <AvatarName />,
      render: (_, avatarChildren) => (
        <AvatarDropdown menu>{avatarChildren}</AvatarDropdown>
      ),
    },
    waterMarkProps: {
      content: initialState?.currentUser?.display_name,
    },
    footerRender: () => <Footer />,
    onPageChange: () => {
      const { location } = history;
      // 如果没有登录，重定向到 login
      if (!initialState?.currentUser && location.pathname !== loginPath) {
        history.push(loginPath);
      }
    },
    // 配置动态菜单 - 完全替换默认菜单
    menuDataRender: (menuData) => {
      console.log('=== Menu Debug Info ===');
      console.log('ProLayout default menuData:', menuData);
      console.log('initialState exists:', !!initialState);
      console.log('userMenus from backend:', initialState?.userMenus);
      console.log('userMenus length:', initialState?.userMenus?.length || 0);
      console.log('transformed userMenus:', userMenus);
      console.log('userMenus length after transform:', userMenus?.length || 0);
      
      // 详细打印每个菜单项
      if (initialState?.userMenus) {
        console.log('Backend menu details:');
        initialState.userMenus.forEach((menu, index) => {
          console.log(`Menu ${index}:`, {
            menu_id: menu.menu_id,
            menu_name: menu.menu_name,
            route_path: menu.route_path,
            icon: menu.icon
          });
        });
      }
      
      if (userMenus && userMenus.length > 0) {
        console.log('Transformed menu details:');
        userMenus.forEach((menu, index) => {
          console.log(`Transformed Menu ${index}:`, {
            key: menu.key,
            name: menu.name,
            path: menu.path,
            icon: menu.icon
          });
        });
      }
      console.log('========================');
      
      // 如果没有用户菜单数据，返回空数组（不显示任何菜单）
      if (!initialState?.userMenus || initialState.userMenus.length === 0) {
        console.log('No user menus available, returning empty array');
        return [];
      }
      
      // 完全使用后端动态菜单数据，忽略静态路由配置
      console.log('Returning transformed userMenus:', userMenus);
      return userMenus;
    },
    bgLayoutImgList: [
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/D2LWSqNny4sAAAAAAAAAAAAAFl94AQBr',
        left: 85,
        bottom: 100,
        height: '303px',
      },
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/C2TWRpJpiC0AAAAAAAAAAAAAFl94AQBr',
        bottom: -68,
        right: -45,
        height: '303px',
      },
      {
        src: 'https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/F6vSTbj8KpYAAAAAAAAAAAAAFl94AQBr',
        bottom: 0,
        left: 0,
        width: '331px',
      },
    ],
    links: isDevOrTest
      ? [
          <Link key="openapi" to="/umi/plugin/openapi" target="_blank">
            <LinkOutlined />
            <span>OpenAPI 文档</span>
          </Link>,
        ]
      : [],
    menuHeaderRender: undefined,
    // 自定义 403 页面
    // unAccessible: <div>unAccessible</div>,
    // 增加一个 loading 的状态
    childrenRender: (children) => {
      // if (initialState?.loading) return <PageLoading />;
      return (
        <>
          {/* 添加标签页管理器 */}
          <TabManagerWrapper />
          {children}
          {isDevOrTest && (
            <SettingDrawer
              disableUrlParams
              enableDarkTheme
              settings={initialState?.settings}
              onSettingChange={(settings) => {
                setInitialState((preInitialState) => ({
                  ...preInitialState,
                  settings,
                }));
              }}
            />
          )}
        </>
      );
    },
    ...initialState?.settings,
  };
};

/**
 * @name request 配置，可以配置错误处理
 * 它基于 axios 和 ahooks 的 useRequest 提供了一套统一的网络请求和错误处理方案。
 * @doc https://umijs.org/docs/max/request#配置
 */
export const request: RequestConfig = {
  baseURL: 'http://106.52.172.124:8088',
  ...errorConfig,
};

// 添加调试日志
console.log('API Base URL:', 'http://106.52.172.124:8088');
