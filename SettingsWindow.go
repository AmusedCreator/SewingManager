package main

import (
	"fmt"
	"io"
	"strings"

	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"path/filepath"
)

var path string

func InitSettings(myApp fyne.App){
	settingsWindow := myApp.NewWindow("Настройки")
	
	settingsWindow.Resize(fyne.NewSize(600, 600))

	backupList := widget.NewList(
		func() int {
			return len(getBackups())
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Item")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(getBackups()[i])
		},
	)

	backupList.OnSelected = func(id widget.ListItemID) {
		path = getBackups()[id]
	}
	
	label := widget.NewLabel("Настройки")

	buttons := container.New(layout.NewGridLayout(1),
		widget.NewButtonWithIcon("Установить выбранную", nil, func() {
			Confirm(myApp, func() {
				ClearDataBase(getDB())
				setBackup(path)
				dialog.ShowInformation("База данных успешно установлена", "Для корректного отображения данных перезагрузите приложение", settingsWindow)
			})
		}),
		widget.NewButtonWithIcon("Очистить базу данных", nil, func() {
			db := getDB()
			ClearDataBase(db)
			dialog.ShowInformation("Очистка базы данных", "База данных успешно очищена", settingsWindow)
		}),
		widget.NewButtonWithIcon("Назад", nil, func() {
			settingsWindow.Close()
		}),
	)


	content := container.NewBorder(label, nil, buttons, nil, backupList)
	settingsWindow.SetContent(content)
	
	settingsWindow.Show()
	return
}

func setBackup(path string) error{
	db := getDB()
	// Открытие файла SQL-скрипта
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Чтение содержимого файла
	script, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Разбиение SQL-скрипта на отдельные запросы
	queries := strings.Split(string(script), ";")

	// Выполнение каждого запроса поочередно
	for _, query := range queries {
		// Игнорируем пустые строки
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err = db.Exec(query)
		if err != nil {
			return fmt.Errorf("ошибка при выполнении запроса: %w, запрос: %s", err, query)
		}
	}

	return nil

}

func getBackups() []string {
	directoryPath := "./dbBackup"

	// Массив для хранения путей к файлам и директориям
	var files []string

	// Используем filepath.Walk для рекурсивного сканирования каталога
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Добавляем путь в массив
		files = append(files, path)
		return nil
	})

	if err != nil {
		fmt.Println("Ошибка при сканировании каталога:", err)
		dialog.ShowError(err, nil)

	}

	return files
}