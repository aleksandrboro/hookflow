package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"

	"github.com/aleksandrboro/hookflow/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	var migrationsPath, migrationsType string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsType, "migrations-type", "", "type of migrations")

	flag.Parse()

	sourcePath := fmt.Sprintf("file://%s", migrationsPath)

	fmt.Println(sourcePath)

	var dbCfg config.DBCfg
	if err := cleanenv.ReadConfig(".env", &dbCfg); err != nil {
		panic(err)
	}

	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbCfg.User, dbCfg.Password),
		Host:   net.JoinHostPort(dbCfg.Host, dbCfg.Port),
		Path:   "/" + dbCfg.Name,
	}

	fmt.Println(databaseURL.String())

	m, err := migrate.New(sourcePath, databaseURL.String())
	if err != nil {
		panic(err)
	}

	switch migrationsType {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no change")
				return
			}
			panic(err)
		}
	case "down":
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no change")
				return
			}
			panic(err)
		}
	default:
		fmt.Println("invalid migration type")
		return
	}

	fmt.Println(m.Version())

	fmt.Printf("%s migrations applied successfully\n", migrationsType)
}
