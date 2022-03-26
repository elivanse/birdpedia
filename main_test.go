package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	hf := http.HandlerFunc(handler)
	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler retorno status code erroneo: obtuvo %v esperando %v", status, http.StatusOK)
	}

	expected := "Hello mister Anderson"
	actual := recorder.Body.String()

	if actual != expected {
		t.Errorf("handler retorno status code erroneo: obtuvo %v esperando %v ", actual, expected)
	}

}

func TestRouter(t *testing.T) {

	r := newRouter()
	mockServer := httptest.NewServer(r)
	resp, err := http.Get(mockServer.URL + "/hello")

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status debe estar bien, ubtuvo %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	respString := string(b)
	expected := "Hello mister Anderson"

	if respString != expected {
		t.Errorf("Respuesta debe ser %s, obtuvo  %s", expected, respString)
	}

}

func TestRouterForNonExistentRoute(t *testing.T) {

	r := newRouter()
	mockServer := httptest.NewServer(r)
	resp, err := http.Post(mockServer.URL+"/hello", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status debe ser 405, obtuvo %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	respString := string(b)
	expected := ""

	if respString != expected {
		t.Errorf("Response deberia haber sido %s, obtuvo %s", expected, respString)
	}

}

func TestStaticFileServer(t *testing.T) {

	r := newRouter()
	mockServer := httptest.NewServer(r)
	resp, err := http.Get(mockServer.URL + "/assets/")

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status debe ser 200, obtuvo %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("tipo de contenido erroneo, esperando %s,obtuvo %s", expectedContentType, contentType)
	}

}
