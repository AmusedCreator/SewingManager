package main

import (
	"log"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func InitSWindow(myApp fyne.App, name string, startDate string, endDate string){
	summaryWindow := myApp.NewWindow("Сводка")
	table := sTableMaker(myApp, nil, name, startDate, endDate)

	summaryWindow.Resize(fyne.NewSize(1000, 600))

	label := widget.NewLabel("Сводка по работнику: " + name + " c " + startDate + " по " + endDate + "\nбыло заработано: " + GetSummarySum(getDB(), name, startDate, endDate))

	content := container.NewBorder(label, nil, nil, nil, table)
	summaryWindow.SetContent(content)

	summaryWindow.Show()
	return
}

func sTableMaker(myApp fyne.App, table *widget.Table, name string, startDate string, endDate string) *widget.Table {
	var summarydata [][]string = nil
	db := getDB()
	summarydata, err := GetSummary(db, name, startDate, endDate)
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{"id", "id задачи", "Работник", "кол-во", "дата", "сумма за день" }

	if table == nil {
		table = widget.NewTable(
			func() (int, int) {
				return len(summarydata) + 1, len(headers)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("Header")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				if i.Row == 0 {
					o.(*widget.Label).SetText(headers[i.Col])
				} else {
					o.(*widget.Label).SetText(summarydata[i.Row-1][i.Col])
				}
			},
		)

		table.SetColumnWidth(0, 30)
		table.SetColumnWidth(1, 80)
		table.SetColumnWidth(2, 150)
		table.SetColumnWidth(3, 100)
		table.SetColumnWidth(4, 100)
		table.SetColumnWidth(5, 100)

		for i := 0; i < len(summarydata); i++ {
			table.SetRowHeight(i, 30)
		}
	}else {
		table.Length = func() (int, int) {
			return len(summarydata) + 1, len(headers)
		}

	}
	return table
}

