package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"

	svg "github.com/swill/svgo"
)

type genTypically struct {
	Element string `json:"element"`
	Rule    string `json:"rule"`
}

type Task struct {
	Name         string         `json:"name"`
	Axiom        string         `json:"axiom"`
	GenTypically []genTypically `json:"genTypically"`
	RotAngle     float64        `json:"rotAngle"`
	Step         float64        `json:"step"`
	Depth        int            `json:"depth"`
}

type Coordinate struct {
	x     float64
	y     float64
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

	lineStyle = `stroke="black" stroke-width="2"`
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

func createTask(task *Task) error {
	buf, err := json.Marshal(task)
	if err != nil {
		return err
	}

	path := filepath.Join("task", task.Name+".json")
	err = ioutil.WriteFile(path, buf, 0777)
	if err != nil {
		return err
	}

	return nil
}

func convertAxiom(task *Task) {
	for j := 1; j <= task.Depth; j++ {
		for i := 0; i < len(task.GenTypically); i++ {
			task.Axiom = strings.ReplaceAll(task.Axiom, task.GenTypically[i].Element, task.GenTypically[i].Rule)
		}
	}

}

func saveSVG(buf *bytes.Buffer, name string) {
	path := filepath.Join("image", name+".svg")
	err := ioutil.WriteFile(path, buf.Bytes(), 0777)
	if err != nil {
		fmt.Print(err)
	}
}

func drawSVG(task *Task) error {
	var buf = new(bytes.Buffer)
	canvas := svg.New(buf)
	canvas.Start(width, height)

	point := Coordinate{
		x:     width / 2,
		y:     height / 2,
		angle: 0,
	}

	convertAxiom(task)
	fmt.Print(task.Axiom)
	for _, i := range task.Axiom {
		switch string(i) {
		case forward:
			x := task.Step * math.Cos(point.angle)
			y := task.Step * math.Sin(point.angle)
			canvas.Line(int(point.x), int(point.y), int(point.x+x), int(point.y+y), lineStyle)
			point.x += x
			point.y += y
		case forwardWithoutDrawing:
			point.x += task.Step
			point.y += task.Step
		case remember:
		case recall:
		case turnClockwise:
			point.angle += task.RotAngle
		case turnCounterClockwise:
			point.angle -= task.RotAngle
		}
	}
	canvas.End()

	saveSVG(buf, task.Name)
	return nil
}

func main() {
	//generateFileTask()
	task, _ := readTask("snowFlake.json")

	drawSVG(&task)
}

func generateFileTask() {
	var task = Task{
		Name:         "snowFlake",
		Axiom:        "F++F++F",
		GenTypically: []genTypically{0: {"F", "F-F++F-F"}},
		RotAngle:     60,
		Step:         700,
		Depth:        5,
	}

	//var task = Task{
	//	Name:         "triangle",
	//	Axiom:        "FXF--FF--FF",
	//	GenTypically: []genTypically{0: {"F", "FF"}, 1: {"X", "--FXF++FXF++FXF--"}},
	//	RotAngle:     60,
	//	Step:         12,
	//	Depth:        5,
	//}

	err := createTask(&task)
	if err != nil {
		fmt.Print(err)
	}
}
