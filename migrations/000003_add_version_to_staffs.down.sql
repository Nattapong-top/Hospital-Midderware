-- Drop version columns from staffs table
ALTER TABLE staffs DROP COLUMN IF EXISTS previous_version;
ALTER TABLE staffs DROP COLUMN IF EXISTS current_version;
