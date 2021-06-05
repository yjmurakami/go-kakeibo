package repository

type Repositories struct {
	Category    CategoryRepository
	User        UserRepository
	Transaction TransactionRepository
}

func NewRepositories() Repositories {
	return Repositories{
		Category:    &categoryRepository{},
		User:        &userRepository{},
		Transaction: &transactionRepository{},
	}
}

func NewMockRepositories() Repositories {
	return Repositories{
		Category:    nil,
		User:        nil,
		Transaction: nil,
	}
}
