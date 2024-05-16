package services

import (
	"fmt"
	"net/http"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/util"
)

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	cartID := r.URL.Query().Get("cart_id")

	var cart models.Cart

	if cartID == "" {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "no cart id sent in the query param", Data: nil, Status: "error"})
		return
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

	result = database.DB.Preload("CartItems").First(&cart, cartID)

	if len(cart.CartItems) < 1 {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "cart is empty", Data: nil, Status: "error"})
	}

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "no cart found with this id", Data: nil, Status: "error"})
		return
	}

	order := models.Order{
		BuyerID: uint8(foundUser.ID),
		CartID:  uint8(cart.ID),
	}

	result = database.DB.Create(&order)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, APIResponse{Message: "unable to create order", Data: nil, Status: "error"})
		return
	}

	for _, product := range cart.CartItems {
		var foundProduct models.Product
		result = database.DB.First(&foundProduct, product.ID)

		newQuantity := foundProduct.Quantity - product.Quantity // (quantity in store - quantity bought)

		result = database.DB.Model(&product).Updates(models.Product{
			Quantity: newQuantity,
		})

		if result.Error != nil {
			fmt.Println(result.Error)
			util.RespondWithJSON(w, http.StatusInternalServerError, "Error updating product")
			return
		}
	}

	// generate receipt
	_, err := util.GeneratePdf(cart, foundUser)

	if err != nil {
		fmt.Println(err)
		util.RespondWithJSON(w, http.StatusInternalServerError, APIResponse{Message: "unable to generate invoice", Data: nil, Status: "error"})
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, APIResponse{Message: "", Data: order, Status: "success"})
}
