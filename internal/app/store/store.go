package store

// store ...
type Store interface {
	User() UserRepository
	Service() ServiceRepository
}
