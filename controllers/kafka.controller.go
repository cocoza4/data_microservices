package controllers

import (
	"net/http"

	"github.com/cocoza4/data_microservices/models"
	"github.com/cocoza4/data_microservices/services"
	"github.com/gin-gonic/gin"
)

type KafkaController struct {
	KafkaService services.KafkaService
}

func NewKafkaController(service services.KafkaService) KafkaController {
	return KafkaController{
		KafkaService: service,
	}
}

func (ctr *KafkaController) getTopics(ctx *gin.Context) {
	topics, err := ctr.KafkaService.GetTopics()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	} else {
		if len(topics) == 0 {
			ctx.JSON(http.StatusOK, gin.H{"message": "topics not found"})
		} else {

			var resp []map[string]string
			for _, topic := range topics {
				item := map[string]string{"name": topic}
				resp = append(resp, item)
			}

			ctx.JSON(http.StatusOK, resp)
		}
	}
}

func (ctr *KafkaController) createTopic(ctx *gin.Context) {
	name := ctx.Param("name")
	err := ctr.KafkaService.CreateTopic(&name)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func (ctr *KafkaController) publishMessage(ctx *gin.Context) {
	topic := ctx.Param("name")
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid message body"})
		return
	}
	err := ctr.KafkaService.PublishMessage(&topic, &product)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func (ctr *KafkaController) getLatest(ctx *gin.Context) {
	topic := ctx.Param("name")
	latest, err := ctr.KafkaService.GetLatest(&topic)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, latest)
	}
}

func (ctr *KafkaController) RegisterRoutes(rg *gin.RouterGroup) {
	route := *rg.Group("/kafka")
	route.GET("/topics", ctr.getTopics)
	route.POST("/topics/:name", ctr.createTopic)
	route.POST("/topics/:name/publish", ctr.publishMessage)
	route.GET("/topics/:name/latest", ctr.getLatest)
}
