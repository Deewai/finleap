# Setting up
Application has been setup to use docker automatically.

Only required thing is to have docker and docker-compose installed.

Then the application can be run with
```
chmod +x ./start.sh
./start.sh
```
The command above builds the application and database image, runs the tests and exposes the api on port :3000 for usage

Application can also be run manually by doing the following.

- Start a mysql server instance
- set environment variables as follows MYSQL_HOST, MYSQL_PORT, MYSQL_USER, MYSQL_PASSWORD, MYSQL_DATABASE
- Run the tests in the application by running
```
go test ./app/
```
- Build application image 
```
go build -o main .
```
- Run the app
```
./main
```

NOTE: Application receives payload of application/json format for POST and PATCH requests