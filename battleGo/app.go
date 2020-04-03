package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"context"

	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/routes"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	addr             = "30124"
	RWTimeout        = time.Second * 15
	IdleTimeout      = time.Second * 60
	ShutdownDeadline = time.Second * 15
	//pem              = "/etc/certs/wildcard_cs_mtech_edu.key"
	//cert             = "/etc/certs/wildcard_cs_mtech_edu.cer"
)

func main() {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	//r.Use(middlewares.SessionResource)

	// Set a timeout value on the request context, to signal when the request has timed out
	r.Use(middleware.Timeout(60 * time.Second))

	events := sse.NewServer(nil)
	defer events.Shutdown()

	session := routes.NewSession()
	session.RegisterEventSource(events)

	/*
		go func() {
			for {
				s.SendMessage("/events/my-channel", sse.SimpleMessage(time.Now().Format("2006/02/01/ 15:04:05")))
				time.Sleep(5 * time.Second)
			}
		}()
	*/

	FileServer(r)
	r.Mount("/events/", events)
	r.Mount("/bsState", routes.BsStateResource{}.Routes())
	r.Mount("/auth", routes.AuthResource{}.Routes())
	r.Mount("/session", session.Routes())
	r.Mount("/target", session.TargetRoute())
	r.Mount("/battle", session.BattleRoute())

	srv := &http.Server{
		Addr: ":" + addr,

		WriteTimeout: RWTimeout,
		ReadTimeout:  RWTimeout,
		IdleTimeout:  IdleTimeout,
		Handler:      r,
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
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownDeadline)
	defer cancel()

	// Shutdown doesn't block if there are no connections, but will otherwise
	// wait until the timeout before closing connections and shutting down.
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
		return
	}

	log.Println("shutting down")
	os.Exit(0)
}

func FileServer(router *chi.Mux) {
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
