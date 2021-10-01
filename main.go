package main

import (
  "net/http"
  "os"
)

func main(){

  //CLEARDB_DATABASE_URL
  http.HandlerFunc("/api/add", func (){
    w.Write([]byte("hello add..."))
  })
  
   http.HandlerFunc("/api/get", func (){
    w.Write([]byte("hello get..."))
  })
  http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
