
CREATE DATABASE stucon;

\c stucon


CREATE TABLE scheme (
    scheme_id SERIAL PRIMARY KEY,
    scheme_name VARCHAR(100) NOT NULL
);

CREATE TABLE courses (
    course_id SERIAL PRIMARY KEY,
    course_name VARCHAR(150) NOT NULL,
    scheme_id INT REFERENCES scheme(scheme_id) ON DELETE CASCADE
);

CREATE TABLE semesters (
    semester_id SERIAL PRIMARY KEY,
    course_id INT REFERENCES courses(course_id) ON DELETE CASCADE,
    semester_number INT NOT NULL CHECK (semester_number BETWEEN 1 AND 8)
);

CREATE TABLE subjects (
    subject_id SERIAL PRIMARY KEY,
    semester_id INT REFERENCES semesters(semester_id) ON DELETE CASCADE,
    subject_name VARCHAR(100) NOT NULL
);

CREATE TABLE admin_users (
    admin_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE normal_users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE materials (
    material_id SERIAL PRIMARY KEY,
    subject_id INT REFERENCES subjects(subject_id) ON DELETE CASCADE,
    uploaded_by_type VARCHAR(10) CHECK (uploaded_by_type IN ('admin','normal')) NOT NULL,
    uploaded_by_admin INT REFERENCES admin_users(admin_id),
    uploaded_by_user INT REFERENCES normal_users(user_id),
    title VARCHAR(200) NOT NULL,
    file_path TEXT NOT NULL, 
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);


CREATE USER appuser WITH PASSWORD 'GTAC';


GRANT CONNECT ON DATABASE stucon TO appuser;
GRANT USAGE ON SCHEMA public TO appuser;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO appuser;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO appuser;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO appuser;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO appuser;

