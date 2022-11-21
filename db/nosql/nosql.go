package nosql

import (
	"context"
	"fmt"
	"log"

	"helpers"
	m "models/nosql"
	mr "models/request"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type DbNoSql struct {
	Connection *mongo.Client
}

// region "Connection"

/* Connect connects to the database */
func (db *DbNoSql) Connect() error {
	connTimeout := helpers.GetEnvVariable("CTX_TIMEOUT")
	clientOptions := options.Client().ApplyURI(helpers.GetEnvVariable("MONGO_CONN"))
	ctx, cancel := helpers.GetTimeoutCtx(connTimeout)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return err
	}

	db.Connection = client

	return nil
}

/* CheckConnection makes a ping to the Database */
func (db *DbNoSql) IsConnection() bool {
	err := db.Connection.Ping(context.TODO(), nil)

	if err != nil {
		fmt.Println(err.Error())

		return false
	}

	return true
}

// endregion

// region "Users"

/* GetProfile gets a profile in the DB */
func (db *DbNoSql) GetProfile(id string) (mr.User, error) {
	var profileRequest mr.User
	var profileModel m.User

	objId, err := getObjectId(id)

	if err != nil {
		return profileRequest, err
	}

	col := getCollection(db, "twittor", "users")

	condition := bson.M{
		"_id": objId,
	}

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancel()

	err = col.FindOne(ctx, condition).Decode(&profileModel)

	if err != nil {
		log.Println("Registry not found: " + err.Error())

		return profileRequest, err
	}

	profileRequest = getUserRequest(profileModel)

	profileRequest.Password = ""

	return profileRequest, nil
}

/* InsertUser inserts an user into de DB */
func (db *DbNoSql) InsertUser(user mr.User) (string, error) {
	col := getCollection(db, "twittor", "users")

	user.Password, _ = encryptPassword(user.Password)
	userModel, err := getUserModel(user)

	if err != nil {
		return "", err
	}

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancel()

	result, err := col.InsertOne(ctx, userModel)

	if err != nil {
		return "", err
	}

	objID, _ := result.InsertedID.(primitive.ObjectID)

	return objID.String(), nil
}

/* IsUser checks that the user already exists in the DB */
func (db *DbNoSql) IsUser(email string) (bool, mr.User, error) {
	var userModel m.User
	var requestModel mr.User

	col := getCollection(db, "twittor", "users")
	condition := bson.M{"email": email}
	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancel()

	err := col.FindOne(ctx, condition).Decode(&userModel)

	if err != nil && err == mongo.ErrNoDocuments {
		return false, requestModel, nil
	}

	requestModel = getUserRequest(userModel)

	return true, requestModel, err
}

