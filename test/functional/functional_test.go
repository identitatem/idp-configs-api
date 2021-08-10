package functional_test

import (
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var resp *http.Response
var err error
var baseURL string

var _ = Describe("HelloWorldService", func() {
	Context("The /api/hello-world-service/v0/ping endpoint responds successfully", func() {
		// Get base URL for server from environment variable; default is set to localhost for local testing
		if (os.Getenv("APP_BASE_URL") != "") {
			baseURL = os.Getenv("APP_BASE_URL")
		} else {
			baseURL = "http://localhost:3000"
		}

		BeforeEach(func() {
			resp, err = http.Get(baseURL + "/api/hello-world-service/v0/ping")
		})

		// Verify status code
		It ("/api/hello-world-service/v0/ping responds with 200", func () {
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(200))
			Expect(err).NotTo(HaveOccurred())
		})		

		// Verify response
		It ("/api/hello-world-service/v0/ping serves up Hello World", func () {
			expectedResponse := "Hello world"
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			Expect(string(body)).To(Equal(expectedResponse))
			Expect(err).NotTo(HaveOccurred())
		})				
	})	
})
