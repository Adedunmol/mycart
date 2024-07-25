package services

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/schema"
	"github.com/Adedunmol/mycart/internal/tasks"
	"github.com/Adedunmol/mycart/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type APIResponse struct {
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
}

type ValidationErrors struct {
	Errors []ValidationErrorItems `json:"errors"`
}

type ValidationErrorItems struct {
	Field  string `json:"field"`
	Detail string `json:"detail"`
}

const OTP_EXPIRATION = 30 * time.Minute

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	data, problems, err := util.DecodeJSON[*schema.CreateUser](r)

	if err != nil {

		if err == util.ErrValidation {
			util.RespondWithJSON(w, http.StatusUnprocessableEntity, util.APIResponse{Status: "error", Message: "error processing data", Data: problems})
			return
		}

		if err == util.ErrDecode {
			logger.Error.Println(err)
			util.RespondWithJSON(w, http.StatusBadRequest, util.APIResponse{Status: "error", Message: "request body needed", Data: nil})
			return
		}
	}

	var result *gorm.DB

	roleToAssign := r.URL.String()

	var role models.Role

	if strings.Contains(roleToAssign, "users") {
		result = database.DB.Where(models.Role{Default: true}).First(&role)
	} else {
		result = database.DB.Where(models.Role{Name: "Vendor"}).First(&role)
	}

	if result.Error != nil {
		logger.Error.Println("error looking for the role")
		logger.Error.Println(result.Error)
		util.RespondWithJSON(w, http.StatusInternalServerError, "error looking a role")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)

	if err != nil {
		logger.Info.Println("could not hash password")
		logger.Error.Println(err)
		util.RespondWithJSON(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	user := models.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  string(hashedPassword),
		Username:  data.Username,
		RoleID:    role.ID,
	}

	result = database.DB.Create(&user)

	if result.Error != nil {
		logger.Error.Println(result.Error)
		util.RespondWithJSON(w, http.StatusInternalServerError, APIResponse{Message: "error creating user", Data: nil, Status: "error"})
		return
	}

	verificationCode := rand.Intn(10000)
	hashedOtp, err := bcrypt.GenerateFromPassword([]byte(strconv.Itoa(verificationCode)), 14)

	if err != nil {
		logger.Error.Println("could not hash verification code")
		logger.Error.Println(err)
		util.RespondWithJSON(w, http.StatusInternalServerError, "Status internal server error")
		return
	}

	otp := models.Otp{
		User:      user,
		Code:      string(hashedOtp),
		ExpiresAt: time.Now().Add(OTP_EXPIRATION),
	}

	result = database.DB.Create(&otp)

	if result.Error != nil {
		logger.Error.Println(result.Error)
		util.RespondWithJSON(w, http.StatusInternalServerError, APIResponse{Message: "error creating otp", Data: nil, Status: "error"})
		return
	}

	emailTask, err := tasks.NewEmailDeliveryTask("verification",
		"Verify your account",
		int(user.ID),
		struct {
			Username string
			Otp      int
		}{
			Username: user.Username,
			Otp:      verificationCode,
		})

	if err != nil {
		logger.Error.Printf("Could not create task for: %d", user.ID)
		logger.Error.Println(err)
	}

	client := tasks.GetClient()

	_, err = client.Enqueue(emailTask)

	if err != nil {
		logger.Error.Printf("Could not enqueue task for: %d", user.ID)
		logger.Error.Println(err)
	}

	// util.SendMailWithTemplate("verification", user.Email, "Verify your account", )

	util.RespondWithJSON(w, http.StatusCreated, APIResponse{Message: "", Data: user, Status: "success"})
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Token      string        `json:"token"`
		Expiration time.Duration `json:"expiration"`
	}

	data, problems, err := util.DecodeJSON[*schema.CreateUser](r)

	if err != nil {

		if err == util.ErrValidation {
			util.RespondWithJSON(w, http.StatusUnprocessableEntity, util.APIResponse{Status: "error", Message: "error processing data", Data: problems})
			return
		}

		if err == util.ErrDecode {
			logger.Error.Println(err)
			util.RespondWithJSON(w, http.StatusBadRequest, util.APIResponse{Status: "error", Message: "request body needed", Data: nil})
			return
		}
	}

	var foundUser models.User

	result := database.DB.Where(models.User{Email: data.Email}).First(&foundUser)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, util.APIResponse{Message: "user does not exist", Data: nil, Status: "error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(data.Password))

	if err != nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, util.APIResponse{Message: "Invalid credentials", Data: nil, Status: "error"})
		return
	}

	accessToken, err := util.GenerateToken(foundUser.Username, util.ACCESS_TOKEN_EXPIRATION)

	if err != nil {
		logger.Error.Println(err)
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
		Name:  "token",
		Value: refreshToken,
		// Expires:  time.Now().Add(util.REFRESH_TOKEN_EXPIRATION),
		HttpOnly: true,
		MaxAge:   1 * 60 * 60,
	}

	database.DB.Model(&foundUser).UpdateColumn("RefreshToken", refreshToken)

	resData := Response{Token: accessToken, Expiration: time.Duration(util.ACCESS_TOKEN_EXPIRATION.Seconds())}

	http.SetCookie(w, &cookie)

	cartUpdateTask, err := tasks.NewCartUpdateTask(int(foundUser.ID))

	if err != nil {
		logger.Error.Printf("Could not create task for: %d", foundUser.ID)
		logger.Error.Println(err)
	}

	client := tasks.GetClient()

	_, err = client.Enqueue(cartUpdateTask)

	if err != nil {
		logger.Error.Printf("Could not enqueue task for: %d", foundUser.ID)
		logger.Error.Println(err)
	}

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: resData, Status: "success"})
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token      string        `json:"token"`
		Expiration time.Duration `json:"expiration"`
	}

	token, err := r.Cookie("token")

	if err != nil {
		util.RespondWithJSON(w, http.StatusForbidden, "Unable to access token from cookie")
		return
	}

	var foundUser models.User

	result := database.DB.Where(models.User{RefreshToken: token.Value}).First(&foundUser)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusForbidden, APIResponse{Message: "user with token does not exist", Data: nil, Status: "error"})
		return
	}

	// validate token
	username, err := util.DecodeToken(token.Value)

	if err != nil || username != foundUser.Username {
		util.RespondWithJSON(w, http.StatusUnauthorized, APIResponse{Message: "bad token", Data: nil, Status: "error"})
		return
	}

	accessToken, err := util.GenerateToken(foundUser.Username, util.ACCESS_TOKEN_EXPIRATION)

	if err != nil {
		fmt.Println(err)
		util.RespondWithJSON(w, http.StatusInternalServerError, "Unable to generate token")
		return
	}

	data := Response{Token: accessToken, Expiration: time.Duration(util.ACCESS_TOKEN_EXPIRATION.Seconds())}

	cartUpdateTask, err := tasks.NewCartUpdateTask(int(foundUser.ID))

	if err != nil {
		logger.Error.Printf("Could not create task for: %d", foundUser.ID)
		logger.Error.Println(err)
	}

	client := tasks.GetClient()

	_, err = client.Enqueue(cartUpdateTask)

	if err != nil {
		logger.Error.Printf("Could not enqueue task for: %d", foundUser.ID)
		logger.Error.Println(err)
	}

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: data, Status: "success"})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")

	if err != nil {
		util.RespondWithJSON(w, http.StatusNoContent, "")
		return
	}

	cookie := http.Cookie{
		Name:  "token",
		Value: "",
		// Expires:  time.Now().Add(util.REFRESH_TOKEN_EXPIRATION),
		HttpOnly: true,
		MaxAge:   -1,
	}

	var foundUser models.User

	result := database.DB.Where(models.User{RefreshToken: token.Value}).First(&foundUser)

	if result.Error != nil {
		http.SetCookie(w, &cookie)
		util.RespondWithJSON(w, http.StatusNoContent, "")
		return
	}

	_ = database.DB.Model(&foundUser).UpdateColumn("RefreshToken", "")

	http.SetCookie(w, &cookie)

	util.RespondWithJSON(w, http.StatusNoContent, "")
}

