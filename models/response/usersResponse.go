package response

import mr "models/request"

/* UsersResponse is the response model for the GetUsers endpoint */
type UsersResponse struct {
	Users []*mr.User `json:"users"`
	Total int64      `json:"total"`
}
