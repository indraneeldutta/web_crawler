package apis

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/web_crawler/models"
	"github.com/web_crawler/services"
)

type ApiController struct{}

func NewApiController(router *gin.RouterGroup) *ApiController {
	controller := ApiController{}
	router.POST("/page/details", controller.getPageDetails)

	return &controller
}

func (a ApiController) getPageDetails(ginCtx *gin.Context) {
	var (
		request                                  models.PageDetailsRequest
		version, title                           string
		wg                                       sync.WaitGroup
		internalLinks, externalLinks, inaccLinks []string
		checkLoginForm                           bool
	)
	err := ginCtx.ShouldBindJSON(&request)
	if err != nil {
		log.Println(err)
		ginCtx.SecureJSON(http.StatusBadRequest, gin.H{"message": "Request is not in correct format"})
		return
	}

	resp, err := http.Get(request.Url)
	if err != nil {
		log.Println(err)
		ginCtx.SecureJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	wg.Add(1)
	go func(w *sync.WaitGroup) {
		version, err = services.DetectHTML(bodyBytes)
		if err != nil {
			log.Println(err)
		}
		w.Done()
	}(&wg)

	wg.Add(1)
	go func(w *sync.WaitGroup) {
		title, _ = services.GetHtmlTitle(bodyBytes)
		w.Done()
	}(&wg)

	wg.Add(1)
	go func(w *sync.WaitGroup) {
		internalLinks, externalLinks, err = services.GetLinks(bodyBytes, request.Url)
		if err != nil {
			log.Println(err)
		}
		w.Done()
	}(&wg)

	var allLinks []string
	allLinks = append(allLinks, internalLinks...)
	allLinks = append(allLinks, externalLinks...)

	wg.Add(1)
	go func(w *sync.WaitGroup) {
		inaccLinks = services.CheckInaccessibleLinks(allLinks)
		w.Done()
	}(&wg)

	wg.Add(1)
	go func(w *sync.WaitGroup) {
		checkLoginForm = services.CheckLoginFormExists(bodyBytes)
		w.Done()
	}(&wg)

	wg.Wait()

	var response models.PageDetailsResponse

	response.HtmlVersion = version
	response.PageTitle = title
	response.LinksCount.InternalLinks = len(internalLinks)
	response.LinksCount.ExternalLinks = len(externalLinks)
	response.InaccessibleLinks = len(inaccLinks)
	response.LoginFormExists = checkLoginForm

	ginCtx.SecureJSON(http.StatusOK, gin.H{
		"data": response,
	})
}
