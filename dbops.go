/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-06-20 09:50:27
 * @LastEditTime: 2025-06-20 11:22:44
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /biblio-matrix/dbops.go
 */

package main

import (
	"database/sql"
	"log"
	"strings"
)

// Open DB using DB file specified in global var.
func DbOpen(conn string) *sql.DB {
	db, err := sql.Open("mssql", conn)
	if err != nil {
		log.Fatalln("Error occurred when opening the database: " + strings.Trim(err.Error(), "\n"))
	}
	return db
}

// Prepare a query.
func DbPrepare(db *sql.DB, query string) *sql.Stmt {
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalln("Error occurred when preparing the SQL statement: " + strings.Trim(err.Error(), "\n"))
	}
	return stmt
}

// In dev...
func Auth(username string, passwd string) bool {
	db := DbOpen(DbConn)
	stmt := DbPrepare(db, "SELECT ")
	_ = stmt
	return true
}

// In dev...
func ChkIsAdmin(username string) bool {
	db := DbOpen(DbConn)
	stmt := DbPrepare(db, "SELECT ")
	_ = stmt
	return true
}
