package rex_test

import (
	"net/http"
	"testing"

	"github.com/cswank/rex"
)

func BenchmarkRex(b *testing.B) {
	r := rex.New("bench")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/v1/{v1}", handler)

	request, _ := http.NewRequest("GET", "/v1/anything", nil)
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(nil, request)
	}
}
