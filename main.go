package main

import (
	"database/sql"

	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// TODO Переделать базу данных - добавить новую сущность ТОВАР

func main() {

	myApp := app.New()
	myWindow := myApp.NewWindow("SewingManager")

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButtonWithIcon("Работники", theme.AccountIcon(), nil),
		widget.NewButtonWithIcon("Добавить Задачу", theme.ContentAddIcon(), nil),
		widget.NewButtonWithIcon("Номенкулатура", theme.DocumentCreateIcon(), nil),
		widget.NewButtonWithIcon("Настройки", theme.SettingsIcon(), nil),
		widget.NewButtonWithIcon("Помощь", theme.QuestionIcon(), nil),
	)

	myWindow.Resize(fyne.NewSize(1050, 600))

	lable := widget.NewLabel("Таблица задач")

	table := tableMaker()

	content := container.NewBorder(lable, nil, buttons, nil, table)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func tableMaker() fyne.CanvasObject {
	dsn := "root:123456789@/SewingDB"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	headers := []string{"id", "дата приема", "дата сдачи", "заказчик", "наименование", "кол-во", "готово", "сумма"}

	data, err := GetTasks(db)
	if err != nil {
		log.Fatal(err)
	}

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
			return
		} else {
			fmt.Println("Selected", id)
		}
	}

	table.SetColumnWidth(0, 30)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 200)
	table.SetColumnWidth(5, 100)
	table.SetColumnWidth(6, 100)
	table.SetColumnWidth(7, 100)

	for i := 0; i < len(data)+1; i++ {
		table.SetRowHeight(i, 30)
	}

	return table
}
