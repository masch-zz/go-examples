package main

//Source: https://talks.golang.org/2013/bestpractices.slide#11

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

//######################################## Basic handler ########################################

func BasicHandler(w http.ResponseWriter, r *http.Request) {
	err := doThis()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("handling %q: %v", r.RequestURI, err)
		return
	}

	err = doThat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("handling %q: %v", r.RequestURI, err)
		return
	}
}

func TestBasicHandler(t *testing.T) {
	doHandler(BasicHandler, t)
}

//################################## Function adapter handler ###################################

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
		}
	}
}

func betterHandler(w http.ResponseWriter, r *http.Request) error {
	if err := doThis(); err != nil {
		return fmt.Errorf("doing this: %v ", err)
	}

	if err := doThat(); err != nil {
		return fmt.Errorf("doing that: %v", err)
	}

	return nil
}

func TestFunctionAdapterHandler(t *testing.T) {
	doHandler(errorHandler(betterHandler), t)
}

//###############################################################################################

func doHandler(handlerFunc http.HandlerFunc, t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func doThis() error {
	return nil
}

func doThat() error {
	return nil
}
