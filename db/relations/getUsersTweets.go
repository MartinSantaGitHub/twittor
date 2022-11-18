package relations

import (
	"db"
	"models"
	mr "models/result"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* GetUsersTweets Returns the followers' tweets */
func GetUsersTweets(id primitive.ObjectID, page int64, limit int64, isOnlyTweets bool) (interface{}, int64, error) {
	var results interface{}
	var total int64
	var err error

	conditions := make([]bson.M, 0)
	conditionsCount := make([]bson.M, 0)
	conditionsAgg := make([]bson.M, 0)

	skip := (page - 1) * limit

	// region "Pipeline"

	conditions = append(conditions, bson.M{"$match": bson.M{"userId": id, "active": true}})
	conditions = append(conditions, bson.M{
		"$lookup": bson.M{
			"from":         "tweet",
			"localField":   "userRelationId",
			"foreignField": "userId",
			"as":           "tweet",
			"pipeline": []bson.M{{
				"$match": bson.M{
					"active": true}},
			},
		}})
	conditions = append(conditions, bson.M{
		"$unwind": bson.M{
			"path":                       "$tweet",
			"preserveNullAndEmptyArrays": false,
		},
	})

	if isOnlyTweets {
		conditions = append(conditions, bson.M{
			"$project": bson.M{
				"_id":     "$tweet._id",
				"userId":  "$tweet.userId",
				"message": "$tweet.message",
				"date":    "$tweet.date",
			}})
		conditions = append(conditions, bson.M{"$sort": bson.M{"date": -1}})
	} else {
		conditions = append(conditions, bson.M{"$sort": bson.M{"tweet.date": -1}})
	}

	conditionsCount = append(conditionsCount, conditions...)
	conditionsCount = append(conditionsCount, bson.M{"$count": "total"})

	conditionsAgg = append(conditionsAgg, conditions...)
	conditionsAgg = append(conditionsAgg, bson.M{"$skip": skip})
	conditionsAgg = append(conditionsAgg, bson.M{"$limit": limit})

	// endregion

	if isOnlyTweets {
		results, total, err = db.GetResults[models.Tweet]("relation", conditionsCount, conditionsAgg)
	} else {
		results, total, err = db.GetResults[mr.UserTweet]("relation", conditionsCount, conditionsAgg)
	}

	return results, total, err
}
