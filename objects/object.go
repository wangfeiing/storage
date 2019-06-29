package objects

import (
	"net/http"
	"os"
	"storage/config"
	"strings"
	"log"
	"io"
)

func Handler(w http.ResponseWriter , r * http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w , r)
	} else if m == http.MethodGet {
		get(w , r)
	}
}

func put(w http.ResponseWriter, r *http.Request) {

	f , e := os.Create(config.STORAGE_ROOT + "/objects/" + strings.Split(r.URL.EscapedPath() , "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()
	io.Copy(f , r.Body)
}

func get(w http.ResponseWriter , r* http.Request )  {
	f , e := os.Open(config.STORAGE_ROOT + "/objects/" + strings.Split(r.URL.EscapedPath() , "/")[2])

	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer f.Close()
	io.Copy(w , f)
}
