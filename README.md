# GO-server-with-concurrent-routes
go lang rest api server with mongodb connection and jwt authentication with concurrency.
Run the server with go run main.go
Start mongodb in cmd by -: "mongod --dbpath=/<yourPath>/data/db --port 27018" in cmd.Port and db directory can be changed as per user.
Data in mongodb shell can be checked by starting mongo db shell in cmd -: "mongo mongodb://localhost:27018"
Config.go can be used to change the db name, collection name, mongodb routes in go.
Use the postman collection, import it to hit the respective routes.
