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

// TODO: handle db connection properly
// make dynamic number of similar films in request

type Embedding struct {
	Idx  int       `json:"id"`
	Data []float32 `json:"embedding"`
}

type SimilarTitle struct {
	Idx   int    `json:"id"`
	Title string `json:"title"`
}

type SearchbarInput struct {
	Data string `json:"searchbar_input"`
}

var fetchedEmbeddings map[int]string

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		resBody, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Println("something went wrong in reading body")
			defer w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var searchbarContent SearchbarInput

		json.Unmarshal(resBody, &searchbarContent)

		ml_body := []byte(fmt.Sprintf("{\"movie_title\": \"%s\"}", searchbarContent.Data))
		bodyReader := bytes.NewReader(ml_body)
		res, err := http.Post(ML_SERVER+"/model/", "application/json", bodyReader)

		if err != nil {
			defer w.WriteHeader(http.StatusInternalServerError)
			return
		}

		emb_resBody, err := io.ReadAll(res.Body)

		if err != nil {
			fmt.Println("something went wrong in reading body")
			defer w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var titleEmbedding Embedding

		json.Unmarshal(emb_resBody, &titleEmbedding)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			defer w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			defer w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rows, err := db.Query(`SELECT id, title, embedding FROM title_embeddings ORDER BY embedding <=> ($1) LIMIT 3`, pgvector.NewVector(titleEmbedding.Data))

		if err != nil {
			fmt.Printf("DB Error: %v\n", err)
			defer w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var similarTitles []SimilarTitle
		fetchedEmbeddings = make(map[int]string)

		for rows.Next() {
			var smimilarT SimilarTitle
			var str_embedding string
			if err := rows.Scan(&smimilarT.Idx, &smimilarT.Title, &str_embedding); err != nil {
				fmt.Printf("Scanning error %v\n", err)
				defer w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			fetchedEmbeddings[smimilarT.Idx] = str_embedding
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
			defer w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		w.Write(similarTitlesBytes)
	}
	defer w.WriteHeader(http.StatusBadRequest)
}

func main() {

	http.HandleFunc("/search", searchHandler)
	fmt.Printf("Server runs on: %s\n", GO_SERVER)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
