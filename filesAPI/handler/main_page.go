package handler

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func (h *Handler) MainPage(c *gin.Context) {
	indexPage, err := ioutil.ReadFile("template/index.html")
	if err != nil {
		log.Printf("Cannot read index page file: %+v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", indexPage)
}
