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

// TODO: handle db connection properly
// make dynamic number of similar films in request

type searchHandler struct {
	db *sql.DB
}

type Embedding struct {
	Idx  int       `json:"id"`
	Data []float32 `json:"embedding"`
}

type SimilarBook struct {
	Idx     int    `json:"id"`
	Title   string `json:"title"`
	ImgPath string `json:"img_path"`
}

type SearchbarInput struct {
	Data       string `json:"searchbar_input"`
	SearchType string `json:"search_type"`
}

func sendSearchResults(w *http.ResponseWriter, rows *sql.Rows) {

	var similarBooks []SimilarBook

	for rows.Next() {
		var smimilarB SimilarBook
		if err := rows.Scan(&smimilarB.Idx, &smimilarB.Title, &smimilarB.ImgPath); err != nil {
			fmt.Printf("Scanning error %v\n", err)
			defer (*w).WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		similarBooks = append(similarBooks, smimilarB)
	}

	similarTitlesMessage := struct {
		Books []SimilarBook `json:"similar_books"`
	}{
		Books: similarBooks,
	}

	similarTitlesBytes, err := json.Marshal(similarTitlesMessage)
	if err != nil {
		fmt.Printf("Problem with marshal: %v\n", err)
		defer (*w).WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	(*w).Write(similarTitlesBytes)
}

func getSearchbarInputEmbedding(w *http.ResponseWriter, data string) (Embedding, error) {
	ml_body := []byte(fmt.Sprintf("{\"book_title\": \"%s\"}", data))
	bodyReader := bytes.NewReader(ml_body)
	res, err := http.Post(ML_SERVER+"/model/", "application/json", bodyReader)

	if err != nil {
		defer (*w).WriteHeader(http.StatusInternalServerError)
		return Embedding{}, err
	}

	emb_resBody, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println("something went wrong in reading body")
		defer (*w).WriteHeader(http.StatusUnprocessableEntity)
		return Embedding{}, err
	}

	var titleEmbedding Embedding
	json.Unmarshal(emb_resBody, &titleEmbedding)

	return titleEmbedding, nil
}

func (sh *searchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		resBody, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Println("something went wrong in reading body")
			defer w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var searchbarContent SearchbarInput
		err = json.Unmarshal(resBody, &searchbarContent)

		if err != nil {
			defer w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		var rows *sql.Rows
		switch searchbarContent.SearchType {
		case "semantic":

			titleEmbedding, err := getSearchbarInputEmbedding(&w, searchbarContent.Data)

			if err != nil {
				fmt.Printf("Search Input Embedding Error: %s\n", err)
				return
			}

			rows, err = sh.db.Query(`SELECT id, title, img_path FROM books ORDER BY title_embedding <=> ($1) LIMIT 3`, pgvector.NewVector(titleEmbedding.Data))

		case "exact":
			rows, err = sh.db.Query(`SELECT id, title, img_path FROM books WHERE lower(title) like  '%'||lower(($1))||'%' LIMIT 3`, searchbarContent.Data)
		default:
			fmt.Printf("Error: Type of search not specified\n")
			defer w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if err != nil {
			fmt.Printf("DB Error: %v\n", err)
			defer w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sendSearchResults(&w, rows)
		return
	}

	defer w.WriteHeader(http.StatusBadRequest)
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return
	}
	http.Handle("/search", &searchHandler{db: db})
	fmt.Printf("Server runs on: %s\n", GO_SERVER)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
