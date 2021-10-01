package main

import (
  "encoding/json"
  "github.com/gorilla/schema"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "log"
  "net/http"
  "os"
)

type Table struct {
  Hello  string `json:"hello" schema:"hello"`
  Name    string `json:"name" schema:"name"`
  gorm.Model
}

func main(){

  //CLEARDB_DATABASE_URL
  dsn := os.Getenv("CLEARDB_DATABASE_URL")
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  if err != nil {
    log.Fatalf(err.Error())
  }

  er := db.AutoMigrate(&Table{})
  if er != nil{
    log.Fatalf(er.Error())
  }


  
  http.HandleFunc("/api/add", func (w http.ResponseWriter, r *http.Request){
    var t Table
    err := schema.NewDecoder().Decode(&t, r.Form)
    if err != nil {
      http.Error(w, err.Error(), 400)
      return
    }

    err = db.Model(Table{}).Where("hello = ?", t.Hello).Updates(&t).FirstOrCreate(&t).First(&t).Error
    if err != nil {
      http.Error(w, err.Error(), 400)
      return
    }
    w.Header().Set("content-type", "application/json")
    json.NewEncoder(w).Encode(t)
  })
  
   http.HandleFunc("/api/get", func (w http.ResponseWriter, r *http.Request){
     var t []Table

     err = db.Model(Table{}).Where("hello = ?", r.FormValue("hello")).Find(&t).Error
     if err != nil {
       http.Error(w, err.Error(), 400)
       return
     }
     w.Header().Set("content-type", "application/json")
     json.NewEncoder(w).Encode(t)
  })
  http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
