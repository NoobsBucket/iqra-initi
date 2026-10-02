package seo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound  = errors.New("page SEO entry not found")
	ErrDuplicate = errors.New("page path already exists")
)

type Page struct {
	ID              string          `json:"id"`
	PagePath        string          `json:"page_path"`
	PageName        string          `json:"page_name"`
	MetaTitle       string          `json:"meta_title"`
	MetaDescription string          `json:"meta_description"`
	MetaKeywords    string          `json:"meta_keywords"`
	OGTitle         string          `json:"og_title"`
	OGDescription   string          `json:"og_description"`
	OGImage         string          `json:"og_image"`
	CanonicalURL    string          `json:"canonical_url"`
	Robots          string          `json:"robots"`
	JSONLD          json.RawMessage `json:"json_ld"`
	UpdatedBy       string          `json:"updated_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type Store interface {
	GetByPath(ctx context.Context, path string) (*Page, error)
	GetAll(ctx context.Context) ([]*Page, error)
	Create(ctx context.Context, page *Page) (*Page, error)
	Update(ctx context.Context, page *Page) (*Page, error)
	Delete(ctx context.Context, id string) error
}

type store struct{ db *pgxpool.Pool }

func NewStore(db *pgxpool.Pool) Store { return &store{db: db} }

const pageColumns = `id, page_path, page_name, COALESCE(meta_title, ''), COALESCE(meta_description, ''), COALESCE(meta_keywords, ''), COALESCE(og_title, ''), COALESCE(og_description, ''), COALESCE(og_image, ''), COALESCE(canonical_url, ''), COALESCE(robots, 'index,follow'), COALESCE(json_ld, 'null'::jsonb), COALESCE(updated_by::text, ''), created_at, updated_at`

func scanPage(row interface{ Scan(...any) error }, page *Page) error {
	return row.Scan(
		&page.ID, &page.PagePath, &page.PageName, &page.MetaTitle, &page.MetaDescription,
		&page.MetaKeywords, &page.OGTitle, &page.OGDescription, &page.OGImage,
		&page.CanonicalURL, &page.Robots, &page.JSONLD, &page.UpdatedBy,
		&page.CreatedAt, &page.UpdatedAt,
	)
}

func (s *store) GetByPath(ctx context.Context, path string) (*Page, error) {
	page := &Page{}
	err := scanPage(s.db.QueryRow(ctx, `SELECT `+pageColumns+` FROM page_seo WHERE page_path = $1`, path), page)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return page, err
}

func (s *store) GetAll(ctx context.Context) ([]*Page, error) {
	rows, err := s.db.Query(ctx, `SELECT `+pageColumns+` FROM page_seo ORDER BY page_name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []*Page
	for rows.Next() {
		page := &Page{}
		if err := scanPage(rows, page); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	return pages, rows.Err()
}

func (s *store) Create(ctx context.Context, page *Page) (*Page, error) {
	err := scanPage(s.db.QueryRow(ctx, `
		INSERT INTO page_seo (
			page_path, page_name, meta_title, meta_description, meta_keywords,
			og_title, og_description, og_image, canonical_url, robots, json_ld, updated_by,
			id, created_at, updated_at
		) VALUES (
			$1, $2, NULLIF($3, ''), NULLIF($4, ''), NULLIF($5, ''),
			NULLIF($6, ''), NULLIF($7, ''), NULLIF($8, ''), NULLIF($9, ''), $10, $11, $12::uuid,
			gen_random_uuid(), NOW(), NOW()
		)
		RETURNING `+pageColumns,
		page.PagePath, page.PageName, page.MetaTitle, page.MetaDescription, page.MetaKeywords,
		page.OGTitle, page.OGDescription, page.OGImage, page.CanonicalURL, page.Robots,
		jsonValue(page.JSONLD), page.UpdatedBy,
	), page)
	return page, mapStoreError(err)
}

func (s *store) Update(ctx context.Context, page *Page) (*Page, error) {
	err := scanPage(s.db.QueryRow(ctx, `
		UPDATE page_seo SET
			page_path = $1, page_name = $2, meta_title = NULLIF($3, ''),
			meta_description = NULLIF($4, ''), meta_keywords = NULLIF($5, ''),
			og_title = NULLIF($6, ''), og_description = NULLIF($7, ''),
			og_image = NULLIF($8, ''), canonical_url = NULLIF($9, ''),
			robots = $10, json_ld = $11, updated_by = $12::uuid, updated_at = NOW()
		WHERE id = $13
		RETURNING `+pageColumns,
		page.PagePath, page.PageName, page.MetaTitle, page.MetaDescription, page.MetaKeywords,
		page.OGTitle, page.OGDescription, page.OGImage, page.CanonicalURL, page.Robots,
		jsonValue(page.JSONLD), page.UpdatedBy, page.ID,
	), page)
	return page, mapStoreError(err)
}

func (s *store) Delete(ctx context.Context, id string) error {
	result, err := s.db.Exec(ctx, `DELETE FROM page_seo WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func jsonValue(value json.RawMessage) any {
	if len(value) == 0 {
		return nil
	}
	return string(value)
}

func mapStoreError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrDuplicate
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
