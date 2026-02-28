CREATE TABLE vouchers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    voucher_code VARCHAR(100) UNIQUE NOT NULL,
    department_id UUID REFERENCES departments(id),
    amount NUMERIC(15,2),
    bank VARCHAR(150),
    created_by UUID REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'pending',
    current_stage INT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
