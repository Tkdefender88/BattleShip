package routes

import (
	"crypto/sha256"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/go-chi/chi"
)

var (
	u = UserForm{
		UName:  "readBattleState",
		Passwd: [sha256.Size]byte{106, 36, 36, 56, 2, 151, 190, 194, 141, 236, 10, 63, 147, 82, 160, 95, 82, 86, 84, 183, 204, 221, 186, 123, 40, 129, 211, 30, 166, 7, 5, 74},
	}
)

type AuthResource struct{}

// UserFrom represents the url form encoded post request
// received from the login page.
type UserForm struct {
	UName  string
	Passwd [sha256.Size]byte
}

//"sv2zSY7WZK3xPPI2FIapggDBcwnUeoAj"

// Store ...
var Store *sessions.CookieStore

func init() {

	authKeyOne := []byte{179, 17, 111, 242, 236, 58, 163, 144, 110, 204, 183, 105, 45, 240, 97, 71, 181, 45, 162, 89, 60, 186, 69, 248, 230, 203, 110, 225, 22, 47, 233, 3, 247, 108, 105, 246, 183, 74, 155, 199, 108, 40, 118, 73, 7, 117, 163, 178, 51, 221, 54, 207, 52, 181, 148, 30, 129, 134, 149, 60, 29, 28, 190, 78}
	encryptionKeyOne := []byte{169, 144, 113, 1, 251, 78, 49, 234, 154, 138, 169, 239, 105, 150, 14, 94, 246, 255, 202, 7, 169, 100, 94, 162, 207, 38, 81, 199, 201, 42, 140, 87}

	Store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	Store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

}

// Authenticated is a middleware to ensure the client has a valid session.
func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "BATTLESHIP")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if auth, ok := session.Values["user"].(bool); !ok || !auth {
			session.Values["user"] = false
			session.AddFlash(r.RequestURI)
			if err := session.Save(r, w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func hashPw(pwd []byte) [sha256.Size]byte {
	return sha256.Sum256(pwd)
}

func (ar AuthResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Route("/login", func(r chi.Router) {
		r.Get("/", ar.LoginGET)
		r.Post("/", ar.LoginPOST)
	})

	return r
}

func (ar AuthResource) LoginPOST(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "BATTLESHIP")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		INTERNALERROR(w)
	}

	usrPass := r.FormValue("password")
	usrName := r.FormValue("username")

	if hashPw([]byte(usrPass)) == u.Passwd && u.UName == usrName {
		flashes := session.Flashes()
		referer := "/bsState/"
		if len(flashes) > 0 {
			ref := flashes[0].(string)
			referer = ref
		}

		session.Values["user"] = true

		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, referer, http.StatusFound)
	} else {
		errorMsg := ErrorMsg{Message: "Login Failed"}
		UNAUTHORIZEDPage(w, errorMsg)
	}
}

func (ar AuthResource) LoginGET(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/base.html", "views/login.html"))
	if err := tmpl.ExecuteTemplate(w, "base.html", nil); err != nil {
		INTERNALERROR(w)
		log.Println(err)
	}
}
