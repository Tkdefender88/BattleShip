package routes

import (
	"crypto/sha256"
	"html/template"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/go-chi/chi"
)

type (
	Credentials struct {
		Password string `json:"password"`
		Username string `json:"username"`
	}
	Claims struct {
		Username string `json:"username"`
		jwt.StandardClaims
	}
	AuthResource struct{}
)

var (
	users = map[string]string{
		"justin": "babylon",
	}

	jwtKey = []byte("my_secret_key")
)

const (
	cookieName = "token"
)

// Routes returns a router with all the login endpoints
func (ar AuthResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/login", func(r chi.Router) {
		r.Get("/", loginPage)
		r.Post("/", signIn)
	})

	return r
}

//Refresh will refresh a valid token if it is within 30 seconds of the expiration time.
func Refresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				log.Printf("err %s\n", err)
				next.ServeHTTP(w, r)
				return
			}
			log.Printf("err %s\n", err)
			next.ServeHTTP(w, r)
			return
		}

		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				unauthorized(w)
				return
			}
			log.Printf("err %s\n", err)
			next.ServeHTTP(w, r)
			return
		}

		if !tkn.Valid {
			log.Printf("token not valid\n")
			next.ServeHTTP(w, r)
			return
		}

		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
			log.Printf("token not expiring\n")
			next.ServeHTTP(w, r)
			return
		}

		expireTime := time.Now().Add(5 * time.Minute)
		claims.ExpiresAt = expireTime.Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			internalError(w)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    cookieName,
			Value:   tokenString,
			Expires: expireTime,
		})
	})
}

// Authenticated is a middleware to ensure the client has a valid session.
func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				log.Printf("error %s\n", err)
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}
			badRequest(w, "")
			return
		}

		tknStr := c.Value

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.Printf("error %s\n", err)
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}
			badRequest(w, "")
			return
		}

		if !tkn.Valid {
			log.Printf("token not valid\n")
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func signIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := r.ParseForm(); err != nil {
		badRequest(w, err.Error())
		return
	}
	creds.Username = r.FormValue("username")
	creds.Password = r.FormValue("password")

	expectedPasswd, ok := users[creds.Username]

	if !ok || expectedPasswd != creds.Password {
		unauthorized(w)
		return
	}

	expireTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		internalError(w)
		return
	}

	log.Println("Setting cookie")
	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Domain:  "localhost",
		Path:    "/",
		Value:   tokenString,
		Expires: expireTime,
	})

	referer := r.Referer()
	if referer == "" {
		referer = "/"
	}

	log.Printf("redirecting to %s\n", referer)

	http.Redirect(w, r, "/", http.StatusFound)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/base.html", "views/login.html"))
	if err := tmpl.ExecuteTemplate(w, "base.html", nil); err != nil {
		internalError(w)
		log.Println(err)
	}
}

func hashPw(pwd []byte) [sha256.Size]byte {
	return sha256.Sum256(pwd)
}
