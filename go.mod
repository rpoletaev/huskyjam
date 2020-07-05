module github.com/rpoletaev/huskyjam

go 1.13

// replace github.com/dre1080/recovr => github.com/dre1080/recover

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dre1080/recovr v1.0.3
	github.com/garyburd/redigo v1.6.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/wire v0.4.0
	github.com/gorilla/mux v1.7.4
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/justinas/alice v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/rs/xid v1.2.1
	github.com/rs/zerolog v1.19.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/appengine v1.6.6 // indirect
)
