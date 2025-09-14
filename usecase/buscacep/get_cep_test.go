package buscacep

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) Do(req *http.Request) (*http.Response, error) { return f(req) }

func newTestClient(fn roundTripFunc) HTTPDoer { return roundTripFunc(fn) }

func TestViaCepClient_Success(t *testing.T) {
	client := NewViaCepClient(newTestClient(func(r *http.Request) (*http.Response, error) {
		body := `{"localidade":"Sao Paulo"}`
		return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(body))}, nil
	}))
	loc, err := client.GetLocationByCep(context.Background(), "01001000")
	if err != nil {
		t.Fatalf("erro: %v", err)
	}
	if loc != "Sao Paulo" {
		t.Fatalf("esperava Sao Paulo got %s", loc)
	}
}

func TestViaCepClient_HTTPError(t *testing.T) {
	client := NewViaCepClient(newTestClient(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("net down")
	}))
	_, err := client.GetLocationByCep(context.Background(), "01001000")
	if err == nil {
		t.Fatalf("esperava erro")
	}
}

func TestViaCepClient_InvalidJSON(t *testing.T) {
	client := NewViaCepClient(newTestClient(func(r *http.Request) (*http.Response, error) {
		body := `{"localidade":` // truncado
		return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(body))}, nil
	}))
	_, err := client.GetLocationByCep(context.Background(), "01001000")
	if err == nil {
		t.Fatalf("esperava erro de JSON")
	}
}
