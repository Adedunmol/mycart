package services

import (
	"fmt"
	"net/http"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/schema"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
)

func CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "no product id sent in the url param", Data: nil, Status: "error"})
		return
	}

	data, problems, err := util.DecodeJSON[*schema.CreateReviewDto](r)

	if err != nil {
		if err == util.ErrValidation {
			util.RespondWithJSON(w, http.StatusUnprocessableEntity, util.APIResponse{Status: "error", Message: "error processing data", Data: problems})
			return
		}

		if err == util.ErrDecode {
			logger.Logger.Error(err.Error())
			util.RespondWithJSON(w, http.StatusBadRequest, util.APIResponse{Status: "error", Message: "request body needed", Data: nil})
			return
		}
	}

	username := r.Context().Value("username")

	if username == nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, "Not authorized")
		return
	}

	var foundUser models.User

	result := database.DB.Where(models.User{Username: username.(string)}).First(&foundUser)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "user does not exist", Data: nil, Status: "error"})
		return
	}

	var product models.Product

	result = database.DB.First(&product, id)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "", Data: nil, Status: "success"})
		return
	}

	review := models.Review{
		Comment:   data.Comment,
		Rating:    data.Rating,
		ProductID: product.ID,
	}

	result = database.DB.Create(&review)

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
	result := database.DB.First(&review, id)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "review not found", Data: nil, Status: "success"})
		return
	}

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: review, Status: "success"})
}
