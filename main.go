/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-06-20 08:28:04
 * @LastEditTime: 2025-06-22 15:41:49
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /biblio-matrix/main.go
 */
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/FunctionSir/goset"
	"github.com/FunctionSir/readini"
	_ "github.com/microsoft/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)

const VER string = "0.1.0"

var DbConn string
var HttpAddr string
var FrontendDir string
var BCryptCost int

func getConf() {
	if len(os.Args) < 2 {
		panic("no config file specified")
	}
	confFile, err := readini.LoadFromFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	if !confFile.HasKey("options", "DB") {
		panic("no database specified")
	}
	if !confFile.HasKey("options", "Addr") {
		panic("no listen addr specified")
	}
	if !confFile.HasKey("options", "Frontend") {
		FrontendDir = ""
	} else {
		FrontendDir = confFile["options"]["Frontend"]
	}
	if !confFile.HasKey("options", "BCryptCost") {
		BCryptCost = bcrypt.DefaultCost
	} else {
		tmp := confFile["options"]["BCryptCost"]
		BCryptCost, err = strconv.Atoi(tmp)
		if err != nil || BCryptCost < bcrypt.MinCost || BCryptCost > bcrypt.MaxCost {
			panic("bcrypt cost found but illegal")
		}
	}
	DbConn = confFile["options"]["DB"]
	HttpAddr = confFile["options"]["Addr"]
}

func main() {
	fmt.Println("Biblio Matrix Library Management System Server")
	fmt.Printf("Version: %s | This is a FOSS under AGPLv3\n", VER)
	getConf()
	log.Println("Testing DB connection...")
	db, err := sql.Open("mssql", DbConn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("MSSQL DB connection OK.")
	log.Println("Init token storage...")
	TokensSet = make(goset.Set[string])
	TokensExp = make(map[string]time.Time)
	TokensUser = make(map[string]string)
	TokensIsAdmin = make(map[string]bool)
	log.Println("Token storage ready.")
	serveHttp(HttpAddr)
}
