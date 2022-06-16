package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/zicops/sidecar-auth-proxy/server"
)

type handlerStub struct {
	caller  string
	handler http.Handler
}

func (h *handlerStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := r.Header.Get("caller")
	w.Header().Set("caller", c+h.caller)

	if h.handler != nil {
		r.Header.Set("caller", c+h.caller)
		h.handler.ServeHTTP(w, r)
	}
}

func check(caller string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler { return &handlerStub{caller, h} }
}

func Test_getHandler(t *testing.T) {
	server.AuthZHandler = check("a")

	type args struct {
		mode string
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantPanic bool
	}{
		{"No mode, panics", args{""}, "", true},
		{"Illegal mode, panics", args{"foo"}, "", true},
		{"authz", args{"authz"}, "aa", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil && !tt.wantPanic {
					t.Errorf("Unexpected panic with %v", err)
				}
			}()
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/echo", nil)
			got := getHandler(tt.args.mode, check("a")(nil))
			got.ServeHTTP(rr, req)
			caller := rr.Header().Get("caller")

			if caller != tt.want {
				t.Errorf("getHandler() = %v, want %v", caller, tt.want)
			}
		})
	}
}

//helper function to test the main
func getHandler(mode string, h http.Handler) http.Handler {
	handler := h
	modes := strings.Split(mode, "|")
	for i := len(modes) - 1; i >= 0; i-- {
		switch modes[i] {
		case "authz":
			handler = server.AuthZHandler(handler)
		default:
			log.Panicf("required -mode not given or illegal value")
		}
	}
	return handler
}
