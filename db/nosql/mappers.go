package nosql

import (
	"fmt"

	m "models/nosql"
	mr "models/request"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// region "Mappers"

/* getUserModel obtains the DB User model */
func getUserModel(requestModel mr.User) (m.User, error) {
	var userModel m.User

	objId, err := getObjectId(requestModel.Id)

	if err != nil {
		return userModel, err
	}

	userModel = m.User{
		Id:        objId,
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

/* getUserRequest obtains the Request User model */
func getUserRequest(userModel m.User) mr.User {
	requestModel := mr.User{
		Id:        userModel.Id.Hex(),
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

/* getTweetModel obtains the DB Tweet model */
func getTweetModel(requestModel mr.Tweet) (m.Tweet, error) {
	var tweetModel m.Tweet

	objId, err := getObjectId(requestModel.Id)

	if err != nil {
		return tweetModel, err
	}

	objUserId, err := getObjectId(requestModel.UserId)

	if err != nil {
		return tweetModel, err
	}

	tweetModel = m.Tweet{
		Id:      objId,
		UserId:  objUserId,
		Message: requestModel.Message,
		Date:    requestModel.Date,
		Active:  requestModel.Active,
	}

	return tweetModel, nil
}

/* getTweetRequest obtains the Request Tweet model */
func getTweetRequest(tweetModel m.Tweet) mr.Tweet {
	requestModel := mr.Tweet{
		Id:      tweetModel.Id.Hex(),
		UserId:  tweetModel.UserId.Hex(),
		Message: tweetModel.Message,
		Date:    tweetModel.Date,
		Active:  tweetModel.Active,
	}

	return requestModel
}

/* getRelationModel obtains the DB Relation model */
func getRelationModel(requestModel mr.Relation) (m.Relation, error) {
	var relationModel m.Relation

	objUserId, err := getObjectId(requestModel.UserId)

	if err != nil {
		return relationModel, err
	}

	objUserRelationId, err := getObjectId(requestModel.UserRelationId)

	if err != nil {
		return relationModel, err
	}

	relationModel = m.Relation{
		UserId:         objUserId,
		UserRelationId: objUserRelationId,
		Active:         requestModel.Active,
	}

	return relationModel, nil
}

/* getRelationRequest obtains the Request Relation model */
func getRelationRequest(relationModel m.Relation) mr.Relation {
	requestModel := mr.Relation{
		UserId:         relationModel.UserId.Hex(),
		UserRelationId: relationModel.UserRelationId.Hex(),
		Active:         relationModel.Active,
	}

	return requestModel
}

/* getUserTweetRequest obtains the Request UserTweet model */
func getUserTweetRequest(userTweetModel m.UserTweet) mr.UserTweet {
	requestModel := mr.UserTweet{
		Id:             userTweetModel.Id.Hex(),
		UserId:         userTweetModel.UserId.Hex(),
		UserRelationId: userTweetModel.UserRelationId.Hex(),
	}

	requestModel.Tweet.Id = userTweetModel.Tweet.Id.Hex()
	requestModel.Tweet.Message = userTweetModel.Tweet.Message
	requestModel.Tweet.Date = userTweetModel.Tweet.Date

	return requestModel
}

// endregion

// region "Helpers"

func getObjectId(id string) (primitive.ObjectID, error) {
	var objId primitive.ObjectID
	var err error

	if len(id) < 1 {
		return objId, nil
	}

	objId, err = primitive.ObjectIDFromHex(id)

	if err != nil {
		return objId, fmt.Errorf("invalid id param")
	}

	return objId, nil
}

// endregion
