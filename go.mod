module github.com/SayIfOrg/say_keeper

go 1.18

require (
	github.com/99designs/gqlgen v0.17.30
	github.com/SayIfOrg/say_protos/packages/go v0.0.0-20230629164057-56a4174e39bf
	github.com/gorilla/websocket v1.5.0
	github.com/graph-gophers/dataloader v5.0.0+incompatible
	github.com/redis/go-redis/v9 v9.0.3
	github.com/vektah/gqlparser/v2 v2.5.1
	golang.org/x/exp v0.0.0-20230425010034-47ecfdc1ba53
	google.golang.org/grpc v1.51.0
	gopkg.in/guregu/null.v4 v4.0.0
	gorm.io/driver/postgres v1.5.0
	gorm.io/gorm v1.25.0
)

//replace github.com/SayIfOrg/say_protos/packages/go => ../say_protos/packages/go

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.25.1 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
