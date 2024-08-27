package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func InitWWindow(myApp fyne.App) {
	wWindow := myApp.NewWindow("Работники")
	table := wTableMaker(myApp, nil)

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButtonWithIcon("Добавить", theme.ContentAddIcon(), func() {
			addWorker(myApp, wWindow, table)
		}),
		widget.NewButtonWithIcon("Назад", theme.NavigateBackIcon(), func() { wWindow.Close() }),
	)

	wWindow.Resize(fyne.NewSize(1000, 600))

	label := widget.NewLabel("Таблица работников")

	content := container.NewBorder(label, nil, buttons, nil, table)
	wWindow.SetContent(content)

	wWindow.Show()
	return
}

func wTableMaker(myApp fyne.App, table *widget.Table) *widget.Table {
	db := getDB()
	workersdata, err := GetWorkers(db)
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{"id", "Имя", "Фамилия", "Прочее"}

	if table == nil {
		table = widget.NewTable(
			func() (int, int) {
				return len(workersdata) + 1, len(headers)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("Header")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				if i.Row == 0 {
					o.(*widget.Label).SetText(headers[i.Col])
				} else {
					o.(*widget.Label).SetText(workersdata[i.Row-1][i.Col])
				}
			},
		)

		table.OnSelected = func(id widget.TableCellID) {
			if id.Row <= 0 {
				return
			}

			// Сначала заново получаем актуальные данные
			db := getDB()
			workersdata, err := GetWorkers(db)
			if err != nil {
				log.Fatal(err)
			}

			if id.Row-1 >= len(workersdata) {
				return // Избегаем выхода за пределы массива
			}

			workerInfoWindow := myApp.NewWindow("Информация о работнике")
			workerInfoWindow.Resize(fyne.NewSize(400, 300))

			fName := widget.NewEntry()
			sName := widget.NewEntry()
			about := widget.NewEntry()

			fName.SetText(workersdata[id.Row-1][1])
			sName.SetText(workersdata[id.Row-1][2])
			about.SetText(workersdata[id.Row-1][3])

			buttons := container.New(layout.NewGridLayout(1),
				widget.NewButton("Сохранить", func() {
					db := getDB()
					_, err := db.Exec("UPDATE workers SET worker_fname = ?, worker_sname = ?, worker_about = ? WHERE worker_id = ?", fName.Text, sName.Text, about.Text, workersdata[id.Row-1][0])
					if err != nil {
						log.Fatal(err)
					}
					workerInfoWindow.Close()
					updateWTable(myApp, table)
				}),
				widget.NewButton("Удалить", func() {
					db := getDB()
					_, err := db.Exec("DELETE FROM workers WHERE worker_id = ?", workersdata[id.Row-1][0])
					if err != nil {
						log.Fatal(err)
					}
					workerInfoWindow.Close()
					updateWTable(myApp, table)
				}),
				widget.NewButton("Отмена", func() {
					workerInfoWindow.Close()
				}),
			)

			content := container.NewVBox(fName, sName, about, buttons)
			workerInfoWindow.SetContent(content)

			workerInfoWindow.Show()
		}

		table.SetColumnWidth(0, 30)
		table.SetColumnWidth(1, 150)
		table.SetColumnWidth(2, 150)
		table.SetColumnWidth(3, 300)

		for i := 0; i < len(workersdata)+1; i++ {
			table.SetRowHeight(i, 30)
		}

	} else {
		table.Length = func() (int, int) {
			return len(workersdata) + 1, len(headers)
		}

		table.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
			if id.Row == 0 {
				obj.(*widget.Label).SetText(headers[id.Col])
			} else {
				obj.(*widget.Label).SetText(workersdata[id.Row-1][id.Col])
			}
		}
		table.Refresh()
	}

	return table
}

func addWorker(myApp fyne.App, wWindow fyne.Window, table *widget.Table) {
	addWWindow := myApp.NewWindow("Добавление работника")

	addWWindow.Resize(fyne.NewSize(400, 300))

	fName := widget.NewEntry()
	sName := widget.NewEntry()
	about := widget.NewEntry()

	errorLabel := canvas.NewText("", color.RGBA{255, 0, 0, 255})

	fName.SetPlaceHolder("Имя")
	sName.SetPlaceHolder("Фамилия")
	about.SetPlaceHolder("Прочее")

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButton("Добавить", func() {
			if fName.Text == "" || sName.Text == "" {
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
			db := getDB()
			_, err := db.Exec("INSERT INTO workers (worker_fname, worker_sname, worker_about) VALUES (?, ?, ?)", fName.Text, sName.Text, about.Text)
			if err != nil {
				errorLabel.Text = "Ошибка добавления"
				return
			}
			addWWindow.Close()
			updateWTable(myApp, table)
		}),
		widget.NewButton("Отмена", func() {
			addWWindow.Close()
		}),
		errorLabel,
	)

	content := container.NewVBox(fName, sName, about, buttons)
	addWWindow.SetContent(content)

	addWWindow.Show()

	return
}

func updateWTable(myApp fyne.App, table *widget.Table) {
	wTableMaker(myApp, table)
	table.Refresh() 
}
