-- 医院应急 HIS 数据库初始化脚本
-- 数据库：his_emergency
-- 适用版本：MySQL 8.4.11
-- 用途：创建数据库、表结构、索引及默认系统配置

CREATE DATABASE IF NOT EXISTS his_emergency 
    CHARACTER SET utf8mb4 
    COLLATE utf8mb4_unicode_ci;

USE his_emergency;

-- ============================================
-- 一、同步表（从主HIS同步，应急期间只读）
-- ============================================

-- 1.1 患者主索引 (sync_patient)
-- 来源: G0003 查看患者主索引信息 + G0042 挂号信息 + G0071 就诊记录
CREATE TABLE IF NOT EXISTS sync_patient (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '自增ID',
    patient_id          VARCHAR(50) NOT NULL COMMENT '患者ID，如 M001521974',
    name                VARCHAR(100) NOT NULL COMMENT '姓名',
    gender              VARCHAR(10) COMMENT '性别: 男/女',
    birth_date          DATE COMMENT '出生日期',
    id_card_no          VARCHAR(50) COMMENT '身份证号',
    phone               VARCHAR(50) COMMENT '手机号',
    charge_type         VARCHAR(50) COMMENT '费别: 医疗保险/公费医疗/自费等',
    identity_type       VARCHAR(50) COMMENT '身份类型: 一般人员/本院职员等',
    address             VARCHAR(500) COMMENT '地址',
    emergency_contact   VARCHAR(100) COMMENT '紧急联系人',
    emergency_phone     VARCHAR(50) COMMENT '紧急联系电话',
    nation              VARCHAR(50) COMMENT '民族',
    card_no             VARCHAR(50) COMMENT '就诊卡号',
    last_sync_time      DATETIME COMMENT '最后同步时间',
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_patient_id (patient_id),
    KEY idx_id_card (id_card_no),
    KEY idx_name (name),
    KEY idx_phone (phone),
    KEY idx_sync_time (last_sync_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-患者主索引';

-- 1.2 科室字典 (sync_dept)
-- 来源: G0013 查科室信息
CREATE TABLE IF NOT EXISTS sync_dept (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    dept_code           VARCHAR(50) NOT NULL COMMENT '科室编码',
    dept_name           VARCHAR(100) NOT NULL COMMENT '科室名称',
    dept_short_name     VARCHAR(100) COMMENT '科室简称',
    dept_type           VARCHAR(50) COMMENT '科室类型',
    is_outpatient       TINYINT DEFAULT 1 COMMENT '是否门诊: 1是 0否',
    is_inpatient        TINYINT DEFAULT 0 COMMENT '是否住院: 1是 0否',
    is_emergency        TINYINT DEFAULT 0 COMMENT '是否急诊: 1是 0否',
    status              TINYINT DEFAULT 1 COMMENT '状态: 1启用 0停用',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_dept_code (dept_code),
    KEY idx_dept_name (dept_name),
    KEY idx_outpatient (is_outpatient, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-科室字典';

-- 1.3 医生字典 (sync_doctor)
-- 来源: G0012 查医生信息 + G0076 门诊护士获取医生信息
CREATE TABLE IF NOT EXISTS sync_doctor (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    doctor_code         VARCHAR(50) NOT NULL COMMENT '医生工号/DB_USER',
    name                VARCHAR(100) NOT NULL COMMENT '姓名',
    title               VARCHAR(50) COMMENT '职称: 主任医师/副主任医师/主治医师/医师',
    title_code          VARCHAR(50) COMMENT '职称编码',
    dept_code           VARCHAR(50) COMMENT '所属科室编码',
    specialty           VARCHAR(200) COMMENT '专长',
    room                VARCHAR(100) COMMENT '诊室',
    room_code           VARCHAR(50) COMMENT '诊室编码',
    phone               VARCHAR(50) COMMENT '联系电话',
    job                 VARCHAR(50) COMMENT '职务: 医生/护士等',
    ca_code             VARCHAR(100) COMMENT 'CA编码',
    status              TINYINT DEFAULT 1 COMMENT '状态: 1启用 0停用',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_doctor_code (doctor_code),
    KEY idx_dept (dept_code),
    KEY idx_name (name),
    KEY idx_title (title)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-医生字典';

-- 1.4 用户字典 (sync_user)
-- 来源: G0038 查用户信息
-- 用途: 应急HIS登录认证
CREATE TABLE IF NOT EXISTS sync_user (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    login_name          VARCHAR(50) NOT NULL COMMENT '登录名/用户唯一标识',
    name                VARCHAR(100) NOT NULL COMMENT '姓名',
    gender              VARCHAR(10) COMMENT '性别',
    dept_code           VARCHAR(50) COMMENT '所属科室编码',
    ca_code             VARCHAR(100) COMMENT 'CA编码',
    student_flag        TINYINT DEFAULT 0 COMMENT '学生标识',
    teacher_flag        TINYINT DEFAULT 0 COMMENT '老师标识',
    hire_date           DATE COMMENT '入职时间',
    status              TINYINT DEFAULT 1 COMMENT '状态: 1启用 0停用',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_login_name (login_name),
    KEY idx_dept (dept_code),
    KEY idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-用户字典';

-- 1.5 药品字典 (sync_drug)
-- 来源: G0062 药品字典数据
CREATE TABLE IF NOT EXISTS sync_drug (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    drug_code           VARCHAR(50) NOT NULL COMMENT '药品编码',
    drug_name           VARCHAR(200) NOT NULL COMMENT '药品名称',
    drug_spec           VARCHAR(200) COMMENT '规格',
    drug_form           VARCHAR(50) COMMENT '剂型: 片剂/注射剂/胶囊等',
    units               VARCHAR(50) COMMENT '单位',
    normal_name         VARCHAR(200) COMMENT '通用名',
    dose_per_unit       DECIMAL(10,4) COMMENT '单次剂量',
    dose_units          VARCHAR(50) COMMENT '剂量单位',
    toxi_property       VARCHAR(50) COMMENT '药品类型: 普通药品/精神药品/麻醉药品/抗生素',
    input_code          VARCHAR(50) COMMENT '拼音码',
    drug_indicator      TINYINT DEFAULT 1 COMMENT '使用标志: 1启用 0停用',
    status              TINYINT DEFAULT 1 COMMENT '状态',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_drug_code (drug_code),
    KEY idx_input_code (input_code),
    KEY idx_name (drug_name),
    KEY idx_toxi (toxi_property)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-药品字典';

-- 1.6 药品别名字典 (sync_drug_alias)
-- 来源: G0057 获取药品别名字典信息
-- 新增表
CREATE TABLE IF NOT EXISTS sync_drug_alias (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    drug_code           VARCHAR(50) NOT NULL COMMENT '药品编码',
    alias               VARCHAR(200) COMMENT '药品别名',
    alias_type          VARCHAR(50) COMMENT '别名类型: 药品通用名/商品名等',
    class_code          VARCHAR(10) COMMENT '项目类别代码: A西药/B中成药/C草药等',
    class_name          VARCHAR(50) COMMENT '项目类别名称',
    stop_time           DATETIME COMMENT '停用时间',
    input_code          VARCHAR(50) COMMENT '拼音码',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    KEY idx_drug_code (drug_code),
    KEY idx_alias (alias),
    KEY idx_input_code (input_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-药品别名字典';

-- 1.7 药房库存 (sync_drug_stock)
-- 来源: G0061 药房库存数据
-- 说明: 每次同步时全量刷新（先删后插）
CREATE TABLE IF NOT EXISTS sync_drug_stock (
    id                      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    storage                 VARCHAR(50) NOT NULL COMMENT '药房编码',
    storage_name            VARCHAR(100) COMMENT '药房名称',
    drug_code               VARCHAR(50) NOT NULL COMMENT '药品编码',
    drug_spec               VARCHAR(200) COMMENT '规格',
    units                   VARCHAR(50) COMMENT '单位',
    batch_no                VARCHAR(100) COMMENT '批次号',
    expire_date             DATE COMMENT '有效期',
    firm_id                 VARCHAR(200) COMMENT '生产厂家',
    purchase_price          DECIMAL(10,4) COMMENT '零售价',
    discount                DECIMAL(10,2) COMMENT '折扣',
    package_spec            VARCHAR(200) COMMENT '包装规格',
    quantity                DECIMAL(10,2) COMMENT '库存数量',
    package_units           VARCHAR(50) COMMENT '包装单位',
    sub_package_1           DECIMAL(10,2) COMMENT '最小包装数量',
    sub_package_units_1     VARCHAR(50) COMMENT '最小包装单位',
    sub_package_spec_1      VARCHAR(200) COMMENT '最小包装规格',
    sub_storage             VARCHAR(100) COMMENT '二级库',
    location                VARCHAR(200) COMMENT '货位',
    document_no             VARCHAR(100) COMMENT '单据号',
    supply_indicator        TINYINT DEFAULT 1 COMMENT '可供标志: 1可供 0不可供',
    outp_supply_indicator   TINYINT DEFAULT 1 COMMENT '门诊可供标志: 1可供 0不可供',
    hlw_supply_indicator    TINYINT COMMENT '互联网可供标志',
    last_sync_time          DATETIME,
    created_at              DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    KEY idx_storage_drug (storage, drug_code),
    KEY idx_drug_code (drug_code),
    KEY idx_expire (expire_date),
    KEY idx_supply (supply_indicator, outp_supply_indicator)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-药房库存';

-- 1.8 收费价表 (sync_price_item)
-- 来源: G0005 获取麻醉模版物价项目服务
CREATE TABLE IF NOT EXISTS sync_price_item (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    item_code           VARCHAR(50) NOT NULL COMMENT '项目代码',
    item_name           VARCHAR(200) NOT NULL COMMENT '项目名称',
    item_class          VARCHAR(10) COMMENT '项目类别: A药品/E检查/C化验/F手术等',
    normal_name         VARCHAR(200) COMMENT '通用名',
    item_spec           VARCHAR(200) COMMENT '规格',
    units               VARCHAR(50) COMMENT '单位',
    price               DECIMAL(10,4) NOT NULL COMMENT '单价',
    prefer_price        DECIMAL(10,4) COMMENT '优惠单价',
    foreigner_price     DECIMAL(10,4) COMMENT '外宾价格',
    class_on_mr         VARCHAR(100) COMMENT '病案首页费用分类',
    status              TINYINT DEFAULT 1 COMMENT '状态: 1启用 0停用',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_item_code (item_code),
    KEY idx_name (item_name),
    KEY idx_class (item_class)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-收费价表';

-- 1.9 诊断字典 (sync_diagnosis)
-- 来源: G0040 获取诊断字典信息
CREATE TABLE IF NOT EXISTS sync_diagnosis (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    diagnose_code       VARCHAR(50) NOT NULL COMMENT '诊断编码(ICD)',
    diagnose_name       VARCHAR(200) NOT NULL COMMENT '诊断名称',
    input_code          VARCHAR(50) COMMENT '拼音码',
    yb_reason_code_v2   VARCHAR(50) COMMENT '医保诊断编码',
    yb_reason_name_v2   VARCHAR(200) COMMENT '医保诊断名称',
    status              TINYINT DEFAULT 1 COMMENT '状态',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_diag_code (diagnose_code),
    KEY idx_input_code (input_code),
    KEY idx_name (diagnose_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-诊断字典';

-- 1.10 医生排班 (sync_schedule)
-- 来源: G0078 门诊护士获取号别排班
-- 说明: 按日期同步，每天重新拉取
CREATE TABLE IF NOT EXISTS sync_schedule (
    id                      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    clinic_date             DATE NOT NULL COMMENT '就诊日期',
    clinic_dept             VARCHAR(50) NOT NULL COMMENT '科室编码',
    dept_name               VARCHAR(100) COMMENT '科室名称',
    clinic_label            VARCHAR(100) NOT NULL COMMENT '号别',
    time_desc               VARCHAR(50) COMMENT '午别: 上午/下午',
    doctor_id               VARCHAR(50) COMMENT '医生工号',
    doctor_name             VARCHAR(100) COMMENT '医生姓名',
    registration_limits     INT DEFAULT 0 COMMENT '限号数',
    registration_num        INT DEFAULT 0 COMMENT '已挂号数',
    regist_price            DECIMAL(10,2) COMMENT '挂号费',
    clinic_type             VARCHAR(50) COMMENT '号类: 普通号/专家号/特需号',
    states                  VARCHAR(50) COMMENT '状态: 正常/停诊',
    last_sync_time          DATETIME,
    created_at              DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    KEY idx_date_dept (clinic_date, clinic_dept),
    KEY idx_doctor (doctor_id, clinic_date),
    KEY idx_label (clinic_label, clinic_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-医生排班';

-- 1.11 通用字典 (sync_common_dict)
-- 来源: G0039 获取通用字典数据
CREATE TABLE IF NOT EXISTS sync_common_dict (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    dict_type           VARCHAR(50) NOT NULL COMMENT '字典类型: nation/ChargeType/MaritalStatus/Occupation',
    dict_name           VARCHAR(100) COMMENT '字典名称',
    dict_key            VARCHAR(50) NOT NULL COMMENT '字典值编码',
    dict_value          VARCHAR(100) COMMENT '字典值',
    input_code          VARCHAR(50) COMMENT '拼音码',
    index_no            INT COMMENT '排序序号',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_type_key (dict_type, dict_key),
    KEY idx_type (dict_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-通用字典';

-- 1.12 手术字典 (sync_surgery_dict)
-- 来源: G0014 查手术信息
CREATE TABLE IF NOT EXISTS sync_surgery_dict (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    oper_code           VARCHAR(50) NOT NULL COMMENT '手术编码',
    oper_name           VARCHAR(200) NOT NULL COMMENT '手术名称',
    oper_level          VARCHAR(10) COMMENT '手术级别: 01/02/03/04',
    oper_clinic         TINYINT DEFAULT 0 COMMENT '是否门诊手术: 1是 0否',
    input_code          VARCHAR(50) COMMENT '拼音码',
    status              TINYINT DEFAULT 1 COMMENT '状态',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_oper_code (oper_code),
    KEY idx_name (oper_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-手术字典';

-- 1.13 药品用法字典 (sync_drug_usage)
-- 来源: G0055 获取药品用法字典信息
CREATE TABLE IF NOT EXISTS sync_drug_usage (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    usage_name          VARCHAR(100) NOT NULL COMMENT '途径名称: 口服/静脉注射等',
    drug_usage_code     VARCHAR(50) COMMENT '途径代码',
    input_code          VARCHAR(50) COMMENT '拼音码',
    optional_data       VARCHAR(100) COMMENT '输液对应途径',
    index_no            INT COMMENT '排序编号',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_usage_code (drug_usage_code),
    KEY idx_name (usage_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-药品用法字典';

-- 1.14 检查模板字典 (sync_exam_template)
-- 来源: G0048 同步HIS检查模板字典信息(门诊)
-- 新增表
CREATE TABLE IF NOT EXISTS sync_exam_template (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    serial_no           VARCHAR(50) NOT NULL COMMENT '模板编号',
    template_name       VARCHAR(200) COMMENT '模板名称',
    performed_by        VARCHAR(50) COMMENT '执行科室',
    pattern_sub_class   VARCHAR(100) COMMENT '模板大类',
    pattern_sub_class_2 VARCHAR(100) COMMENT '模板子类',
    exam_position       VARCHAR(200) COMMENT '检查部位',
    exam_method         VARCHAR(200) COMMENT '检查方法',
    item_no             INT COMMENT '模板序号',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_serial_no (serial_no),
    KEY idx_name (template_name),
    KEY idx_performed (performed_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-检查模板字典';

-- 1.15 检查模板明细 (sync_exam_template_item)
-- 来源: G0048 检查模板明细
-- 新增表
CREATE TABLE IF NOT EXISTS sync_exam_template_item (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    template_id         BIGINT UNSIGNED NOT NULL COMMENT '模板ID',
    item_no             INT COMMENT '明细序号',
    item_code           VARCHAR(50) COMMENT '项目编码',
    item_name           VARCHAR(200) COMMENT '项目名称',
    package_spec        VARCHAR(200) COMMENT '包装规格',
    package_units       VARCHAR(50) COMMENT '包装单位',
    firm_id             VARCHAR(200) COMMENT '厂家',
    quantity            DECIMAL(10,2) COMMENT '数量',
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    KEY idx_template (template_id),
    KEY idx_item_code (item_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-检查模板明细';

-- 1.16 化验模板字典 (sync_lab_template)
-- 来源: G0049 同步HIS化验模板字典信息(门诊)
-- 新增表
CREATE TABLE IF NOT EXISTS sync_lab_template (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    serial_no           VARCHAR(50) NOT NULL COMMENT '模板编号',
    template_name       VARCHAR(200) COMMENT '模板名称',
    performed_by        VARCHAR(50) COMMENT '执行科室',
    item_no             INT COMMENT '模板序号',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_serial_no (serial_no),
    KEY idx_name (template_name),
    KEY idx_performed (performed_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-化验模板字典';

-- 1.17 化验模板明细 (sync_lab_template_item)
-- 来源: G0049 化验模板明细
-- 新增表
CREATE TABLE IF NOT EXISTS sync_lab_template_item (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    template_id         BIGINT UNSIGNED NOT NULL COMMENT '模板ID',
    item_no             INT COMMENT '明细序号',
    item_code           VARCHAR(50) COMMENT '项目编码',
    item_name           VARCHAR(200) COMMENT '项目名称',
    serial_no           VARCHAR(50) COMMENT '模板编号',
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    KEY idx_template (template_id),
    KEY idx_item_code (item_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-化验模板明细';

-- 1.18 检验标本字典 (sync_lab_specimen)
-- 来源: G0045 检验标本字典数据
-- 新增表
CREATE TABLE IF NOT EXISTS sync_lab_specimen (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    specimen_code       VARCHAR(50) NOT NULL COMMENT '标本代码',
    specimen_name       VARCHAR(100) COMMENT '标本名称',
    input_code          VARCHAR(50) COMMENT '输入码',
    is_valid            TINYINT DEFAULT 1 COMMENT '是否可用: 1是 0否',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_specimen_code (specimen_code),
    KEY idx_name (specimen_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-检验标本字典';

-- 1.19 检验容器字典 (sync_lab_container)
-- 来源: G0047 检验容器字典数据
-- 新增表
CREATE TABLE IF NOT EXISTS sync_lab_container (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    container_code      VARCHAR(50) NOT NULL COMMENT '容器代码',
    container_name      VARCHAR(100) COMMENT '容器名称',
    lab_item_code       VARCHAR(50) COMMENT '关联项目代码',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_container_code (container_code),
    KEY idx_name (container_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-检验容器字典';

-- 1.20 4+7药品目录 (sync_drug_4p7)
-- 来源: G0058 4+7药品目录字典
-- 新增表
CREATE TABLE IF NOT EXISTS sync_drug_4p7 (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    drug_code           VARCHAR(50) NOT NULL COMMENT '药品编码',
    drug_spec           VARCHAR(200) COMMENT '规格',
    firm_id             VARCHAR(200) COMMENT '生产厂家',
    units               VARCHAR(50) COMMENT '单位',
    drug_name           VARCHAR(200) COMMENT '药品名称',
    normal_name         VARCHAR(200) COMMENT '通用名称',
    oper_date           DATETIME COMMENT '操作日期',
    operator            VARCHAR(100) COMMENT '操作员',
    create_date         DATETIME COMMENT '创建时间',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_drug_code (drug_code),
    KEY idx_name (drug_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-4+7药品目录';

-- 1.21 抗菌药物权限 (sync_antibiotic_perm)
-- 来源: G0059 抗菌药物权限字典
-- 新增表
CREATE TABLE IF NOT EXISTS sync_antibiotic_perm (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    emp_no              VARCHAR(50) NOT NULL COMMENT '工号',
    title               VARCHAR(50) COMMENT '职称',
    user_name           VARCHAR(100) COMMENT '用户名',
    memo                VARCHAR(200) COMMENT '权限: 限制类+非限制类/特殊使用级等',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_emp_no (emp_no),
    KEY idx_name (user_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-抗菌药物权限';

-- 1.22 医嘱执行科室对照 (sync_exec_dept_map)
-- 来源: G0044 医嘱项目执行科室对照字典
-- 新增表
CREATE TABLE IF NOT EXISTS sync_exec_dept_map (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    clinic_item_class   VARCHAR(10) COMMENT '医嘱项目类别',
    clinic_item_code    VARCHAR(50) COMMENT '医嘱项目编码',
    exec_dept_code      VARCHAR(50) COMMENT '执行科室代码',
    last_sync_time      DATETIME,
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    KEY idx_item (clinic_item_class, clinic_item_code),
    KEY idx_exec_dept (exec_dept_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='同步-医嘱执行科室对照';

-- ============================================
-- 二、应急业务表（应急期间读写）
-- ============================================

-- 2.1 应急就诊记录 (emergency_encounter)
-- 说明: 应急期间产生的就诊记录
CREATE TABLE IF NOT EXISTS emergency_encounter (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    visit_no            VARCHAR(50) NOT NULL COMMENT '应急就诊号: E+YYYYMMDD+4位序号',
    patient_id          VARCHAR(50) NOT NULL COMMENT '患者ID',
    patient_name        VARCHAR(100) COMMENT '患者姓名',
    patient_age         VARCHAR(20) COMMENT '患者年龄',
    gender              VARCHAR(10) COMMENT '性别',
    id_card_no          VARCHAR(50) COMMENT '身份证号',
    phone               VARCHAR(50) COMMENT '手机号',
    charge_type         VARCHAR(50) COMMENT '费别',
    encounter_type      VARCHAR(20) DEFAULT 'OPD' COMMENT '就诊类型: OPD门诊/EMG急诊',
    dept_id             VARCHAR(50) COMMENT '科室ID',
    dept_name           VARCHAR(100) COMMENT '科室名称',
    doctor_id           VARCHAR(50) COMMENT '医生ID',
    doctor_name         VARCHAR(100) COMMENT '医生姓名',
    doctor_title        VARCHAR(50) COMMENT '医生职称',
    reg_time            DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '挂号时间',
    start_time          DATETIME COMMENT '看诊开始时间',
    end_time            DATETIME COMMENT '看诊结束时间',
    chief_complaint     VARCHAR(500) COMMENT '主诉',
    present_illness     TEXT COMMENT '现病史',
    diagnosis           VARCHAR(500) COMMENT '诊断',
    diagnosis_icd       VARCHAR(100) COMMENT 'ICD诊断编码',
    encounter_status    VARCHAR(20) DEFAULT 'REGISTERED' COMMENT '状态: REGISTERED已挂号/IN_TREATMENT就诊中/COMPLETED已完成/CANCELLED已取消',
    total_amount        DECIMAL(10,2) DEFAULT 0 COMMENT '总费用',
    paid_amount         DECIMAL(10,2) DEFAULT 0 COMMENT '已支付金额',
    payment_status      VARCHAR(20) DEFAULT 'UNPAID' COMMENT '支付状态: UNPAID未支付/PARTIAL部分支付/PAID已支付/REFUNDED已退费',
    is_emergency        TINYINT DEFAULT 0 COMMENT '是否急诊: 1是 0否',
    is_green_channel    TINYINT DEFAULT 0 COMMENT '是否绿色通道',
    is_return_visit     TINYINT DEFAULT 0 COMMENT '是否复诊',
    is_synced_back      TINYINT DEFAULT 0 COMMENT '是否已回写主HIS: 1是 0否',
    sync_back_time      DATETIME COMMENT '回写时间',
    remark              VARCHAR(500) COMMENT '备注',
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_visit_no (visit_no),
    KEY idx_patient (patient_id),
    KEY idx_status (encounter_status),
    KEY idx_reg_time (reg_time),
    KEY idx_dept (dept_id, encounter_status),
    KEY idx_doctor (doctor_id, encounter_status),
    KEY idx_sync (is_synced_back)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应急-就诊记录';

-- 2.2 应急处方主表 (emergency_prescription)
-- 说明: 应急期间医生开具的处方
-- 修正版: 补充I0003回写所需字段
CREATE TABLE IF NOT EXISTS emergency_prescription (
    id                      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    prescription_no         VARCHAR(50) NOT NULL COMMENT '应急处方号: ERX+YYYYMMDD+4位序号',
    encounter_id            BIGINT UNSIGNED NOT NULL COMMENT '关联就诊ID',
    patient_id              VARCHAR(50) NOT NULL COMMENT '患者ID',
    patient_name            VARCHAR(100) COMMENT '患者姓名',
    p_name                  VARCHAR(100) COMMENT '处方打印姓名(可能与patient_name一致)',
    sex                     VARCHAR(10) COMMENT '性别',
    age                     VARCHAR(20) COMMENT '年龄',
    identity                VARCHAR(50) COMMENT '人员类别',
    charge_type             VARCHAR(50) COMMENT '支付类别',
    presc_type              VARCHAR(20) DEFAULT 'WESTERN' COMMENT '处方类型: WESTERN西药/HERBAL草药',
    presc_attr              VARCHAR(50) COMMENT '处方属性: 普通药品/精麻药品等',
    status                  VARCHAR(20) DEFAULT 'DRAFT' COMMENT '状态: DRAFT草稿/SIGNED已签名/PAID已收费/DISPENSED已发药/CANCELLED已作废',
    prescriber_id           VARCHAR(50) COMMENT '开方医生ID',
    prescriber_name         VARCHAR(100) COMMENT '开方医生姓名',
    prescriber_dept         VARCHAR(50) COMMENT '开方科室',
    prescriber_sign_time    DATETIME COMMENT '签名时间',
    reviewer_id             VARCHAR(50) COMMENT '审核人ID',
    reviewer_name           VARCHAR(100) COMMENT '审核人姓名',
    review_time             DATETIME COMMENT '审核时间',
    total_amount            DECIMAL(10,2) DEFAULT 0 COMMENT '处方总金额',
    drug_amount             DECIMAL(10,2) DEFAULT 0 COMMENT '药品金额',
    diagnosis               VARCHAR(500) COMMENT '诊断',
    diagnosis_icd           VARCHAR(100) COMMENT 'ICD诊断编码',
    memo                    VARCHAR(500) COMMENT '建议/备注',
    perform_dept            VARCHAR(50) COMMENT '执行科室',
    order_dept              VARCHAR(50) COMMENT '开单科室',
    doctor                  VARCHAR(100) COMMENT '开单医生姓名',
    doctor_id               VARCHAR(50) COMMENT '开单医生工号',
    visit_no                VARCHAR(50) COMMENT '就诊号',
    visit_date              DATE COMMENT '就诊日期',
    card_id                 VARCHAR(50) COMMENT '就诊卡ID',
    sys_id                  VARCHAR(50) COMMENT '系统ID',
    visit_times             VARCHAR(10) COMMENT '就诊标识: 1初诊/2复诊',
    presc_sort              TINYINT DEFAULT 1 COMMENT '非适应症标志: 2非适应症 1正常',
    order_type              VARCHAR(10) DEFAULT '1' COMMENT '处方来源: 1医生工作站',
    resure_presc_type       VARCHAR(10) COMMENT '处方类型: 1保内/2医保患者自费/NULL其他',
    oper_code               VARCHAR(50) COMMENT '操作员',
    is_psychotropic         TINYINT DEFAULT 0 COMMENT '是否精神药品',
    is_narcotic             TINYINT DEFAULT 0 COMMENT '是否麻醉药品',
    is_antibiotic           TINYINT DEFAULT 0 COMMENT '是否抗生素',
    is_iv                   TINYINT DEFAULT 0 COMMENT '是否输液',
    is_synced_back          TINYINT DEFAULT 0 COMMENT '是否已回写主HIS',
    sync_back_time          DATETIME COMMENT '回写时间',
    created_at              DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_presc_no (prescription_no),
    KEY idx_encounter (encounter_id),
    KEY idx_patient (patient_id),
    KEY idx_status (status),
    KEY idx_prescriber (prescriber_id),
    KEY idx_sync (is_synced_back)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应急-处方主表';

-- 2.3 应急处方明细 (emergency_prescription_item)
-- 说明: 处方中的药品/项目明细
-- 修正版: 补充I0003回写所需字段
CREATE TABLE IF NOT EXISTS emergency_prescription_item (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    prescription_id     BIGINT UNSIGNED NOT NULL COMMENT '处方ID',
    item_no             INT NOT NULL COMMENT '处方内序号',
    item_class          VARCHAR(10) COMMENT '项目类型: A药品/E检查/C化验等',
    drug_code           VARCHAR(50) COMMENT '药品编码',
    drug_name           VARCHAR(200) NOT NULL COMMENT '药品名称',
    drug_spec           VARCHAR(200) COMMENT '规格',
    manufacturer        VARCHAR(200) COMMENT '生产厂家',
    firm_id             VARCHAR(200) COMMENT '生产厂商(回写主HIS用)',
    dose                DECIMAL(10,4) COMMENT '单次剂量',
    dose_unit           VARCHAR(50) COMMENT '剂量单位',
    frequency           VARCHAR(100) COMMENT '频次: BID/TID/Qd等',
    usage_method        VARCHAR(200) COMMENT '用法: 口服/静脉注射等',
    administration      VARCHAR(200) COMMENT '给药途径',
    course_days         INT DEFAULT 1 COMMENT '疗程天数',
    total_qty           DECIMAL(10,2) COMMENT '总量',
    total_unit          VARCHAR(50) COMMENT '总量单位',
    package_count       DECIMAL(10,2) COMMENT '包装数量',
    package_units       VARCHAR(50) COMMENT '包装单位',
    unit_price          DECIMAL(10,4) COMMENT '单价',
    amount              DECIMAL(10,2) COMMENT '金额',
    performed_dept      VARCHAR(50) COMMENT '执行科室',
    mi_reimburse        TINYINT DEFAULT 1 COMMENT '是否医保报销: 1是 0否',
    mi_self_ratio       DECIMAL(5,4) DEFAULT 0 COMMENT '自付比例',
    herb_role           VARCHAR(50) COMMENT '草药角色: 君/臣/佐/使',
    decoction           VARCHAR(200) COMMENT '煎法',
    allergy_flag        TINYINT DEFAULT 0 COMMENT '过敏标志',
    interact_flag       TINYINT DEFAULT 0 COMMENT '相互作用标志',
    contraind_flag      TINYINT DEFAULT 0 COMMENT '禁忌标志',
    skin_test_req       TINYINT DEFAULT 0 COMMENT '是否需要皮试',
    skin_test_result    VARCHAR(50) COMMENT '皮试结果',
    dispense_status     VARCHAR(20) DEFAULT 'PENDING' COMMENT '发药状态: PENDING待发药/DISPENSED已发药/RETURNED已退药',
    dispense_time       DATETIME COMMENT '发药时间',
    dispenser_id        VARCHAR(50) COMMENT '发药人ID',
    sort_no             INT DEFAULT 0 COMMENT '排序号',
    remark              VARCHAR(200) COMMENT '备注',
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    KEY idx_prescription (prescription_id),
    KEY idx_drug (drug_code),
    KEY idx_dispense (dispense_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应急-处方明细';

-- 2.4 应急处方诊断 (emergency_prescription_diagnosis)
-- 说明: 处方关联的诊断
CREATE TABLE IF NOT EXISTS emergency_prescription_diagnosis (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    prescription_id     BIGINT UNSIGNED NOT NULL COMMENT '处方ID',
    diagnosis_serial    INT NOT NULL COMMENT '诊断序号',
    diagnosis_code      VARCHAR(50) COMMENT '诊断编码(ICD)',
    diagnosis_name      VARCHAR(200) COMMENT '诊断名称',
    yb_reason_code_v2   VARCHAR(50) COMMENT '医保诊断编码',
    yb_reason_name_v2   VARCHAR(200) COMMENT '医保诊断名称',
    is_main             TINYINT DEFAULT 0 COMMENT '是否主要诊断: 1是 0否',
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    KEY idx_prescription (prescription_id),
    KEY idx_code (diagnosis_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应急-处方诊断';

-- 2.5 应急收费记录 (emergency_charge)
-- 说明: 应急期间的收费/退费记录
CREATE TABLE IF NOT EXISTS emergency_charge (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    charge_no           VARCHAR(50) NOT NULL COMMENT '收费单号: EC+YYYYMMDD+4位序号',
    encounter_id        BIGINT UNSIGNED NOT NULL COMMENT '关联就诊ID',
    prescription_id     BIGINT UNSIGNED COMMENT '关联处方ID',
    patient_id          VARCHAR(50) NOT NULL COMMENT '患者ID',
    patient_name        VARCHAR(100) COMMENT '患者姓名',
    item_type           VARCHAR(20) NOT NULL COMMENT '项目类型: DRUG药品/ITEM收费项目',
    item_id             BIGINT UNSIGNED COMMENT '关联ID',
    item_code           VARCHAR(50) COMMENT '项目编码',
    item_name           VARCHAR(200) NOT NULL COMMENT '项目名称',
    item_spec           VARCHAR(200) COMMENT '规格',
    qty                 DECIMAL(10,2) COMMENT '数量',
    unit                VARCHAR(50) COMMENT '单位',
    unit_price          DECIMAL(10,4) COMMENT '单价',
    amount              DECIMAL(10,2) COMMENT '金额',
    discount_amount     DECIMAL(10,2) DEFAULT 0 COMMENT '优惠金额',
    mi_amount           DECIMAL(10,2) DEFAULT 0 COMMENT '医保金额',
    self_amount         DECIMAL(10,2) NOT NULL COMMENT '自付金额',
    charge_status       VARCHAR(20) DEFAULT 'UNPAID' COMMENT '收费状态: UNPAID未收费/PAID已收费/REFUNDED已退费/CANCELLED已取消',
    pay_method          VARCHAR(50) COMMENT '支付方式: CASH现金/WECHAT微信/ALIPAY支付宝/CARD刷卡',
    pay_time            DATETIME COMMENT '支付时间',
    cashier_id          VARCHAR(50) COMMENT '收费员ID',
    cashier_name        VARCHAR(100) COMMENT '收费员姓名',
    invoice_no          VARCHAR(100) COMMENT '发票号',
    invoice_prefix      VARCHAR(20) COMMENT '发票前缀(区分应急)',
    rcpt_no             VARCHAR(50) COMMENT '收据号',
    is_refund           TINYINT DEFAULT 0 COMMENT '是否退费: 1是 0否',
    original_charge_id  BIGINT UNSIGNED COMMENT '原收费记录ID(退费用)',
    is_synced_back      TINYINT DEFAULT 0 COMMENT '是否已回写主HIS',
    sync_back_time      DATETIME COMMENT '回写时间',
    remark              VARCHAR(200) COMMENT '备注',
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_charge_no (charge_no),
    KEY idx_encounter (encounter_id),
    KEY idx_prescription (prescription_id),
    KEY idx_patient (patient_id),
    KEY idx_status (charge_status),
    KEY idx_pay_time (pay_time),
    KEY idx_sync (is_synced_back)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应急-收费记录';

-- 2.6 应急药房发药记录 (emergency_dispense)
-- 说明: 药房发药/退药记录
CREATE TABLE IF NOT EXISTS emergency_dispense (
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    prescription_id     BIGINT UNSIGNED NOT NULL COMMENT '处方ID',
    prescription_no     VARCHAR(50) COMMENT '处方号',
    patient_id          VARCHAR(50) COMMENT '患者ID',
    patient_name        VARCHAR(100) COMMENT '患者姓名',
    drug_code           VARCHAR(50) COMMENT '药品编码',
    drug_name           VARCHAR(200) COMMENT '药品名称',
    drug_spec           VARCHAR(200) COMMENT '规格',
    batch_no            VARCHAR(100) COMMENT '批次号',
    dispense_qty        DECIMAL(10,2) COMMENT '发药数量(正数发药/负数退药)',
    dispense_unit       VARCHAR(50) COMMENT '单位',
    storage             VARCHAR(50) COMMENT '药房编码',
    dispenser_id        VARCHAR(50) COMMENT '发药人ID',
    dispenser_name      VARCHAR(100) COMMENT '发药人姓名',
    dispense_time       DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '发药时间',
    window_no           VARCHAR(20) COMMENT '发药窗口',
    remark              VARCHAR(200) COMMENT '备注',
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    KEY idx_prescription (prescription_id),
    KEY idx_patient (patient_id),
    KEY idx_drug (drug_code),
    KEY idx_time (dispense_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应急-药房发药记录';

-- 2.7 应急检查申请 (emergency_exam_request)
-- 说明: 应急期间医生开的检查申请
-- 修正版: 补充G0052所需字段
CREATE TABLE IF NOT EXISTS emergency_exam_request (
    id                      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    request_no              VARCHAR(50) NOT NULL COMMENT '申请单号: EEX+YYYYMMDD+4位序号',
    encounter_id            BIGINT UNSIGNED NOT NULL COMMENT '就诊ID',
    patient_id              VARCHAR(50) COMMENT '患者ID',
    patient_name            VARCHAR(100) COMMENT '患者姓名',
    patient_type            VARCHAR(10) COMMENT '患者来源: 1门诊/2住院/3留观',
    visit_number            VARCHAR(50) COMMENT '就诊号码',
    medical_record_id       VARCHAR(50) COMMENT '病案号',
    exam_no                 VARCHAR(50) COMMENT '检查申请单编号(回写主HIS后)',
    exam_name               VARCHAR(200) COMMENT '检查名称',
    exam_code               VARCHAR(50) COMMENT '检查项目编码',
    exam_class              VARCHAR(50) COMMENT '检查类别: CT/MRI/超声/X光等',
    exam_class_code         VARCHAR(50) COMMENT '检查类别编码',
    body_part               VARCHAR(200) COMMENT '检查部位',
    exam_method             VARCHAR(200) COMMENT '检查方法',
    expense_desc            VARCHAR(50) COMMENT '检查费别',
    apply_dept              VARCHAR(50) COMMENT '申请科室',
    apply_doctor            VARCHAR(100) COMMENT '申请医生',
    apply_time              DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    submit_doc_code         VARCHAR(50) COMMENT '送检医生代码',
    submit_dept_code        VARCHAR(50) COMMENT '申请科室代码',
    perform_dept            VARCHAR(50) COMMENT '执行科室',
    accept_dept_code        VARCHAR(50) COMMENT '接收科室代码',
    accept_dept_desc        VARCHAR(100) COMMENT '接收科室名称',
    request_status          VARCHAR(20) DEFAULT 'REQUESTED' COMMENT '状态: REQUESTED已申请/SAMPLED已采样/TESTING检查中/COMPLETED已完成/CANCELLED已取消',
    is_emergency            TINYINT DEFAULT 0 COMMENT '是否加急: 1是 0否',
    diag_desc               VARCHAR(500) COMMENT '诊断描述',
    expense                 DECIMAL(10,2) COMMENT '检查费用',
    orderitemlist           JSON COMMENT '医嘱明细列表(JSON)',
    is_synced_back          TINYINT DEFAULT 0 COMMENT '是否已回写主HIS',
    remark                  VARCHAR(500) COMMENT '备注',
    created_at              DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_request_no (request_no),
    KEY idx_encounter (encounter_id),
    KEY idx_patient (patient_id),
    KEY idx_status (request_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应急-检查申请';

-- 2.8 应急检验申请 (emergency_lab_request)
-- 说明: 应急期间医生开的检验申请
-- 修正版: 补充G0065所需字段
CREATE TABLE IF NOT EXISTS emergency_lab_request (
    id                      BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    request_no              VARCHAR(50) NOT NULL COMMENT '申请单号: EL+YYYYMMDD+4位序号',
    encounter_id            BIGINT UNSIGNED NOT NULL COMMENT '就诊ID',
    patient_id              VARCHAR(50) COMMENT '患者ID',
    patient_name            VARCHAR(100) COMMENT '患者姓名',
    patient_type            VARCHAR(10) COMMENT '患者来源: 1门诊/2住院/3留观',
    identity_num            VARCHAR(50) COMMENT '身份证号',
    card_num                VARCHAR(50) COMMENT '卡号',
    visit_number            VARCHAR(50) COMMENT '就诊号码',
    patient_age             DECIMAL(5,1) COMMENT '年龄',
    patient_sex_code        VARCHAR(10) COMMENT '性别',
    test_no                 VARCHAR(50) COMMENT '检验申请单号(回写主HIS后)',
    specimen_code           VARCHAR(50) COMMENT '标本代码',
    specimen_name           VARCHAR(100) COMMENT '标本名称',
    specimen_id             VARCHAR(50) COMMENT '申请单号/条码号',
    specimen_muban          VARCHAR(50) COMMENT '模版号',
    lab_item_code           VARCHAR(50) COMMENT '检验项目编码',
    lab_item_name           VARCHAR(200) COMMENT '检验项目名称',
    diag_desc               VARCHAR(500) COMMENT '诊断描述',
    apply_dept              VARCHAR(50) COMMENT '申请科室',
    apply_doctor            VARCHAR(100) COMMENT '申请医生',
    apply_time              DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    submit_doc_code         VARCHAR(50) COMMENT '开单医生代码',
    submit_dept_code        VARCHAR(50) COMMENT '开单科室代码',
    submit_ward_code        VARCHAR(50) COMMENT '病区代码',
    perform_dept            VARCHAR(50) COMMENT '执行科室',
    accept_dept_code        VARCHAR(50) COMMENT '执行科室编码',
    accept_dept_desc        VARCHAR(100) COMMENT '执行科室名称',
    request_status          VARCHAR(20) DEFAULT 'REQUESTED' COMMENT '状态',
    charging_code           TINYINT DEFAULT 0 COMMENT '计费标识: 0未收费 1已收费',
    submit_time             DATETIME COMMENT '开单时间',
    sample_time             DATETIME COMMENT '采样时间',
    orderitemlist           JSON COMMENT '检验项目明细(JSON)',
    is_synced_back          TINYINT DEFAULT 0 COMMENT '是否已回写主HIS',
    remark                  VARCHAR(500) COMMENT '备注',
    created_at              DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_request_no (request_no),
    KEY idx_encounter (encounter_id),
    KEY idx_patient (patient_id),
    KEY idx_status (request_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='应急-检验申请';

-- ============================================
-- 三、系统配置与审计表
-- ============================================

-- 3.1 系统配置 (sys_config)
CREATE TABLE IF NOT EXISTS sys_config (
    config_key      VARCHAR(100) PRIMARY KEY COMMENT '配置键',
    config_value    VARCHAR(500) COMMENT '配置值',
    remark          VARCHAR(200) COMMENT '备注',
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置';

-- 3.2 审计日志 (sys_audit_log)
CREATE TABLE IF NOT EXISTS sys_audit_log (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id         VARCHAR(50) COMMENT '用户ID',
    user_name       VARCHAR(100) COMMENT '用户姓名',
    operation       VARCHAR(100) COMMENT '操作类型: LOGIN/SAVE_PRESC/CHARGE/DISPENSE等',
    target_type     VARCHAR(50) COMMENT '操作对象类型',
    target_id       VARCHAR(100) COMMENT '操作对象ID',
    detail          TEXT COMMENT '详情JSON',
    ip_address      VARCHAR(50) COMMENT 'IP地址',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    KEY idx_user (user_id, created_at),
    KEY idx_operation (operation, created_at),
    KEY idx_time (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='审计日志';

-- 3.3 数据同步日志 (sync_log)
CREATE TABLE IF NOT EXISTS sync_log (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    sync_type       VARCHAR(50) NOT NULL COMMENT '同步类型: DEPT/DOCTOR/DRUG/STOCK/PRICE/DIAGNOSIS/SCHEDULE等',
    sync_status     VARCHAR(20) COMMENT '状态: SUCCESS/FAILED',
    record_count    INT DEFAULT 0 COMMENT '同步记录数',
    error_msg       TEXT COMMENT '错误信息',
    start_time      DATETIME COMMENT '开始时间',
    end_time        DATETIME COMMENT '结束时间',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    KEY idx_type (sync_type, created_at),
    KEY idx_status (sync_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据同步日志';

-- ============================================
-- 四、初始化默认配置
-- ============================================

INSERT INTO sys_config (config_key, config_value, remark) VALUES
('emergency_mode', '0', '是否应急模式: 0=正常(静默同步) 1=应急(停止同步)'),
('sync_interval_min', '15', '数据同步间隔(分钟)'),
('visit_no_prefix', 'E', '应急就诊号前缀'),
('presc_no_prefix', 'ERX', '应急处方号前缀'),
('charge_no_prefix', 'EC', '应急收费单号前缀'),
('exam_no_prefix', 'EEX', '应急检查申请号前缀'),
('lab_no_prefix', 'EL', '应急检验申请号前缀'),
('invoice_prefix', 'E', '应急发票前缀(与主HIS区分)'),
('last_sync_time', '', '最后一次数据同步时间'),
('main_his_url', '', '主HIS接口地址'),
('main_his_meskey', '', '主HIS MESKEY'),
('enable_encrypt', '0', '是否启用接口加密: 0=否 1=是'),
('data_sync_enabled', '1', '是否启用数据同步: 1是 0否'),
('max_sync_retry', '3', '同步失败最大重试次数')
ON DUPLICATE KEY UPDATE config_value = VALUES(config_value);

