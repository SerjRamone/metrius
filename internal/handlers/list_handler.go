package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// List handles GET requests to the root address (/) of the project, returning an HTML table with the current metric values.
// Possible response status codes:
//   - 500 in case of a service error.
//   - 200 for a successful request.
func (bHandler baseHandler) List() http.HandlerFunc {
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
		for name, value := range bHandler.storage.Gauges(r.Context()) {
			metrics[name] = fmt.Sprintf("%v", value)
		}
		for name, value := range bHandler.storage.Counters(r.Context()) {
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
		_, err = w.Write(buf.Bytes())
		if err != nil {
			log.Println("can't write response:", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
