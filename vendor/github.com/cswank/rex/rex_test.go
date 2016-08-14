package rex_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"

	"github.com/cswank/rex"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gadgets", func() {
	var (
		r      *rex.Router
		w      *httptest.ResponseRecorder
		req    *http.Request
		method string
		data   string
	)

	Context("fileserver", func() {
		var (
			tmp string
		)
		BeforeEach(func() {
			var err error
			tmp, err = ioutil.TempDir("", "")
			Expect(err).To(BeNil())
			Expect(ioutil.WriteFile(path.Join(tmp, "hello.html"), []byte("hello, world!"), 777)).To(BeNil())
			method = ""
			data = ""

			pals := http.HandlerFunc(func(ww http.ResponseWriter, rr *http.Request) {
				method = rr.Method
				ww.Write([]byte("pals"))
			})

			r = rex.New("my other router")
			r.Get("/pals", pals)

			r.ServeFiles(http.FileServer(http.Dir(tmp)))
			http.Handle("/", r)

			w = httptest.NewRecorder()
		})

		AfterEach(func() {
			os.RemoveAll(tmp)
		})

		It("gets a file", func() {
			var err error
			req, err = http.NewRequest("GET", "/hello.html", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("hello, world!"))
		})

	})

	Context("routes", func() {
		BeforeEach(func() {
			method = ""
			data = ""
			root := http.HandlerFunc(func(ww http.ResponseWriter, rr *http.Request) {
				method = rr.Method
				ww.Write([]byte("root"))
			})

			pals := http.HandlerFunc(func(ww http.ResponseWriter, rr *http.Request) {
				method = rr.Method
				ww.Write([]byte("pals"))
			})

			post := http.HandlerFunc(func(ww http.ResponseWriter, rr *http.Request) {
				method = rr.Method
				d, _ := ioutil.ReadAll(rr.Body)
				data = string(d)
			})

			delete := http.HandlerFunc(func(ww http.ResponseWriter, rr *http.Request) {
				method = rr.Method
			})

			pal := http.HandlerFunc(func(ww http.ResponseWriter, rr *http.Request) {
				method = rr.Method
				ww.Write([]byte("pal"))
			})

			pets := http.HandlerFunc(func(ww http.ResponseWriter, rr *http.Request) {
				method = rr.Method
				ww.Write([]byte("pets"))
			})

			pet := http.HandlerFunc(func(ww http.ResponseWriter, rr *http.Request) {
				method = rr.Method
				ww.Write([]byte("pet"))
			})

			colors := http.HandlerFunc(func(ww http.ResponseWriter, rr *http.Request) {
				method = rr.Method
				ww.Write([]byte("colors"))
			})

			color := http.HandlerFunc(func(ww http.ResponseWriter, rr *http.Request) {
				method = rr.Method
				ww.Write([]byte("color"))
			})

			r = rex.New("my router")
			r.Get("/", root)
			r.Get("/pals", pals)
			r.Post("/pals", post)
			r.Get("/pals/{id}", pal)
			r.Get("/pals/{id}/pets", pets)
			r.Get("/pals/{id}/pets/{pet}", pet)
			r.Get("/pals/{id}/colors", colors)
			r.Post("/pals/{id}/colors", post)
			r.Get("/pals/{id}/colors/{color}", color)
			r.Delete("/pals/{id}/colors/{color}", delete)

			w = httptest.NewRecorder()
		})

		AfterEach(func() {
		})

		It("gets root", func() {
			var err error
			req, err = http.NewRequest("GET", "/", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("root"))
		})

		It("gets the first collection", func() {
			var err error
			req, err = http.NewRequest("GET", "/pals", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("pals"))
		})

		It("gets the first resource", func() {
			var err error
			req, err = http.NewRequest("GET", "/pals/1", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("pal"))
		})

		It("gets the pets collection", func() {
			var err error
			req, err = http.NewRequest("GET", "/pals/1/pets", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("pets"))
		})

		It("gets a pet resource", func() {
			var err error
			req, err = http.NewRequest("GET", "/pals/1/pets/5", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("pet"))
		})

		It("gets a pet resource with params", func() {
			var err error
			req, err = http.NewRequest("GET", "/pals/1/pets/5?alive=true", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("pet"))
		})

		It("gets the colors collection", func() {
			var err error
			req, err = http.NewRequest("GET", "/pals/1/colors", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("colors"))
			vars := rex.Vars(req, "my router")
			Expect(vars["id"]).To(Equal("1"))
		})

		It("gets a colors resource", func() {
			var err error
			req, err = http.NewRequest("GET", "/pals/3/colors/red", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("color"))
			vars := rex.Vars(req, "my router")
			Expect(vars["id"]).To(Equal("3"))
			Expect(vars["color"]).To(Equal("red"))
		})

		It("POSTS to pals", func() {
			var err error
			buf := bytes.NewBuffer([]byte("stu"))
			req, err = http.NewRequest("POST", "/pals", buf)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(method).To(Equal("POST"))
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal(""))
			Expect(data).To(Equal("stu"))
		})

		It("POSTS to colors", func() {
			var err error
			buf := bytes.NewBuffer([]byte("green"))
			req, err = http.NewRequest("POST", "/pals/55/colors", buf)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(method).To(Equal("POST"))
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal(""))
			Expect(data).To(Equal("green"))
		})

		It("DELETES a color", func() {
			var err error
			req, err = http.NewRequest("DELETE", "/pals/55/colors/red", nil)
			Expect(err).To(BeNil())
			r.ServeHTTP(w, req)
			Expect(method).To(Equal("DELETE"))
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal(""))
			vars := rex.Vars(req, "my router")
			Expect(vars["id"]).To(Equal("55"))
			Expect(vars["color"]).To(Equal("red"))
		})
	})
})
