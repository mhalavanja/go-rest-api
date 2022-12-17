package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func ExecuteStoredProcedures(db *sql.DB) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	pwd = filepath.Join(pwd, "/db/functions/")
	f, err := os.Open(pwd)
	if err != nil {
		panic(err)
	}

	fileInfo, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		panic(err)
	}

	for _, file := range fileInfo {
		b, err := os.ReadFile(filepath.Join(pwd, file.Name()))
		if err != nil {
			panic(err)
		}

		_, err = db.Exec(string(b))
		if err != nil {
			fmt.Println(file.Name())
			panic(err)
		}
	}
}
