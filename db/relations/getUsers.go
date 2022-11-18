package relations

import (
	"context"
	"fmt"

	"db"
	"helpers"
	"models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* GetUsers Gets a list of users */
func GetUsers(id primitive.ObjectID, page int64, limit int64, search string, searchType string) ([]*models.User, int64, error) {
	var results []*models.User
	var total int64

	col := db.GetCollection("twittor", "users")

	query := bson.M{
		"name": bson.M{"$regex": search, "$options": "im"},
	}

	findOpts := options.Find()

	findOpts.SetSort(bson.D{{Key: "birthDate", Value: -1}})
	findOpts.SetSkip((page - 1) * limit)
	findOpts.SetLimit(limit)

	ctxFind, cancelFind := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))
	cursor, err := col.Find(ctxFind, query, findOpts)

	defer cancelFind()

	if err != nil {
		return results, total, err
	}

	ctxCursor := context.TODO()

	defer cursor.Close(ctxCursor)

	var found, include bool

	for cursor.Next(ctxCursor) {
		var result models.User

		err := cursor.Decode(&result)

		if err != nil {
			return results, total, err
		}

		relation := models.Relation{
			UserId:         id,
			UserRelationId: result.Id,
			Active:         true,
		}

		found, _, err = IsRelation(relation)

		if err != nil {
			return results, total, err
		}

		include = false

		if relation.UserRelationId == id {
			include = false
		} else if (searchType == "new" && !found) || (searchType == "follow" && found) {
			include = true
		}

		if include {
			result.Email = ""
			result.Password = ""
			result.Avatar = ""
			result.Banner = ""
			result.Biography = ""
			result.Location = ""
			result.WebSite = ""

			results = append(results, &result)
		}
	}

	err = cursor.Err()

	if err != nil {
		fmt.Println(err.Error())

		return results, total, err
	}

	// This total only reflects the total returned registries
	total = int64(len(results))

	return results, total, nil
}

/* GetFollowers Gets an user's followers list */
func GetFollowers(id primitive.ObjectID, page int64, limit int64, search string) ([]*models.User, int64, error) {

	// region Pipeline

	matchId := bson.M{"$match": bson.M{"userId": id, "active": true}}
	lookupUsers := bson.M{"$lookup": bson.M{
		"from":         "users",
		"localField":   "userRelationId",
		"foreignField": "_id",
		"as":           "result"}}
	projectResult := bson.M{"$project": bson.M{
		"user": bson.M{"$arrayElemAt": []interface{}{"$result", 0}},
		"_id":  0}}
	projectUser := bson.M{"$project": bson.M{
		"_id":       "$user._id",
		"name":      "$user.name",
		"lastName":  "$user.lastName",
		"birthDate": "$user.birthDate"}}
	matchName := bson.M{"$match": bson.M{"name": bson.M{"$regex": search, "$options": "im"}}}

	sort := bson.M{"$sort": bson.M{"birthDate": -1}}
	skip := bson.M{"$skip": (page - 1) * limit}
	agLimit := bson.M{"$limit": limit}
	count := bson.M{"$count": "total"}
	basePipeline := []bson.M{matchId, lookupUsers, projectResult, projectUser, matchName}
	countPipeline := append(basePipeline, count)
	aggPipeline := append(basePipeline, sort, skip, agLimit)

	// endregion

	results, total, err := db.GetResults[models.User]("relation", countPipeline, aggPipeline)

	return results, total, err
}

/* GetNotFollowers Gets an user's not followers list */
func GetNotFollowers(id primitive.ObjectID, page int64, limit int64, search string) ([]*models.User, int64, error) {

	// region Pipeline

	matchId := bson.M{"$match": bson.M{"_id": id}}
	lookupRelation := bson.M{"$lookup": bson.M{
		"from":         "relation",
		"localField":   "_id",
		"foreignField": "userId",
		"as":           "r",
		"pipeline": []interface{}{bson.M{
			"$match": bson.M{
				"$expr": bson.M{
					"$eq": []interface{}{"$active", true},
				},
			}},
		}},
	}
	lookupUsers := bson.M{"$lookup": bson.M{
		"from": "users",
		"as":   "u",
		"let":  bson.M{"userId": "$r.userRelationId", "id": "$r.userId"},
		"pipeline": []interface{}{bson.M{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": []interface{}{bson.M{
						"$not": bson.M{
							"$in": []string{"$_id", "$$userId"}}},
						bson.M{
							"$not": bson.M{
								"$in": []string{"$_id", "$$id"}}},
					},
				}},
		}},
	}}
	projectArray := bson.M{"$project": bson.M{"u": "$u", "_id": 0}}
	unwind := bson.M{"$unwind": bson.M{
		"path":                       "$u",
		"preserveNullAndEmptyArrays": false}}
	projectUser := bson.M{"$project": bson.M{
		"_id":       "$u._id",
		"name":      "$u.name",
		"lastName":  "$u.lastName",
		"birthDate": "$u.birthDate"}}
	matchName := bson.M{"$match": bson.M{"name": bson.M{"$regex": search, "$options": "im"}}}
	count := bson.M{"$count": "total"}
	sort := bson.M{"$sort": bson.M{"birthDate": -1}}
	skip := bson.M{"$skip": (page - 1) * limit}
	agLimit := bson.M{"$limit": limit}

	basePipeline := []bson.M{matchId, lookupRelation, lookupUsers, projectArray, unwind, projectUser, matchName}
	countPipeline := append(basePipeline, count)
	aggPipeline := append(basePipeline, sort, skip, agLimit)

	// endregion

	results, total, err := db.GetResults[models.User]("users", countPipeline, aggPipeline)

	return results, total, err
}
