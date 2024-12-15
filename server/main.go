package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select name from beers limit 10")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var val string
		err := rows.Scan(&val)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		fmt.Println("Row value:", val)
	}

	fmt.Println("Successfully connected!")
	// router := gin.Default()
	// router.GET("/beers", getBeers)
	// router.GET("/beers/:name", getBeerByName)
	// router.POST("/beers", postBeers)

	// router.Run("localhost:8080")
}
