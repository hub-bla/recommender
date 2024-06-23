package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

const EMB_SIZE = 384
const GO_SERVER = "http://127.0.0.1:8080"
const ML_SERVER = "http://127.0.0.1:8000"
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "movie_recommender"
)

var db sql.DB

type Embedding struct {
	Idx   int       `json:"id"`
	Title string    `json:"title"`
	Data  []float32 `json:"embedding"`
}

type SimilarTitle struct {
	Idx   int    `json:"id"`
	Title string `json:"title"`
}

type SearchbarInput struct {
	Data string `json:"searchbar_input"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// GET TEXT TO MAKE EMBEDDING FROM IT
		fmt.Printf("test")

		resBody, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Printf("Couldn't read a response Body!\n")
			return
		}

		var searchbarContent SearchbarInput

		json.Unmarshal(resBody, &searchbarContent)

		ml_body := []byte(fmt.Sprintf("{\"movie_title\": \"%s\"}", searchbarContent.Data))
		bodyReader := bytes.NewReader(ml_body)
		res, err := http.Post(ML_SERVER+"/model/", "application/json", bodyReader)

		if err != nil {
			fmt.Println("something went wrong in making request")
			return
		}

		emb_resBody, err := io.ReadAll(res.Body)

		if err != nil {
			fmt.Println("something went wrong in reading body")
		}

		var titleEmbedding Embedding

		json.Unmarshal(emb_resBody, &titleEmbedding)
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

		rows, err := db.Query(`SELECT id, title FROM title_embeddings ORDER BY embedding <=> ($1) LIMIT 3`, pgvector.NewVector(titleEmbedding.Data))

		if err != nil {
			fmt.Printf("DB Error: %v\n", err)
			return
		}

		var similarTitles []SimilarTitle

		for rows.Next() {
			var smimilarT SimilarTitle
			fmt.Println(rows)
			if err := rows.Scan(&smimilarT.Idx, &smimilarT.Title); err != nil {
				fmt.Printf("Scanning error %v\n", err)
				return
			}
			similarTitles = append(similarTitles, smimilarT)
		}

		similarTitlesMessage := struct {
			Titles []SimilarTitle `json:"titles"`
		}{
			Titles: similarTitles,
		}

		similarTitlesBytes, err := json.Marshal(similarTitlesMessage)
		if err != nil {
			fmt.Printf("Problem with marshal: %v\n", err)
			return
		}
		w.Write(similarTitlesBytes)
	}
}

func main() {

	http.HandleFunc("/search", searchHandler)

	fmt.Printf("Server runs on: %s\n", GO_SERVER)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
