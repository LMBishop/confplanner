package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/LMBishop/confplanner/api"
	config "github.com/LMBishop/confplanner/internal"
	"github.com/LMBishop/confplanner/pkg/calendar"
	"github.com/LMBishop/confplanner/pkg/database"
	"github.com/LMBishop/confplanner/pkg/favourites"
	"github.com/LMBishop/confplanner/pkg/ical"
	"github.com/LMBishop/confplanner/pkg/schedule"
	"github.com/LMBishop/confplanner/pkg/user"
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
	scheduleService, err := schedule.NewService(c.Conference.ScheduleURL)
	if err != nil {
		return fmt.Errorf("failed to create schedule service: %w", err)
	}
	calendarService := calendar.NewService(pool)
	icalService := ical.NewService(favouritesService, scheduleService)

	app := api.NewServer(api.ApiServices{
		UserService:       userService,
		FavouritesService: favouritesService,
		ScheduleService:   scheduleService,
		CalendarService:   calendarService,
		IcalService:       icalService,
	}, c.BaseURL)

	slog.Info("Server is listening", "host", c.Server.Host, "port", c.Server.Port)

	if err := app.Listen(fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
