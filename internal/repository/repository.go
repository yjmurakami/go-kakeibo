package repository

type Repositories struct {
	User UserRepository
}

func NewRepositories() Repositories {
	return Repositories{
		User: &userRepository{},
	}
}

func NewMockRepositories() Repositories {
	return Repositories{
		User: nil,
	}
}
