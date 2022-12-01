package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service services.ProductService
}

func NewProductController(
	service services.ProductService,
) ProductController {

	return ProductController{
		service: service,
	}

}

func (cc ProductController) AddProducts(ctx *gin.Context) {
	responses.SuccessJSON(ctx, http.StatusOK, "Product Added successfully")
}

func (cc ProductController) GetAllProducts(ctx *gin.Context) {
	responses.SuccessJSON(ctx, http.StatusOK, "Get All Products")
}
