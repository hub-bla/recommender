module movie_recommender/main

go 1.23

require github.com/lib/pq v1.10.9

require github.com/pgvector/pgvector-go v0.1.1

require example.com/http_request_handler v0.0.0-00010101000000-000000000000 // indirect

replace example.com/http_request_handler => ./http_request_handler
