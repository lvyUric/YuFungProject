import React, { useState } from 'react';
import { PageContainer } from '@ant-design/pro-components';
import {
  Card,
  Button,
  Space,
  Popconfirm,
  message,
  Tag,
  Modal,
  Form,
  Input,
  Select,
  Switch,
  InputNumber,
  TreeSelect,
  Divider,
  Table,
  Tooltip,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  ReloadOutlined,
  MenuOutlined,
  FolderOutlined,
  FileOutlined,
  ControlOutlined,
  SearchOutlined,
} from '@ant-design/icons';
import { useRequest } from 'ahooks';
import type { ColumnsType } from 'antd/es/table';
import {
  getMenuList,
  createMenu,
  updateMenu,
  deleteMenu,
  batchUpdateMenuStatus,
  getMenuStats,
  type MenuInfo,
  type MenuStatsResponse,
} from '../../services/menu';
import IconSelector, { renderMenuIcon } from '../../components/IconSelector';

const { Option } = Select;

// 菜单类型选项
const menuTypeOptions = [
  { label: '目录', value: 'directory', icon: <FolderOutlined /> },
  { label: '菜单', value: 'menu', icon: <FileOutlined /> },
  { label: '按钮', value: 'button', icon: <ControlOutlined /> },
];

// 菜单管理页面
const MenuManagement: React.FC = () => {
  const [form] = Form.useForm();
  const [searchForm] = Form.useForm();
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'create' | 'edit' | 'view'>('create');
  const [currentMenu, setCurrentMenu] = useState<MenuInfo | null>(null);
  const [expandedRowKeys, setExpandedRowKeys] = useState<React.Key[]>([]);
  const [menuTreeData, setMenuTreeData] = useState<MenuInfo[]>([]);
  const [filteredData, setFilteredData] = useState<MenuInfo[]>([]);
  const [searchValues, setSearchValues] = useState<{
    menu_name?: string;
    status?: string;
  }>({});

  // 获取菜单统计信息
  const { data: statsData, refresh: refreshStats } = useRequest(getMenuStats, {});

  // 获取菜单列表
  const { data: menuData, loading, refresh } = useRequest(
    () => getMenuList({}),
    {
      onSuccess: (data) => {
        if (data.code === 200) {
          const menus = data.data?.menus || [];
          setMenuTreeData(menus);
          setFilteredData(menus);
          // 默认展开所有节点
          const allKeys = getAllKeys(menus);
          setExpandedRowKeys(allKeys);
        }
      },
    }
  );

  // 获取所有节点的key（用于展开）
  const getAllKeys = (menus: MenuInfo[]): string[] => {
    let keys: string[] = [];
    menus.forEach(menu => {
      keys.push(menu.menu_id);
      if (menu.children && menu.children.length > 0) {
        keys = keys.concat(getAllKeys(menu.children));
      }
    });
    return keys;
  };

  // 构建树选择器数据
  const buildTreeSelectData = (menus: MenuInfo[]): any[] => {
    return menus.map(menu => ({
      title: menu.menu_name,
      value: menu.menu_id,
      key: menu.menu_id,
      children: menu.children ? buildTreeSelectData(menu.children) : undefined,
      disabled: menu.menu_type === 'button', // 按钮类型不能作为父菜单
    }));
  };

  // 搜索过滤菜单
  const filterMenus = (menus: MenuInfo[], searchValues: any): MenuInfo[] => {
    return menus.filter(menu => {
      // 检查当前菜单是否匹配搜索条件
      let matches = true;
      
      if (searchValues.menu_name && !menu.menu_name.toLowerCase().includes(searchValues.menu_name.toLowerCase())) {
        matches = false;
      }
      
      if (searchValues.status && menu.status !== searchValues.status) {
        matches = false;
      }
      
      // 如果当前菜单匹配，则包含它和所有子菜单
      if (matches) {
        return true;
      }
      
      // 如果当前菜单不匹配，检查是否有子菜单匹配
      if (menu.children && menu.children.length > 0) {
        const filteredChildren = filterMenus(menu.children, searchValues);
        if (filteredChildren.length > 0) {
          // 如果有子菜单匹配，则包含当前菜单但只显示匹配的子菜单
          return true;
        }
      }
      
      return false;
    }).map(menu => {
      if (menu.children && menu.children.length > 0) {
        const filteredChildren = filterMenus(menu.children, searchValues);
        return {
          ...menu,
          children: filteredChildren
        };
      }
      return menu;
    });
  };

  // 处理搜索
  const handleSearch = () => {
    const values = searchForm.getFieldsValue();
    setSearchValues(values);
    
    if (!values.menu_name && !values.status) {
      setFilteredData(menuTreeData);
    } else {
      const filtered = filterMenus(menuTreeData, values);
      setFilteredData(filtered);
      // 搜索时展开所有节点
      const allKeys = getAllKeys(filtered);
      setExpandedRowKeys(allKeys);
    }
  };

  // 重置搜索
  const handleReset = () => {
    searchForm.resetFields();
    setSearchValues({});
    setFilteredData(menuTreeData);
    // 重置时展开所有节点
    const allKeys = getAllKeys(menuTreeData);
    setExpandedRowKeys(allKeys);
  };

  // 表格列定义 - 树形表格
  const columns: ColumnsType<MenuInfo> = [
    {
      title: '菜单名称',
      dataIndex: 'menu_name',
      key: 'menu_name',
      width: 250,
      render: (text: string, record: MenuInfo) => (
        <Space>
          {/* 优先显示数据库中的自定义图标，否则显示默认类型图标 */}
          {record.icon ? renderMenuIcon(record.icon) : (
            <>
              {record.menu_type === 'directory' && <FolderOutlined style={{ color: '#1890ff' }} />}
              {record.menu_type === 'menu' && <FileOutlined style={{ color: '#52c41a' }} />}
              {record.menu_type === 'button' && <ControlOutlined style={{ color: '#faad14' }} />}
            </>
          )}
          <span>{text}</span>
        </Space>
      ),
    },
    {
      title: '排序',
      dataIndex: 'sort_order',
      key: 'sort_order',
      width: 80,
      align: 'center',
    },
    {
      title: '请求地址',
      dataIndex: 'route_path',
      key: 'route_path',
      width: 200,
      render: (text: string) => text ? <code style={{ fontSize: '12px' }}>{text}</code> : '#',
    },
    {
      title: '类型',
      dataIndex: 'menu_type',
      key: 'menu_type',
      width: 80,
      render: (type: string) => {
        const config = {
          directory: { color: 'blue', text: '目录' },
          menu: { color: 'green', text: '菜单' },
          button: { color: 'orange', text: '按钮' },
        };
        return <Tag color={config[type as keyof typeof config]?.color}>
          {config[type as keyof typeof config]?.text}
        </Tag>;
      },
    },
    {
      title: '可见',
      dataIndex: 'visible',
      key: 'visible',
      width: 80,
      align: 'center',
      render: (visible: boolean) => (
        <Tag color={visible ? 'green' : 'red'}>
          {visible ? '是' : '否'}
        </Tag>
      ),
    },
    {
      title: '权限标识',
      dataIndex: 'permission_code',
      key: 'permission_code',
      width: 150,
      render: (text: string) => text ? <code style={{ fontSize: '12px' }}>{text}</code> : '-',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 150,
      render: (text: string) => new Date(text).toLocaleString(),
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="查看详情">
            <Button
              type="link"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => handleView(record)}
            />
          </Tooltip>
          <Tooltip title="修改">
            <Button
              type="link"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            />
          </Tooltip>
          <Tooltip title="删除">
            <Popconfirm
              title="确定要删除这个菜单吗？"
              description="删除后将无法恢复，且会影响相关权限配置"
              onConfirm={() => handleDelete(record.menu_id)}
              okText="确定"
              cancelText="取消"
            >
              <Button
                type="link"
                size="small"
                danger
                icon={<DeleteOutlined />}
              />
            </Popconfirm>
          </Tooltip>
        </Space>
      ),
    },
  ];

  // 处理查看
  const handleView = (record: MenuInfo) => {
    setCurrentMenu(record);
    setModalType('view');
    setModalVisible(true);
    form.setFieldsValue({
      ...record,
    });
  };

  // 处理编辑
  const handleEdit = (record: MenuInfo) => {
    setCurrentMenu(record);
    setModalType('edit');
    setModalVisible(true);
    form.setFieldsValue({
      ...record,
    });
  };

  // 处理新建
  const handleCreate = () => {
    setCurrentMenu(null);
    setModalType('create');
    setModalVisible(true);
    form.resetFields();
    form.setFieldsValue({
      menu_type: 'menu',
      visible: true,
      status: 'enable',
      sort_order: 1,
    });
  };

  // 处理删除
  const handleDelete = async (menuId: string) => {
    try {
      const response = await deleteMenu(menuId);
      if (response.code === 200) {
        message.success('删除成功');
        refresh();
        refreshStats();
      } else {
        message.error(response.message);
      }
    } catch (error) {
      message.error('删除失败');
    }
  };

  // 批量操作
  const handleBatchStatus = async (status: string) => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要操作的菜单');
      return;
    }

    try {
      const response = await batchUpdateMenuStatus({
        menu_ids: selectedRowKeys.map(key => String(key)),
        status,
      });

      if (response.code === 200) {
        message.success(`批量${status === 'enable' ? '启用' : '禁用'}成功`);
        setSelectedRowKeys([]);
        refresh();
        refreshStats();
      } else {
        message.error(response.message);
      }
    } catch (error) {
      message.error('批量操作失败');
    }
  };

  // 提交表单
  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();

      let response;
      if (modalType === 'create') {
        response = await createMenu(values);
      } else if (modalType === 'edit' && currentMenu) {
        response = await updateMenu(currentMenu.menu_id, values);
      }

      if (response && response.code === 200) {
        message.success(modalType === 'create' ? '创建成功' : '更新成功');
        setModalVisible(false);
        refresh();
        refreshStats();
      } else {
        message.error(response?.message || '操作失败');
      }
    } catch (error) {
      console.error('表单验证失败:', error);
    }
  };

  // 统计卡片
  const StatsCards = () => {
    const stats = (statsData as any)?.data as MenuStatsResponse | undefined;
    return (
      <div style={{ marginBottom: 16 }}>
        <div style={{ display: 'flex', gap: 16 }}>
          <Card size="small" style={{ flex: 1 }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 'bold', color: '#1890ff' }}>
                {stats?.total_menus || 0}
              </div>
              <div>总菜单数</div>
            </div>
          </Card>
          <Card size="small" style={{ flex: 1 }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 'bold', color: '#52c41a' }}>
                {stats?.enabled_menus || 0}
              </div>
              <div>启用菜单</div>
            </div>
          </Card>
          <Card size="small" style={{ flex: 1 }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 'bold', color: '#f5222d' }}>
                {stats?.disabled_menus || 0}
              </div>
              <div>禁用菜单</div>
            </div>
          </Card>
          <Card size="small" style={{ flex: 1 }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 'bold', color: '#722ed1' }}>
                {stats?.directory_menus || 0}
              </div>
              <div>目录</div>
            </div>
          </Card>
          <Card size="small" style={{ flex: 1 }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 'bold', color: '#fa8c16' }}>
                {stats?.page_menus || 0}
              </div>
              <div>页面</div>
            </div>
          </Card>
          <Card size="small" style={{ flex: 1 }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 'bold', color: '#13c2c2' }}>
                {stats?.button_menus || 0}
              </div>
              <div>按钮</div>
            </div>
          </Card>
        </div>
      </div>
    );
  };

  return (
    <PageContainer>
      <StatsCards />
      
      <Card>
        {/* 搜索表单 */}
        <Form
          form={searchForm}
          layout="inline"
          style={{ marginBottom: 16 }}
          onFinish={handleSearch}
        >
          <Form.Item name="menu_name" label="菜单名称">
            <Input placeholder="请输入菜单名称" allowClear style={{ width: 200 }} />
          </Form.Item>
          <Form.Item name="status" label="菜单状态">
            <Select placeholder="菜单状态" allowClear style={{ width: 120 }}>
              <Option value="enable">正常</Option>
              <Option value="disable">停用</Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
                搜索
              </Button>
              <Button onClick={handleReset}>
                重置
              </Button>
            </Space>
          </Form.Item>
        </Form>

        {/* 操作按钮 */}
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
          <Space>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleCreate}
            >
              新增
            </Button>
            <Button
              onClick={() => handleBatchStatus('enable')}
              disabled={selectedRowKeys.length === 0}
            >
              批量启用
            </Button>
            <Button
              onClick={() => handleBatchStatus('disable')}
              disabled={selectedRowKeys.length === 0}
            >
              批量停用
            </Button>
            <Button onClick={() => setExpandedRowKeys(getAllKeys(filteredData))}>
              展开所有
            </Button>
            <Button onClick={() => setExpandedRowKeys([])}>
              折叠所有
            </Button>
          </Space>
          <Button
            icon={<ReloadOutlined />}
            onClick={() => {
              refresh();
              refreshStats();
            }}
          >
            刷新
          </Button>
        </div>

        {/* 树形表格 */}
        <Table
          columns={columns}
          dataSource={filteredData}
          rowKey="menu_id"
          loading={loading}
          pagination={false}
          expandable={{
            expandedRowKeys: expandedRowKeys,
            onExpandedRowsChange: (keys) => setExpandedRowKeys([...keys]),
            childrenColumnName: 'children',
          }}
          rowSelection={{
            selectedRowKeys,
            onChange: setSelectedRowKeys,
          }}
          scroll={{ x: 1200 }}
          size="small"
        />
      </Card>

      {/* 菜单详情/编辑/新建模态框 */}
      <Modal
        title={
          modalType === 'create'
            ? '新增菜单'
            : modalType === 'edit'
            ? '修改菜单'
            : '菜单详情'
        }
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={modalType === 'view' ? () => setModalVisible(false) : handleSubmit}
        okText={modalType === 'view' ? '关闭' : '确定'}
        cancelText="取消"
        width={600}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          disabled={modalType === 'view'}
        >
          <Form.Item
            name="parent_id"
            label="上级菜单"
            tooltip="选择上级菜单，留空表示顶级菜单"
          >
            <TreeSelect
              placeholder="请选择上级菜单"
              allowClear
              treeData={buildTreeSelectData(menuTreeData)}
              treeDefaultExpandAll
            />
          </Form.Item>

          <Form.Item
            name="menu_name"
            label="菜单名称"
            rules={[{ required: true, message: '请输入菜单名称' }]}
          >
            <Input placeholder="请输入菜单名称" />
          </Form.Item>

          <Form.Item
            name="menu_type"
            label="菜单类型"
            rules={[{ required: true, message: '请选择菜单类型' }]}
          >
            <Select placeholder="请选择菜单类型">
              {menuTypeOptions.map(option => (
                <Option key={option.value} value={option.value}>
                  <Space>
                    {option.icon}
                    {option.label}
                  </Space>
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="route_path"
            label="路由地址"
            tooltip="访问的路由地址，如：/user-management"
          >
            <Input placeholder="请输入路由地址" />
          </Form.Item>

          <Form.Item
            name="component"
            label="组件路径"
            tooltip="组件路径，如：./user-management"
          >
            <Input placeholder="请输入组件路径" />
          </Form.Item>

          <Form.Item
            name="permission_code"
            label="权限标识"
            tooltip="权限标识，如：system:user:list"
          >
            <Input placeholder="请输入权限标识" />
          </Form.Item>

          <Form.Item
            name="icon"
            label="菜单图标"
            tooltip="菜单图标名称"
          >
            <IconSelector
              value={form.getFieldValue('icon')}
              onChange={(icon) => form.setFieldsValue({ icon })}
            />
          </Form.Item>

          <Form.Item
            name="sort_order"
            label="显示排序"
            rules={[{ required: true, message: '请输入显示排序' }]}
          >
            <InputNumber
              min={0}
              max={9999}
              placeholder="请输入显示排序"
              style={{ width: '100%' }}
            />
          </Form.Item>

          <Form.Item
            name="visible"
            label="显示状态"
            valuePropName="checked"
            tooltip="选择隐藏时，路由将不会出现在侧边栏，但仍然可以访问"
          >
            <Switch checkedChildren="显示" unCheckedChildren="隐藏" />
          </Form.Item>

          <Form.Item
            name="status"
            label="菜单状态"
            valuePropName="checked"
            getValueFromEvent={(checked) => checked ? 'enable' : 'disable'}
            getValueProps={(value) => ({ checked: value === 'enable' })}
          >
            <Switch checkedChildren="正常" unCheckedChildren="停用" />
          </Form.Item>

          {modalType === 'view' && (
            <>
              <Divider />
              <Form.Item name="menu_id" label="菜单ID">
                <Input disabled />
              </Form.Item>
              <Form.Item name="created_at" label="创建时间">
                <Input 
                  disabled 
                  value={currentMenu?.created_at ? new Date(currentMenu.created_at).toLocaleString() : ''} 
                />
              </Form.Item>
              <Form.Item name="updated_at" label="更新时间">
                <Input 
                  disabled 
                  value={currentMenu?.updated_at ? new Date(currentMenu.updated_at).toLocaleString() : ''} 
                />
              </Form.Item>
            </>
          )}
        </Form>
      </Modal>
    </PageContainer>
  );
};

export default MenuManagement; 