package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func CacheIndex(w http.ResponseWriter, r *http.Request) {
	dbcaches, err := getCaches()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var apicaches APICaches
	for db := range dbcaches {
		apicaches = append(apicaches, DBToAPI(dbcaches[db]))
	}

	if err := json.NewEncoder(w).Encode(apicaches); err != nil {
		panic(err)
	}
}

func CacheShow(w http.ResponseWriter, r *http.Request) {
	cacheId, err := strconv.ParseUint(mux.Vars(r)["cacheId"], 10, 64)
	if err != nil {
		panic(err)
	}

	cache, err := getCache(cacheId)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(DBToAPI(cache)); err != nil {
		panic(err)
	}
}

func CacheCreate(w http.ResponseWriter, r *http.Request) {
	var postcache PostCache

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &postcache); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	filename := uuid.New() + ".mp3"
	dbcache := PostToDB(postcache, filename)

	err = postCache(dbcache)
	if err != nil {
		panic(err)
	}

	err = writeFile(postcache, filename)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	fmt.Fprintln(w, "OK")
}