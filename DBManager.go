package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var gdb *sql.DB

func getDB() *sql.DB {
	return gdb
}

func readConfig() (string, string, string, error){
	file, err := os.Open("dbConfig.txt")
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	var user, password, database string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "= ", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]

		switch key {
		case "user":
			user = value
		case "password":
			password = value
		case "database":
			database = value
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", "", err
	}

	return user, password, database, nil
}

func BackupDB(user, password, dbName, backupDir, host string) error{
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}
	defer db.Close()
	
	
	// Формируем имя файла для бэкапа с указанием текущей даты
	timestamp := time.Now().Format("20060102_150405")
	backupFile := fmt.Sprintf("%s/%s_backup_%s.sql", backupDir, dbName, timestamp)

	// Создаем файл для сохранения бэкапа
	outfile, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("не удалось создать файл бэкапа: %w", err)
	}
	defer outfile.Close()

	// Пишем информацию о бэкапе
	outfile.WriteString(fmt.Sprintf("-- Бэкап базы данных %s\n-- Время: %s\n\n", dbName, timestamp))

	// Таблицы для бэкапа
	tables := []string{"Nomenclature", "Tasks", "Workers", "Task_Workers"}

	// Бэкап каждой таблицы
	for _, table := range tables {
		err := backupTable(db, outfile, table)
		if err != nil {
			return fmt.Errorf("ошибка при создании бэкапа таблицы %s: %w", table, err)
		}
	}

	fmt.Printf("Бэкап базы данных %s успешно создан в файле: %s\n", dbName, backupFile)
	return nil
}

func backupTable(db *sql.DB, outfile *os.File, table string) error {
	// Получаем схему таблицы
	rows, err := db.Query(fmt.Sprintf("SHOW CREATE TABLE %s", table))
	if err != nil {
		return fmt.Errorf("не удалось получить схему таблицы %s: %w", table, err)
	}
	defer rows.Close()

	var tableName, createTableSQL string
	if rows.Next() {
		if err := rows.Scan(&tableName, &createTableSQL); err != nil {
			return fmt.Errorf("ошибка при получении схемы таблицы %s: %w", table, err)
		}
	}

	// Записываем схему таблицы в файл
	outfile.WriteString(fmt.Sprintf("-- Схема таблицы %s\n", table))
	outfile.WriteString(createTableSQL + ";\n\n")

	// Получаем данные из таблицы
	dataRows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return fmt.Errorf("не удалось получить данные из таблицы %s: %w", table, err)
	}
	defer dataRows.Close()

	// Получаем столбцы таблицы
	columns, err := dataRows.Columns()
	if err != nil {
		return fmt.Errorf("ошибка при получении столбцов для таблицы %s: %w", table, err)
	}

	// Подготовка для сканирования данных
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Пишем данные таблицы в файл
	outfile.WriteString(fmt.Sprintf("-- Данные таблицы %s\n", table))
	for dataRows.Next() {
		err = dataRows.Scan(scanArgs...)
		if err != nil {
			return fmt.Errorf("ошибка при чтении данных из таблицы %s: %w", table, err)
		}

		// Генерируем строку для вставки данных
		insertSQL := fmt.Sprintf("INSERT INTO %s VALUES (", table)
		for i, col := range values {
			if col == nil {
				insertSQL += "NULL"
			} else {
				insertSQL += fmt.Sprintf("'%s'", string(col))
			}
			if i < len(values)-1 {
				insertSQL += ", "
			}
		}
		insertSQL += ");\n"

		// Записываем в файл
		outfile.WriteString(insertSQL)
	}
	outfile.WriteString("\n")

	return nil
}

func dbInit() (*sql.DB, error) {
	user, password, database, err := readConfig()

	dsn := fmt.Sprintf("%s:%s@/%s", user, password, database)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	gdb = db
	return db, nil
}

var queryCheck = []bool{false, false, false, false, false, false, false, false}

