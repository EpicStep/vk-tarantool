package shorter

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	shortDB "github.com/EpicStep/vk-tarantool/internal/shorter/database"
	"github.com/EpicStep/vk-tarantool/internal/shorter/model"
	"github.com/EpicStep/vk-tarantool/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type Service struct {
	db *shortDB.ShorterDB
}

func New(db *database.DB) *Service {
	return &Service{db: shortDB.New(db)}
}

func (s *Service) Short(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")

	if url == "" {
		return
	}

	hasher := md5.New()

	hasher.Write([]byte(url))
	hexStr := hex.EncodeToString(hasher.Sum(nil))

	ub64 := base64.StdEncoding.EncodeToString([]byte(hexStr))

	ipArr := strings.Split(r.RemoteAddr, ":")

	err := s.db.CreateShort(&model.Short{
		Shorted:  ub64[0:12],
		Original: url,
		CreatedBy: ipArr[0],
	})

	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte("http://37.139.34.190/" + ub64[0:12]))
}

func (s *Service) Transition(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	short, err := s.db.GetShortByHash(hash)
	if err != nil {
		return
	}

	if len(short) <= 0 {
		return
	}

	ipArr := strings.Split(r.RemoteAddr, ":")

	_ = s.db.InsertAnalytics(&model.Transition{
		ID: uuid.New().String(),
		Shorted: hash,
		IP:      ipArr[0],
		UA:      r.UserAgent(),
	})

	http.Redirect(w, r, short[0].Original, http.StatusMovedPermanently)
}
