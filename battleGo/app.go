package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"context"

	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	addr             = "30124"
	shutdownDeadline = time.Second * 15
	//pem  = "/etc/certs/wildcard_cs_mtech_edu.key"
	//cert = "/etc/certs/wildcard_cs_mtech_edu.cer"
)

func main() {

	r := chi.NewRouter()

	r.Mount("/events/", routes.EventBroker)

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)

		// Set a timeout value on the request context, to signal when the request has timed out
		r.Use(middleware.Timeout(60 * time.Second))

		fileServer(r.(*chi.Mux))

		session := routes.NewSession()
		r.Mount("/bsState", routes.BsStateResource{}.Routes())
		r.Mount("/auth", routes.AuthResource{}.Routes())
		r.Mount("/session", session.Routes())

		r.With(session.BattlePhase, session.ActiveSessionCheck).Post("/target", session.PostTarget)
		r.Get("/battle/{filename}", session.Get)
		r.Get("/battle/{filename}/{url}", session.URLParam(session.Get))
	})

	//r.Use(middlewares.SessionResource)

	srv := &http.Server{
		Addr:    ":" + addr,
		Handler: r,
	}

	go func() {
		log.Printf("Listening on Port:%s", addr)
		//log.Fatal(srv.ListenAndServeTLS(cert, pem))
		log.Fatal(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive a signal from the OS to shutdown
	<-c

	// Create a deadline for tear down
	ctx, cancel := context.WithTimeout(context.Background(), shutdownDeadline)
	defer cancel()

	// Shutdown doesn't block if there are no connections, but will otherwise
	// wait until the timeout before closing connections and shutting down.
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}

	log.Println("shutting down")
	os.Exit(0)
}

func fileServer(router *chi.Mux) {
	root := "./public"

	fs := http.FileServer(http.Dir(root))
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
