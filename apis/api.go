package apis

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/web_crawler/models"
	"github.com/web_crawler/services"
)

type ApiController struct{}

func NewApiController(router *gin.RouterGroup) {
	controller := ApiController{}
	router.POST("/page/details", controller.getPageDetails)
}

func (a ApiController) getPageDetails(ginCtx *gin.Context) {
	var request models.PageDetailsRequest
	err := ginCtx.ShouldBindJSON(&request)
	if err != nil {
		log.Println(err)
		ginCtx.SecureJSON(http.StatusBadRequest, gin.H{"message": "Request is not in correct format"})
		return
	}

	url := "https://play.golang.org/"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	version, err := services.DetectHTML(bodyBytes)
	if err != nil {
		log.Println(err)
	}
	title, _ := services.GetHtmlTitle(bodyBytes)
	fmt.Println(version, title)

	internalLinks, externalLinks, err := services.GetLinks(bodyBytes, url)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(len(internalLinks), len(externalLinks))

	var allLinks []string
	allLinks = append(allLinks, internalLinks...)
	allLinks = append(allLinks, externalLinks...)

	inaccLinks := services.CheckInaccessibleLinks(allLinks)
	fmt.Println(len(inaccLinks))

	checkLoginForm := services.CheckLoginFormExists(bodyBytes)
	fmt.Println(checkLoginForm)
}
