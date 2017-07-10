package main

import (
    _ "github.com/mattn/go-sqlite3"
    "database/sql"
    "log"
    "fmt"
)

type TestItem struct {
    MessageID   string
    MessageFrom string
    MessageText string
    MessageTime string
}

type ReportItem struct {
    MessageID string
}

func InitDB(filepath string) *sql.DB {
    file := fmt.Sprintf("file:%v?cache=shared&mode=rwc", filepath)
    db, err := sql.Open("sqlite3", file)
    check(err)
    if db == nil { panic("db nil") }
    return db
}

func CreateTableReport(db *sql.DB) {
    // create table if not exists
    sql_table := `
    CREATE TABLE IF NOT EXISTS report(
        messageid string primary key
    );
    `
    _, err := db.Exec(sql_table)
    check(err)
}

func CreateTable(db *sql.DB) {
    // create table if not exists
    sql_table := `
    CREATE TABLE IF NOT EXISTS messages(
        messageid string primary key,
        messagefrom string,
        messagetext string,
        messagetime string
    );
    `
    _, err := db.Exec(sql_table)
    check(err)
}

func ReadItemsReport(db *sql.DB) (rows *sql.Rows) {
    rows, err :=  db.Query("SELECT * FROM report")
    check(err)
    return rows
}

func ReadItems(db *sql.DB) (rows *sql.Rows) {
    rows, err := db.Query("SELECT * FROM messages")
    check(err)
    return rows
}

func StoreItemReport(db *sql.DB, items []ReportItem) {
    sql_additem := `
    INSERT INTO report(
        MessageID
    ) values(?)
    `

    routeSQL, err := db.Prepare(sql_additem)
    check(err)

    tx, err := db.Begin()
    check(err)
    _, err = tx.Stmt(routeSQL).Exec(items[0].MessageID)
    if err != nil {
        log.Println(err)
        log.Println("doing rollback")
        tx.Rollback()
    } else {
        err = tx.Commit()
        check(err)
    }
}

func StoreItem(db *sql.DB, items []TestItem) {
    sql_additem := `
    INSERT OR REPLACE INTO messages(
        MessageID,
        MessageFrom,
        MessageText,
        MessageTime
    ) values(?, ?, ?, ?)
    `

    routeSQL, err := db.Prepare(sql_additem)
    check(err)

    for _, item := range items {
        tx, err := db.Begin()
        check(err)
        _, err = tx.Stmt(routeSQL).Exec(item.MessageID, item.MessageFrom, item.MessageText, item.MessageTime)
        if err != nil {
            log.Println(err)
            log.Println("doing rollback")
            tx.Rollback()
        } else {
            err = tx.Commit()
            check(err)
        }
    }
}
