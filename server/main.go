package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Beer struct {
	ID             sql.NullString  `json:"id"`
	Name           sql.NullString  `json:"name"`
	Style          sql.NullString  `json:"style"`
	Description    sql.NullString  `json:"description"`
	ABV            sql.NullFloat64 `json:"abv"`
	IBU            sql.NullInt32   `json:"ibu"`
	BPVerified     sql.NullBool    `json:"bpVerified"`
	BrewerVerified sql.NullBool    `json:"brewerVerified"`
	LastModified   sql.NullInt32   `json:"lastModified"`
	BrewerID       sql.NullString  `json:"brewerId"`
}

type BeerResponse struct {
	ID             *string  `json:"id,omitempty"`
	Name           *string  `json:"name,omitempty"`
	Style          *string  `json:"style,omitempty"`
	Description    *string  `json:"description,omitempty"`
	ABV            *float64 `json:"abv,omitempty"`
	IBU            *int     `json:"ibu,omitempty"`
	BPVerified     *bool    `json:"bpVerified,omitempty"`
	BrewerVerified *bool    `json:"brewerVerified,omitempty"`
	LastModified   *int     `json:"lastModified,omitempty"`
	BrewerID       *string  `json:"brewerId,omitempty"`
}

func convertToResponse(beer Beer) BeerResponse {
	return BeerResponse{
		ID:             toPtrString(beer.ID),
		Name:           toPtrString(beer.Name),
		Style:          toPtrString(beer.Style),
		Description:    toPtrString(beer.Description),
		ABV:            toPtrFloat64(beer.ABV),
		IBU:            toPtrInt(beer.IBU),
		BPVerified:     toPtrBool(beer.BPVerified),
		BrewerVerified: toPtrBool(beer.BrewerVerified),
		LastModified:   toPtrInt(beer.LastModified),
		BrewerID:       toPtrString(beer.BrewerID),
	}
}

func toPtrString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func toPtrFloat64(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}

func toPtrInt(ni sql.NullInt32) *int {
	if ni.Valid {
		val := int(ni.Int32)
		return &val
	}
	return nil
}

func toPtrBool(nb sql.NullBool) *bool {
	if nb.Valid {
		return &nb.Bool
	}
	return nil
}

func initDatabase() *sql.DB {
	err := godotenv.Load("../.env")
	checkError(err)

	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	checkError(err)

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		host, port, user, password, dbName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	checkError(err)

	err = db.Ping()
	checkError(err)

	return db
}

func mapRowsToResponses(rows *sql.Rows) []BeerResponse {
	beers := []Beer{}

	for rows.Next() {
		var beer Beer
		err := rows.Scan(
			&beer.ID, &beer.Name, &beer.Style,
			&beer.Description, &beer.ABV, &beer.IBU,
			&beer.BPVerified, &beer.BrewerVerified,
			&beer.LastModified, &beer.BrewerID,
		)
		checkError(err)

		beers = append(beers, beer)
	}

	beerResponses := make([]BeerResponse, len(beers))
	for i, beer := range beers {
		beerResponses[i] = convertToResponse(beer)
	}
	return beerResponses
}

func queryDatabase(db *sql.DB, search string) []BeerResponse {
	query := `
        SELECT * FROM beers
        WHERE name ILIKE '%' || $1 || '%'
        OR style ILIKE '%' || $1 || '%'
        OR description ILIKE '%' || $1 || '%'
        LIMIT 10`
	rows, err := db.Query(query, search)
	checkError(err)
	defer rows.Close()

	return mapRowsToResponses(rows)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db := initDatabase()
	defer db.Close()

	router := gin.Default()

	router.GET("/beers/", func(c *gin.Context) {
		search := c.Query("search")
		beers := queryDatabase(db, search)
		c.IndentedJSON(http.StatusOK, beers)
	})

	router.Run(":8080")
}
