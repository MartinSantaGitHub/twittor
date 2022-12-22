package relational

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"helpers"
	m "models/relational"
	mr "models/request"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbSql struct {
	Connection *gorm.DB
}

// region "Connection"

/* Connect connects to the database */
func (db *DbSql) Connect() error {
	host := os.Getenv("DB_REL_HOST")
	port := os.Getenv("DB_REL_PORT")
	user := os.Getenv("DB_REL_USER")
	pass := os.Getenv("DB_REL_PASS")
	dbName := os.Getenv("DB_REL_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, pass, dbName, port)
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	db.Connection = client

	client.AutoMigrate(&m.User{})
	client.AutoMigrate(&m.Relation{})
	client.AutoMigrate(&m.Tweet{})

	return nil
}

/* IsConnection makes a ping to the Database */
func (db *DbSql) IsConnection() bool {
	sqlDb, err := db.Connection.DB()

	if err != nil {
		fmt.Println(err.Error())

		return false
	}

	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	err = sqlDb.PingContext(ctx)

	if err != nil {
		fmt.Println(err.Error())

		return false
	}

	return true
}

// endregion

// region "Users"

/* GetProfile gets a profile in the DB */
func (db *DbSql) GetProfile(id string) (mr.User, bool, error) {
	var profileRequest mr.User
	var profileModel m.User

	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result := db.Connection.WithContext(ctx).First(&profileModel, id)
	err := result.Error

	if err != nil && err == gorm.ErrRecordNotFound {
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
func (db *DbSql) InsertUser(user mr.User) (string, error) {
	user.Password, _ = encryptPassword(user.Password)
	userModel, err := getUserModel(user)

	if err != nil {
		return "", err
	}

	tx := db.Connection.Begin()
	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result := tx.WithContext(ctx).Create(&userModel)
	err = result.Error

	if err != nil {
		tx.Rollback()

		return "", err
	}

	tx.Commit()

	return strconv.FormatUint(userModel.Id, 10), nil
}

/* IsUser checks that the user already exists in the DB */
func (db *DbSql) IsUser(email string) (bool, mr.User, error) {
	var userModel m.User
	var requestModel mr.User

	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result := db.Connection.WithContext(ctx).
		Where(&m.User{Email: email}).
		First(&userModel)
	err := result.Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return false, requestModel, nil
	}

	requestModel = getUserRequest(userModel)

	return true, requestModel, err
}

/* ModifyRegistry modifies a registry in the DB */
func (db *DbSql) ModifyRegistry(id string, user mr.User) error {
	registry := make(map[string]any)

	if len(user.Name) > 0 {
		registry["name"] = user.Name
	}

	if len(user.LastName) > 0 {
		registry["last_name"] = user.LastName
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
		registry["web_site"] = user.WebSite
	}

	if !user.BirthDate.IsZero() {
		registry["birth_date"] = user.BirthDate
	}

	userId, _ := getUintId(id)
	tx := db.Connection.Begin()
	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result := tx.WithContext(ctx).Model(&m.User{Id: userId}).Updates(registry)
	err := result.Error

	if err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()

	return nil
}

/* TryLogin makes the login to the DB */
func (db *DbSql) TryLogin(email string, password string) (mr.User, bool) {
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
func (db *DbSql) DeleteTweet(id string, userId string) error {
	err := db.deleteTweetLogical(id, userId)

	return err
}

/* Get gets an user's tweets from the DB */
func (db *DbSql) GetTweets(id string, page int64, limit int64) ([]*mr.Tweet, int64, error) {
	var results []*mr.Tweet
	var tweetsDbResults []m.Tweet
	var total int64

	uintUserId, err := getUintId(id)

	if err != nil {
		return results, 0, err
	}

	filter := &m.Tweet{UserId: uintUserId, Active: true}

	ctxCount, cancelCount := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelCount()

	result := db.Connection.WithContext(ctxCount).
		Model(&m.Tweet{}).
		Where(filter).
		Count(&total)
	err = result.Error

	if err != nil {
		return results, 0, err
	}

	offset := int((page - 1) * limit)
	limitInt := int(limit)
	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result = db.Connection.WithContext(ctx).
		Where(filter).
		Order("date desc").
		Offset(offset).
		Limit(limitInt).
		Find(&tweetsDbResults)
	err = result.Error

	if err != nil {
		return results, 0, err
	}

	for _, tweetModel := range tweetsDbResults {
		tweetRequest := getTweetRequest(tweetModel)
		results = append(results, &tweetRequest)
	}

	return results, total, nil
}

/* InsertTweet inserts a tweet in the DB */
func (db *DbSql) InsertTweet(tweet mr.Tweet) (string, error) {
	tweetModel, err := getTweetModel(tweet)

	if err != nil {
		return "", err
	}

	tx := db.Connection.Begin()
	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	err = tx.WithContext(ctx).Model(&m.User{Id: tweetModel.UserId}).Association("Tweets").Append(&tweetModel)

	if err != nil {
		tx.Rollback()

		return "", err
	}

	tx.Commit()

	return strconv.FormatUint(tweetModel.Id, 10), nil
}

// endregion

// region "Relations"

/* IsRelation verifies if exist a relation in the DB */
func (db *DbSql) IsRelation(relation mr.Relation) (bool, mr.Relation, error) {
	var result mr.Relation
	var relationModel m.Relation

	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	results := db.Connection.WithContext(ctx).
		Where(map[string]string{"user_id": relation.UserId, "following_id": relation.UserRelationId}).
		First(&relationModel)
	err := results.Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return false, result, nil
	}

	result = getRelationRequest(relationModel)

	return true, result, err
}

/* InsertRelation creates a relation into the DB */
func (db *DbSql) InsertRelation(relation mr.Relation) error {
	_, isFound, err := db.GetProfile(relation.UserRelationId)

	if err != nil {
		return err
	}

	if !isFound {
		return errors.New("no registry found in the DB")
	}

	isFound, relationDb, err := db.IsRelation(relation)

	if err != nil {
		return err
	}

	if !isFound {
		userId, _ := getUintId(relation.UserId)
		userRelationId, err := getUintId(relation.UserRelationId)

		if err != nil {
			return err
		}

		tx := db.Connection.Begin()
		ctxInsert, cancelInsert := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

		defer cancelInsert()

		err = tx.WithContext(ctxInsert).Model(&m.User{Id: userId}).Association("Following").Append(&m.User{Id: userRelationId})

		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}

		return err
	}

	if relationDb.Active {
		return fmt.Errorf("the relation with the user id = %s already exists", relation.UserRelationId)
	}

	err = db.updateRelation(relation, true)

	return err
}

/* DeleteRelation deletes a relation in the DB */
func (db *DbSql) DeleteRelation(relation mr.Relation) error {
	err := db.updateRelation(relation, false)

	return err
}

/* GetUsers gets an user's following or not following list */
func (db *DbSql) GetUsers(id string, page int64, limit int64, search string, searchType string) ([]*mr.User, int64, error) {
	var results []*mr.User
	var users []m.User
	var total int64

	offset := int((page - 1) * limit)
	limitInt := int(limit)
	ctxFind, cancelFind := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelFind()

	result := db.Connection.WithContext(ctxFind).
		Where("id <> ? AND lower(name) LIKE ?", id, "%"+strings.ToLower(search)+"%").
		Order("birth_date desc").Offset(offset).
		Limit(limitInt).
		Find(&users)
	err := result.Error

	if err != nil {
		return results, total, err
	}

	var include, isRelation bool

	for _, user := range users {
		isRelation = false
		include = false

		if err != nil {
			return results, total, err
		}

		userRequest := getUserRequest(user)

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

	// This total only reflects the total returned registries
	total = int64(len(results))

	return results, total, nil
}

/* GetFollowing gets an user's following list */
func (db *DbSql) GetFollowing(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	var results []*mr.User
	var user m.User
	var total int64

	ctxCount, cancelCount := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelCount()

	filter := "lower(name) LIKE ?"
	condition := "%" + strings.ToLower(search) + "%"
	result := db.Connection.WithContext(ctxCount).
		Preload("Following", filter, condition).
		First(&user, id)
	err := result.Error

	if err != nil {
		return results, total, err
	}

	total = int64(len(user.Following))

	offset := int((page - 1) * limit)
	limitInt := int(limit)
	ctxFind, cancelFind := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelFind()

	result = db.Connection.WithContext(ctxFind).
		Preload("Following", func(db *gorm.DB) *gorm.DB {
			return db.Where(filter, condition).
				Order("birth_date desc").
				Offset(offset).
				Limit(limitInt).
				Select("id", "name", "last_name", "birth_date")
		}).First(&user, id)
	err = result.Error

	if err == nil {
		for _, userModel := range user.Following {
			userRequest := getUserRequest(userModel)
			results = append(results, &userRequest)
		}
	}

	return results, total, err
}

/* GetNotFollowing gets an user's not following list */
func (db *DbSql) GetNotFollowing(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	var results []*mr.User
	var users []m.User
	var total int64

	uintId, _ := getUintId(id)

	ctxNewUsers, cancelNewUsers := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelNewUsers()

	err := db.Connection.WithContext(ctxNewUsers).
		Model(&m.User{Id: uintId}).
		Select("id").
		Association("Following").
		Find(&users)

	if err != nil {
		return results, total, err
	}

	notIds := []uint64{uintId}

	for _, user := range users {
		notIds = append(notIds, user.Id)
	}

	ctxCount, cancelCount := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelCount()

	filter := "lower(name) LIKE ?"
	condition := "%" + strings.ToLower(search) + "%"
	result := db.Connection.WithContext(ctxCount).
		Model(&m.User{}).
		Not(notIds).
		Where(filter, condition).
		Count(&total)
	err = result.Error

	if err != nil {
		return results, total, err
	}

	offset := int((page - 1) * limit)
	limitInt := int(limit)
	ctxFind, cancelFind := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelFind()

	result = db.Connection.WithContext(ctxFind).
		Model(&m.User{}).
		Not(notIds).
		Where(filter, condition).
		Order("birth_date desc").
		Offset(offset).
		Limit(limitInt).
		Select("id", "name", "last_name", "birth_date").
		Find(&users)
	err = result.Error

	if err == nil {
		for _, userModel := range users {
			userRequest := getUserRequest(userModel)
			results = append(results, &userRequest)
		}
	}

	return results, total, err
}

/* GetFollowingTweets returns the following's tweets */
func (db *DbSql) GetFollowingTweets(id string, page int64, limit int64, isOnlyTweets bool) (any, int64, error) {
	var results any
	var total int64
	var err error

	log.Fatal("Method not implemented")

	return results, total, err
}

// endregion

// region "Helpers"

func (db *DbSql) deleteTweetFisical(id string, userId string) error {
	var tweetModel m.Tweet

	uintId, err := getUintId(id)

	if err != nil {
		return err
	}

	ctxFind, cancelFind := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelFind()

	result := db.Connection.WithContext(ctxFind).First(&tweetModel, uintId)
	err = result.Error

	if err != nil {
		return err
	}

	uintUserId, _ := getUintId(userId)

	if uintUserId != tweetModel.UserId {
		return errors.New("invalid operation - cannot delete a non-owner tweet")
	}

	tx := db.Connection.Begin()
	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result = tx.WithContext(ctx).Delete(&m.Tweet{}, uintId)
	err = result.Error

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return err
}

func (db *DbSql) deleteTweetLogical(id string, userId string) error {
	var tweetModel m.Tweet

	uintId, err := getUintId(id)

	if err != nil {
		return err
	}

	ctxFind, cancelFind := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancelFind()

	result := db.Connection.WithContext(ctxFind).First(&tweetModel, uintId)
	err = result.Error

	if err != nil {
		return err
	}

	uintUserId, _ := getUintId(userId)

	if uintUserId != tweetModel.UserId {
		return errors.New("invalid operation - cannot delete a non-owner tweet")
	}

	tx := db.Connection.Begin()
	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result = tx.WithContext(ctx).Model(&tweetModel).Update("active", false)
	err = result.Error

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return err
}

func (db *DbSql) deleteRelationFisical(relation mr.Relation) error {
	tx := db.Connection.Begin()
	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result := tx.WithContext(ctx).Where(map[string]string{"user_id": relation.UserId, "following_id": relation.UserRelationId}).Delete(&m.Relation{})
	err := result.Error

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return err
}

func (db *DbSql) updateRelation(relation mr.Relation, value bool) error {
	tx := db.Connection.Begin()
	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result := tx.WithContext(ctx).Model(&m.Relation{}).Where(map[string]string{"user_id": relation.UserId, "following_id": relation.UserRelationId}).Update("active", value)
	err := result.Error

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

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

// func getResults[T any](db *DbSql, colName string, countPipeline []primitive.M, aggPipeline []primitive.M) ([]*T, int64, error) {
// 	var results []*T
// 	var totalResult m.TotalResult

// 	col := getCollection(db, "twittor", colName)

// 	// region "Count"

// 	ctxCount, cancelCount := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

// 	defer cancelCount()

// 	curCount, err := col.Aggregate(ctxCount, countPipeline)

// 	if err != nil {
// 		return results, totalResult.Total, err
// 	}

// 	ctxCurCount := context.TODO()

// 	defer curCount.Close(ctxCurCount)

// 	if curCount.Next(ctxCurCount) {
// 		err = curCount.Decode(&totalResult)
// 	}

// 	if err != nil {
// 		return results, totalResult.Total, err
// 	}

// 	// endregion

// 	// region "Aggregate"

// 	ctxAggregate, cancelAggregate := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

// 	defer cancelAggregate()

// 	curAgg, err := col.Aggregate(ctxAggregate, aggPipeline)

// 	if err != nil {
// 		return results, totalResult.Total, err
// 	}

// 	ctxCurAgg := context.TODO()

// 	defer curAgg.Close(ctxCurAgg)

// 	err = curAgg.All(ctxCurAgg, &results)

// 	if err != nil {
// 		return results, totalResult.Total, err
// 	}

// 	// endregion

// 	return results, totalResult.Total, nil
// }

// endregion
