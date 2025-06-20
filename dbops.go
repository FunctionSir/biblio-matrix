/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-06-20 09:50:27
 * @LastEditTime: 2025-06-20 16:28:57
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

	"golang.org/x/crypto/bcrypt"
)

type Book struct {
	Id     string
	Name   string
	Author string
	Price  int
	Count  int
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

func AuthAdmin(username string, passwd string) bool {
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

func Borrow(ctx context.Context, username string, bookId string, borrowedAt time.Time, returnAt time.Time) string {
	db := DbOpen(DbConn)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "无法启动事务. 请联系管理员."
	}
	defer tx.Rollback() // If anything fail, rollback the transaction.
	row := tx.QueryRowContext(ctx, "SELECT COUNT(*)>0 FROM RECORDS WHERE ID=? AND USERNAME=?", bookId, username)
	var alreadyBorrowed bool
	err = row.Scan(&alreadyBorrowed)
	if err != nil {
		return "无法完成查询. 请联系管理员."
	}
	if alreadyBorrowed {
		return "您已经借过该书了."
	}
	row = tx.QueryRowContext(ctx, "SELECT COUNT(*)=1 FROM BOOKS WHERE ID=? AND CNT>=1", bookId)
	var canBorrow bool
	err = row.Scan(&canBorrow)
	if err != nil {
		return "无法完成查询. 请联系管理员."
	}
	if !canBorrow {
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
	row := tx.QueryRowContext(ctx, "SELECT COUNT(*)>0 FROM RECORDS WHERE ID=? AND USERNAME=?", bookId, username)
	var alreadyBorrowed bool
	err = row.Scan(&alreadyBorrowed)
	if err != nil {
		return "无法完成查询. 请联系管理员."
	}
	if !alreadyBorrowed {
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

func List() []Book {
	books := make([]Book, 0)
	db := DbOpen(DbConn)
	stmt := DbPrepare(db, "SELECT ID,NAME,AUTHOR,PRICE,CNT FROM READERS FROM BOOKS")
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
	stmt := DbPrepare(db, "SELECT COUNT(*)>0 FROM BOOKS WHERE ID=?")
	res := stmt.QueryRow(id)
	var bookExists bool
	err := res.Scan(&bookExists)
	if err != nil {
		return false, err
	}
	return bookExists, nil
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
	stmt.Exec(id, count)
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
	tmp, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	passwd = string(tmp)
	_, err = stmt.Exec(username, passwd, name, 0)
	return err
}
