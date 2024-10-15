package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var taskID int
var custom int

func InitTWindow(myApp fyne.App, id int, cust int) {
	taskID = id
	custom = cust

	twWindow := myApp.NewWindow("Задачи работников")
	table := tTableMaker(myApp, nil)

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButtonWithIcon("Добавить", theme.ContentAddIcon(), func() {
			addTaskWorker(myApp, twWindow, table)
		}),
		widget.NewButtonWithIcon("Назад", theme.NavigateBackIcon(), func() { twWindow.Close() }),
	)

	twWindow.Resize(fyne.NewSize(1000, 600))

	taskName := GetTaskNameByID(getDB(), taskID)

	label := widget.NewLabel(taskName + " (id: " + strconv.Itoa(custom) + ")")

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

	headers := []string{"id", "id задачи", "Работник", "кол-во", "дата", "сумма за день" }

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

		table.OnSelected = func(id widget.TableCellID){
			if id.Row <= 0{return}
			tasksdata, err := GetTaskByID(getDB(), taskID)
			if err != nil {
				log.Fatal(err)
			}

			if id.Row -1 >= len(tasksdata){return}

			workersTaskInfoWindow := myApp.NewWindow("Информация о задаче работника")
			workersTaskInfoWindow.Resize(fyne.NewSize(400, 300))

			workerName := widget.NewEntry()
			workerName.SetText(tasksdata[id.Row-1][2])
			workerName.Disable()

			workerDone := widget.NewEntry()
			workerDone.SetText(tasksdata[id.Row-1][3])
			
			

			workerDate := widget.NewEntry()
			workerDate.SetText(tasksdata[id.Row-1][4])

			

			
			buttons := container.New(layout.NewGridLayout(1),
				widget.NewButton("Сохранить", func() {
					Confirm(myApp, func() {
						wDone, err := strconv.Atoi(workerDone.Text)
						if err != nil{
							log.Fatal(err)
							return
						}
						workerSum := CalculateDaySum(db, taskID, wDone)

						_, terr := db.Exec("UPDATE Task_Workers JOIN Workers ON Task_Workers.worker_id = Workers.worker_id SET tw_done = ?, tw_date = ?, tw_day_sum = ? WHERE task_worker_id = ?", workerDone.Text, workerDate.Text, workerSum, tasksdata[id.Row-1][0])
						if terr != nil {
							log.Fatal(err)
						}
						workersTaskInfoWindow.Close()
						UpdateDataBase(db)
						updateTWTable(myApp, table)
						updateTable(mainTableMaker(0, myApp).(*widget.Table), myApp)
					})
				}), 
				widget.NewButton("Удалить", func() {
					Confirm(myApp, func() {
						_, err := db.Exec("DELETE Task_Workers FROM Task_Workers JOIN Workers ON Task_Workers.worker_id = Workers.worker_id WHERE task_worker_id = ?", taskworkerdata[id.Row-1][0])
						if err != nil{
							log.Fatal(err)
						}
						workersTaskInfoWindow.Close()
						UpdateDataBase(db)
						updateTWTable(myApp, table)
						updateTable(mainTableMaker(0, myApp).(*widget.Table), myApp)
					})
				}),
				widget.NewButton("Отмена", func() {
					workersTaskInfoWindow.Close()
				}),
			)
			
			twDaySum := widget.NewLabel("Сумма за день будет рассчитана автоматически")

			content := container.NewVBox(workerName, workerDone, workerDate, twDaySum, buttons)
			workersTaskInfoWindow.SetContent(content)

			workersTaskInfoWindow.Show()
		}

		table.SetColumnWidth(0, 30)
		table.SetColumnWidth(1, 80)
		table.SetColumnWidth(2, 150)
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
		workerNames = append(workerNames, worker[1] + " " + worker[2])
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
			Confirm(myApp, func() {
				{if  workerSelect.Selected == "" || twDone.Text == "" || twDate.Text == "" {
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

				parsedDate, err := time.Parse("02-01-2006", twDate.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("некорректная дата приема"), addTWWindow)
					return
				}
				twDate.SetText(parsedDate.Format("02-01-2006"))
				
				done, _ := strconv.Atoi(twDone.Text)
				daySum := CalculateDaySum(db, taskID, done)

				workerID := GetWorkerID(workerSelect.Selected)
				_, err = AddTaskWorker(db, taskID, workerID, done, parsedDate.Format("02-01-2006"), daySum)

				if err != nil {
					log.Fatal(err)
					return
				}
				addTWWindow.Close()
				updateTWTable(myApp, table)
				UpdateDataBase(db)
				updateTable(mainTableMaker(0, myApp).(*widget.Table), myApp)}
			})
		}),
		widget.NewButton("Отмена", func() {
			addTWWindow.Close()
		}),
	)

	taskName := widget.NewLabel(GetTaskNameByID(db, taskID) + " (id: " + strconv.Itoa(custom) + ")")
	content := container.NewVBox(taskName, workerSelect, twDone, twDate, twDaySum, buttons)
	addTWWindow.SetContent(content)

	addTWWindow.Show()

	return
}


