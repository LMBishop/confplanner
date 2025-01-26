package api

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"time"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/api/handlers"
	"github.com/LMBishop/confplanner/api/middleware"
	"github.com/LMBishop/confplanner/pkg/calendar"
	"github.com/LMBishop/confplanner/pkg/favourites"
	"github.com/LMBishop/confplanner/pkg/ical"
	"github.com/LMBishop/confplanner/pkg/schedule"
	"github.com/LMBishop/confplanner/pkg/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type ApiServices struct {
	UserService       user.Service
	FavouritesService favourites.Service
	ScheduleService   schedule.Service
	CalendarService   calendar.Service
	IcalService       ical.Service
}

func NewServer(apiServices ApiServices, baseURL string) *fiber.App {
	sessionStore := session.New(session.Config{
		Expiration:     7 * 24 * time.Hour,
		KeyGenerator:   generateSessionToken,
		KeyLookup:      "cookie:confplanner_session",
		CookieSameSite: "Strict",
		CookieSecure:   true,
	})

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			var ok *dto.OkResponse
			if errors.As(err, &ok) {
				return ctx.Status(ok.Code).JSON(ok)
			}

			var e *dto.ErrorResponse
			if errors.As(err, &e) {
				return ctx.Status(e.Code).JSON(e)
			}

			slog.Error("fiber runtime error", "error", err, "URL", ctx.OriginalURL())
			return ctx.Status(500).JSON(dto.ErrorResponse{
				Code:    500,
				Message: "Internal Server Error",
			})
		},
		AppName: "confplanner",
	})

	// app.Use(cors.New())

	requireAuthenticated := middleware.RequireAuthenticated(apiServices.UserService, sessionStore)

	app.Post("/register", handlers.Register(apiServices.UserService))
	app.Post("/login", handlers.Login(apiServices.UserService, sessionStore))
	app.Post("/logout", requireAuthenticated, handlers.Logout(sessionStore))

	app.Get("/favourites", requireAuthenticated, handlers.GetFavourites(apiServices.FavouritesService))
	app.Post("/favourites", requireAuthenticated, handlers.CreateFavourite(apiServices.FavouritesService))
	app.Delete("/favourites", requireAuthenticated, handlers.DeleteFavourite(apiServices.FavouritesService))

	app.Get("/schedule", requireAuthenticated, handlers.GetSchedule(apiServices.ScheduleService))

	app.Get("/calendar", requireAuthenticated, handlers.GetCalendar(apiServices.CalendarService, baseURL))
	app.Post("/calendar", requireAuthenticated, handlers.CreateCalendar(apiServices.CalendarService, baseURL))
	app.Delete("/calendar", requireAuthenticated, handlers.DeleteCalendar(apiServices.CalendarService))
	app.Use("/calendar/ical", handlers.GetIcal(apiServices.IcalService, apiServices.CalendarService))

	return app
}

func generateSessionToken() string {
	b := make([]byte, 100)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
