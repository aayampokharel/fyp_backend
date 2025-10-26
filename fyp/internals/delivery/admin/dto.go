package admin

type AdminLoginRequest struct {
	AdminEmail string `json:"admin_email"`
	Password   string `json:"password"`
}

type AdminLoginResponse struct {
	SSEToken  string `json:"sse_token"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`

	////many more nested maps
}
