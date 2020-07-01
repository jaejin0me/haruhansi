package main

import (
	"net/http"
	"os"
	"time"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{os.Getenv("DBHOST")},
		Timeout:  60 * time.Second,
		Database: os.Getenv("DBAUTHDB"),
		Username: os.Getenv("DBUSER"),
		Password: os.Getenv("DBPW"),
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	s, err := mgo.DialWithInfo(mongoDBDialInfo)
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

	router.GET("/apoem/:id", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		session := mongoSession.Copy()
		defer session.Close()

		var poem Poem
		id := ps.ByName("id")
		var pipeline []bson.M

		collection := session.DB("kosmos").C("poems")
		if id != "empty" {
			pipeline = []bson.M{
				{"$match": bson.M{"_id": bson.M{"$ne": bson.ObjectId(id)}}},
				{"$sample": bson.M{"size": 1}},
			}
		} else {
			pipeline = []bson.M{{"$sample": bson.M{"size": 1}}}
		}

		err := collection.Pipe(pipeline).One(&poem)
		if err != nil {
			renderer.JSON(w, http.StatusInternalServerError, err)
			return
		}

		renderer.JSON(w, http.StatusOK, poem)
	})

	router.GET("/auth/:action/:provider", loginHandler)
	router.POST("/rooms", createRoom)
	router.GET("/rooms", retrieveRooms)
	//router.GET("/rooms/:id/messages", retrieveMessages)
	//router.GET("/ws/:room_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//	socket, err := upgrader.Upgrade(w, r, nil)
	//	if err != nil {
	//		log.Fatal("ServeHTTP:", err)
	//		return
	//	}
	//	newClient(socket, ps.ByName("room_id"), GetCurrentUser(r))
	//})

	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
		OptionsPassthrough: true,
	})

	handler := c.Handler(router)

	n := negroni.Classic()
	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))
	n.Use(LoginRequired("/login", "/auth", "/apoem"))
	n.UseHandler(handler)

	n.Run(":3000")
}