/* ModifyRegistry modifies a registry in the DB */
func (db *DbNoSql) ModifyRegistry(id string, user mr.User) error {
	objId, err := getObjectId(id)

	if err != nil {
		return err
	}

	col := getCollection(db, "twittor", "users")
	registry := make(map[string]interface{})

	if len(user.Name) > 0 {
		registry["name"] = user.Name
	}

	if len(user.LastName) > 0 {
		registry["lastName"] = user.LastName
	}

	if len(user.Avatar) > 0 {
		registry["Avatar"] = user.Avatar
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

	registry["birthDate"] = user.BirthDate

	updateString := bson.M{
		"$set": registry,
	}

	filter := bson.M{"_id": objId}
	//filter := bson.M{"_id": bson.M{"$eq": objId}}

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancel()

	_, err = col.UpdateOne(ctx, filter, updateString)
	//_, err := col.UpdateByID(ctx, objId, updateString)

	if err != nil {
		return err
	}

	return nil
}

/* TryLogin makes the login to the DB */
func (db *DbNoSql) TryLogin(email string, password string) (mr.User, bool) {
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

/* Delete deletes a tweet in the DB */
func (db *DbNoSql) DeleteTweetFisical(id string) error {
	objId, err := getObjectId(id)

	if err != nil {
		return err
	}

	col := getCollection(db, "twittor", "tweet")

	condition := bson.M{
		"_id": objId,
	}

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancel()

	_, err = col.DeleteOne(ctx, condition)

	return err
}

/* DeleteLogical inactivates a tweet in the DB */
func (db *DbNoSql) DeleteTweetLogical(id string) error {
	objId, err := getObjectId(id)

	if err != nil {
		return err
	}

	col := getCollection(db, "twittor", "tweet")

	condition := bson.M{
		"_id": objId,
	}
	updateString := bson.M{
		"$set": bson.M{"active": false},
	}

	// Also map[string]map[string]bool{"$set": {"active": false}} in the updateString

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancel()

	_, err = col.UpdateOne(ctx, condition, updateString)

	return err
}

/* Get gets an user's tweets from the DB */
func (db *DbNoSql) GetTweets(id string, page int64, limit int64) ([]*mr.Tweet, int64, error) {
	var results []*mr.Tweet
	var tweetsDbResults []*m.Tweet

	objId, err := getObjectId(id)

	if err != nil {
		return results, 0, err
	}

	col := getCollection(db, "twittor", "tweet")
	condition := bson.M{
		"userId": objId,
		"active": true,
	}

	ctxCount, cancelCount := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancelCount()

	total, err := col.CountDocuments(ctxCount, condition)

	if err != nil {
		return results, total, err
	}

	opts := options.Find()

	opts.SetSort(bson.D{{Key: "date", Value: -1}})
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)

	ctxFind, cancelFind := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancelFind()

	cursor, err := col.Find(ctxFind, condition, opts)

	if err != nil {
		return results, total, err
	}

	ctxCursor := context.TODO()

	defer cursor.Close(ctxCursor)

	err = cursor.All(ctxCursor, &tweetsDbResults)

	if err != nil {
		return results, total, err
	}

	for _, tweetModel := range tweetsDbResults {
		tweetRequest := getTweetRequest(*tweetModel)
		results = append(results, &tweetRequest)
	}

	return results, total, nil
}

/* InsertTweet inserts a tweet in the DB */
func (db *DbNoSql) InsertTweet(tweet mr.Tweet) (string, error) {
	tweetModel, err := getTweetModel(tweet)

	if err != nil {
		return "", err
	}

	col := getCollection(db, "twittor", "tweet")
	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancel()

	result, err := col.InsertOne(ctx, tweetModel)

	if err != nil {
		return "", err
	}

	objId := result.InsertedID.(primitive.ObjectID)

	// The same goes with objId.Hex()
	return objId.String(), nil
}

// endregion

// region "Relations"

/* IsRelation verifies if exist a relation in the DB */
func (db *DbNoSql) IsRelation(relation mr.Relation) (bool, mr.Relation, error) {
	var result mr.Relation
	var relationModel m.Relation

	col := getCollection(db, "twittor", "relation")
	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancel()

	err := col.FindOne(ctx, relation).Decode(&relationModel)

	if err != nil && err == mongo.ErrNoDocuments {
		return false, result, nil
	}

	result = getRelationRequest(relationModel)

	return true, result, err
}

/* InsertRelation creates a relation into the DB */
func (db *DbNoSql) InsertRelation(relation mr.Relation) error {
	col := getCollection(db, "twittor", "relation")
	isRelation, relationDb, err := db.IsRelation(relation)

	if err != nil {
		return err
	}

	if !isRelation {
		relationModel, err := getRelationModel(relation)

		if err != nil {
			return err
		}

		ctxInsert, cancelInsert := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

		defer cancelInsert()

		_, err = col.InsertOne(ctxInsert, relationModel)

		return err
	}

	if relationDb.Active {
		return fmt.Errorf("the relation with the user id: %s already exists", relation.UserRelationId)
	}

	updateString := bson.M{
		"$set": bson.M{"active": true},
	}

	id, err := getObjectId(relationDb.Id)

	if err != nil {
		return err
	}

	ctxUpdate, cancelUpdate := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancelUpdate()

	_, err = col.UpdateByID(ctxUpdate, id, updateString)

	return err
}

/* Delete deletes a relation in the DB */
func (db *DbNoSql) DeleteRelationFisical(relation mr.Relation) error {
	relationModel, err := getRelationModel(relation)

	if err != nil {
		return err
	}

	col := getCollection(db, "twittor", "relation")
	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancel()

	_, err = col.DeleteOne(ctx, relationModel)

	return err
}

/* DeleteLogical inactivates a relation in the DB */
func (db *DbNoSql) DeleteRelationLogical(relation mr.Relation) error {
	relationModel, err := getRelationModel(relation)

	if err != nil {
		return err
	}

	col := getCollection(db, "twittor", "relation")
	updateString := map[string]map[string]bool{"$set": {"active": false}}

	// Also bson.M{"$set": bson.M{"active": false},} in the updateString

	ctx, cancel := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancel()

	_, err = col.UpdateOne(ctx, relationModel, updateString)

	return err
}

/* GetUsers gets a list of users */
func (db *DbNoSql) GetUsers(id string, page int64, limit int64, search string, searchType string) ([]*mr.User, int64, error) {
	var results []*mr.User
	var total int64

	col := getCollection(db, "twittor", "users")

	query := bson.M{
		"name": bson.M{"$regex": search, "$options": "im"},
	}

	findOpts := options.Find()

	findOpts.SetSort(bson.D{{Key: "birthDate", Value: -1}})
	findOpts.SetSkip((page - 1) * limit)
	findOpts.SetLimit(limit)

	ctxFind, cancelFind := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

	defer cancelFind()

	cursor, err := col.Find(ctxFind, query, findOpts)

	if err != nil {
		return results, total, err
	}

	ctxCursor := context.TODO()

	defer cursor.Close(ctxCursor)

	var found, include bool

	for cursor.Next(ctxCursor) {
		var result m.User

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

		found, _, err = db.IsRelation(relationRequest)

		if err != nil {
			return results, total, err
		}

		include = false

		if relationRequest.UserRelationId == id {
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

/* GetFollowers gets an user's followers list */
func (db *DbNoSql) GetFollowers(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	var results []*mr.User

	objId, err := getObjectId(id)

	if err != nil {
		return nil, 0, err
	}

	// region Pipeline

	matchId := bson.M{"$match": bson.M{"userId": objId, "active": true}}
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

	dbResults, total, err := getResults[m.User](db, "relation", countPipeline, aggPipeline)

	if err == nil {
		for _, userModel := range dbResults {
			userRequest := getUserRequest(*userModel)
			results = append(results, &userRequest)
		}
	}

	return results, total, err
}

/* GetNotFollowers gets an user's not followers list */
func (db *DbNoSql) GetNotFollowers(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	var results []*mr.User

	objId, err := getObjectId(id)

	if err != nil {
		return nil, 0, err
	}

	// region Pipeline

	matchId := bson.M{"$match": bson.M{"_id": objId}}
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
					"$not": bson.M{"$in": []interface{}{
						"$_id",
						"$$userId",
					}},
				}},
		}},
	}}
	projectArray := bson.M{"$project": bson.M{"u": "$u", "_id": 0}}
	unwind := bson.M{"$unwind": bson.M{
		"path":                       "$u",
		"preserveNullAndEmptyArrays": false,
	}}
	projectUser := bson.M{"$project": bson.M{
		"_id":       "$u._id",
		"name":      "$u.name",
		"lastName":  "$u.lastName",
		"birthDate": "$u.birthDate",
	}}
	matchName := bson.M{"$match": bson.M{
		"_id":  bson.M{"$ne": id},
		"name": bson.M{"$regex": search, "$options": "im"},
	}}
	count := bson.M{"$count": "total"}
	sort := bson.M{"$sort": bson.M{"birthDate": -1}}
	skip := bson.M{"$skip": (page - 1) * limit}
	agLimit := bson.M{"$limit": limit}

	basePipeline := []bson.M{matchId, lookupRelation, lookupUsers, projectArray, unwind, projectUser, matchName}
	countPipeline := append(basePipeline, count)
	aggPipeline := append(basePipeline, sort, skip, agLimit)

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

