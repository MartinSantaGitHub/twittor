package relational

import (
	"fmt"
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

	id, err := getUintId(requestModel.Id)

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

/* getTweetRequest obtains the Request Tweet model */
func getTweetRequest(tweetModel m.Tweet) mr.Tweet {
	requestModel := mr.Tweet{
		Id:      strconv.FormatUint(tweetModel.Id, 10),
		UserId:  strconv.FormatUint(tweetModel.UserId, 10),
		Message: tweetModel.Message,
		Date:    tweetModel.Date,
		Active:  tweetModel.Active,
	}

	return requestModel
}

/* getTweetModel obtains the DB Tweet model */
func getTweetModel(requestModel mr.Tweet) (m.Tweet, error) {
	var tweetModel m.Tweet

	uintId, err := getUintId(requestModel.Id)

	if err != nil {
		return tweetModel, err
	}

	uintUserId, err := getUintId(requestModel.UserId)

	if err != nil {
		return tweetModel, err
	}

	tweetModel = m.Tweet{
		Id:      uintId,
		Message: requestModel.Message,
		Date:    requestModel.Date,
		Active:  requestModel.Active,
		UserId:  uintUserId,
	}

	return tweetModel, nil
}

/* getRelationRequest obtains the Request Relation model */
func getRelationRequest(relationModel m.Relation) mr.Relation {
	requestModel := mr.Relation{
		Active: relationModel.Active,
	}

	return requestModel
}

// endregion

// region "Helpers"

func getUintId(id string) (uint64, error) {
	var uintId uint64
	var err error

	if len(id) < 1 {
		return uintId, nil
	}

	uintId, err = strconv.ParseUint(id, 10, 64)

	if err != nil {
		return uintId, fmt.Errorf("invalid id param")
	}

	return uintId, nil
}

// endregion
