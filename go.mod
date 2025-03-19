module github.com/avgra3/gator

go 1.23.0

require internal/config v1.0.0

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

replace internal/config => ./internal/config
