package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func InitNWindow(myApp fyne.App) {
	nWindow := myApp.NewWindow("Номенклатура")
	table := nTableMaker(myApp, nil)

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButtonWithIcon("Добавить", theme.ContentAddIcon(), func() {
			addNomenclature(myApp, nWindow, table)
		}),
		widget.NewButtonWithIcon("Назад", theme.NavigateBackIcon(), func() { nWindow.Close() }),
	)

	nWindow.Resize(fyne.NewSize(1000, 600))

	label := widget.NewLabel("Таблица номенклатуры")

	content := container.NewBorder(label, nil, buttons, nil, table)
	nWindow.SetContent(content)

	nWindow.Show()
	return
}

func nTableMaker(myApp fyne.App, table *widget.Table) *widget.Table {
	var nomenclaturedata [][]string = nil
	db := getDB()
	nomenclaturedata, err := GetNomenclature(db)
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{"id", "Название", "Цена шт."}

	if table == nil {
		table = widget.NewTable(
			func() (int, int) {
				return len(nomenclaturedata) + 1, len(headers)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("Header")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				if i.Row == 0 {
					o.(*widget.Label).SetText(headers[i.Col])
				} else {
					o.(*widget.Label).SetText(nomenclaturedata[i.Row-1][i.Col])
				}
			},
		)

		table.OnSelected = func(id widget.TableCellID) {
			if id.Row <= 0 {
				return
			}
	
			// Сначала заново получаем актуальные данные
			db := getDB()
			nomenclaturedata, err := GetNomenclature(db)
			if err != nil {
				log.Fatal(err)
			}
	
			if id.Row-1 >= len(nomenclaturedata) {
				return // Избегаем выхода за пределы массива
			}
	
			nomInfoWindow := myApp.NewWindow("Информация о номенклатуре")
			nomInfoWindow.Resize(fyne.NewSize(200, 200))
	
			nName := widget.NewEntry()
			nPrice := widget.NewEntry()
	
			nName.SetText(nomenclaturedata[id.Row-1][1])
			nPrice.SetText(nomenclaturedata[id.Row-1][2])
	
			buttons := container.New(layout.NewGridLayout(1),
				widget.NewButton("Сохранить", func() {
					db := getDB()
					_, err := db.Exec("UPDATE Nomenclature SET nom_name = ?, nom_price = ? WHERE nom_id = ?", nName.Text, nPrice.Text, nomenclaturedata[id.Row-1][0])
					if err != nil {
						log.Fatal(err)
						return
					}
					nomInfoWindow.Close()
					UpdateDataBase(db)
					updateNTable(myApp, table)
				}),
				widget.NewButton("Удалить", func() {
					db := getDB()
					_, err := db.Exec("DELETE FROM Nomenclature WHERE nom_id = ?", nomenclaturedata[id.Row-1][0])
					if err != nil {
						log.Fatal(err)
						return
					}
					nomInfoWindow.Close()
					updateNTable(myApp, table)
				}),
				widget.NewButton("Отмена", func() {
					nomInfoWindow.Close()
				}),
			)

			content := container.NewVBox(nName, nPrice, buttons)
			nomInfoWindow.SetContent(content)

			nomInfoWindow.Show()
		}

		table.SetColumnWidth(0, 30)
		table.SetColumnWidth(1, 200)
		table.SetColumnWidth(2, 50)

		for i := 0; i < len(nomenclaturedata)+1; i++ {
			table.SetRowHeight(i, 30)
		}

	} else {
		table.Length = func() (int, int) {
			return len(nomenclaturedata) + 1, len(headers)
		}

		table.UpdateCell = func(id widget.TableCellID, obj fyne.CanvasObject) {
			if id.Row == 0 {
				obj.(*widget.Label).SetText(headers[id.Col])
			} else {
				obj.(*widget.Label).SetText(nomenclaturedata[id.Row-1][id.Col])
			}
		}
		table.Refresh()
	}
	
	return table
}

func addNomenclature(myApp fyne.App, w fyne.Window, table *widget.Table) {
	addNWindow := myApp.NewWindow("Добавить номенклатуру")

	addNWindow.Resize(fyne.NewSize(200, 200))

	nName := widget.NewEntry()
	nPrice := widget.NewEntry()

	nName.SetPlaceHolder("Название")
	nPrice.SetPlaceHolder("Цена шт.")

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButton("Добавить", func() {
			if nName.Text == "" || nPrice.Text == "" {
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
			_, err := db.Exec("INSERT INTO Nomenclature (nom_name, nom_price) VALUES (?, ?)", nName.Text, nPrice.Text)
			if err != nil {
				log.Fatal(err)
				return
			}
			addNWindow.Close()
			updateNTable(myApp, table)
		}),
		widget.NewButton("Отмена", func() {
			addNWindow.Close()
		}),
	)
	
	content := container.NewVBox(nName, nPrice, buttons)
	addNWindow.SetContent(content)

	addNWindow.Show()

	return
}

func updateNTable(myApp fyne.App, table *widget.Table) {
	table = nTableMaker(myApp, table)
	table.Refresh()
}