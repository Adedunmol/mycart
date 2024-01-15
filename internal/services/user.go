package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserDto struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userDto CreateUserDto
	err := json.NewDecoder(r.Body).Decode(&userDto)

	if _, ok := err.(*json.InvalidUnmarshalError); ok {
		util.RespondWithJSON(w, http.StatusInternalServerError, "Unable to format the request body")
		return
	}
	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var foundUser models.User

	result := database.Database.DB.Where(models.User{Username: userDto.Username}).First(&foundUser)

	if result.Error == nil {
		util.RespondWithJSON(w, http.StatusConflict, APIResponse{Message: "username already exists", Data: nil, Status: "error"})
		return
	}

	roleToAssign := r.URL.String()

	var role models.Role

	if strings.Contains(roleToAssign, "users") {
		result = database.Database.DB.Where(models.Role{Default: true}).First(&role)
	} else {
		result = database.Database.DB.Where(models.Role{Name: "Event-Organizer"}).First(&role)
	}

	if result.Error != nil {
		fmt.Println("error looking for the role")
		util.RespondWithJSON(w, http.StatusInternalServerError, "error looking a role")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), 14)

	if err != nil {
		fmt.Println("could not hash password", err)
		util.RespondWithJSON(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	user := models.User{
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
		Email:     userDto.Email,
		Password:  string(hashedPassword),
		Username:  userDto.Username,
		RoleID:    role.ID,
	}

	result = database.Database.DB.Create(&user)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusInternalServerError, APIResponse{Message: "error creating user", Data: nil, Status: "error"})
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, APIResponse{Message: "", Data: user, Status: "success"})
	return
}
