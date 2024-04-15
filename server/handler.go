package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Banner struct {
	ID        int                    `json:"banner_id"`
	Tags      []int                  `json:"tag_ids"`
	Featue    int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func (s *Server) userBannerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		token := r.Header.Get("token")
		if token == "" {
			fmt.Fprintf(w, "Incorrect token")
			return
		}
		s.getBanners(w, r)
	default:
		fmt.Fprintf(w, "Incorrect request")
	}
}

func (s *Server) bannerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		admin := r.Header.Get("token")
		if admin == "" || admin != adminToken {
			fmt.Fprintf(w, "Incorrect token")
			return
		}
		s.getBanners(w, r)
	case "POST":
		admin := r.Header.Get("token")
		if admin == "" || admin != adminToken {
			fmt.Fprintf(w, "Incorrect token")
			return
		}
		var banner map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&banner)
		if err != nil {
			fmt.Fprintf(w, "Something went wrong with json")
			return
		}
		query := ` insert into "banners"("id", "tag_ids", "feature_id", "content", "is_active", "created_at", "updated_at")
		values($1, $2, $3, $4, $5, $6, $7)
		`
		_, error := s.db.Exec(query, len(banner["tag_ids"].([]int)) + 1,banner["tag_ids"].([]int), banner["feature_id"].(int), banner["content"].(map[string]interface{}), banner["is_active"].(bool), time.Now(), time.Now())
		if error != nil {
			fmt.Fprintf(w, "Something went wrong with execution")
		}
		w.Header().Set("Content")
		json.NewEncoder(w).Encode(banner)
	default:
		fmt.Fprintf(w, "Incorrect request")
	}
}

func (s *Server) idBannerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PATCH":
		admin := r.Header.Get("token")
		if admin == "" || admin != adminToken {
			fmt.Fprintf(w, "Incorrect token")
			return
		}
		var banner map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&banner)
		if err != nil {
			fmt.Fprintf(w, "Something went wrong with json")
			return
		}
		query := `update "banners" set "content"=$1, "updated_at"=$2
		where "id"=$3
		`
		_, error := s.db.Exec(query, banner["content"].(map[string]interface{}), time.Now(), banner["tag_id"].(int))
		if error != nil {
			fmt.Fprintf(w, "Something went wrong with execution")
		}
		w.Header().Set("Content")
		json.NewEncoder(w).Encode(banner)
	case "DELETE":
		admin := r.Header.Get("token")
		if admin == "" || admin != adminToken {
			fmt.Fprintf(w, "Incorrect token")
			return
		}
		var banner map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&banner)
		if err != nil {
			fmt.Fprintf(w, "Something went wrong with json")
			return
		}
		query := `delete "banners" where "id"=$3`
		_, error := s.db.Exec(query, banner["tag_id"].(int))
		if error != nil {
			fmt.Fprintf(w, "Something went wrong with execution")
		}
		w.Header().Set("Content")
		json.NewEncoder(w).Encode(banner)
	default:
		fmt.Fprintf(w, "Incorrect request")
	}
}

func (s *Server) getBanners(w http.ResponseWriter, r *http.Request) {
	tagID, _ := strconv.Atoi(r.URL.Query().Get("tag_id"))
	featureID, _ := strconv.Atoi(r.URL.Query().Get("feature_id"))
	last, _ := strconv.ParseBool(r.URL.Query().Get("use_last_version"))

	query := `SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at
		FROM banners
		WHERE tag_id = $1 AND feature_id = $2`
	rows, err := s.db.Query(query, tagID, featureID)
	if err != nil {
		fmt.Fprintf(w, "Something went wrong")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var banner Banner
		err = rows.Scan(&banner)
		if err != nil {
			fmt.Fprintf(w, "Incorrect token")
			return
		}
		json.NewEncoder(w).Encode(banner)
	}
}
