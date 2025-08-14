package api

import (
	"net/http"

	"github.com/LMBishop/confplanner/api/handlers"
	"github.com/LMBishop/confplanner/api/middleware"
	"github.com/LMBishop/confplanner/pkg/calendar"
	"github.com/LMBishop/confplanner/pkg/favourites"
	"github.com/LMBishop/confplanner/pkg/ical"
	"github.com/LMBishop/confplanner/pkg/schedule"
	"github.com/LMBishop/confplanner/pkg/session"
	"github.com/LMBishop/confplanner/pkg/user"
)

type ApiServices struct {
	UserService       user.Service
	FavouritesService favourites.Service
	ScheduleService   schedule.Service
	CalendarService   calendar.Service
	IcalService       ical.Service
	SessionService    session.Service
}

func NewServer(apiServices ApiServices, baseURL string) *http.ServeMux {
	mustAuthenticate := middleware.MustAuthenticate(apiServices.UserService, apiServices.SessionService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", handlers.Register(apiServices.UserService))
	mux.HandleFunc("POST /login", handlers.Login(apiServices.UserService, apiServices.SessionService))
	mux.HandleFunc("POST /logout", mustAuthenticate(handlers.Register(apiServices.UserService)))

	mux.HandleFunc("GET /favourites", mustAuthenticate(handlers.GetFavourites(apiServices.FavouritesService)))
	mux.HandleFunc("POST /favourites", mustAuthenticate(handlers.CreateFavourite(apiServices.FavouritesService)))
	mux.HandleFunc("DELETE /favourites", mustAuthenticate(handlers.DeleteFavourite(apiServices.FavouritesService)))

	mux.HandleFunc("GET /schedule", mustAuthenticate(handlers.GetSchedule(apiServices.ScheduleService)))

	mux.HandleFunc("GET /calendar", mustAuthenticate(handlers.GetCalendar(apiServices.CalendarService, baseURL)))
	mux.HandleFunc("POST /calendar", mustAuthenticate(handlers.CreateCalendar(apiServices.CalendarService, baseURL)))
	mux.HandleFunc("DELETE /calendar", mustAuthenticate(handlers.DeleteCalendar(apiServices.CalendarService)))
	mux.HandleFunc("/calendar/ical", handlers.GetIcal(apiServices.IcalService, apiServices.CalendarService))

	return mux
}
