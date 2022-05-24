package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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

//func init() {
//	data, err := storage.LoadResourceFromURI(storage.NewFileURI(dataUrl))
//	if err != nil {
//		log.Print(err)
//		return
//	}
//	if err := json.Unmarshal(data.Content(), &Sets); err != nil {
//		log.Print(err)
//	}
//}

func main() {
	app := app.New()
	win := app.NewWindow("Dial Habit Tracker")

	//var tabs []*container.TabItem
	//for _, color := range []Color{Red, Orange, Yellow, Green} {
	//	vbox := container.NewVBox(widget.NewLabel(string(color) + " Tasks"))
	//	tabs = append(tabs, &container.TabItem{
	//		Text:    string(color),
	//		Icon:    nil,
	//		Content: vbox,
	//	})
	//}
	//tab := container.NewAppTabs(tabs...)
	//win.SetContent(tab)
	win.SetContent(container.NewVBox(widget.NewLabel("hi")))
	win.ShowAndRun()
}
