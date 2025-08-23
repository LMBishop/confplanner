package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/LMBishop/confplanner/api"
	"github.com/LMBishop/confplanner/internal/config"
	"github.com/LMBishop/confplanner/pkg/auth"
	"github.com/LMBishop/confplanner/pkg/calendar"
	"github.com/LMBishop/confplanner/pkg/conference"
	"github.com/LMBishop/confplanner/pkg/database"
	"github.com/LMBishop/confplanner/pkg/favourites"
	"github.com/LMBishop/confplanner/pkg/ical"
	"github.com/LMBishop/confplanner/pkg/session"
	"github.com/LMBishop/confplanner/pkg/user"
	"github.com/LMBishop/confplanner/web"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Unhandled error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	c := &config.Config{}
	err := config.ReadConfig("config.yaml", c)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	pool, err := database.Connect(c.Database.ConnString)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := database.Migrate(pool); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	userService := user.NewService(pool, c.AcceptRegistrations)
	favouritesService := favourites.NewService(pool)
	conferenceService, err := conference.NewService(pool)
	if err != nil {
		return fmt.Errorf("failed to create schedule service: %w", err)
	}
	calendarService := calendar.NewService(pool)
	icalService := ical.NewService(favouritesService, conferenceService)
	sessionService := session.NewMemoryStore()
	authService := auth.NewService()

	if c.Auth.EnableBasicAuth {
		authService.RegisterAuthProvider("basic", auth.NewBasicAuthProvider(userService))
	}
	for _, authProvider := range c.Auth.AuthProviders {
		provider, err := auth.NewOIDCAuthProvider(
			userService,
			authProvider.Name,
			authProvider.ClientID,
			authProvider.ClientSecret,
			authProvider.Endpoint,
			fmt.Sprintf("%s/login/%s", c.BaseURL, authProvider.Identifier),
			authProvider.LoginFilter,
			authProvider.UserSyncFilter,
			authProvider.LoginFilterAllowedValues,
		)
		if err != nil {
			return fmt.Errorf("failed to create OIDC auth provider: %w", err)
		}

		err = authService.RegisterAuthProvider(authProvider.Identifier, provider)
		if err != nil {
			return fmt.Errorf("failed to register OIDC auth provider: %w", err)
		}
	}

	mux := http.NewServeMux()
	api := api.NewServer(api.ApiServices{
		UserService:       userService,
		FavouritesService: favouritesService,
		ConferenceService: conferenceService,
		CalendarService:   calendarService,
		IcalService:       icalService,
		SessionService:    sessionService,
		AuthService:       authService,
	}, c.BaseURL)
	web := web.NewWebFileServer()

	mux.Handle("/api/", http.StripPrefix("/api", api))
	mux.Handle("/", web)

	slog.Info("starting HTTP server", "host", c.Server.Host, "port", c.Server.Port)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port), mux); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
