package entity

// Setable ...
type Setable interface {
	GetDataToSet() *map[string]interface{}
}
