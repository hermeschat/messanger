package pkg

//Model ....
type Model interface {
	ToMap() (map[string]interface{}, error)
}
