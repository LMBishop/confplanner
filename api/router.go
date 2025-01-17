package api

import (
	"errors"
	"log/slog"
	"time"

	"github.com/LMBishop/confplanner/api/dto"
	"github.com/LMBishop/confplanner/api/handlers"
	"github.com/LMBishop/confplanner/api/middleware"
	"github.com/LMBishop/confplanner/pkg/favourites"
	"github.com/LMBishop/confplanner/pkg/schedule"
	"github.com/LMBishop/confplanner/pkg/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type ApiServices struct {
	UserService       user.Service
	FavouritesService favourites.Service
	ScheduleService   schedule.Service
}

func NewServer(apiServices ApiServices) *fiber.App {
	sessionStore := session.New(session.Config{
		Expiration:     24 * time.Hour,
		KeyLookup:      "cookie:confplanner_session",
		CookieSameSite: "None",
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

	return app
}
