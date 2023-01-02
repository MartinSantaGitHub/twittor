package mocks

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"db"
	mr "models/request"

	"golang.org/x/crypto/bcrypt"
)

type DbMock struct {
	Users          map[string]*mr.User
	Tweets         []*mr.Tweet
	Relations      []*mr.Relation
	IsError        bool
	IsConnected    bool
	IdUserCounter  int
	IdTweetCounter int
}

// region "Connection"

func (db *DbMock) Connect() error {
	if db.IsError {
		return fmt.Errorf("Connection Error")
	}

	db.Users = make(map[string]*mr.User)

	return nil
}

func (db *DbMock) IsConnection() bool {
	return db.IsConnected
}

// endregion

// region "Users"

func (db *DbMock) GetProfile(id string) (mr.User, bool, error) {
	var profileModel *mr.User
	var profileRequest mr.User

	if db.IsError {
		return profileRequest, false, fmt.Errorf("Error!")
	}

	profileModel = db.Users[id]

	if exists := profileModel == nil; exists {
		return *profileModel, exists, nil
	}

	return profileRequest, false, nil
}

func (db *DbMock) InsertUser(user mr.User) (string, error) {
	cost := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), cost)

	if err != nil {
		return "", fmt.Errorf("Error!")
	}

	db.IdUserCounter++

	user.Id = strconv.Itoa(db.IdUserCounter)
	user.Password = string(bytes)

	db.Users[user.Id] = &user

	return user.Id, nil
}

func (db *DbMock) IsUser(email string) (bool, mr.User, error) {
	var profileRequest mr.User
	var isUser bool

	if db.IsError {
		return false, profileRequest, fmt.Errorf("Error!")
	}

	for _, value := range db.Users {
		isUser = value.Email == email

		if isUser {
			profileRequest = *value
			break
		}
	}

	return isUser, profileRequest, nil
}

func (db *DbMock) ModifyRegistry(id string, user mr.User) error {
	if db.IsError {
		return fmt.Errorf("Error!")
	}

	userDb := db.Users[id]

	userDb.Name = user.Name

	return nil
}

