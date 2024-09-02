package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var gdb *sql.DB

func getDB() *sql.DB {
	return gdb
}

func dbInit() (*sql.DB, error) {
	dsn := "root:123456789@/SewingDB"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	gdb = db
	return db, nil
}

var queryCheck = []bool{false, false, false, false, false, false, false, false}

func GetTasks(db *sql.DB, tSort int) ([][]string, error) {
	query := "SELECT task_id, task_accept, task_deli, task_client, nom_name, task_count, task_done, task_about FROM Tasks JOIN Nomenclature ON Tasks.task_name = Nomenclature.nom_id ORDER BY"
	switch tSort {
	case 0:
		if queryCheck[0] {
			query += " task_id DESC"
		} else {
			query += " task_id"
		}
		queryCheck[0] = !queryCheck[0]
	case 1:
		if queryCheck[1] {
			query += " task_accept DESC"
		} else {
			query += " task_accept"
		}
		queryCheck[1] = !queryCheck[1]
	case 2:
		if queryCheck[2] {
			query += " task_deli DESC"
		} else {
			query += " task_deli"
		}
		queryCheck[2] = !queryCheck[2]
	case 3:
		if queryCheck[3] {
			query += " task_client DESC"
		} else {
			query += " task_client"
		}
		queryCheck[3] = !queryCheck[3]
	case 4:
		if queryCheck[4] {
			query += " nom_name DESC"
		} else {
			query += " nom_name"
		}
		queryCheck[4] = !queryCheck[4]
	case 5:
		if queryCheck[5] {
			query += " task_count DESC"
		} else {
			query += " task_count"
		}
		queryCheck[5] = !queryCheck[5]
	case 6:
		if queryCheck[6] {
			query += " task_done DESC"
		} else {
			query += " task_done"
		}
		queryCheck[6] = !queryCheck[6]
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks [][]string
	for rows.Next() {
		var taskID, taskCount, taskDone int
		var taskAccept, taskDeli, taskClient string
		var taskName string
		var taskAbout string
		err := rows.Scan(&taskID, &taskAccept, &taskDeli, &taskClient, &taskName, &taskCount, &taskDone, &taskAbout)
		if err != nil {
			return nil, err
		}
		task := []string{
			fmt.Sprintf("%d", taskID),
			taskAccept,
			taskDeli,
			taskClient,
			taskName,
			fmt.Sprintf("%d", taskCount),
			fmt.Sprintf("%d", taskDone),
			taskAbout,
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetWorkers(db *sql.DB) ([][]string, error) {
	rows, err := db.Query("SELECT worker_id, worker_fname, worker_sname, worker_about FROM Workers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workers [][]string
	for rows.Next() {
		var workerID int
		var workerFname, workerSname, workerAbout string
		err := rows.Scan(&workerID, &workerFname, &workerSname, &workerAbout)
		if err != nil {
			return nil, err
		}
		worker := []string{
			fmt.Sprintf("%d", workerID),
			workerFname,
			workerSname,
			workerAbout,
		}
		workers = append(workers, worker)
	}
	return workers, nil
}

func GetTaskWorkers(db *sql.DB) ([][]string, error) {
	rows, err := db.Query("SELECT task_worker_id, task_id, worker_id, tw_done, tw_date, tw_day_sum FROM Task_Workers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskWorkers [][]string
	for rows.Next() {
		var taskWorkerID, taskID, workerID, twDone int
		var twDate, tw_day_sum string
		err := rows.Scan(&taskWorkerID, &taskID, &workerID, &twDone, &twDate, &tw_day_sum)
		if err != nil {
			return nil, err
		}
		taskWorker := []string{
			fmt.Sprintf("%d", taskWorkerID),
			fmt.Sprintf("%d", taskID),
			fmt.Sprintf("%d", workerID),
			fmt.Sprintf("%d", twDone),
			twDate,
			tw_day_sum,
		}
		taskWorkers = append(taskWorkers, taskWorker)
	}
	return taskWorkers, nil
}

func GetNomenclature(db *sql.DB) ([][]string, error) {
	rows, err := db.Query("SELECT nom_id, nom_name, nom_price FROM Nomenclature")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nomenclature [][]string
	for rows.Next() {
		var nomID int
		var nomName string
		var nomPrice string // Используем string для хранения DECIMAL
		err := rows.Scan(&nomID, &nomName, &nomPrice)
		if err != nil {
			return nil, err
		}
		nom := []string{
			fmt.Sprintf("%d", nomID),
			nomName,
			nomPrice,
		}
		nomenclature = append(nomenclature, nom)
	}
	return nomenclature, nil
}

func GetNomenclatureID(db *sql.DB, name string) (int, error) {
	var id int
	err := db.QueryRow("SELECT nom_id FROM Nomenclature WHERE nom_name = ?", name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SaveTask(db *sql.DB, id, acceptDate string, deliDate, client string, nomID int, quantity string) error {
	_, err := db.Exec("INSERT INTO Tasks (task_id ,task_accept, task_deli, task_client, task_name, task_count, task_done) VALUES (?, ?, ?, ?, ?, ?, ?)",
	id ,acceptDate, deliDate, client, nomID, quantity, 0)
	return err
}

func GetTaskByID(db *sql.DB, id int) ([][]string, error) {
	rows, err := db.Query("SELECT task_worker_id, task_id, worker_id, tw_done, tw_date, tw_day_sum FROM Task_Workers WHERE task_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskWorkers [][]string
	for rows.Next() {
		var taskWorkerID, taskID, workerID, twDone int
		var twDate, tw_day_sum string
		err := rows.Scan(&taskWorkerID, &taskID, &workerID, &twDone, &twDate, &tw_day_sum)
		if err != nil {
			return nil, err
		}
		taskWorker := []string{
			fmt.Sprintf("%d", taskWorkerID),
			fmt.Sprintf("%d", taskID),
			fmt.Sprintf("%d", workerID),
			fmt.Sprintf("%d", twDone),
			twDate,
			tw_day_sum,
		}
		taskWorkers = append(taskWorkers, taskWorker)
	}
	return taskWorkers, nil
}

func GetTaskNameByID(db *sql.DB, id int) string {
	var name string
	err := db.QueryRow("SELECT task_name FROM Tasks WHERE task_id = ?", id).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	return name
}

func DeleteTask(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM Tasks WHERE task_id = ?", id)
	return err
}

// calculateDaySum вычисляет сумму за день на основе выполненного и цены за штуку
func CalculateDaySum(db *sql.DB, taskID int, done int) float64 {
	var pricePerUnit float64
	err := db.QueryRow("SELECT n.nom_price FROM Tasks t JOIN Nomenclature n ON t.task_name = n.nom_id WHERE t.task_id = ?", taskID).Scan(&pricePerUnit)
	if err != nil {
		log.Fatal(err)
	}
	return pricePerUnit * float64(done)
}



// getWorkerID извлекает ID работника из строки выбранной в списке
func GetWorkerID(selectedWorker string) string {
	// Разбираем строку и извлекаем ID работника (предполагается, что формат строки имеет ID в начале)
	return selectedWorker[:strings.Index(selectedWorker, " ")]
}
