CREATE TABLE candidates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(150) UNIQUE NOT NULL,
    phone VARCHAR(20),
    status VARCHAR(50) DEFAULT 'pending'
);

CREATE TABLE offers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    candidate_id UUID REFERENCES candidates(id) ON DELETE CASCADE,
    department_id UUID REFERENCES departments(id),
    designation_id UUID REFERENCES designations(id),
    level_id UUID REFERENCES levels(id),
    proposed_start_date DATE,
    new_start_date DATE,
    offer_letter_url TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
