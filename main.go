package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GO-server-with-concurrent-routes/controllers/Employee"
	"github.com/gorilla/mux"
)

//var ID_array []string

func main() {
	router := mux.NewRouter().StrictSlash(true)
	fmt.Println("API SERVER RUNNING ON localhost:8080")

	router.HandleFunc("/add", Employee.AddEmployeeData).Methods("POST")
	router.HandleFunc("/list", Employee.GetAllEmployeeData).Methods("GET")
	router.HandleFunc("/search/{id}", Employee.GetEmployeeDataByID).Methods("GET")      //get employee data by ID
	router.HandleFunc("/update/{id}", Employee.UpdateEmployeeDataByID).Methods("PATCH") //update employee data
	router.HandleFunc("/delete/{id}", Employee.DeactivateEmployee).Methods("PATCH")     //deactivate employee isActive = false, if permanentlyDelete parameter passed then delete
	router.HandleFunc("/restore/{id}", Employee.ActivateEmployee).Methods("PATCH")      //activate employee isActive = true
	log.Fatal(http.ListenAndServe(":8080", router))
}
