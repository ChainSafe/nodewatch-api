package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"eth2-crawler/crawler"
	"eth2-crawler/graph"
	"eth2-crawler/graph/generated"
	"eth2-crawler/utils/config"
	"eth2-crawler/utils/server"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	cfgPath := flag.String("p", "./cmd/config/config.dev.yaml", "The configuration path")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("error reading configuration files: %s", err.Error())
	}

	// TODO collect config from a config files or from command args and pass to Start()
	go crawler.Start()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	// TODO: make playground accessible only in Dev mode
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	router := http.NewServeMux()
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	server.Start(context.TODO(), cfg.Server, router)
}
