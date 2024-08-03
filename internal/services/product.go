package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/logger"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/schema"
	"github.com/Adedunmol/mycart/internal/util"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm/clause"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	data, problems, err := util.DecodeJSON[*schema.CreateProductDto](r)

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

	product := models.Product{
		Name:     data.Name,
		Details:  data.Details,
		Price:    data.Price,
		Category: data.Category,
		Quantity: uint(data.Quantity),
		Vendor:   foundUser.ID,
	}

	result = database.DB.Create(&product)

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
	result := database.DB.First(&product, id)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "product not found", Data: nil, Status: "success"})
		return
	}

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: product, Status: "success"})
}

func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	clauses := make([]clause.Expression, 0)

	// filters
	category := r.URL.Query().Get("category")
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")
	minRating := r.URL.Query().Get("min_rating")
	maxRating := r.URL.Query().Get("max_rating")

	// sorting
	sortBy := r.URL.Query().Get("sort_by")

	// pagination
	// page := r.URL.Query().Get("page")
	// pageSize := r.URL.Query().Get("page_size")

	if category != "" {
		clauses = append(clauses, clause.Eq{Column: "category", Value: category})
	}

	if minPrice != "" {
		clauses = append(clauses, clause.Gte{Column: "price", Value: minPrice})
	}

	if maxPrice != "" {
		clauses = append(clauses, clause.Lte{Column: "price", Value: maxPrice})
	}

	if minRating != "" {
		clauses = append(clauses, clause.Gte{Column: "rating", Value: minRating})
	}

	if maxRating != "" {
		clauses = append(clauses, clause.Lte{Column: "rating", Value: maxRating})
	}

	if sortBy != "" {
		condition := strings.Split(sortBy, "-")

		var orderDesc bool

		if strings.ToLower(condition[1]) == "desc" {
			orderDesc = true
		}

		switch strings.ToLower(condition[0]) {

		case "price":
			clauses = append(clauses, clause.OrderBy{Columns: []clause.OrderByColumn{{Column: clause.Column{Name: condition[0]}, Desc: orderDesc, Reorder: false}}})

		default:
		}
	}

	// newPage, err := strconv.ParseUint(page, 10, 8)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// newPageSize, _ := strconv.ParseUint(pageSize, 10, 8)

	// offset := (newPage - 1) * newPageSize

	// intPageSize := int(newPageSize)

	// if pageSize != "" {
	// 	clauses = append(clauses, clause.Limit{Limit: &intPageSize, Offset: int(offset)})
	// }

	fmt.Println(clauses)

	database.DB = database.DB.Where("deleted_at is null").Clauses(clauses...) // .Find(&products)
	database.DB = database.DB.Scopes(util.Paginate(r))

	database.DB.Find(&products)

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: products, Status: "success"})
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	if id == "" {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "no product id sent in the url param", Data: nil, Status: "error"})
		return
	}

	var product models.Product

	result := database.DB.Where("deleted_at is null").First(&product, id)

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

	result = database.DB.Where(models.User{Username: username.(string)}).First(&foundUser)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "user does not exist", Data: nil, Status: "error"})
		return
	}

	if product.Vendor != foundUser.ID {
		util.RespondWithJSON(w, http.StatusForbidden, APIResponse{Message: "not the owner of the product", Data: nil, Status: "error"})
		return
	}

	result = database.DB.Delete(&product)

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

	result := database.DB.Where("deleted_at is null").First(&product, id)

	if result.Error != nil {
		fmt.Println(result.Error)
		util.RespondWithJSON(w, http.StatusNotFound, APIResponse{Message: "product not found", Data: nil, Status: "success"})
		return
	}

	data, problems, err := util.DecodeJSON[*schema.UpdateProductDto](r)

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

	result = database.DB.Where(models.User{Username: username.(string)}).First(&foundUser)

	if result.Error != nil {
		util.RespondWithJSON(w, http.StatusBadRequest, APIResponse{Message: "user does not exist", Data: nil, Status: "error"})
		return
	}

	if product.Vendor == foundUser.ID {
		util.RespondWithJSON(w, http.StatusForbidden, APIResponse{Message: "not the owner of the product", Data: nil, Status: "error"})
		return
	}

	result = database.DB.Model(&product).Updates(models.Product{
		Name:     data.Name,
		Details:  data.Details,
		Price:    data.Price,
		Category: data.Category,
	})

	if result.Error != nil {
		fmt.Println(err)
		util.RespondWithJSON(w, http.StatusInternalServerError, "Error updating product")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, APIResponse{Message: "", Data: product, Status: "success"})
}
