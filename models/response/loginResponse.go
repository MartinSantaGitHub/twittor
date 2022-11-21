package response

/* LoginResponse is the response model for the login endpoint */
type LoginResponse struct {
	Token string `json:"token,omitempty"`
}
