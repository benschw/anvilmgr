package anvilrepo

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var _ = log.Printf

type AnvilRepoService struct {
	Bind  string
	Path  string
	WebFS http.FileSystem
}

func NewAnvilRepoService(bind string, path string, webFs http.FileSystem) (*AnvilRepoService, error) {
	s := &AnvilRepoService{
		Bind:  bind,
		Path:  path,
		WebFS: webFs,
	}

	return s, nil
}

func (s *AnvilRepoService) Run() error {

	// route handlers
	resource := &AnvilArtifactResource{
		Repo: &RepoClient{Path: s.Path},
	}

	// Configure Routes
	r := mux.NewRouter()

	r.HandleFunc("/api/repo/{user}/{module}", resource.GetModuleArtifacts).Methods("Get")
	r.HandleFunc("/api/repo/{user}/{module}", resource.AddArtifact).Methods("POST")
	r.HandleFunc("/api/repo", resource.FindAll).Methods("GET")

	// r.HandleFunc("/", s.HandleWebRequest).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(s.WebFS))
	// http.Handle("/", http.FileServer(rice.MustFindBox("http-files").HTTPBox()))

	// http.FileServer(
	// 	&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "dist"}))

	http.Handle("/", r)

	// Start HTTP Server
	return http.ListenAndServe(s.Bind, nil)
}

// func (s *AnvilRepoService) HandleWebRequest(res http.ResponseWriter, req *http.Request) error {
// 	data, err := Asset("pub/style/foo.css")
// 	if err != nil {
// 		// Asset was not found.
// 	}

// }
