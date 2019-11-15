package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	svg "github.com/swill/svgo"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Task struct {
	Axiom        string `json:"axiom"`
	GenTypically map[int]struct {
		element string `json:"element"`
		rule    string `json:"rule"`
	} `json:"genTypically"`
	RotAngle float64 `json:"rotAngle"`
	Step     int `json:"step"`
	Depth    int `json:"depth"`
}

type Coordinate struct {
	x int
	y int
	angle float64
}

const (
	forward               = "F"
	forwardWithoutDrawing = "b"
	remember              = "["
	recall                = "]"
	turnClockwise         = "+"
	turnCounterClockwise  = "-"

	height = 1500
	width  = 1500

	lineStyle = `stroke="green" stroke-width="2"`
)

func readTask(name string) (Task, error) {
	path := filepath.Join("task", name)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Task{}, err
	}

	if !json.Valid(data) {
		err = errors.New("invalid json")
		return Task{}, err
	}

	var task Task
	err = json.Unmarshal(data, &task)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func saveSVG(buf *bytes.Buffer) {
	name := "snowFlake.svg"
	path := filepath.Join("image", name)
	err := ioutil.WriteFile(path, buf.Bytes(), 0777)
	if err != nil {
		fmt.Print(err)
	}
}

func createTask(task *Task) error {
	var buf = new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(task)
	if err != nil {
		return err
	}

	name := "snowFlake.json"
	path := filepath.Join("task", name)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	io.Copy(f, buf)
	return nil
}

func convertAxiom(task *Task) {
	for j:= 1; j <= task.Depth; j++ {
		for i := 1; i <= len(task.GenTypically); i++ {
			//fmt.Print(task)
			//fmt.Print("\n")
			//fmt.Print(task.GenTypically[1].element)
			//fmt.Print("\n")
			//fmt.Print(task.GenTypically[i].rule)
			//fmt.Print("\n")
			task.Axiom = strings.ReplaceAll(task.Axiom, "F", "F-F++F-F")
			//fmt.Print("\n")
			//fmt.Print(task.Axiom)
		}
	}

}

func drawSVG(task *Task) error {
	var buf = new(bytes.Buffer)
	canvas := svg.New(buf)
	canvas.Start(width, height)

	point := Coordinate{
		x: width/2,
		y: height/2,
		angle: 0,
	}

	convertAxiom(task)
	for _, i := range task.Axiom {
		switch string(i) {
		case forward:
			canvas.Rotate(point.angle)
			canvas.Line(point.x, point.y, point.x+task.Step, point.y+task.Step, lineStyle)
			canvas.Gend()
			point.x += task.Step
			point.y += task.Step
		case forwardWithoutDrawing:
			point.x += task.Step
			point.y += task.Step
		case remember:
		case recall:
		case turnClockwise:
			point.angle += task.RotAngle
			//canvas.Rotate(task.RotAngle)
		case turnCounterClockwise:
			//canvas.Rotate(-task.RotAngle)
			point.angle -= task.RotAngle
		}
	}
	canvas.End()

	saveSVG(buf)
	return nil
}

func generateFileTask(){
	var task = Task{
		Axiom: "F++F++F",
		GenTypically: map[int]struct {
			element string `json:"element"`
			rule    string `json:"rule"`
		}{1: {"F", "F-F++F-F"},},
		RotAngle: 60,
		Step:     5,
		Depth:    3,
	}

	err := createTask(&task)
	if err != nil {
		fmt.Print(err)
	}
}

func main() {
	generateFileTask()
	task, _ := readTask("snowFlake.json")

	drawSVG(&task)
}
