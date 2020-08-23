package config

const MONGO_URL = "mongodb://localhost:27018/"

const DB_NAME = "EMPLOYEEDB"
const COLLECTION_NAME = "EMP_DATA"

var ID_array []interface{}

var Users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var MYSECRETKEY = []byte("my_secret_key")
