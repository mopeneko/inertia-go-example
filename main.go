package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/petaki/inertia-go"
)

var inertiaManager *inertia.Inertia

func main() {
	manifest, err := LoadManifest()
	if err != nil {
		log.Fatal(err)
	}

	asset := func(path string) string {
		return "/" + manifest[path].File
	}

	inertiaManager = inertia.New("http://localhost:4321", "./app.gohtml", "")
	inertiaManager.ShareFunc("asset", asset)

	mux := http.NewServeMux()
	mux.Handle("/", inertiaManager.Middleware(http.HandlerFunc(homeHandler)))
	mux.Handle("/name/{name}", inertiaManager.Middleware(http.HandlerFunc(nameHandler)))
	mux.Handle("/assets/", http.FileServer(http.Dir("./front/dist")))

	http.ListenAndServe(":4321", mux)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := inertiaManager.Render(w, r, "home/Index", map[string]any{
		"name": "World",
	})
	if err != nil {
		log.Printf("Failed to render page: %+v", err)
	}
}

func nameHandler(w http.ResponseWriter, r *http.Request) {
	err := inertiaManager.Render(w, r, "name/Echo", map[string]any{
		"name": r.PathValue("name"),
	})
	if err != nil {
		log.Printf("Failed to render page: %+v", err)
	}
}

type Manifest struct {
	File    string `json:"file"`
	Name    string `json:"name"`
	Src     string `json:"src"`
	IsEntry bool   `json:"isEntry"`
}

func LoadManifest() (map[string]Manifest, error) {
	manifest, err := os.ReadFile("front/dist/.vite/manifest.json")
	if err != nil {
		return nil, err
	}

	var manifestData map[string]Manifest
	if err := json.Unmarshal(manifest, &manifestData); err != nil {
		return nil, err
	}

	return manifestData, nil
}
