package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Task struct {
	Color Color
	Label string
	Done  bool
}

func (t *Task) Render(index int, m *MainWindow) app.UI {
	return app.Button().
		Class("btn", "btn-outline-secondary", "col-sm").
		Body(
			app.Text(t.Label),
			app.Button().
				DataSet("color", t.Color).
				DataSet("index", index).
				Class("btn", "btn-outline-danger", "float-end").
				OnClick(m.deleteTask).
				Body(
					app.If(m.Armed.C == t.Color && m.Armed.I == index, icon("delete", "Delete?")),
					app.If(m.Armed.C != t.Color || m.Armed.I != index, icon("delete", "")),
				))
}
