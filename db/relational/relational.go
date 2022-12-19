package relational

import (
	"fmt"
	"log"
	"os"

	"helpers"
	m "models/relational"
	mr "models/request"

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
	var profile mr.User

	log.Fatal("Method not implemented")

	return profile, false, nil
}

/* InsertUser inserts an user into de DB */
func (db *DbSql) InsertUser(user mr.User) (string, error) {
	log.Fatal("Method not implemented")

	return "", nil
}

/* IsUser checks that the user already exists in the DB */
func (db *DbSql) IsUser(email string) (bool, mr.User, error) {
	var result mr.User
	var err error

	log.Fatal("Method not implemented")

	return false, result, err
}

/* ModifyRegistry modifies a registry in the DB */
func (db *DbSql) ModifyRegistry(id string, user mr.User) error {
	var err error

	log.Fatal("Method not implemented")

	return err
}

/* TryLogin makes the login to the DB */
func (db *DbSql) TryLogin(email string, password string) (mr.User, bool) {
	var user mr.User

	log.Fatal("Method not implemented")

	return user, false
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

// endregion
