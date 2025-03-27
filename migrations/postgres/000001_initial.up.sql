CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "gender_type" AS ENUM (
    'male',
    'female',
    'unknown'
);

CREATE TABLE IF NOT EXISTS "roles" (
    "id" UUID NOT NULL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "translation" JSONB NOT NULL
);
INSERT INTO "roles" ("id", "name", "translation") VALUES
    ('0195d786-257b-75bc-9901-52a3cf27e1ed', 'admin', '{"en":"Admin", "ru":"Админ", "uz":"Admin"}'),
    ('0195d786-5a00-71f1-be49-b7cd48c150d2', 'doctor', '{"en":"Doctor", "ru":"Врач", "uz":"Shifokor"}');

CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID NOT NULL PRIMARY KEY,
    "phone_number" VARCHAR(12) NOT NULL UNIQUE,
    "is_validated" BOOLEAN NOT NULL DEFAULT FALSE,
    "first_name" VARCHAR(255) DEFAULT '',
    "last_name" VARCHAR(255) DEFAULT '',
    "gender" gender_type NOT NULL DEFAULT 'unknown',
    "specialty" VARCHAR(255) DEFAULT '',
    "description" TEXT DEFAULT '',
    "password" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX "user_specialty_idx" ON "users" ("specialty");
CREATE INDEX "user_phone_deleted_idx" ON "users" ("phone_number", "deleted_at");
CREATE INDEX "user_id_deleted_idx" ON "users" ("id", "deleted_at");

INSERT INTO "users" ("id", "phone_number", "is_validated", "first_name", "last_name", "password") VALUES
    ('0195d790-c1fb-768f-a452-6278a74e35bd', '998331202124', TRUE, 'Asliddin', 'Berdiev', '$2y$10$sGtkwdbKL.0Ggb67m6OVBOlX77Mxb.qsyYWJH3Dd4LrFWqc1LZ/fa'), -- password: 123123
    ('0195d791-1546-7cee-8d89-60c1072d44a3', '998901112233', TRUE, 'Alisher', 'Berdiev', '$2y$10$eKFsgOgMpiJn0EHFgAj4IeQdtuPIQP33W8OetvGzua7lFtyq9umMW'); -- password: 123321


CREATE TABLE IF NOT EXISTS "user_roles" (
    "user_id" UUID NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "role_id" UUID NOT NULL REFERENCES "roles"("id") ON DELETE CASCADE,
    PRIMARY KEY ("user_id", "role_id")
);

INSERT INTO "user_roles" ("user_id", "role_id") VALUES
    ('0195d790-c1fb-768f-a452-6278a74e35bd', '0195d786-257b-75bc-9901-52a3cf27e1ed'),
    ('0195d791-1546-7cee-8d89-60c1072d44a3', '0195d786-5a00-71f1-be49-b7cd48c150d2');

CREATE TABLE IF NOT EXISTS "user_work_times" (
    "user_id" UUID NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "week_day" INT CHECK ("week_day" BETWEEN 0 AND 6),
    "start_time" TIME NOT NULL DEFAULT '09:00:00',
    "finish_time" TIME NOT NULL DEFAULT '17:00:00',
    CHECK ("start_time" < "finish_time"),
    PRIMARY KEY ("user_id", "week_day")
);

INSERT INTO "user_work_times" ("user_id", "week_day") VALUES
    ('0195d790-c1fb-768f-a452-6278a74e35bd', 0),
    ('0195d791-1546-7cee-8d89-60c1072d44a3', 1);

CREATE TABLE IF NOT EXISTS "patients" (
    "id" UUID NOT NULL PRIMARY KEY,
    "phone_number" VARCHAR(12) NOT NULL UNIQUE,
    "is_validated" BOOLEAN NOT NULL DEFAULT FALSE,
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "gender" gender_type NOT NULL DEFAULT 'unknown',
    "password" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" BIGINT NOT NULL DEFAULT 0
);

CREATE TYPE "appointment_status" AS ENUM (
    'pending',
    'confirmed',
    'rejected',
    'cancelled'
);

CREATE TABLE IF NOT EXISTS "appointments" (
    "id" UUID NOT NULL PRIMARY KEY,
    "patient_id" UUID REFERENCES "patients" ("id") ON DELETE SET NULL,
    "doctor_id" UUID REFERENCES "users" ("id") ON DELETE SET NULL,
    "appointment_date" DATE NOT NULL,
    "appointment_time" TIME NOT NULL,
    "status" appointment_status DEFAULT 'pending',
    "doctor_approval" BOOLEAN DEFAULT FALSE,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "unique_appointment" UNIQUE ("doctor_id", "appointment_date", "appointment_time")
);

CREATE TYPE "notification_status" AS ENUM (
    'pending',
    'sented'
);

CREATE TABLE IF NOT EXISTS "notifications" (
    "id" UUID NOT NULL PRIMARY KEY,
    "appointment_id" UUID REFERENCES "appointments" ("id") ON DELETE CASCADE,
    "sent_at" TIMESTAMP,
    "status" notification_status DEFAULT 'pending',
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);