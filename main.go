/*
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-06-20 08:28:04
 * @LastEditTime: 2025-06-20 10:06:35
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /biblio-matrix/main.go
 */
package main

import (
	"fmt"
	"os"

	"github.com/FunctionSir/readini"
)

const VER string = "0.1.0"

var DbConn string
var HttpAddr string
var FrontendDir string

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
	DbConn = confFile["options"]["DB"]
	HttpAddr = confFile["options"]["Addr"]
}

func main() {
	fmt.Println("Biblio Matrix Library Management System Server")
	fmt.Printf("Version: %s | This is a FOSS under AGPLv3\n", VER)
	getConf()
	serveHttp(HttpAddr)
}
