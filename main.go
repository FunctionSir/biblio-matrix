/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-06-20 08:28:04
 * @LastEditTime: 2025-06-23 10:15:58
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
var TlsCert string
var TlsKey string

func getConf() {
	if len(os.Args) < 2 {
		panic("no config file or operation specified")
	}
	if os.Args[1] == "hash-passwd" {
		fmt.Println("Please make sure that NO ONE ELSE can see your terminal!")
		fmt.Printf("BCrypt cost (%d~%d): ", bcrypt.MinCost, bcrypt.MaxCost)
		var bcCost int
		fmt.Scanf("%d\n", &bcCost)
		fmt.Println("Again! Please make sure that NO ONE ELSE can see your terminal!")
		fmt.Print("Input password (will be echoed!): ")
		var passwd string
		fmt.Scanln(&passwd)
		fmt.Print("Hashed password: ")
		hashed, err := bcrypt.GenerateFromPassword([]byte(passwd), bcCost)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(hashed))
		os.Exit(0)
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
	if !confFile.HasKey("options", "Cert") {
		TlsCert = ""
	} else {
		TlsCert = confFile["options"]["Cert"]
	}
	if !confFile.HasKey("options", "Key") {
		TlsKey = ""
	} else {
		TlsKey = confFile["options"]["Key"]
	}
	if TlsCert == "" || TlsKey == "" {
		log.Println("Warning: Incomplete TLS config, using HTTP instead of HTTPS!")
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
