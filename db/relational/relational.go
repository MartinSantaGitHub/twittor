package relational

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result := db.Connection.WithContext(ctx).Create(&userModel)
	err = result.Error

	if err != nil {
		return "", err
	}

	return strconv.FormatUint(userModel.Id, 10), nil
}

/* IsUser checks that the user already exists in the DB */
func (db *DbSql) IsUser(email string) (bool, mr.User, error) {
	var userModel m.User
	var requestModel mr.User

	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result := db.Connection.WithContext(ctx).Where(&m.User{Email: email}).First(&userModel)
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

	if !user.BirthDate.IsZero() {
		registry["birthDate"] = user.BirthDate
	}

	user.Id = id
	userModel, err := getUserModel(user)

	if err != nil {
		return err
	}

	ctx, cancel := helpers.GetTimeoutCtx(os.Getenv("CTX_TIMEOUT"))

	defer cancel()

	result := db.Connection.WithContext(ctx).Model(&userModel).Updates(registry)
	err = result.Error

	if err != nil {
		return err
	}

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

	if err != nil {
		return err
	}

	return nil
}

/* Get gets an user's tweets from the DB */
func (db *DbSql) GetTweets(id string, page int64, limit int64) ([]*mr.Tweet, int64, error) {
	var results []*mr.Tweet

	log.Fatal("Method not implemented")

	return results, 0, nil
}

/* InsertTweet inserts a tweet in the DB */
func (db *DbSql) InsertTweet(tweet mr.Tweet) (string, error) {
	log.Fatal("Method not implemented")

	// The same goes with objId.Hex()
	return "", nil
}

// endregion

// region "Relations"

/* IsRelation verifies if exist a relation in the DB */
func (db *DbSql) IsRelation(relation mr.Relation) (bool, mr.Relation, error) {
	var result mr.Relation

	log.Fatal("Method not implemented")

	return false, result, nil
}

/* InsertRelation creates a relation into the DB */
func (db *DbSql) InsertRelation(relation mr.Relation) error {
	var err error

	log.Fatal("Method not implemented")

	return err
}

/* DeleteRelation deletes a relation in the DB */
func (db *DbSql) DeleteRelation(relation mr.Relation) error {
	err := db.deleteRelationLogical(relation)

	if err != nil {
		return err
	}

	return nil
}

/* GetUsers gets an user's following or not following list */
func (db *DbSql) GetUsers(id string, page int64, limit int64, search string, searchType string) ([]*mr.User, int64, error) {
	var results []*mr.User

	log.Fatal("Method not implemented")

	return results, 0, nil
}

/* GetFollowing gets an user's following list */
func (db *DbSql) GetFollowing(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	var results []*mr.User

	log.Fatal("Method not implemented")

	return results, 0, nil
}

/* GetNotFollowing gets an user's not following list */
func (db *DbSql) GetNotFollowing(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	var results []*mr.User

	log.Fatal("Method not implemented")

	return results, 0, nil
}

/* GetFollowingTweets returns the following's tweets */
func (db *DbSql) GetFollowingTweets(id string, page int64, limit int64, isOnlyTweets bool) (any, int64, error) {
	var results any
	var total int64
	var err error

	log.Fatal("Method not implemented")

	return results, total, err
}

// func (db *DbSql) GetUsers(id string, page int64, limit int64, search string, searchType string) ([]*mr.User, int64, error) {
// 	var results []*mr.User

// 	log.Fatal("Method not implemented")

// 	return results, 0, nil
// }

// endregion

// region "Helpers"

func (db *DbSql) deleteTweetFisical(id string, userId string) error {
	log.Fatal("Method not implemented")

	return nil
}

func (db *DbSql) deleteTweetLogical(id string, userId string) error {
	log.Fatal("Method not implemented")

	return nil
}

func (db *DbSql) deleteRelationFisical(relation mr.Relation) error {
	log.Fatal("Method not implemented")

	return nil
}

func (db *DbSql) deleteRelationLogical(relation mr.Relation) error {
	log.Fatal("Method not implemented")

	return nil
}

func encryptPassword(password string) (string, error) {
	// Minimum - cost: 6
	// Common user - cost: 6
	// Admin user - cost: 8

	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	return string(bytes), err
}

// endregion
