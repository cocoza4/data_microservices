package controllers

import (
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
	ctx.JSON(200, "nil")
}

func (ctr *ProductController) GetProduct(ctx *gin.Context) {
	ctx.JSON(200, "nil")
}

func (ctr *ProductController) GetAll(ctx *gin.Context) {
	ctx.JSON(200, "nil")
}

func (ctr *ProductController) RegisterRoutes(rg *gin.RouterGroup) {
	route := *rg.Group("/products")
	route.GET("/", ctr.GetAll)
	route.POST("/create", ctr.CreateProduct)
	route.GET("/:name", ctr.GetProduct)
}
