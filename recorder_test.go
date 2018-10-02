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
		r     *gogadgets.Recorder
		ts    *httptest.Server
		posts []string
		urls  []string
	)

	BeforeEach(func() {
		posts = []string{}
		urls = []string{}
		ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf := &bytes.Buffer{}
			io.Copy(buf, r.Body)
			posts = append(posts, strings.TrimSpace(buf.String()))
			urls = append(urls, r.URL.Path)
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
		r.On(nil)
	})

	AfterEach(func() {
		ts.Close()
	})

	Describe("when all's good", func() {

		It("saves output values", func() {
			msg := &gogadgets.Message{
				Type:     gogadgets.UPDATE,
				Location: "lab",
				Name:     "thermometer",
				Value: gogadgets.Value{
					Value: 411.0,
				},
			}
			r.Update(msg)
			Expect(posts).To(ConsistOf(`{"value":411}`))
			Expect(urls).To(ConsistOf("/api/fakeid/locations/lab/devices/thermometer/datapoints"))
		})

		It("saves input values", func() {
			msg := &gogadgets.Message{
				Type:     gogadgets.UPDATE,
				Location: "lab",
				Name:     "thermostat",
				Value: gogadgets.Value{
					Output: map[string]bool{
						"heat": true,
						"cool": false,
					},
				},
				Info: gogadgets.Info{
					Direction: "output",
				},
			}
			r.Update(msg)
			Expect(posts).To(ConsistOf(`{"value":1}`, `{"value":0}`))
			Expect(urls).To(ConsistOf("/api/fakeid/locations/lab/devices/thermostat heat/datapoints", "/api/fakeid/locations/lab/devices/thermostat cool/datapoints"))
		})
	})
})
