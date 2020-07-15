module github.com/rpoletaev/huskyjam

go 1.13

// replace github.com/dre1080/recovr => github.com/dre1080/recover
replace gopkg.in/gavv/httpexpect.v2 => github.com/gavv/httpexpect/v2 v2.1.0

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dre1080/recovr v1.0.3
	github.com/garyburd/redigo v1.6.0
	github.com/go-mixins/log v0.2.1 // indirect
	github.com/go-openapi/spec v0.19.8 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/golang/mock v1.4.3
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/wire v0.4.0
	github.com/gorilla/mux v1.7.4
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/justinas/alice v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.0.0
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/rs/xid v1.2.1
	github.com/rs/zerolog v1.19.0
	github.com/swaggo/http-swagger v0.0.0-20200308142732-58ac5e232fba
	github.com/swaggo/swag v1.6.7
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/tools v0.0.0-20200713190748-01425d701627 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	gopkg.in/gavv/httpexpect.v2 v2.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
