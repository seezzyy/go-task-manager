package service

import (
	"fmt"
	"strings"
	"task-manager/internal/model"
	"unicode/utf8"
)

func validateTitle(title string) error {
	cleanTitle := strings.TrimSpace(title)

	if cleanTitle == "" {
		return &model.ValidationError{
			Field:   "title",
			Message: "название пустое",
		}
	}

	if utf8.RuneCountInString(cleanTitle) > 100 {
		return fmt.Errorf("%w: Название слишком большое", model.ErrInvalidTask)
	}
	return nil
}

func validateDescription(description string) error {
	cleanDescription := strings.TrimSpace(description)

	if utf8.RuneCountInString(cleanDescription) > 500 {
		return fmt.Errorf("%w: Описание слишком большое", model.ErrInvalidTask)
	}

	return nil
}
