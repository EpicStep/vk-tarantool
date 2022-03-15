package database

import (
	"context"
	"errors"
	"github.com/EpicStep/vk-tarantool/internal/shorter/model"
	"github.com/EpicStep/vk-tarantool/pkg/database"
	"github.com/tarantool/go-tarantool"
	"go.opentelemetry.io/otel/trace"
)

const (
	SpaceShort     = "short"
	SpaceAnalytics = "transitions"
)

var ErrNotFound = errors.New("element not found")
var ErrAlreadyExists = errors.New("element already exists")

type ShorterDB struct {
	db *database.DB
}

func New(db *database.DB) *ShorterDB {
	return &ShorterDB{db: db}
}

func (db *ShorterDB) CreateShort(ctx context.Context, s *model.Short) error {
	span := trace.SpanFromContext(ctx)
	span.SetName("db.create_short")
	defer span.End()

	resp, err := db.db.DB.Insert(SpaceShort, []interface{}{s.Shorted, s.Original, s.CreatedBy})
	if err != nil {
		if resp.Code == tarantool.ErrTupleFound {
			return ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (db *ShorterDB) GetShortByHash(ctx context.Context, hash string) ([]model.Short, error) {
	var s []model.Short

	span := trace.SpanFromContext(ctx)
	span.SetName("db.get_short_by_hash")
	defer span.End()

	err := db.db.DB.SelectTyped(SpaceShort, "primary", 0, 1, tarantool.IterEq, tarantool.StringKey{S: hash}, &s)
	if err != nil {
		return nil, err
	}

	if len(s) <= 0 {
		return nil, ErrNotFound
	}

	return s, nil
}

func (db *ShorterDB) InsertAnalytics(ctx context.Context, t *model.Transition) error {
	span := trace.SpanFromContext(ctx)
	span.SetName("db.insert_analytics")
	defer span.End()

	_ = db.db.DB.InsertAsync(SpaceAnalytics, []interface{}{t.ID, t.Shorted, t.IP, t.UA})

	return nil
}

func (db *ShorterDB) GetTransitionByShort(ctx context.Context, short string) ([]model.Transition, error) {
	var s []model.Transition

	span := trace.SpanFromContext(ctx)
	span.SetName("db.get_transition_by_short")
	defer span.End()

	var i []int64
	err := db.db.DB.EvalTyped("return box.space."+SpaceAnalytics+":count()", []interface{}{}, &i)
	if err != nil {
		return nil, err
	}

	if len(i) <= 0 {
		return nil, ErrNotFound
	}

	err = db.db.DB.SelectTyped(SpaceAnalytics, "shorted_idx", 0, uint32(i[0]), tarantool.IterEq, tarantool.StringKey{S: short}, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
