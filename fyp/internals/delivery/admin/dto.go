package admin

import "project/internals/domain/entity"

type AdminLoginRequest struct {
	AdminEmail string `json:"admin_email"`
	Password   string `json:"password"`
}

type AdminLoginResponse struct {
	UserID                     string                            `json:"user_id"`
	CreatedAt                  string                            `json:"created_at"`
	InstitutionList            []entity.Institution              `json:"institution_list"`
	AdminDashboardCountDetails entity.AdminDashboardCountsEntity `json:"admin_dashboard_count_details"`

	////many more nested maps
}

type GetAllPendingInstitutionsRequest []string

const AdminID = "admin_id"

var GetAllPendingInstitutionsQuery = GetAllPendingInstitutionsRequest{AdminID}

type GetAllPendingInstitutionsResponse struct {
	PendingInstitutionList []entity.Institution `json:"pending_institution_list"`
}
