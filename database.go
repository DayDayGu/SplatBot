// Package main provides ...
package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	// "os"
)

// DefaultDB 默认数据库
var DefaultDB *sql.DB

// VolumnPath 获取存储路径，区分环境
func VolumnPath(file string) string {
	s := os.Getenv("SPLAT_ENV")
	if s == "release" {
		return "/SplatBot/" + file
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir + "/" + file
}

// Delete tmp dir.
func ClearTmpPath() {
	s := os.Getenv("SPLAT_ENV")
	if s != "release" {
		os.RemoveAll(VolumnPath("tmp"))
		os.Mkdir(VolumnPath("tmp"), os.ModePerm)
	}
}

// InitDatabase init database for splat
func InitDatabase() {
	dbPath := VolumnPath("splat.db")
	db, err := sql.Open("sqlite3", dbPath)
	DefaultDB = db
	if err != nil {
		log.Fatal(err)
	}
	createLeagueTable()
}

func createLeagueTable() {
	sql := `
    create table if not exists league (id integer primary key, 
                                       member1 text default '', 
                                       member2 text default '', 
                                       member3 text default '', 
                                       member4 text default '', 
                                       memberid1 integer default 0,
                                       memberid2 integer default 0,
                                       memberid3 integer default 0,
                                       memberid4 integer default 0,
                                       start_time integer default 0,
                                       rule text default '',
                                       type integer default 0,
                                       status integer default 0,
                                       create_date datetime default current_timestamp);`
	_, err := DefaultDB.Exec(sql)
	if err != nil {
		log.Printf("%q: %s\n", err, sql)
		return
	}
}
