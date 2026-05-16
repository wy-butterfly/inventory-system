-- Supabase PostgreSQL 数据库结构
-- 适配 PostgreSQL 语法

-- 启用 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'viewer',
    status INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 库存表
CREATE TABLE IF NOT EXISTS inventories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT,
    quantity INTEGER NOT NULL DEFAULT 0,
    unit VARCHAR(20) NOT NULL DEFAULT '个',
    min_quantity INTEGER NOT NULL DEFAULT 0,
    max_quantity INTEGER NOT NULL DEFAULT 0,
    location VARCHAR(100),
    supplier VARCHAR(100),
    purchase_price DECIMAL(10,2) DEFAULT 0.00,
    sale_price DECIMAL(10,2) DEFAULT 0.00,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    version INTEGER NOT NULL DEFAULT 1,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 附件表
CREATE TABLE IF NOT EXISTS attachments (
    id SERIAL PRIMARY KEY,
    inventory_id INTEGER NOT NULL REFERENCES inventories(id) ON DELETE CASCADE,
    display_name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(20) NOT NULL DEFAULT 'document',
    extension VARCHAR(10) NOT NULL,
    mime_type VARCHAR(100),
    thumbnail_path VARCHAR(500),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    user_name VARCHAR(50) NOT NULL,
    action VARCHAR(20) NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    target_id INTEGER,
    description TEXT,
    ip_address VARCHAR(50) NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 变更日志表
CREATE TABLE IF NOT EXISTS change_logs (
    id SERIAL PRIMARY KEY,
    inventory_id INTEGER NOT NULL REFERENCES inventories(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    user_name VARCHAR(50) NOT NULL,
    action VARCHAR(20) NOT NULL,
    field_name VARCHAR(50) NOT NULL,
    old_value TEXT,
    new_value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_inventories_code ON inventories(code);
CREATE INDEX IF NOT EXISTS idx_inventories_category ON inventories(category);
CREATE INDEX IF NOT EXISTS idx_inventories_status ON inventories(status);
CREATE INDEX IF NOT EXISTS idx_inventories_is_deleted ON inventories(is_deleted);
CREATE INDEX IF NOT EXISTS idx_attachments_inventory_id ON attachments(inventory_id);
CREATE INDEX IF NOT EXISTS idx_attachments_is_deleted ON attachments(is_deleted);
CREATE INDEX IF NOT EXISTS idx_operation_logs_user_id ON operation_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_created_at ON operation_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_change_logs_inventory_id ON change_logs(inventory_id);
CREATE INDEX IF NOT EXISTS idx_change_logs_user_id ON change_logs(user_id);

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为需要的表添加更新时间触发器
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_inventories_updated_at BEFORE UPDATE ON inventories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_attachments_updated_at BEFORE UPDATE ON attachments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 插入初始数据
INSERT INTO users (username, password_hash, display_name, role) VALUES
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '系统管理员', 'admin'),
('operator', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '操作员', 'operator'),
('viewer', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '查看者', 'viewer')
ON CONFLICT (username) DO NOTHING;

-- 插入示例库存数据
INSERT INTO inventories (name, code, category, description, quantity, unit, min_quantity, max_quantity, location, supplier, purchase_price, sale_price) VALUES
('笔记本电脑', 'NB-001', '电子产品', 'ThinkPad T14 商务笔记本', 10, '台', 2, 50, '仓库A', '联想官方', 5800.00, 6800.00),
('无线鼠标', 'WM-002', '电子产品', '罗技 MX Master 3', 25, '个', 5, 100, '仓库B', '罗技官方', 699.00, 899.00),
('A4纸', 'PA-003', '办公用品', '得力 A4 复印纸 500张/包', 100, '包', 20, 200, '仓库C', '得力文具', 18.50, 25.00),
('办公椅', 'OC-004', '家具', '人体工学办公椅', 15, '把', 3, 30, '仓库D', '家具厂', 450.00, 650.00)
ON CONFLICT (code) DO NOTHING;
