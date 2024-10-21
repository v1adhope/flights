# Swagger usage

http://0.0.0.0:8081/v1/swagger/index.html

# PgAdmin 4 usage

http://0.0.0.0:8082

Username: test@mail.com
Password: admin

# Structure

```bash
.
├── bin
│   └── service_start.sh
├── cmd
│   └── main.go
├── compose.yml
├── db
│   ├── flights_schema.pdf
│   └── migrations
│       ├── 000001_init.down.sql
│       ├── 000001_init.up.sql
│       ├── 000002_init.down.sql
│       ├── 000002_init.up.sql
│       ├── 000003_init.down.sql
│       ├── 000003_init.up.sql
│       ├── 000004_init.down.sql
│       └── 000004_init.up.sql
├── dockerfile
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── configs
│   │   └── config.go
│   ├── controllers
│   │   └── http
│   │       └── v1
│   │           ├── document.go
│   │           ├── errors.go
│   │           ├── general.go
│   │           ├── interfaces.go
│   │           ├── passenger.go
│   │           ├── report.go
│   │           ├── router.go
│   │           ├── ticket.go
│   │           ├── v1_test.go
│   │           └── validations.go
│   ├── entities
│   │   ├── documents.go
│   │   ├── errors.go
│   │   ├── general.go
│   │   ├── passenger.go
│   │   ├── report.go
│   │   └── ticket.go
│   ├── testhelpers
│   │   ├── postgres.go
│   │   └── utils.go
│   └── usecases
│       ├── constructor.go
│       ├── document.go
│       ├── infrastructure
│       │   └── repository
│       │       ├── constructor.go
│       │       ├── document.go
│       │       ├── dtos.go
│       │       ├── errors.go
│       │       ├── passenger.go
│       │       ├── readers.go
│       │       ├── report.go
│       │       └── ticket.go
│       ├── interfaces.go
│       ├── passenger.go
│       ├── report.go
│       └── ticket.go
├── pkg
│   ├── httpsrv
│   │   └── httpsrv
│   │       ├── httpsrv.go
│   │       └── opts.go
│   ├── logger
│   │   ├── logger.go
│   │   └── opts.go
│   └── postgresql
│       ├── opts.go
│       └── postgresql.go
├── readme.md
├── scripts
│   └── tasks
│       ├── migrate_create.sh
│       ├── migrate_down.sh
│       ├── migrate_force.sh
│       └── migrate_up.sh
└── taskfile.yaml
```
