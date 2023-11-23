package controllers

import (
	"net/http"
	"strings"

	"github.com/cocoza4/data_microservices/models"
	"github.com/cocoza4/data_microservices/services"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	ProductService services.ProductService
}

func NewProductController(service services.ProductService) ProductController {
	return ProductController{
		ProductService: service,
	}
}

func (ctr *ProductController) CreateProduct(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := ctr.ProductService.CreateProduct(&product)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (ctr *ProductController) GetProduct(ctx *gin.Context) {
	name := ctx.Param("name")
	product, err := ctr.ProductService.GetProduct(&name)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	ctx.JSON(http.StatusOK, product)
}

func (ctr *ProductController) GetLatest(ctx *gin.Context) {
	product, err := ctr.ProductService.GetLatest()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	ctx.JSON(http.StatusOK, product)
}

func (ctr *ProductController) GetAll(ctx *gin.Context) {
	products, err := ctr.ProductService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	ctx.JSON(http.StatusOK, products)
}

func (ctr *ProductController) GetIndexes(ctx *gin.Context) {
	indexes := ctr.ProductService.GetIndexes()
	var resp []map[string]string
	for _, index := range indexes {
		item := map[string]string{"name": *index}
		resp = append(resp, item)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (ctr *ProductController) CreateIndex(ctx *gin.Context) {
	field := ctx.Query("field")
	if len(strings.TrimSpace(field)) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "'field' can't be empty"})
		return
	}
	err := ctr.ProductService.CreateIndex(&field)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	}

}

func (ctr *ProductController) RegisterRoutes(rg *gin.RouterGroup) {
	route := *rg.Group("/products")
	route.GET("/", ctr.GetAll)
	route.POST("/", ctr.CreateProduct)
	route.GET("/latest", ctr.GetLatest)
	route.GET("/indexes", ctr.GetIndexes)
	route.POST("/indexes", ctr.CreateIndex)
	route.GET("/:name", ctr.GetProduct)
}
