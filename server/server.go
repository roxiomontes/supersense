package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/graph"
	"github.com/minskylab/supersense/graph/generated"
	"github.com/minskylab/supersense/persistence"
	"github.com/minskylab/supersense/sources"
	log "github.com/sirupsen/logrus"
)

const defaultPort = 8080

// LaunchServer launch the graphQL server
func LaunchServer(mux *supersense.Mux, port int64, withGraphQLPlayground bool, spokesman *sources.Spokesman, store persistence.Store) error {
	if port <= 0 {
		port = defaultPort
	}

	resolver := graph.NewResolver(mux, spokesman, store)

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	// Allowing introspection
	srv.Use(extension.Introspection{})

	// Serving static from observer build
	fs := http.FileServer(http.Dir("./observer_static"))
	http.Handle("/", fs)
	log.Infof("connect to http://localhost:%d/ for the observer", port)

	// Setting up graphql endpoint
	http.Handle("/graphql", srv)

	// If GraphQL Playground is enabled
	if withGraphQLPlayground {
		http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
		log.Infof("connect to http://localhost:%d/playground for GraphQL playground", port)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
