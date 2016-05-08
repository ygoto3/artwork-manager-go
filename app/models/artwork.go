package models

import (
  "github.com/revel/revel"
)

type Artwork struct {
  Id              int64 `db:"id" json:"id"`
  Title           string `db:"title" json:"title"`
  Description     string `db:"description" json:"description"`
  ArtworkImageId  int64 `db:"artwork_image_id" json:"artwork_image_id"`
}

type ArtworkWithImage struct {
  Artwork
  ArtworkImage    ArtworkImage `json:"artwork_image"`
}

func (m *Artwork) Validate(v *revel.Validation) {

  v.Check(m.Title,
    revel.ValidRequired(),
    revel.ValidMaxSize(25))

}
