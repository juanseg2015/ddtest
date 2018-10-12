# ddtest
this is a code for evaluation

## precondition
- glide install 
or
- go-sqlite3
    go get github.com/mattn/go-sqlite3
pat (for routing)
    go get github.com/bmizerany/pat

## REST API
## url
localhost:12345
## return all songs
- GET songs
## return songs by artist
- GET artist/:artist
## return  songs by song name
- GET song/:song
## return songs by genre
- GET genre/:genre
## return song between min a max length
- GET length/:min/:max
## return totals by genre
- GET totals/


