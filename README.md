# Example Go API 

## Description
This is an example reference from https://github.com/bxcodec/go-clean-arch.

This project using
* Gorm https://gorm.io/ for database ORM.
* Echo Framework https://echo.labstack.com/.
* Support Database mysql, mssql, postgres.

This project has  4 Domain layer :
* Models Layer
* Repository Layer
* Usecase Layer
* Delivery Layer

### How To Run This Project
> Make Sure you install go and set up the environtment

#### Run the Testing

```bash
$ go test -v -cover -covermode=atomic ./...
```

#### Run the Applications

```bash

# Clone this repo
$ git clone https://github.com/alpakih/go-api.git

#move to project
$ cd go-api

#configure copy and rename env.json.example to env.json
$ cp env.json.example env.json

# Run the application
$ go run ./cmd/go-api/main.go

```

### Tools Used:

- All libraries listed in [`go.mod`](https://github.com/bxcodec/go-clean-arch/blob/master/go.mod)
- ["github.com/vektra/mockery".](https://github.com/vektra/mockery) To Generate Mocks for testing needs.