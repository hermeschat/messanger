package db

import "hermes/pkg/db/mongo"

type Repository interface {
	//Name returns name of the repository
	Name() string
	//Find runs a query using primary key
	Find(id string) (map[string]interface{}, error)
	//Get runs query on db
	Get(query map[string]interface{}) ([]map[string]interface{}, error)
	//Update runs an update query with given selector and update map
	Update(selector map[string]interface{}, update map[string]interface{}) (int, error)
	//Add adds a new object to db
	Add(object map[string]interface{}) (string, error)
	//Delete deletes a record
	Delete(query map[string]interface{}) error
}

type HermesDatabase struct {
	Messages Repository
	Channels Repository
	Users    Repository
}

func initMongoRepos(h *HermesDatabase) {
	h.Channels = mongo.NewChannelrepository()
	h.Messages = mongo.NewMessagerepository()
	h.Users = mongo.NewUserrepository()
}

var hermesDatabaseInstance *HermesDatabase

//Gets a Instance instance
func Instance() HermesDatabase {
	return *hermesDatabaseInstance
}

func Init() {
	hermesDatabaseInstance = new(HermesDatabase)
	initMongoRepos(hermesDatabaseInstance)
}
