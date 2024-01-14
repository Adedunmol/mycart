package database

import (
	"fmt"

	"github.com/Adedunmol/mycart/internal/models"
	"gorm.io/gorm"
)

const (
	USER         uint8 = 1
	CREATE_EVENT uint8 = 2
	MODIFY_EVENT uint8 = 4
	DELETE_EVENT uint8 = 8
	MODERATE     uint8 = 16
	ADMIN        uint8 = 32
)

func InsertRoles() {
	roles := make(map[string][]uint8)

	roles["User"] = []uint8{USER}
	roles["Event-Organizer"] = []uint8{USER, CREATE_EVENT, MODIFY_EVENT, DELETE_EVENT}
	roles["Moderator"] = []uint8{USER, CREATE_EVENT, MODIFY_EVENT, DELETE_EVENT, MODERATE}
	roles["Admin"] = []uint8{USER, CREATE_EVENT, MODIFY_EVENT, DELETE_EVENT, MODERATE, ADMIN}

	default_role := "User"

	for r := range roles {
		var role models.Role
		result := Database.DB.Where(&models.Role{Name: r}).First(&role)

		if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
			role = models.Role{Name: r}

			exec_result := Database.DB.Create(&role)

			if exec_result.Error != nil {
				fmt.Println("Unable to create role")
			}
		}
		role.ResetPermissions()

		for _, perm := range roles[r] {
			role.AddPermission(uint8(perm))
		}

		role.Default = (role.Name == default_role)

		Database.DB.Save(role)
	}

	fmt.Println(roles, default_role)
}
