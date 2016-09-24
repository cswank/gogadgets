package gogadgets_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {

	Context("success", func() {

		var (
			ts   *httptest.Server
			addr string
			req  map[string]string
			cli  *gogadgets.Client
		)

		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var m map[string]string
				dec := json.NewDecoder(r.Body)
				Expect(dec.Decode(&m)).To(BeNil())
				req = m
			}))
			addr = ts.URL

			cli = gogadgets.NewClient("127.0.0.1:", addr)
		})

		AfterEach(func() {
			ts.Close()
		})

		FIt("connects to a gadget", func() {
			cli.Connect()
			expected := map[string]string{"address": "127.0.0.1:/messages", "token": "n/a"}
			Expect(req).To(Equal(expected))
		})
	})
})
