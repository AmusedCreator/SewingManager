package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {

	myApp := app.New()
	mainWindow := myApp.NewWindow("SewingManager")

	mainWindow.SetMaster()

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButtonWithIcon("Работники", theme.AccountIcon(), func() {
			InitWWindow(myApp)
		}),
		widget.NewButtonWithIcon("Добавить Задачу", theme.ContentAddIcon(), nil),
		widget.NewButtonWithIcon("Номенкулатура", theme.DocumentCreateIcon(), nil),
		widget.NewButtonWithIcon("Настройки", theme.SettingsIcon(), nil),
		widget.NewButtonWithIcon("Помощь", theme.QuestionIcon(), nil),
	)

	mainWindow.Resize(fyne.NewSize(1050, 600))

	lable := widget.NewLabel("Таблица задач")

	table := mainTableMaker(0)

	content := container.NewBorder(lable, nil, buttons, nil, table)
	mainWindow.SetContent(content)
	mainWindow.ShowAndRun()
}

var data [][]string = nil

func mainTableMaker(tSort int) fyne.CanvasObject {
	db, err := dbInit()
	if err != nil {
		log.Fatal(err)
	}

	data, err = GetTasks(db, tSort)
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{"id", "дата приема", "дата сдачи", "заказчик", "наименование", "кол-во", "готово"}

	table := widget.NewTable(
		func() (int, int) {
			return len(data) + 1, len(headers)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(" ")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row == 0 {
				o.(*widget.Label).SetText(headers[i.Col])
			} else {
				o.(*widget.Label).SetText(data[i.Row-1][i.Col])
			}
		})
	table.OnSelected = func(id widget.TableCellID) {
		if id.Row == 0 {
			if id.Col >= 0 {
				mainTableMaker(id.Col)
				table.Refresh()
				fmt.Println("Selected", id.Col)
				table.Unselect(id)
			} else {
				return
			}
		} else {
			if id.Col >= 0 || id.Row > 0 {
				fmt.Println("Selected", id)
				table.Unselect(id)
			}

		}
	}

	table.SetColumnWidth(0, 30)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 200)
	table.SetColumnWidth(4, 200)
	table.SetColumnWidth(5, 100)
	table.SetColumnWidth(6, 100)

	for i := 0; i < len(data)+1; i++ {
		table.SetRowHeight(i, 30)
	}

	return table
}
