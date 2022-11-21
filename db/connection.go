package db

import (
	"helpers"
	"log"

	mr "models/request"

	dns "db/nosql"
	dnsv2 "db/nosqlv2"
	dr "db/relational"
)

type DbAdapter interface {
	// Connection
	Connect() error
	IsConnection() bool

	// Users
	GetProfile(id string) (mr.User, error)
	InsertUser(user mr.User) (string, error)
	IsUser(email string) (bool, mr.User, error)
	ModifyRegistry(id string, user mr.User) error
	TryLogin(email string, password string) (mr.User, bool)

	// Tweets
	DeleteTweetFisical(id string, userId string) error
	DeleteTweetLogical(id string, userId string) error
	GetTweets(id string, page int64, limit int64) ([]*mr.Tweet, int64, error)
	InsertTweet(tweet mr.Tweet) (string, error)

	// Relations
	GetRelation(relation mr.Relation) (bool, mr.Relation, error)
	InsertRelation(relation mr.Relation) error
	DeleteRelationFisical(relation mr.Relation) error
	DeleteRelationLogical(relation mr.Relation) error
	GetFollowers(id string, page int64, limit int64, search string) ([]*mr.User, int64, error)
	GetNotFollowers(id string, page int64, limit int64, search string) ([]*mr.User, int64, error)
	GetUsersTweets(id string, page int64, limit int64, isOnlyTweets bool) (interface{}, int64, error)
	//GetUsers(id string, page int64, limit int64, search string, searchType string) ([]*mr.User, int64, error)
}

func getDataBaseConnector(dbType string) DbAdapter {
	switch dbType {
	case "NoSql":
		dbNoSql := new(dns.DbNoSql)

		dbNoSql.Connect()

		return dbNoSql
	case "NoSqlV2":
		dbNoSqlV2 := new(dnsv2.DbNoSqlV2)

		dbNoSqlV2.Connect()

		return dbNoSqlV2
	case "Sql":
		dbSql := new(dr.DbSql)

		dbSql.Connect()

		return dbSql
	default:
		log.Fatal("No database connector selected")

		return nil
	}
}

var dbType string = helpers.GetEnvVariable("DB_TYPE")

/* DbConn is the connection to the database */
var DbConn DbAdapter = getDataBaseConnector(dbType)
