package gogadgets_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"

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

		It("connects to a gadget", func() {
			ch := cli.Connect()
			expected := map[string]string{"address": fmt.Sprintf("%s/messages", cliAddr), "token": "n/a"}
			Expect(req).To(Equal(expected))

			buf := bytes.Buffer{}
			enc := json.NewEncoder(&buf)
			msg := gogadgets.Message{Type: gogadgets.UPDATE, Value: gogadgets.Value{Value: 33.3}}
			enc.Encode(msg)

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				msg2 := <-ch
				Expect(msg2.Value).To(Equal(msg.Value))
				wg.Done()
			}()

			resp, err := http.Post(fmt.Sprintf("http://%s", expected["address"]), "application/json", &buf)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			wg.Wait()
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
