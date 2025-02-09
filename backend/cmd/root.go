package cmd

import (
	"backend/internal/controller"
	"backend/internal/db/queries"
	"context"
	"errors"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

var (
	port          string
	dbSource      string
	migrationsDir string

	rootCmd = &cobra.Command{
		Use:   "backend",
		Short: "backend of the pinger application",
		Run:   Run,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&dbSource, "db_source",
		"postgres://root:secret@localhost:5432/example?sslmode=disable",
		"Database connection string")
	rootCmd.PersistentFlags().StringVar(&migrationsDir, "migrations",
		"migrations",
		"Path to the migrations folder")
	rootCmd.PersistentFlags().StringVar(&port, "port", "8091", "HTTP Server port")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func Run(cmd *cobra.Command, args []string) {
	ctx := context.TODO()

	m, err := migrate.New("file://"+migrationsDir, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries := queries.New(pool)
	controller := controller.NewContainerController(queries, ctx)

	server := gin.Default()

	server.Use(cors.Default())

	server.GET("/containers", controller.ListContainerStatuses)
	server.POST("/containers/ping", controller.PingHandler)

	log.Fatal(server.Run(":" + port))
}
