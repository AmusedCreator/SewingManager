package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func GetTasks(db *sql.DB) ([][]string, error) {
    rows, err := db.Query("SELECT task_id, task_accept, task_deli, task_client, nom_name, task_count, task_done, task_sum, task_about FROM Tasks JOIN Nomenclature ON Tasks.task_name = Nomenclature.nom_id;")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks [][]string
    for rows.Next() {
        var taskID, taskCount, taskDone, taskSum int
        var taskAccept, taskDeli, taskClient string
        var taskName string
        var taskAbout string
        err := rows.Scan(&taskID, &taskAccept, &taskDeli, &taskClient, &taskName, &taskCount, &taskDone, &taskSum, &taskAbout)
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
            fmt.Sprintf("%d", taskSum),
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
    rows, err := db.Query("SELECT task_worker_id, task_id, worker_id, tw_done, tw_date FROM Task_Workers")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var taskWorkers [][]string
    for rows.Next() {
        var taskWorkerID, taskID, workerID, twDone int
        var twDate string
        err := rows.Scan(&taskWorkerID, &taskID, &workerID, &twDone, &twDate)
        if err != nil {
            return nil, err
        }
        taskWorker := []string{
            fmt.Sprintf("%d", taskWorkerID),
            fmt.Sprintf("%d", taskID),
            fmt.Sprintf("%d", workerID),
            fmt.Sprintf("%d", twDone),
            twDate,
        }
        taskWorkers = append(taskWorkers, taskWorker)
    }
    return taskWorkers, nil
}

func GetNomenclature(db *sql.DB) ([][]string, error) {
    rows, err := db.Query("SELECT nom_id, nom_name FROM Nomenclature")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var nomenclature [][]string
    for rows.Next() {
        var nomID int
        var nomName string
        err := rows.Scan(&nomID, &nomName)
        if err != nil {
            return nil, err
        }
        nom := []string{
            fmt.Sprintf("%d", nomID),
            nomName,
        }
        nomenclature = append(nomenclature, nom)
    }
    return nomenclature, nil
}