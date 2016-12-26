package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func registerUserHandler(formatter *render.Render, database Database) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var user User
		payload, _ := ioutil.ReadAll(req.Body)
		err := json.Unmarshal(payload, &user)
		if err != nil || (user == User{}) {
			formatter.JSON(w, http.StatusBadRequest, "Failed to parse user.")
			return
		}
		err = user.Save(database)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Failed to create user.")
			return
		}
		formatter.JSON(w, http.StatusCreated, "User succesfully created.")
	}
}

func loginUserHandler(formatter *render.Render, database Database) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var user User
		payload, _ := ioutil.ReadAll(req.Body)
		err := json.Unmarshal(payload, &user)
		if err != nil || (user == User{}) {
			formatter.JSON(w, http.StatusBadRequest, "Failed to parse user.")
			return
		}
		token, err := UserLogin(user.Username, user.Password, database)

		if err != nil {
			formatter.JSON(w, http.StatusBadRequest, "Failed to login")
			return
		}
		formatter.JSON(w, http.StatusOK, token)
	}
}

func tokenValidatorHandler(formatter *render.Render, database Database) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		key := vars["key"]
		if key == "" {
			formatter.JSON(w, http.StatusNotFound, "No key sent.")
			return
		}
		validToken, err := CheckTokenKey(key, database)
		if err != nil {
			formatter.JSON(w, http.StatusNotFound, "Token key not found.")
			return
		}
		formatter.JSON(w, http.StatusOK, validToken)
	}
}
