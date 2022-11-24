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
		Tweets:    []m.Tweet{},
		Following: []primitive.ObjectID{},
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

/* GetTweetRequest obtains the Request Tweet model */
func GetTweetRequest(tweetModel m.Tweet) mr.Tweet {
	requestModel := mr.Tweet{
		Id:      tweetModel.Id.Hex(),
		UserId:  tweetModel.UserId.Hex(),
		Message: tweetModel.Message,
		Date:    tweetModel.Date,
		Active:  tweetModel.Active,
	}

	return requestModel
}

/* GetUserTweetRequest obtains the Request UserTweet model */
func GetUserTweetRequest(userTweetModel m.UserTweet) mr.UserTweet {
	requestModel := mr.UserTweet{
		UserId:         userTweetModel.UserId.Hex(),
		UserRelationId: userTweetModel.UserFollowingId.Hex(),
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
