package main

import (
	"log"
	"net/http"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2"
)

const (
	sessionKey       = "simple_chat_session"
	sessionSecret    = "simple_chat_session_secret"
	socketBufferSize = 1024
)

var (
	renderer     *render.Render
	mongoSession *mgo.Session
	upgrader     = &websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
	}
)

func init() {
	renderer = render.New()

	s, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	mongoSession = s
}

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// 기본페이지
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "Simple chat!"})
	})

	router.GET("/login", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// 로그인 페이지
		renderer.HTML(w, http.StatusOK, "login", nil)
	})

	router.GET("/logout", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// 로그아웃 페이지
		sessions.GetSession(req).Delete(currentUserKey)
		http.Redirect(w, req, "/login", http.StatusFound)
	})

	router.GET("/auth/:action/:provider", loginHandler)
	router.POST("/rooms", createRoom)
	router.GET("/rooms", retrieveRooms)
	router.GET("/rooms/:id/messages", retrieveMessages)
	router.GET("/ws/:room_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("ServeHTTP:", err)
			return
		}
		newClient(socket, ps.ByName("room_id"), GetCurrentUser(r))
	})

	n := negroni.Classic()
	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))
	n.Use(LoginRequired("/login", "/auth"))
	n.UseHandler(router)

	n.Run(":3000")
}
