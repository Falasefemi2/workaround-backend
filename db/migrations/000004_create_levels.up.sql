CREATE TABLE levels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    annual_leave_days INT NOT NULL,
    annual_gross NUMERIC(15,2) DEFAULT 0,
    support_total NUMERIC(15,2) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE level_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    level_id UUID REFERENCES levels(id) ON DELETE CASCADE,
    component_name VARCHAR(100) NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    category VARCHAR(50) -- gross, allowance, support
);
