package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stockyard-dev/stockyard-apothecary/internal/server"
	"github.com/stockyard-dev/stockyard-apothecary/internal/store"
)

var version = "dev"

func main() {
	portFlag := flag.String("port", "", "HTTP port (overrides PORT env var)")
	dataFlag := flag.String("data", "", "Data directory (overrides DATA_DIR env var)")
	flag.Parse()

	port := *portFlag
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "9700"
	}

	dataDir := *dataFlag
	if dataDir == "" {
		dataDir = os.Getenv("DATA_DIR")
	}
	if dataDir == "" {
		dataDir = "./apothecary-data"
	}

	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("apothecary: %v", err)
	}
	defer db.Close()

	srv := server.New(db, server.DefaultLimits(), dataDir)

	fmt.Printf("\n  Apothecary v%s — Self-hosted health and medication tracker\n", version)
	fmt.Printf("  Dashboard:  http://localhost:%s/ui\n", port)
	fmt.Printf("  API:        http://localhost:%s/api\n", port)
	fmt.Printf("  Data:       %s\n", dataDir)
	fmt.Printf("  Questions?  hello@stockyard.dev — I read every message\n\n")

	log.Printf("apothecary: listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}
