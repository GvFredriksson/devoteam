package database

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/lib/pq"
)

func Init_test_db(m *testing.M) (code int, err error) {
	// 1. create test.db if it does not exist
	// 2. run our migrations to create the required tables if they do not exist
	// 3. run our tests
	// 4. truncate the test db tables
	testdb := "test"
	conf = GetConfig()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.Username, conf.Password, conf.Dbname)
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return -1, fmt.Errorf("could not connect to database: %w", err)
	}
	_, err = DB.Exec(fmt.Sprintf("CREATE DATABASE %s", testdb))
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() != "duplicate_database" {
				log.Fatal(err)
			} else {
				log.Printf("duplicate database")
			}
		}
	} else {
		log.Printf("Created database " + testdb)
	}
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.Username, conf.Password, testdb)
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return -1, fmt.Errorf("could not connect to database: %w", err)
	}
	// Run Migrations for each table
	migrateInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Username, conf.Password, conf.Host, conf.Port, testdb)
	mig, err := migrate.New(
		"file:///usr/src/app/db/migrations",
		migrateInfo)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Migrating up")
	if err := mig.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Println("Up " + err.Error())
		}
		log.Println("Up " + err.Error())
	}
	return 1, nil
}

func BoardFactory(board *Board) *Board {
	var out Board
	sqlStatement := `
	INSERT INTO board (size_x, size_y)
	VALUES ($1, $2)
	RETURNING *`
	err = DB.QueryRow(
		sqlStatement,
		board.SizeX,
		board.SizeY).Scan(
		&out.Id,
		&out.SizeX,
		&out.SizeY)
	if err != nil {
		log.Printf("Error in BoardFactory : %s", err)
	}
	return &out
}

func RobotFactory(robot *Robot) *Robot {
	var out Robot
	out.Board = robot.Board
	sqlStatement := `
	INSERT INTO robot (x, y, direction, board_id)
	VALUES ($1, $2, $3, $4)
	RETURNING *`
	err = DB.QueryRow(
		sqlStatement,
		robot.X,
		robot.Y,
		robot.Direction,
		robot.Board.Id).Scan(
		&out.Id,
		&out.X,
		&out.Y,
		&out.Direction,
		&out.Board.Id)
	if err != nil {
		log.Printf("Error in RobotFactory : %s", err)
	}
	return &out
}
