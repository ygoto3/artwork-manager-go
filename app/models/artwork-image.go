package models

import (
  "github.com/revel/revel"
)

type ArtworkImage struct {
  Id              int64 `db:"id" json:"id"`
  Path            string `db:"path" json:"path"`
}

func (m *ArtworkImage) Validate(v *revel.Validation) {

  v.Check(m.Path,
    revel.ValidRequired(),
    revel.ValidMaxSize(100))

}
