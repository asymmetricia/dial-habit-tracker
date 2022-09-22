package main

import (
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Color string

var (
	Red    Color = "Red"
	Orange Color = "Orange"
	Yellow Color = "Yellow"
	Green  Color = "Green"
)

type MainWindow struct {
	app.Compo
	SelectedTab Color
	CanUpdate   bool
	Tasks       map[int]*Task
	Armed       int
	LastVisit   time.Time
	NextId      int
}

func (m *MainWindow) OnAppUpdate(ctx app.Context) {
	m.CanUpdate = ctx.AppUpdateAvailable()
	m.Update()
}

var _ app.AppUpdater = (*MainWindow)(nil)

const stateKeySelectedTab = "selected-tab"
const stateKeyTaskList = "task-list"
const stateKeyArmed = "armed"
const stateKeyLastVisit = "last-visit"
const stateKeyNextId = "next-id"

func (m *MainWindow) OnMount(ctx app.Context) {
	ctx.ObserveState(stateKeySelectedTab).Value(&m.SelectedTab)
	if m.SelectedTab == "" {
		ctx.SetState(stateKeySelectedTab, Red, app.Persist)
	}

	ctx.ObserveState(stateKeyTaskList).Value(&m.Tasks)
	ctx.ObserveState(stateKeyArmed).Value(&m.Armed)
	ctx.ObserveState(stateKeyLastVisit).Value(&m.LastVisit)
	ctx.ObserveState(stateKeyNextId).Value(&m.NextId)

	if m.NextId == 0 {
		ctx.SetState(stateKeyNextId, 1, app.Persist)
	}

	log.Printf("welcome back! It's been %s since your last visit", time.Since(m.LastVisit).String())

	taskModal := app.Window().GetElementByID(taskModalId)
	if !taskModal.IsNull() {
		taskModal.Call(
			"addEventListener",
			"shown.bs.modal",
			app.FuncOf(func(this app.Value, args []app.Value) interface{} {
				app.Window().GetElementByID(taskModalDescriptionId).Call("focus")
				return nil
			}))
	}
}

func (m *MainWindow) Render() app.UI {
	return app.If(
		m.LastVisit.IsZero() || time.Since(m.LastVisit).Hours() >= 23, m.GoodMorning()).
		Else(m.MainActivity())
}

func (m *MainWindow) SelectMood(mood string) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		ctx.SetState(stateKeySelectedTab, Color(mood), app.Persist)
		ctx.SetState(stateKeyLastVisit, time.Now(), app.Persist)
		for _, task := range m.Tasks {
			task.Done = false
		}
		ctx.SetState(stateKeyTaskList, m.Tasks, app.Persist)
		m.Update()
	}
}

func (m *MainWindow) ResetDay() app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		ctx.SetState(stateKeyLastVisit, time.Time{}, app.Persist)
		m.Update()
	}
}

func (m *MainWindow) GoodMorning() app.UI {
	var greeting app.UI
	if time.Now().Hour() < 12 {
		greeting = symbol("clear_day", "Good morning! You look nice today.")
	} else if time.Now().Hour() < 17 {
		greeting = symbol("schedule", "Good afternoon! You look nice today.")
	} else {
		greeting = symbol("nights_stay", "Good evening! You look nice today.")
	}

	colors := [][3]string{
		{"Red", "thunderstorm", "btn-danger"},
		{"Orange", "rainy", "btn-orange"},
		{"Yellow", "partly_cloudy_day", "btn-warning"},
		{"Green", "sunny", "btn-success"},
	}

	moodBtn := func(offset int) func(int) app.UI {
		return func(i int) app.UI {
			i += offset
			return app.Button().
				Class(colors[i][2], "btn", "p-4", "m-2").
				Body(
					app.Div().Class("text-center").
						Text(colors[i][0]),
					app.Div().Class("material-symbols-round").
						Style("font-size", "4em").
						Text(colors[i][1]),
				).OnClick(m.SelectMood(colors[i][0]))
		}
	}

	return m.withPreamble(
		app.Div().Class("d-flex", "justify-content-center").Body(greeting),
		app.H2().Class("text-center").
			Text("How are you feeling?"),
		app.Div().Class("d-flex", "flex-wrap", "justify-content-center", "p-2").Body(
			app.Div().Body(app.Range(colors[0:2]).Slice(moodBtn(0))),
			app.Div().Body(app.Range(colors[2:]).Slice(moodBtn(2))),
		))
}

