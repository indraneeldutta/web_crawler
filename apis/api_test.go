package apis

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/web_crawler/models"
)

func TestNewApiController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	g := gin.Default()
	NewApiController(g.Group("/v1"))
}

type ApiControllerTestSuite struct {
	suite.Suite
	ctrl       *gomock.Controller
	controller *ApiController
	gin        *gin.Engine
}

func (suite *ApiControllerTestSuite) BeforeTest(suiteName, testName string) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(suite.T())
	suite.ctrl = ctrl
	suite.gin = gin.Default()
	suite.controller = NewApiController(suite.gin.Group("/v1"))
}

func (suite *ApiControllerTestSuite) AfterTest(suiteName, testName string) {
	suite.ctrl.Finish()
}

func TestApiControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ApiControllerTestSuite))
}

func (suite *ApiControllerTestSuite) TestGetPageDetails() {
	request := models.PageDetailsRequest{
		Url: "https://google.com",
	}

	data, _ := json.Marshal(request)

	r, _ := http.NewRequest(http.MethodPost, "/v1/page/details", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	suite.gin.ServeHTTP(w, r)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *ApiControllerTestSuite) TestGetPageDetails_Error() {
	request := models.PageDetailsRequest{
		Url: "bad url",
	}

	data, _ := json.Marshal(request)

	r, _ := http.NewRequest(http.MethodPost, "/v1/page/details", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	suite.gin.ServeHTTP(w, r)

	suite.Equal(http.StatusInternalServerError, w.Code)
}

func (suite *ApiControllerTestSuite) TestGetPageDetails_BadRequest() {
	data := `{"url": 12}`

	r, _ := http.NewRequest(http.MethodPost, "/v1/page/details", bytes.NewBuffer([]byte(data)))
	w := httptest.NewRecorder()
	suite.gin.ServeHTTP(w, r)

	suite.Equal(http.StatusBadRequest, w.Code)
}
