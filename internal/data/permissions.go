package data

import (
	"database/sql"
	"slices"

	"github.com/google/uuid"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	return slices.Contains(p, code)
}

type PermissionModel struct {
	DB *sql.DB
}

func (m PermissionModel) GetAllForUser(userID uuid.UUID) (Permissions, error) {

}
