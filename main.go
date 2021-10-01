package main

import (
  "encoding/json"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "log"
  "net/http"
  "os"
)

type Table struct {
  PowerStatus  bool `json:"powerStatus" schema:"powerStatus"`
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
    s := r.FormValue("powerStatus")
    if s != "0" && s != "1" {
      http.Error(w, "powerStatus is invalid", 400)
      return
    }

    switch s {
    case "0":
      t.PowerStatus = false
    default:
      t.PowerStatus = true
    }

    err = db.Model(Table{}).Save(&t).First(&t).Error
    if err != nil {
      http.Error(w, err.Error(), 400)
      return
    }
    w.Header().Set("content-type", "application/json")
    json.NewEncoder(w).Encode(t)
  })
  
   http.HandleFunc("/api/get", func (w http.ResponseWriter, r *http.Request){
     var t []Table
     err = db.Model(Table{}).Find(&t).Error
     if err != nil {
       http.Error(w, err.Error(), 400)
       return
     }
     w.Header().Set("content-type", "application/json")
     json.NewEncoder(w).Encode(t)
  })
  http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
