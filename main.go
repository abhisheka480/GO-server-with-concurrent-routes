package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GO-server-with-concurrent-routes/controllers/Employee"
	"github.com/GO-server-with-concurrent-routes/controllers/JwtAuthentication"
	"github.com/gorilla/mux"
)

//var ID_array []string

func main() {
	router := mux.NewRouter().StrictSlash(true)
	fmt.Println("API SERVER RUNNING ON localhost:8080")
	router.HandleFunc("/setJwtCookie", JwtAuthentication.JwtTokenSet)
	subrouter := router.PathPrefix("/").Subrouter()
	// subrouter.HandleFunc("/setJwtCookie", JwtAuthentication.JwtTokenSet).Methods("POST") //subrouter used to disable middleware when setting jwt token

	subrouter.Use(JwtAuthentication.AuthenticateUser) //middleware for jwt authentication
	subrouter.HandleFunc("/add", Employee.AddEmployeeData).Methods("POST")
	subrouter.HandleFunc("/list", Employee.GetAllEmployeeData).Methods("GET")
	subrouter.HandleFunc("/search/{id}", Employee.GetEmployeeDataByID).Methods("GET")      //get employee data by ID
	subrouter.HandleFunc("/update/{id}", Employee.UpdateEmployeeDataByID).Methods("PATCH") //update employee data
	subrouter.HandleFunc("/delete/{id}", Employee.DeactivateEmployee).Methods("PATCH")     //deactivate employee isActive = false, if permanentlyDelete parameter passed then delete
	subrouter.HandleFunc("/restore/{id}", Employee.ActivateEmployee).Methods("PATCH")      //activate employee isActive = true
	subrouter.HandleFunc("/getAllEmployeeID", Employee.GetAllEmployeeID).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
