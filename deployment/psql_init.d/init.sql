CREATE TABLE schemes (
    scheme_id SERIAL PRIMARY KEY,
    scheme_code VARCHAR(16) UNIQUE NOT NULL,
    scheme_name VARCHAR(100) NOT NULL
);

CREATE TABLE branches (
    branch_id SERIAL PRIMARY KEY,
    branch_code VARCHAR(16) UNIQUE NOT NULL,
    branch_name VARCHAR(100) NOT NULL
);

CREATE TABLE subjects (
    subject_id SERIAL PRIMARY KEY,
    subject_code VARCHAR(12) UNIQUE NOT NULL,
    subject_name VARCHAR(100) NOT NULL
);

CREATE TABLE subject_offerings (
    offering_id SERIAL PRIMARY KEY,
    subject_id INT NOT NULL REFERENCES subjects(subject_id) ON DELETE CASCADE,
    scheme_id INT NOT NULL REFERENCES schemes(scheme_id) ON DELETE CASCADE,
    branch_id INT NOT NULL REFERENCES branches(branch_id) ON DELETE CASCADE,
    semester INT NOT NULL,
    UNIQUE(subject_id, scheme_id, branch_id, semester)
);

CREATE TABLE normal_users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    email VARCHAR(128) UNIQUE NOT NULL,
    password_hash CHAR(60) NOT NULL, -- for bcrypt/argon2
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE materials (
    material_id SERIAL PRIMARY KEY,
    offering_id INT NOT NULL REFERENCES subject_offerings(offering_id) ON DELETE CASCADE,
    uploaded_by_user INT NOT NULL REFERENCES normal_users(user_id) ON DELETE SET NULL,
    title VARCHAR(256) NOT NULL,
    file_path TEXT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

