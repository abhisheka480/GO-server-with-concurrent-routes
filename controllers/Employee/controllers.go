package Employee

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GO-server-with-concurrent-routes/config"
	"github.com/GO-server-with-concurrent-routes/controllers/mongoDB"
	"github.com/GO-server-with-concurrent-routes/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddEmployeeData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	employeeData := models.Employee{}

	// we decode our body request params
	err := json.NewDecoder(r.Body).Decode(&employeeData)
	if err != nil {
		fmt.Println("Error in reading employee data")
		resData := createResponse(http.StatusBadRequest, "Error in reading employee data")
		json.NewEncoder(w).Encode(resData)
	}

	employeeData.IsActive = true
	fmt.Println(employeeData)

	if employeeData.Name == "" {
		fmt.Println("Name of employee cannot be empty")
		resData := createResponse(http.StatusBadRequest, "Name of employee cannot be empty")
		json.NewEncoder(w).Encode(resData)
		return
	}

	// connect db
	collection := mongoDB.ConnectDB()

	result, err := collection.InsertOne(context.TODO(), employeeData)
	if err != nil {
		mongoDB.GetError(err, w)
		return
	}
	config.ID_array = append(config.ID_array, result.InsertedID)

	json.NewEncoder(w).Encode(result)
	fmt.Println("POST succesfull")
}

func UpdateEmployeeDataByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	fmt.Println("params:", params)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])
	fmt.Println("id:", id)

	employeeData := &models.Employee{}
	incoingEmployeeData := &models.Employee{}

	collection := mongoDB.ConnectDB()

	// Create filter
	filter := bson.M{"_id": id}

	// get stored employee data
	err := collection.FindOne(context.TODO(), filter).Decode(&employeeData)
	if err != nil {
		mongoDB.GetError(err, w)
		return
	}

	// Read update model from body request
	err = json.NewDecoder(r.Body).Decode(&incoingEmployeeData)
	if err != nil {
		fmt.Println("Error in reading employee data")
		resData := createResponse(http.StatusBadRequest, "Error in reading employee data")
		json.NewEncoder(w).Encode(resData)
		return
	}

	if incoingEmployeeData.Name != "" {
		employeeData.Name = incoingEmployeeData.Name
	}
	if incoingEmployeeData.Address.HouseNumber != 0 {
		employeeData.Address.HouseNumber = incoingEmployeeData.Address.HouseNumber
	}
	if incoingEmployeeData.Address.City != "" {
		employeeData.Address.City = incoingEmployeeData.Address.City
	}
	if incoingEmployeeData.Address.State != "" {
		employeeData.Address.State = incoingEmployeeData.Address.State
	}
	if incoingEmployeeData.Address.Pincode != "" {
		employeeData.Address.Pincode = incoingEmployeeData.Address.Pincode
	}
	if incoingEmployeeData.Address.Street != "" {
		employeeData.Address.Street = incoingEmployeeData.Address.Street
	}
	if incoingEmployeeData.Department != "" {
		employeeData.Department = incoingEmployeeData.Department
	}
	if len(incoingEmployeeData.Skills) > 0 {
		employeeData.Skills = incoingEmployeeData.Skills
	}

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"name", employeeData.Name},
			{"department", employeeData.Department},
			{"address", bson.D{
				{"houseNumber", employeeData.Address.HouseNumber},
				{"street", employeeData.Address.Street},
				{"city", employeeData.Address.City},
				{"state", employeeData.Address.State},
				{"pincode", employeeData.Address.Pincode},
			}},
			{"skills", employeeData.Skills},
		}},
	}

	err = collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&employeeData)
	if err != nil {
		mongoDB.GetError(err, w)
		return
	}

	employeeData.ID = id

	resData := createResponse(http.StatusOK, "Successfully Updated employee data")
	json.NewEncoder(w).Encode(resData)
}

func ActivateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	fmt.Println("params:", params)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])
	fmt.Println("id:", id)

	employeeData := &models.Employee{}
	employeeData.IsActive = true

	collection := mongoDB.ConnectDB()

	// Create filter
	filter := bson.M{"_id": id}

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"isActive", employeeData.IsActive},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&employeeData)
	if err != nil {
		mongoDB.GetError(err, w)
		return
	}

	employeeData.ID = id
	resData := createResponse(http.StatusOK, "Successfully Activated employee")

	json.NewEncoder(w).Encode(resData)
}

func DeactivateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)
	fmt.Println("params:", params)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])
	fmt.Println("id:", id)
	keys, ok := r.URL.Query()["permanentlyDelete"]
	fmt.Println(keys)
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}
	queryParam := keys[0]
	//queryParam := r.URL.Query().Get("permanentlyDelete") //can be true or false
	fmt.Println("queryParam:", queryParam)

	employeeData := &models.Employee{}
	employeeData.IsActive = false

	collection := mongoDB.ConnectDB()

	// Create filter
	filter := bson.M{"_id": id}

	if queryParam == "false" { //permanentlyDelete = false
		update := bson.D{
			{"$set", bson.D{
				{"isActive", employeeData.IsActive},
			}},
		}

		err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&employeeData)
		if err != nil {
			mongoDB.GetError(err, w)
			return
		}
	} else { //permanentlyDelete = true
		deleteCount, err := collection.DeleteOne(context.TODO(), filter)
		fmt.Println("deleteCount:", deleteCount)
		if err != nil {
			mongoDB.GetError(err, w)
			return
		}
		if deleteCount.DeletedCount == 0 {
			resData := createResponse(http.StatusNotFound, "Employee data could not be deleted as it is not present in mongodb")
			json.NewEncoder(w).Encode(resData)
			return
		}
		for i := range config.ID_array {
			if config.ID_array[i] == id {
				config.ID_array = append(config.ID_array[:i], config.ID_array[i+1:]...)
			}
		}
	}
	employeeData.ID = id
	resData := createResponse(http.StatusOK, "Successfully Deactivated employee")
	json.NewEncoder(w).Encode(resData)
}

func GetAllEmployeeData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	employees := []models.Employee{}

	//Connection mongoDB with helper class
	collection := mongoDB.ConnectDB()

	// Create filter for get all employee data of active employees
	filter := bson.M{"isActive": true}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		mongoDB.GetError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		employeeData := models.Employee{}
		err := cur.Decode(&employeeData)
		if err != nil {
			fmt.Println("Error in decoding employee data")
			resData := createResponse(http.StatusBadRequest, "Error in decoding employee data")
			json.NewEncoder(w).Encode(resData)
			return
		}

		employees = append(employees, employeeData)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Error in decoding employee data")
		resData := createResponse(http.StatusBadRequest, "Error in decoding employee data")
		json.NewEncoder(w).Encode(resData)
	}

	json.NewEncoder(w).Encode(employees)
}

func GetEmployeeDataByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	employeeData := &models.Employee{}
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := mongoDB.ConnectDB()

	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&employeeData)
	if err != nil {
		mongoDB.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(employeeData)
}

func GetAllEmployeeID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(config.ID_array)
}

func createResponse(code int, message string) models.ResponseData {
	responseData := models.ResponseData{}
	responseData.StatusCode = code
	responseData.Message = message
	return responseData
}
