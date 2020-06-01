package main

import (
	"log"
	"net/http"
	"strings"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/urfave/negroni"
)

const (
	nextPageKey = "next_page"
	authSecurityKey = "auth_security_key"
)

func init() {
	gomniauth.SetSecurityKey(authSecurityKey)
	gomniauth.WithProviders(
		google.New("671717321525-677jl2c8etmr564otgfr9sa1625bvbm9.apps.googleusercontent.com","kbj7usLLLOEk0aI4KMaGR43_","http://127.0.0.1:3000/auth/callback/google"),
	)
}

func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	provider := ps.ByName("provider")
	s := sessions.GetSession(r)

	switch action {
	case "login":
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}

		loginUrl, err := p.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln(err)
		}

		http.Redirect(w, r, loginUrl, http.StatusFound)
	case "callback":
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}

		creds, err := p.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln(err)
		}

		user, err := p.GetUser(creds)
		if err != nil {
			log.Fatalln(err)
		}

		u := &User {
			Uid: user.Data().Get("id").MustStr(),
			Name: user.Name(),
			Email: user.Email(),
			AvatarUrl: user.AvatarURL(),
		}

		SetCurrentUser(r, u)
		http.Redirect(w, r, s.Get(nextPageKey).(string), http.StatusFound)
	default:
		http.Error(w, "Auth action '" + action + "'is not supported", http.StatusNotFound)
	}
}

func LoginRequired(ignore ... string) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		for _, s := range ignore {
			if strings.HasPrefix(r.URL.Path, s) {
				next(w, r)
				return
			}
		}
		
		u := GetCurrentUser(r)

		if u != nil && u.Valid() {
			SetCurrentUser(r, u)
			next(w, r)
			return
		}

		SetCurrentUser(r, nil)

		sessions.GetSession(r).Set(nextPageKey, r.URL.RequestURI())

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}