package gogadgets_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {

	Context("success", func() {

		var (
			ts      *httptest.Server
			addr    string
			req     map[string]string
			cli     *gogadgets.Client
			cliAddr string
		)

		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var m map[string]string
				dec := json.NewDecoder(r.Body)
				Expect(dec.Decode(&m)).To(BeNil())
				req = m
			}))
			addr = ts.URL

			cliAddr = fmt.Sprintf("127.0.0.1:%d", getPort())
			cli = gogadgets.NewClient(cliAddr, addr)
		})

		AfterEach(func() {
			ts.Close()
		})

		FIt("connects to a gadget", func() {
			cli.Connect()
			expected := map[string]string{"address": fmt.Sprintf("%s/messages", cliAddr), "token": "n/a"}
			Expect(req).To(Equal(expected))
		})
	})
})

func getPort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}
