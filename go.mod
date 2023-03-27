module github.com/dmitrymomot/solana-wallets

go 1.20

require (
	filippo.io/edwards25519 v1.0.0
	github.com/dmitrymomot/go-env v1.0.2
	github.com/dmitrymomot/oauth2-server v0.1.1-rc
	github.com/dmitrymomot/random v1.0.6
	github.com/dmitrymomot/solana v0.1.2-alpha
	github.com/fatih/color v1.15.0
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-chi/cors v1.2.1
	github.com/go-kit/kit v0.12.0
	github.com/go-redis/cache/v8 v8.4.4
	github.com/google/go-querystring v1.1.0
	github.com/gookit/validate v1.4.6
	github.com/joho/godotenv v1.5.1
	github.com/labstack/gommon v0.4.0
	github.com/lib/pq v1.10.7
	github.com/magefile/mage v1.14.0
	github.com/mcnijman/go-emailaddress v1.1.0
	github.com/mr-tron/base58 v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/portto/solana-go-sdk v1.23.1
	github.com/rubenv/sql-migrate v1.4.0
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.2
	github.com/tyler-smith/go-bip39 v1.1.0
	golang.org/x/net v0.8.0
	golang.org/x/sync v0.1.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-gorp/gorp/v3 v3.1.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-oauth2/oauth2/v4 v4.5.2 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang-jwt/jwt/v4 v4.0.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.0.0-rc.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gookit/filter v1.1.4 // indirect
	github.com/gookit/goutil v0.6.0 // indirect
	github.com/klauspost/compress v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/near/borsh-go v0.3.2-0.20220516180422-1ff87d108454 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/vmihailenco/go-tinylfu v0.2.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.4 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/oauth2 v0.6.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210917145530-b395a37504d4 // indirect
	google.golang.org/grpc v1.40.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/portto/solana-go-sdk => github.com/dmitrymomot/solana-go-sdk v1.23.6
