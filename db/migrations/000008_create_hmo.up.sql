CREATE TABLE hmos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    new_employee_form_url TEXT,
    existing_employee_form_url TEXT
);

CREATE TABLE employee_hmos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    employee_id UUID REFERENCES employees(id) ON DELETE CASCADE,
    hmo_id UUID REFERENCES hmos(id),
    form_url TEXT,
    status VARCHAR(50) DEFAULT 'submitted'
);
