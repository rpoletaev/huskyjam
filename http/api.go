package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rpoletaev/huskyjam/internal"
	httpauth "github.com/rpoletaev/huskyjam/pkg/auth/http"
	"github.com/rs/zerolog"

	_ "github.com/rpoletaev/huskyjam/cmd/service/docs"
	"github.com/rpoletaev/huskyjam/pkg/auth"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Port string `envconfig:"PORT"`
}

type Api struct {
	*Config
	Accounts   *AccountHandler
	Goods      *GoodsHandler
	Tokens     auth.Tokens
	Store      internal.Store
	HashHelper PassHashHelper
	KVStore    internal.KVStore
	Logger     zerolog.Logger
	server     *http.Server `wire:"-"`
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @BasePath /
// @securityDefinitions.apikey X-Auth-Key
// @in header
// @name Authorization

// Server builds router and returns ready for listening server
func (api *Api) Server() *http.Server {
	if api.server != nil {
		return api.server
	}

	r := mux.NewRouter()

	// accrouter handles auth endpoints
	accrouter := r.PathPrefix("/auth").Subrouter()
	accrouter.HandleFunc("/signup", api.Accounts.Signup).Methods(http.MethodPost)
	accrouter.HandleFunc("/signin", api.Accounts.Signin).Methods(http.MethodPost)
	accrouter.HandleFunc("/refresh", api.Accounts.Refresh).Methods(http.MethodPost)

	// admrouter handles endpoints for a goods management
	admrouter := r.PathPrefix("/admin").Subrouter()
	admrouter.Use(httpauth.WithAuth(api.Tokens))
	admrouter.HandleFunc("/goods", api.Goods.CreateGood).Methods(http.MethodPost)
	admrouter.HandleFunc("/goods", api.Goods.UpdateGood).Methods(http.MethodPut)
	admrouter.HandleFunc("/goods/{id}", api.Goods.DeleteGood).Methods(http.MethodDelete)
	admrouter.HandleFunc("/categories", api.Goods.CreateCategory).Methods(http.MethodPost)
	admrouter.HandleFunc("/categories", api.Goods.UpdateCategory).Methods(http.MethodPut)
	admrouter.HandleFunc("/categories/{id}", api.Goods.DeleteCategory).Methods(http.MethodDelete)

	// shoprouter provides access to view products
	shoprouter := r.PathPrefix("/shop").Subrouter()
	shoprouter.HandleFunc("/categories", api.Goods.CategoriesList).Methods(http.MethodGet)
	shoprouter.HandleFunc("/categories/{category}/goods", api.Goods.CategoryGoods).Methods(http.MethodGet)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost%s/swagger/doc.json", api.Port)), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	api.server = &http.Server{
		Addr:    api.Port,
		Handler: setupGlobalMiddleware(r),
	}

	return api.server
}
