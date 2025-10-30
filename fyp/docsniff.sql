
drop table if exists institution_user;
drop table if exists user_accounts;
drop table if exists pdf_files;

drop table if exists certificates;
drop table if exists pdf_file_categories;
drop table if exists blocks;
drop table if exists institution_faculty;
drop table if exists institutions;


CREATE TABLE IF NOT EXISTS  institutions(
    institution_id VARCHAR(16) PRIMARY KEY,
    institution_name VARCHAR(300) NOT NULL,
    ward_number VARCHAR(16) NOT NULL,
    tole_address VARCHAR(250) NOT NULL,
    district_address VARCHAR(250) NOT NULL,
	is_active BOOLEAN DEFAULT NULL,
    is_signup_completed BOOLEAN DEFAULT FALSE
    -- UNIQUE(institution_name, ward_number, tole_address, district_address, is_active)
-- remove unique constraint for above . user cant insert into table only if : is_active =true
);

CREATE UNIQUE INDEX unique_active_institutions 
ON institutions (institution_name, ward_number, tole_address, district_address) 
WHERE is_active = TRUE;

CREATE TABLE IF NOT EXISTS  user_accounts(
    id VARCHAR(16) PRIMARY KEY,
    system_role VARCHAR(16) NOT NULL Check (system_role IN ('ADMIN', 'INSTITUTE')),  -- Added NOT NULL
    institution_role VARCHAR(16) , 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    email VARCHAR(255) NOT NULL, 
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS  institution_user(
    institution_id VARCHAR(16)  UNIQUE REFERENCES institutions(institution_id) ON DELETE CASCADE,
    user_id VARCHAR(16) REFERENCES user_accounts(id) ON DELETE CASCADE,
    public_key TEXT DEFAULT NULL,  -- Changed from VARCHAR(100) to TEXT
    institution_logo_base64 TEXT NOT NULL,
    PRIMARY KEY (institution_id, user_id) 
);

CREATE TABLE IF NOT EXISTS  institution_faculty(
    institution_faculty_id VARCHAR(16) PRIMARY KEY,
    institution_id VARCHAR(16) REFERENCES institutions(institution_id) ON DELETE CASCADE,
    faculty_name VARCHAR(200) NOT NULL,
    university_affiliation VARCHAR(100) NOT NULL,
    university_college_code VARCHAR(20) NOT NULL,
    faculty_authority_with_signature JSONB NOT NULL DEFAULT '[]'
);

alter table institution_faculty
add constraint unique_faculty_per_institution
unique (institution_id, faculty_name);

CREATE TABLE IF NOT EXISTS  pdf_file_categories(
    category_id VARCHAR(16) PRIMARY KEY,
    institution_id VARCHAR(16) REFERENCES institutions(institution_id) ON DELETE CASCADE,
    institution_faculty_id VARCHAR(16) REFERENCES institution_faculty(institution_faculty_id) ON DELETE CASCADE,
    category_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(institution_id, institution_faculty_id, category_name)
);

CREATE TABLE IF NOT EXISTS  pdf_files(
    file_id VARCHAR(16) PRIMARY KEY,
    category_id VARCHAR(16) REFERENCES pdf_file_categories(category_id) ON DELETE CASCADE,
    pdf_data BYTEA NOT NULL, -- Using BYTEA for binary data instead of BLOB
    file_name VARCHAR(255) NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(file_name, category_id)
);

-- Blocks table (header only)
CREATE TABLE IF NOT EXISTS  blocks(
    block_number INTEGER PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    previous_hash VARCHAR(255) NOT NULL,
    nonce VARCHAR(255) NOT NULL,
    current_hash VARCHAR(255) UNIQUE NOT NULL,
    merkle_root VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Enhanced Certificates table with flexible authority structure
CREATE TABLE IF NOT EXISTS  certificates(
     certificate_id VARCHAR(255) PRIMARY KEY,
    block_number INTEGER NOT NULL,
    position INTEGER NOT NULL CHECK (position BETWEEN 1 AND 4),
    
    -- Student Information (Required)
    student_id VARCHAR(255) NOT NULL,
    student_name VARCHAR(255) NOT NULL,
    
    -- Institution & Faculty Information
    institution_id VARCHAR(16) REFERENCES institutions(institution_id) ON DELETE CASCADE,
    institution_faculty_id VARCHAR(16) REFERENCES institution_faculty(institution_faculty_id) ON DELETE CASCADE,
     pdf_category_id VARCHAR(16) REFERENCES pdf_file_categories(category_id) ON DELETE SET NULL,


    
    
    -- Certificate type and basic info
    certificate_type VARCHAR(50) NOT NULL CHECK (certificate_type IN (
        'COURSE_COMPLETION', 
        'CHARACTER', 
        'LEAVING', 
        'TRANSFER',
        'PROVISIONAL'
    )),
    
    -- Academic information (Optional fields)
    degree VARCHAR(100),
    college VARCHAR(255),
    major VARCHAR(255),
    gpa VARCHAR(10),
    percentage DECIMAL(5,2),
    division VARCHAR(50),
    university_name VARCHAR(255),
    
    -- Date information
    issue_date TIMESTAMP NOT NULL,
    enrollment_date TIMESTAMP,
    completion_date TIMESTAMP,
    leaving_date TIMESTAMP,
    
    -- Reason fields (for leaving/character certificates)
    reason_for_leaving TEXT,
    character_remarks TEXT,
    general_remarks TEXT,
    
    -- Cryptographic verification
    data_hash VARCHAR(255) NOT NULL,
    issuer_public_key VARCHAR(255) NOT NULL,
    
   
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    
    FOREIGN KEY (block_number) REFERENCES blocks(block_number),
    UNIQUE(block_number, position)
);

ALTER TABLE certificates ADD COLUMN certificate_hash VARCHAR(255) NOT NULL DEFAULT 'temp_hash';

