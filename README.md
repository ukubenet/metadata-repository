# metadata-repository
Repository that allows to create, manage, and store metadata

# Running the app in 'dev' mode locally
* Navigate to the project root in Terminal and run the following to start the server:
```go
go run cmd/api/*.go
```
* Go to the browser and open Status page
```go
http://localhost:4000/status
```
# Packages
### HTTP Router
* Docs: https://pkg.go.dev/github.com/julienschmidt/httprouter#section-readme
* Install:
```shell
go get -u github.com/julienschmidt/httprouter
```
### Redis
* Docs: https://pkg.go.dev/github.com/go-redis/redis#section-readme
* Install:
```shell
go get -u github.com/go-redis/redis
```
### Uuid
* Docs: https://pkg.go.dev/github.com/google/uuid#section-readme
* Install:
```shell
go get github.com/google/uuid
```