package controllers

import (
  "github.com/revel/revel"
  "artwork-manager/app/models"
  "encoding/json"
  "strconv"
)

type Artworks struct {
  GorpController
}

func (c Artworks) parseArtwork() (models.Artwork, error) {
  artwork := models.Artwork{}
  err := json.NewDecoder(c.Request.Body).Decode(&artwork)
  return artwork, err
}

func (c Artworks) List() revel.Result {
  query := `SELECT * FROM Artwork WHERE id > ? LIMIT ?`
  artworks, err := c.Txn.Select(models.Artwork{}, query, 0, 10)

  if err != nil {
    return c.RenderText("Error trying to get records from DB.")
  }

  query = `SELECT * FROM ArtworkImage WHERE`
  for i, v := range artworks {
    aw := v.(*models.Artwork)
    if i > 0 {
      query = query + " OR"
    }
    query = query + " id = " + strconv.FormatInt(aw.ArtworkImageId, 10)
  }

  artworkImages, err := c.Txn.Select(models.ArtworkImage{}, query)
  if err != nil {
    return c.RenderText("Error trying to get records from DB.")
  }

  artworkWithImages := make([]models.ArtworkWithImage, len(artworks))
  for _, v := range artworkImages {
    awi := v.(*models.ArtworkImage)
    for j, w := range artworks {
      aw := w.(*models.Artwork)
      artworkWithImages[j].Artwork = *aw
      if awi.Id == aw.ArtworkImageId {
        artworkWithImages[j].ArtworkImage = *awi
      }
    }
  }
  return c.RenderJson(artworkWithImages)
}

func (c Artworks) Show(artworkId int64) revel.Result {
  artwork, err := c.Txn.Get(models.Artwork{}, artworkId)
  if err != nil {
    return c.RenderText("Error trying to get record from DB.")
  }

  aw := artwork.(*models.Artwork)
  artworkImageId := aw.ArtworkImageId
  artworkImage, err := c.Txn.Get(models.ArtworkImage{}, artworkImageId)

  if err != nil {
    return c.RenderText("Error trying to get record from DB.")
  }

  artworkWithImage := models.ArtworkWithImage{Artwork: *aw, ArtworkImage: *artworkImage.(*models.ArtworkImage)}

  return c.RenderJson(artworkWithImage)
}

func (c Artworks) Create() revel.Result {
  if artwork, err := c.parseArtwork(); err != nil {
    return c.RenderText("Unable to parse the ")
  } else {
    artwork.Validate(c.Validation)
    if c.Validation.HasErrors() {
      return c.RenderText("You have error in your artwork")
    } else {
      if err := c.Txn.Insert(&artwork); err != nil {
        return c.RenderText("Error inserting record into database!")
      } else {
        return c.RenderJson(artwork)
      }
    }
  }
}

func (c Artworks) Update(artworkId int64) revel.Result {
  artwork, err := c.parseArtwork()
  if err != nil {
    return c.RenderText("Unable to parse the Artwork from JSON.")
  }
  artwork.Id = artworkId
  success, err := c.Txn.Update(&artwork)
  if err != nil || success == 0 {
    return c.RenderText("Unable to update the artwork")
  }
  return c.RenderText("Updated %v", artworkId)
}

func (c Artworks) Destroy(artworkId int64) revel.Result {
  success, err := c.Txn.Delete(&models.Artwork{Id: artworkId})
  if err != nil || success == 0 {
    c.RenderText("Failed to remove Artwork")
  }
  return c.RenderText("Deleted %v", artworkId)
}
