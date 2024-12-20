package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/liushuangls/go-anthropic/v2"
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

// Utility functions for SQL-to-JSON conversion
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

// Database initialization
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

// LLM Client initialization
func initLLMClient() *anthropic.Client {
	err := godotenv.Load("../.env")
	checkError(err)
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	return anthropic.NewClient(apiKey)
}

// Query database for beers
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

func mapRowsToResponses(rows *sql.Rows) []BeerResponse {
	var beers []Beer
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

// Extract beer name using LLM
func extractBeerName(llmClient *anthropic.Client, mediaType string, imageData []byte) (string, error) {
	llmResp, err := llmClient.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Opus20240229,
		Messages: []anthropic.Message{
			{
				Role: anthropic.RoleUser,
				Content: []anthropic.MessageContent{
					anthropic.NewImageMessageContent(
						anthropic.NewMessageContentSource(
							anthropic.MessagesContentSourceTypeBase64,
							mediaType,
							imageData,
						),
					),
					anthropic.NewTextMessageContent("Respond only with the name of this beer."),
				},
			},
		},
		MaxTokens: 1000,
	})
	if err != nil {
		return "", err
	}

	if len(llmResp.Content) == 0 {
		return "", errors.New("No response from LLM")
	}

	return strings.TrimSpace(llmResp.Content[0].GetText()), nil
}

// Identify beer by querying database
func identifyBeer(db *sql.DB, beerName string) (Beer, error) {
	query := `
        SELECT * FROM beers
        WHERE name ILIKE '%' || $1 || '%'
        OR style ILIKE '%' || $1 || '%'
        OR description ILIKE '%' || $1 || '%'
        LIMIT 1`
	var beer Beer
	err := db.QueryRow(query, beerName).Scan(
		&beer.ID, &beer.Name, &beer.Style,
		&beer.Description, &beer.ABV, &beer.IBU,
		&beer.BPVerified, &beer.BrewerVerified,
		&beer.LastModified, &beer.BrewerID,
	)
	return beer, err
}

// Error handling
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db := initDatabase()
	defer db.Close()

	llmClient := initLLMClient()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/beers", func(c *gin.Context) {
		search := c.Query("search")
		beers := queryDatabase(db, search)
		c.IndentedJSON(http.StatusOK, beers)
	})

	router.POST("/identify-beer", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
			return
		}

		imageFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process image"})
			return
		}
		defer imageFile.Close()

		imageData, err := io.ReadAll(imageFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image"})
			return
		}

		imageMediaType := file.Header.Get("Content-Type")
		beerName, err := extractBeerName(llmClient, imageMediaType, imageData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		beer, err := identifyBeer(db, beerName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Beer not found in the database"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
			}
			return
		}

		c.JSON(http.StatusOK, convertToResponse(beer))
	})

	router.Run(":8080")
}
