-- flowershop 初始化建表脚本（后端启动时也会通过 GORM AutoMigrate 自动建表）
-- 此脚本供手动初始化参考。

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    nickname VARCHAR(64),
    email VARCHAR(128),
    phone VARCHAR(32),
    avatar VARCHAR(255),
    role VARCHAR(16) DEFAULT 'USER',
    status VARCHAR(16) DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT DEFAULT 0,
    name VARCHAR(64) NOT NULL,
    level INT DEFAULT 1,
    sort INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL,
    name VARCHAR(128) NOT NULL,
    sub_title VARCHAR(255),
    description TEXT,
    detail TEXT,
    cover_image VARCHAR(255),
    images TEXT,
    price NUMERIC(10,2) NOT NULL,
    original_price NUMERIC(10,2),
    stock INT DEFAULT 0,
    sales INT DEFAULT 0,
    rating NUMERIC(3,2) DEFAULT 5,
    rating_count INT DEFAULT 0,
    shipping_from VARCHAR(64),
    free_shipping BOOLEAN DEFAULT FALSE,
    status VARCHAR(16) DEFAULT 'ON_SALE',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(32) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    status VARCHAR(24) DEFAULT 'PENDING_PAYMENT',
    total_amount NUMERIC(10,2) NOT NULL,
    discount_amount NUMERIC(10,2) DEFAULT 0,
    shipping_fee NUMERIC(10,2) DEFAULT 0,
    pay_amount NUMERIC(10,2) NOT NULL,
    coupon_id BIGINT DEFAULT 0,
    receiver_name VARCHAR(64) NOT NULL,
    receiver_phone VARCHAR(32) NOT NULL,
    receiver_addr VARCHAR(255) NOT NULL,
    remark VARCHAR(255),
    paid_at TIMESTAMPTZ,
    shipped_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    cancelled_by VARCHAR(32),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    product_name VARCHAR(128) NOT NULL,
    product_image VARCHAR(255),
    price NUMERIC(10,2) NOT NULL,
    quantity INT NOT NULL,
    total_price NUMERIC(10,2) NOT NULL,
    reviewed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    order_no VARCHAR(32) NOT NULL,
    pay_no VARCHAR(48) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    channel VARCHAR(32) DEFAULT 'ALIPAY_SANDBOX',
    amount NUMERIC(10,2) NOT NULL,
    status VARCHAR(16) DEFAULT 'PENDING',
    transaction_id VARCHAR(64),
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS logistics (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL UNIQUE,
    order_no VARCHAR(32) NOT NULL,
    tracking_no VARCHAR(64) NOT NULL UNIQUE,
    carrier VARCHAR(32) DEFAULT 'SF_EXPRESS',
    status VARCHAR(24) DEFAULT 'PENDING',
    events TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reviews (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    order_item_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    rating INT NOT NULL,
    content TEXT,
    images TEXT,
    reply TEXT,
    replied_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS coupons (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    template_id BIGINT DEFAULT 0,
    name VARCHAR(64) NOT NULL,
    type VARCHAR(24) NOT NULL,
    threshold NUMERIC(10,2) DEFAULT 0,
    amount NUMERIC(10,2) DEFAULT 0,
    discount_rate NUMERIC(4,2) DEFAULT 0,
    status VARCHAR(16) DEFAULT 'UNUSED',
    used_at TIMESTAMPTZ,
    order_id BIGINT DEFAULT 0,
    valid_from TIMESTAMPTZ,
    valid_to TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    username VARCHAR(64),
    action VARCHAR(24),
    method VARCHAR(16),
    path VARCHAR(255),
    entity VARCHAR(64),
    entity_id VARCHAR(64),
    detail TEXT,
    request_id VARCHAR(64),
    ip VARCHAR(64),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
