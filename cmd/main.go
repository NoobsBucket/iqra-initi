package main

import (
	"log"
	"os"
	"time"

	"github.com/NoobsBucket/iqra-initi/internal/auth"
	"github.com/NoobsBucket/iqra-initi/internal/db"
	"github.com/NoobsBucket/iqra-initi/internal/mailer"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("failed to load .env: %v", err)
	}
	cfg := config{
		addr: ":4000",
		env:  os.Getenv("IQRA_INITI_ENV"),
		db: dbConfig{
			dsn:          os.Getenv("IQRA_INITI_DB_DSN"),
			maxOpenConns: 25,
			maxIdleConns: 25,
			maxIdleTime:  15 * time.Minute,
		},
		jwt: jwtConfig{
			secret: os.Getenv("IQRA_INITI_JWT_SECRET"),
			exp:    24 * time.Hour,
		},

		mailer: mailerConfig{
			apiKey: os.Getenv("IQRA_INITI_RESEND_API_KEY"),
			from:   os.Getenv("IQRA_INITI_RESEND_FROM"),
		},
	}
	dbPool, err := db.NewPostgresDB(
		cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	log.Print("database connected")
	m := mailer.New(cfg.mailer.apiKey, cfg.mailer.from)
	authStore := auth.NewStore(dbPool)
	app := application{
		config:    cfg,
		db:        dbPool,
		mailer:    m,
		authStore: authStore,
	}
	if err := app.run(app.mount()); err != nil {
		log.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
