package relational

import (
	"log"
	mr "models/request"
)

type DbSql struct {
	Connection interface{}
}

// region "Connection"

/* Connect connects to the database */
func (db *DbSql) Connect() error {
	log.Fatal("Method not implemented")

	return nil
}

/* CheckConnection makes a ping to the Database */
func (db *DbSql) IsConnection() bool {
	log.Fatal("Method not implemented")

	return false
}

// endregion

// region "Users"

/* GetProfile gets a profile in the DB */
func (db *DbSql) GetProfile(id string) (mr.User, error) {
	var profile mr.User

	log.Fatal("Method not implemented")

	return profile, nil
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

/* Delete deletes a tweet in the DB */
func (db *DbSql) DeleteTweetFisical(id string, userId string) error {
	log.Fatal("Method not implemented")

	return nil
}

/* DeleteLogical inactivates a tweet in the DB */
func (db *DbSql) DeleteTweetLogical(id string, userId string) error {
	log.Fatal("Method not implemented")

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
func (db *DbSql) GetRelation(relation mr.Relation) (bool, mr.Relation, error) {
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

/* Delete deletes a relation in the DB */
func (db *DbSql) DeleteRelationFisical(relation mr.Relation) error {
	log.Fatal("Method not implemented")

	return nil
}

/* DeleteLogical inactivates a relation in the DB */
func (db *DbSql) DeleteRelationLogical(relation mr.Relation) error {
	log.Fatal("Method not implemented")

	return nil
}

/* GetFollowers gets an user's followers list */
func (db *DbSql) GetFollowers(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	var results []*mr.User

	log.Fatal("Method not implemented")

	return results, 0, nil
}

/* GetNotFollowers gets an user's not followers list */
func (db *DbSql) GetNotFollowers(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	var results []*mr.User

	log.Fatal("Method not implemented")

	return results, 0, nil
}

/* GetUsersTweets returns the followers' tweets */
func (db *DbSql) GetUsersTweets(id string, page int64, limit int64, isOnlyTweets bool) (interface{}, int64, error) {
	var results []interface{}
	var total int64
	var err error

	log.Fatal("Method not implemented")

	return results, total, err
}

// endregion
