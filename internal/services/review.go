package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type CreateReviewDto struct {
	Comment string `json:"comment" validate:"required"`
	Rating  uint   `json:"rating" validate:"required"`
}

func CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "no product id sent in the url param", Data: nil, Status: "error"})
		return
	}

	var reviewDto CreateReviewDto
	err := json.NewDecoder(r.Body).Decode(&reviewDto)

	if _, ok := err.(*json.InvalidUnmarshalError); ok {
		util.RespondWithJSON(w, http.StatusInternalServerError, "Unable to format the request body")
		return
	}
	if err != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := util.Validator.Struct(reviewDto); err != nil {

		validationErrors := ValidationErrors{}

		for _, err := range err.(validator.ValidationErrors) {

			errorItem := ValidationErrorItems{Field: err.Field(), Detail: err.ActualTag()}

			validationErrors.Errors = append(validationErrors.Errors, errorItem)
		}

		util.RespondWithJSON(w, http.StatusUnprocessableEntity, APIResponse{Message: validationErrors, Data: nil, Status: "error"})
		return
	}

	username := r.Context().Value("username")

	if username == nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, "Not authorized")
		return
	}

	var foundUser models.User

	result := database.Database.DB.Where(models.User{Username: username.(string)}).First(&foundUser)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "user does not exist", Data: nil, Status: "error"})
		return
	}

	var product models.Product

	result = database.Database.DB.First(&product, id)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "", Data: nil, Status: "success"})
		return
	}

	review := models.Review{
		Comment:   reviewDto.Comment,
		Rating:    reviewDto.Rating,
		ProductID: product.ID,
	}

	result = database.Database.DB.Create(&review)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusInternalServerError, APIResponse{Message: "error creating review", Data: nil, Status: "error"})
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, APIResponse{Message: "", Data: review, Status: "success"})
}

func GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "review_id")

	if id == "" {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "no review id sent in the url param", Data: nil, Status: "error"})
		return
	}

	var review models.Review
	result := database.Database.DB.First(&review, id)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "review not found", Data: nil, Status: "success"})
		return
	}

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: review, Status: "success"})
}
