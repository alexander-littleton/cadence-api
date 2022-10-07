package controllers_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "internal/controllers"
	"internal/models"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

var _ = Describe("Main", func() {
	var r *gin.Engine
	var w *httptest.ResponseRecorder
	var req *http.Request
	var data []byte
	BeforeEach(func() {
		r = SetUpRouter()
		w = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		r.ServeHTTP(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, _ = ioutil.ReadAll(res.Body)
	})

	Context("CreateUser", func() {
		BeforeEach(func() {
			requestUser := models.User{Email: "test@test.com"}
			r.POST("/", CreateUser())
			jsonValue, _ := json.Marshal(requestUser)
			req, _ = http.NewRequest("POST", "/", bytes.NewBuffer(jsonValue))
		})
		Context("the request is valid", func() {
			It("returns a 201", func() {
				Expect(w.Code).To(Equal(201))
				Expect(string(data)).To(ContainSubstring("InsertedID"))
			})
		})
		Context("the request is invalid", func() {
			BeforeEach(func() {
				req, _ = http.NewRequest("POST", "/", nil)
			})
			It("returns a 400 error", func() {
				Expect(w.Code).To(Equal(400))
				Expect(string(data)).To(ContainSubstring("invalid request"))
			})
		})
	})
})
