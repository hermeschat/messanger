package db

import "hermes/pkg/db/mongo"

type Mongo struct {
	repos map[string]Repository
}

func (m *Mongo) Repo(name string) Repository {
	repo, exists := m.repos[name]
	if !exists {
		return nil
	}
	return repo
}

func (m *Mongo) SetRepo(name string, repo Repository) {
	m.repos[name] = repo
}
func initMongoRepos(h HermesDatabase) {
	h.SetRepo("channels", &mongo.Channelrepository{})
	h.SetRepo("messages", &mongo.Messagerepository{})
	h.SetRepo("users", mongo.Userrepository{})
}
