package api

import (
	"net/http"

	"github.com/LMBishop/confplanner/api/handlers"
	"github.com/LMBishop/confplanner/api/middleware"
	"github.com/LMBishop/confplanner/pkg/auth"
	"github.com/LMBishop/confplanner/pkg/calendar"
	"github.com/LMBishop/confplanner/pkg/conference"
	"github.com/LMBishop/confplanner/pkg/favourites"
	"github.com/LMBishop/confplanner/pkg/ical"
	"github.com/LMBishop/confplanner/pkg/session"
	"github.com/LMBishop/confplanner/pkg/user"
)

type ApiServices struct {
	UserService       user.Service
	FavouritesService favourites.Service
	ConferenceService conference.Service
	CalendarService   calendar.Service
	IcalService       ical.Service
	SessionService    session.Service
	AuthService       auth.Service
}

func NewServer(apiServices ApiServices, baseURL string) *http.ServeMux {
	mustAuthenticate := middleware.MustAuthenticate(apiServices.UserService, apiServices.SessionService)
	admin := middleware.MustAuthoriseAdmin(apiServices.UserService, apiServices.SessionService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", handlers.Register(apiServices.UserService, apiServices.AuthService))
	mux.HandleFunc("GET /login", handlers.GetLoginOptions(apiServices.AuthService))
	mux.HandleFunc("POST /login/{provider}", handlers.Login(apiServices.AuthService, apiServices.SessionService))
	mux.HandleFunc("POST /logout", mustAuthenticate(handlers.Logout(apiServices.SessionService)))

	mux.HandleFunc("GET /conference", mustAuthenticate(handlers.GetConferences(apiServices.ConferenceService)))
	mux.HandleFunc("GET /conference/{id}", mustAuthenticate(handlers.GetSchedule(apiServices.ConferenceService)))
	mux.HandleFunc("POST /conference", mustAuthenticate(admin(handlers.CreateConference(apiServices.ConferenceService))))
	mux.HandleFunc("DELETE /conference", mustAuthenticate(admin(handlers.DeleteConference(apiServices.ConferenceService))))

	mux.HandleFunc("GET /favourites/{id}", mustAuthenticate(handlers.GetFavourites(apiServices.FavouritesService)))
	mux.HandleFunc("POST /favourites", mustAuthenticate(handlers.CreateFavourite(apiServices.FavouritesService)))
	mux.HandleFunc("DELETE /favourites", mustAuthenticate(handlers.DeleteFavourite(apiServices.FavouritesService)))

	mux.HandleFunc("GET /calendar", mustAuthenticate(handlers.GetCalendar(apiServices.CalendarService, baseURL)))
	mux.HandleFunc("POST /calendar", mustAuthenticate(handlers.CreateCalendar(apiServices.CalendarService, baseURL)))
	mux.HandleFunc("DELETE /calendar", mustAuthenticate(handlers.DeleteCalendar(apiServices.CalendarService)))
	mux.HandleFunc("/calendar/ical", handlers.GetIcal(apiServices.IcalService, apiServices.CalendarService))

	return mux
}
