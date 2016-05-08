package controllers

import (
  "github.com/revel/revel"
  "github.com/coopernurse/gorp"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "fmt"
  "strings"
  "artwork-manager/app/models"
)

func init() {
  revel.OnAppStart(InitDb)
  revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
  revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
  revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}

func getParamString(param string, defaultValue string) string {
  p, found := revel.Config.String(param)
  if !found {
    if defaultValue == "" {
      revel.ERROR.Fatal("Cound not find parameter: " + param)
    } else {
      return defaultValue
    }
  }
  return p
}

func getConnectionString() string {
  host := getParamString("db.host", "")
  port := getParamString("db.port", "3306")
  user := getParamString("db.user", "")
  dbname := getParamString("db.name", "auction")
  protocol := getParamString("db.protocol", "tcp")
  dbargs := getParamString("dbargs", " ")

  if strings.Trim(dbargs, " ") != "" {
    dbargs = "?" + dbargs
  } else {
    dbargs = ""
  }
  return fmt.Sprintf("%s@%s([%s]:%s)/%s%s", 
  user, protocol, host, port, dbname, dbargs)
}

var InitDb func() = func() {
  connectionString := getConnectionString()
  if db, err := sql.Open("mysql", connectionString); err != nil {
    revel.ERROR.Fatal(err)
  } else {
    Dbm = &gorp.DbMap{ Db: db,
      Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"} }
  }

  defineArtworkTable(Dbm)
  defineArtworkImageTable(Dbm)
  if err := Dbm.CreateTablesIfNotExists(); err != nil {
    revel.ERROR.Fatal(err)
  }
}

func defineArtworkTable(dbm *gorp.DbMap) {
  t := dbm.AddTable(models.Artwork{}).SetKeys(true, "id")
  t.ColMap("title").SetMaxSize(25)
}

func defineArtworkImageTable(dbm *gorp.DbMap) {
  t := dbm.AddTable(models.ArtworkImage{}).SetKeys(true, "id")
  t.ColMap("path").SetMaxSize(100)
}