func GetTasks(db *sql.DB, tSort int) ([][]string, error) {
	query := "SELECT task_id, task_custom_id, task_accept, task_deli, task_client, nom_name, task_count, task_done, task_about FROM Tasks JOIN Nomenclature ON Tasks.task_name = Nomenclature.nom_id ORDER BY"
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
			query += " task_custom_id DESC"
		} else {
			query += " task_custom_id"
		}
		queryCheck[1] = !queryCheck[1]
	case 2:
		if queryCheck[1] {
			query += " task_accept DESC"
		} else {
			query += " task_accept"
		}
		queryCheck[1] = !queryCheck[1]
	case 3:
		if queryCheck[2] {
			query += " task_deli DESC"
		} else {
			query += " task_deli"
		}
		queryCheck[2] = !queryCheck[2]
	case 4:
		if queryCheck[3] {
			query += " task_client DESC"
		} else {
			query += " task_client"
		}
		queryCheck[3] = !queryCheck[3]
	case 5:
		if queryCheck[4] {
			query += " nom_name DESC"
		} else {
			query += " nom_name"
		}
		queryCheck[4] = !queryCheck[4]
	case 6:
		if queryCheck[5] {
			query += " task_count DESC"
		} else {
			query += " task_count"
		}
		queryCheck[5] = !queryCheck[5]
	case 7:
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
		var taskDone sql.NullInt64
		var taskID, taskcustomID, taskCount int
		var taskAccept, taskDeli, taskClient string
		var taskName string
		var taskAbout string
		err := rows.Scan(&taskID, &taskcustomID, &taskAccept, &taskDeli, &taskClient, &taskName, &taskCount, &taskDone, &taskAbout)
		if err != nil {
			return nil, err
		}
		task := []string{
			fmt.Sprintf("%d", taskID),
			fmt.Sprintf("%d", taskcustomID),
			taskAccept,
			taskDeli,
			taskClient,
			taskName,
			fmt.Sprintf("%d", taskCount),
			fmt.Sprintf("%d", taskDone.Int64),
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
	var taskDone int = 0
	_, err := db.Exec("INSERT INTO Tasks (task_custom_id ,task_accept, task_deli, task_client, task_name, task_count, task_done) VALUES (?, ?, ?, ?, ?, ?, COALESCE(?, 0))",
	id ,acceptDate, deliDate, client, nomID, quantity, taskDone)
	return err
}

