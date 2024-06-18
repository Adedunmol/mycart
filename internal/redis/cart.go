package redis

import (
	"context"
	"strconv"

	"github.com/Adedunmol/mycart/internal/logger"
)

type CartItem struct {
	ItemId int
	Count  int
}

func AddItemToCart(userId int, itemId int, count int64) {
	ctx := context.Background()
	var itemCount int64

	if count == 0 {
		itemCount = 1
	} else {
		itemCount = count
	}

	_, err := redisClient.HIncrBy(ctx, "cart:"+strconv.Itoa(userId), strconv.Itoa(itemId), itemCount).Result()

	if err != nil {
		logger.Error.Println("error adding item to cart", err)
	}
}

func GetCart(userId int) []CartItem {

	ctx := context.Background()

	cart, err := redisClient.HGetAll(ctx, "cart:"+strconv.Itoa(userId)).Result()

	if err != nil {
		logger.Error.Println("error getting cart", err)
		return nil
	}

	if len(cart) == 0 {
		logger.Info.Println("Cart is empty")
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

func RemoveItemFromCart(userId int, itemId int) {
	ctx := context.Background()
	itemExists, err := redisClient.HExists(ctx, "cart:"+strconv.Itoa(userId), strconv.Itoa(itemId)).Result()

	if err != nil {
		logger.Error.Println("error checking if item exists in cart", err)
		return
	}

	if !itemExists {
		logger.Info.Printf("item %d does not exist in cart", itemId)
	}

	_, err = redisClient.HIncrBy(ctx, "cart:"+strconv.Itoa(userId), strconv.Itoa(itemId), -1).Result()

	if err != nil {
		logger.Error.Printf("error removing item %d from cart", itemId)
	}

	itemCount, err := redisClient.HGet(ctx, "cart:"+strconv.Itoa(userId), strconv.Itoa(itemId)).Result()

	if err != nil {
		logger.Error.Printf("error getting item %d from cart", itemId)
	}

	itemCountVal, _ := strconv.Atoi(itemCount)
	if itemCount != "" && itemCountVal == 0 {
		_, err := redisClient.HDel(ctx, "cart:"+strconv.Itoa(userId), strconv.Itoa(itemId)).Result()

		if err != nil {
			logger.Error.Printf("error deleting item %d from cart", itemId)
		}
	}
}

func DeleteCart(userId int) {
	ctx := context.Background()
	_, err := redisClient.Exists(ctx, "cart:"+strconv.Itoa(userId)).Result()

	if err != nil {
		logger.Error.Println("error getting cart", err)
	}

	// if result == nil {

	// }

	_, err = redisClient.Del(ctx, "cart:"+strconv.Itoa(userId)).Result()

	if err != nil {
		logger.Error.Printf("error deleting cart %d", userId)
	}
}
