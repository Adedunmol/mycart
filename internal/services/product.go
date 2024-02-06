package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
)

type CreateProductDto struct {
	Name     string `json:"name"`
	Details  string `json:"details"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Category string `json:"category"`
	// Date     time.Time `json:"date"`
}

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var productDto CreateProductDto
	err := json.NewDecoder(r.Body).Decode(&productDto)

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

	product := models.Product{
		Name:     productDto.Name,
		Details:  productDto.Details,
		Price:    productDto.Price,
		Category: productDto.Category,
		Quantity: uint(productDto.Quantity),
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

func GetProductHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	if id == "" {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "no product id sent in the url param", Data: nil, Status: "error"})
		return
	}

	var product models.Product
	result := database.Database.DB.First(&product, id)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "product not found", Data: nil, Status: "success"})
		return
	}

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: product, Status: "success"})
}

func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	database.Database.DB.Where("deleted_at is null").Find(&products)

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: products, Status: "success"})
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	if id == "" {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "no product id sent in the url param", Data: nil, Status: "error"})
		return
	}

	var product models.Product

	result := database.Database.DB.Where("deleted_at is null").First(&product, id)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "product not found", Data: nil, Status: "error"})
		return
	}

	username := r.Context().Value("username")

	if username == nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, "Not authorized")
		return
	}

	var foundUser models.User

	result = database.Database.DB.Where(models.User{Username: username.(string)}).First(&foundUser)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "user does not exist", Data: nil, Status: "error"})
		return
	}

	if product.Vendor == foundUser.ID {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "not the owner of the product", Data: nil, Status: "error"})
		return
	}

	result = database.Database.DB.Delete(&product)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "error deleting product", Data: nil, Status: "error"})
		return
	}

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: product, Status: "success"})
}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	if id == "" {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "no product id sent in the url param", Data: nil, Status: "error"})
		return
	}

	var product models.Product

	result := database.Database.DB.Where("deleted_at is null").First(&product, id)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "product not found", Data: nil, Status: "success"})
		return
	}

	var productDto CreateProductDto
	err := json.NewDecoder(r.Body).Decode(&productDto)

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

	if username == nil {
		util.RespondWithJSON(w, http.StatusUnauthorized, "Not authorized")
		return
	}

	var foundUser models.User

	result = database.Database.DB.Where(models.User{Username: username.(string)}).First(&foundUser)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "user does not exist", Data: nil, Status: "error"})
		return
	}

	if product.Vendor == foundUser.ID {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "not the owner of the product", Data: nil, Status: "error"})
		return
	}

	result = database.Database.DB.Model(&product).Updates(models.Product{
		Name:     productDto.Name,
		Details:  product.Details,
		Price:    productDto.Price,
		Category: productDto.Category,
	})

	if result.Error != nil {
		fmt.Println(err)
		util.RespondWithJSON(w, http.StatusInternalServerError, "Error updating product")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: product, Status: "success"})
}
