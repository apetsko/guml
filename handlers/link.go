package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"oss.terrastruct.com/util-go/go2"
)

func Link(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UML string `json:"link"`
	}

	link := r.URL.Query().Get("link")
	fmt.Println(link)
	client := http.Client{}
	resp, err := client.Get("https://google.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	return
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	ruler, _ := textmeasure.NewRuler()
	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return d2dagrelayout.DefaultLayout, nil
	}
	renderOpts := &d2svg.RenderOpts{
		Pad:     go2.Pointer(int64(5)),
		ThemeID: &d2themescatalog.GrapeSoda.ID,
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler:          ruler,
	}
	ctx := log.WithDefault(context.Background())
	diagram, _, err := d2lib.Compile(ctx, req.UML, compileOpts, renderOpts)
	if err != nil {
		http.Error(w, "D2 compile error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	svg, err := d2svg.Render(diagram, renderOpts)
	if err != nil {
		http.Error(w, "SVG render error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	if _, err = w.Write(svg); err != nil {
		http.Error(w, "SVG write error: "+err.Error(), http.StatusInternalServerError)
	}
}
