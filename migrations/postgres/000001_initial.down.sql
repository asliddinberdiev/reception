DROP TABLE IF EXISTS "notifications";

DROP TYPE IF EXISTS "notification_status";

DROP TABLE IF EXISTS "appointments";

DROP TYPE IF EXISTS "appointment_status";

DROP TABLE IF EXISTS "patients";

DROP TABLE IF EXISTS "user_work_times";

DROP TABLE IF EXISTS "user_roles";

DROP INDEX IF EXISTS "patient_phone_deleted_idx";
DROP INDEX IF EXISTS "patient_id_deleted_idx";

DROP INDEX IF EXISTS "user_specialty_idx";
DROP INDEX IF EXISTS "user_phone_deleted_idx";
DROP INDEX IF EXISTS "user_id_deleted_idx";

DROP TABLE IF EXISTS "users";

DROP TABLE IF EXISTS "roles";
DROP TYPE IF EXISTS "gender_type";

DROP EXTENSION IF EXISTS "uuid-ossp";