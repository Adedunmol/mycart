package database

import (
	"fmt"

	"github.com/Adedunmol/mycart/internal/models"
	"gorm.io/gorm"
)

const (
	USER           uint8 = 1
	CREATE_PRODUCT uint8 = 2
	MODIFY_PRODUCT uint8 = 4
	DELETE_PRODUCT uint8 = 8
	MODERATE       uint8 = 16
	ADMIN          uint8 = 32
)

func InsertRoles() {
	roles := make(map[string][]uint8)

	roles["User"] = []uint8{USER}
	roles["Vendor"] = []uint8{USER, CREATE_PRODUCT, MODIFY_PRODUCT, DELETE_PRODUCT}
	roles["Moderator"] = []uint8{USER, CREATE_PRODUCT, MODIFY_PRODUCT, DELETE_PRODUCT, MODERATE}
	roles["Admin"] = []uint8{USER, CREATE_PRODUCT, MODIFY_PRODUCT, DELETE_PRODUCT, MODERATE, ADMIN}

	default_role := "User"

	for r := range roles {
		var role models.Role
		result := DB.Where(&models.Role{Name: r}).First(&role)

		if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
			role = models.Role{Name: r}

			exec_result := DB.Create(&role)

			if exec_result.Error != nil {
				fmt.Println("Unable to create role")
			}
		}
		role.ResetPermissions()

		for _, perm := range roles[r] {
			role.AddPermission(uint8(perm))
		}

		role.Default = (role.Name == default_role)

		DB.Save(role)
	}

	fmt.Println(roles, default_role)
}
