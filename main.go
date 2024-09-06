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

func main() {

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
		widget.NewButtonWithIcon("Настройки", theme.SettingsIcon(), func() {
			// InitSWindow(myApp)
		}),
		widget.NewButtonWithIcon("Помощь", theme.QuestionIcon(), func() {
			// InitHWindow(myApp)
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

	headers := []string{"id", "дата приема", "дата сдачи", "заказчик", "наименование", "кол-во", "готово"}

	table := widget.NewTable(
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

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row == 0 {
			if id.Col >= 0 {
				mainTableMaker(id.Col, myApp)
				table.Refresh()
				table.Unselect(id)
			} else {
				return
			}
		} else {
			if id.Col >= 0 || id.Row > 0 {
				taskID, err := strconv.Atoi(tasksdata[id.Row-1][0])
				if err != nil {
					log.Fatal(err)
				}
				//создай окно в котором будет выбор, просмотреть информацию о задаче или удалить задачу
				askWindow := myApp.NewWindow("")
				askWindow.Resize(fyne.NewSize(100, 50))

				buttons := container.New(layout.NewGridLayout(1),
					widget.NewButtonWithIcon("Просмотреть", theme.ContentAddIcon(), func() {
						InitTWindow(myApp, taskID)
						defer updateTable(mainTableMaker(0, myApp).(*widget.Table), myApp)
						
						askWindow.Close()
					}),
					widget.NewButtonWithIcon("Удалить", theme.DeleteIcon(), func() {
						err := DeleteTask(db, taskID)
						if err != nil {
							dialog.ShowError(err, askWindow)
							return
						}

						updateTable(mainTableMaker(0, myApp).(*widget.Table), myApp)
						askWindow.Close()
					}),
				)

				content := container.NewBorder(nil, nil, buttons, nil)
				askWindow.SetContent(content)
				askWindow.Show()

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

	for i := 0; i < len(tasksdata)+1; i++ {
		table.SetRowHeight(i, 30)
	}

	return table
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

func updateTable(table *widget.Table, myApp fyne.App) {
	table = mainTableMaker(0, myApp).(*widget.Table)
	table.Refresh()
}
