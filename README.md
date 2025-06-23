<!--
 * @Author: FunctionSir
 * @License: AGPLv3
 * @Date: 2025-06-21 11:42:21
 * @LastEditTime: 2025-06-23 10:39:54
 * @LastEditors: FunctionSir
 * @Description: -
 * @FilePath: /biblio-matrix/README.md
-->
# biblio-matrix

Library management system. A DB course project at SDTBU.

This is a free (liber) software under AGPLv3.

## How to deploy

### Build the backend

You need the Go programming language dev environment and you might need Internet connection.

``` bash
go build -o biblio-matrix -ldflags '-s -w'
```

### Build the frontend

You need Node.js related env.

``` bash
cd frontend
npm install
npm run build
```

### Run it

``` bash
./biblio-matrix CONFIG-FILE
```

### Conf example

``` ini
[options]
DB = "server=MS-SQL-SERVER[,PORT];user id=USER;password=PASSWD;database=DB"
Addr = "127.0.0.1:8080"
Frontend = "./frontend/dist"
BCryptCost = 14
Cert = "TLS-CERT"
Key = "TLS-KEY"
```

### Tips

You can use mkdb.sql to help you to create the DB and tables, etc.

## Authors

Backend, SQL, and a really little bit of frontend by @FunctionSir.

Frontend mainly by @shangaoyan.

## Repo avaliable at

FunctionSir's: <https://github.com/FunctionSir/biblio-matrix>

shangaoyan's: <https://github.com/shangaoyan/2025_CS_DB_Library_Manage_System?tab=readme-ov-file>
