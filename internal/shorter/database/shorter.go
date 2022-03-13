package database

import (
	"errors"
	"github.com/EpicStep/vk-tarantool/internal/shorter/model"
	"github.com/EpicStep/vk-tarantool/pkg/database"
	"github.com/tarantool/go-tarantool"
)

const SpaceShort = "short"
const SpaceAnalytics = "transitions"

var ErrNotFound = errors.New("element not found")

type ShorterDB struct {
	db *database.DB
}

func New(db *database.DB) *ShorterDB {
	return &ShorterDB{db: db}
}

func (db *ShorterDB) CreateShort(s *model.Short) error {
	_, err := db.db.DB.Insert(SpaceShort, []interface{}{s.Shorted, s.Original, s.CreatedBy})
	if err != nil {
		return err
	}

	return nil
}

func (db *ShorterDB) GetShortByHash(hash string) ([]model.Short, error) {
	var s []model.Short

	err := db.db.DB.SelectTyped(SpaceShort, "primary", 0, 1, tarantool.IterEq, tarantool.StringKey{S: hash}, &s)
	if err != nil {
		return nil, err
	}

	if len(s) <= 0 {
		return nil, ErrNotFound
	}

	return s, nil
}

func (db *ShorterDB) InsertAnalytics(t *model.Transition) error {
	_ = db.db.DB.InsertAsync(SpaceAnalytics, []interface{}{t.ID, t.Shorted, t.IP, t.UA})

	return nil
}

func (db *ShorterDB) GetTransitionByShort(short string) ([]model.Transition, error) {
	var s []model.Transition

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
