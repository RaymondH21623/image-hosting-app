package utils

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateID() (string, error) {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 9)
	if err != nil {
		return "", err
	}
	return id, nil
}
