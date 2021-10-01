package main

import (
  "github.com/gin-gonic/gin"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "log"
  "net/http"
  "os"
  "time"
)

type Table struct {
  ID        uint `gorm:"primary_key"`
  PowerStatus  bool `json:"powerStatus" schema:"powerStatus"`
  CreatedAt time.Time
  UpdatedAt time.Time `json:"_"`
  DeletedAt *time.Time `json:"_" sql:"index"`
}

func main(){

  //CLEARDB_DATABASE_URL
  dsn := os.Getenv("DSN")
  if dsn == "" {
   //dsn = "b0c18710b783ac:7e825d69@eu-cdbr-west-01.cleardb.com/heroku_25194f2ea429b61?reconnect=true"
   dsn = "b0c18710b783ac:7e825d69@tcp(eu-cdbr-west-01.cleardb.com)/heroku_25194f2ea429b61?parseTime=true"
    }

  db, err := gorm.Open( "mysql", dsn)
  if err != nil {
    log.Fatalf(err.Error())
  }

  er := db.AutoMigrate(&Table{})
  if er.Error != nil{
    log.Fatalf(er.Error.Error())
  }

  r := gin.Default()

  r.GET("/api/hiii", func(ctx *gin.Context) {
    ctx.JSON(200 ,  gin.H{"hi": "successful"})
    return
  })
  r.GET("/api/add", add(db))
  r.GET("/api/get",get(db))

  port := ":"+os.Getenv("PORT")
  if port == ":" {
    port = ":3000"
  }

  log.Printf("server running on %s", port)
  err  = r.Run(port)
    if err != nil {
      log.Fatalf(err.Error())
    }

    log.Println("server stoping")

}

func get(db *gorm.DB) func(ctx *gin.Context) {
  return func (ctx *gin.Context){
    w := ctx.Writer

    t := make([]Table,1)
    err := db.Model(Table{}).Find(&t).Error
    if err != nil {
      http.Error(w, err.Error(), 400)
      return
    }
    w.WriteHeader(200)
    w.Header().Set("content-type", "application/json")
    ctx.JSON(200, t)
  }
}

func add(db *gorm.DB) func(ctx *gin.Context) {
  return func(ctx *gin.Context) {
    w := ctx.Writer
    r := ctx.Request
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

    err := db.Model(Table{}).Create(&t).First(&t).Error
    if err != nil {
      http.Error(w, err.Error(), 400)
      return
    }

    w.WriteHeader(201)
    w.Header().Set("content-type", "application/json")
    ctx.JSON(200, t)
  }

}