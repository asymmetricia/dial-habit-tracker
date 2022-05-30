package main

type Task struct {
	Label string
	Done  bool
}

type Color string

var (
	Red    Color = "Red"
	Orange Color = "Orange"
	Yellow Color = "Yellow"
	Green  Color = "Green"
)

var Sets = map[Color]Task{}

const dataUrl = "dialhabittracker.dat"

func main() {
}
