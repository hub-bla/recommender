package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	httprequesthandler "example.com/http_request_handler"
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

const DEFAULT_LIMIT = 50

type httpHandler struct {
	db     *sql.DB
	httpRH httprequesthandler.HTTPRequestHandler
}

type searchHandler struct {
	httpHandler
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
	Limit      int    `json:"limit"`
}

func searchResults(rows *sql.Rows) (struct {
	Books []SimilarBook "json:\"similar_books\""
}, error) {

	var similarBooks []SimilarBook

	for rows.Next() {
		var smimilarB SimilarBook
		if err := rows.Scan(&smimilarB.Idx, &smimilarB.Title, &smimilarB.ImgPath); err != nil {
			return struct {
				Books []SimilarBook `json:"similar_books"`
			}{}, fmt.Errorf("scanning error %v", err)
		}
		similarBooks = append(similarBooks, smimilarB)
	}

	similarBooksMessage := struct {
		Books []SimilarBook `json:"similar_books"`
	}{
		Books: similarBooks,
	}

	return similarBooksMessage, nil
}

func (hh *httpHandler) handleDbError(w *http.ResponseWriter) {
	message := "unable to retrieve data from database"
	errorMessage := httprequesthandler.ErrorMessage{
		Error:   "internal server error",
		Message: message,
	}
	hh.httpRH.SendResponse(w, errorMessage, message, http.StatusInternalServerError, "error")
}

func (hh *httpHandler) handleReadBodyError(w *http.ResponseWriter) {
	message := "could not read request body"
	errorMessage := httprequesthandler.ErrorMessage{
		Error:   "unprocessable entity",
		Message: message,
	}
	hh.httpRH.SendResponse(w, errorMessage, message, http.StatusUnprocessableEntity, "error")
}

func (hh *httpHandler) handleBadRequestError(w *http.ResponseWriter) {
	message := "wrong request method"
	errorMessage := httprequesthandler.ErrorMessage{
		Error:   "bad request",
		Message: message,
	}
	hh.httpRH.SendResponse(w, errorMessage, message, http.StatusBadRequest, "error")
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
			sh.handleReadBodyError(&w)
			return
		}

		searchbarContent := SearchbarInput{Limit: DEFAULT_LIMIT}
		err = json.Unmarshal(resBody, &searchbarContent)

		if err != nil || searchbarContent.Data == "" {
			message := "searchbar data is corrupted or missing"
			errorMessage := httprequesthandler.ErrorMessage{
				Error:   "unprocessable entity",
				Message: message,
			}
			sh.httpRH.SendResponse(&w, errorMessage, message, http.StatusUnprocessableEntity, "error")
			return
		}

		var rows *sql.Rows
		switch searchbarContent.SearchType {
		case "semantic":

			titleEmbedding, err := getSearchbarInputEmbedding(&w, searchbarContent.Data)

			if err != nil {
				message := "unable to retrieve embedding from server"
				errorMessage := httprequesthandler.ErrorMessage{
					Error:   "internal server error",
					Message: message,
				}
				sh.httpRH.SendResponse(&w, errorMessage, message, http.StatusInternalServerError, "error")
				return
			}

			rows, err = sh.db.Query(`SELECT 
									id, title, img_path 
									FROM books 
									ORDER BY title_embedding <=> ($1) 
									LIMIT ($2);`, pgvector.NewVector(titleEmbedding.Data), searchbarContent.Limit)

		case "exact":
			rows, err = sh.db.Query(`SELECT 
									id, title, img_path 
									FROM books 
									WHERE lower(title) like '%'||lower(($1))||'%' 
									LIMIT ($2);`, searchbarContent.Data, searchbarContent.Limit)
		default:
			message := "value for 'search_type' is not specified or incorrect"
			errorMessage := httprequesthandler.ErrorMessage{
				Error:   "bad request",
				Message: message,
			}
			sh.httpRH.SendResponse(&w, errorMessage, message, http.StatusBadRequest, "error")
			return
		}

		if err != nil {
			sh.handleDbError(&w)
			return
		}

		similarBooksMessage, err := searchResults(rows)
		if err != nil {
			errorMessage := httprequesthandler.ErrorMessage{
				Error:   "unprocessable entity",
				Message: fmt.Sprintf("%v", err),
			}
			sh.httpRH.SendResponse(&w, errorMessage, errorMessage.Message, http.StatusUnprocessableEntity, "error")
			return
		}

		sh.httpRH.SendResponse(&w, similarBooksMessage, "POST /search/ - sending similar Books", http.StatusOK, "success")
		return
	}
	sh.handleBadRequestError(&w)
}

type recommendHandler struct {
	httpHandler
}

type BookIdxMessage struct {
	Idx   string `json:"book_id"`
	Limit int    `json:"limit"`
}

func (rh *recommendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		resBody, err := io.ReadAll(r.Body)

		if err != nil {
			rh.handleReadBodyError(&w)
			return
		}

		bookIdxMessage := BookIdxMessage{Limit: DEFAULT_LIMIT}
		err = json.Unmarshal(resBody, &bookIdxMessage)

		if err != nil || bookIdxMessage.Idx == "" {
			message := "book id is corrupted or missing"
			errorMessage := httprequesthandler.ErrorMessage{
				Error:   "unprocessable entity",
				Message: message,
			}
			rh.httpRH.SendResponse(&w, errorMessage, message, http.StatusUnprocessableEntity, "error")
			return
		}

		rows, err := rh.db.Query(`
			SELECT 
			id, title, img_path 
			FROM books 
			where id != CAST(($1) AS INTEGER) 
			ORDER BY book_embedding <=> 
			(SELECT book_embedding 
			FROM books WHERE id = CAST(($1) AS INTEGER)) 
			LIMIT ($2);`, bookIdxMessage.Idx, bookIdxMessage.Limit)

		if err != nil {
			rh.handleDbError(&w)
			return
		}

		similarBooksMessage, err := searchResults(rows)

		if err != nil {
			errorMessage := httprequesthandler.ErrorMessage{
				Error:   "unprocessable entity",
				Message: fmt.Sprintf("%v", err),
			}
			rh.httpRH.SendResponse(&w, errorMessage, errorMessage.Message, http.StatusUnprocessableEntity, "error")
			return
		}

		rh.httpRH.SendResponse(&w, similarBooksMessage, "POST /recommend/ sending recommended Books", http.StatusOK, "success")
		return
	}
	rh.handleBadRequestError(&w)
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
	http.Handle("/search", &searchHandler{httpHandler{db: db}})
	http.Handle("/recommend", &recommendHandler{httpHandler{db: db}})
	fmt.Printf("Server runs on: %s\n", GO_SERVER)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
