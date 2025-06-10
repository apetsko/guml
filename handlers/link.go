package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/apetsko/guml/uml"
)

func Link(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		link := r.URL.Query().Get("link")
		uri, err := url.Parse(link)
		if err != nil {
			logger.Error("failed to parse link")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if uri == nil {
			logger.Error("Invalid URL: " + link)
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		fmt.Println(link)
		client := http.Client{}
		resp, err := client.Get(uri.String())
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			logger.Error("failed do request", "uri", uri.String(), "error", err)
			return
		}
		defer resp.Body.Close()

		byteUML, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("failed read response body", "uri", uri.String(), "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		u, err := uml.UMLtoSVG(string(byteUML))

		w.Header().Set("Content-Type", "image/svg+xml")
		if _, err = w.Write(u); err != nil {
			logger.Error("SVG write error: " + err.Error())
			http.Error(w, "SVG write error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
