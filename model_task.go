package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"log"
	"strconv"
)

type Task struct {
	Color Color
	Label string
	Done  bool
	Id    int
}

func toggleTask(m *MainWindow) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		taskId, _ := strconv.Atoi(ctx.JSSrc().Get("dataset").Get("id").String())
		log.Printf("toggle task %d", taskId)
		m.Tasks[taskId].Done = !m.Tasks[taskId].Done
		ctx.SetState(stateKeyTaskList, m.Tasks, app.Persist)
		m.Update()
	}
}

func deleteTask(m *MainWindow) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		taskId, _ := strconv.Atoi(ctx.JSSrc().Get("dataset").Get("id").String())
		if m.Armed != taskId {
			ctx.SetState(stateKeyArmed, taskId)
			m.Update()
			return
		}

		log.Printf("delete task %d", taskId)
		delete(m.Tasks, taskId)
		ctx.SetState(stateKeyTaskList, m.Tasks, app.Persist)
		ctx.SetState(stateKeyArmed, -1)
		m.Update()
	}
}

func (t *Task) Render(m *MainWindow) app.UI {
	return app.Div().
		Class("d-flex", "justify-content-between").
		Role("group").
		Body(
			app.Button().
				DataSet("id", t.Id).
				Class("btn", "btn-outline-dark").
				Style("border-top-right-radius", "0").
				Style("border-bottom-right-radius", "0").
				Style("border-right", "none").
				Body(
					app.If(t.Done, symbol("select_check_box", "")).
						Else(symbol("check_box_outline_blank", ""))).
				OnClick(toggleTask(m)),
			app.Button().
				DataSet("id", t.Id).
				Class("btn", "flex-grow-1", "btn-outline-dark").
				Style("border-radius", "0").
				Style("border-left", "none").
				Style("border-right", "none").
				Text(t.Label).
				OnClick(toggleTask(m)),
			app.Button().
				DataSet("id", t.Id).
				Class("btn", "btn-outline-danger").
				Style("border-top-left-radius", "0").
				Style("border-bottom-left-radius", "0").
				OnClick(deleteTask(m)).
				Body(
					app.
						If(m.Armed == t.Id, icon("delete", "Delete?")).
						Else(icon("delete", "")),
				),
		)
}
