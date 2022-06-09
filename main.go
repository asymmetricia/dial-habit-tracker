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

var Sets = map[Color]Task{}

const dataUrl = "dialhabittracker.dat"

type MainWindow struct {
	app.Compo
	SelectedTab Color
	CanUpdate   bool
}

func (m *MainWindow) OnAppUpdate(ctx app.Context) {
	m.CanUpdate = ctx.AppUpdateAvailable()
	m.Update()
}

var _ app.AppUpdater = (*MainWindow)(nil)

func (m *MainWindow) OnMount(ctx app.Context) {
	ctx.ObserveState("selected-tab").Value(&m.SelectedTab)
	if m.SelectedTab == "" {
		m.SelectedTab = Red
	}
	println(m.SelectedTab)
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
				m.SelectedTab = Color(ctx.JSSrc().Get("dataset").Get("color").String())
				ctx.SetState("selected-tab", m.SelectedTab, app.Persist)
				m.Update()
			})

		pane := app.Div().
			Class("tab-pane", "fade").
			ID(string(color) + "-tab-pane").
			Role("tabpanel").
			Text("herein lies " + color + " data.")

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
				app.A().Class("navbar-brand").Href("#").Text("Dial Habit Tracker")),
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
				app.Div().Class("tab-content").Body(panes...)))

	return doc
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
			"/web/css/bootstrap.min.css",
		},
	})
	log.Print("Starting up!")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
