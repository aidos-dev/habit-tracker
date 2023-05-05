package main

import (
	"github.com/aidos-dev/habit-tracker/internal/app"
	_ "github.com/lib/pq"
)

// TODO: Изучи драйвер pgx
func main() {
	app.Run()
}
