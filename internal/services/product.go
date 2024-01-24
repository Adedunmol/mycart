package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/util"
)

type CreateProductDto struct {
	Name     string `json:"name"`
	Details  string `json:"details"`
	Price    int    `json:"price"`
	Category string `json:"category"`
	// Date     time.Time `json:"date"`
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var eventDto CreateProductDto
	err := json.NewDecoder(r.Body).Decode(&eventDto)

	if _, ok := err.(*json.InvalidUnmarshalError); ok {
		fmt.Println(err)
		util.RespondWithJSON(w, http.StatusInternalServerError, "Unable to format the request body")
		return
	}

	if err != nil {
		fmt.Println(err)
		util.RespondWithJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	username := r.Context().Value("username")

	var foundUser models.User

	result := database.Database.DB.Where(models.User{Username: username.(string)}).First(&foundUser)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "user does not exist", Data: nil, Status: "error"})
		return
	}

	product := models.Product{
		Name:     eventDto.Name,
		Details:  eventDto.Details,
		Price:    eventDto.Price,
		Category: eventDto.Category,
		Vendor:   foundUser.ID,
	}

	result = database.Database.DB.Create(&product)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusInternalServerError, APIResponse{Message: "error creating product", Data: nil, Status: "error"})
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, APIResponse{Message: "", Data: product, Status: "success"})
}
