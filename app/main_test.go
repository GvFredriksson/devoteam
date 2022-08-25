package main

import (
	"devoteam-api/app/database"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// os.Exit skips defer calls
	// so we need to call another function
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	// 1. create test.db if it does not exist
	// 2. run migrations to create the required tables if they do not exist
	// 3. run tests
	// 4. close the db
	code, err = database.Init_test_db(m)
	if err != nil {
		return code, err
	}

	// truncates all test data after the tests are run
	defer func() {
		database.DB.Close()
	}()

	return m.Run(), nil
}

func TestInitBoard(t *testing.T) {
	form := url.Values{}
	form.Set("SizeX", "5")
	form.Set("SizeY", "5")
	req, err := http.NewRequest("POST", "/init_board", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	InitiateBoard(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	log.Printf(string(data))
	assert.Containsf(t, string(data), "New Board created with id", "formatted")
}

func TestInitBoardNegative(t *testing.T) {
	form := url.Values{}
	form.Set("SizeX", "5")
	form.Set("SizeY", "-5")
	req, err := http.NewRequest("POST", "/init_board", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	InitiateBoard(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Containsf(t, string(data), "The size of each side needs to be larger than 0", "formatted")
}

func TestInitBoardZero(t *testing.T) {
	form := url.Values{}
	form.Set("SizeX", "0")
	form.Set("SizeY", "0")
	req, err := http.NewRequest("POST", "/init_board", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	InitiateBoard(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Containsf(t, string(data), "The size of each side needs to be larger than 0", "formatted")
}

func TestInitRobot(t *testing.T) {
	board := &database.Board{
		SizeX: 5,
		SizeY: 5,
	}
	board = database.BoardFactory(board)
	form := url.Values{}
	form.Set("X", "5")
	form.Set("Y", "5")
	form.Set("Direction", "E")
	form.Set("Board", board.Id)
	req, err := http.NewRequest("POST", "/init_robot", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	InitiateRobot(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Containsf(t, string(data), "New Robot created with id", "formatted")
}

func TestInitRobotNegative(t *testing.T) {
	board := &database.Board{
		SizeX: 5,
		SizeY: 5,
	}
	board = database.BoardFactory(board)
	form := url.Values{}
	form.Set("X", "-5")
	form.Set("Y", "-5")
	form.Set("Direction", "E")
	form.Set("Board", board.Id)
	req, err := http.NewRequest("POST", "/init_robot", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	InitiateRobot(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Containsf(t, string(data), "The Robot can not start outside the board", "formatted")
}

func TestInitRobotOutsideBoard(t *testing.T) {
	board := &database.Board{
		SizeX: 5,
		SizeY: 5,
	}
	board = database.BoardFactory(board)
	form := url.Values{}
	form.Set("X", "6")
	form.Set("Y", "6")
	form.Set("Direction", "E")
	form.Set("Board", board.Id)
	req, err := http.NewRequest("POST", "/init_robot", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	InitiateRobot(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Containsf(t, string(data), "The Robot can not start outside the board", "formatted")
}

func TestInitRobotWrongDirection(t *testing.T) {
	board := &database.Board{
		SizeX: 5,
		SizeY: 5,
	}
	board = database.BoardFactory(board)
	form := url.Values{}
	form.Set("X", "6")
	form.Set("Y", "-5")
	form.Set("Direction", "K")
	form.Set("Board", board.Id)
	req, err := http.NewRequest("POST", "/init_robot", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	InitiateRobot(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	assert.Containsf(t, string(data), "The Direction can only be N or E or S or W", "formatted")
}

func TestMoveRobot(t *testing.T) {
	board := &database.Board{
		SizeX: 5,
		SizeY: 5,
	}
	board = database.BoardFactory(board)
	robot := &database.Robot{
		X:         1,
		Y:         2,
		Direction: "N",
		Board:     board,
	}
	robot = database.RobotFactory(robot)
	form := url.Values{}
	form.Set("Moves", "RFRFFRFRF")
	req, err := http.NewRequest("POST", "/move_robot", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	MoveRobot(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	log.Printf(string(data))
	assert.Equal(t, string(data), "X:1, Y:3, Direction:N\n")
}

func TestMoveRobotBadInput(t *testing.T) {
	board := &database.Board{
		SizeX: 5,
		SizeY: 5,
	}
	board = database.BoardFactory(board)
	robot := &database.Robot{
		X:         1,
		Y:         2,
		Direction: "N",
		Board:     board,
	}
	robot = database.RobotFactory(robot)
	form := url.Values{}
	form.Set("Moves", "RFRFFKRFRF")
	req, err := http.NewRequest("POST", "/move_robot", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	MoveRobot(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	log.Printf(string(data))
	assert.Equal(t, string(data), "Only the letters L(eft) R(ight) F(orward) can be in the string\n")
}
