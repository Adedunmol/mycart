package util

import (
	"errors"

	"github.com/Adedunmol/mycart/internal/database"
	"github.com/Adedunmol/mycart/internal/models"

	"fmt"
	"time"

	generator "github.com/angelodlfrtr/go-invoice-generator"
)

// func GeneratePdf(product models.product, user models.User, fileStr *string, mux *sync.Mutex, wg *sync.WaitGroup) (string, error) {
func GeneratePdf(product models.Product, user models.User, fileStr *string) (string, error) {

	var organizer models.User

	result := database.Database.DB.First(&organizer, product.Vendor)

	if result.Error != nil {
		return "", errors.New("no user found with this id")
	}

	doc, _ := generator.New(generator.Invoice, &generator.Options{
		TextTypeInvoice: "FACTURE",
		AutoPrint:       true,
	})

	doc.SetHeader(&generator.HeaderFooter{
		Text:       "<center>Receipt from Eve",
		Pagination: true,
	})

	doc.SetFooter(&generator.HeaderFooter{
		Text:       "<center>Receipt from Eve",
		Pagination: true,
	})

	doc.SetRef("random")
	doc.SetVersion("1.0")

	doc.SetDescription(fmt.Sprintf("An invoice for the purchase %s", product.Name))
	doc.SetDate(fmt.Sprintf("%d/%d/%d", time.Now().Day(), time.Now().Month(), time.Now().Year()))

	doc.SetCompany(&generator.Contact{
		Name: organizer.Username,
		Address: &generator.Address{
			Address:    "123 test str",
			Address2:   "Apartment 2",
			PostalCode: "12345",
			City:       "Test",
			Country:    "Test",
		},
	})

	doc.SetCustomer(&generator.Contact{
		Name: organizer.Username,
		Address: &generator.Address{
			Address:    "123 test str",
			Address2:   "Apartment 2",
			PostalCode: "12345",
			City:       "Test",
			Country:    "Test",
		},
	})

	doc.AppendItem(&generator.Item{
		Name:     product.Name,
		UnitCost: fmt.Sprintf("%d", product.Price),
		Quantity: fmt.Sprintf("%d", 0),
		Tax: &generator.Tax{
			Amount: "0",
		},
		Discount: &generator.Discount{
			Percent: "0",
		},
	})

	pdf, err := doc.Build()
	if err != nil {
		return "", err
	}

	*fileStr = fmt.Sprintf("%s.pdf", user.Username)
	err = pdf.OutputFileAndClose(*fileStr)

	if err != nil {
		return "", err
	}

	return *fileStr, nil
}
