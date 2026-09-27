-- เพิ่มคอลัมน์ current_version และ previous_version ในตาราง staffs สำหรับ Version Value Object (Optimistic Locking / Conflict Handling)
ALTER TABLE staffs ADD COLUMN IF NOT EXISTS current_version INT NOT NULL DEFAULT 1;
ALTER TABLE staffs ADD COLUMN IF NOT EXISTS previous_version INT NOT NULL DEFAULT 1;
