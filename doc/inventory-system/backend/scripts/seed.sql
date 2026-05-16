-- 插入测试库存数据
USE inventory_db;

INSERT INTO inventories (material_code, material_name, category, specification, unit, quantity, safety_stock, location, status, created_by) VALUES
('YL-202603-0001', '铝合金板材', 'raw_material', '1000x2000x2mm', '张', 500.00, 100.00, 'A区-01-01', 1, 1),
('YL-202603-0002', '不锈钢管', 'raw_material', 'DN50 壁厚3mm', '米', 200.00, 50.00, 'A区-01-02', 1, 1),
('CP-202603-0001', '电控柜成品', 'finished_product', 'GGD型 800x600x2200', '台', 15.00, 5.00, 'B区-02-01', 1, 1),
('CP-202603-0002', '配电箱', 'finished_product', 'XL-21型', '台', 30.00, 10.00, 'B区-02-02', 1, 1),
('BJ-202603-0001', 'ABB断路器', 'spare_part', 'S203-C32', '个', 80.00, 20.00, 'C区-03-01', 1, 1),
('BJ-202603-0002', '继电器', 'spare_part', 'MY2NJ DC24V', '个', 3.00, 50.00, 'C区-03-02', 1, 1),
('QT-202603-0001', '标签纸', 'other', 'A4热敏纸', '包', 25.00, 10.00, 'D区-04-01', 1, 1);
