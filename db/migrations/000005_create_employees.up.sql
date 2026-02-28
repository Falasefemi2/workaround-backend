CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    employee_number VARCHAR(50) UNIQUE NOT NULL,
    department_id UUID REFERENCES departments(id),
    unit_id UUID REFERENCES units(id),
    designation_id UUID REFERENCES designations(id),
    level_id UUID REFERENCES levels(id),
    employment_type VARCHAR(50),
    employment_status VARCHAR(50) DEFAULT 'active',
    date_of_employment DATE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_employees_department ON employees(department_id);
CREATE INDEX idx_employees_number ON employees(employee_number);
