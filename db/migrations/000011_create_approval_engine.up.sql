CREATE TABLE approval_setups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    module_type VARCHAR(50), -- leave, memo, voucher, exit
    department_id UUID REFERENCES departments(id),
    level_order INT NOT NULL,
    role_id UUID REFERENCES roles(id)
);

CREATE TABLE approvals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    module_type VARCHAR(50),
    reference_id UUID NOT NULL,
    approval_level INT,
    approver_id UUID REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'pending',
    comment TEXT,
    acted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_approval_reference ON approvals(reference_id);
