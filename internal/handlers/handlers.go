// Package handlers ...
package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/SerjRamone/metrius/internal/metrics"
	"github.com/SerjRamone/metrius/internal/storage"
	"github.com/go-chi/chi/v5"
)

// Router returns chi.Router
func Router(mS storage.MemStorage) chi.Router {
	r := chi.NewRouter()

	r.Get("/", List(mS))
	r.Get("/value/{type}/{name}", Value(mS))
	r.Post("/update/{type}/{name}/{value}", Update(mS))

	return r
}

// ComingSoon func
func ComingSoon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("it's works"))
	}
}

// List ...
func List(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const tmpl = `<html>
			<head><title>Metrius</title></head>
			<body>
				<table border="1" cellspacing="0">
					<thead><tr><th>Metrics</th><th>Value</th></tr></thead>
					<tbody>
						{{range $name, $value := .}}
							<tr><td>{{$name}}</td><td>{{$value}}</td></tr>
						{{end}}
					</tbody>
				</table>
			</body>
		</html>`
		metrics := map[string]string{}
		for name, value := range s.Gauges() {
			metrics[name] = fmt.Sprintf("%v", value)
		}
		for name, value := range s.Counters() {
			metrics[name] = fmt.Sprintf("%v", value)
		}

		t := template.New("metrics tpl")
		t, err := t.Parse(tmpl)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		buf := new(bytes.Buffer)
		err = t.Execute(buf, metrics)
		if err != nil {
			http.Error(w, "500", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(buf.Bytes())
	}
}

// Value handler handls GET "/value/counter/foo" requests
func Value(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			mType  = chi.URLParam(r, "type")
			mName  = chi.URLParam(r, "name")
			mValue string
		)

		switch mType {
		case "gauge":
			v, ok := s.Gauge(mName)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			mValue = fmt.Sprint(v)
		case "counter":
			v, ok := s.Counter(mName)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			mValue = fmt.Sprint(v)
		default:
			http.Error(w, "unknown type", http.StatusBadRequest)
			return
		}

		_, _ = w.Write([]byte(mValue))
	}
}

// Update is a /update/ handler
func Update(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			mValue = chi.URLParam(r, "value")
			mType  = chi.URLParam(r, "type")
			mName  = chi.URLParam(r, "name")
		)

		// @todo
		// if r.Header.Get("Content-Type") != "text/plain" {
		// 	http.Error(w, "Bad content-type", http.StatusBadRequest)
		// 	return
		// }

		log.Println("📨  ", r.Method, r.URL.Path)

		if mName == "" {
			http.Error(w, "Metrics name not set", http.StatusNotFound)
			return
		}

		if mType != "counter" && mType != "gauge" {
			http.Error(w, "Metrics type not set or unknown", http.StatusBadRequest)
			return
		}

		if mValue == "" {
			http.Error(w, "Metrics value not set or invalid", http.StatusBadRequest)
			return
		}

		switch mType {
		case "counter":
			c, err := strconv.ParseInt(mValue, 10, 64)
			if err != nil {
				http.Error(w, "Invalid metrics value", http.StatusBadRequest)
				return
			}

			if err := s.SetCounter(mName, metrics.Counter(c)); err != nil {
				log.Fatal("can't set counter", err)
				return
			}

		case "gauge":
			g, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				http.Error(w, "Invalid metrics value", http.StatusBadRequest)
				return
			}

			if err := s.SetGauge(mName, metrics.Gauge(g)); err != nil {
				log.Fatal("can't set gauge", err)
				return
			}
		}

		_, _ = w.Write([]byte("OK"))
	}
}
