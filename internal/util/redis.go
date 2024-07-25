package util

import (
	"errors"
	"fmt"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/redis"
)

func WriteCartToDB(userId int) error {
	cartItems := redis.GetCart(userId)

	for _, item := range cartItems {

		if int(item.Count) < 1 {
			return errors.New("quantity can't be less than 1")
		}

		var cart models.Cart

		result := database.DB.First(&cart, userId)

		if result.Error != nil {
			cart = models.Cart{
				BuyerID: uint(userId),
			}

			database.DB.Create(&cart)
		}

		var product models.Product
		result = database.DB.First(&product, item.ItemId)

		if result.Error != nil {
			fmt.Println(result.Error)
			return errors.New("product not found")
		}

		if uint64(item.Count) > uint64(product.Quantity) {
			return errors.New("product quantity more than what's available")
		}

		cartItem := models.CartItem{
			CartID:      cart.ID,
			ProductName: product.Name,
			ProductID:   product.ID,
			UnitPrice:   uint(product.Price),
			Quantity:    uint(item.Count),
			TotalPrice:  uint(item.Count) * uint(product.Price),
		}

		result = database.DB.Create(&cartItem)

		if result.Error != nil {
			fmt.Println(result.Error)
			return errors.New("error creating cart item")
		}

		newCartPrice := cart.TotalPrice + cartItem.TotalPrice

		database.DB.Model(&cart).Updates(models.Cart{
			TotalPrice: newCartPrice,
		})

		result = database.DB.Preload("CartItems").First(&cart, cart.ID)

		if result.Error != nil {
			fmt.Println(result.Error)
			return errors.New("error updating cart")
		}

	}
	return nil
}

func UpdateCartFromDB(userId int) {
	var cart models.Cart

	result := database.DB.Preload("CartItems").First(&cart, userId)

	if result.Error != nil {
		cart = models.Cart{
			BuyerID: uint(userId),
		}

		database.DB.Create(&cart)
	}

	if len(cart.CartItems) < 1 {
		return
	}

	for _, product := range cart.CartItems {
		var foundProduct models.Product
		database.DB.First(&foundProduct, product.ID)

		redis.AddItemToCart(userId, int(foundProduct.ID), int64(product.Quantity))
	}

}

func ClearCartAndDB(userId int) {
	redis.DeleteCart(userId)

	var cart models.Cart

	result := database.DB.Preload("CartItems").First(&cart, userId)

	if result.Error != nil {
		return
	}

	if len(cart.CartItems) < 1 {
		return
	}

	for _, product := range cart.CartItems {
		var foundProduct models.Product
		database.DB.First(&foundProduct, product.ID)

		database.DB.Delete(&foundProduct)
	}

}
