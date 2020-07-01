package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rpoletaev/huskyjam/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type goodsHandler struct {
	Store internal.GoodsRepository
	Log   zerolog.Logger
}

func (h *goodsHandler) logger(ctx context.Context) *zerolog.Logger {
	id, ok := hlog.IDFromCtx(ctx)
	var l zerolog.Logger
	if ok {
		l = h.Log.With().Str("req_id", id.String()).Logger()
		return &l
	}

	l = h.Log.With().Logger()
	return &l
}

type createGoodRequest struct {
	Name       string `json:"name"`
	Categories []uint `json:"categories"`
}

func (req *createGoodRequest) Validate() error {
	if len(req.Name) == 0 {
		return errors.Wrap(errValidateError, "length of good name shold be greater than 0")
	}
	return nil
}

func (h *goodsHandler) CreateGood(w http.ResponseWriter, r *http.Request) {
	logger := h.logger(r.Context())
	req := &createGoodRequest{}
	if err := unmarshal(r, req); err != nil {
		logger.Error().Err(err).Msg("on get createGoodRequest")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	good := &internal.Good{
		Name:       req.Name,
		Categories: make([]internal.Category, 0, len(req.Categories)),
	}

	for _, c := range req.Categories {
		good.Categories = append(good.Categories, internal.Category{ID: c})
	}

	if err := h.Store.CreateGood(good); err != nil {
		logger.Error().Err(err).Msg("on create new good")
		e := parseInternalError(err)
		http.Error(w, e.Message, e.Code)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

type updateGoodRequest struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Categories []uint `json:"categories"`
}

func (req *updateGoodRequest) Validate() error {
	if req.ID <= 0 {
		return errors.Wrap(errValidateError, "id should be greather than 0")
	}
	if len(req.Name) == 0 {
		return errors.Wrap(errValidateError, "length of good name shold be greater than 0")
	}
	return nil
}

func (h *goodsHandler) UpdateGood(w http.ResponseWriter, r *http.Request) {
	logger := h.logger(r.Context())
	req := &updateGoodRequest{}
	if err := unmarshal(r, req); err != nil {
		logger.Error().Err(err).Msg("on get updateGoodRequest")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	good := &internal.Good{
		ID:         req.ID,
		Name:       req.Name,
		Categories: make([]internal.Category, 0, len(req.Categories)),
	}

	for _, c := range req.Categories {
		good.Categories = append(good.Categories, internal.Category{ID: c})
	}

	if err := h.Store.UpdateGood(good); err != nil {
		logger.Error().Err(err).Uint("id", req.ID).Msg("on update good")
		e := parseInternalError(err)
		http.Error(w, e.Message, e.Code)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Good struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}

type categoryGoodsResponse struct {
	Goods []*Good `json:"goods"`
}

func (h *goodsHandler) CategoryGoods(w http.ResponseWriter, r *http.Request) {
	logger := h.logger(r.Context())
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["category"])
	if err != nil {
		logger.Error().Err(err).Str("category_id", vars["category"]).Msg("on convert category id")
		http.Error(w, "category_id should be an integer", http.StatusBadRequest)
		return
	}

	goods, err := h.Store.GoodsByCategory(uint(categoryID))
	if err != nil {
		logger.Error().Err(err).Int("category_id", categoryID).Msg("on get category goods")
		e := parseInternalError(err)
		http.Error(w, e.Message, e.Code)
		return
	}

	resp := &categoryGoodsResponse{
		Goods: GoodListFromModel(goods),
	}

	writeJSON(w, resp)
}

func (h *goodsHandler) DeleteGood(w http.ResponseWriter, r *http.Request) {
	logger := h.logger(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.Error().Err(err).Str("id", vars["id"]).Msg("on convert removable good id")
		http.Error(w, "id should be an integer", http.StatusBadRequest)
		return
	}

	if err := h.Store.DeleteGood(uint(id)); err != nil {
		logger.Error().Err(err).Int("good_id", id).Msg("on delete good")
		e := parseInternalError(err)
		http.Error(w, e.Message, e.Code)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type categoriesResponse struct {
	Categories []*Category `json:"categories"`
}

func (h *goodsHandler) CategoriesList(w http.ResponseWriter, r *http.Request) {
	logger := h.logger(r.Context())
	categories, err := h.Store.ListCategories()
	if err != nil {
		logger.Error().Err(err).Msg("on get categories list")
		e := parseInternalError(err)
		http.Error(w, e.Message, e.Code)
		return
	}

	resp := &categoriesResponse{
		Categories: CategoriesListFromModel(categories),
	}

	writeJSON(w, resp)
}

type createCategoryRequest struct {
	Name string
}

func (c *createCategoryRequest) Validate() error {
	if len(c.Name) == 0 {
		return errors.Wrap(errValidateError, "length of category name shold be greater than 0")
	}
	return nil
}

func (h *goodsHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	logger := h.logger(r.Context())
	req := &createCategoryRequest{}
	if err := unmarshal(r, req); err != nil {
		logger.Error().Err(err).Msg("on get createGoodRequest")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cat := &internal.Category{
		Name: req.Name,
	}

	if err := h.Store.CreateCategory(cat); err != nil {
		logger.Error().Err(err).Msg("on create new category")
		e := parseInternalError(err)
		http.Error(w, e.Message, e.Code)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

type updateCategoryRequest struct {
	Category *Category
}

func (req *updateCategoryRequest) Validate() error {
	if req.Category == nil {
		return fmt.Errorf("empty request")
	}

	if req.Category.ID < 1 {
		return fmt.Errorf("wrong category id")
	}

	if len(req.Category.Name) == 0 {
		return fmt.Errorf("length of category name should be greater than 0")
	}
	return nil
}

func (h *goodsHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	logger := h.logger(r.Context())
	req := &updateCategoryRequest{}
	if err := unmarshal(r, req); err != nil {
		logger.Error().Err(err).Msg("on parse updateCategoryRequest")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category := CategoryToModel(req.Category)

	if err := h.Store.UpdateCategory(category); err != nil {
		logger.Error().Err(err).Uint("id", req.Category.ID).Msg("on update category")
		e := parseInternalError(err)
		http.Error(w, e.Message, e.Code)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *goodsHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	logger := h.logger(r.Context())
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.Error().Err(err).Str("id", vars["id"]).Msg("on convert removable category id")
		http.Error(w, "id should be an integer", http.StatusBadRequest)
		return
	}

	if err := h.Store.DeleteCategory(uint(id)); err != nil {
		logger.Error().Err(err).Int("category_id", id).Msg("on delete category")
		e := parseInternalError(err)
		http.Error(w, e.Message, e.Code)
		return
	}
}
