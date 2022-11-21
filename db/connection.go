package db

import (
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
	GetProfile(id string) (mr.User, bool, error)
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

/* DbConn is the connection to the database */
var DbConn DbAdapter

/* SetDataBaseConnector sets the connector to the database type */
func SetDataBaseConnector(dbType string) {
	switch dbType {
	case "NoSql":
		dbNoSql := new(dns.DbNoSql)

		dbNoSql.Connect()

		DbConn = dbNoSql
	case "NoSqlV2":
		dbNoSqlV2 := new(dnsv2.DbNoSqlV2)

		dbNoSqlV2.Connect()

		DbConn = dbNoSqlV2
	case "Sql":
		dbSql := new(dr.DbSql)

		dbSql.Connect()

		DbConn = dbSql
	default:
		log.Fatal("No database connector selected")
	}
}
