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
  ID        uint `gorm:"primary_key" json:"id"`
  PowerStatus  bool `json:"power_status" schema:"power_status"`
  CreatedAt time.Time  `json:"created_at"`
  UpdatedAt time.Time `json:"_"`
  DeletedAt *time.Time `json:"_" sql:"index"`
}

func main(){

  //CLEARDB_DATABASE_URL
  dsn := os.Getenv("DSN")
  if dsn == "" {
    panic("DSN not found in environment")
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
    s := r.FormValue("power_status")
    if s != "0" && s != "1" {
      http.Error(w, "power_status is invalid", 400)
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
