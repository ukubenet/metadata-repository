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

# Docker instructions
## Prepare environment
Create a folder for the project and switch to it.
Clone this repository under the `engine` name.
Clone the `app_without_coding_configuration` under the `data` name.
## Run Docker image and application
```
docker run --rm -it -v /path/to/ptoject:/go -w /go/engine/cmd/api -p 4000:4000 golang:latest /bin/bash
```
Start application inside the container:
```
go run *.go
```
