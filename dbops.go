/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-06-20 09:50:27
 * @LastEditTime: 2025-06-21 16:26:38
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /biblio-matrix/dbops.go
 */

package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/microsoft/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)

type Book struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Price  int    `json:"price"`
	Count  int    `json:"count"`
}

type Record struct {
	Username   string    `json:"username"` // Who borrowed the book
	Id         string    `json:"id"`       // Which book was borrowed
	BorrowedAt time.Time `json:"borrowed"` // Borrowed time
	ReturnAt   time.Time `json:"return"`   // Must return before or at this time
}

func DbOpen(conn string) *sql.DB {
	db, err := sql.Open("mssql", conn)
	if err != nil {
		log.Fatalln("Error occurred when opening the database: " + strings.Trim(err.Error(), "\n"))
	}
	return db
}

func DbPrepare(db *sql.DB, query string) *sql.Stmt {
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalln("Error occurred when preparing the SQL statement: " + strings.Trim(err.Error(), "\n"))
	}
	return stmt
}

func AuthReader(username string, passwd string) bool {
	db := DbOpen(DbConn)
	stmt := DbPrepare(db, "SELECT PASSWD FROM READERS WHERE USERNAME=?")
	row := stmt.QueryRow(username)
	var hashedPasswd string
	err := row.Scan(&hashedPasswd)
	if err != nil {
		return false
	}
	if bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(passwd)) == nil {
		return true
	}
	return false
}

func AuthAdmin(username string, passwd string) bool {
	db := DbOpen(DbConn)
	stmt := DbPrepare(db, "SELECT PASSWD FROM ADMINS WHERE USERNAME=?")
	row := stmt.QueryRow(username)
	var hashedPasswd string
	err := row.Scan(&hashedPasswd)
	if err != nil {
		return false
	}
	if bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(passwd)) == nil {
		return true
	}
	return false
}

func Borrow(ctx context.Context, username string, bookId string, borrowedAt time.Time, returnAt time.Time) string {
	db := DbOpen(DbConn)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "无法启动事务. 请联系管理员."
	}
	defer tx.Rollback() // If anything fail, rollback the transaction.
	row := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM RECORDS WHERE ID=? AND USERNAME=?", bookId, username)
	var alreadyBorrowed int
	err = row.Scan(&alreadyBorrowed)
	if err != nil {
		return "无法完成查询. 请联系管理员."
	}
	if alreadyBorrowed > 0 {
		return "您已经借过该书了."
	}
	row = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM BOOKS WHERE ID=? AND CNT>=1", bookId)
	var bookEntryCnt int
	err = row.Scan(&bookEntryCnt)
	if err != nil {
		return "无法完成查询. 请联系管理员."
	}
	if bookEntryCnt != 1 {
		return "该书已无剩余库存, 或不存在, 或存在ID冲突."
	}
	_, err = tx.ExecContext(ctx, "UPDATE BOOKS SET CNT=CNT-1 WHERE ID=?", bookId)
	if err != nil {
		return "无法完成借书, 请联系管理员."
	}
	_, err = tx.ExecContext(ctx, "UPDATE READERS SET CNT=CNT+1 WHERE USERNAME=?", username)
	if err != nil {
		return "无法完成借书, 请联系管理员."
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO RECORDS VALUES (?,?,?,?)", username, bookId, borrowedAt, returnAt)
	if err != nil {
		return "无法完成借书, 请联系管理员."
	}
	if err = tx.Commit(); err != nil {
		return "无法成功提交事务. 请联系管理员."
	}
	return ""
}

