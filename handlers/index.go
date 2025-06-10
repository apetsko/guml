package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/apetsko/guml/uml"
)

func Index(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<form method="POST" enctype="multipart/form-data">
					<input type="file" name="d2file">
					<button type="submit">Upload</button>
				</form>`)
			return
		}
		file, _, err := r.FormFile("d2file")
		if err != nil {
			logger.Error("Error reading file: %v", err)
			http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()
		d2data, err := io.ReadAll(file)
		if err != nil {
			logger.Error("Failed to read file: %v", err)
			http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		svg, err := uml.UMLtoSVG(string(d2data))

		w.Header().Set("Content-Type", "image/svg+xml")
		if _, err = w.Write(svg); err != nil {
			logger.Error("Error writing response: %v", err)
			http.Error(w, "SVG write error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
