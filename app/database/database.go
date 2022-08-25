package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Dbname   string `json:"dbname"`
}

type Board struct {
	Id    string
	SizeX int
	SizeY int
}

type Robot struct {
	Id        string
	X         int
	Y         int
	Direction string
	Board     *Board
}

var DB *sql.DB
var err error
var conf *Config

func GetConfig() *Config {
	// ONLY FOR TESTING PURPOSES
	out := Config{
		Username: "theuser",
		Password: "supersecretpassword",
		Host:     "postgres-devoteam",
		Port:     5432,
		Dbname:   "devoteam",
	}
	return &out
}

func Connection() {
	log.Println("Connecting to database")
	conf = GetConfig()
	if err != nil {
		panic(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.Username, conf.Password, conf.Dbname)
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}
}

func RunMigrations() {
	log.Println("Running migrations")
	migrateInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Username, conf.Password, conf.Host, conf.Port, conf.Dbname)
	m, err := migrate.New(
		"file://db/migrations",
		migrateInfo)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Dropping DB")
	//m.Down()
	//m.Drop()
	log.Println("Applying migrations")
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No changes found")
		}
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