func (m *MainWindow) MainActivity() app.UI {
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

		var taskIds []int
		for _, task := range m.Tasks {
			if task.Color == color {
				taskIds = append(taskIds, task.Id)
			}
		}
		sort.Ints(taskIds)

		var tasks []app.UI
		for _, taskId := range taskIds {
			tasks = append(tasks, m.Tasks[taskId].Render(m))
		}

		tasks = append(tasks,
			app.Button().
				Type("button").
				Class("btn", "btn-primary", "col-sm").
				Body(icon("add_task", "Add New Task")).
				DataSet("bs-toggle", "modal").
				DataSet("bs-target", "#"+taskModalId))

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

	doc := m.withPreamble(
		app.Div().Class("container-xxl").Body(
			app.Ul().Class("nav", "justify-content-center", "nav-tabs").Body(tabs...),
			app.Div().Class("tab-content").Body(panes...),
			m.taskModal(),
			m.DebugMenu()))

	return doc
}

func (m *MainWindow) DebugMenu() app.UI {
	return app.Div().
		Class("btn-group", "dropup", "fixed-bottom", "float-right").
		Body(
			app.Button().
				Class("btn", "dropdown-toggle").
				DataSet("bs-toggle", "dropdown").
				Aria("expanded", "false").
				Body(icon("bug_report")),
			app.Div().
				Class("dropdown-menu").
				Body(
					app.Div().Class("dropdown-item").Text("test")))
}

func (m *MainWindow) withPreamble(body ...app.UI) app.UI {
	return app.Div().
		Body(
			append([]app.UI{
				app.Nav().Class("navbar", "navbar-expand-sm", "bg-light").Body(
					app.Div().Class("container-xxl").Body(
						app.A().Href("#").Class("navbar-brand", "me-auto").Body(
							icon("done_all", "Dial Habit Tracker")),
						app.Button().
							Class("navbar-toggler").
							DataSet("bs-toggle", "collapse").
							DataSet("bs-target", "#navbarCollapse").
							Aria("controls", "navbarCollapse").
							Aria("expanded", "false").
							Aria("label", "toggle navigation").
							Body(icon("more")),
						app.Div().Class("collapse", "navbar-collapse").ID("navbarCollapse").Body(
							app.Ul().Class("navbar-nav", "ms-auto").Body(
								app.Li().Class("nav-item", "dropdown").Body(
									app.A().Href("#").
										Class("nav-link", "dropdown-toggle").
										DataSet("bs-toggle", "dropdown").
										Aria("expanded", "false").
										Text("Debug"),
									app.Ul().Class("dropdown-menu", "dropdown-menu-end").Body(
										app.Li().Body(
											app.A().Href("#").
												Class("dropdown-item").
												Text("Reset Day"))).OnClick(m.ResetDay()),
								))),
					)),
				app.If(m.CanUpdate,
					app.Div().Class("container").Body(
						app.Div().Class("alert", "alert-warning", "row", "justify-content-md-center", "align-items-center").
							Role("alert").
							Body(
								app.Div().Class("col-md-2").Body(
									symbol("update", "Update available!"),
								), app.Div().Class("col-md-2").Body(
									app.Button().
										Class("btn", "btn-success").
										Text("Reload").
										OnClick(func(ctx app.Context, e app.Event) { ctx.Reload() }))))),
			},
				body...)...)
}

func symbol(name string, label string) app.UI {
	if label == "" {
		return app.Span().Class("material-symbols-round").Text(name)
	}
	return app.Div().Class("d-flex", "align-items-center").Body(
		app.Span().
			Class("material-icons-round").
			Text(name),
		app.Span().
			Style("margin-left", ".5em").
			Text(label))
}

// icon generates an [app.UI] element for the [Material Icon] named by `name`, with
// an optional label `label`.
//
// [Material Icon]: https://fonts.google.com/icons
func icon(name string, label ...string) app.UI {
	if len(label) == 0 {
		return app.Span().Class("material-icons-round").Text(name)
	}
	return app.Div().Class("d-flex", "align-items-center").Body(
		app.Span().
			Class("material-icons-round").
			Text(name),
		app.Span().
			Style("margin-left", ".5em").
			Text(label[0]))
}

func main() {
	app.Route("/", &MainWindow{})
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Name:            "Dial Habit Tracker",
		ShortName:       "DialHabit",
		Author:          "Tricia Bogen (@tricia@tech.lgbt)",
		Title:           "Dial Habit Tracker is a tool for tracking habits using the Dial Method",
		BackgroundColor: "#000000",
		Icon: app.Icon{
			Default: "/web/Icon.png",
		},
		Scripts: []string{
			"/web/js/bootstrap.bundle.min.js",
		},
		Styles: []string{
			"/web/css/dial.css",
			"/web/css/bootstrap.min.css",
		},
	})
	log.Print("Starting up! Listening on port 8000.")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
