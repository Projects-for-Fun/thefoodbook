package webservice

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (w *Webservice) GetDummyData(rw http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		panic("error")
		//rw.Header().Set("random-header", "dummy error")
		//http.Error(rw, "generic error", http.StatusInternalServerError)
		//return
	}

	rw.Header().Set("random-header", "dummy data")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("dummy search " + id))
}

func (w *Webservice) GetDummyDataForId(rw http.ResponseWriter, r *http.Request) {
	// Context cancelled
	if r.Context().Err() != nil {
		return
	}

	//ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	//defer cancel()
	//id := r.URL.Query().Get(":id")
	id := chi.URLParam(r, "id")

	rw.Header().Set("random-header", "dummy data")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("dummy " + id))
}
