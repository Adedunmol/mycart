package services_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Adedunmol/mycart/internal/app"
	"github.com/Adedunmol/mycart/internal/config"
	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"
	"github.com/Adedunmol/mycart/internal/redis"
	"github.com/Adedunmol/mycart/internal/tasks"
	jwt "github.com/golang-jwt/jwt/v5"
)

type APIResponse struct {
	Message string `json:"message"`
	Data    struct {
		Token      string        `json:"token"`
		Expiration time.Duration `json:"expiration"`
	}
	Status string `json:"status"`
}

const redisAddress = "127.0.0.1:6379"

func TestMain(m *testing.M) {

	go tasks.Init(redisAddress)

	go tasks.Run()

	defer tasks.Close()

	go redis.Init(redisAddress)
	defer redis.Close()

	code := m.Run()

	// drop table(s) here
	clearTables()

	os.Exit(code)
}

func clearTables() {
	database.DB.Migrator().DropTable(&models.User{}, &models.Product{})
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func createUser() models.User {
	var role models.Role

	_ = database.DB.Where(models.Role{Default: true}).First(&role)

	user := models.User{
		FirstName: "test",
		LastName:  "test",
		Email:     "test@test.com",
		Password:  "$2a$14$qyjua7SGlBgCN9/sV9eowOuVVOyQe2hwpHrZ.rKMXdZxlSvi/ubXe",
		Username:  "testusername",
		RoleID:    role.ID,
	}

	_ = database.DB.Create(&user)

	return user
}

func createVendor() models.User {
	var role models.Role

	_ = database.DB.Where(models.Role{Name: "Vendor"}).First(&role)

	user := models.User{
		FirstName: "testvendor",
		LastName:  "testvendor",
		Email:     "testvendor@test.com",
		Password:  "$2a$14$qyjua7SGlBgCN9/sV9eowOuVVOyQe2hwpHrZ.rKMXdZxlSvi/ubXe",
		Username:  "testvendorusername",
		RoleID:    role.ID,
	}

	_ = database.DB.Create(&user)

	return user
}

func createProduct() (models.Product, models.User) {
	user := createVendor()

	product := models.Product{
		Name:     "test product",
		Details:  "test product details",
		Price:    10,
		Category: "test",
		Quantity: uint(10),
		Vendor:   user.ID,
	}

	database.DB.Create(&product)

	return product, user
}

func generateToken(username string, expiration time.Duration) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(expiration).Unix(), // jwt.NewNumericDate(time.Now().Add(expiration)),
		"iat":      time.Now(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.EnvConfig.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func createReview() models.Review {
	user := createUser()
	product, _ := createProduct()

	review := models.Review{
		Comment:   "some random comment",
		Rating:    4,
		ProductID: product.ID,
		UserID:    user.ID,
	}

	database.DB.Create(&review)

	return review
}

func addItemToCart(userID int, productID int) {
	product, _ := createProduct()

	if productID == 0 {
		redis.AddItemToCart(userID, int(product.ID), 1)
	} else {
		redis.AddItemToCart(userID, int(productID), 1)
	}
}