func VerifyUserHandler(w http.ResponseWriter, r *http.Request) {
	data, problems, err := util.DecodeJSON[*schema.VerifyUserSchema](r)

	if err != nil {

		if err == util.ErrValidation {
			util.RespondWithJSON(w, http.StatusUnprocessableEntity, util.APIResponse{Status: "error", Message: "error processing data", Data: problems})
			return
		}

		if err == util.ErrDecode {
			logger.Error.Println(err)
			util.RespondWithJSON(w, http.StatusBadRequest, util.APIResponse{Status: "error", Message: "request body needed", Data: nil})
			return
		}
	}

	var otp models.Otp

	result := database.DB.First(&otp, data.UserId)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, util.APIResponse{Message: "Account record doesn't exist or has been verified already. Please sign up or log in.", Data: nil, Status: "error"})
		return
	}

	if time.Now().Unix() > otp.ExpiresAt.Unix() {
		util.RespondWithJSON(w, http.StatusBadRequest, util.APIResponse{Message: "Code expired. Please request a new code", Data: nil, Status: "error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(otp.Code), []byte(data.Otp))

	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, util.APIResponse{Message: "Invalid otp", Data: nil, Status: "error"})
		return
	}

	var user models.User

	result = database.DB.First(&user, data.UserId)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, util.APIResponse{Message: "Account record doesn't exist. Please sign up.", Data: nil, Status: "error"})
		return
	}

	database.DB.Model(&user).UpdateColumn("Verified", true)

	util.RespondWithJSON(w, http.StatusOK, util.APIResponse{Message: "user verified successfully", Data: nil, Status: "success"})
}
