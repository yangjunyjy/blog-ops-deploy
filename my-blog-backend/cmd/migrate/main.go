package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 检查参数
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <db-path> [sql-file]")
		os.Exit(1)
	}

	dbPath := os.Args[1]
	sqlFile := "migrations/sqlite/init.sql"
	if len(os.Args) > 2 {
		sqlFile = os.Args[2]
	}

	// 打开数据库连接
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 读取SQL文件
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	// 执行SQL
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		log.Fatalf("Failed to execute SQL: %v", err)
	}

	fmt.Printf("Migration completed successfully! Database: %s\n", dbPath)
}
