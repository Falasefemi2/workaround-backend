CREATE TABLE leave_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    max_days INT NOT NULL,
    description TEXT
);

CREATE TABLE leave_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID REFERENCES employees(id) ON DELETE CASCADE,
    leave_type_id UUID REFERENCES leave_types(id),
    days INT NOT NULL,
    start_date DATE NOT NULL,
    resumption_date DATE NOT NULL,
    relief_officer_id UUID REFERENCES employees(id),
    handover_note_url TEXT,
    leave_allowance BOOLEAN DEFAULT FALSE,
    status VARCHAR(50) DEFAULT 'initiated',
    current_approval_level INT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_leave_employee ON leave_requests(employee_id);
