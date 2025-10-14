CREATE DATABASE docsniff_db;
\c docsniff_db;  

CREATE TABLE institutions(
    institution_id VARCHAR(16) PRIMARY KEY,
    institution_name VARCHAR(300) NOT NULL,
    tole_address VARCHAR(250) NOT NULL,
    district_address VARCHAR(250) NOT NULL,
	is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE user_accounts(
    id VARCHAR(16) PRIMARY KEY,
    role VARCHAR(16) NOT NULL,  -- Added NOT NULL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    email VARCHAR(255) UNIQUE NOT NULL, 
    password VARCHAR(255) NOT NULL
);

CREATE TABLE institution_user(
    institution_id VARCHAR(16)  UNIQUE REFERENCES institutions(institution_id) ON DELETE CASCADE,
    user_id VARCHAR(16) REFERENCES user_accounts(id) ON DELETE CASCADE,
    public_key TEXT DEFAULT NULL,  -- Changed from VARCHAR(100) to TEXT
    principal_name VARCHAR(300) NOT NULL,
    principal_signature_base64 TEXT NOT NULL,
    institution_logo_base64 TEXT NOT NULL,
    PRIMARY KEY (institution_id, user_id) 
);

CREATE TABLE institution_faculty(
    institution_faculty_id VARCHAR(16) PRIMARY KEY,
    institution_id VARCHAR(16) REFERENCES institutions(institution_id) ON DELETE CASCADE,
    faculty VARCHAR(200) NOT NULL,
    faculty_hod_name VARCHAR(300) NOT NULL,
    faculty_hod_signature_base64 TEXT NOT NULL
);


