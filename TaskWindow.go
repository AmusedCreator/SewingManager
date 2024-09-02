package main

import (
	"fmt"
	"log"
	"strconv"
	//"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var taskID int

func InitTWindow(myApp fyne.App, id int){

	db, err := dbInit()
	if err != nil {
		log.Fatal(err)
	}
	taskName := GetTaskNameByID(db, id)
	delText := "Удалить задачу " + taskName

	taskID = id
	twWindow := myApp.NewWindow("Задачи работников")
	table := tTableMaker(myApp, nil)

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButtonWithIcon("Добавить", theme.ContentAddIcon(), func() {
			addTaskWorker(myApp, twWindow, table)
		}),
		widget.NewButtonWithIcon("Назад", theme.NavigateBackIcon(), func() { twWindow.Close() }),
		widget.NewButtonWithIcon(delText, theme.DeleteIcon(), func() {
			DeleteTask(db, taskID)
			twWindow.Close()
		}),

	)

	twWindow.Resize(fyne.NewSize(1000, 600))

	label := widget.NewLabel("Task Workers Table")

	content := container.NewBorder(label, nil, buttons, nil, table)
	twWindow.SetContent(content)

	twWindow.Show()
	return
}

func tTableMaker(myApp fyne.App, table *widget.Table) *widget.Table {
	var taskworkerdata [][]string = nil
	db := getDB()
	taskworkerdata, err := GetTaskByID(db, taskID)
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{"id", "id задачи", "id работника", "кол-во", "дата", "сумма за день" }

	if table == nil {
		table = widget.NewTable(
			func() (int, int) {
				return len(taskworkerdata) + 1, len(headers)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("Header")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				if i.Row == 0 {
					o.(*widget.Label).SetText(headers[i.Col])
				} else {
					o.(*widget.Label).SetText(taskworkerdata[i.Row-1][i.Col])
				}
			},
		)
		table.SetColumnWidth(0, 50)
		table.SetColumnWidth(1, 100)
		table.SetColumnWidth(2, 100)
		table.SetColumnWidth(3, 100)
		table.SetColumnWidth(4, 100)
		table.SetColumnWidth(5, 100)

		for i := 0; i < len(taskworkerdata); i++ {
			table.SetRowHeight(i, 30)
		}
	} else {
		table.Length = func() (int, int) {
			return len(taskworkerdata) + 1, len(headers)
		}

		table.UpdateCell = func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row == 0 {
				o.(*widget.Label).SetText(headers[i.Col])
			} else {
				o.(*widget.Label).SetText(taskworkerdata[i.Row-1][i.Col])
			}
		}
		table.Refresh()
	}
	return table
}

func updateTWTable(myApp fyne.App, table *widget.Table) {
	table = tTableMaker(myApp, table)
	table.Refresh()
}

func addTaskWorker(myApp fyne.App, w fyne.Window, table *widget.Table) {
	addTWWindow := myApp.NewWindow("Добавить связь задачи и работника")

	addTWWindow.Resize(fyne.NewSize(300, 300))

	db := getDB()
	workers, err := GetWorkers(db)
	if err != nil {
		log.Fatal(err)
	}
 
	var workerNames []string
	for _, worker := range workers {
		workerNames = append(workerNames, worker[0])
	}
	workerSelect := widget.NewSelect(workerNames, func(value string) {})
	workerSelect.PlaceHolder = "Выберите работника"

	twDone := widget.NewEntry()
	twDone.SetPlaceHolder("Выполнено")

	twDate := widget.NewEntry()
	twDate.SetPlaceHolder("дд-мм-гггг")

	twDate.OnChanged = func(text string) {
		if len(text) > 10 {
			twDate.SetText(text[:10])
			twDate.CursorRow = 10
		}

		if len(text) == 2 || len(text) == 5 {
			twDate.SetText(text + "-")
			twDate.CursorRow = len(text) + 1
		}
	}

	twDaySum := widget.NewLabel("Сумма за день будет рассчитана автоматически")

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButton("Добавить", func() {
			if  workerSelect.Selected == "" || twDone.Text == "" || twDate.Text == "" {
				error := myApp.NewWindow("Ошибка!")
				error.Resize(fyne.NewSize(200, 100))
				error.SetContent(widget.NewLabel("Заполните все поля!"))
				okbutton := widget.NewButton("OK", func() {
					error.Close()
				})
				error.SetContent(container.NewVBox(widget.NewLabel("Заполните все поля!"), okbutton))
				error.Show()
				return
			}

			//twDate, err := time.Parse("02-01-2006", twDate.Text)
			if err != nil {
			dialog.ShowError(fmt.Errorf("некорректная дата приема"), addTWWindow)
			return
		}

			// Расчет суммы
			//done, _ := strconv.Atoi(twDone.Text)
			//daySum := CalculateDaySum(db, taskID, done)

			//_, err := db.Exec("INSERT INTO Task_Workers (task_id, worker_id, tw_done, tw_date, tw_day_sum) VALUES (?, ?, ?, ?, ?)",
			//	taskID, GetWorkerID(workerSelect.Selected), twDone.Text, twDate.Text, daySum)
			if err != nil {
				log.Fatal(err)
				return
			}
			addTWWindow.Close()
			updateTaskWorkersTable(myApp, table)
		}),
		widget.NewButton("Отмена", func() {
			addTWWindow.Close()
		}),
	)

	taskIDLabel := widget.NewLabel(strconv.Itoa(taskID))
	content := container.NewVBox(taskIDLabel, workerSelect, twDone, twDate, twDaySum, buttons)
	addTWWindow.SetContent(content)

	addTWWindow.Show()

	return
}

func updateTaskWorkersTable(myApp fyne.App, table *widget.Table) {
	table = tTableMaker(myApp, table)
	table.Refresh()
}


