package model

// Modeler is the interface that all database types will need to implement
// the return values need to be a generic so we type assert them
type Modeler[T any, V any] interface {
	GetEntry(key V) (T, error)
	// filter will have all the search parameters
	GetEntryAdvanced(filter V, sort V) ([]T, error)
	AddEntry(val V) error
	RemoveEntry(val V) error
	ModifyEntry(filter V, val V) error
}
