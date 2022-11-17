package response

import "models"

type RelationUsersResponse struct {
	Users []*models.User `json:"users"`
	Total int64          `json:"total"`
}
