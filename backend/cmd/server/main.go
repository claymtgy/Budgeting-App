package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/claym/budgeting-app/internal/auth"
	"github.com/claym/budgeting-app/internal/config"
	"github.com/claym/budgeting-app/internal/db"
	"github.com/claym/budgeting-app/internal/handler"
	"github.com/claym/budgeting-app/internal/middleware"
	"github.com/claym/budgeting-app/internal/repository"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := db.WaitForDatabase(context.Background(), cfg.DatabaseURL); err != nil {
		log.Fatalf("database: %v", err)
	}

	if err := db.RunMigrations(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	repo := repository.New(pool)
	tokens := auth.NewTokenService(cfg.JWTSecret)

	if cfg.CORSAllowLocal {
		log.Printf("CORS: local dev mode (localhost/127.0.0.1 on any port); defaults: %v", cfg.CORSOrigins)
	} else {
		log.Printf("CORS allowed origins: %v", cfg.CORSOrigins)
	}

	authHandler := handler.NewAuthHandler(repo, tokens)
	incomeHandler := handler.NewIncomeHandler(repo)
	incomeReceiptHandler := handler.NewIncomeReceiptHandler(repo)
	envelopeHandler := handler.NewEnvelopeHandler(repo)
	expenseHandler := handler.NewExpenseHandler(repo)
	summaryHandler := handler.NewSummaryHandler(repo)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(buildCORSOptions(cfg)))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(tokens))
			r.Get("/auth/me", authHandler.Me)

			r.Group(func(r chi.Router) {
				r.Use(middleware.Household(repo))

				r.Get("/incomes", incomeHandler.List)
			r.Post("/incomes", incomeHandler.Create)
			r.Put("/incomes/{id}", incomeHandler.Update)
			r.Delete("/incomes/{id}", incomeHandler.Delete)

			r.Get("/envelopes", envelopeHandler.List)
			r.Post("/envelopes", envelopeHandler.Create)
			r.Put("/envelopes/{id}", envelopeHandler.Update)
			r.Delete("/envelopes/{id}", envelopeHandler.Delete)

			r.Get("/income-receipts", incomeReceiptHandler.List)
			r.Post("/income-receipts", incomeReceiptHandler.Create)
			r.Put("/income-receipts/{id}/void", incomeReceiptHandler.Void)

			r.Get("/expenses", expenseHandler.List)
			r.Get("/expenses/places", expenseHandler.ListPlaces)
			r.Post("/expenses", expenseHandler.Create)
			r.Put("/expenses/{id}/void", expenseHandler.Void)

			r.Get("/summary/trends", summaryHandler.GetTrends)
			r.Get("/summary/monthly", summaryHandler.GetMonthly)
			r.Get("/summary", summaryHandler.Get)
			})
		})
	})

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	fmt.Println("server stopped")
}

func buildCORSOptions(cfg *config.Config) cors.Options {
	opts := cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}

	if cfg.CORSAllowLocal {
		allowed := make(map[string]struct{}, len(cfg.CORSOrigins))
		for _, origin := range cfg.CORSOrigins {
			allowed[origin] = struct{}{}
		}
		opts.AllowOriginFunc = func(_ *http.Request, origin string) bool {
			if origin == "" {
				return true
			}
			if _, ok := allowed[origin]; ok {
				return true
			}
			return config.IsLocalOrigin(origin)
		}
	} else {
		opts.AllowedOrigins = cfg.CORSOrigins
	}

	return opts
}
