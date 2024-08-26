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

	table := wTableMaker()

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButtonWithIcon("Добавить", theme.ContentAddIcon(), func() {
			addWorker(myApp)
		}),
		widget.NewButtonWithIcon("Назад", theme.NavigateBackIcon(), func() {wWindow.Close()}),
	)

	wWindow.Resize(fyne.NewSize(1000, 600))

	lable := widget.NewLabel("Таблица работников")

	content := container.NewBorder(lable, nil, buttons, nil, table)
	wWindow.SetContent(content)

	wWindow.Show()
	return
}

func wTableMaker() fyne.CanvasObject {
	headers := []string{"id", "Имя", "Фамилия", "Прочее"}
	
	db := getDB()

	rows, err := db.Query("SELECT worker_id, worker_fname, worker_sname, worker_about FROM workers")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var data [][]string

	for rows.Next() {
		var id, firstName, lastName, otherInfo string
		err := rows.Scan(&id, &firstName, &lastName, &otherInfo)
		if err != nil {
			log.Fatal(err)		
		}

		data = append(data, []string{id, firstName, lastName, otherInfo})
	}

	table := widget.NewTable(
		func() (int, int) {
			return len(data) + 1, len(headers)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Header")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row == 0 {
				o.(*widget.Label).SetText(headers[i.Col])
			} else {
				o.(*widget.Label).SetText(data[i.Row-1][i.Col])
			}
		},
	)

	table.SetColumnWidth(0, 30)
	table.SetColumnWidth(1, 150)
	table.SetColumnWidth(2, 150)
	table.SetColumnWidth(3, 300)
	

	for i := 0; i < len(data)+1; i++ {
		table.SetRowHeight(i, 30)
	}

	return table
}

func addWorker( myApp fyne.App) {
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
			InitWWindow(myApp)
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

