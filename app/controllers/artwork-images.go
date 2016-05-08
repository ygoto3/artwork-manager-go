package controllers

import (
  "github.com/revel/revel"
  "artwork-manager/app/models"
  "os"
  "io"
  "strings"
  "strconv"
)

type ArtworkImages struct {
  GorpController
}

func getProjectPath() string {
  dir, _ := os.Getwd()
  return strings.Replace(dir, " ", "\\ ", -1)
}

func (c ArtworkImages) Create() revel.Result {
  err := c.Request.ParseMultipartForm(32 << 20)
  if err != nil {
    return c.RenderText("Error")
  }

  file, _, err := c.Request.FormFile("image")
  if err != nil {
    return c.RenderText("Error getting file")
  }
  defer file.Close()

  count, err := c.Txn.SelectInt(`SELECT COUNT(*) FROM ArtworkImage`)
  if err != nil {
    return c.RenderText("Error counting")
  }

  projPath := getProjectPath()
  filePath := "/artwork-images/" + strconv.FormatInt(count + 1, 10) + ".png"

  f, err := os.Create(projPath + "/app/storage/" + filePath)
  if err != nil {
    return c.RenderText("Error creating file")
  }
  defer f.Close()

  io.Copy(f, file)

  artworkImage := models.ArtworkImage{ Path: filePath }
  if err := c.Txn.Insert(&artworkImage); err != nil {
    return c.RenderText("Error inserting record into database!")
  } else {
    return c.RenderJson(artworkImage)
  }
}
