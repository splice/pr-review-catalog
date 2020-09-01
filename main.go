package main

import (
	"context"
	"flag"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/splice/catalog-interview/libs/golang/healthcheck"
	"github.com/splice/catalog-interview/libs/golang/requestid"
	"github.com/splice/catalog-interview/libs/golang/requestlogger"
	"github.com/splice/catalog-interview/server"
	"github.com/splice/catalog-interview/storage"
)

var (
	argPort = flag.Int("port", 0, "the port to run server on")
)

// Build is set at compile time with the git SHA.
var Build string

// These values will be replaced automatically during service creation.
const (
	serviceName = "catalog"
	serviceTeam = "service-team-name"
)

func main() {
	flag.Parse()

	lg := logrus.New()
	lg.SetFormatter(&logrus.JSONFormatter{})
	log := lg.WithFields(logrus.Fields{"service": serviceName})
	ctx := context.Background()
	ctx = requestlogger.ContextWithLogger(ctx, log)

	hc := healthcheck.New(serviceName)

	params := &server.PageControllerParams{
		Products: &storage.ProductRepository{
			StorageLoc: "storage/seeds/products.json",
		},
	}
	mux := server.NewPageController(params)

	mux.Handle("/health", hc).Methods(http.MethodGet)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := httputil.DumpRequest(r, true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(body)
	})

	handler := func(mux http.Handler) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Build", Build)
			w.Header().Set("X-Service", serviceName)

			handler := requestid.Middleware(requestlogger.Middleware(mux, serviceName))
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	}

	hc.Start(ctx)
	serverPort := ""
	if *argPort != 0 {
		serverPort = ":" + strconv.Itoa(*argPort)
		log.Printf("%q owned by %q running on port %s", serviceName, serviceTeam, serverPort)
	}
	log.Fatal(http.ListenAndServe(serverPort, http.HandlerFunc(handler(mux))))
}
