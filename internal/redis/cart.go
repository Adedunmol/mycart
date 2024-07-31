package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/Adedunmol/mycart/internal/logger"

	"errors"
	"fmt"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
)

type CartItem struct {
	ItemId int
	Count  int
}

func AddItemToCart(userId int, itemId int, count int64) {
	ctx := context.Background()

	ttl, err := redisClient.TTL(ctx, "shadowKey:cart:"+strconv.Itoa(userId)).Result()

	if err != nil {
		logger.Logger.Error("error getting shadow key's TTL for cart")
		logger.Logger.Error(err.Error())
	}

	if ttl.Nanoseconds() < 0 {
		_, err = redisClient.Set(ctx, "shadowKey:cart:"+strconv.Itoa(userId), "", time.Duration(1)*time.Hour).Result()

		if err != nil {
			logger.Logger.Error("error adding shadow key to cart")
			logger.Logger.Error(err.Error())
		}
	}

	var itemCount int64

	if count == 0 {
		itemCount = 1
	} else {
		itemCount = count
	}

	_, err = redisClient.HIncrBy(ctx, "cart:"+strconv.Itoa(userId), strconv.Itoa(itemId), itemCount).Result()

	if err != nil {
		logger.Logger.Error("error adding item to cart")
		logger.Logger.Error(err.Error())
	}

	_, err = redisClient.HSet(ctx, "cart:"+strconv.Itoa(userId)+".meta", "updatedAt", time.Now().Unix()).Result()

	if err != nil {
		logger.Logger.Error("error adding cart's updatedAt to redis")
		logger.Logger.Error(err.Error())
	}
}

func GetCart(userId int) []CartItem {

	ctx := context.Background()

	cart, err := redisClient.HGetAll(ctx, "cart:"+strconv.Itoa(userId)).Result()

	if err != nil {
		logger.Logger.Error("error getting cart")
		logger.Logger.Error(err.Error())
		return nil
	}

	if len(cart) == 0 {
		logger.Logger.Info("cart is empty")
		return nil
	}

	var cartItems []CartItem

	for key, val := range cart {
		itemId, _ := strconv.Atoi(key)
		count, _ := strconv.Atoi(val)

		cartItem := CartItem{ItemId: itemId, Count: count}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems
}

func GetCartAndUpdatedAt(userId int) ([]CartItem, string) {

	ctx := context.Background()

	cart, err := redisClient.HGetAll(ctx, "cart:"+strconv.Itoa(userId)).Result()

	fmt.Println("cart: ", cart)

	if err != nil {
		logger.Logger.Error("error getting cart")
		logger.Logger.Error(err.Error())
		return nil, ""
	}

	updatedAt, err := redisClient.HGet(ctx, "cart:"+strconv.Itoa(userId)+".meta", "updatedAt").Result()

	fmt.Println("updatedAt: ", updatedAt)

	if err != nil {
		logger.Logger.Error("error getting cart's updatedAt")
		logger.Logger.Error(err.Error())
		return nil, ""
	}

	if len(cart) == 0 {
		logger.Logger.Info("cart is empty")
		return nil, ""
	}

	var cartItems []CartItem

	for key, val := range cart {
		itemId, _ := strconv.Atoi(key)
		count, _ := strconv.Atoi(val)

		cartItem := CartItem{ItemId: itemId, Count: count}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, updatedAt
}

func RemoveItemFromCart(userId int, itemId int) {
	ctx := context.Background()
	itemExists, err := redisClient.HExists(ctx, "cart:"+strconv.Itoa(userId), strconv.Itoa(itemId)).Result()

	if err != nil {
		logger.Logger.Error("error checking if item exists in cart")
		logger.Logger.Error(err.Error())
		return
	}

	if !itemExists {
		msg := fmt.Sprintf("item %d does not exist in cart", itemId)
		logger.Logger.Info(msg)
	}

	_, err = redisClient.HIncrBy(ctx, "cart:"+strconv.Itoa(userId), strconv.Itoa(itemId), -1).Result()

	if err != nil {
		msg := fmt.Sprintf("error removing item %d from cart", itemId)

		logger.Logger.Error(msg)
		logger.Logger.Error(err.Error())
	}

	itemCount, err := redisClient.HGet(ctx, "cart:"+strconv.Itoa(userId), strconv.Itoa(itemId)).Result()

	if err != nil {
		msg := fmt.Sprintf("error getting item %d from cart", itemId)

		logger.Logger.Error(msg)
		logger.Logger.Error(err.Error())
	}

	itemCountVal, _ := strconv.Atoi(itemCount)
	if itemCount != "" && itemCountVal == 0 {
		_, err := redisClient.HDel(ctx, "cart:"+strconv.Itoa(userId), strconv.Itoa(itemId)).Result()

		if err != nil {
			msg := fmt.Sprintf("error deleting item %d from cart", itemId)

			logger.Logger.Error(msg)
			logger.Logger.Error(err.Error())
		}
	}
}

func DeleteCart(userId int) {
	ctx := context.Background()
	result, err := redisClient.Exists(ctx, "cart:"+strconv.Itoa(userId)).Result()

	if err != nil {
		logger.Logger.Error("error getting cart")
		logger.Logger.Error(err.Error())
	}

	if result == 0 {
		msg := fmt.Sprintf("cart %d does not exist", userId)

		logger.Logger.Error(msg)
	}

	_, err = redisClient.Del(ctx, "cart:"+strconv.Itoa(userId)).Result()

	if err != nil {
		msg := fmt.Sprintf("error deleting cart %d", userId)

		logger.Logger.Error(msg)
		logger.Logger.Error(err.Error())
	}
}

func WriteCartToDB(userId int) error {
	cartItems := GetCart(userId)

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

	fmt.Println(cart)

	for _, product := range cart.CartItems {
		var foundProduct models.Product
		database.DB.First(&foundProduct, product.ID)

		AddItemToCart(userId, int(foundProduct.ID), int64(product.Quantity))
	}

}

func ClearCartAndDB(userId int) {
	DeleteCart(userId)

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
