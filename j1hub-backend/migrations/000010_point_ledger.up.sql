CREATE TABLE point_ledger (
    ledger_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    source_type TEXT NOT NULL CHECK (source_type IN ('Mission_Base', 'Speed_Bonus', 'Streak_Bonus', 'First_Completer', 'Expense_Penalty', 'Admin_Adjust')),
    source_id TEXT,
    delta INT NOT NULL,
    lifetime_balance_after INT NOT NULL,
    phase_balance_after INT NOT NULL,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
