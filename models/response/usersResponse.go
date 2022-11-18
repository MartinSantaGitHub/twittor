package response

import "models"

/* UsersResponse Response model for the GetUsers endpoint */
type UsersResponse struct {
	Users []*models.User `json:"users"`
	Total int64          `json:"total"`
}
