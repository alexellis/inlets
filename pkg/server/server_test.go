package server

import (
	"net/http"
	"testing"
)

func Test_tokenValid_Valid(t *testing.T) {
	token := "abcdefg"
	s := &Server{Token: token}
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Authorization", "Bearer "+token)
	if !s.tokenValid(r) {
		t.Error("expected isTokenValid to be true")
	}
}

func Test_tokenValid_Invalid(t *testing.T) {
	token := "abcdefg"
	s := &Server{Token: token}
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Authorization", "Bearer tuvwxyz")
	if s.tokenValid(r) {
		t.Error("expected isTokenValid to be false")
	}
}
