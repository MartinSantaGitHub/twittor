package relational

import (
	"strconv"

	m "models/relational"
	mr "models/request"
)

// region "Mappers"

/* getUserRequest obtains the Request User model */
func getUserRequest(userModel m.User) mr.User {
	requestModel := mr.User{
		Id:        strconv.FormatUint(userModel.Id, 10),
		Name:      userModel.Name,
		LastName:  userModel.LastName,
		Email:     userModel.Email,
		BirthDate: userModel.BirthDate,
		Avatar:    userModel.Avatar,
		Banner:    userModel.Banner,
		Biography: userModel.Biography,
		Location:  userModel.Location,
		WebSite:   userModel.WebSite,
		Password:  userModel.Password,
	}

	return requestModel
}

/* getUserModel obtains the DB User model */
func getUserModel(requestModel mr.User) (m.User, error) {
	var userModel m.User

	id, err := strconv.ParseUint(requestModel.Id, 10, 64)

	if err != nil {
		return userModel, err
	}

	userModel = m.User{
		Id:        id,
		Name:      requestModel.Name,
		LastName:  requestModel.LastName,
		Email:     requestModel.Email,
		BirthDate: requestModel.BirthDate,
		Avatar:    requestModel.Avatar,
		Banner:    requestModel.Banner,
		Biography: requestModel.Biography,
		Location:  requestModel.Location,
		WebSite:   requestModel.WebSite,
		Password:  requestModel.Password,
	}

	return userModel, nil
}

// endregion