func (db *DbMock) TryLogin(email string, password string) (mr.User, bool) {
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

func (db *DbMock) DeleteTweet(id string, userId string) error {
	var tweetModel *mr.Tweet

	if db.IsError {
		return fmt.Errorf("Error!")
	}

	for _, t := range db.Tweets {
		if t.Id == id {
			tweetModel = t
			break
		}
	}

	if userId != tweetModel.UserId {
		return fmt.Errorf("invalid operation - cannot delete a non-owner tweet")
	}

	tweetModel.Active = false

	return nil
}

func (db *DbMock) GetTweets(id string, page int64, limit int64) ([]*mr.Tweet, int64, error) {
	tweets := []*mr.Tweet{}

	if db.IsError {
		return tweets, 0, fmt.Errorf("Error!")
	}

	offset := (page - 1) * limit
	count := int64(0)
	total := int64(0)

	for _, tweet := range db.Tweets {
		if tweet.UserId == id && tweet.Active {
			total++
		}
	}

	sort.Slice(db.Tweets, func(i, j int) bool {
		return db.Tweets[i].Date.After(db.Tweets[j].Date)
	})

	for _, tweet := range db.Tweets[offset:] {
		if tweet.UserId == id && tweet.Active && count < limit {
			tweets = append(tweets, tweet)
			count++
		}
	}

	return tweets, total, nil
}

func (db *DbMock) InsertTweet(tweet mr.Tweet) (string, error) {
	db.IdTweetCounter++

	tweet.Id = strconv.Itoa(db.IdTweetCounter)
	db.Tweets = append(db.Tweets, &tweet)

	return tweet.Id, nil
}

// endregion

// region "Relations"

func (db *DbMock) IsRelation(relation mr.Relation) (bool, mr.Relation, error) {
	var relationModel mr.Relation
	var isFound bool

	if db.IsError {
		return false, relationModel, fmt.Errorf("Error!")
	}

	for _, r := range db.Relations {
		if relation.UserId == r.UserId && relation.UserRelationId == r.UserRelationId {
			isFound = true
			relationModel = *r
		}
	}

	return isFound, relationModel, nil
}

func (db *DbMock) InsertRelation(relation mr.Relation) error {
	_, isFound, err := db.GetProfile(relation.UserRelationId)

	if err != nil {
		return err
	}

	if !isFound {
		return fmt.Errorf("no registry found in the DB")
	}

	isFound, relationDb, err := db.IsRelation(relation)

	if err != nil {
		return err
	}

	if !isFound {
		db.Relations = append(db.Relations, &mr.Relation{
			UserId:         relationDb.UserId,
			UserRelationId: relationDb.UserRelationId})

		return nil
	}

	if relationDb.Active {
		return fmt.Errorf("the relation with the user id = %s already exists", relation.UserRelationId)
	}

	db.updateRelation(relationDb, true)

	return nil
}

func (db *DbMock) DeleteRelation(relation mr.Relation) error {
	if db.IsError {
		return fmt.Errorf("Error!")
	}

	db.updateRelation(relation, false)

	return nil
}

func (db *DbMock) GetUsers(id string, page int64, limit int64, search string, searchType string) ([]*mr.User, int64, error) {
	var results []*mr.User
	var users []*mr.User
	var total int64

	if db.IsError {
		return results, total, fmt.Errorf("Error!")
	}

	for _, u := range db.Users {
		if strings.Contains(strings.ToLower(u.Name), strings.ToLower(search)) {
			users = append(users, u)
		}
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].BirthDate.After(users[j].BirthDate)
	})

	offset := (page - 1) * limit
	usersCount := int64(0)
	usersFiltered := []*mr.User{}

	for _, u := range users[offset:] {
		if usersCount < limit {
			usersFiltered = append(usersFiltered, u)
			usersCount++
		}
	}

	var include, isRelation bool

	for _, user := range usersFiltered {
		isRelation = false
		include = false

		relationRequest := mr.Relation{
			UserId:         id,
			UserRelationId: user.Id,
			Active:         true,
		}

		isFound, relationDb, _ := db.IsRelation(relationRequest)

		if isFound {
			isRelation = relationDb.Active
		}

		if relationRequest.UserRelationId == id {
			include = false
		} else if (searchType == "new" && !isRelation) || (searchType == "follow" && isRelation) {
			include = true
		}

		if include {
			userRequest := *user

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

func (db *DbMock) GetFollowing(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	return db.getUsers(id, page, limit, search, true)
}

func (db *DbMock) GetNotFollowing(id string, page int64, limit int64, search string) ([]*mr.User, int64, error) {
	return db.getUsers(id, page, limit, search, false)
}

func (db *DbMock) GetFollowingTweets(id string, page int64, limit int64, isOnlyTweets bool) (any, int64, error) {
	var results any
	var userRelationIds []string
	var tweets []*mr.Tweet
	var total int64

	if db.IsError {
		return results, total, fmt.Errorf("Error!")
	}

	for _, r := range db.Relations {
		if r.UserId == id {
			userRelationIds = append(userRelationIds, r.UserRelationId)
		}
	}

	for _, id := range userRelationIds {
		for _, tweet := range db.Tweets {
			if tweet.UserId == id {
				tweets = append(tweets, tweet)
			}
		}
	}

	total = int64(len(tweets))

	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].Date.After(tweets[j].Date)
	})

	offset := (page - 1) * limit
	tweetsCount := int64(0)

	if isOnlyTweets {
		var reqResults []*mr.Tweet

		for _, t := range tweets[offset:] {
			if tweetsCount < limit {
				reqResults = append(reqResults, t)

				tweetsCount++
			}
		}

		results = reqResults
	} else {
		var reqResults []*mr.UserTweet

		for _, t := range tweets[offset:] {
			if tweetsCount < limit {
				ut := &mr.UserTweet{}

				ut.UserId = id
				ut.UserRelationId = t.UserId
				ut.Tweet.Id = t.Id
				ut.Tweet.Message = t.Message
				ut.Tweet.Date = t.Date

				reqResults = append(reqResults, ut)

				tweetsCount++
			}
		}

		results = reqResults
	}

	return results, total, nil
}

// endregion

// region "Helpers"

func (db *DbMock) updateRelation(relation mr.Relation, value bool) {
	for _, r := range db.Relations {
		if relation.UserId == r.UserId && relation.UserRelationId == r.UserRelationId {
			r.Active = value
		}
	}
}

func (db *DbMock) getUsers(id string, page int64, limit int64, search string, isFollowing bool) ([]*mr.User, int64, error) {
	var results []*mr.User
	var users []mr.User
	var userRelationIds []string
	var total int64

	if db.IsError {
		return results, total, fmt.Errorf("Error!")
	}

	for _, r := range db.Relations {
		if r.UserId == id {
			userRelationIds = append(userRelationIds, r.UserRelationId)
		}
	}

	for _, id := range userRelationIds {
		for _, u := range db.Users {
			if ((id == u.Id && isFollowing) || (id != u.Id && !isFollowing)) && strings.Contains(strings.ToLower(u.Name), strings.ToLower(search)) {
				users = append(users, *u)
				break
			}
		}
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].BirthDate.After(users[j].BirthDate)
	})

	total = int64(len(users))
	offset := (page - 1) * limit
	usersCount := int64(0)

	for _, u := range users[offset:] {
		if usersCount < limit {
			results = append(results, &u)
			usersCount++
		}
	}

	return results, total, nil
}

// endregion

func Init() {
	db.DbConn = &DbMock{}
}
