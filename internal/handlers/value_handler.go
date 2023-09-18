package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Value handler handls GET "/value/counter/foo" requests
func (bHandler baseHandler) Value() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			mType  = chi.URLParam(r, "type")
			mName  = chi.URLParam(r, "name")
			mValue string
		)

		switch mType {
		case "gauge":
			v, ok := bHandler.storage.Gauge(mName)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			mValue = fmt.Sprint(v)
		case "counter":
			v, ok := bHandler.storage.Counter(mName)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			mValue = fmt.Sprint(v)
		default:
			http.Error(w, "unknown type", http.StatusBadRequest)
			return
		}

		_, err := w.Write([]byte(mValue))
		if err != nil {
			log.Println("can't write response:", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
