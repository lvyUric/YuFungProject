import React, { useRef, useState, useEffect } from 'react';
import {
  StepsForm,
  ProFormText,
  ProFormSelect,
  ProFormDatePicker,
  ProFormDigit,
  ProFormSwitch,
  ProFormTextArea,
  ProForm,
  type ProFormInstance,
} from '@ant-design/pro-components';
import { message, Card, Typography, Space, Modal } from 'antd';
import { 
  UserOutlined, 
  TeamOutlined, 
  MoneyCollectOutlined, 
  ShoppingOutlined, 
  SettingOutlined
} from '@ant-design/icons';
import { 
  createPolicy, 
  updatePolicy, 
  type PolicyCreateRequest, 
  type PolicyInfo 
} from '@/services/policy';
import { getSystemConfigOptions, type SystemConfigInfo } from '@/services/system-config';
import { getCompanyList } from '@/services/ant-design-pro/company';

const { Title, Text } = Typography;

interface PolicyStepFormProps {
  visible: boolean;
  onVisibleChange: (visible: boolean) => void;
  onSuccess: () => void;
  initialValues?: PolicyInfo;
}

const PolicyStepForm: React.FC<PolicyStepFormProps> = ({
  visible,
  onVisibleChange,
  onSuccess,
  initialValues,
}) => {
  const formRef = useRef<ProFormInstance>(null);
  const isEdit = !!initialValues;
  
  // 下拉框选项状态
  const [hkManagerOptions, setHkManagerOptions] = useState<Array<{label: string, value: string}>>([]);
  const [referralBranchOptions, setReferralBranchOptions] = useState<Array<{label: string, value: string}>>([]);
  const [partnerOptions, setPartnerOptions] = useState<Array<{label: string, value: string}>>([]);
  const [companyOptions, setCompanyOptions] = useState<Array<{label: string, value: string}>>([]);

  // 币种选项
  const currencyOptions = [
    { label: 'USD - 美元', value: 'USD' },
    { label: 'HKD - 港币', value: 'HKD' },
    { label: 'CNY - 人民币', value: 'CNY' },
  ];

  // 缴费方式选项
  const paymentMethodOptions = [
    { label: '期缴 - 定期缴费', value: '期缴' },
    { label: '趸缴 - 一次性缴费', value: '趸缴' },
    { label: '预缴 - 预先缴费', value: '预缴' },
  ];

  // 加载下拉框选项数据
  useEffect(() => {
    const loadOptions = async () => {
      try {
        // 加载港分客户经理选项
        const hkManagerRes = await getSystemConfigOptions('hk_manager');
        if (hkManagerRes.code === 200 && hkManagerRes.data && Array.isArray(hkManagerRes.data)) {
          setHkManagerOptions(hkManagerRes.data.map((item: SystemConfigInfo) => ({
            label: item.display_name,
            value: item.config_value,
          })));
        } else {
          console.warn('港分客户经理数据格式不正确:', hkManagerRes);
          setHkManagerOptions([]);
        }

        // 加载转介分行选项
        const referralBranchRes = await getSystemConfigOptions('referral_branch');
        if (referralBranchRes.code === 200 && referralBranchRes.data && Array.isArray(referralBranchRes.data)) {
          setReferralBranchOptions(referralBranchRes.data.map((item: SystemConfigInfo) => ({
            label: item.display_name,
            value: item.config_value,
          })));
        } else {
          console.warn('转介分行数据格式不正确:', referralBranchRes);
          setReferralBranchOptions([]);
        }

        // 加载合作伙伴选项
        const partnerRes = await getSystemConfigOptions('partner');
        if (partnerRes.code === 200 && partnerRes.data && Array.isArray(partnerRes.data)) {
          setPartnerOptions(partnerRes.data.map((item: SystemConfigInfo) => ({
            label: item.display_name,
            value: item.config_value,
          })));
        } else {
          console.warn('合作伙伴数据格式不正确:', partnerRes);
          setPartnerOptions([]);
        }

        // 加载承保公司选项
        const companyRes = await getCompanyList({ page: 1, page_size: 1000 });
        if (companyRes.code === 200 && companyRes.data?.companies && Array.isArray(companyRes.data.companies)) {
          setCompanyOptions(companyRes.data.companies.map((item: any) => ({
            label: item.company_name,
            value: item.company_name,
          })));
        } else {
          console.warn('公司列表数据格式不正确:', companyRes);
          setCompanyOptions([]);
        }
      } catch (error) {
        console.error('加载下拉框选项失败:', error);
      }
    };

    if (visible) {
      loadOptions();
    }
  }, [visible]);

  // 表单提交处理
  const handleFinish = async (values: PolicyCreateRequest) => {
    try {
      if (isEdit && initialValues) {
        await updatePolicy(initialValues.policy_id, values);
        message.success('保单更新成功');
      } else {
        await createPolicy(values);
        message.success('保单创建成功');
      }
      onSuccess();
      return true;
    } catch (error) {
      message.error(isEdit ? '保单更新失败' : '保单创建失败');
      return false;
    }
  };

  // 格式化初始值用于StepsForm
  const getFormattedInitialValues = () => {
    if (!initialValues) return undefined;
    
    const formattedValues = {
      ...initialValues,
      referral_date: initialValues.referral_date ? new Date(initialValues.referral_date) : undefined,
      payment_date: initialValues.payment_date ? new Date(initialValues.payment_date) : undefined,
      effective_date: initialValues.effective_date ? new Date(initialValues.effective_date) : undefined,
      payment_pay_date: initialValues.payment_pay_date ? new Date(initialValues.payment_pay_date) : undefined,
    };
    
    console.log('PolicyStepForm - 编辑模式，设置初始值:', {
      policy_id: initialValues.policy_id,
      account_number: formattedValues.account_number,
      customer_name_cn: formattedValues.customer_name_cn,
      proposal_number: formattedValues.proposal_number,
      totalFields: Object.keys(formattedValues).length
    });
    
    return formattedValues;
  };

  return (
    <StepsForm
      key={initialValues?.policy_id || 'new'}
      formRef={formRef}
      onFinish={handleFinish}
      stepsProps={{
        size: 'default',
      }}
      stepsFormRender={(dom, submitter) => {
        return (
          <Modal
            title={
              <div style={{ textAlign: 'center', padding: '8px 0' }}>
                <Title level={3} style={{ margin: 0, color: '#1890ff' }}>
                  {isEdit ? '编辑保单' : '新建保单'}
                </Title>
                <Text type="secondary">请按步骤填写保单信息</Text>
              </div>
            }
            open={visible}
            onCancel={() => onVisibleChange(false)}
            footer={submitter}
            width={1000}
            destroyOnClose
            maskClosable={false}
          >
            <Card 
              style={{ 
                minHeight: '500px',
              }}
              bodyStyle={{ padding: '24px' }}
            >
              {dom}
            </Card>
          </Modal>
        );
      }}
    >
      {/* 第一步：基本信息 */}
      <StepsForm.StepForm
        name="basic"
        title="基本信息"
        initialValues={getFormattedInitialValues()}
        stepProps={{
          description: '填写客户和保单基本信息',
          icon: <UserOutlined />,
        }}
      >
        <div style={{ marginBottom: '24px' }}>
          <Space align="center" style={{ marginBottom: '16px' }}>
            <UserOutlined style={{ fontSize: '20px', color: '#1890ff' }} />
            <Title level={4} style={{ margin: 0 }}>客户基本信息</Title>
          </Space>
          <Text type="secondary">请填写客户的基本信息和保单标识</Text>
        </div>
        
        <ProForm.Group>
          <ProFormText
            width="md"
            name="account_number"
            label="账户号"
            placeholder="请输入客户账户号"
          />
          <ProFormText
            width="md"
            name="customer_number"
            label="客户号"
            rules={[{ required: true, message: '请输入客户号' }]}
            placeholder="请输入客户编号"
          />
        </ProForm.Group>
        
        <ProForm.Group>
          <ProFormText
            width="md"
            name="customer_name_cn"
            label="客户中文名"
            rules={[
              { required: true, message: '请输入客户中文名' },
              { min: 2, message: '姓名至少2个字符' }
            ]}
            placeholder="请输入客户中文姓名"
          />
          <ProFormText
            width="md"
            name="customer_name_en"
            label="客户英文名"
            placeholder="请输入客户英文姓名（可选）"
          />
        </ProForm.Group>
        
        <ProForm.Group>
          <ProFormText
            width="md"
            name="proposal_number"
            label="投保单号"
            rules={[{ required: true, message: '请输入投保单号' }]}
            placeholder="请输入投保单号"
          />
          <ProFormSelect
            width="md"
            name="policy_currency"
            label="保单币种"
            options={currencyOptions}
            rules={[{ required: true, message: '请选择保单币种' }]}
            placeholder="请选择保单使用的币种"
          />
        </ProForm.Group>
      </StepsForm.StepForm>

      {/* 第二步：转介信息 */}
      <StepsForm.StepForm
        name="referral"
        title="转介信息"
        initialValues={getFormattedInitialValues()}
        stepProps={{
          description: '填写合作伙伴和转介相关信息',
          icon: <TeamOutlined />,
        }}
      >
        <div style={{ marginBottom: '24px' }}>
          <Space align="center" style={{ marginBottom: '16px' }}>
            <TeamOutlined style={{ fontSize: '20px', color: '#1890ff' }} />
            <Title level={4} style={{ margin: 0 }}>合作伙伴信息</Title>
          </Space>
          <Text type="secondary">请填写转介相关的合作伙伴和经理信息</Text>
        </div>

        <ProForm.Group>
          <ProFormSelect
            width="md"
            name="partner"
            label="合作伙伴"
            options={partnerOptions}
            placeholder="请选择合作伙伴"
            fieldProps={{
              showSearch: true,
              optionFilterProp: 'label',
            }}
            rules={[{ required: true, message: '请选择合作伙伴' }]}
          />
          <ProFormText
            width="md"
            name="referral_code"
            label="转介编号"
            placeholder="请输入转介编号"
          />
        </ProForm.Group>
        
        <ProForm.Group>
          <ProFormSelect
            width="md"
            name="hk_manager"
            label="港分客户经理"
            options={hkManagerOptions}
            placeholder="请选择港分客户经理"
            fieldProps={{
              showSearch: true,
              optionFilterProp: 'label',
            }}
            rules={[{ required: true, message: '请选择港分客户经理' }]}
          />
          <ProFormText
            width="md"
            name="referral_pm"
            label="转介理财经理"
            placeholder="请输入转介理财经理姓名"
          />
        </ProForm.Group>
        
        <ProForm.Group>
          <ProFormSelect
            width="md"
            name="referral_branch"
            label="转介分行"
            options={referralBranchOptions}
            placeholder="请选择转介分行"
            fieldProps={{
              showSearch: true,
              optionFilterProp: 'label',
            }}
            rules={[{ required: true, message: '请选择转介分行' }]}
          />
          <ProFormText
            width="md"
            name="referral_sub_branch"
            label="转介支行"
            placeholder="请输入转介支行名称"
          />
        </ProForm.Group>
        
        <ProForm.Group>
          <ProFormDatePicker
            width="md"
            name="referral_date"
            label="转介日期"
            placeholder="请选择转介日期"
          />
        </ProForm.Group>
      </StepsForm.StepForm>

      {/* 第三步：缴费信息 */}
      <StepsForm.StepForm
        name="payment"
        title="缴费信息"
        initialValues={getFormattedInitialValues()}
        stepProps={{
          description: '填写缴费相关信息',
          icon: <MoneyCollectOutlined />,
        }}
      >
        <div style={{ marginBottom: '24px' }}>
          <Space align="center" style={{ marginBottom: '16px' }}>
            <MoneyCollectOutlined style={{ fontSize: '20px', color: '#1890ff' }} />
            <Title level={4} style={{ margin: 0 }}>缴费信息</Title>
          </Space>
          <Text type="secondary">请填写保费缴纳和费用相关信息</Text>
        </div>

        <ProForm.Group>
          <ProFormDatePicker
            width="md"
            name="payment_date"
            label="缴费日期"
            placeholder="请选择缴费日期"
          />
          <ProFormDatePicker
            width="md"
            name="effective_date"
            label="生效日期"
            placeholder="请选择生效日期"
          />
        </ProForm.Group>

        <ProForm.Group>
          <ProFormSelect
            width="md"
            name="payment_method"
            label="缴费方式"
            options={paymentMethodOptions}
            placeholder="请选择缴费方式"
          />
          <ProFormDigit
            width="md"
            name="payment_years"
            label="缴费年期"
            min={0}
            placeholder="请输入缴费年期"
          />
        </ProForm.Group>
        
        <ProForm.Group>
          <ProFormDigit
            width="md"
            name="actual_premium"
            label="实际缴纳保费"
            min={0}
            fieldProps={{ precision: 2 }}
            placeholder="请输入实际缴纳保费"
          />
          <ProFormDigit
            width="md"
            name="aum"
            label="AUM"
            min={0}
            fieldProps={{ precision: 2 }}
            placeholder="请输入AUM"
          />
        </ProForm.Group>
        
        <ProForm.Group>
          <ProFormDigit
            width="md"
            name="referral_rate"
            label="转介费率(%)"
            min={0}
            max={100}
            fieldProps={{ precision: 2 }}
            placeholder="请输入转介费率"
          />
          <ProFormDigit
            width="md"
            name="exchange_rate"
            label="汇率"
            min={0}
            fieldProps={{ precision: 4 }}
            placeholder="请输入汇率"
          />
        </ProForm.Group>
        
        <ProForm.Group>
          <ProFormDigit
            width="md"
            name="expected_fee"
            label="预计转介费"
            min={0}
            fieldProps={{ precision: 2 }}
            placeholder="请输入预计转介费"
          />
        </ProForm.Group>
      </StepsForm.StepForm>

      {/* 第四步：产品信息 */}
      <StepsForm.StepForm
        name="product"
        title="产品信息"
        initialValues={getFormattedInitialValues()}
        stepProps={{
          description: '填写保险产品相关信息',
          icon: <ShoppingOutlined />,
        }}
      >
        <div style={{ marginBottom: '24px' }}>
          <Space align="center" style={{ marginBottom: '16px' }}>
            <ShoppingOutlined style={{ fontSize: '20px', color: '#1890ff' }} />
            <Title level={4} style={{ margin: 0 }}>产品信息</Title>
          </Space>
          <Text type="secondary">请填写保险产品的详细信息</Text>
        </div>

        <ProForm.Group>
          <ProFormSelect
            width="md"
            name="insurance_company"
            label="承保公司"
            options={companyOptions}
            placeholder="请选择承保公司"
            fieldProps={{
              showSearch: true,
              optionFilterProp: 'label',
            }}
            rules={[{ required: true, message: '请选择承保公司' }]}
          />
          <ProFormText
            width="md"
            name="product_name"
            label="保险产品名称"
            rules={[{ required: true, message: '请输入保险产品名称' }]}
            placeholder="请输入保险产品名称"
          />
        </ProForm.Group>
        
        <ProForm.Group>
          <ProFormText
            width="md"
            name="product_type"
            label="产品类型"
            rules={[{ required: true, message: '请输入产品类型' }]}
            placeholder="请输入产品类型"
          />
        </ProForm.Group>
        
        <ProFormTextArea
          name="remark"
          label="备注说明"
          placeholder="请输入备注信息（可选）"
          fieldProps={{
            rows: 4,
            showCount: true,
            maxLength: 500,
          }}
        />
      </StepsForm.StepForm>

      {/* 第五步：状态设置 */}
      <StepsForm.StepForm
        name="status"
        title="状态设置"
        initialValues={getFormattedInitialValues()}
        stepProps={{
          description: '设置保单相关状态',
          icon: <SettingOutlined />,
        }}
      >
        <div style={{ marginBottom: '24px' }}>
          <Space align="center" style={{ marginBottom: '16px' }}>
            <SettingOutlined style={{ fontSize: '20px', color: '#1890ff' }} />
            <Title level={4} style={{ margin: 0 }}>状态设置</Title>
          </Space>
          <Text type="secondary">请设置保单的相关状态信息</Text>
        </div>

        <ProForm.Group>
          <ProFormSwitch
            name="is_surrendered"
            label="是否退保"
            tooltip="签单后是否已经退保"
          />
          <ProFormSwitch
            name="past_cooling_period"
            label="是否已过冷静期"
            tooltip="保单是否已过冷静期"
          />
        </ProForm.Group>

        <ProForm.Group>
          <ProFormSwitch
            name="is_paid_commission"
            label="是否支付佣金"
            tooltip="是否已支付相关佣金"
          />
          <ProFormSwitch
            name="is_employee"
            label="是否员工"
            tooltip="客户是否为公司员工"
          />
        </ProForm.Group>
      </StepsForm.StepForm>
    </StepsForm>
  );
};

export default PolicyStepForm; 