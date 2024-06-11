package server

import (
	"log"
	"net/http"
	"postsandcomments/internal/graph"
	"postsandcomments/internal/db"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sirupsen/logrus"
)

func StartServer(port string, db db.Database) {
	cfg := graph.Config{
		Resolvers: &graph.Resolver{
			DataBase:            db,
			SubscriptionManager: graph.NewSubscriptionManager(),
			Logger:              logrus.New(),
		},
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(cfg))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

