package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//TODO 
//1. Сделать обновление главного списка ++
//2. Доделать кнопки 
//3. Сделать везде адекватный контроль ввода 
//4. Сделать подтверждение действий с помощью пароля ++
//5. Создание резервных копий базы данных ++
//6. Сводку за выбранный период ++
//7. Одинаковые ID задач ++
//8. Обновление суммы в задаче(опять исправить нужно) ++
//9. Конфиг ++


func main() {

	user, pass, dbname, err := readConfig()
	if err !=nil {
		fmt.Println(err)
	}

	err = BackupDB( user, pass, dbname, "./dbBackup", "localhost")
	if err != nil{
		fmt.Println("Error:", err)
	}

	myApp := app.New()
	mainWindow := myApp.NewWindow("SewingManager")

	mainWindow.SetMaster()

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButtonWithIcon("Работники", theme.AccountIcon(), func() {
			InitWWindow(myApp)
		}),
		widget.NewButtonWithIcon("Добавить Задачу", theme.ContentAddIcon(), func() {
			AddTask(myApp)
		}),
		widget.NewButtonWithIcon("Номенкулатура", theme.DocumentCreateIcon(), func() {
			InitNWindow(myApp)
		}),
		widget.NewButtonWithIcon("Сводка", theme.DocumentIcon(), func() {
			summaryPrep(myApp)
		}),
		widget.NewButtonWithIcon("Настройки", theme.SettingsIcon(), func() {
			// InitSWindow(myApp)
		}),
		widget.NewButtonWithIcon("Помощь", theme.QuestionIcon(), func() {
			dialog.ShowInformation("Помощь", "Для получения помощи обратитесь к разработчику в телеграмме @AmusedCreator", mainWindow)
		}),
	)

	mainWindow.Resize(fyne.NewSize(1050, 600))

	label := widget.NewLabel("Таблица задач")

	table := mainTableMaker(0, myApp)

	content := container.NewBorder(label, nil, buttons, nil, table)
	mainWindow.SetContent(content)
	mainWindow.ShowAndRun()
}



var tasksdata [][]string = nil

var mainTable *widget.Table

func mainTableMaker(tSort int, myApp fyne.App) fyne.CanvasObject {
	db, err := dbInit()
	if err != nil {
		log.Fatal(err)
	}

	UpdateDataBase(db)

	tasksdata, err = GetTasks(db, tSort)
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{"№", "id", "дата приема", "дата сдачи", "заказчик", "наименование", "кол-во", "готово"}

	mainTable = widget.NewTable(
		func() (int, int) {
			return len(tasksdata) + 1, len(headers)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(" ")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row == 0 {
				o.(*widget.Label).SetText(headers[i.Col])
			} else {
				o.(*widget.Label).SetText(tasksdata[i.Row-1][i.Col])
			}
		})

	mainTable.OnSelected = func(id widget.TableCellID) {
		if id.Row == 0 {
			if id.Col >= 0 {
				mainTableMaker(id.Col, myApp)
				mainTable.Refresh()
				mainTable.Unselect(id)
			} else {
				return
			}
		} else {
			fmt.Println(id.Col, id.Row)
			if id.Col >= 0 && id.Row > 0 {
				taskID, err := strconv.Atoi(tasksdata[id.Row-1][0])
				if err != nil {
					log.Fatal(err)
				}

				askWindow := myApp.NewWindow("")
				askWindow.Resize(fyne.NewSize(100, 50))

				custom, err := strconv.Atoi(tasksdata[id.Row-1][1])
				if err != nil {
					log.Fatal(err)
				}

				buttons := container.New(layout.NewGridLayout(1),
					widget.NewButtonWithIcon("Просмотреть", theme.ContentAddIcon(), func() {
						InitTWindow(myApp, taskID, custom)
						askWindow.Close()
					}),
					widget.NewButtonWithIcon("Удалить", theme.DeleteIcon(), func() {
						Confirm(myApp, func() {
							err := DeleteTask(db, taskID)
							if err != nil {
								dialog.ShowError(err, askWindow)
								return
							}

							updateTable(mainTableMaker(0, myApp).(*widget.Table), myApp)
							askWindow.Close()
						})
					}),
				)

				content := container.NewBorder(nil, nil, buttons, nil)
				askWindow.SetContent(content)
				askWindow.Show()

				mainTable.Unselect(id)
			} else {
				return
			}
		}
	}

	mainTable.SetColumnWidth(0, 50)
	mainTable.SetColumnWidth(1, 100)
	mainTable.SetColumnWidth(2, 100)
	mainTable.SetColumnWidth(3, 100)
	mainTable.SetColumnWidth(4, 200)
	mainTable.SetColumnWidth(5, 150)
	mainTable.SetColumnWidth(6, 80)
	mainTable.SetColumnWidth(7, 80)

	for i := 0; i < len(tasksdata)+1; i++ {
		mainTable.SetRowHeight(i, 30)
	}

	return mainTable
}

