package api

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	database "url-shortener/internal/db"
	"github.com/akamensky/base58"
	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ctxKey string

func handleCreate(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, ok := ctx.Value(ctxKey("db")).(*database.DB)
		if !ok {
			panic("wrong context key")
		}
		apiUrl, ok := ctx.Value(ctxKey("url")).(string)
		if !ok {
			panic("wrong context key")
		}
		req := CreateRequest{}
		raw := make([]byte, 1024)
		r.Body.Read(raw)
		json.Unmarshal(raw, &req)
		short := generateShortLink(req.Url)
		err := db.AddLink(req.Url, short)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		resp := CreateResponse{ShortUrl: apiUrl + "/" + short}
		res, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		w.Write(res)
	}
}

func handleRedirect(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, ok := ctx.Value(ctxKey("db")).(*database.DB)
		if !ok {
			panic("wrong context key")
		}
		id := chi.URLParam(r, "id")
		url, err := db.GetInitialLink(id)
		if err != nil {
			return
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func generateShortLink(url string) string {
	hash := sha256.Sum256([]byte(url))
	short := base58.Encode(hash[:])
	return short[:10]
}

func Init(r chi.Router, db *database.DB, apiUrl string) {
	ctx := context.WithValue(context.Background(), ctxKey("db"), db)
	ctx = context.WithValue(ctx, ctxKey("url"), apiUrl)
	r.Post("/create", handleCreate(ctx))
	r.Get("/{id:.{10}}", handleRedirect(ctx))
}