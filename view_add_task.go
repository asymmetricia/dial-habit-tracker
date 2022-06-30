package main

import (
	"fmt"
	"strconv"
	"strings"

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
	e.PreventDefault()

	if m.Tasks == nil {
		m.Tasks = map[Color][]Task{}
	}

	desc := app.Window().GetElementByID(taskModalDescriptionId)
	newTask := desc.Get("value").String()
	newTask = strings.TrimSpace(newTask)

	var err error
	if len(newTask) == 0 {
		err = fmt.Errorf("Please provide a short task description!")
	}

	if err == nil {
		for _, t := range m.Tasks[m.SelectedTab] {
			if strings.ToLower(t.Label) == strings.ToLower(newTask) {
				err = fmt.Errorf("Please provide a task description that is unique!")
				break
			}
		}
	}

	if err != nil {
		desc.Call("setCustomValidity", err.Error())
		app.Window().GetElementByID(taskModalValidationErrorId).Set("innerText", err.Error())
		app.Window().GetElementByID(taskModalFormId).Get("classList").Call("add", "was-validated")
		return
	}

	desc.Call("setCustomValidity", "")
	desc.Set("value", "")
	app.Window().GetElementByID(taskModalFormId).Get("classList").Call("remove", "was-validated")

	m.Tasks[m.SelectedTab] = append(m.Tasks[m.SelectedTab], Task{Label: newTask, Color: m.SelectedTab})
	ctx.SetState(stateKeyTaskList, m.Tasks, app.Persist)
	m.Update()
	app.Window().Get("bootstrap").Get("Modal").Call("getOrCreateInstance", "#"+taskModalId).Call("hide")
}

const taskModalId = "task-modal"
const taskModalDescriptionId = "task-modal-description"
const taskModalFormId = "task-modal-form"
const taskModalValidationErrorId = "task-modal-validation-error"

func (m *MainWindow) taskModal() app.UI {
	return app.Div().Class("modal", "fade").
		ID(taskModalId).
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
					app.Form().
						ID(taskModalFormId).
						OnSubmit(m.saveTask).
						Body(
							app.Label().
								Class("form-label").
								For(taskModalDescriptionId).
								Text("Task Description"),
							app.Input().
								Class("form-control").
								ID(taskModalDescriptionId).
								Placeholder("What would you like to work on?").
								Required(true),
							app.Div().
								Class("invalid-feedback").
								ID(taskModalValidationErrorId).
								Text("Please provide a short description of your new task!"),
							app.Button().
								Class("btn", "btn-success", "mt-2").
								Body(icon("done", "Save Task")).
								OnClick(m.saveTask))))))
}
