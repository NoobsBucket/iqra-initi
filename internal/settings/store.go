package settings

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Settings struct {
	ID                 string `json:"id"`
	SiteName           string `json:"site_name"`
	SiteLogo           string `json:"site_logo"`
	SiteFavicon        string `json:"site_favicon"`
	SiteEmail          string `json:"site_email"`
	SitePhone          string `json:"site_phone"`
	SiteAddress        string `json:"site_address"`
	NavbarAnnouncement string `json:"navbar_announcement"`
	FooterTagline      string `json:"footer_tagline"`
	FooterCopyright    string `json:"footer_copyright"`
	FacebookURL        string `json:"facebook_url"`
	InstagramURL       string `json:"instagram_url"`
	TwitterURL         string `json:"twitter_url"`
	YoutubeURL         string `json:"youtube_url"`
	TiktokURL          string `json:"tiktok_url"`
	WhatsappNumber     string `json:"whatsapp_number"`
	LinkedinURL        string `json:"linkedin_url"`
	TelegramURL        string `json:"telegram_url"`
	DefaultMetaTitle   string `json:"default_meta_title"`
	DefaultMetaDesc    string `json:"default_meta_description"`
}

type Store interface {
	Get(ctx context.Context) (*Settings, error)
	Update(ctx context.Context, s *Settings) (*Settings, error)
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{db: db}
}

func (s *store) Get(ctx context.Context) (*Settings, error) {
	settings := &Settings{}
	err := s.db.QueryRow(ctx, `
		SELECT id, site_name, site_logo, site_favicon, site_email, site_phone, site_address,
		       navbar_announcement, footer_tagline, footer_copyright,
		       facebook_url, instagram_url, twitter_url, youtube_url, tiktok_url,
		       whatsapp_number, linkedin_url, telegram_url,
		       default_meta_title, default_meta_description
		FROM site_settings LIMIT 1
	`).Scan(
		&settings.ID, &settings.SiteName, &settings.SiteLogo, &settings.SiteFavicon,
		&settings.SiteEmail, &settings.SitePhone, &settings.SiteAddress,
		&settings.NavbarAnnouncement, &settings.FooterTagline, &settings.FooterCopyright,
		&settings.FacebookURL, &settings.InstagramURL, &settings.TwitterURL,
		&settings.YoutubeURL, &settings.TiktokURL, &settings.WhatsappNumber,
		&settings.LinkedinURL, &settings.TelegramURL,
		&settings.DefaultMetaTitle, &settings.DefaultMetaDesc,
	)
	return settings, err
}

func (s *store) Update(ctx context.Context, settings *Settings) (*Settings, error) {
	err := s.db.QueryRow(ctx, `
		UPDATE site_settings SET
		    site_name = $1, site_logo = $2, site_favicon = $3, site_email = $4,
		    site_phone = $5, site_address = $6, navbar_announcement = $7,
		    footer_tagline = $8, footer_copyright = $9, facebook_url = $10,
		    instagram_url = $11, twitter_url = $12, youtube_url = $13,
		    tiktok_url = $14, whatsapp_number = $15, linkedin_url = $16,
		    telegram_url = $17, default_meta_title = $18, default_meta_description = $19,
		    updated_at = NOW()
		RETURNING id, site_name, site_logo, site_favicon, site_email, site_phone, site_address,
		          navbar_announcement, footer_tagline, footer_copyright,
		          facebook_url, instagram_url, twitter_url, youtube_url, tiktok_url,
		          whatsapp_number, linkedin_url, telegram_url,
		          default_meta_title, default_meta_description
	`,
		settings.SiteName, settings.SiteLogo, settings.SiteFavicon, settings.SiteEmail,
		settings.SitePhone, settings.SiteAddress, settings.NavbarAnnouncement,
		settings.FooterTagline, settings.FooterCopyright, settings.FacebookURL,
		settings.InstagramURL, settings.TwitterURL, settings.YoutubeURL,
		settings.TiktokURL, settings.WhatsappNumber, settings.LinkedinURL,
		settings.TelegramURL, settings.DefaultMetaTitle, settings.DefaultMetaDesc,
	).Scan(
		&settings.ID, &settings.SiteName, &settings.SiteLogo, &settings.SiteFavicon,
		&settings.SiteEmail, &settings.SitePhone, &settings.SiteAddress,
		&settings.NavbarAnnouncement, &settings.FooterTagline, &settings.FooterCopyright,
		&settings.FacebookURL, &settings.InstagramURL, &settings.TwitterURL,
		&settings.YoutubeURL, &settings.TiktokURL, &settings.WhatsappNumber,
		&settings.LinkedinURL, &settings.TelegramURL,
		&settings.DefaultMetaTitle, &settings.DefaultMetaDesc,
	)
	return settings, err
}
