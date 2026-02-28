CREATE TABLE resignations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID REFERENCES employees(id),
    reason TEXT,
    disengagement_date DATE,
    letter_url TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE handovers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resignation_id UUID REFERENCES resignations(id) ON DELETE CASCADE,
    document_url TEXT,
    status VARCHAR(50) DEFAULT 'pending'
);

CREATE TABLE exit_interviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID REFERENCES employees(id),
    response JSONB,
    submitted_at TIMESTAMPTZ DEFAULT NOW()
);
