package admin

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginResponse struct {
	SSEToken  string `json:"sse_token"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`

	////many more nested maps
}
