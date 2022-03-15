package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_CachedHandler(t *testing.T) {
	cache, shutdownFunc := NewCache(30 * time.Second)
	defer shutdownFunc()

	h := http.HandlerFunc(doExpensiveWork)
	h = CachedHandler(cache, h)

	req, err := http.NewRequest(http.MethodGet, "http://ticketure/test", nil)
	if err != nil {
		t.Error("fail to create new request")
	}
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	sameReq, err := http.NewRequest(http.MethodGet, "http://ticketure/test", nil)
	if err != nil {
		t.Error("fail to create new request")
	}
	expectedSameRes := httptest.NewRecorder()
	h.ServeHTTP(expectedSameRes, sameReq)
	if string(res.Body.Bytes()) != string(expectedSameRes.Body.Bytes()) {
		t.Errorf("response different but expect same: %q %q", string(res.Body.Bytes()), string(expectedSameRes.Body.Bytes()))
	} else if res.Code != expectedSameRes.Code {
		t.Errorf("response code different but expect same: %d %d", res.Code, expectedSameRes.Code)
	}

	reqWithDiffURI, err := http.NewRequest(http.MethodGet, "http://ticketure/test/1", nil)
	if err != nil {
		t.Error("fail to create new request")
	}
	resWithDiffURI := httptest.NewRecorder()
	h.ServeHTTP(resWithDiffURI, reqWithDiffURI)
	if string(res.Body.Bytes()) == string(resWithDiffURI.Body.Bytes()) {
		t.Errorf("URI path differnt: response same but expect different: %q %q", string(res.Body.Bytes()), string(resWithDiffURI.Body.Bytes()))
	}

	reqWithDiffMethod, err := http.NewRequest(http.MethodPost, "http://ticketure/test", nil)
	if err != nil {
		t.Error("fail to create new request")
	}
	resWithDiffMethod := httptest.NewRecorder()
	h.ServeHTTP(resWithDiffMethod, reqWithDiffMethod)
	if string(res.Body.Bytes()) == string(resWithDiffMethod.Body.Bytes()) {
		t.Errorf("http method differnt: response same but expect different: %q %q", string(res.Body.Bytes()), string(resWithDiffMethod.Body.Bytes()))
	}
}
