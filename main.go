package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

// Verse : struct for json
type Verse struct {
	Code1   string `json:"code1"`
	Code2   string `json:"code2"`
	Book1   string `json:"book1"`
	Book2   string `json:"book2"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Words1  string `json:"words1"`
	Words2  string `json:"words2"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{code}/{chapter}/{from}/{to}", multiVerseHandler)
	r.HandleFunc("/{code}/{chapter}/{verse}", oneVerseHandler)
	r.HandleFunc("/{code}/{chapter}", chapterHandler)
	r.HandleFunc("/", homeHandler)

	handler := cors.Default().Handler(r)
	fmt.Println("Bible API server started ... ")
	http.ListenAndServe(":3001", handler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	usage := ` Usage:
   - localhost:3000/Gn/1     : Genesis Chapter 1
   - localhost:3000/Gn/1/1   : Genesis Chapter 1 Verse 1
   - localhost:3000/Gn/1/1/5 : Genesis Chapter 1 Verse 1 ~ 5`
	fmt.Fprintf(w, "%s\n", usage)
}

func multiVerseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	verses := getVerses(vars["code"], vars["chapter"], vars["from"], vars["to"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(verses)
}

func oneVerseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	verse := getVerses(vars["code"], vars["chapter"], vars["verse"], vars["verse"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(verse)
}

func chapterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	verse := getVerses(vars["code"], vars["chapter"], "1", "1000")
	json.NewEncoder(w).Encode(verse)
}

func getVerses(bookcode string, chap string, from string, to string) []Verse {
	db, err := sql.Open("sqlite3", "./bible.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql := `
	select BOOK.e_abbv1 as e_code, BOOK.k_abbv as k_code, KJV.book_name as e_bookname, HKJV.book_name as k_bookname, 
       KJV.chapter, KJV.verse, KJV.words as e_words, HKJV.words as k_words
    from KJV, HKJV, BOOK 
	where KJV.book_id = BOOK.book_id and 
	   KJV.book_id = HKJV.book_id and KJV.chapter = HKJV.chapter and KJV.verse = HKJV.verse and
	   BOOK.e_abbv1 = '%s' and KJV.chapter = %s and KJV.verse >= %s and KJV.verse <= %s`

	rows, err := db.Query(fmt.Sprintf(sql, bookcode, chap, from, to))
	if err != nil {
		log.Fatal(err)
	}

	var vv []Verse
	var v Verse
	for rows.Next() {
		err := rows.Scan(&v.Code1, &v.Code2, &v.Book1, &v.Book2,
			&v.Chapter, &v.Verse, &v.Words1, &v.Words2)
		if err != nil {
			log.Fatal(err)
		}
		vv = append(vv, v)
	}
	return vv
}
