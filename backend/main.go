package main

import (
	"database/sql"
	"fmt"
	"os"

	appRouter "flowsync-pulse/backend/internal/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := connectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := appRouter.New(db)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}

func connectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		dbUser,
		password,
		host,
		port,
		name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
