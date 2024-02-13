package util

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"

	"fmt"
	"time"

	generator "github.com/angelodlfrtr/go-invoice-generator"
)

// func GeneratePdf(product models.product, user models.User, fileStr *string, mux *sync.Mutex, wg *sync.WaitGroup) (string, error) {
func GeneratePdf(cart models.Cart, user models.User) (string, error) {

	var foundUser models.User

	result := database.Database.DB.First(&foundUser, user.ID)

	if result.Error != nil {
		return "", errors.New("no user found with this id")
	}

	doc, _ := generator.New(generator.Invoice, &generator.Options{
		TextTypeInvoice: "FACTURE",
		AutoPrint:       true,
	})

	doc.SetHeader(&generator.HeaderFooter{
		Text:       "<center>Receipt from mycart",
		Pagination: true,
	})

	doc.SetFooter(&generator.HeaderFooter{
		Text:       "<center>Receipt from mycart",
		Pagination: true,
	})

	doc.SetRef("random")
	doc.SetVersion("1.0")

	doc.SetDescription(fmt.Sprintf("A description"))
	doc.SetDate(fmt.Sprintf("%d/%d/%d", time.Now().Day(), time.Now().Month(), time.Now().Year()))

	doc.SetCompany(&generator.Contact{
		Name: "mycart",
		Address: &generator.Address{
			Address:    "123 test str",
			Address2:   "Apartment 2",
			PostalCode: "12345",
			City:       "Test",
			Country:    "Test",
		},
	})

	doc.SetCustomer(&generator.Contact{
		Name: foundUser.Username,
		Address: &generator.Address{
			Address:    "123 test str",
			Address2:   "Apartment 2",
			PostalCode: "12345",
			City:       "Test",
			Country:    "Test",
		},
	})

	for _, cartItem := range cart.CartItems {
		doc.AppendItem(&generator.Item{
			Name:     cartItem.ProductName,
			UnitCost: fmt.Sprintf("%d", cartItem.UnitPrice),
			Quantity: fmt.Sprintf("%d", cartItem.Quantity),
			Tax: &generator.Tax{
				Amount: "0",
			},
			Discount: &generator.Discount{
				Percent: "0",
			},
		})
	}

	pdf, err := doc.Build()
	if err != nil {
		return "", err
	}

	currentDir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	currentTime := time.Now()
	userInvoice := fmt.Sprintf("%s-%s.pdf", user.Username, currentTime)
	filePath := filepath.Join(currentDir, "..", "..", "internal", "order_invoices", userInvoice)

	err = pdf.OutputFileAndClose(filePath)

	if err != nil {
		return "", err
	}

	return filePath, nil
}
