package services

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/util"
)

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product_id")
	quantity := r.URL.Query().Get("quantity")

	var cart models.Cart

	if productID == "" {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "no product id sent in the query param", Data: nil, Status: "error"})
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

	result = database.Database.DB.First(&cart, foundUser.ID)

	if result.Error != nil {
		cart = models.Cart{
			BuyerID: foundUser.ID,
		}

		result = database.Database.DB.Create(&cart)
	}

	var product models.Product
	result = database.Database.DB.First(&product, productID)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "product not found", Data: nil, Status: "error"})
		return
	}

	newQuantity, err := strconv.ParseUint(quantity, 10, 8)

	if err != nil {
		log.Fatal(err)
	}

	if newQuantity > uint64(product.Quantity) {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "product quantity more than what's available", Data: nil, Status: "error"})
		return
	}

	cartItem := models.CartItem{
		CartID:      cart.ID,
		ProductName: product.Name,
		ProductID:   product.ID,
		Quantity:    uint(newQuantity),
		TotalPrice:  uint(newQuantity) * uint(product.Price),
	}

	result = database.Database.DB.Create(&cartItem)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusInternalServerError, APIResponse{Message: "error creating cart item", Data: nil, Status: "error"})
		return
	}

	newCartPrice := cart.TotalPrice + cartItem.TotalPrice

	result = database.Database.DB.Model(&cart).Updates(models.Cart{
		TotalPrice: newCartPrice,
	})

	result = database.Database.DB.Preload("CartItems").First(&cart, cart.ID)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusInternalServerError, "Error updating cart")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: cart, Status: "success"})
}

func RemoveFromCartHandler(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product_id")

	var cart models.Cart

	if productID == "" {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "no user id sent in the query param", Data: nil, Status: "error"})
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

	result = database.Database.DB.First(&cart, foundUser.ID)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "user does not have a cart", Data: nil, Status: "error"})
		return
	}

	newProductID, err := strconv.ParseUint(productID, 10, 8)

	if err != nil {
		log.Fatal(err)
	}

	for _, cartItem := range cart.CartItems {
		if cartItem.ProductID == uint(newProductID) {

			result = database.Database.DB.Delete(&cartItem)

			if result.Error != nil {
				fmt.Println(result.Error)
				util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "error deleting product", Data: nil, Status: "error"})
				return
			}

			util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: cartItem, Status: "success"})
		}
	}
}
