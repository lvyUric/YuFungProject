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
  Divider,
  Table,
  Tooltip,
  Tree,
  Row,
  Col,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  ReloadOutlined,
  TeamOutlined,
  SearchOutlined,
} from '@ant-design/icons';
import { useRequest } from 'ahooks';
import type { ColumnsType } from 'antd/es/table';
import type { DataNode } from 'antd/es/tree';
import {
  getRoleList,
  createRole,
  updateRole,
  deleteRole,
  batchUpdateRoleStatus,
  getRoleStats,
  type RoleInfo,
  type RoleStatsResponse,
} from '../../services/role';
import { getMenuTree, type MenuInfo } from '../../services/menu';

const { Option } = Select;

// 数据权限范围选项
const dataScopeOptions = [
  { label: '全部数据权限', value: 'all' },
  { label: '本公司数据权限', value: 'company' },
  { label: '仅本人数据权限', value: 'self' },
];

// 角色管理页面
const RoleManagement: React.FC = () => {
  const [form] = Form.useForm();
  const [searchForm] = Form.useForm();
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [modalVisible, setModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'create' | 'edit' | 'view'>('create');
  const [currentRole, setCurrentRole] = useState<RoleInfo | null>(null);
  const [permissionModalVisible, setPermissionModalVisible] = useState(false);
  const [menuTreeData, setMenuTreeData] = useState<MenuInfo[]>([]);
  const [checkedMenuKeys, setCheckedMenuKeys] = useState<React.Key[]>([]);
  const [expandedMenuKeys, setExpandedMenuKeys] = useState<React.Key[]>([]);
  const [filteredData, setFilteredData] = useState<RoleInfo[]>([]);
  const [searchValues, setSearchValues] = useState<{
    role_name?: string;
    status?: string;
  }>({});
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);

  // 获取角色统计信息
  const { data: statsData, refresh: refreshStats } = useRequest(
    () => getRoleStats({ company_id: '' }),
    {}
  );

  // 获取角色列表
  const { data: roleData, loading, refresh } = useRequest(
    () => getRoleList({ page: 1, page_size: 100 }),
    {
      onSuccess: (data) => {
        if (data.code === 200) {
          const roles = data.data?.roles || [];
          setFilteredData(roles);
        }
      },
    }
  );

  // 获取菜单树数据
  const { data: menuData } = useRequest(
    () => getMenuTree({}),
    {
      onSuccess: (data) => {
        if (data.code === 200) {
          const menus = data.data || [];
          setMenuTreeData(menus);
          // 默认展开所有节点
          const allKeys = getAllMenuKeys(menus);
          setExpandedMenuKeys(allKeys);
        }
      },
    }
  );

  // 获取所有菜单节点的key（用于展开）
  const getAllMenuKeys = (menus: MenuInfo[]): string[] => {
    let keys: string[] = [];
    menus.forEach(menu => {
      keys.push(menu.menu_id);
      if (menu.children && menu.children.length > 0) {
        keys = keys.concat(getAllMenuKeys(menu.children));
      }
    });
    return keys;
  };

  // 构建菜单树数据结构
  const buildMenuTreeData = (menus: MenuInfo[]): DataNode[] => {
    return menus.map(menu => ({
      title: menu.menu_name,
      key: menu.menu_id,
      children: menu.children && menu.children.length > 0 ? buildMenuTreeData(menu.children) : undefined,
    }));
  };

  // 搜索过滤角色
  const filterRoles = (roles: RoleInfo[], searchValues: any): RoleInfo[] => {
    return roles.filter(role => {
      let matches = true;
      
      if (searchValues.role_name && !role.role_name.toLowerCase().includes(searchValues.role_name.toLowerCase())) {
        matches = false;
      }
      
      if (searchValues.status && role.status !== searchValues.status) {
        matches = false;
      }
      
      return matches;
    });
  };

  // 处理搜索
  const handleSearch = () => {
    const values = searchForm.getFieldsValue();
    setSearchValues(values);
    
    const roles = roleData?.data?.roles || [];
    if (!values.role_name && !values.status) {
      setFilteredData(roles);
    } else {
      const filtered = filterRoles(roles, values);
      setFilteredData(filtered);
    }
  };

  // 重置搜索
  const handleReset = () => {
    searchForm.resetFields();
    setSearchValues({});
    setFilteredData(roleData?.data?.roles || []);
  };

  // 表格列定义
  const columns: ColumnsType<RoleInfo> = [
    {
      title: '角色名称',
      dataIndex: 'role_name',
      key: 'role_name',
      width: 150,
    },
    {
      title: '权限字符',
      dataIndex: 'role_key',
      key: 'role_key',
      width: 150,
      render: (text: string) => <code style={{ fontSize: '12px' }}>{text}</code>,
    },
    {
      title: '显示顺序',
      dataIndex: 'sort_order',
      key: 'sort_order',
      width: 100,
      align: 'center',
    },
    {
      title: '数据权限',
      dataIndex: 'data_scope',
      key: 'data_scope',
      width: 120,
      render: (scope: string) => {
        const config = {
          all: { color: 'blue', text: '全部数据' },
          company: { color: 'green', text: '本公司数据' },
          self: { color: 'orange', text: '仅本人数据' },
        };
        return <Tag color={config[scope as keyof typeof config]?.color}>
          {config[scope as keyof typeof config]?.text}
        </Tag>;
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (status: string) => (
        <Tag color={status === 'enable' ? 'green' : 'red'}>
          {status === 'enable' ? '正常' : '停用'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 150,
      render: (text: string) => new Date(text).toLocaleString(),
    },
    {
      title: '备注',
      dataIndex: 'remark',
      key: 'remark',
      width: 150,
      ellipsis: true,
      render: (text: string) => text || '-',
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="分配权限">
            <Button
              type="link"
              size="small"
              onClick={() => handleAssignPermission(record)}
            >
              分配权限
            </Button>
          </Tooltip>
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
              title="确定要删除这个角色吗？"
              description="删除后将无法恢复，且会影响相关用户权限"
              onConfirm={() => handleDelete(record.role_id)}
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

  // 处理分配权限
  const handleAssignPermission = (record: RoleInfo) => {
    setCurrentRole(record);
    setCheckedMenuKeys(record.menu_ids || []);
    setPermissionModalVisible(true);
  };

  // 处理查看
  const handleView = (record: RoleInfo) => {
    setCurrentRole(record);
    setModalType('view');
    setModalVisible(true);
    form.setFieldsValue({
      ...record,
    });
  };

  // 处理编辑
  const handleEdit = (record: RoleInfo) => {
    setCurrentRole(record);
    setModalType('edit');
    setModalVisible(true);
    form.setFieldsValue({
      ...record,
    });
  };

  // 处理新建
  const handleCreate = () => {
    setCurrentRole(null);
    setModalType('create');
    setModalVisible(true);
    form.resetFields();
    form.setFieldsValue({
      data_scope: 'company',
      status: 'enable',
      sort_order: 1,
    });
  };

  // 处理删除
  const handleDelete = async (roleId: string) => {
    try {
      const response = await deleteRole(roleId);
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
      message.warning('请选择要操作的角色');
      return;
    }

    try {
      const response = await batchUpdateRoleStatus({
        role_ids: selectedRowKeys.map(key => String(key)),
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
        response = await createRole(values);
      } else if (modalType === 'edit' && currentRole) {
        response = await updateRole(currentRole.role_id, values);
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

  // 提交权限分配
  const handlePermissionSubmit = async () => {
    if (!currentRole) return;

    try {
      const response = await updateRole(currentRole.role_id, {
        ...currentRole,
        menu_ids: checkedMenuKeys.map(key => String(key)),
      });

      if (response.code === 200) {
        message.success('权限分配成功');
        setPermissionModalVisible(false);
        refresh();
      } else {
        message.error(response.message);
      }
    } catch (error) {
      message.error('权限分配失败');
    }
  };

  // 统计卡片
  const StatsCards = () => {
    const stats = (statsData as any)?.data as RoleStatsResponse | undefined;
    return (
      <div style={{ marginBottom: 16 }}>
        <div style={{ display: 'flex', gap: 16 }}>
          <Card size="small" style={{ flex: 1 }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 'bold', color: '#1890ff' }}>
                {stats?.total_roles || 0}
              </div>
              <div>总角色数</div>
            </div>
          </Card>
          <Card size="small" style={{ flex: 1 }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 'bold', color: '#52c41a' }}>
                {stats?.enabled_roles || 0}
              </div>
              <div>启用角色</div>
            </div>
          </Card>
          <Card size="small" style={{ flex: 1 }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 'bold', color: '#f5222d' }}>
                {stats?.disabled_roles || 0}
              </div>
              <div>禁用角色</div>
            </div>
          </Card>
          <Card size="small" style={{ flex: 1 }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 'bold', color: '#722ed1' }}>
                {stats?.platform_roles || 0}
              </div>
              <div>平台角色</div>
            </div>
          </Card>
          <Card size="small" style={{ flex: 1 }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: 24, fontWeight: 'bold', color: '#fa8c16' }}>
                {stats?.company_roles || 0}
              </div>
              <div>公司角色</div>
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
          <Form.Item name="role_name" label="角色名称">
            <Input placeholder="请输入角色名称" allowClear style={{ width: 200 }} />
          </Form.Item>
          <Form.Item name="status" label="角色状态">
            <Select placeholder="角色状态" allowClear style={{ width: 120 }}>
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

        {/* 角色表格 */}
        <Table
          columns={columns}
          dataSource={filteredData}
          rowKey="role_id"
          loading={loading}
          onChange={(pagination, filters, sorter) => {
            // 处理分页变化
            console.log('分页变化:', pagination);
            setCurrentPage(pagination.current || 1);
            setPageSize(pagination.pageSize || 20);
          }}
          pagination={{
            current: currentPage,
            pageSize: pageSize,
            showQuickJumper: true,
            showSizeChanger: true,
            pageSizeOptions: [10, 20, 50, 100],
            showTotal: (total, range) =>
              `第 ${range[0]}-${range[1]} 条/总共 ${total} 条`,
          }}
          rowSelection={{
            selectedRowKeys,
            onChange: setSelectedRowKeys,
          }}
          scroll={{ x: 1200 }}
          size="small"
        />
      </Card>

      {/* 角色详情/编辑/新建模态框 */}
      <Modal
        title={
          modalType === 'create'
            ? '新增角色'
            : modalType === 'edit'
            ? '修改角色'
            : '角色详情'
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
            name="role_name"
            label="角色名称"
            rules={[{ required: true, message: '请输入角色名称' }]}
          >
            <Input placeholder="请输入角色名称" />
          </Form.Item>

          <Form.Item
            name="role_key"
            label="权限字符"
            rules={[{ required: true, message: '请输入权限字符' }]}
            tooltip="控制器中定义的权限字符，如：system:user:list"
          >
            <Input placeholder="请输入权限字符" />
          </Form.Item>

          <Form.Item
            name="sort_order"
            label="显示排序"
            rules={[{ required: true, message: '请输入显示排序' }]}
          >
            <InputNumber
              min={1}
              max={9999}
              placeholder="请输入显示排序"
              style={{ width: '100%' }}
            />
          </Form.Item>

          <Form.Item
            name="data_scope"
            label="数据权限"
            rules={[{ required: true, message: '请选择数据权限' }]}
          >
            <Select placeholder="请选择数据权限">
              {dataScopeOptions.map(option => (
                <Option key={option.value} value={option.value}>
                  {option.label}
                </Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="status"
            label="角色状态"
            valuePropName="checked"
            getValueFromEvent={(checked) => checked ? 'enable' : 'disable'}
            getValueProps={(value) => ({ checked: value === 'enable' })}
          >
            <Switch checkedChildren="正常" unCheckedChildren="停用" />
          </Form.Item>

          <Form.Item
            name="remark"
            label="备注"
          >
            <Input.TextArea rows={3} placeholder="请输入备注信息" />
          </Form.Item>

          {modalType === 'view' && (
            <>
              <Divider />
              <Form.Item name="role_id" label="角色ID">
                <Input disabled />
              </Form.Item>
              <Form.Item name="created_at" label="创建时间">
                <Input 
                  disabled 
                  value={currentRole?.created_at ? new Date(currentRole.created_at).toLocaleString() : ''} 
                />
              </Form.Item>
              <Form.Item name="updated_at" label="更新时间">
                <Input 
                  disabled 
                  value={currentRole?.updated_at ? new Date(currentRole.updated_at).toLocaleString() : ''} 
                />
              </Form.Item>
            </>
          )}
        </Form>
      </Modal>

      {/* 权限分配模态框 */}
      <Modal
        title={`分配权限 - ${currentRole?.role_name}`}
        open={permissionModalVisible}
        onCancel={() => setPermissionModalVisible(false)}
        onOk={handlePermissionSubmit}
        width={800}
        destroyOnClose
      >
        <Row gutter={16}>
          <Col span={12}>
            <Card size="small" title="菜单权限">
              <Tree
                checkable
                showLine
                expandedKeys={expandedMenuKeys}
                onExpand={setExpandedMenuKeys}
                checkedKeys={checkedMenuKeys}
                onCheck={(keys) => setCheckedMenuKeys(keys as React.Key[])}
                treeData={buildMenuTreeData(menuTreeData)}
                height={400}
              />
            </Card>
          </Col>
          <Col span={12}>
            <Card size="small" title="角色信息">
              <div style={{ padding: '16px 0' }}>
                <p><strong>角色名称：</strong>{currentRole?.role_name}</p>
                <p><strong>权限字符：</strong><code>{currentRole?.role_key}</code></p>
                <p><strong>数据权限：</strong>
                  <Tag color="blue">
                    {dataScopeOptions.find(opt => opt.value === currentRole?.data_scope)?.label}
                  </Tag>
                </p>
                <p><strong>状态：</strong>
                  <Tag color={currentRole?.status === 'enable' ? 'green' : 'red'}>
                    {currentRole?.status === 'enable' ? '正常' : '停用'}
                  </Tag>
                </p>
                <p><strong>备注：</strong>{currentRole?.remark || '无'}</p>
              </div>
              <Divider />
              <div>
                <p><strong>已选择权限：</strong></p>
                <div style={{ maxHeight: 300, overflow: 'auto' }}>
                  {checkedMenuKeys.length === 0 ? (
                    <p style={{ color: '#999' }}>未选择任何权限</p>
                  ) : (
                    checkedMenuKeys.map(key => (
                      <Tag key={key} style={{ margin: '2px' }}>
                        {key}
                      </Tag>
                    ))
                  )}
                </div>
              </div>
            </Card>
          </Col>
        </Row>
      </Modal>
    </PageContainer>
  );
};

export default RoleManagement; 