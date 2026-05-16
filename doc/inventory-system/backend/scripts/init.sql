-- ============================
-- 库存管理系统 - 数据库初始化脚本
-- ============================
-- 使用方法：mysql -u root -p < scripts/init.sql

CREATE DATABASE IF NOT EXISTS inventory_db
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE inventory_db;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    username        VARCHAR(50)  NOT NULL UNIQUE COMMENT '用户名',
    password_hash   VARCHAR(200) NOT NULL COMMENT '密码哈希',
    display_name    VARCHAR(100) DEFAULT '' COMMENT '显示名称',
    role            VARCHAR(20)  NOT NULL DEFAULT 'viewer' COMMENT '角色: admin/operator/viewer',
    status          TINYINT      NOT NULL DEFAULT 1 COMMENT '状态: 1=启用, 0=停用',
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 库存表
CREATE TABLE IF NOT EXISTS inventories (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    material_code   VARCHAR(50)   NOT NULL UNIQUE COMMENT '物料编码',
    material_name   VARCHAR(200)  NOT NULL COMMENT '物料名称',
    category        VARCHAR(50)   NOT NULL COMMENT '分类',
    specification   VARCHAR(200)  DEFAULT '' COMMENT '规格型号',
    unit            VARCHAR(20)   NOT NULL COMMENT '单位',
    quantity        DECIMAL(12,2) NOT NULL DEFAULT 0 COMMENT '库存数量',
    safety_stock    DECIMAL(12,2) DEFAULT 0 COMMENT '安全库存',
    location        VARCHAR(200)  DEFAULT '' COMMENT '仓库位置',
    status          TINYINT       NOT NULL DEFAULT 1 COMMENT '状态: 1=启用, 0=停用',
    remark          VARCHAR(1000) DEFAULT '' COMMENT '备注',
    version         INT           NOT NULL DEFAULT 1 COMMENT '乐观锁版本号',
    is_deleted      TINYINT       NOT NULL DEFAULT 0 COMMENT '软删除: 0=正常, 1=已删除',
    created_by      BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    updated_by      BIGINT UNSIGNED DEFAULT NULL COMMENT '更新人ID',
    created_at      DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_material_code (material_code),
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_is_deleted (is_deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存表';

-- 附件表
CREATE TABLE IF NOT EXISTS attachments (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    inventory_id    BIGINT UNSIGNED NOT NULL COMMENT '关联库存ID',
    file_name       VARCHAR(200) NOT NULL COMMENT '原始文件名',
    display_name    VARCHAR(200) NOT NULL COMMENT '显示文件名',
    file_path       VARCHAR(500) NOT NULL COMMENT '存储路径',
    file_size       BIGINT       NOT NULL COMMENT '文件大小(字节)',
    file_type       VARCHAR(20)  NOT NULL COMMENT '文件类型: image/document',
    mime_type       VARCHAR(100) NOT NULL COMMENT 'MIME类型',
    extension       VARCHAR(10)  NOT NULL COMMENT '扩展名',
    thumbnail_path  VARCHAR(500) DEFAULT '' COMMENT '缩略图路径',
    uploaded_by     BIGINT UNSIGNED NOT NULL COMMENT '上传人ID',
    is_deleted      TINYINT      NOT NULL DEFAULT 0,
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_inventory_id (inventory_id),
    INDEX idx_file_type (file_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='附件表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id         BIGINT UNSIGNED NOT NULL COMMENT '操作人ID',
    user_name       VARCHAR(50)  NOT NULL COMMENT '操作人用户名',
    action          VARCHAR(50)  NOT NULL COMMENT '操作类型',
    target_type     VARCHAR(50)  NOT NULL COMMENT '对象类型',
    target_id       BIGINT UNSIGNED DEFAULT NULL COMMENT '对象ID',
    description     VARCHAR(500) DEFAULT '' COMMENT '操作描述',
    ip_address      VARCHAR(50)  DEFAULT '' COMMENT 'IP地址',
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- 变更日志表
CREATE TABLE IF NOT EXISTS change_logs (
    id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    inventory_id    BIGINT UNSIGNED NOT NULL COMMENT '库存记录ID',
    user_id         BIGINT UNSIGNED NOT NULL COMMENT '操作人ID',
    user_name       VARCHAR(50)  NOT NULL COMMENT '操作人用户名',
    field_name      VARCHAR(50)  NOT NULL COMMENT '变更字段名',
    field_label     VARCHAR(50)  NOT NULL COMMENT '字段中文名',
    old_value       VARCHAR(500) DEFAULT '' COMMENT '旧值',
    new_value       VARCHAR(500) DEFAULT '' COMMENT '新值',
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_inventory_id (inventory_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='变更日志表';

-- 插入默认管理员账号（密码: admin123）
-- 此bcrypt哈希值对应密码 admin123
INSERT INTO users (username, password_hash, display_name, role, status)
VALUES ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '系统管理员', 'admin', 1)
ON DUPLICATE KEY UPDATE username=username;

INSERT INTO users (username, password_hash, display_name, role, status)
VALUES ('operator', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '操作员张三', 'operator', 1)
ON DUPLICATE KEY UPDATE username=username;

INSERT INTO users (username, password_hash, display_name, role, status)
VALUES ('viewer', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '查看者李四', 'viewer', 1)
ON DUPLICATE KEY UPDATE username=username;
