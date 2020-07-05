package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rpoletaev/huskyjam/internal"
	httpauth "github.com/rpoletaev/huskyjam/pkg/auth/http"
	"github.com/rs/zerolog"

	// "github.com/rpoletaev/huskyjam/pkg/auth/jwt"
	"github.com/rpoletaev/huskyjam/pkg/auth"
)

type Config struct {
	Port string `envconfig:"PORT"`
}

type Api struct {
	*Config
	accounts   *AccountHandler `wire:"-"`
	goods      *GoodsHandler   `wire:"-"`
	server     *http.Server    `wire:"-"`
	Tokens     auth.Tokens
	Store      internal.Store
	HashHelper PassHashHelper
	KVStore    internal.KVStore
	Logger     zerolog.Logger
}

// Accounts builds and returns AccountsHandler
func (api *Api) Accounts() *AccountHandler {
	if api.accounts != nil {
		return api.accounts
	}

	api.accounts = &AccountHandler{
		Store:          api.Store.Accounts(),
		PassHashHelper: api.HashHelper,
		RefreshRepo:    api.KVStore.Tokens(),
		Auth:           api.Tokens,
		Log:            api.Logger.With().Str("service", "account").Logger(),
	}

	return api.accounts
}

// Goods builds and returns GoodsHandler
func (api *Api) Goods() *GoodsHandler {
	if api.goods != nil {
		return api.goods
	}

	api.goods = &GoodsHandler{
		Store: api.Store.Goods(),
		Log:   api.Logger.With().Str("service", "goods").Logger(),
	}

	return api.goods
}

// Server builds router and returns ready for listening server
func (api *Api) Server() *http.Server {
	if api.server != nil {
		return api.server
	}

	r := mux.NewRouter()

	// accrouter handles auth endpoints
	accrouter := r.PathPrefix("/auth").Subrouter()
	accrouter.HandleFunc("/signup", api.Accounts().Signup).Methods(http.MethodPost)
	accrouter.HandleFunc("/signin", api.Accounts().Signin).Methods(http.MethodPost)
	accrouter.HandleFunc("/refresh", api.Accounts().Refresh).Methods(http.MethodPost)

	// admrouter handles endpoints for a goods management
	admrouter := r.PathPrefix("/admin").Subrouter()
	admrouter.HandleFunc("/goods", api.Goods().CreateGood).Methods(http.MethodPost)
	admrouter.HandleFunc("/goods", api.Goods().UpdateGood).Methods(http.MethodPut)
	admrouter.HandleFunc("/goods/{id}", api.Goods().DeleteGood).Methods(http.MethodDelete)
	admrouter.HandleFunc("/categories", api.Goods().CreateCategory).Methods(http.MethodPost)
	admrouter.HandleFunc("/categories", api.Goods().UpdateCategory).Methods(http.MethodPut)
	admrouter.HandleFunc("/categories/{id}", api.Goods().DeleteCategory).Methods(http.MethodDelete)
	admrouter = httpauth.WithAuth(api.Tokens)(admrouter).(*mux.Router)

	// shoprouter provides access to view products
	shoprouter := r.PathPrefix("/shop").Subrouter()
	shoprouter.HandleFunc("/categories", api.Goods().CategoriesList).Methods(http.MethodGet)
	shoprouter.HandleFunc("/categories/{category}/goods", api.Goods().CategoryGoods).Methods(http.MethodGet)

	api.server = &http.Server{
		Addr:    api.Port,
		Handler: setupGlobalMiddleware(r),
	}

	return api.server
}
