-- Standardize DB constraints to lowercase snake_case and current_coordinates to GEOGRAPHY

-- 1. profiles table
ALTER TABLE profiles DROP CONSTRAINT IF EXISTS profiles_radar_visibility_check;
UPDATE profiles SET radar_visibility = LOWER(radar_visibility);
ALTER TABLE profiles ADD CONSTRAINT profiles_radar_visibility_check 
    CHECK (radar_visibility IN ('show_anonymous', 'show_friends', 'hidden'));

-- Alter profiles.current_coordinates type from GEOMETRY to GEOGRAPHY
ALTER TABLE profiles ALTER COLUMN current_coordinates TYPE GEOGRAPHY(POINT, 4326) USING current_coordinates::geography;

-- 2. friendships table
ALTER TABLE friendships DROP CONSTRAINT IF EXISTS friendships_status_check;
UPDATE friendships SET status = LOWER(status);
ALTER TABLE friendships ADD CONSTRAINT friendships_status_check 
    CHECK (status IN ('pending', 'accepted', 'blocked'));

-- 3. missions table
ALTER TABLE missions DROP CONSTRAINT IF EXISTS missions_verification_type_check;
ALTER TABLE missions DROP CONSTRAINT IF EXISTS missions_due_date_type_check;
UPDATE missions SET verification_type = LOWER(verification_type), due_date_type = LOWER(due_date_type);
ALTER TABLE missions ADD CONSTRAINT missions_verification_type_check 
    CHECK (verification_type IN ('none', 'upload', 'admin'));
ALTER TABLE missions ADD CONSTRAINT missions_due_date_type_check 
    CHECK (due_date_type IN ('relative', 'fixed'));

-- 4. user_missions table
ALTER TABLE user_missions DROP CONSTRAINT IF EXISTS user_missions_status_check;
UPDATE user_missions SET status = LOWER(status);
ALTER TABLE user_missions ADD CONSTRAINT user_missions_status_check 
    CHECK (status IN ('not_started', 'in_progress', 'pending_verification', 'completed', 'overdue'));

-- 5. point_ledger table
ALTER TABLE point_ledger DROP CONSTRAINT IF EXISTS point_ledger_source_type_check;
UPDATE point_ledger SET source_type = LOWER(source_type);
ALTER TABLE point_ledger ADD CONSTRAINT point_ledger_source_type_check 
    CHECK (source_type IN ('mission_base', 'speed_bonus', 'streak_bonus', 'first_completer', 'expense_penalty', 'admin_adjust'));

-- 6. badges table
ALTER TABLE badges DROP CONSTRAINT IF EXISTS badges_trigger_type_check;
UPDATE badges SET trigger_type = LOWER(trigger_type);
ALTER TABLE badges ADD CONSTRAINT badges_trigger_type_check 
    CHECK (trigger_type IN ('speed', 'streak', 'first_completer', 'phase_complete', 'manual'));

-- 7. expense_splits table
ALTER TABLE expense_splits DROP CONSTRAINT IF EXISTS expense_splits_payment_status_check;
ALTER TABLE expense_splits DROP CONSTRAINT IF EXISTS expense_splits_approval_status_check;
UPDATE expense_splits SET payment_status = LOWER(payment_status), approval_status = LOWER(approval_status);
ALTER TABLE expense_splits ADD CONSTRAINT expense_splits_payment_status_check 
    CHECK (payment_status IN ('pending', 'submitted', 'approved', 'overdue'));
ALTER TABLE expense_splits ADD CONSTRAINT expense_splits_approval_status_check 
    CHECK (approval_status IN ('pending_approval', 'approved', 'rejected'));

-- 8. user_carts table
ALTER TABLE user_carts DROP CONSTRAINT IF EXISTS user_carts_status_check;
UPDATE user_carts SET status = LOWER(status);
ALTER TABLE user_carts ADD CONSTRAINT user_carts_status_check 
    CHECK (status IN ('saved', 'viewed', 'applied', 'removed'));
