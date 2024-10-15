package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func Confirm(myApp fyne.App, onConfirm func()) {
	confWindow := myApp.NewWindow("Подтвержите действие")
	confWindow.Resize(fyne.NewSize(300, 100))
	confWindow.CenterOnScreen()

	label := widget.NewLabel("Введите пароль")
	password := widget.NewPasswordEntry()

	buttons := container.New(layout.NewGridLayout(2),
		widget.NewButton("Подтвердить", func() {
			if password.Text == "Caochue4" {
				onConfirm() // вызываем callback при правильном пароле
				confWindow.Close() // закрываем окно после подтверждения
			}else{
				dialog.ShowError(fmt.Errorf("Неверный пароль"), confWindow)
			}
		}),
		widget.NewButton("Отмена", func() {
			confWindow.Close()
		}),
	)

	content := container.NewVBox(label, password, buttons)
	confWindow.SetContent(content)

	confWindow.Show() // используем Show, чтобы не блокировать главный поток
}



