package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	svg "github.com/swill/svgo"
	"io/ioutil"
	"math"
	"path/filepath"
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
	Width        float64            `json:"width"`
	Height       float64            `json:"height"`
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

	snowFlake = "snowFlake.json"
	triangle  = "triangle.json"
	dragon = "dragon.json"
	brokenLine = "brokenLine.json"
	tree_1 = "tree_1.json"
	tree_2 = "tree_2.json"
	tree_3 = "tree_3.json"
	tree_4 = "tree_4.json"

	height = 1500
	width  = 1500

	//heightStart = 500
	//widthStart = 20

	lineStyle = `stroke="yellow" stroke-width="2"`
)

func getRadian(g float64) float64 {
	return (g * math.Pi) / 180
}

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
		for i := 0; i < len(task.Axiom); i++ {
			for g := 0; g < len(task.GenTypically); g++ {

				if string(task.Axiom[i]) == task.GenTypically[g].Element {
					task.Axiom = task.Axiom[:i] + task.GenTypically[g].Rule + task.Axiom[i + 1:]
					i += len(task.GenTypically[g].Rule) - 1

				}
			}
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
	var (
		buf = new(bytes.Buffer)
		list []Coordinate
		pop Coordinate
	)
	canvas := svg.New(buf)
	canvas.Start(width, height)

	point := Coordinate{
		x:     task.Width,
		y:     task.Height,
		angle: 0,
	}
	list = append(list, point)

	convertAxiom(task)
	for _, i := range task.Axiom {
		switch string(i) {
		case forward:
			x := task.Step * math.Cos(getRadian(point.angle))
			y := task.Step * math.Sin(getRadian(point.angle))
			canvas.Line(int(point.x), int(point.y), int(point.x+x), int(point.y+y), lineStyle)
			point.x += x
			point.y += y
		case forwardWithoutDrawing:
			point.x += task.Step
			point.y += task.Step
		case remember:
			list = append(list, point)
		case recall:
			pop, list = list[len(list)-1], list[:len(list)-1]
			point.x = pop.x
			point.y = pop.y
			point.angle = pop.angle
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
	task, _ := readTask(dragon)

	err := drawSVG(&task)
	if err != nil {
		fmt.Print(err)
	}
}

func generateFileTask() {
	var task = Task{
		Name:         "snowFlake",
		Axiom:        "F++F++F",
		GenTypically: []genTypically{0: {"F", "F-F++F-F"}},
		RotAngle:     60,
		Step:         700,
		Depth:        5,
		Width:        300,
		Height:       500,
	}

	//var task = Task{
	//	Name:         "dragon",
	//	Axiom:        "FX",
	//	GenTypically: []genTypically{
	//		0:{"X", "X+YF+"},
	//		1:{"Y", "-FX-Y"}},
	//	RotAngle:     90,
	//	Step:         20,
	//	Depth:        5,
	//}

	//var task = Task{
	//	Name:         "brokenLine",
	//	Axiom:        "X",
	//	GenTypically: []genTypically{
	//		0:{"X", "-YF+XFX+FY-"},
	//		1:{"Y", "+XF-YFY-FX+"}},
	//	RotAngle:     90,
	//	Step:         20,
	//	Depth:        4,
	//}

	//var task = Task{
	//	Name:         "tree_1",
	//	Axiom:        "F",
	//	GenTypically: []genTypically{
	//		0:{"F", "F[+F]F[-F]F"}},
	//	RotAngle:     25.7,
	//	Step:         20,
	//	Depth:        6,
	//}

	//var task = Task{
	//	Name:         "tree_2",
	//	Axiom:        "F",
	//	GenTypically: []genTypically{
	//		0:{"F", "F[+F]F[-F][F]"}},
	//	RotAngle:     20,
	//	Step:         20,
	//	Depth:        6,
	//}

	//var task = Task{
	//	Name:         "tree_3",
	//	Axiom:        "X",
	//	GenTypically: []genTypically{
	//		0:{"F", "FF"},
	//		2:{"X", "F[+X][-X]FX"}},
	//	RotAngle:     25.7,
	//	Step:         8,
	//	Depth:        6,
	//}

	//var task = Task{
	//	Name:         "tree_4",
	//	Axiom:        "F",
	//	GenTypically: []genTypically{
	//		0:{"F", "-F[-F+F-F]+[+F-F-F]"}},
	//	RotAngle:     20,
	//	Step:         30,
	//	Depth:        5,
	//}

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
