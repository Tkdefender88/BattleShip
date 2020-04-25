package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Tkdefender88/BattleShip/battleGo/battlestate"

	"context"

	"github.com/Tkdefender88/BattleShip/battleGo/bsprotocol"
	"github.com/Tkdefender88/BattleShip/battleGo/routes"
	"github.com/Tkdefender88/BattleShip/battleGo/sse"

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

	r := router()

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

func fileServer(router chi.Router) {
	root := "./public"

	fs := http.FileServer(http.Dir(root))
	router.With(routes.Refresh, routes.Authenticated).Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}

func router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)                    // add IP address headers
	r.Use(middleware.Recoverer)                 // Recover gracefully and print trace
	r.Use(middleware.Logger)                    // Logging
	r.Use(middleware.Timeout(60 * time.Second)) // Set a timeout for requests

	eventBroker := sse.NewBroker()

	fileServer(r)

	r.Route("/", func(r chi.Router) {

		session := bsprotocol.NewSession(eventBroker)
		r.Mount("/bsState", battlestate.BsStateResource{}.Routes())
		r.Mount("/bsProtocol", session.Routes())

	})

	return r
}
