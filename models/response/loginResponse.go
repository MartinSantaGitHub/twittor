package response

/* LoginResponse Login response model for the login endpoint */
type LoginResponse struct {
	Token string `json:"token,omitempty"`
}
