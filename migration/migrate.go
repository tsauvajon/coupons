package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tsauvajon/coupons/database"
)

func main() {
	client, err := database.NewClient()
	if err != nil {
		fmt.Println("connecting to the database failed: ", err)
		return
	}

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("please provide sql files")
		return
	}

	for _, arg := range args {
		path := filepath.Join(arg)
		fi, err := os.Stat(path)
		if err != nil {
			fmt.Println("couldn't open file or directory: ", err)
			return
		}
		switch mode := fi.Mode(); {
		case mode.IsRegular():
			migrate(path, client)
		default:
			log.Println(fi.Name(), "is not a file")
		}
	}
}

func migrate(filename string, client *database.Client) {
	queries, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	log.Println("running migration", filename)

	for _, query := range strings.Split(string(queries), ";") {
		if _, err := client.Connection.Exec(query); err != nil {
			panic(fmt.Sprintf("migration failed: %v", err))
		}
	}
}
