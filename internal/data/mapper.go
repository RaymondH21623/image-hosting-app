package data

import "shareapp/internal/domain"

func mapUserDomainToDB(u *domain.User) *User {
	return &User{
		ID:           u.ID,
		PublicID:     u.PublicID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash(),
		Activated:    u.Activated,
		Version:      u.Version,
	}
}

func MapUserDBToDomain(u *User) *domain.User {
	return domain.NewUserFromDB(u.ID, u.PublicID, u.Username, u.Email, u.PasswordHash, u.Activated, u.Version)
}
