package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jieggii/softrust/services/backend/internal/adapters/mongoadapter"
	"github.com/jieggii/softrust/services/backend/internal/oapi"
	"github.com/jieggii/softrust/services/backend/internal/server"
	"github.com/jieggii/softrust/services/backend/internal/usecases"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/genai"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const (
	listenPort = 8080

	envVarGenAIAPIKey = "GENERATIVE_AI_API_KEY"
)

const (
	mongoDBName = "softrust"
	mongoURI    = "mongodb://mongo:27017"
)

func newGenAIClient(apiKey string) (*genai.Client, error) {
	ctx := context.Background()
	return genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
}

func newMongoDatabase(uri string, dbName string) (*mongo.Database, error) {
	// connect to mongoDB:
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("mongo connect error: %w", err)
	}

	// check the connection:
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping error: %w", err)
	}

	return client.Database(dbName), nil
}

func setupSwagger(r *chi.Mux) {
	// Add swagger routes BEFORE passing to oapi.HandlerFromMux
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		swagger, err := oapi.GetSwagger()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(swagger)
	})

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	// create mongo database client:
	db, err := newMongoDatabase(mongoURI, mongoDBName)
	if err != nil {
		logger.Fatalf("failed to connect to mongoDB: %v", err)
	}

	// create repo:
	repo := &mongoadapter.Repo{
		DB: db,
	}

	// create gen AI client:
	genAIAPIKey, ok := os.LookupEnv(envVarGenAIAPIKey)
	if !ok {
		logger.Fatalf("environment variable %s is not set", envVarGenAIAPIKey)
	}

	genAIClient, err := newGenAIClient(genAIAPIKey)
	if err != nil {
		logger.Fatalf("failed to create GenAI client: %v", err)
	}

	// create product meta resolver service:
	productMetaResolver := &usecases.ProductMetaResolver{
		Client: genAIClient,
	}

	// create product security assessor service:
	productSecurityAssessor := &usecases.ProductSecurityAssessor{
		Client: genAIClient,
	}

	// create report generator service:
	reportGeneratorSvc := &usecases.ReportContentGenerator{
		MetaResolver:     productMetaResolver,
		SecurityAssessor: productSecurityAssessor,
		Log:              logger,
	}

	// create the main service:
	svc := usecases.NewService(repo, reportGeneratorSvc, logger)

	// create service server:
	serviceServer := server.NewHTTP(svc, logger)

	// create HTTP router:
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// set up CORS:
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allow all origins for dev/testing
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(c.Handler)

	setupSwagger(r)

	// create HTTP handler:
	handler := oapi.NewStrictHandler(serviceServer, nil)

	// create and start HTTP server:
	addr := fmt.Sprintf("0.0.0.0:%d", listenPort)
	httpServer := &http.Server{
		Handler: oapi.HandlerFromMux(handler, r),
		Addr:    addr,
	}

	logger.Printf("HTTP server is listening on %v", addr)
	if err := httpServer.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
