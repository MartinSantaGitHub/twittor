package nosqlv2

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"helpers"
	m "models/nosqlv2"
	mr "models/request"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type DbNoSqlV2 struct {
	Connection *mongo.Client
}

// region "Connection"

/* Connect connects to the database */
func (db *DbNoSqlV2) Connect() error {
	connTimeout := os.Getenv("CTX_TIMEOUT")
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_CONN"))
	ctx, cancel := helpers.GetTimeoutCtx(connTimeout)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return err
	}

	db.Connection = client

	return nil
}

/* IsConnection makes a ping to the Database */
func (db *DbNoSqlV2) IsConnection() bool {
	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	err := db.Connection.Ping(ctx, nil)

	if err != nil {
		fmt.Println(err.Error())

		return false
	}

	return true
}

// endregion

// region "Users"

/* GetProfile gets a profile in the DB */
func (db *DbNoSqlV2) GetProfile(id string) (mr.User, bool, error) {
	var profileRequest mr.User
	var profileModel m.User

	objId, err := getObjectId(id)

	if err != nil {
		return profileRequest, false, err
	}

	col := getCollection(db, "twitton", "users")

	condition := bson.M{
		"_id": objId,
	}

	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	err = col.FindOne(ctx, condition).Decode(&profileModel)

	if err != nil && err == mongo.ErrNoDocuments {
		log.Println("Registry not found: " + err.Error())

		return profileRequest, false, nil
	} else if err != nil {
		return profileRequest, false, err
	}

	profileRequest = getUserRequest(profileModel)

	profileRequest.Password = ""

	return profileRequest, true, nil
}

/* InsertUser inserts an user into de DB */
func (db *DbNoSqlV2) InsertUser(user mr.User) (string, error) {
	col := getCollection(db, "twitton", "users")

	user.Password, _ = encryptPassword(user.Password)
	userModel, err := getUserModel(user)

	if err != nil {
		return "", err
	}

	callback := func(sessCtx mongo.SessionContext) (any, error) {
		result, err := col.InsertOne(sessCtx, userModel)

		return result, err
	}

	res, err := db.executeTransaction(callback)

	if err != nil {
		return "", err
	}

	result := res.(*mongo.InsertOneResult)
	objID, _ := result.InsertedID.(primitive.ObjectID)

	return objID.String(), nil
}

/* IsUser checks that the user already exists in the DB */
func (db *DbNoSqlV2) IsUser(email string) (bool, mr.User, error) {
	var userModel m.User
	var requestModel mr.User

	col := getCollection(db, "twitton", "users")
	condition := bson.M{"email": email}
	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	err := col.FindOne(ctx, condition).Decode(&userModel)

	if err != nil && err == mongo.ErrNoDocuments {
		return false, requestModel, nil
	}

	requestModel = getUserRequest(userModel)

	return true, requestModel, err
}

/* ModifyRegistry modifies a registry in the DB */
func (db *DbNoSqlV2) ModifyRegistry(id string, user mr.User) error {
	objId, err := getObjectId(id)

	if err != nil {
		return err
	}

	col := getCollection(db, "twitton", "users")
	registry := make(map[string]any)

	if len(user.Name) > 0 {
		registry["name"] = user.Name
	}

	if len(user.LastName) > 0 {
		registry["lastName"] = user.LastName
	}

	if len(user.Avatar) > 0 {
		registry["avatar"] = user.Avatar
	}

	if len(user.Banner) > 0 {
		registry["banner"] = user.Banner
	}

	if len(user.Biography) > 0 {
		registry["biography"] = user.Biography
	}

	if len(user.Location) > 0 {
		registry["location"] = user.Location
	}

	if len(user.WebSite) > 0 {
		registry["webSite"] = user.WebSite
	}

	if !user.BirthDate.IsZero() {
		registry["birthDate"] = user.BirthDate
	}

	updateString := bson.M{
		"$set": registry,
	}

	filter := bson.M{"_id": objId}
	//filter := bson.M{"_id": bson.M{"$eq": objId}}

	callback := func(sessCtx mongo.SessionContext) (any, error) {
		result, err := col.UpdateOne(sessCtx, filter, updateString)
		// result, err := col.UpdateByID(ctx, objId, updateString)

		return result, err
	}

	_, err = db.executeTransaction(callback)

	return err
}

