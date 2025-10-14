package models

import (
	"time"
)

// All GORM models for database persistence

type InstitutionModel struct {
	InstitutionID   string `gorm:"primaryKey;column:institution_id;size:16"`
	InstitutionName string `gorm:"column:institution_name;size:300;not null"`
	ToleAddress     string `gorm:"column:tole_address;size:250;not null"`
	DistrictAddress string `gorm:"column:district_address;size:250;not null"`
	IsActive        bool   `gorm:"column:is_active;default:true"`
}

func (InstitutionModel) TableName() string { return "institutions" }

type UserAccountModel struct {
	ID        string     `gorm:"primaryKey;column:id;size:16"`
	Role      string     `gorm:"column:role;size:16;not null"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at;default:NULL"`
	Email     string     `gorm:"column:email;size:255;uniqueIndex;not null"`
	Password  string     `gorm:"column:password;size:255;not null"`
}

func (UserAccountModel) TableName() string { return "user_accounts" }

type InstitutionUserModel struct {
	InstitutionID            string `gorm:"primaryKey;column:institution_id;size:16"`
	UserID                   string `gorm:"primaryKey;column:user_id;size:16"`
	PublicKey                string `gorm:"column:public_key;type:text"`
	PrincipalName            string `gorm:"column:principal_name;size:300;not null"`
	PrincipalSignatureBase64 string `gorm:"column:principal_signature_base64;type:text;not null"`
	InstitutionLogoBase64    string `gorm:"column:institution_logo_base64;type:text;not null"`

	User        UserAccountModel `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Institution InstitutionModel `gorm:"foreignKey:InstitutionID;references:InstitutionID;constraint:OnDelete:CASCADE"`
}

func (InstitutionUserModel) TableName() string { return "institution_user" }

type InstitutionFacultyModel struct {
	InstitutionFacultyID      string `gorm:"primaryKey;column:institution_faculty_id;size:16"`
	InstitutionID             string `gorm:"column:institution_id;size:16;not null"`
	Faculty                   string `gorm:"column:faculty;size:200;not null"`
	FacultyHODName            string `gorm:"column:faculty_hod_name;size:300;not null"`
	FacultyHODSignatureBase64 string `gorm:"column:faculty_hod_signature_base64;type:text;not null"`

	//foreign key
	Institution InstitutionModel `gorm:"foreignKey:InstitutionID;references:InstitutionID;constraint:OnDelete:CASCADE"`
}

func (InstitutionFacultyModel) TableName() string { return "institution_faculty" }

type BlockModel struct {
	BlockNumber  int       `gorm:"primaryKey;column:block_number"`
	Timestamp    time.Time `gorm:"column:timestamp;not null"`
	PreviousHash string    `gorm:"column:previous_hash;size:255;not null"`
	Nonce        string    `gorm:"column:nonce;size:255;not null"`
	CurrentHash  string    `gorm:"column:current_hash;size:255;uniqueIndex;not null"`
	MerkleRoot   string    `gorm:"column:merkle_root;size:255;not null"`
	Status       string    `gorm:"column:status;size:50;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`

	Certificates []CertificateModel `gorm:"foreignKey:BlockNumber"`
}

func (BlockModel) TableName() string { return "blocks" }

type CertificateModel struct {
	ID                 uint      `gorm:"primaryKey;column:id;autoIncrement"`
	CertificateID      string    `gorm:"column:certificate_id;size:255;not null"`
	BlockNumber        int       `gorm:"column:block_number;not null;uniqueIndex:idx_block_position"`
	Position           int       `gorm:"column:position;not null;uniqueIndex:idx_block_position"` // enforce 1-4 in app
	StudentID          string    `gorm:"column:student_id;size:255;not null;index"`
	StudentName        string    `gorm:"column:student_name;size:255;not null"`
	UniversityName     string    `gorm:"column:university_name;size:255;not null"`
	Degree             string    `gorm:"column:degree;size:100;not null"`
	College            string    `gorm:"column:college;size:255;not null"`
	Major              string    `gorm:"column:major;size:255;not null"`
	GPA                string    `gorm:"column:gpa;size:10"`
	Percentage         float64   `gorm:"column:percentage;type:decimal(5,2)"`
	Division           string    `gorm:"column:division;size:50;not null"`
	IssueDate          time.Time `gorm:"column:issue_date;not null"`
	EnrollmentDate     time.Time `gorm:"column:enrollment_date;not null"`
	CompletionDate     time.Time `gorm:"column:completion_date;not null"`
	PrincipalSignature string    `gorm:"column:principal_signature;size:255;not null"`
	DataHash           string    `gorm:"column:data_hash;size:255;not null;index"`
	IssuerPublicKey    string    `gorm:"column:issuer_public_key;size:255;not null"`
	CertificateType    string    `gorm:"column:certificate_type;size:50;not null"`

	Block BlockModel `gorm:"foreignKey:BlockNumber;references:BlockNumber;constraint:OnDelete:CASCADE"`
}

func (CertificateModel) TableName() string { return "certificates" }
