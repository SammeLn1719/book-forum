package handler

import (
	//"encoding/json"
	"net/http"
)

//type Book struct {
//	ID 	   int     `json:"id"`
//	Title  string  `json:"title"`
//	Author string  `json:"authar"`
//}

func GetBooksHandler(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Getbooks"))
}
