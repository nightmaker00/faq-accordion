package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nightmaker00/accordion-go/internal/domain"
)

type FAQRepository struct {
	db *sql.DB
}

func NewFAQRepository(db *sql.DB) *FAQRepository {
	return &FAQRepository{db: db}
}

func (r *FAQRepository) ListActive(ctx context.Context) ([]domain.FAQ, error) {
	const q = `
		SELECT id, title, content, position
		FROM faqs
		WHERE is_active = true
		ORDER BY position ASC
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list active faqs: %w", err)
	}
	defer rows.Close()

	out := make([]domain.FAQ, 0)
	for rows.Next() {
		var (
			idRaw    string
			title    string
			content  string
			position int
		)
		if err := rows.Scan(&idRaw, &title, &content, &position); err != nil {
			return nil, fmt.Errorf("scan faq: %w", err)
		}
		id, err := uuid.Parse(idRaw)
		if err != nil {
			return nil, fmt.Errorf("parse faq id: %w", err)
		}
		out = append(out, domain.FAQ{
			ID:       id,
			Title:    title,
			Content:  content,
			Position: position,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate faqs: %w", err)
	}
	return out, nil
}

func (r *FAQRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.FAQ, error) {
	if err := validateFAQID(id); err != nil {
		return domain.FAQ{}, err
	}
	const q = `
		SELECT id, title, content, position, is_active, created_at, updated_at
		FROM faqs
		WHERE id = $1
	`

	var out domain.FAQ
	var idRaw string
	err := r.db.QueryRowContext(ctx, q, id.String()).
		Scan(&idRaw, &out.Title, &out.Content, &out.Position, &out.IsActive, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.FAQ{}, domain.ErrNotFound
		}
		return domain.FAQ{}, fmt.Errorf("get faq: %w", err)
	}
	parsedID, err := uuid.Parse(idRaw)
	if err != nil {
		return domain.FAQ{}, fmt.Errorf("parse faq id: %w", err)
	}
	out.ID = parsedID
	return out, nil
}

func (r *FAQRepository) Create(ctx context.Context, in domain.CreateFAQInput) (domain.FAQ, error) {
	if err := validateFAQInput(in.Title, in.Content, in.Position); err != nil {
		return domain.FAQ{}, err
	}
	const q = `
		INSERT INTO faqs (title, content, position, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, content, position, is_active, created_at, updated_at
	`

	var out domain.FAQ
	var idRaw string
	err := r.db.QueryRowContext(ctx, q, in.Title, in.Content, in.Position, in.IsActive).
		Scan(&idRaw, &out.Title, &out.Content, &out.Position, &out.IsActive, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return domain.FAQ{}, fmt.Errorf("create faq: %w", err)
	}

	parsedID, err := uuid.Parse(idRaw)
	if err != nil {
		return domain.FAQ{}, fmt.Errorf("parse faq id: %w", err)
	}
	out.ID = parsedID
	return out, nil
}

func (r *FAQRepository) Update(ctx context.Context, id uuid.UUID, in domain.UpdateFAQInput) (domain.FAQ, error) {
	if err := validateFAQID(id); err != nil {
		return domain.FAQ{}, err
	}
	if err := validateFAQInput(in.Title, in.Content, in.Position); err != nil {
		return domain.FAQ{}, err
	}
	const q = `
		UPDATE faqs
		SET title = $2, content = $3, position = $4, is_active = $5, updated_at = now()
		WHERE id = $1
		RETURNING id, title, content, position, is_active, created_at, updated_at
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.FAQ{}, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var out domain.FAQ
	var idRaw string
	err = tx.QueryRowContext(ctx, q, id.String(), in.Title, in.Content, in.Position, in.IsActive).
		Scan(&idRaw, &out.Title, &out.Content, &out.Position, &out.IsActive, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.FAQ{}, domain.ErrNotFound
		}
		return domain.FAQ{}, fmt.Errorf("update faq: %w", err)
	}

	parsedID, err := uuid.Parse(idRaw)
	if err != nil {
		return domain.FAQ{}, fmt.Errorf("parse faq id: %w", err)
	}
	out.ID = parsedID

	if err := tx.Commit(); err != nil {
		return domain.FAQ{}, fmt.Errorf("commit tx: %w", err)
	}
	return out, nil
}

func (r *FAQRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := validateFAQID(id); err != nil {
		return err
	}
	const q = `DELETE FROM faqs WHERE id = $1`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(ctx, q, id.String())
	if err != nil {
		return fmt.Errorf("delete faq: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete faq: rows affected: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	_ = affected
	return nil
}

func validateFAQID(id uuid.UUID) error {
	if id == uuid.Nil {
		return domain.ValidationError{Message: "id is required"}
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
