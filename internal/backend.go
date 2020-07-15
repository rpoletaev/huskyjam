package internal

import (
	"time"

	"github.com/rpoletaev/huskyjam/pkg/auth"
)

// Store combines all rdb repositories
type Store interface {
	Accounts() AccountsRepository
	Goods() GoodsRepository
}

type KVStore interface {
	Tokens() TokensRepository
}

// Account struct
type Account struct {
	ID        uint
	Email     string
	Pass      string
	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// AccountsRepository allows interact with Accounts storage
type AccountsRepository interface {
	// Init creates schema. In most projects we will use migration instead
	Init() error
	Create(a *Account) error
	GetByEmail(email string) (*Account, error)
}

// Category struct
type Category struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	DeletedAt *time.Time
}

// Good struct
type Good struct {
	ID         uint
	Name       string
	CreatedAt  time.Time
	DeletedAt  *time.Time
	Categories []Category
}

// GoodsRepository allows interact with Goods storage
type GoodsRepository interface {
	// Init creates schema. In most projects we will use migration instead
	Init() error
	CreateCategory(c *Category) error
	UpdateCategory(c *Category) error
	DeleteCategory(id uint) error
	ListCategories() ([]*Category, error)

	CreateGood(c *Good) error
	UpdateGood(c *Good) error
	DeleteGood(id uint) error
	GoodsByCategory(catID uint) ([]*Good, error)
}

type TokensRepository interface {
	New(claims *auth.SystemClaims) (string, error)
	Get(token string) (*auth.SystemClaims, error)
	Delete(token string) error
}