/* TryLogin makes the login to the DB */
func (db *DbNoSqlV2) TryLogin(email string, password string) (mr.User, bool) {
	var requestModel mr.User

	isFound, requestModel, err := db.IsUser(email)

	if err != nil || !isFound {
		return requestModel, false
	}

	passwordBytes := []byte(password)
	passwordDB := []byte(requestModel.Password)
	err = bcrypt.CompareHashAndPassword(passwordDB, passwordBytes)

	if err != nil {
		return requestModel, false
	}

	return requestModel, true
}

// endregion

// region "Tweets"

/* DeleteTweet deletes a tweet in the DB */
func (db *DbNoSqlV2) DeleteTweet(id string, userId string) error {
	err := db.deleteTweetLogical(id, userId)

	return err
}

/* Get gets an user's tweets from the DB */
func (db *DbNoSqlV2) GetTweets(id string, page int64, limit int64) ([]*mr.Tweet, int64, error) {
	var results []*mr.Tweet

	objId, err := getObjectId(id)

	if err != nil {
		return results, 0, err
	}

	// region Pipeline

	matchId := bson.M{"$match": bson.M{"_id": objId}}
	projectTweets := bson.M{"$project": bson.M{
		"t":   "$tweets",
		"_id": 0}}
	unwindTweets := bson.M{"$unwind": bson.M{
		"path":                       "$t",
		"preserveNullAndEmptyArrays": false}}
	filterTweets := bson.M{"$match": bson.M{"t.active": true}}

	count := bson.M{"$count": "total"}

	sort := bson.M{"$sort": bson.M{"t.date": -1}}
	skip := bson.M{"$skip": (page - 1) * limit}
	agLimit := bson.M{"$limit": limit}
	projectResult := bson.M{"$project": bson.M{
		"_id":     "$t._id",
		"message": "$t.message",
		"date":    "$t.date",
		"active":  "$t.active"}}

	basePipeline := []bson.M{matchId, projectTweets, unwindTweets, filterTweets}
	countPipeline := append(basePipeline, count)
	aggPipeline := append(basePipeline, sort, skip, agLimit, projectResult)

	// endregion

	dbResults, total, err := getResults[m.Tweet](db, "users", countPipeline, aggPipeline)

	if err == nil {
		for _, tweetModel := range dbResults {
			tweetRequest := getTweetRequest(*tweetModel)
			results = append(results, &tweetRequest)

			tweetRequest.UserId = ""
		}
	}

	return results, total, err
}

/* InsertTweet inserts a tweet in the DB */
func (db *DbNoSqlV2) InsertTweet(tweet mr.Tweet) (string, error) {
	objId, _ := getObjectId(tweet.UserId)
	update := bson.M{
		"$push": bson.M{
			"tweets": bson.D{
				primitive.E{Key: "_id", Value: primitive.NewObjectID()},
				primitive.E{Key: "message", Value: tweet.Message},
				primitive.E{Key: "date", Value: tweet.Date},
				primitive.E{Key: "active", Value: tweet.Active},
			},
		},
	}

	col := getCollection(db, "twitton", "users")
	callback := func(sessCtx mongo.SessionContext) (any, error) {
		result, err := col.UpdateByID(sessCtx, objId, update)

		return result, err
	}

	res, err := db.executeTransaction(callback)

	if err != nil {
		return "", err
	}

	result := res.(*mongo.UpdateResult)
	objID, _ := result.UpsertedID.(primitive.ObjectID)

	// The same goes with objId.String()
	return objID.Hex(), nil
}

// endregion

// region "Relations"

/* IsRelation obtains a relation from the DB if exist */
func (db *DbNoSqlV2) IsRelation(relation mr.Relation) (bool, mr.Relation, error) {
	objId, _ := getObjectId(relation.UserId)
	objUserFollowingId, err := getObjectId(relation.UserRelationId)

	if err != nil {
		return false, relation, err
	}

	col := getCollection(db, "twitton", "users")
	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	filter := bson.M{
		"_id": objId,
		"following": bson.M{
			"$in": [1]primitive.ObjectID{
				objUserFollowingId,
			}},
	}

	err = col.FindOne(ctx, filter).Err()

	if err != nil && err == mongo.ErrNoDocuments {
		return false, relation, nil
	} else if err != nil && err != mongo.ErrNoDocuments {
		return false, relation, err
	}

	relation.Active = true

	return true, relation, nil
}

