package main

import (
	"github.com/aidos-dev/habit-tracker/internal/app"
	_ "github.com/jackc/pgx/v5"
)

func main() {
	app.Run()
}