func Return(ctx context.Context, username string, bookId string) string {
	db := DbOpen(DbConn)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "无法启动事务. 请联系管理员."
	}
	defer tx.Rollback() // If anything fail, rollback the transaction.
	row := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM RECORDS WHERE ID=? AND USERNAME=?", bookId, username)
	var alreadyBorrowed int
	err = row.Scan(&alreadyBorrowed)
	if err != nil {
		return "无法完成查询. 请联系管理员."
	}
	if alreadyBorrowed <= 0 {
		return "您还没有借过过该书."
	}
	_, err = tx.ExecContext(ctx, "DELETE FROM RECORDS WHERE ID=? AND USERNAME=?", bookId, username)
	if err != nil {
		return "无法完成还书, 请联系管理员."
	}
	_, err = tx.ExecContext(ctx, "UPDATE BOOKS SET CNT=CNT+1 WHERE ID=?", bookId)
	if err != nil {
		return "无法完成还书, 请联系管理员."
	}
	_, err = tx.ExecContext(ctx, "UPDATE READERS SET CNT=CNT-1 WHERE USERNAME=?", username)
	if err != nil {
		return "无法完成还书, 请联系管理员."
	}
	if err = tx.Commit(); err != nil {
		return "无法成功提交事务. 请联系管理员."
	}
	return ""
}

func ListBooks() []Book {
	books := make([]Book, 0)
	db := DbOpen(DbConn)
	stmt := DbPrepare(db, "SELECT ID,NAME,AUTHOR,PRICE,CNT FROM BOOKS")
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	var tmp Book
	for rows.Next() {
		rows.Scan(&tmp.Id, &tmp.Name, &tmp.Author, &tmp.Price, &tmp.Count)
		books = append(books, tmp)
	}
	return books
}

func IsBookExists(id string) (bool, error) {
	db := DbOpen(DbConn)
	stmt := DbPrepare(db, "SELECT COUNT(*) FROM BOOKS WHERE ID=?")
	res := stmt.QueryRow(id)
	var booksEntriesExists int
	err := res.Scan(&booksEntriesExists)
	if err != nil {
		return false, err
	}
	return booksEntriesExists > 0, nil
}

func AddCnt(id string, count int) {
	db := DbOpen(DbConn)
	var stmt *sql.Stmt = nil
	if count == 0 {
		stmt = DbPrepare(db, "UPDATE BOOKS SET CNT=0 WHERE ID=?")
		stmt.Exec(id)
		return
	}
	stmt = DbPrepare(db, "UPDATE BOOKS SET CNT=CNT+? WHERE ID=?")
	stmt.Exec(count, id)
}

func AddBook(b Book) error {
	db := DbOpen(DbConn)
	stmt := DbPrepare(db, "INSERT INTO BOOKS VALUES (?,?,?,?,?)")
	_, err := stmt.Exec(b.Id, b.Name, b.Author, b.Price, b.Count)
	return err
}

func AddReader(username, passwd, name string) error {
	db := DbOpen(DbConn)
	stmt := DbPrepare(db, "INSERT INTO READERS VALUES (?,?,?,?)")
	tmp, err := bcrypt.GenerateFromPassword([]byte(passwd), BCryptCost)
	if err != nil {
		return err
	}
	passwd = string(tmp)
	_, err = stmt.Exec(username, passwd, name, 0)
	return err
}

func AddAdmin(username, passwd, name string) error {
	db := DbOpen(DbConn)
	stmt := DbPrepare(db, "INSERT INTO ADMINS VALUES (?,?,?)")
	tmp, err := bcrypt.GenerateFromPassword([]byte(passwd), BCryptCost)
	if err != nil {
		return err
	}
	passwd = string(tmp)
	_, err = stmt.Exec(username, passwd, name)
	return err
}

func ListRecords(username string) []Record {
	result := make([]Record, 0)
	db := DbOpen(DbConn)
	var rows *sql.Rows
	var err error
	if username != "*" {
		stmt := DbPrepare(db, "SELECT * FROM RECORDS WHERE USERNAME=? ORDER BY \"RETURN\"")
		rows, err = stmt.Query(username)
	} else {
		stmt := DbPrepare(db, "SELECT * FROM RECORDS ORDER BY \"RETURN\"")
		rows, err = stmt.Query()
	}
	if err != nil {
		return nil
	}
	var tmp Record
	for rows.Next() {
		rows.Scan(&tmp.Username, &tmp.Id, &tmp.BorrowedAt, &tmp.ReturnAt)
		result = append(result, tmp)
	}
	return result
}
