package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/nightmaker00/accordion-go/internal/domain"
)

type FAQService struct {
	repo FAQRepository
}

func NewFAQService(repo FAQRepository) *FAQService {
	return &FAQService{repo: repo}
}

func (s *FAQService) ListActive(ctx context.Context) ([]domain.FAQ, error) {
	items, err := s.repo.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *FAQService) GetByID(ctx context.Context, id uuid.UUID) (domain.FAQ, error) {
	if id == uuid.Nil {
		return domain.FAQ{}, domain.ValidationError{Message: "id is required"}
	}
	out, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.FAQ{}, err
	}
	return out, nil
}

func (s *FAQService) Create(ctx context.Context, in domain.CreateFAQInput) (domain.FAQ, error) {
	if err := validateFAQInput(in.Title, in.Content, in.Position); err != nil {
		return domain.FAQ{}, err
	}
	out, err := s.repo.Create(ctx, in)
	if err != nil {
		return domain.FAQ{}, err
	}
	return out, nil
}

func (s *FAQService) Update(ctx context.Context, id uuid.UUID, in domain.UpdateFAQInput) (domain.FAQ, error) {
	if id == uuid.Nil {
		return domain.FAQ{}, domain.ValidationError{Message: "id is required"}
	}
	if err := validateFAQInput(in.Title, in.Content, in.Position); err != nil {
		return domain.FAQ{}, err
	}
	out, err := s.repo.Update(ctx, id, in)
	if err != nil {
		return domain.FAQ{}, err
	}
	return out, nil
}

func (s *FAQService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return domain.ValidationError{Message: "id is required"}
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func validateFAQInput(title, content string, position int) error {
	if strings.TrimSpace(title) == "" {
		return domain.ValidationError{Message: "title is required"}
	}
	if strings.TrimSpace(content) == "" {
		return domain.ValidationError{Message: "content is required"}
	}
	if position <= 0 {
		return domain.ValidationError{Message: "position must be greater than 0"}
	}
	return nil
}
