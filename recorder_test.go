package gogadgets_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Recorder", func() {
	var (
		out   chan gogadgets.Message
		in    chan gogadgets.Value
		r     *gogadgets.Recorder
		ts    *httptest.Server
		posts []string
	)

	BeforeEach(func() {
		posts = []string{}
		ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf := &bytes.Buffer{}
			io.Copy(buf, r.Body)
			posts = append(posts, buf.String())
			r.Body.Close()
		}))
		cfg := fmt.Sprintf(`{
  "args": {
    "host": "%s/api/fakeid/locations/%%s/devices/%%s/datapoints",
    "token": "xyz"
  }
}`, ts.URL)
		p := &gogadgets.Pin{}
		err := json.Unmarshal([]byte(cfg), p)
		Expect(err).To(BeNil())
		x, err := gogadgets.NewRecorder(p)
		Expect(err).To(BeNil())
		r = x.(*gogadgets.Recorder)
		out = make(chan gogadgets.Message)
		in = make(chan gogadgets.Value)
		r.On(nil)
	})

	AfterEach(func() {
		ts.Close()
	})

	Describe("when all's good", func() {
		It("does it's thing", func() {
			msg := &gogadgets.Message{
				Type:     gogadgets.UPDATE,
				Location: "lab",
				Name:     "thermometer",
				Value: gogadgets.Value{
					Value: 411.0,
				},
			}
			r.Update(msg)
			Expect(len(posts)).To(Equal(1))
			Expect(strings.TrimSpace(posts[0])).To(Equal(`{"value":411}`))
		})
	})
})
