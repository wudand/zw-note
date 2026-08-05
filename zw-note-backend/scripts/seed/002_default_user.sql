-- Seed: 笔记业务默认用户（仅用于 dev/test）
-- username: admin
-- password: password
-- 生产环境务必修改密码

INSERT INTO users (username, password_hash, display_name, status)
VALUES (
    'admin',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    '管理员',
    1
)
ON CONFLICT (username) DO NOTHING;
