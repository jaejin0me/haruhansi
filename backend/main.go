package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	type User struct {
		Id        string
		AddressId string
	}

	s := NewServer()

	s.HandleFunc("GET", "/", func(c *Context) {
		c.RenderTemplate("/public/index.html", map[string]interface{}{"time": time.Now()})
	})

	s.HandleFunc("GET", "/about", func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, "about")
	})

	s.HandleFunc("GET", "/login", func(c *Context) {
		c.RenderTemplate("public/login.html", map[string]interface{}{"message": "로그인이 필요합니다"})
	})

	s.HandleFunc("POST", "/login", func(c *Context) {
		if CheckLogin(c.Params["username"].(string), c.Params["password"].(string)) {
			http.SetCookie(c.ResponseWriter, &http.Cookie{
				Name:  "X_AUTH",
				Value: Sign(VerifyMessage),
				Path:  "/",
			})
			c.Redirect("/")
		}
		c.RenderTemplate("/public/login.html", map[string]interface{}{"message": "invalid info"})
	})

	s.HandleFunc("GET", "/users/:id", func(c *Context) {
		u := User{Id: c.Params["id"].(string)}
		c.RenderXml(u)
		fmt.Fprintf(c.ResponseWriter, "/users/:id %v\n", c.Params["id"])
	})

	s.HandleFunc("GET", "/users/:id/addresses/:address_id", func(c *Context) {
		u := User{Id: c.Params["id"].(string), AddressId: c.Params["address_id"].(string)}
		c.RenderXml(u)
		fmt.Fprintf(c.ResponseWriter, "/users/:id/addresses/:address_id %v %v\n", c.Params["id"], c.Params["address_id"])
	})

	s.HandleFunc("POST", "/users", func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, "create user")
	})

	s.HandleFunc("POST", "/users/:user_id/addresses", func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, "create user's address")
	})

	s.Use(AuthHandler)

	s.Run(":8080")
}

const VerifyMessage = "verified"

func AuthHandler(next HandlerFunc) HandlerFunc {
	ignore := []string{"/login", "public/index.html"}
	return func(c *Context) {
		for _, s := range ignore {
			if strings.HasPrefix(c.Request.URL.Path, s) {
				next(c)
				return
			}
		}

		if v, err := c.Request.Cookie("X_AUTH"); err == http.ErrNoCookie {
			c.Redirect("/login")
		} else if err != nil {
			c.RenderErr(http.StatusInternalServerError, err)
			return
		} else if Verify(VerifyMessage, v.Value) {
			next(c)
			return
		}
		c.Redirect("/login")
	}
}

func Verify(message, sig string) bool {
	return hmac.Equal([]byte(sig), []byte(Sign(message)))
}

func CheckLogin(username, password string) bool {
	const (
		USERNAME = "test"
		PASSWORD = "12345"
	)

	return username == USERNAME && password == PASSWORD
}

func Sign(message string) string {
	secreKey := []byte("golang-book-secret-key2")
	if len(secreKey) == 0 {
		return ""
	}
	mac := hmac.New(sha1.New, secreKey)
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}