func AddTask(myApp fyne.App) {
	taskWindow := myApp.NewWindow("Добавить Задачу")
	taskWindow.Resize(fyne.NewSize(400, 400))

	db := getDB()

	nomenclature, err := GetNomenclature(db)
	if err != nil {
		dialog.ShowError(err, taskWindow)
		return
	}

	nomenclatureMap := make(map[string]string)
	var nomenclatureList []string
	for _, nom := range nomenclature {
		nomID := nom[0]
		nomName := nom[1]
		nomenclatureList = append(nomenclatureList, nomName)
		nomenclatureMap[nomName] = nomID
	}

	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("ID")

	acceptEntry := widget.NewEntry()
	acceptEntry.SetPlaceHolder("дд-мм-гггг")

	acceptEntry.OnChanged = func(text string) {
		if len(text) > 10 {
			acceptEntry.SetText(text[:10])
			acceptEntry.CursorRow = 10
		}

		if len(text) == 2 || len(text) == 5 {
			acceptEntry.SetText(text + "-")
			acceptEntry.CursorRow = len(text) + 1
		}
	}

	deliEntry := widget.NewEntry()
	deliEntry.SetPlaceHolder("дд-мм-гггг")

	deliEntry.OnChanged = func(text string) {
		if len(text) > 10 {
			deliEntry.SetText(text[:10])
			deliEntry.CursorRow = 10
		}

		if len(text) == 2 || len(text) == 5 {
			deliEntry.SetText(text + "-")
			deliEntry.CursorRow = len(text) + 1
		}
	}

	clientEntry := widget.NewEntry()
	clientEntry.SetPlaceHolder("Имя заказчика")

	nomenclatureSelect := widget.NewSelect(nomenclatureList, func(value string) {
		fmt.Println("Выбрано:", value)
	})

	quantityEntry := widget.NewEntry()
	quantityEntry.SetPlaceHolder("Количество")

	saveButton := widget.NewButton("Сохранить", func() {
		Confirm(myApp, func() {
			acceptDate, err := time.Parse("02-01-2006", acceptEntry.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("некорректная дата приема: %s", acceptEntry.Text), taskWindow)
				return
			}

			deliDate, err := time.Parse("02-01-2006", deliEntry.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("некорректная дата сдачи: %s", deliEntry.Text), taskWindow)
				return
			}

			if acceptDate.After(deliDate) {
				dialog.ShowError(fmt.Errorf("дата приема не может быть позже даты сдачи"), taskWindow)
				return
			}

			formattedAcceptDate := acceptDate.Format("01-02-2006")
			formattedDeliDate := deliDate.Format("01-02-2006")

			quantity, err := strconv.Atoi(quantityEntry.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("ошибка перевода"), taskWindow)
				return
			}

			if quantity <= 0 {
				dialog.ShowError(fmt.Errorf("некорректное количество: %s", quantityEntry.Text), taskWindow)
				return
			}

			nomID := nomenclatureMap[nomenclatureSelect.Selected]

			nomIDInt, err := strconv.Atoi(nomID)
			if err != nil {
				dialog.ShowError(err, taskWindow)
				return
			}

			err = SaveTask(db, idEntry.Text, formattedAcceptDate, formattedDeliDate, clientEntry.Text, nomIDInt, quantityEntry.Text)
			if err != nil {
				dialog.ShowError(err, taskWindow)
				return
			}

			dialog.ShowInformation("Успех", "Задача добавлена", taskWindow)
			taskWindow.Close()
			updateTable(mainTableMaker(0, myApp).(*widget.Table), myApp)
		})
	})

	form := container.NewVBox(
		widget.NewLabel("ID заказа:"),
		idEntry,
		widget.NewLabel("Дата приема:"),
		acceptEntry,
		widget.NewLabel("Дата сдачи:"),
		deliEntry,
		widget.NewLabel("Заказчик:"),
		clientEntry,
		widget.NewLabel("Наименование:"),
		nomenclatureSelect,
		widget.NewLabel("Количество:"),
		quantityEntry,
		saveButton,
	)

	taskWindow.SetContent(form)
	taskWindow.Show()
}

func summaryPrep(myApp fyne.App){
	summaryPrepWindow := myApp.NewWindow("Выберите период")
	summaryPrepWindow.Resize(fyne.NewSize(400, 400))

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

	startDate := widget.NewEntry()
	startDate.SetPlaceHolder("дд-мм-гггг")
	startDate.OnChanged = func(text string) {
		if len(text) > 10 {
			startDate.SetText(text[:10])
			startDate.CursorRow = 10
		}

		if len(text) == 2 || len(text) == 5 {
			startDate.SetText(text + "-")
			startDate.CursorRow = len(text) + 1
		}
	}

	endDate := widget.NewEntry()
	endDate.SetPlaceHolder("дд-мм-гггг")
	endDate.OnChanged = func(text string) {
		if len(text) > 10 {
			endDate.SetText(text[:10])
			endDate.CursorRow = 10
		}

		if len(text) == 2 || len(text) == 5 {
			endDate.SetText(text + "-")
			endDate.CursorRow = len(text) + 1
		}
	}

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButton("Показать", func() {
			InitSWindow(myApp, workerSelect.Selected, startDate.Text, endDate.Text)
			summaryPrepWindow.Close()
		}),
	)


	form := container.NewVBox(
		widget.NewLabel("Работник:"),
		workerSelect,
		widget.NewLabel("Дата начала:"),
		startDate,
		widget.NewLabel("Дата окончания:"),
		endDate,
		buttons,
	)


	summaryPrepWindow.SetContent(form)
	summaryPrepWindow.Show()
}

func updateTable(mainTable *widget.Table, myApp fyne.App) {
	mainTable = mainTableMaker(0, myApp).(*widget.Table)
	mainTable.Refresh()
}
