package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
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

const dataUrl = "dialhabittracker.dat"

type MainWindow struct {
	app.Compo
	SelectedTab Color
	CanUpdate   bool
	Tasks       map[Color][]Task
	Armed       struct {
		C Color
		I int
	}
}

func (m *MainWindow) OnAppUpdate(ctx app.Context) {
	m.CanUpdate = ctx.AppUpdateAvailable()
	m.Update()
}

var _ app.AppUpdater = (*MainWindow)(nil)

const stateKeySelectedTab = "selected-tab"
const stateKeyTaskList = "task-list"
const stateKeyArmed = "armed"

func (m *MainWindow) OnMount(ctx app.Context) {
	ctx.ObserveState(stateKeySelectedTab).Value(&m.SelectedTab)
	if m.SelectedTab == "" {
		m.SelectedTab = Red
	}
	ctx.ObserveState(stateKeyTaskList).Value(&m.Tasks)
	ctx.ObserveState(stateKeyArmed).Value(&m.Armed)
}

func (m *MainWindow) Render() app.UI {
	println("rendering")
	var tabs []app.UI
	var panes []app.UI
	for _, color := range []Color{Red, Orange, Yellow, Green} {
		tab := app.A().Class("nav-link").Href("#")

		tab = tab.
			Text(color).
			DataSet("color", color).
			DataSet("bs-toggle", "tab").
			DataSet("bs-target", "#"+string(color)+"-tab-pane").
			OnClick(func(ctx app.Context, e app.Event) {
				ctx.SetState(
					stateKeySelectedTab,
					Color(ctx.JSSrc().Get("dataset").Get("color").String()),
					app.Persist)
				m.Update()
			})

		pane := app.Div().Class("d-grid", "gap-2", "mt-1")

		var tasks []app.UI

		for i, task := range m.Tasks[color] {
			tasks = append(tasks, app.Button().
				Class("btn", "btn-outline-secondary", "col-sm").
				Body(
					app.Text(task.Label),
					app.Button().
						DataSet("color", color).
						DataSet("index", i).
						Class("btn", "btn-outline-danger", "float-end").
						OnClick(m.deleteTask).
						Body(
							app.If(m.Armed.C == color && m.Armed.I == i, icon("delete", "Delete?")),
							app.If(m.Armed.C != color || m.Armed.I != i, icon("delete", "")),
						)))
		}

		tasks = append(tasks,
			app.Button().
				Type("button").
				Class("btn", "btn-primary", "col-sm").
				Body(icon("add_task", "Add New Task")).
				DataSet("bs-toggle", "modal").
				DataSet("bs-target", "#taskModal"))

		pane = pane.Body(tasks...)

		pane = app.Div().
			Class("tab-pane", "fade").
			ID(string(color) + "-tab-pane").
			Role("tabpanel").
			Body(pane)

		if m.SelectedTab == color {
			tab = tab.Class("active")
			pane = pane.Class("show", "active")
		}

		tabs = append(tabs, app.Li().Class("nav-item").Body(tab))
		panes = append(panes, pane)
	}

	doc := app.Div().
		Body(
			app.Nav().Class("navbar", "navbar-expand-lg", "bg-light").Body(
				app.Div().Class("container-xxl").Body(
					app.A().Class("navbar-brand").Href("#").Body(
						icon("done_all", "Dial Habit Tracker")))),
			app.Div().Class("container").Body(
				app.If(m.CanUpdate, app.Div().Class("alert", "alert-warning").Role("alert").Body(
					app.Span().
						Text("Update available!").
						Style("margin-right", "1em"),
					app.Button().
						Class("btn", "btn-success").
						Text("Reload").
						OnClick(func(ctx app.Context, e app.Event) { ctx.Reload() }))),
				app.Ul().Class("nav", "justify-content-center", "nav-tabs").Body(tabs...),
				app.Div().Class("tab-content").Body(panes...)),
			m.taskModal())

	return doc
}

func icon(name string, label string) app.UI {
	container := app.Div().Class("d-flex", "align-items-center")
	return container.Body(
		app.Span().
			Class("material-icons").
			Text(name),
		app.If(label != "", app.Span().Style("margin-left", ".5em").Text(label)))
}

func main() {
	app.Route("/", &MainWindow{})
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:            "Dial Habit Tracker",
		ShortName:       "DialHabit",
		Author:          "Tricia Bogen (@tricia@tech.lgbt)",
		BackgroundColor: "#000000",
		Icon: app.Icon{
			Default: "/web/Icon.png",
		},
		Scripts: []string{
			"/web/js/bootstrap.bundle.min.js",
		},
		Styles: []string{
			"https://fonts.googleapis.com/icon?family=Material+Icons",
			"/web/css/bootstrap.min.css",
		},
	})
	log.Print("Starting up!")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
