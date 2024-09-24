package config

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/marin-vukojevic/todos/generated/database"
	"github.com/marin-vukojevic/todos/todos"
	"github.com/pressly/goose/v3"

	_ "github.com/lib/pq"
)

func Initialize() {
	godotenv.Load()

	db := createDB()

	executeMigrations(db)

	todoHandler := todos.NewTodoHandler(todos.NewTodoRepository(database.NewQueries(db)))

	router := createRouter(todoHandler)

	startAndWaitForInterrupt(router)
}

func createDB() *sql.DB {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("db url not found in the environment")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Cannot connect to db")
	}
	db.SetMaxOpenConns(5)
	return db
}

func executeMigrations(db *sql.DB) {
	migrationsDir := "./sql/schema"

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("Failed to set goose dialect:", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Migrations executed successfully")
}

func createRouter(todoHandler *todos.TodoHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)

	router.Get("/", todoHandler.Index)
	router.Route("/todo", func(r chi.Router) {
		r.Post("/", todoHandler.CreateTodo)
		r.Post("/{todoUuid}:complete", todoHandler.CompleteTodo)
	})

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	return router
}

func startAndWaitForInterrupt(router *chi.Mux) {
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port not found in the environment")
	}

	server := &http.Server{
		Addr:    ":" + portString,
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
