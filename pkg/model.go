package pkg

//Model ....
type Model interface {
	ConstructFromMap(map[string]interface{}) (Model, error)

	Add() error
	Delete()

	Update() error
	FindAll()
	FindOne()
}
