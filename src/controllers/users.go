package controllers

import "net/http"

//CreateUser inserts a new user to the database
func CreateUser(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Creates a new user"))
}

//FetchUsers fetch all users from the database
func FetchUsers(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Fetches all users"))
}

// FetchUser fetch an un user from the database by userID
func FetchUser(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Fetch an user"))
}

//UpdateUser updates as user from the database by userID
func UpdateUser(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Updates user"))
}

//DeleteUser deletes an usser from the database by userID
func DeleteUser(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Deletes a user"))
}