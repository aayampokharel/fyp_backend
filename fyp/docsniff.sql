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


-- Blocks table (header only)
CREATE TABLE blocks (
    block_number INTEGER PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    previous_hash VARCHAR(255) NOT NULL,
    nonce VARCHAR(255) NOT NULL,
    current_hash VARCHAR(255) UNIQUE NOT NULL,
    merkle_root VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Certificates table (all 4 certificates linked to block)
CREATE TABLE certificates (
    id SERIAL PRIMARY KEY,
    certificate_id VARCHAR(255) NOT NULL,
    block_number INTEGER NOT NULL,
    position INTEGER NOT NULL CHECK (position BETWEEN 1 AND 4),
    student_id VARCHAR(255) NOT NULL,
    student_name VARCHAR(255) NOT NULL,
    university_name VARCHAR(255) NOT NULL,
    degree VARCHAR(100) NOT NULL,
    college VARCHAR(255) NOT NULL,
    major VARCHAR(255) NOT NULL,
    gpa VARCHAR(10),
    percentage DECIMAL(5,2),
    division VARCHAR(50) NOT NULL,
    issue_date TIMESTAMP NOT NULL,
    enrollment_date TIMESTAMP NOT NULL,
    completion_date TIMESTAMP NOT NULL,
    principal_signature VARCHAR(255) NOT NULL,
    data_hash VARCHAR(255) NOT NULL,
    issuer_public_key VARCHAR(255) NOT NULL,
    certificate_type VARCHAR(50) NOT NULL,
    
    FOREIGN KEY (block_number) REFERENCES blocks(block_number),
    UNIQUE(block_number, position) -- Ensures exactly 4 per block
);