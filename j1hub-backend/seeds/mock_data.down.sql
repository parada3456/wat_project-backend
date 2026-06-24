-- 000023_seed_mock_data.down.sql

-- Clear junction, tracking, and transactional data first to avoid FK constraint violations
DELETE FROM user_jobs WHERE user_id IN ('usr_somchai', 'usr_alice');
DELETE FROM job_reviews WHERE review_id IN ('rev_001');
DELETE FROM job_overall_ratings WHERE rating_summary_id IN ('rtg_job_001');
DELETE FROM user_carts WHERE cart_id IN ('crt_som_job_001', 'crt_som_job_002');
DELETE FROM job_housings WHERE housing_id IN ('hsg_001', 'hsg_002');
DELETE FROM job_postings WHERE job_id IN ('job_001', 'job_002', 'job_003');
DELETE FROM expense_splits WHERE split_id IN ('spl_001');
DELETE FROM expense_transactions WHERE transaction_id IN ('tx_001');
DELETE FROM credit_scores WHERE credit_id IN ('crd_admin', 'crd_somchai', 'crd_alice');
DELETE FROM user_badges WHERE user_badge_id IN ('ubg_som_1', 'ubg_som_2', 'ubg_ali_1');
DELETE FROM badges WHERE badge_id IN ('bdg_001', 'bdg_002', 'bdg_003');
DELETE FROM point_ledger WHERE ledger_id LIKE 'plg_%';
DELETE FROM user_tasks WHERE user_task_id LIKE 'utk_%';
DELETE FROM tasks WHERE task_id IN ('tsk_001', 'tsk_002', 'tsk_003');
DELETE FROM user_missions WHERE user_mission_id LIKE 'ums_%';
DELETE FROM missions WHERE mission_id IN ('mis_001', 'mis_002', 'mis_003', 'mis_004', 'mis_005');
DELETE FROM user_phase_history WHERE history_id LIKE 'uph_%';
DELETE FROM friendships WHERE friendship_id IN ('frn_somchai_alice');
DELETE FROM profiles WHERE profile_id IN ('prf_admin', 'prf_somchai', 'prf_alice');
DELETE FROM users WHERE user_id IN ('usr_admin', 'usr_somchai', 'usr_alice');
DELETE FROM journey_phases WHERE phase_id IN ('phs_001', 'phs_002', 'phs_003', 'phs_004');