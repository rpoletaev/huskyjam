package internal

// Store combines all rdb repositories
type Store interface {
	Init() error
	Accounts() AccountsRepository
	Goods() GoodsRepository
}

// Account struct
type Account struct {
	ID    string
	Email string
	Pass  string
}

// AccountsRepository allows interact with Accounts storage
type AccountsRepository interface {
	Create(a *Account) error
	GetByEmail(email string) (*Account, error)
}

// Category struct
type Category struct {
	ID   string
	Name string
}

// Good struct
type Good struct {
	ID   string
	Name string
}

// GoodsRepository allows interact with Goods storage
type GoodsRepository interface {
	CreateCategory(c *Category) error
	UpdateCategory(c *Category) error
	DeleteCategory(id string) error
	ListCategories() ([]*Category, error)

	CreateGood(c *Good) error
	UpdateGood(c *Good) error
	DeleteGood(id string) error
	GoodsByCategory(catID string) ([]*Good, error)
	SetGoodCategories(goodID string, categories []string) error
}
