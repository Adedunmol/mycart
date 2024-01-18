package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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

type UserLoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var userDto UserLoginDto

	type Response struct {
		Token      string        `json:"token"`
		Expiration time.Duration `json:"expiration"`
	}

	err := json.NewDecoder(r.Body).Decode(&userDto)

	if _, ok := err.(*json.InvalidUnmarshalError); ok {
		util.RespondWithJSON(w, http.StatusInternalServerError, "Unable to format the request body")
		return
	}

	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if userDto.Email == "" && userDto.Password == "" {
		util.RespondWithJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var foundUser models.User

	result := database.Database.DB.Where(models.User{Email: userDto.Email}).First(&foundUser)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "user does not exist", Data: nil, Status: "error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userDto.Password))

	if err != nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, APIResponse{Message: "Invalid credentials", Data: nil, Status: "error"})
		return
	}

	accessToken, err := util.GenerateToken(foundUser.Username, util.ACCESS_TOKEN_EXPIRATION)

	if err != nil {
		fmt.Println(err)
		util.RespondWithJSON(w, http.StatusInternalServerError, "Unable to generate token")
		return
	}

	refreshToken, err := util.GenerateToken(foundUser.Username, util.REFRESH_TOKEN_EXPIRATION)

	if err != nil {
		fmt.Println(err)
		util.RespondWithJSON(w, http.StatusInternalServerError, "Unable to generate token")
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    refreshToken,
		Expires:  time.Now().Add(util.REFRESH_TOKEN_EXPIRATION),
		HttpOnly: true,
	}

	result = database.Database.DB.Model(&foundUser).UpdateColumn("RefreshToken", refreshToken)

	data := Response{Token: accessToken, Expiration: time.Duration(util.ACCESS_TOKEN_EXPIRATION.Seconds())}

	http.SetCookie(w, &cookie)
	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: data, Status: "success"})
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")

	if err != nil {
		util.RespondWithJSON(w, http.StatusForbidden, "Unable to access token from cookie")
		return
	}

	var foundUser models.User

	result := database.Database.DB.Where(models.User{RefreshToken: token.Value}).First(&foundUser)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusForbidden, APIResponse{Message: "user with token does not exist", Data: nil, Status: "error"})
		return
	}

	// validate token
}
