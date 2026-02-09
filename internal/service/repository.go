package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nightmaker00/accordion-go/internal/domain"
)

type FAQRepository interface {
	ListActive(ctx context.Context) ([]domain.FAQ, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.FAQ, error)
	Create(ctx context.Context, in domain.CreateFAQInput) (domain.FAQ, error)
	Update(ctx context.Context, id uuid.UUID, in domain.UpdateFAQInput) (domain.FAQ, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