/* InsertRelation creates a relation into the DB */
func (db *DbNoSqlV2) InsertRelation(relation mr.Relation) error {
	_, isFound, err := db.GetProfile(relation.UserRelationId)

	if err != nil {
		return err
	}

	if !isFound {
		return errors.New("no registry found in the DB")
	}

	objId, _ := getObjectId(relation.UserId)
	objUserFollowingId, _ := getObjectId(relation.UserRelationId)

	update := bson.M{
		"$addToSet": bson.M{
			"following": objUserFollowingId,
		},
	}

	col := getCollection(db, "twitton", "users")
	callback := func(sessCtx mongo.SessionContext) (any, error) {
		result, err := col.UpdateByID(sessCtx, objId, update)

		return result, err
	}

	_, err = db.executeTransaction(callback)

	return err
}

/* DeleteRelation deletes a relation in the DB */
func (db *DbNoSqlV2) DeleteRelation(relation mr.Relation) error {
	err := db.deleteRelationFisical(relation)

	return err
}

/* GetUsers gets a list of users */
func (db *DbNoSqlV2) GetUsers(id string, page int64, limit int64, search string, searchType string) ([]*mr.User, int64, error) {
	var results []*mr.User
	var total int64

	col := getCollection(db, "twitton", "users")

	query := bson.M{
		"name": bson.M{"$regex": search, "$options": "im"},
	}

	findOpts := options.Find()

	findOpts.SetSort(bson.D{{Key: "birthDate", Value: -1}})
	findOpts.SetSkip((page - 1) * limit)
	findOpts.SetLimit(limit)

	ctxFind, cancelFind := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelFind()

	cursor, err := col.Find(ctxFind, query, findOpts)

	if err != nil {
		return results, total, err
	}

	ctxCursor := context.TODO()

	defer cursor.Close(ctxCursor)

	var include, isRelation bool

	for cursor.Next(ctxCursor) {
		var result m.User

		isRelation = false
		include = false

		err := cursor.Decode(&result)

		if err != nil {
			return results, total, err
		}

		userRequest := getUserRequest(result)

		relationRequest := mr.Relation{
			UserId:         id,
			UserRelationId: userRequest.Id,
			Active:         true,
		}

		isFound, relationDb, err := db.IsRelation(relationRequest)

		if err != nil {
			return results, total, err
		}

		if isFound {
			isRelation = relationDb.Active
		}

		if relationRequest.UserRelationId == id {
			include = false
		} else if (searchType == "new" && !isRelation) || (searchType == "follow" && isRelation) {
			include = true
		}

		if include {
			userRequest.Email = ""
			userRequest.Password = ""
			userRequest.Avatar = ""
			userRequest.Banner = ""
			userRequest.Biography = ""
			userRequest.Location = ""
			userRequest.WebSite = ""

			results = append(results, &userRequest)
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

/* GetFollowing gets an user's following list */
func (db *DbNoSqlV2) GetFollowing(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	var results []*mr.User

	objId, _ := getObjectId(id)

	// region Pipeline

	matchId := bson.M{"$match": bson.M{"_id": objId}}
	lookupUsers := bson.M{"$lookup": bson.M{
		"from":         "users",
		"localField":   "following",
		"foreignField": "_id",
		"as":           "r"}}
	projectResult := bson.M{"$project": bson.M{
		"u":   "$r",
		"_id": 0}}
	unwind := bson.M{"$unwind": bson.M{
		"path":                       "$u",
		"preserveNullAndEmptyArrays": false}}
	matchName := bson.M{"$match": bson.M{"u.name": bson.M{"$regex": search, "$options": "im"}}}

	count := bson.M{"$count": "total"}

	sort := bson.M{"$sort": bson.M{"u.birthDate": -1}}
	skip := bson.M{"$skip": (page - 1) * limit}
	agLimit := bson.M{"$limit": limit}
	projectUser := bson.M{"$project": bson.M{
		"_id":       "$u._id",
		"name":      "$u.name",
		"lastName":  "$u.lastName",
		"birthDate": "$u.birthDate"}}

	basePipeline := []bson.M{matchId, lookupUsers, projectResult, unwind, matchName}
	countPipeline := append(basePipeline, count)
	aggPipeline := append(basePipeline, sort, skip, agLimit, projectUser)

	// endregion

	dbResults, total, err := getResults[m.User](db, "users", countPipeline, aggPipeline)

	if err == nil {
		for _, userModel := range dbResults {
			userRequest := getUserRequest(*userModel)
			results = append(results, &userRequest)
		}
	}

	return results, total, err
}

/* GetNotFollowing gets an user's not following list */
func (db *DbNoSqlV2) GetNotFollowing(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	var results []*mr.User

	objId, _ := getObjectId(id)

	// region Pipeline

	matchId := bson.M{"$match": bson.M{"_id": objId}}
	lookupUsers := bson.M{"$lookup": bson.M{
		"from": "users",
		"as":   "r",
		"let":  bson.M{"userId": "$following"},
		"pipeline": [1]bson.M{{
			"$match": bson.M{
				"$expr": bson.M{
					"$not": bson.M{"$in": [2]string{
						"$_id",
						"$$userId",
					}},
				}},
		}},
	}}
	projectArray := bson.M{"$project": bson.M{"u": "$r", "_id": 0}}
	unwind := bson.M{"$unwind": bson.M{
		"path":                       "$u",
		"preserveNullAndEmptyArrays": false,
	}}
	matchName := bson.M{"$match": bson.M{
		"u._id":  bson.M{"$ne": objId},
		"u.name": bson.M{"$regex": search, "$options": "im"},
	}}

	count := bson.M{"$count": "total"}

	sort := bson.M{"$sort": bson.M{"u.birthDate": -1}}
	skip := bson.M{"$skip": (page - 1) * limit}
	agLimit := bson.M{"$limit": limit}
	projectUser := bson.M{"$project": bson.M{
		"_id":       "$u._id",
		"name":      "$u.name",
		"lastName":  "$u.lastName",
		"birthDate": "$u.birthDate"}}

	basePipeline := []bson.M{matchId, lookupUsers, projectArray, unwind, matchName}
	countPipeline := append(basePipeline, count)
	aggPipeline := append(basePipeline, sort, skip, agLimit, projectUser)

	// endregion

	dbResults, total, err := getResults[m.User](db, "users", countPipeline, aggPipeline)

	if err == nil {
		for _, userModel := range dbResults {
			userRequest := getUserRequest(*userModel)
			results = append(results, &userRequest)
		}
	}

	return results, total, err
}

/* GetFollowingTweets returns the following's tweets */
func (db *DbNoSqlV2) GetFollowingTweets(id string, page int64, limit int64, isOnlyTweets bool) (any, int64, error) {
	var results any
	var total int64
	var err error

	objId, err := getObjectId(id)

	if err != nil {
		return results, 0, err
	}

	conditions := make([]bson.M, 0)
	conditionsCount := make([]bson.M, 0)
	conditionsAgg := make([]bson.M, 0)

	skip := (page - 1) * limit

	conditions = append(conditions, bson.M{"$match": bson.M{"_id": objId}})
	conditions = append(conditions, bson.M{
		"$lookup": bson.M{
			"from":         "users",
			"localField":   "following",
			"foreignField": "_id",
			"as":           "r",
		}})
	conditions = append(conditions, bson.M{
		"$unwind": bson.M{
			"path":                       "$r",
			"preserveNullAndEmptyArrays": false,
		},
	})
	conditions = append(conditions, bson.M{
		"$project": bson.M{
			"_id":             0,
			"userId":          "$_id",
			"userFollowingId": "$r._id",
			"tweet": bson.M{
				"$filter": bson.M{
					"input": "$r.tweets",
					"as":    "tweet",
					"cond":  bson.M{"$eq": [2]any{"$$tweet.active", true}},
				},
			},
		},
	})
	conditions = append(conditions, bson.M{
		"$unwind": bson.M{
			"path":                       "$tweet",
			"preserveNullAndEmptyArrays": false,
		},
	})

	conditionsCount = append(conditionsCount, conditions...)
	conditionsCount = append(conditionsCount, bson.M{"$count": "total"})

	conditionsAgg = append(conditionsAgg, conditions...)
	conditionsAgg = append(conditionsAgg, bson.M{"$sort": bson.M{"tweet.date": -1}})
	conditionsAgg = append(conditionsAgg, bson.M{"$skip": skip})
	conditionsAgg = append(conditionsAgg, bson.M{"$limit": limit})

	if isOnlyTweets {
		conditionsAgg = append(conditionsAgg, bson.M{
			"$project": bson.M{
				"_id":     "$tweet._id",
				"userId":  "$userFollowingId",
				"message": "$tweet.message",
				"date":    "$tweet.date",
			}})

		var dbResults []*m.Tweet
		var reqResults []*mr.Tweet

		dbResults, total, err = getResults[m.Tweet](db, "users", conditionsCount, conditionsAgg)

		if err == nil {
			for _, tweetModel := range dbResults {
				tweetRequest := getTweetRequest(*tweetModel)
				reqResults = append(reqResults, &tweetRequest)
			}

			results = reqResults
		}
	} else {
		var dbResults []*m.UserTweet
		var reqResults []*mr.UserTweet

		dbResults, total, err = getResults[m.UserTweet](db, "users", conditionsCount, conditionsAgg)

		if err == nil {
			for _, userTweetModel := range dbResults {
				userTweetRequest := getUserTweetRequest(*userTweetModel)
				reqResults = append(reqResults, &userTweetRequest)
			}

			results = reqResults
		}
	}

	return results, total, err
}

// endregion

// region "Helpers"

func (db *DbNoSqlV2) deleteTweetLogical(id string, userId string) error {
	objId, err := getObjectId(id)

	if err != nil {
		return err
	}

	objUserId, _ := getObjectId(userId)
	filter := bson.M{
		"_id":        objUserId,
		"tweets._id": objId,
	}
	update := bson.M{
		"$set": bson.M{
			"tweets.$.active": false,
		},
	}

	col := getCollection(db, "twitton", "users")
	callback := func(sessCtx mongo.SessionContext) (any, error) {
		result, err := col.UpdateOne(sessCtx, filter, update)

		return result, err
	}

	_, err = db.executeTransaction(callback)

	return err
}

func (db *DbNoSqlV2) deleteTweetFisical(id string, userId string) error {
	objId, err := getObjectId(id)

	if err != nil {
		return err
	}

	objUserId, _ := getObjectId(userId)
	update := bson.M{
		"$pull": bson.M{
			"tweets": bson.D{primitive.E{Key: "_id", Value: objId}},
		},
	}

	col := getCollection(db, "twitton", "users")
	callback := func(sessCtx mongo.SessionContext) (any, error) {
		result, err := col.UpdateByID(sessCtx, objUserId, update)

		return result, err
	}

	_, err = db.executeTransaction(callback)

	return err
}

func (db *DbNoSqlV2) deleteRelationFisical(relation mr.Relation) error {
	objId, _ := getObjectId(relation.UserId)
	objUserFollowingId, err := getObjectId(relation.UserRelationId)

	if err != nil {
		return err
	}

	update := bson.M{
		"$pull": bson.M{
			"following": objUserFollowingId,
		},
	}

	col := getCollection(db, "twitton", "users")
	callback := func(sessCtx mongo.SessionContext) (any, error) {
		result, err := col.UpdateByID(sessCtx, objId, update)

		return result, err
	}

	_, err = db.executeTransaction(callback)

	return err
}

func encryptPassword(password string) (string, error) {
	// Minimum - cost: 6
	// Common user - cost: 6
	// Admin user - cost: 8

	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	return string(bytes), err
}

func (db *DbNoSqlV2) executeTransaction(callback func(sessCtx mongo.SessionContext) (any, error)) (any, error) {
	session, err := db.Connection.StartSession()

	if err != nil {
		return nil, err
	}

	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer session.EndSession(ctx)
	defer cancel()

	result, err := session.WithTransaction(ctx, callback)

	return result, err
}

func getCollection(db *DbNoSqlV2, dbName string, colName string) *mongo.Collection {
	database := db.Connection.Database(dbName)
	collection := database.Collection(colName)

	return collection
}

func getResults[T any](db *DbNoSqlV2, colName string, countPipeline []primitive.M, aggPipeline []primitive.M) ([]*T, int64, error) {
	var results []*T
	var totalResult m.TotalResult

	col := getCollection(db, "twitton", colName)

	// region "Count"

	ctxCount, cancelCount := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelCount()

	curCount, err := col.Aggregate(ctxCount, countPipeline)

	if err != nil {
		return results, totalResult.Total, err
	}

	ctxCurCount := context.TODO()

	defer curCount.Close(ctxCurCount)

	if curCount.Next(ctxCurCount) {
		err = curCount.Decode(&totalResult)
	}

	if err != nil {
		return results, totalResult.Total, err
	}

	// endregion

	// region "Aggregate"

	ctxAggregate, cancelAggregate := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelAggregate()

	curAgg, err := col.Aggregate(ctxAggregate, aggPipeline)

	if err != nil {
		return results, totalResult.Total, err
	}

	ctxCurAgg := context.TODO()

	defer curAgg.Close(ctxCurAgg)

	err = curAgg.All(ctxCurAgg, &results)

	if err != nil {
		return results, totalResult.Total, err
	}

	// endregion

	return results, totalResult.Total, nil
}

// endregion
