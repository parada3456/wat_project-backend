-- Revert standardized constraints to original casing and current_coordinates to GEOMETRY

-- 1. profiles table
ALTER TABLE profiles DROP CONSTRAINT IF EXISTS profiles_radar_visibility_check;
UPDATE profiles SET radar_visibility = CASE 
    WHEN radar_visibility = 'show_anonymous' THEN 'Show_Anonymous'
    WHEN radar_visibility = 'show_friends' THEN 'Show_Friends'
    WHEN radar_visibility = 'hidden' THEN 'Hidden'
    ELSE radar_visibility
END;
ALTER TABLE profiles ADD CONSTRAINT profiles_radar_visibility_check 
    CHECK (radar_visibility IN ('Show_Anonymous', 'Show_Friends', 'Hidden'));

-- Revert profiles.current_coordinates type from GEOGRAPHY to GEOMETRY
ALTER TABLE profiles ALTER COLUMN current_coordinates TYPE GEOMETRY(POINT, 4326) USING current_coordinates::geometry;

-- 2. friendships table
ALTER TABLE friendships DROP CONSTRAINT IF EXISTS friendships_status_check;
UPDATE friendships SET status = CASE 
    WHEN status = 'pending' THEN 'Pending'
    WHEN status = 'accepted' THEN 'Accepted'
    WHEN status = 'blocked' THEN 'Blocked'
    ELSE status
END;
ALTER TABLE friendships ADD CONSTRAINT friendships_status_check 
    CHECK (status IN ('Pending', 'Accepted', 'Blocked'));

-- 3. missions table
ALTER TABLE missions DROP CONSTRAINT IF EXISTS missions_verification_type_check;
ALTER TABLE missions DROP CONSTRAINT IF EXISTS missions_due_date_type_check;
UPDATE missions SET verification_type = CASE 
    WHEN verification_type = 'none' THEN 'None'
    WHEN verification_type = 'upload' THEN 'Upload'
    WHEN verification_type = 'admin' THEN 'Admin'
    ELSE verification_type
END, due_date_type = CASE
    WHEN due_date_type = 'relative' THEN 'Relative'
    WHEN due_date_type = 'fixed' THEN 'Fixed'
    ELSE due_date_type
END;
ALTER TABLE missions ADD CONSTRAINT missions_verification_type_check 
    CHECK (verification_type IN ('None', 'Upload', 'Admin'));
ALTER TABLE missions ADD CONSTRAINT missions_due_date_type_check 
    CHECK (due_date_type IN ('Relative', 'Fixed'));

-- 4. user_missions table
ALTER TABLE user_missions DROP CONSTRAINT IF EXISTS user_missions_status_check;
UPDATE user_missions SET status = CASE 
    WHEN status = 'not_started' THEN 'Not_Started'
    WHEN status = 'in_progress' THEN 'In_Progress'
    WHEN status = 'pending_verification' THEN 'Pending_Verification'
    WHEN status = 'completed' THEN 'Completed'
    WHEN status = 'overdue' THEN 'Overdue'
    ELSE status
END;
ALTER TABLE user_missions ADD CONSTRAINT user_missions_status_check 
    CHECK (status IN ('Not_Started', 'In_Progress', 'Pending_Verification', 'Completed', 'Overdue'));

-- 5. point_ledger table
ALTER TABLE point_ledger DROP CONSTRAINT IF EXISTS point_ledger_source_type_check;
UPDATE point_ledger SET source_type = CASE 
    WHEN source_type = 'mission_base' THEN 'Mission_Base'
    WHEN source_type = 'speed_bonus' THEN 'Speed_Bonus'
    WHEN source_type = 'streak_bonus' THEN 'Streak_Bonus'
    WHEN source_type = 'first_completer' THEN 'First_Completer'
    WHEN source_type = 'expense_penalty' THEN 'Expense_Penalty'
    WHEN source_type = 'admin_adjust' THEN 'Admin_Adjust'
    ELSE source_type
END;
ALTER TABLE point_ledger ADD CONSTRAINT point_ledger_source_type_check 
    CHECK (source_type IN ('Mission_Base', 'Speed_Bonus', 'Streak_Bonus', 'First_Completer', 'Expense_Penalty', 'Admin_Adjust'));

-- 6. badges table
ALTER TABLE badges DROP CONSTRAINT IF EXISTS badges_trigger_type_check;
UPDATE badges SET trigger_type = CASE 
    WHEN trigger_type = 'speed' THEN 'Speed'
    WHEN trigger_type = 'streak' THEN 'Streak'
    WHEN trigger_type = 'first_completer' THEN 'First_Completer'
    WHEN trigger_type = 'phase_complete' THEN 'Phase_Complete'
    WHEN trigger_type = 'manual' THEN 'Manual'
    ELSE trigger_type
END;
ALTER TABLE badges ADD CONSTRAINT badges_trigger_type_check 
    CHECK (trigger_type IN ('Speed', 'Streak', 'First_Completer', 'Phase_Complete', 'Manual'));

-- 7. expense_splits table
ALTER TABLE expense_splits DROP CONSTRAINT IF EXISTS expense_splits_payment_status_check;
ALTER TABLE expense_splits DROP CONSTRAINT IF EXISTS expense_splits_approval_status_check;
UPDATE expense_splits SET payment_status = CASE 
    WHEN payment_status = 'pending' THEN 'Pending'
    WHEN payment_status = 'submitted' THEN 'Submitted'
    WHEN payment_status = 'approved' THEN 'Approved'
    WHEN payment_status = 'overdue' THEN 'Overdue'
    ELSE payment_status
END, approval_status = CASE
    WHEN approval_status = 'pending_approval' THEN 'Pending_Approval'
    WHEN approval_status = 'approved' THEN 'Approved'
    WHEN approval_status = 'rejected' THEN 'Rejected'
    ELSE approval_status
END;
ALTER TABLE expense_splits ADD CONSTRAINT expense_splits_payment_status_check 
    CHECK (payment_status IN ('Pending', 'Submitted', 'Approved', 'Overdue'));
ALTER TABLE expense_splits ADD CONSTRAINT expense_splits_approval_status_check 
    CHECK (approval_status IN ('Pending_Approval', 'Approved', 'Rejected'));

-- 8. user_carts table
ALTER TABLE user_carts DROP CONSTRAINT IF EXISTS user_carts_status_check;
UPDATE user_carts SET status = CASE 
    WHEN status = 'saved' THEN 'Saved'
    WHEN status = 'viewed' THEN 'Viewed'
    WHEN status = 'applied' THEN 'Applied'
    WHEN status = 'removed' THEN 'Removed'
    ELSE status
END;
ALTER TABLE user_carts ADD CONSTRAINT user_carts_status_check 
    CHECK (status IN ('Saved', 'Viewed', 'Applied', 'Removed'));
