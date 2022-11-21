package nosqlv2

import (
	"fmt"

	m "models/nosqlv2"
	mr "models/request"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// region "Mappers"

/* GetUserModel obtains the DB User model */
func GetUserModel(requestModel mr.User) (m.User, error) {
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

/* GetUserRequest obtains the Request User model */
func GetUserRequest(userModel m.User) mr.User {
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

/* GetTweetModel obtains the DB Tweet model */
func GetTweetModel(requestModel mr.Tweet) (m.Tweet, error) {
	var tweetModel m.Tweet

	objId, err := getObjectId(requestModel.Id)

	if err != nil {
		return tweetModel, err
	}

	tweetModel = m.Tweet{
		Id:      objId,
		Message: requestModel.Message,
		Date:    requestModel.Date,
		Active:  requestModel.Active,
	}

	return tweetModel, nil
}

/* GetTweetRequest obtains the Request Tweet model */
func GetTweetRequest(tweetModel m.Tweet) mr.Tweet {
	requestModel := mr.Tweet{
		Id:      tweetModel.Id.Hex(),
		Message: tweetModel.Message,
		Date:    tweetModel.Date,
		Active:  tweetModel.Active,
	}

	return requestModel
}

/* GetRelationModel obtains the DB Relation model */
// func GetRelationModel(requestModel mr.Relation) (m.Relation, error) {
// 	var relationModel m.Relation

// 	objId, err := getObjectId(requestModel.Id)

// 	if err != nil {
// 		return relationModel, err
// 	}

// 	objUserId, err := getObjectId(requestModel.UserId)

// 	if err != nil {
// 		return relationModel, err
// 	}

// 	objUserRelationId, err := getObjectId(requestModel.UserRelationId)

// 	if err != nil {
// 		return relationModel, err
// 	}

// 	relationModel = m.Relation{
// 		Id:             objId,
// 		UserId:         objUserId,
// 		UserRelationId: objUserRelationId,
// 		Active:         requestModel.Active,
// 	}

// 	return relationModel, nil
// }

/* GetRelationRequest obtains the Request Relation model */
// func GetRelationRequest(relationModel m.Relation) mr.Relation {
// 	requestModel := mr.Relation{
// 		Id:             relationModel.Id.Hex(),
// 		UserId:         relationModel.UserId.Hex(),
// 		UserRelationId: relationModel.UserRelationId.Hex(),
// 		Active:         relationModel.Active,
// 	}

// 	return requestModel
// }

/* GetUserTweetRequest obtains the Request UserTweet model */
func GetUserTweetRequest(userTweetModel m.UserTweet) mr.UserTweet {
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
