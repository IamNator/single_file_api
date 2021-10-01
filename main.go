package main

import (
  "net/http"
)

func main(){

  http.HandlerFunc("/api/add", func (){
    w.Write([]byte("hello ..."))
  })
  http.ListenAndServe(":4500", nil)
}
