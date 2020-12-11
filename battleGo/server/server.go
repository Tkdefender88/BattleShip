package server

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/repository"
	"gitea.justinbak.com/juicetin/bsStatePersist/battleGo/routes"
)

const (
	addr             = ":30124"
	shutdownDeadline = time.Second * 15
	pem              = "/etc/certs/wildcard_cs_mtech_edu.key"
	cert             = "/etc/certs/wildcard_cs_mtech_edu.cer"
)

func CreateRouter(repo repository.ModelRepository) *chi.Mux {
	controller := routes.NewBsStateController(repo)
	r := chi.NewRouter()
	r.Mount("/events", routes.EventBroker)
	r.Mount("/auth", routes.AuthResource{}.Routes())

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)
		// Set a timeout value on the request context, to signal when the request has timed out
		r.Use(middleware.Timeout(60 * time.Second))

		session := routes.NewSession(repo)
		// Unsecured routes
		r.With(session.BattlePhase).Mount("/bsProtocol/session", session.Routes())
		r.With(session.BattlePhase, session.ActiveSessionCheck).Post("/bsProtocol/target", session.PostTarget)

		// Secured routes
		r.Route("/", func(r chi.Router) {
			r.Use(routes.Refresh)
			r.Use(routes.Authenticated)
			fileServer(r)
			r.With(routes.Refresh, routes.Authenticated).Mount("/bsState", controller.Routes())
			r.With(routes.Refresh, routes.Authenticated).Get("/battle/{filename}", session.Get)
			r.With(routes.Refresh, routes.Authenticated).Get("/battle/{filename}/{url}", session.BattleURL(session.Get))
		})
	})
	return r
}

type ServeFunc func(s *http.Server) error

func StartServer(s *http.Server) error {
	return s.ListenAndServeTLS(cert, pem)
}

func StartDevServer(s *http.Server) error {
	return s.ListenAndServe()
}

func Start(serve ServeFunc, addr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("APIURI")))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	// Set up dependency injection
	database := client.Database("battleStatedb")
	bsStateRepo := repository.NewModelRepo(database)

	r := CreateRouter(bsStateRepo)

	srv := &http.Server{
		Addr:    ":" + addr,
		Handler: r,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	go func() {
		log.Printf("Listening on Port :%s", addr)
		log.Fatal(serve(srv))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive a signal from the OS to shutdown
	<-c

	// Create a deadline for tear down
	ctx, cancel = context.WithTimeout(context.Background(), shutdownDeadline)
	defer cancel()

	// Shutdown doesn't block if there are no connections, but will otherwise
	// wait until the timeout before closing connections and shutting down.
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}

	log.Println("shutting down")
	return nil
}

func fileServer(router chi.Router) {
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
