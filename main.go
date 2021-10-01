package main

import (
  "net/http"
  "os"
  "log"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

func main(){

  //CLEARDB_DATABASE_URL
  dsn := os.Getenv("CLEARDB_DATABASE_URL")
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  if err != nil {
    log.Println(err.Error())
  }
  
  http.HandleFunc("/api/add", func (w http.ResponseWriter, r *http.Request){
    w.Write([]byte("hello add..."))
  })
  
   http.HandleFunc("/api/get", func (w http.ResponseWriter, r *http.Request){
    w.Write([]byte("hello get..."))
  })
  http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
