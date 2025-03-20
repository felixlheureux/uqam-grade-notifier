package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/db"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/domain"
	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type BaseModel struct {
	bun.BaseModel
	db *bun.DB
}

func (m *BaseModel) SetDB(db *bun.DB) {
	m.db = db
}

func (m *BaseModel) GetDB() *bun.DB {
	return m.db
}

type imodel[T any, E any] interface {
	id() string
	idPrefix() string
	tableName() string
	postprocess(T) E
}

// T is the model type
// E is the entity type
// C is the create input type
// U is the update input type
// F is the filters type
type model[T imodel[T, E], E any, C any, U any, F any] struct {
	db        *bun.DB
	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time `bun:",soft_delete"`
}

func (m *model[T, E, C, U, F]) Get(id string, c context.Context, trx bun.IDB) (E, error) {
	if trx == nil {
		trx = m.db
	}

	var model T
	query := trx.NewSelect().Model(&model).Where("id = ?", id).Limit(1)

	err := query.Scan(c)

	switch err {
	case nil:
		return model.postprocess(model), nil
	case sql.ErrNoRows:
		return *new(E), domain.ErrNotFound(err)
	default:
		return *new(E), db.QueryExecuteError(err, fmt.Sprintf("%+v", query))
	}
}

func (m *model[T, E, C, U, F]) FindOne(filters F, c context.Context, trx bun.IDB) (E, error) {
	if trx == nil {
		trx = m.db
	}

	var model T
	query := trx.NewSelect().Model(&model)

	var mp map[string]interface{}
	f, _ := json.Marshal(filters)
	json.Unmarshal(f, &mp)

	for k, v := range mp {
		if v == nil {
			continue
		}
		query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := query.Limit(1).Scan(c)

	switch err {
	case nil:
		return model.postprocess(model), nil
	case sql.ErrNoRows:
		return *new(E), nil
	default:
		return *new(E), db.QueryExecuteError(err, fmt.Sprintf("%+v", query))
	}
}

func (m *model[T, E, C, U, F]) Find(filters F, pagination db.Pagination, c context.Context, trx bun.IDB) (db.Result[E], error) {
	if trx == nil {
		trx = m.db
	}

	pagination.ApplyDefaults()

	var models []T
	query := trx.NewSelect().Model(&models)

	var mp map[string]interface{}
	f, _ := json.Marshal(filters)
	json.Unmarshal(f, &mp)

	var model T
	db.AppendFiltersToQuery(query, model.tableName(), mp)

	if pagination.Cursor.Current != "" {
		cursorID, cursorDirection, err := db.DecodeCursor(pagination.Cursor.Current)
		if err != nil {
			// TODO: Return proper error
			return *new(db.Result[E]), err
		}
		if cursorDirection == "next" {
			query.Where("id > ?", cursorID)
		} else {
			query.Where("id < ?", cursorID)
			if pagination.OrderDirection == "ASC" {
				pagination.OrderDirection = "DESC"
			} else {
				pagination.OrderDirection = "ASC"
			}
		}
	}

	query.Order(fmt.Sprintf("%s %s", pagination.OrderBy, pagination.OrderDirection))
	err := query.Limit(pagination.Limit + 1).Scan(c)

	fmt.Println(query)
	var prevCursor, nextCursor string
	if len(models) > 0 {
		if pagination.Cursor.Current == "" || strings.HasPrefix(pagination.Cursor.Current, "next_") {
			prevCursor = db.EncodeCursor(models[0].id(), "prev")
		}
		if len(models) > pagination.Limit {
			nextCursor = db.EncodeCursor(models[len(models)-2].id(), "next")
			models = models[:pagination.Limit]
		}
	}

	pagination.Cursor.Prev = prevCursor
	pagination.Cursor.Next = nextCursor
	pagination.Cursor.Current = ""

	switch err {
	case nil:
		var results db.Result[E]
		results.Pagination = pagination
		for _, model := range models {
			results.Data = append(results.Data, model.postprocess(model))
		}
		return results, nil
	case sql.ErrNoRows:
		return *new(db.Result[E]), nil
	default:
		return *new(db.Result[E]), db.QueryExecuteError(err, fmt.Sprintf("%+v", query))
	}
}

// TODO: Add returning *
func (m *model[T, E, C, U, F]) Create(input C, c context.Context, trx bun.IDB) (E, error) {
	if trx == nil {
		trx = m.db
	}

	var model T
	id := fmt.Sprintf("%s%s", model.idPrefix(), ksuid.New().String())
	now := time.Now()

	var mp map[string]interface{}
	i, _ := json.Marshal(input)
	json.Unmarshal(i, &mp)

	mp["id"] = id
	mp["created_at"] = now
	mp["updated_at"] = now

	insertQuery := trx.NewInsert().
		Model(&mp).
		TableExpr(model.tableName())

	_, err := insertQuery.Exec(c)

	if err != nil {
		return *new(E), db.QueryExecuteError(err, fmt.Sprintf("%+v", insertQuery))
	}

	selectQuery := trx.NewSelect().Model(&model).Where("id = ?", id).Limit(1)

	err = selectQuery.Scan(c)

	switch err {
	case nil:
		return model.postprocess(model), nil
	case sql.ErrNoRows:
		return *new(E), domain.ErrNotFound(err)
	default:
		return *new(E), db.QueryExecuteError(err, fmt.Sprintf("%+v", selectQuery))
	}
}

// TODO: Add returning *
func (m *model[T, E, C, U, F]) Upsert(input C, conflict string, c context.Context, trx bun.IDB) (E, error) {
	if trx == nil {
		trx = m.db
	}

	var model T
	id := fmt.Sprintf("%s%s", model.idPrefix(), ksuid.New().String())
	now := time.Now()

	var mp map[string]interface{}
	i, _ := json.Marshal(input)
	json.Unmarshal(i, &mp)

	mp["id"] = id
	mp["created_at"] = now
	mp["updated_at"] = now

	insertQuery := trx.NewInsert().
		Model(&mp).
		On(fmt.Sprintf("CONFLICT (%s) DO NOTHING", conflict)).
		TableExpr(model.tableName())

	_, err := insertQuery.Exec(c)

	if err != nil {
		return *new(E), db.QueryExecuteError(err, fmt.Sprintf("%+v", insertQuery))
	}

	selectQuery := trx.NewSelect().Model(&model).Where(fmt.Sprintf("%s = ?", conflict), mp[conflict]).Limit(1)

	err = selectQuery.Scan(c)

	switch err {
	case nil:
		return model.postprocess(model), nil
	case sql.ErrNoRows:
		return *new(E), domain.ErrNotFound(err)
	default:
		return *new(E), db.QueryExecuteError(err, fmt.Sprintf("%+v", selectQuery))
	}
}

func (m *model[T, E, C, U, F]) Update(input U, c context.Context, trx bun.IDB) (E, error) {
	if trx == nil {
		trx = m.db
	}

	var model T
	query := trx.NewUpdate().
		Model(&model).
		Set("updated_at = ?", time.Now()).
		OmitZero().
		WherePK().
		Returning("*")

	res, err := query.Exec(c)

	if err != nil {
		return *new(E), db.QueryExecuteError(err, fmt.Sprintf("%+v", query))
	}

	affectedRows, _ := res.RowsAffected()

	if affectedRows == 0 {
		panic("Not found")
	}

	return model.postprocess(model), nil
}

func (m *model[T, E, C, U, F]) Delete(id string, c context.Context, trx bun.IDB) (string, error) {
	if trx == nil {
		trx = m.db
	}

	query := trx.NewUpdate().
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id)

	res, err := query.Exec(c)

	if err != nil {
		return "", db.QueryExecuteError(err, fmt.Sprintf("%+v", query))
	}

	affectedRows, err := res.RowsAffected()

	if affectedRows == 0 || err != nil {
		return "", domain.ErrNotFound(err)
	}

	return id, nil
}

func (m *model[T, E, C, U, F]) Destroy(id string, c context.Context, trx bun.IDB) (string, error) {
	if trx == nil {
		trx = m.db
	}

	query := trx.NewDelete().
		Where("id = ?", id)

	res, err := query.Exec(c)

	if err != nil {
		return "", db.QueryExecuteError(err, fmt.Sprintf("%+v", query))
	}

	affectedRows, err := res.RowsAffected()

	if affectedRows == 0 || err != nil {
		return "", domain.ErrNotFound(err)
	}

	return id, nil
}