/* GetUsersTweets returns the followers' tweets */
func (db *DbNoSql) GetUsersTweets(id string, page int64, limit int64, isOnlyTweets bool) (interface{}, int64, error) {
	var results []interface{}
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

	// region "Pipeline"

	conditions = append(conditions, bson.M{"$match": bson.M{"userId": objId, "active": true}})
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
		var dbResults []*m.Tweet

		dbResults, total, err = getResults[m.Tweet](db, "relation", conditionsCount, conditionsAgg)

		if err == nil {
			for _, tweetModel := range dbResults {
				tweetRequest := getTweetRequest(*tweetModel)
				results = append(results, &tweetRequest)
			}
		}
	} else {
		var dbResults []*m.UserTweet

		dbResults, total, err = getResults[m.UserTweet](db, "relation", conditionsCount, conditionsAgg)

		if err == nil {
			for _, userTweetModel := range dbResults {
				userTweetRequest := getUserTweetRequest(*userTweetModel)
				results = append(results, &userTweetRequest)
			}
		}
	}

	return results, total, err
}

// endregion

// region "Helpers"

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

func getRelationModel(requestModel mr.Relation) (m.Relation, error) {
	var relationModel m.Relation

	objId, err := getObjectId(requestModel.Id)

	if err != nil {
		return relationModel, err
	}

	objUserId, err := getObjectId(requestModel.UserId)

	if err != nil {
		return relationModel, err
	}

	objUserRelationId, err := getObjectId(requestModel.UserRelationId)

	if err != nil {
		return relationModel, err
	}

	relationModel = m.Relation{
		Id:             objId,
		UserId:         objUserId,
		UserRelationId: objUserRelationId,
		Active:         requestModel.Active,
	}

	return relationModel, nil
}

func getRelationRequest(relationModel m.Relation) mr.Relation {
	requestModel := mr.Relation{
		Id:             relationModel.Id.Hex(),
		UserId:         relationModel.UserId.Hex(),
		UserRelationId: relationModel.UserRelationId.Hex(),
		Active:         relationModel.Active,
	}

	return requestModel
}

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

func encryptPassword(password string) (string, error) {
	// Minimum - cost: 6
	// Common user - cost: 6
	// Admin user - cost: 8

	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	return string(bytes), err
}

func getCollection(db *DbNoSql, dbName string, colName string) *mongo.Collection {
	database := db.Connection.Database(dbName)
	collection := database.Collection(colName)

	return collection
}

func getResults[T any](db *DbNoSql, colName string, countPipeline []primitive.M, aggPipeline []primitive.M) ([]*T, int64, error) {
	var results []*T
	var totalResult m.TotalResult

	col := getCollection(db, "twittor", colName)

	// region "Count"

	ctxCount, cancelCount := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

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

	ctxAggregate, cancelAggregate := helpers.GetTimeoutCtx(helpers.GetEnvVariable("CTX_TIMEOUT"))

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
