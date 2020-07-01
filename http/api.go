package http

import (
	"net/http"

	"github.com/gorilla/mux"
	// auth "github.com/rpoletaev/huskyjam/pkg/auth/http"
)

type Config struct {
	Port string `envconfig:"PORT"`
}

type Api struct {
	accounts accountHandler
	goods    goodsHandler
}

func (api *Api) Connect(c *Config) *http.Server {
	r := mux.NewRouter()

	// accrouter handles auth endpoints
	accrouter := r.PathPrefix("/auth").Subrouter()
	accrouter.HandleFunc("/signup", api.accounts.Signup).Methods(http.MethodPost)
	accrouter.HandleFunc("/signin", api.accounts.Signin).Methods(http.MethodPost)
	accrouter.HandleFunc("/refresh", api.accounts.Refresh).Methods(http.MethodPost)

	// admrouter handles endpoints for a goods management
	admrouter := r.PathPrefix("/admin").Subrouter()
	admrouter.HandleFunc("/goods", api.goods.CreateGood).Methods(http.MethodPost)
	admrouter.HandleFunc("/goods", api.goods.UpdateGood).Methods(http.MethodPut)
	admrouter.HandleFunc("/goods/{id}", api.goods.DeleteGood).Methods(http.MethodDelete)
	admrouter.HandleFunc("/categories", api.goods.CreateCategory).Methods(http.MethodPost)
	admrouter.HandleFunc("/categories", api.goods.UpdateCategory).Methods(http.MethodPut)
	admrouter.HandleFunc("/categories/{id}", api.goods.DeleteCategory).Methods(http.MethodDelete)
	// admrouter = auth.WithAuth(admrouter)

	// shoprouter provides access to view products
	shoprouter := r.PathPrefix("/shop").Subrouter()
	shoprouter.HandleFunc("/categories", api.goods.CategoriesList).Methods(http.MethodGet)
	shoprouter.HandleFunc("/categories/{category}/goods", api.goods.CategoryGoods).Methods(http.MethodGet)

	return &http.Server{
		Addr:    c.Port,
		Handler: setupGlobalMiddleware(r),
	}
}
