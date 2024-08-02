package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string
	Permissions uint8
	Default     bool
	Users       []User
}

func (r *Role) AddPermission(perm uint8) {
	if !r.HasPermission(perm) {
		r.Permissions += perm
	}
}

func (r *Role) RemovePermission(perm uint8) {
	if r.HasPermission(perm) {
		r.Permissions -= perm
	}
}

func (r *Role) ResetPermissions() {
	r.Permissions = 0
}

func (r *Role) HasPermission(perm uint8) bool {
	return (r.Permissions & perm) == perm
}
