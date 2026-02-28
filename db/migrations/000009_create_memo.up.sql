CREATE TABLE memos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    memo_code VARCHAR(100) UNIQUE NOT NULL,
    department_id UUID REFERENCES departments(id),
    title VARCHAR(255),
    notes TEXT,
    amount NUMERIC(15,2),
    beneficiary VARCHAR(255),
    memo_type VARCHAR(50), -- general, asset
    created_by UUID REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'pending',
    current_stage INT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_memo_code ON memos(memo_code);

CREATE TABLE memo_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    memo_id UUID REFERENCES memos(id) ON DELETE CASCADE,
    document_name VARCHAR(255),
    file_url TEXT
);
