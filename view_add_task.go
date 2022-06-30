package main

import (
	"fmt"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func (m *MainWindow) saveTask(ctx app.Context, e app.Event) {
	e.PreventDefault()

	if m.Tasks == nil {
		m.Tasks = map[int]*Task{}
	}

	desc := app.Window().GetElementByID(taskModalDescriptionId)
	newTask := desc.Get("value").String()
	newTask = strings.TrimSpace(newTask)

	var err error
	if len(newTask) == 0 {
		err = fmt.Errorf("Please provide a short task description!")
	}

	if err == nil {
		for _, t := range m.Tasks {
			if t.Color == m.SelectedTab && strings.ToLower(t.Label) == strings.ToLower(newTask) {
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

	task := &Task{Label: newTask, Color: m.SelectedTab, Id: m.NextId}
	m.Tasks[task.Id] = task
	ctx.SetState(stateKeyNextId, m.NextId+1, app.Persist)
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