func GetTaskByID(db *sql.DB, id int) ([][]string, error) {//SELECT task_worker_id, task_id, CONCAT(worker_fname, ' ', worker_sname) AS worker_name, tw_done, tw_date, tw_day_sum FROM Task_Workers JOIN Workers ON Task_Workers.worker_id = Workers.worker_id WHERE task_id = ?
	rows, err := db.Query("SELECT task_worker_id, task_id, CONCAT(worker_fname, ' ', worker_sname) AS worker_name, tw_done, tw_date, tw_day_sum FROM Task_Workers JOIN Workers ON Task_Workers.worker_id = Workers.worker_id WHERE task_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskWorkers [][]string
	for rows.Next() {
		var taskWorkerID, taskID, twDone int
		var twDate, tw_day_sum, workerName string
		err := rows.Scan(&taskWorkerID, &taskID, &workerName, &twDone, &twDate, &tw_day_sum)
		if err != nil {
			return nil, err
		}
		taskWorker := []string{
			fmt.Sprintf("%d", taskWorkerID),
			fmt.Sprintf("%d", taskID),
			workerName,
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
	//запрос который выдает название номенкулатуры по id задачи
	err := db.QueryRow("SELECT nom_name FROM Nomenclature WHERE nom_id = (SELECT task_name FROM Tasks WHERE task_id = ?)", id).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	return name
}

func DeleteTask(db *sql.DB, id int) error {
	//удаление всех записей из таблицы Task_Workers, где task_id = id
	
	_, err := db.Exec("DELETE FROM Task_Workers WHERE task_id = ?", id)
	_, err = db.Exec("DELETE FROM Tasks WHERE task_id = ?", id)
	
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




func GetWorkerID(selectedWorker string) int {
	db := getDB()
	var id int
	err := db.QueryRow("SELECT worker_id FROM Workers WHERE CONCAT(worker_fname, ' ', worker_sname) = ?", selectedWorker).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func AddTaskWorker(db *sql.DB, taskID int, workerID int, done int, date string, daySum float64) (int64, error) {
	fmt.Println(taskID, workerID, done)
	res, err := db.Exec("INSERT INTO Task_Workers (task_id, worker_id, tw_done, tw_date, tw_day_sum) VALUES (?, ?, ?, ?, ?)", taskID, workerID, done, date, daySum)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}


//обновить выполненное в tasks основываясь на сумме выполненного в task_workers во всей базе дынн
func UpdateDataBase(db *sql.DB) {
	rows, err := db.Query("SELECT task_id FROM Tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var taskID int
		err := rows.Scan(&taskID)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec("UPDATE Tasks SET task_done = (SELECT SUM(tw_done) FROM Task_Workers WHERE task_id = ?) WHERE task_id = ?", taskID, taskID)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Updated db")
}

func GetSummary(db *sql.DB, name string, startDate string, endDate string) ([][]string, error) {
	// Разделяем имя на fname и sname
	names := strings.SplitN(name, " ", 2)
	if len(names) != 2 {
		return nil, fmt.Errorf("incorrect name format, expected 'fname sname'")
	}
	fname := names[0]
	sname := names[1]

	// Подготовка SQL-запроса
	query := `
		SELECT 
			tw.task_worker_id,
			t.task_name,
			CONCAT(w.worker_fname, ' ', w.worker_sname) AS worker_name,
			tw.tw_done AS task_count,
			tw.tw_day_sum AS total_sum,
			tw.tw_date AS task_date
		FROM 
			Workers w
			JOIN Task_Workers tw ON w.worker_id = tw.worker_id
			JOIN Tasks t ON t.task_id = tw.task_id
		WHERE 
			w.worker_fname = ? AND w.worker_sname = ? AND tw.tw_date BETWEEN ? AND ?
	`

	// Выполнение запроса
	rows, err := db.Query(query, fname, sname, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taskWorkers [][]string
	for rows.Next() {
		var taskWorkerID, taskID, twDone int
		var twDate, tw_day_sum, workerName string
		err := rows.Scan(&taskWorkerID, &taskID, &workerName, &twDone, &twDate, &tw_day_sum)
		if err != nil {
			return nil, err
		}
		taskWorker := []string{
			fmt.Sprintf("%d", taskWorkerID),
			fmt.Sprintf("%d", taskID),
			workerName,
			fmt.Sprintf("%d", twDone),
			twDate,
			tw_day_sum,
		}
		taskWorkers = append(taskWorkers, taskWorker)
	}
	
	return taskWorkers, nil
}

func GetSummarySum(db *sql.DB, name string, startDate string, endDate string) string {
	// Разделяем имя на fname и sname
	names := strings.SplitN(name, " ", 2)
	if len(names) != 2 {
		return ""
	}
	fname := names[0]
	sname := names[1]

	// Подготовка SQL-запроса
	query := `
		SELECT 
			SUM(tw.tw_day_sum)
		FROM 
			Workers w
			JOIN Task_Workers tw ON w.worker_id = tw.worker_id
			JOIN Tasks t ON t.task_id = tw.task_id
		WHERE 
			w.worker_fname = ? AND w.worker_sname = ? AND tw.tw_date BETWEEN ? AND ?
	`

	// Выполнение запроса
	var sum string
	err := db.QueryRow(query, fname, sname, startDate, endDate).Scan(&sum)
	if err != nil {
		return ""
	}
	return sum
}