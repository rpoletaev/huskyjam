package pg

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/rpoletaev/huskyjam/internal"
)

const (
	initCategories = `CREATE TABLE IF NOT EXISTS categories (
		id integer PRIMARY KEY,
		name text NOT NULL,
		created_at timestamp DEFAULT current_timestamp
		deleted_at timestamp 
		CONSTRAINT unique_cat_name UNIQUE(name)
		)`

	initGoods = `CREATE TABLE IF NOT EXISTS goods (
		id string PRIMARY KEY,
		name text NOT NULL,
		created_at timestamp DEFAULT current_timestamp
		deleted_at timestamp 
		CONSTRAINT unique_good_name UNIQUE(name)
		)`

	initGoodsToCats = `CREATE TABLE IF NOT EXISTS goods_to_cats (
		good_id integer NOT NULL,
		category_id integer NOT NULL,
		CONSTRAINT unique_good_cat UNIQUE(category_id, good_id),
		CONSTRAINT good_to_cat_good_fk FOREIGN KEY (good_id) REFERENCES goods (id),
		CONSTRAINT good_to_cat_cat_fk FOREIGN KEY (good_id) REFERENCES categories (id)
		)`

	initGoodsIndex    = `CREATE INDEX IF NOT EXISTS ON goods_to_cats (good_id)`
	initCategoryIndex = `CREATE INDEX IF NOT EXISTS ON goods_to_cats (category_id)`
)

type GoodsRepository Store

func (s *Store) Goods() internal.GoodsRepository {
	return (*GoodsRepository)(s)
}

var _ internal.GoodsRepository = (*GoodsRepository)(nil)

func (s *GoodsRepository) Init() error {
	if _, err := s.db.Exec(initCategories); err != nil {
		return err
	}
	if _, err := s.db.Exec(initGoods); err != nil {
		return err
	}
	if _, err := s.db.Exec(initGoodsToCats); err != nil {
		return err
	}
	if _, err := s.db.Exec(initGoodsIndex); err != nil {
		return err
	}

	return nil
}

func (s *GoodsRepository) CreateCategory(c *internal.Category) error {

	if _, err := s.db.Exec("INSERT INTO categories (name) VALUES (:email)", c); err != nil {
		if uniqueViolation(err) {
			return internal.ErrAlreadyExists
		}

		return accessDB(err)
	}

	return nil

}

func (s *GoodsRepository) UpdateCategory(c *internal.Category) error {
	if _, err := s.db.Exec("UPDATE categories SET name = :name WHERE id = :id", c); err != nil {
		if uniqueViolation(err) {
			return internal.ErrAlreadyExists
		}

		return accessDB(err)
	}
	return nil
}

func (s *GoodsRepository) DeleteCategory(id uint) (txError error) {
	tx, err := s.db.Begin()
	if err != nil {
		return accessDB(err)
	}

	defer func() {
		if txError == nil {
			txError = tx.Commit()
			return
		}

		tx.Rollback()
	}()

	pars := map[string]interface{}{
		"deleted": time.Now(),
		"id":      id,
	}

	if _, err := tx.Exec("UPDATE categories SET deleted_at = :deleted WHERE id = :id", pars); err != nil {
		txError = accessDB(err)
		return
	}

	return deleteGoodCatsByCategory(tx, id)
}

func deleteGoodCatsByCategory(tx *sql.Tx, categoryID uint) error {
	_, err := tx.Exec("DELETE FROM goods_to_cats WHERE category_id = $1", categoryID)
	return errors.Wrapf(err, "delete link between goods and category: %d", categoryID)
}

func (s *GoodsRepository) ListCategories() ([]*internal.Category, error) {
	list := []*internal.Category{}
	if err := s.db.Select(&list, "SELECT * from categories WHERE deleted_at IS NULL"); err != nil {
		if notFound(err) {
			return nil, internal.ErrNotFound
		}
		return nil, accessDB(err)
	}
	return list, nil
}

func (s *GoodsRepository) CreateGood(g *internal.Good) (txError error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if txError == nil {
			txError = tx.Commit()
			return
		}

		tx.Rollback()
	}()

	var newId uint
	pars := map[string]interface{}{
		"name":    g.Name,
		"created": time.Now(),
	}
	if err := tx.QueryRow("INSERT INTO goods (name, created_at) VALUES(:name, :created) RETURNING id", pars).Scan(&newId); err != nil {
		if uniqueViolation(err) {
			txError = internal.ErrAlreadyExists
			return
		}

		txError = accessDB(err)
		return
	}

	g.ID = newId
	return linkGoodToCategories(tx, g)
}

func linkGoodToCategories(tx *sql.Tx, good *internal.Good) error {
	if len(good.Categories) == 0 {
		return nil
	}

	for _, category := range good.Categories {
		prs := map[string]interface{}{
			"goodId": good.ID,
			"catId":  category.ID,
		}
		if _, err := tx.Exec("INSERT INTO goods_to_cats (good_id, category_id) VALUES (:goodId, :catId)", prs); err != nil {
			if uniqueViolation(err) {
				return errors.Wrapf(internal.ErrAlreadyExists, "link good to category: %d", category.ID)
			}

			return accessDB(err)
		}
	}

	return nil
}

func (s *GoodsRepository) UpdateGood(g *internal.Good) (txError error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if txError == nil {
			txError = tx.Commit()
			return
		}

		tx.Rollback()
	}()

	pars := map[string]interface{}{
		"id":   g.ID,
		"name": g.Name,
	}

	if _, err := tx.Exec("UPDATE goods SET name = :name WHERE id = :id", pars); err != nil {
		if notFound(err) {
			return errors.Wrapf(internal.ErrNotFound, "good with id: %d", g.ID)
		}

		if uniqueViolation(err) {
			return errors.Wrapf(internal.ErrAlreadyExists, "good with name: %s", g.Name)
		}

		return accessDB(err)
	}

	if err := deleteCategoriesForGood(tx, g.ID); err != nil {
		return errors.Wrapf(err, "delete categories for good: %d", g.ID)
	}

	if err := linkGoodToCategories(tx, g); err != nil {
		return err
	}

	return nil
}

func deleteCategoriesForGood(tx *sql.Tx, goodID uint) error {
	if _, err := tx.Exec("DELETE FROM goods_to_cats WHERE good_id = $1", goodID); err != nil {
		return accessDB(err)
	}

	return nil
}

func (s *GoodsRepository) DeleteGood(id uint) (txError error) {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if txError == nil {
			txError = tx.Commit()
			return
		}

		tx.Rollback()
	}()

	pars := map[string]interface{}{
		"id":      id,
		"deleted": time.Now(),
	}
	if _, err := tx.Exec("UPDATE goods SET deleted_at = :deleted WHERE id = :id", pars); err != nil {
		if notFound(err) {
			return errors.Wrapf(internal.ErrNotFound, "good with id: %d", id)
		}

		return accessDB(err)
	}

	return deleteCategoriesForGood(tx, id)
}

func (s *GoodsRepository) GoodsByCategory(catID uint) ([]*internal.Good, error) {
	goods := []*internal.Good{}
	query := `SELECT g.* FROM goods g 
		JOIN goods_to_categories gtc
		ON g.id = gtc.good_id AND gtc.category_id = $1`

	if err := s.db.Select(&goods, query, catID); err != nil {
		return nil, accessDB(err)
	}

	return goods, nil
}
