CREATE TABLE levels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    annual_leave_days INT NOT NULL,
    minimum_leave_days INT NOT NULL DEFAULT 0,
    total_annual_leave_days INT NOT NULL DEFAULT 0,
    leave_expiration_interval INT NOT NULL DEFAULT 0,
    annual_gross NUMERIC(15,2) DEFAULT 0,
    basic_salary NUMERIC(15,2) DEFAULT 0,
    transport_allowance NUMERIC(15,2) DEFAULT 0,
    domestic_allowance NUMERIC(15,2) DEFAULT 0,
    utility_allowance NUMERIC(15,2) DEFAULT 0,
    lunch_subsidy NUMERIC(15,2) DEFAULT 0,
    support_total NUMERIC(15,2) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE level_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    level_id UUID REFERENCES levels(id) ON DELETE CASCADE,
    component_name VARCHAR(100) NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    category VARCHAR(50)
);
