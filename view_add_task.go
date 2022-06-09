package main

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func (m *MainWindow) deleteTask(ctx app.Context, e app.Event) {
	color := Color(ctx.JSSrc().Get("dataset").Get("color").String())
	indexStr := ctx.JSSrc().Get("dataset").Get("index").String()

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return
	}

	if m.Armed.C != color || m.Armed.I != index {
		m.Armed.C = color
		m.Armed.I = index
		ctx.SetState(stateKeyArmed, m.Armed)
	} else {
		m.Tasks[color] = append(m.Tasks[color][:index], m.Tasks[color][index+1:]...)
		ctx.SetState(stateKeyTaskList, m.Tasks, app.Persist)
		m.Armed.C = ""
		ctx.SetState(stateKeyArmed, m.Armed)
	}

	m.Update()
}

func (m *MainWindow) saveTask(ctx app.Context, e app.Event) {
	if m.Tasks == nil {
		m.Tasks = map[Color][]Task{}
	}

	m.Tasks[m.SelectedTab] = append(m.Tasks[m.SelectedTab],
		Task{Label: app.Window().GetElementByID("addTaskDescription").Get("value").String()})
	ctx.SetState(stateKeyTaskList, m.Tasks, app.Persist)
	e.PreventDefault()
	m.Update()
}

func (m *MainWindow) taskModal() app.UI {
	return app.Div().Class("modal", "fade").
		ID("taskModal").
		TabIndex(-1).
		Role("dialog").
		Aria("labelledby", "taskModalLabel").
		Aria("hidden", "true").Body(
		app.Div().Class("modal-dialog").Role("document").Body(
			app.Div().Class("modal-content").Body(
				app.Div().Class("modal-header").Body(
					app.H5().Class("modal-title").ID("taskModalLabel").Text("Add Task"),
					app.Button().
						Class("btn", "close").
						DataSet("bs-dismiss", "modal").
						Aria("label", "close").
						Body(icon("close", ""))),
				app.Div().Class("modal-body").Body(
					app.Form().OnSubmit(m.saveTask).Body(
						app.Label().
							Class("form-label").
							For("addTaskDescription").
							Text("Task Description"),
						app.Input().
							Class("form-control").
							ID("addTaskDescription").
							Placeholder("What would you like to work on?")),
					app.Button().
						Class("btn", "btn-success", "mt-2").
						DataSet("bs-dismiss", "modal").
						DataSet("bs-target", "#taskModal").
						Body(icon("done", "Save Task")).
						OnClick(m.saveTask)))))
}
