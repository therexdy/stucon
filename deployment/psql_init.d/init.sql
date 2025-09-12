CREATE DATABASE stucon;

\c stucon

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE schemes (
    scheme_id TEXT PRIMARY KEY
);

CREATE TABLE branches (
    branch_id TEXT PRIMARY KEY,
    branch_name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE subjects (
    subject_id TEXT PRIMARY KEY,
    scheme_id TEXT NOT NULL REFERENCES schemes(scheme_id) ON DELETE CASCADE,
    branch_id TEXT NOT NULL REFERENCES branches(branch_id) ON DELETE CASCADE,
    sem INT NOT NULL,
    subject_name VARCHAR(100) NOT NULL,
    UNIQUE (scheme_id, branch_id, sem, subject_name)
);

CREATE TABLE users (
    user_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE materials (
    material_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    uploaded_by_user INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    scheme_id TEXT NOT NULL REFERENCES schemes(scheme_id) ON DELETE CASCADE,
    branch_id TEXT NOT NULL REFERENCES branches(branch_id) ON DELETE CASCADE,
    subject_id TEXT NOT NULL REFERENCES subjects(subject_id) ON DELETE CASCADE,
    sem INT NOT NULL CHECK (sem BETWEEN 1 AND 8),
    title VARCHAR(256) NOT NULL,
    file_path TEXT NOT NULL,
    file_type TEXT NOT NULL,
    uploaded_at TIMESTAMPZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);



CREATE VIEW materials_with_user AS
SELECT 
    m.material_id,
    u.name AS uploaded_by,
    m.scheme_id,
    m.branch_id,
    m.subject_id,
    m.sem,
    m.title,
    m.file_path,
    m.file_type,
    m.uploaded_at
FROM materials m
JOIN users u ON m.uploaded_by_user = u.user_id;


CREATE VIEW subjects_with_branch AS
SELECT 
    s.subject_id,
    s.subject_name,
    s.scheme_id,
    s.sem,
    b.branch_id AS branch
FROM subjects s
JOIN branches b ON s.branch_id = b.branch_id;

CREATE USER appuser WITH PASSWORD 'GTAC';

GRANT CONNECT, TEMPORARY ON DATABASE stucon TO appuser;
GRANT USAGE ON SCHEMA public TO appuser;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO appuser;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO appuser;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO appuser;

