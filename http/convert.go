package http

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/rpoletaev/huskyjam/internal"
)

const (
	contentTypeKey   = "content-type"
	contentTypeValue = "application/json"
)

var (
	errValidateError error = errors.New("validation error")
)

type Validator interface {
	Validate() error
}

func unmarshal(r *http.Request, v Validator) error {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	if err := json.Unmarshal(raw, v); err != nil {
		return err
	}

	return v.Validate()
}

func writeJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Add(contentTypeKey, contentTypeValue)
	return json.NewEncoder(w).Encode(v)
}

func CategoryToModel(c *Category) *internal.Category {
	return &internal.Category{
		ID:   c.ID,
		Name: c.Name,
	}
}

func CategoryFromModel(c *internal.Category) *Category {
	return &Category{
		ID:   c.ID,
		Name: c.Name,
	}
}

func CategoriesListFromModel(list []*internal.Category) []*Category {
	res := make([]*Category, 0, len(list))
	for _, c := range list {
		res = append(res, CategoryFromModel(c))
	}
	return res
}

func GoodToModel(g *Good) *internal.Good {
	good := &internal.Good{
		ID:         g.ID,
		Name:       g.Name,
		Categories: make([]internal.Category, 0, len(g.Categories)),
	}

	for _, c := range g.Categories {
		good.Categories = append(good.Categories, *CategoryToModel(&c))
	}

	return good
}

func GoodFromModel(g *internal.Good) *Good {
	good := &Good{
		ID:         g.ID,
		Name:       g.Name,
		Categories: make([]Category, 0, len(g.Categories)),
	}

	for _, c := range g.Categories {
		good.Categories = append(good.Categories, *CategoryFromModel(&c))
	}

	return good
}

func GoodListFromModel(list []*internal.Good) []*Good {
	res := make([]*Good, 0, len(list))
	for _, g := range list {
		res = append(res, GoodFromModel(g))
	}
	return res
}
