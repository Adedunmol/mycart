package services

import (
	"fmt"
	"net/http"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/redis"
	"github.com/Adedunmol/mycart/internal/tasks"
	"github.com/Adedunmol/mycart/internal/util"
)

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
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

	updatedCart := redis.GetCart(int(foundUser.ID))
	err := redis.WriteCartToDB(int(foundUser.ID))

	if err != nil {
		logger.Logger.Error(err.Error())
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "error updating cart", Data: nil, Status: "error"})
		return
	}

	// cartID := r.URL.Query().Get("cart_id")

	// var cart models.Cart

	// if cartID == "" {
	// 	util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "no cart id sent in the query param", Data: nil, Status: "error"})
	// 	return
	// }

	// result = database.DB.Preload("CartItems").First(&cart, cartID)

	// if result.Error != nil {
	// 	util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "no cart found with this id", Data: nil, Status: "error"})
	// 	return
	// }

	if len(updatedCart) < 1 {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "cart is empty", Data: nil, Status: "error"})
	}

	order := models.Order{
		BuyerID: uint8(foundUser.ID),
		CartID:  uint8(foundUser.ID),
	}

	result = database.DB.Create(&order)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusInternalServerError, APIResponse{Message: "unable to create order", Data: nil, Status: "error"})
		return
	}

	for _, product := range updatedCart {
		var foundProduct models.Product
		database.DB.First(&foundProduct, product.ItemId)

		if product.Count > int(foundProduct.Quantity) {
			util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "product quantity more than what's available", Data: struct {
				Available int
				Order     int
			}{Available: int(foundProduct.Quantity), Order: product.Count}, Status: "error"})
			return
		}

		newQuantity := foundProduct.Quantity - uint(product.Count) // (quantity in store - quantity bought)

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
	// _, err := util.GeneratePdf(cart, foundUser)
	invoiceTask, err := tasks.NewInvoiceGenerationTask(int(foundUser.ID), int(foundUser.ID), int(foundUser.ID), struct{}{})

	if err != nil {
		msg := fmt.Sprintf("could not create task for: %d", foundUser.ID)

		logger.Logger.Error(msg)
		logger.Logger.Error(err.Error())
	}

	client := tasks.GetClient()

	_, err = client.Enqueue(invoiceTask)

	if err != nil {
		msg := fmt.Sprintf("could not enqueue task for: %d", foundUser.ID)

		logger.Logger.Error(msg)
		logger.Logger.Error(err.Error())
	}

	if err != nil {
		fmt.Println(err)
		util.RespondWithJSON(w, http.StatusInternalServerError, APIResponse{Message: "unable to generate invoice", Data: nil, Status: "error"})
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, APIResponse{Message: "", Data: order, Status: "success"})
}
