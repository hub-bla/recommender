package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/lib/pq"
)

const EMB_SIZE int = 384

type Embedding struct {
	Data []float32 `json:"embedding"`
}

func main() {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "admin"
		dbname   = "movie_recommender"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
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

	test_message := []byte("{\"movie_title\": \"string\"}")
	bodyReader := bytes.NewReader(test_message)
	res, err := http.Post("http://127.0.0.1:8000/model/", "application/json", bodyReader)

	if err != nil {
		fmt.Println("something went wrong in making request")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("something went wrong in reading body")
	}

	var titleEmbedding Embedding

	json.Unmarshal(resBody, &titleEmbedding)

	fmt.Printf("Response: %v", titleEmbedding.Data)

}
