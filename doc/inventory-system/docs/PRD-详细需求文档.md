# 库存管理系统 — 产品需求文档（PRD）

> 版本：v1.0  
> 日期：2026-02-24  
> 基于：库存系统需求文档（简版）

---

## 目录

1. [项目概述](#1-项目概述)
2. [术语定义](#2-术语定义)
3. [用户角色与权限矩阵](#3-用户角色与权限矩阵)
4. [功能需求详细说明](#4-功能需求详细说明)
5. [前台页面需求](#5-前台页面需求)
6. [后台管理页面需求](#6-后台管理页面需求)
7. [图片与文件管理需求](#7-图片与文件管理需求)
8. [数据模型设计](#8-数据模型设计)
9. [接口需求清单](#9-接口需求清单)
10. [非功能需求](#10-非功能需求)
11. [验收标准](#11-验收标准)

---

## 1. 项目概述

### 1.1 项目背景

企业内部仓储/物料管理场景中，需要一套系统来管理库存信息，包括物料的增删改查、附件（图片/文件）管理、操作日志审计等基础能力。

### 1.2 项目目标

- 实现库存信息的完整生命周期管理（新增→查询→编辑→删除/停用）
- 提供图片和文件的上传、预览、下载、管理能力
- 支持多角色权限控制（管理员、操作员、只读用户）
- 提供操作日志与变更审计追溯
- 前台（业务使用端）与后台（管理运营端）分离

### 1.3 项目范围

| 范围内 | 范围外 |
|--------|--------|
| 单仓库基础库存管理 | 多仓库间调拨 |
| 库存CRUD与附件管理 | 采购管理 |
| 角色权限控制 | 销售管理 |
| 操作日志审计 | 财务对接 |
| 数据导出（Excel/CSV） | 数据BI分析 |
| 库存预警（安全库存阈值） | 自动补货 |

### 1.4 适用用户

- 仓库管理员
- 仓库操作员
- 业务查询人员（只读）

---

## 2. 术语定义

| 术语 | 说明 |
|------|------|
| 物料编码 | 唯一标识一种物料的编码，可手动输入或按规则自动生成 |
| 安全库存 | 库存预警阈值，当前库存数量低于此值时触发预警提示 |
| 软删除 | 逻辑删除，数据不从数据库物理移除，通过标记字段（如 `is_deleted`）实现，可恢复 |
| 附件 | 与库存记录关联的图片或文件，如产品图片、质检报告、说明书等 |
| 前台 | 面向业务操作人员的用户界面，以查询和基本操作为主 |
| 后台 | 面向管理员的运营管理界面，包含全部管理能力 |

---

## 3. 用户角色与权限矩阵

### 3.1 角色定义

| 角色 | 说明 |
|------|------|
| **管理员（Admin）** | 拥有系统全部权限，可管理用户、角色、系统配置 |
| **操作员（Operator）** | 可执行库存的新增、编辑、查询操作；删除权限可配置 |
| **只读用户（Viewer）** | 仅可查询与导出库存数据 |

### 3.2 权限矩阵

| 功能模块 | 具体操作 | 管理员 | 操作员 | 只读用户 |
|----------|----------|:------:|:------:|:--------:|
| **用户管理** | 创建/编辑/停用账号 | ✅ | ❌ | ❌ |
| | 角色分配 | ✅ | ❌ | ❌ |
| **库存管理** | 查看列表/详情 | ✅ | ✅ | ✅ |
| | 新增库存记录 | ✅ | ✅ | ❌ |
| | 编辑库存记录 | ✅ | ✅ | ❌ |
| | 删除库存记录 | ✅ | ⚙️ 可配置 | ❌ |
| | 批量删除 | ✅ | ❌ | ❌ |
| | 导出数据 | ✅ | ✅ | ✅ |
| **附件管理** | 上传图片/文件 | ✅ | ✅ | ❌ |
| | 下载文件 | ✅ | ✅ | ✅ |
| | 删除附件 | ✅ | ⚙️ 可配置 | ❌ |
| | 重命名附件 | ✅ | ✅ | ❌ |
| **日志审计** | 查看操作日志 | ✅ | ❌ | ❌ |
| **系统配置** | 修改系统参数 | ✅ | ❌ | ❌ |

> ⚙️ 表示该权限可通过后台配置开启或关闭

### 3.3 数据权限

- 前台用户仅可查看自己有权限的仓库/分类下的数据
- 后台管理员可查看所有数据
- 数据权限通过 **用户-仓库/分类** 关联关系实现

---

## 4. 功能需求详细说明

### 4.1 库存基础信息管理

#### 4.1.1 数据字段定义

| 字段名 | 字段标识 | 类型 | 必填 | 唯一 | 说明 |
|--------|----------|------|:----:|:----:|------|
| 库存ID | `id` | 整数/UUID | 系统生成 | ✅ | 主键，自动递增或UUID |
| 物料编码 | `material_code` | 字符串(50) | ✅ | ✅ | 支持手动输入或规则生成 |
| 物料名称 | `material_name` | 字符串(200) | ✅ | ❌ | — |
| 分类 | `category` | 枚举/字符串 | ✅ | ❌ | 原料/成品/备件/其他 |
| 规格型号 | `specification` | 字符串(200) | ❌ | ❌ | — |
| 单位 | `unit` | 字符串(20) | ✅ | ❌ | 件、箱、kg、米等 |
| 当前库存数量 | `quantity` | 整数/小数 | ✅ | ❌ | ≥ 0 |
| 安全库存 | `safety_stock` | 整数/小数 | ❌ | ❌ | 预警阈值，默认0 |
| 仓库位置 | `location` | 字符串(200) | ❌ | ❌ | 库区/货架编号 |
| 状态 | `status` | 枚举 | ✅ | ❌ | 启用(active)/停用(inactive) |
| 备注 | `remark` | 文本(1000) | ❌ | ❌ | — |
| 创建人 | `created_by` | 字符串 | 系统填充 | ❌ | 关联用户ID |
| 创建时间 | `created_at` | 日期时间 | 系统填充 | ❌ | — |
| 更新时间 | `updated_at` | 日期时间 | 系统填充 | ❌ | — |
| 更新人 | `updated_by` | 字符串 | 系统填充 | ❌ | — |
| 是否删除 | `is_deleted` | 布尔 | 系统填充 | ❌ | 软删除标记，默认false |

#### 4.1.2 物料编码生成规则（建议）

- **手动输入**：用户自行填写，提交时检查唯一性
- **规则生成（可选）**：`{分类前缀}-{年月}-{4位流水号}`，如 `YL-202602-0001`
- 编码一旦创建不可修改（如需修改，需管理员操作）

### 4.2 新增库存

#### 用例描述

| 项目 | 内容 |
|------|------|
| 用例名称 | 新增库存记录 |
| 参与者 | 管理员、操作员 |
| 前置条件 | 用户已登录且拥有新增权限 |
| 主流程 | 1. 用户点击"新增库存"按钮<br>2. 系统展示新增表单<br>3. 用户填写必填字段，可上传附件<br>4. 用户点击"提交"<br>5. 系统校验数据<br>6. 校验通过，保存数据，返回成功提示<br>7. 页面跳转至列表或详情页 |
| 异常流程 | 5a. 物料编码已存在 → 提示"物料编码已存在，请修改"<br>5b. 必填字段为空 → 高亮标记并提示具体字段<br>5c. 附件上传失败 → 提示失败原因，允许重试，库存信息可先保存 |
| 后置条件 | 库存记录创建成功，操作日志记录 |

#### 校验规则

| 字段 | 校验规则 |
|------|----------|
| 物料编码 | 非空、长度1-50、唯一性校验、仅允许字母/数字/中划线 |
| 物料名称 | 非空、长度1-200 |
| 分类 | 非空、必须为预定义枚举值 |
| 单位 | 非空、长度1-20 |
| 当前库存数量 | 非空、数值类型、≥ 0、最多两位小数 |
| 安全库存 | 数值类型、≥ 0 |

### 4.3 查询库存

#### 4.3.1 列表查询

| 项目 | 内容 |
|------|------|
| 默认展示 | 分页列表，每页20条（可配置），按更新时间倒序 |
| 展示字段 | 物料编码、物料名称、分类、规格型号、单位、当前库存数量、安全库存、状态、更新时间 |
| 特殊标记 | 当前库存数量 < 安全库存时，数量列标红显示⚠️预警 |

#### 4.3.2 筛选条件

| 筛选条件 | 类型 | 说明 |
|----------|------|------|
| 物料编码/名称 | 文本输入 | 模糊匹配（LIKE %keyword%） |
| 分类 | 下拉选择 | 支持多选 |
| 状态 | 下拉选择 | 启用/停用/全部 |
| 库存数量区间 | 数值区间 | 最小值-最大值 |
| 仓库位置 | 文本输入 | 模糊匹配 |

#### 4.3.3 详情页

- 展示库存记录全部字段信息
- 图片区域：缩略图网格展示，点击可查看大图（支持左右切换）
- 文件区域：文件列表，展示文件名、大小、上传人、上传时间，支持下载
- 变更历史：展示该记录的变更日志（最近N条）

#### 4.3.4 数据导出

| 项目 | 内容 |
|------|------|
| 导出格式 | Excel (.xlsx) / CSV |
| 导出范围 | 当前筛选条件下的全部数据（非仅当前页） |
| 导出字段 | 与列表展示字段一致，可勾选 |
| 数量限制 | 单次导出上限 10,000 条 |
| 交互方式 | 点击"导出"按钮 → 后台生成文件 → 浏览器自动下载 |

### 4.4 编辑库存

| 项目 | 内容 |
|------|------|
| 用例名称 | 编辑库存记录 |
| 可编辑字段 | 除 id、created_by、created_at 外的所有字段 |
| 物料编码 | 默认不可修改；管理员可开启修改权限 |
| 校验规则 | 与新增一致 |
| 附件操作 | 支持追加新附件、替换已有附件、删除已有附件 |
| 变更记录 | 系统自动记录字段变更前后值（old_value → new_value） |
| 并发控制 | 使用乐观锁（版本号）防止并发编辑冲突 |

### 4.5 删除库存

| 项目 | 内容 |
|------|------|
| 删除方式 | 软删除（设置 `is_deleted = true`） |
| 单条删除 | 点击删除按钮 → 弹窗二次确认"确定要删除该库存记录？" → 确认后执行 |
| 批量删除 | 勾选多条 → 点击批量删除 → 弹窗显示将删除的数量 → 确认后执行 |
| 数据恢复 | 管理员可在后台"已删除数据"列表中恢复 |
| 关联处理 | 删除库存记录时，关联附件同步标记为已删除（不物理删除文件） |
| 日志记录 | 记录删除操作人、时间、被删除记录的关键信息 |

---

## 5. 前台页面需求

### 5.1 登录页

| 项目 | 内容 |
|------|------|
| 页面路径 | `/login` |
| 功能 | 账号密码登录 |
| 输入字段 | 用户名/手机号、密码 |
| 校验 | 用户名非空、密码非空（6-20位） |
| 登录成功 | 跳转至库存列表页，存储登录Token |
| 登录失败 | 提示"用户名或密码错误"（不区分具体原因） |
| 退出登录 | 右上角用户头像下拉菜单 → 退出登录 → 清除Token → 跳转登录页 |
| 登录状态保持 | Token有效期 24小时，支持自动续期 |

### 5.2 库存列表页

| 项目 | 内容 |
|------|------|
| 页面路径 | `/inventory` |
| 布局 | 顶部搜索栏 + 筛选区域 + 数据表格 + 分页器 |
| 搜索栏 | 关键词搜索（物料编码/名称），回车或点击搜索触发 |
| 筛选区域 | 分类下拉、状态下拉、库存数量区间 |
| 表格列 | 物料编码、物料名称、分类、规格型号、单位、库存数量、安全库存、状态、操作 |
| 操作列 | "查看详情"按钮；操作员额外展示"编辑"按钮 |
| 分页 | 底部分页器，显示总条数，支持切换每页条数（10/20/50） |
| 预警提示 | 库存数量 < 安全库存的行高亮（浅红底色），数量列显示⚠️图标 |
| 空状态 | 无数据时展示"暂无库存数据"插画+提示 |
| 导出按钮 | 右上角，根据当前筛选条件导出 |

### 5.3 库存详情页

| 项目 | 内容 |
|------|------|
| 页面路径 | `/inventory/:id` |
| 基础信息区 | 卡片形式展示所有字段（两列/三列自适应布局） |
| 图片区 | 缩略图网格展示，点击弹出Lightbox大图查看器（支持左右翻页） |
| 文件区 | 文件列表表格：文件名、大小、上传人、上传时间、操作（下载） |
| 返回 | 左上角返回按钮，返回列表页（保留之前的筛选条件） |

---

## 6. 后台管理页面需求

### 6.1 总体布局

- **侧边栏导航**：左侧固定菜单栏
- **顶部栏**：显示当前页面标题、用户信息、退出登录
- **内容区域**：右侧主体内容
- **路径前缀**：`/admin/*`

### 6.2 菜单结构

```
后台管理
├── 仪表盘（Dashboard）         /admin/dashboard
├── 库存管理                     /admin/inventory
│   ├── 库存列表                 /admin/inventory/list
│   ├── 新增库存                 /admin/inventory/create
│   └── 已删除数据               /admin/inventory/deleted
├── 附件管理                     /admin/attachments
├── 用户管理                     /admin/users
│   ├── 用户列表                 /admin/users/list
│   └── 角色管理                 /admin/users/roles
├── 日志审计                     /admin/logs
│   ├── 操作日志                 /admin/logs/operations
│   └── 变更日志                 /admin/logs/changes
└── 系统配置                     /admin/settings
```

### 6.3 仪表盘（Dashboard）

| 展示内容 | 说明 |
|----------|------|
| 库存总数 | 启用状态的库存记录总条数 |
| 分类统计 | 按分类展示库存数量（饼图/柱状图） |
| 预警列表 | 当前库存 < 安全库存的记录列表（Top 10） |
| 近期操作 | 最近的操作日志（最新10条） |

### 6.4 库存管理页

在前台列表功能基础上，增加以下能力：

| 功能 | 说明 |
|------|------|
| 新增按钮 | 跳转新增库存表单页 |
| 编辑按钮 | 每行操作列，跳转编辑页 |
| 删除按钮 | 单条删除，二次确认弹窗 |
| 批量操作栏 | 勾选多条后，底部浮出操作栏：批量删除、批量修改状态 |
| 已删除数据 | 可查看软删除的记录，支持恢复操作 |
| 数据导入 | 支持 Excel 模板导入库存数据（可选功能） |

### 6.5 附件管理页

| 功能 | 说明 |
|------|------|
| 附件列表 | 展示所有附件：文件名、类型、大小、关联库存记录、上传人、上传时间 |
| 筛选 | 按文件类型（图片/文档）、上传时间范围、关联库存记录筛选 |
| 操作 | 预览（图片）、下载、重命名、删除 |

### 6.6 用户管理页

| 功能 | 说明 |
|------|------|
| 用户列表 | 用户名、角色、状态（启用/停用）、创建时间 |
| 新增用户 | 用户名、初始密码、角色选择 |
| 编辑用户 | 修改角色、重置密码、启用/停用 |
| 角色管理 | 查看角色列表、编辑角色权限配置 |

### 6.7 日志审计页

| 功能 | 说明 |
|------|------|
| 操作日志 | 展示字段：操作时间、操作人、操作类型（新增/编辑/删除/上传/下载）、操作对象、操作详情 |
| 变更日志 | 展示字段：变更时间、操作人、库存记录、变更字段、旧值、新值 |
| 筛选 | 按时间范围、操作类型、操作人筛选 |

### 6.8 系统配置页（可选）

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| 附件大小限制 | 单文件最大上传大小 | 20MB |
| 附件格式白名单 | 允许上传的文件扩展名 | jpg,jpeg,png,webp,pdf,doc,docx,xls,xlsx,txt |
| 每页条数 | 列表默认分页大小 | 20 |
| Token有效期 | 登录Token过期时间 | 24小时 |
| 安全库存预警 | 全局默认安全库存值 | 0 |
| 操作员删除权限 | 是否允许操作员删除库存 | 关闭 |

---

## 7. 图片与文件管理需求

### 7.1 上传需求

| 项目 | 内容 |
|------|------|
| 上传方式 | 拖拽上传 + 点击选择文件 |
| 多文件上传 | 支持，单次最多上传 10 个文件 |
| 图片格式 | jpg, jpeg, png, webp |
| 文件格式 | pdf, doc, docx, xls, xlsx, txt |
| 单文件大小 | ≤ 20MB（可配置） |
| 上传进度 | 展示每个文件的上传进度条 |
| 上传结果 | 每个文件独立反馈成功/失败，失败显示原因 |

#### 上传失败原因提示

| 场景 | 提示信息 |
|------|----------|
| 格式不支持 | "不支持的文件格式，请上传 jpg/png/pdf 等格式的文件" |
| 文件过大 | "文件大小超过限制（最大20MB），请压缩后重新上传" |
| 网络失败 | "网络异常，上传失败，请重试" |
| 服务器错误 | "服务器繁忙，请稍后重试" |

### 7.2 展示与操作

| 功能 | 说明 |
|------|------|
| 图片缩略图 | 上传后自动生成缩略图（宽度200px），列表中展示缩略图 |
| 图片大图查看 | 点击缩略图弹出Lightbox，支持缩放、左右切换 |
| 文件下载 | 点击下载按钮，浏览器直接下载原始文件 |
| 文件重命名 | 仅修改显示名称，存储文件名不变 |
| 附件删除 | 二次确认后标记删除（软删除） |
| 信息展示 | 每个附件展示：文件名、文件大小、上传人、上传时间 |

### 7.3 存储安全

| 项目 | 内容 |
|------|------|
| 存储路径 | `uploads/{年月}/{库存ID}/{UUID}.{扩展名}` |
| 文件命名 | 使用UUID重命名，避免冲突和中文路径问题 |
| 安全校验 | 后端校验文件后缀 + MIME类型双重验证 |
| 访问控制 | 文件下载/预览需登录鉴权，权限与库存记录一致 |
| 防盗链 | 文件URL使用签名Token，限时有效 |

---

## 8. 数据模型设计

### 8.1 库存表（inventories）

```sql
CREATE TABLE inventories (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    material_code   VARCHAR(50)  NOT NULL UNIQUE COMMENT '物料编码',
    material_name   VARCHAR(200) NOT NULL COMMENT '物料名称',
    category        VARCHAR(50)  NOT NULL COMMENT '分类: raw_material/finished_product/spare_part/other',
    specification   VARCHAR(200) DEFAULT '' COMMENT '规格型号',
    unit            VARCHAR(20)  NOT NULL COMMENT '单位',
    quantity         DECIMAL(12,2) NOT NULL DEFAULT 0 COMMENT '当前库存数量',
    safety_stock    DECIMAL(12,2) DEFAULT 0 COMMENT '安全库存',
    location        VARCHAR(200) DEFAULT '' COMMENT '仓库位置',
    status          TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1=启用, 0=停用',
    remark          VARCHAR(1000) DEFAULT '' COMMENT '备注',
    version         INT NOT NULL DEFAULT 1 COMMENT '乐观锁版本号',
    is_deleted      TINYINT NOT NULL DEFAULT 0 COMMENT '软删除: 0=正常, 1=已删除',
    created_by      BIGINT NOT NULL COMMENT '创建人ID',
    updated_by      BIGINT DEFAULT NULL COMMENT '更新人ID',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_material_code (material_code),
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存表';
```

### 8.2 附件表（attachments）

```sql
CREATE TABLE attachments (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    inventory_id    BIGINT NOT NULL COMMENT '关联库存ID',
    file_name       VARCHAR(200) NOT NULL COMMENT '原始文件名',
    display_name    VARCHAR(200) NOT NULL COMMENT '显示文件名（可重命名）',
    file_path       VARCHAR(500) NOT NULL COMMENT '存储路径',
    file_size       BIGINT NOT NULL COMMENT '文件大小(bytes)',
    file_type       VARCHAR(20) NOT NULL COMMENT '文件类型: image/document',
    mime_type       VARCHAR(100) NOT NULL COMMENT 'MIME类型',
    extension       VARCHAR(10) NOT NULL COMMENT '文件扩展名',
    thumbnail_path  VARCHAR(500) DEFAULT '' COMMENT '缩略图路径（仅图片）',
    uploaded_by     BIGINT NOT NULL COMMENT '上传人ID',
    is_deleted      TINYINT NOT NULL DEFAULT 0,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_inventory_id (inventory_id),
    INDEX idx_file_type (file_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='附件表';
```

### 8.3 用户表（users）

```sql
CREATE TABLE users (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    username        VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password_hash   VARCHAR(200) NOT NULL COMMENT '密码哈希',
    display_name    VARCHAR(100) DEFAULT '' COMMENT '显示名称',
    role            VARCHAR(20) NOT NULL DEFAULT 'viewer' COMMENT '角色: admin/operator/viewer',
    status          TINYINT NOT NULL DEFAULT 1 COMMENT '状态: 1=启用, 0=停用',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### 8.4 操作日志表（operation_logs）

```sql
CREATE TABLE operation_logs (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id         BIGINT NOT NULL COMMENT '操作人ID',
    user_name       VARCHAR(50) NOT NULL COMMENT '操作人用户名',
    action          VARCHAR(50) NOT NULL COMMENT '操作类型: create/update/delete/upload/download/export',
    target_type     VARCHAR(50) NOT NULL COMMENT '对象类型: inventory/attachment/user',
    target_id       BIGINT DEFAULT NULL COMMENT '对象ID',
    description     VARCHAR(500) DEFAULT '' COMMENT '操作描述',
    ip_address      VARCHAR(50) DEFAULT '' COMMENT '操作IP',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';
```

### 8.5 变更日志表（change_logs）

```sql
CREATE TABLE change_logs (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    inventory_id    BIGINT NOT NULL COMMENT '库存记录ID',
    user_id         BIGINT NOT NULL COMMENT '操作人ID',
    user_name       VARCHAR(50) NOT NULL COMMENT '操作人用户名',
    field_name      VARCHAR(50) NOT NULL COMMENT '变更字段名',
    field_label     VARCHAR(50) NOT NULL COMMENT '变更字段中文名',
    old_value       VARCHAR(500) DEFAULT '' COMMENT '旧值',
    new_value       VARCHAR(500) DEFAULT '' COMMENT '新值',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_inventory_id (inventory_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='变更日志表';
```

### 8.6 ER 关系图（文字描述）

```
users (1) ──→ (N) inventories     [created_by]
users (1) ──→ (N) attachments     [uploaded_by]
users (1) ──→ (N) operation_logs  [user_id]
users (1) ──→ (N) change_logs     [user_id]
inventories (1) ──→ (N) attachments    [inventory_id]
inventories (1) ──→ (N) change_logs    [inventory_id]
```

---

## 9. 接口需求清单

### 9.1 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/login` | 用户登录，返回JWT Token |
| POST | `/api/v1/auth/logout` | 用户登出，Token失效 |
| GET  | `/api/v1/auth/me` | 获取当前登录用户信息 |

### 9.2 库存接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET    | `/api/v1/inventories` | 分页查询库存列表 | 全部角色 |
| GET    | `/api/v1/inventories/:id` | 获取库存详情 | 全部角色 |
| POST   | `/api/v1/inventories` | 新增库存记录 | 管理员、操作员 |
| PUT    | `/api/v1/inventories/:id` | 编辑库存记录 | 管理员、操作员 |
| DELETE | `/api/v1/inventories/:id` | 删除库存记录（软删除） | 管理员、操作员(可配置) |
| POST   | `/api/v1/inventories/batch-delete` | 批量删除 | 管理员 |
| PUT    | `/api/v1/inventories/:id/restore` | 恢复已删除记录 | 管理员 |
| GET    | `/api/v1/inventories/export` | 导出库存数据 | 全部角色 |

### 9.3 附件接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST   | `/api/v1/inventories/:id/attachments` | 上传附件 | 管理员、操作员 |
| GET    | `/api/v1/inventories/:id/attachments` | 获取库存附件列表 | 全部角色 |
| GET    | `/api/v1/attachments/:id/download` | 下载附件 | 全部角色 |
| GET    | `/api/v1/attachments/:id/preview` | 预览附件（图片） | 全部角色 |
| PUT    | `/api/v1/attachments/:id` | 修改附件信息（重命名） | 管理员、操作员 |
| DELETE | `/api/v1/attachments/:id` | 删除附件 | 管理员、操作员(可配置) |

### 9.4 用户管理接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET    | `/api/v1/users` | 用户列表 | 管理员 |
| POST   | `/api/v1/users` | 创建用户 | 管理员 |
| PUT    | `/api/v1/users/:id` | 编辑用户 | 管理员 |
| PUT    | `/api/v1/users/:id/reset-password` | 重置密码 | 管理员 |

### 9.5 日志接口

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/logs/operations` | 查询操作日志 | 管理员 |
| GET | `/api/v1/logs/changes` | 查询变更日志 | 管理员 |

### 9.6 统一响应格式

```json
// 成功
{
    "code": 0,
    "message": "success",
    "data": { ... }
}

// 分页
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [ ... ],
        "total": 100,
        "page": 1,
        "page_size": 20
    }
}

// 错误
{
    "code": 40001,
    "message": "物料编码已存在",
    "data": null
}
```

### 9.7 错误码规范

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 40001 | 参数校验失败 |
| 40002 | 数据重复（如编码重复） |
| 40003 | 数据不存在 |
| 40101 | 未登录/Token过期 |
| 40102 | 权限不足 |
| 40301 | 文件格式不支持 |
| 40302 | 文件大小超限 |
| 50001 | 服务器内部错误 |
| 50002 | 数据库错误 |

---

## 10. 非功能需求

### 10.1 性能要求

| 场景 | 指标 |
|------|------|
| 列表查询（万级数据） | 响应时间 < 2秒 |
| 详情查询 | 响应时间 < 1秒 |
| 新增/编辑/删除 | 响应时间 < 1秒 |
| 文件上传（20MB） | 完成时间 < 30秒（取决于网络） |
| 数据导出（1万条） | 完成时间 < 10秒 |

### 10.2 安全要求

| 项目 | 说明 |
|------|------|
| 身份认证 | 基于JWT Token，所有API接口（除登录外）需携带Token |
| 密码存储 | 使用 bcrypt 加密存储，不可逆 |
| 接口权限 | 每个接口根据角色进行权限校验 |
| 防注入 | 所有SQL查询使用参数化查询，防止SQL注入 |
| 文件安全 | 文件上传后缀+MIME双重校验 |
| CORS | 配置允许的前端域名，禁止跨域滥用 |
| 日志审计 | 所有关键操作记录日志 |

### 10.3 可用性要求

- 关键操作（新增/编辑/删除）有明确的成功/失败反馈
- 表单提交时展示Loading状态，防止重复提交
- 网络异常时给出友好提示
- 页面适配主流浏览器（Chrome、Firefox、Edge、Safari）

### 10.4 可维护性要求

- 代码分层清晰（Controller → Service → Repository）
- 统一错误处理与日志记录
- 配置项集中管理（数据库连接、文件路径、Token过期时间等）
- API接口文档自动生成（Swagger）

---

## 11. 验收标准

### 11.1 功能验收

| 编号 | 验收项 | 验收标准 |
|------|--------|----------|
| F01 | 用户登录/登出 | 正确的账号密码可登录，Token过期后需重新登录 |
| F02 | 库存新增 | 必填项校验通过后可创建成功；物料编码重复时报错 |
| F03 | 库存查询 | 分页正确；各筛选条件生效；详情页展示完整 |
| F04 | 库存编辑 | 修改保存成功；变更日志正确记录字段前后值 |
| F05 | 库存删除 | 软删除后列表不可见；管理员可恢复 |
| F06 | 批量删除 | 勾选多条删除成功；需二次确认 |
| F07 | 附件上传 | 多文件上传成功；格式/大小校验生效；进度展示 |
| F08 | 附件操作 | 图片可预览大图；文件可下载；可重命名和删除 |
| F09 | 数据导出 | Excel/CSV导出内容与列表筛选结果一致 |
| F10 | 权限控制 | 三种角色权限边界正确，越权操作被拒绝 |
| F11 | 操作日志 | 关键操作均有日志记录，可按条件查询 |
| F12 | 前后台分离 | 前台/后台页面功能边界清晰 |

### 11.2 非功能验收

| 编号 | 验收项 | 验收标准 |
|------|--------|----------|
| N01 | 查询性能 | 万级数据量下列表查询响应 < 2秒 |
| N02 | 安全性 | 无Token无法访问API；密码加密存储 |
| N03 | 容错性 | 网络异常/服务异常有友好提示 |
| N04 | 浏览器兼容 | Chrome/Firefox/Edge/Safari正常运行 |

---

> **文档结束**  
> 下一步：技术架构与实现方案文档、项目代码组织结构文档
