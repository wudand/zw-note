-- Seed: 仅用于 dev/test 环境
-- password: Admin@123456（superadmin）/ password（admin、viewer）— 生产环境务必修改

INSERT INTO admin_users (username, password, role, status)
VALUES
    ('superadmin', '$2a$10$ZAlyqcQT9IF.7w.si9mE2ueSK9v867iotA4CvPY4KFmB1jdCgRZty', 'super_admin', 1),
    ('admin',      '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin',       1),
    ('viewer',     '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'viewer',      1)
ON CONFLICT (username) DO NOTHING;
