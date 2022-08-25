package main

import (
	"database/sql"
	"devoteam-api/app/database"
	"fmt"
	"log"
	"net/http"
	"strings"

	"strconv"
)

func ClearBoards(res http.ResponseWriter, req *http.Request) error {
	// Task specification limits the request data so id can not be used
	// Clean all data so there is only ever one board in the database
	log.Println("Clearing away the old board")
	_, err := database.DB.Query("DELETE FROM robot")
	_, err = database.DB.Query("DELETE FROM board")
	return err
}

func CreateBoard(res http.ResponseWriter, req *http.Request) (*database.Board, error) {
	err := ClearBoards(res, req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	board := database.Board{}
	board.SizeX, _ = strconv.Atoi(req.FormValue("SizeX"))
	board.SizeY, _ = strconv.Atoi(req.FormValue("SizeY"))
	if board.SizeX < 1 || board.SizeY < 1 {
		fmt.Fprintf(res, "The size of each side needs to be larger than 0\n")
	}

	sqlStatement := `
	INSERT INTO board (size_x, size_y)
	VALUES ($1, $2)
	RETURNING *`
	err = database.DB.QueryRow(
		sqlStatement,
		board.SizeX,
		board.SizeY).Scan(
		&board.Id,
		&board.SizeX,
		&board.SizeY)

	return &board, err
}

func CreateRobot(res http.ResponseWriter, req *http.Request, board *database.Board) (*database.Robot, error) {
	robot := database.Robot{}
	robot.X, _ = strconv.Atoi(req.FormValue("X"))
	robot.Y, _ = strconv.Atoi(req.FormValue("Y"))
	robot.Direction = req.FormValue("Direction")
	robot.Board = board

	if robot.X < 0 || robot.Y < 0 || robot.X > board.SizeX || robot.Y > board.SizeY {
		fmt.Fprintf(res, "The Robot can not start outside the board\n")
	}
	dir := strings.Contains("NESW", robot.Direction)
	if len(robot.Direction) > 1 || dir == false {
		fmt.Fprintf(res, "The Direction can only be N or E or S or W\n")
	}

	sqlStatement := `
	INSERT INTO robot (x, y, direction, board_id)
	VALUES ($1, $2, $3, $4)
	RETURNING *`

	err := database.DB.QueryRow(
		sqlStatement,
		robot.X,
		robot.Y,
		robot.Direction,
		robot.Board.Id).Scan(
		&robot.Id,
		&robot.X,
		&robot.Y,
		&robot.Direction,
		&robot.Board.Id)

	return &robot, err
}

func InitiateBoard(res http.ResponseWriter, req *http.Request) {
	// If not post request return explanation else create the board
	if req.Method != "POST" {
		fmt.Fprintf(res, "Post request needed\n")
		fmt.Fprintf(res, "SizeX positive integer\nSizeY positive integer\n")
		return
	}
	board, err := CreateBoard(res, req)
	if err != nil {
		fmt.Fprintf(res, "Error occurred %s!\n", err.Error())
	}
	fmt.Fprintf(res, "New Board created with id %s!\n", board.Id)
}

func ParseBoard(rows *sql.Rows) *database.Board {
	// Helper function to parse Board Response data from the db
	var board database.Board

	err := rows.Scan(
		&board.Id,
		&board.SizeX,
		&board.SizeY,
	)
	if err != nil {
		log.Println(err)
	}
	return &board
}

func ParseRobot(rows *sql.Rows, board *database.Board) *database.Robot {
	// Helper function to parse Robot Response data from the db
	var robot database.Robot
	robot.Board = board
	err := rows.Scan(
		&robot.Id,
		&robot.X,
		&robot.Y,
		&robot.Direction,
		&robot.Board.Id,
	)
	if err != nil {
		log.Println(err)
	}
	return &robot
}

func GetBoard() *database.Board {
	var board *database.Board
	// Task specification limits the request data so id can not be used
	// If the database can contain more than one object expand select to use board id
	rows, err := database.DB.Query("SELECT * FROM board")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			board = ParseBoard(rows)
		}
		err = rows.Err()
	}
	if err != nil {
		log.Println(err)
	}
	return board
}

func GetRobot(board *database.Board) *database.Robot {
	var robot *database.Robot
	// Task specification limits the request data so id can not be used
	// If the database can contain more than one object expand select to use id
	rows, err := database.DB.Query("SELECT * FROM robot WHERE board_id = $1", board.Id)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			robot = ParseRobot(rows, board)
		}
		err = rows.Err()
	}
	if err != nil {
		log.Println(err)
	}
	return robot
}

func InitiateRobot(res http.ResponseWriter, req *http.Request) {
	// If not post request return explanation else create the robot
	if req.Method != "POST" {
		fmt.Fprintf(res, "Post request needed\n")
		fmt.Fprintf(res, "X positive integer\nY positive integer\nDirection N,E,S,W\n")
		return
	}
	board := GetBoard()
	robot, err := CreateRobot(res, req, board)
	if err != nil {
		fmt.Fprintf(res, "Error occurred %s!\n", err.Error())
	}
	fmt.Fprintf(res, "New Robot created with id %s!\n", robot.Id)
}

func MoveRobot(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fmt.Fprintf(res, "Post request needed\n")
		fmt.Fprintf(res, "A string with the letters L(eft) R(ight) F(orward), any order, any amount\n")
		return
	}
	board := GetBoard()
	robot := GetRobot(board)
	moves := req.FormValue("Moves")
	// L Turn left
	// R Turn right
	// F Walk forward
	for _, char := range moves {
		c := string(char)
		if !strings.Contains("LRF", c) {
			fmt.Fprintf(res, "Only the letters L(eft) R(ight) F(orward) can be in the string\n")
			return
		}
		switch robot.Direction {
		case "N":
			switch c {
			case "L":
				robot.Direction = "W"
			case "R":
				robot.Direction = "E"
			case "F":
				if robot.Y > 0 {
					robot.Y -= 1
				}
			}
		case "E":
			switch c {
			case "L":
				robot.Direction = "N"
			case "R":
				robot.Direction = "S"
			case "F":
				if robot.X < board.SizeX {
					robot.X += 1
				}
			}
		case "S":
			switch c {
			case "L":
				robot.Direction = "E"
			case "R":
				robot.Direction = "W"
			case "F":
				if robot.Y < board.SizeY {
					robot.Y += 1
				}
			}
		case "W":
			switch c {
			case "L":
				robot.Direction = "S"
			case "R":
				robot.Direction = "N"
			case "F":
				if robot.X > 0 {
					robot.X -= 1
				}
			}
		}
	}
	fmt.Fprintf(res, "X:%d, Y:%d, Direction:%s\n", robot.X, robot.Y, robot.Direction)
}

func main() {
	log.Println("Starting Robot Service")
	database.Connection()
	database.RunMigrations()
	http.HandleFunc("/", InitiateBoard)
	http.HandleFunc("/initiate_board", InitiateBoard)
	http.HandleFunc("/initiate_robot", InitiateRobot)
	http.HandleFunc("/move_robot", MoveRobot)
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
