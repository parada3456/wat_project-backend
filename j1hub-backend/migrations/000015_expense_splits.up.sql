CREATE TABLE expense_splits (
    split_id TEXT PRIMARY KEY,
    transaction_id TEXT NOT NULL REFERENCES expense_transactions(transaction_id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    owe_amount DECIMAL(12,2) NOT NULL,
    payment_status TEXT NOT NULL CHECK (payment_status IN ('Pending', 'Submitted', 'Approved', 'Overdue')),
    payment_method TEXT,
    payslip_url TEXT,
    approval_status TEXT NOT NULL CHECK (approval_status IN ('Pending_Approval', 'Approved', 'Rejected')),
    settled_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
