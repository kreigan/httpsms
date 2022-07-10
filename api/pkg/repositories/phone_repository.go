package repositories

import (
	"context"

	"github.com/NdoleStudio/http-sms-manager/pkg/entities"
)

// PhoneRepository loads and persists an entities.Phone
type PhoneRepository interface {
	// Save Upsert a new entities.Phone
	Save(ctx context.Context, phone *entities.Phone) error

	// Index entities.Phone of a user
	Index(ctx context.Context, userID entities.UserID, params IndexParams) (*[]entities.Phone, error)

	// Load a phone by user and phone number
	Load(ctx context.Context, userID entities.UserID, phoneNumber string) (*entities.Phone, error)
}
