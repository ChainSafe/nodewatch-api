package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"eth2-crawler/crawler"
	"eth2-crawler/graph"
	"eth2-crawler/graph/generated"
	ipResolver "eth2-crawler/resolver"
	"eth2-crawler/store"
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
		log.Fatalf("error loading configuration: %s", err.Error())
	}

	peerStore, err := store.New(cfg.Database)
	if err != nil {
		log.Fatalf("error Initializing the peer store: %s", err.Error())
	}

	resolverService := ipResolver.New(cfg.Resolver.APIKey, time.Duration(cfg.Resolver.Timeout)*time.Second)

	// TODO collect config from a config files or from command args and pass to Start()
	go crawler.Start(peerStore, resolverService)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(peerStore)}))

	// TODO: make playground accessible only in Dev mode
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	router := http.NewServeMux()
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	server.Start(context.TODO(), cfg.Server, router)
}
