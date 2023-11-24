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

func (ctr *ProductController) createProduct(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := ctr.ProductService.CreateProduct(&product)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func (ctr *ProductController) getProduct(ctx *gin.Context) {
	name := ctx.Param("name")
	product, err := ctr.ProductService.GetProduct(&name)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, product)
	}
}

func (ctr *ProductController) getLatest(ctx *gin.Context) {
	product, err := ctr.ProductService.GetLatest()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, product)
	}

}

func (ctr *ProductController) getAll(ctx *gin.Context) {
	products, err := ctr.ProductService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	} else {
		if len(products) == 0 {
			ctx.JSON(http.StatusOK, gin.H{"message": "documents not found"})
		} else {
			ctx.JSON(http.StatusOK, products)
		}
	}
}

func (ctr *ProductController) getIndexes(ctx *gin.Context) {
	indexes := ctr.ProductService.GetIndexes()
	var resp []map[string]string
	for _, index := range indexes {
		item := map[string]string{"name": *index}
		resp = append(resp, item)
	}
	ctx.JSON(http.StatusOK, resp)
}

func (ctr *ProductController) createIndex(ctx *gin.Context) {
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
	route.GET("/", ctr.getAll)
	route.POST("/", ctr.createProduct)
	route.GET("/latest", ctr.getLatest)
	route.GET("/indexes", ctr.getIndexes)
	route.POST("/indexes", ctr.createIndex)
	route.GET("/:name", ctr.getProduct)
}
