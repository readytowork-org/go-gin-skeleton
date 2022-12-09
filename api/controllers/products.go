package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	logger  infrastructure.Logger
	service services.ProductService
}

func NewProductController(
	service services.ProductService,
	logger infrastructure.Logger,
) ProductController {

	return ProductController{
		service: service,
		logger:  logger,
	}

}

func (cc ProductController) AddProducts(ctx *gin.Context) {
	var newProduct models.ProductCreateInput
	cc.logger.Zap.Info(newProduct)
	if err := ctx.ShouldBindJSON(&newProduct); err != nil {
		cc.logger.Zap.Error("Error Create user Should bind Json error")
		responses.HandleError(ctx, err)
		return
	}
	if err := cc.service.AddProduct(newProduct); err != nil {
		cc.logger.Zap.Error("Addproduct Error")
		responses.HandleError(ctx, err)
		return
	}
	responses.SuccessJSON(ctx, http.StatusOK, "Product Added successfully")
}

func (cc ProductController) GetAllProducts(ctx *gin.Context) {
	pagination := utils.BuildPagination(ctx)
	allProducts, count, err := cc.service.GetAllProduct(pagination)
	if err != nil {
		responses.HandleError(ctx, err)
		return
	}
	cc.logger.Zap.Info(allProducts)
	responses.JSONCount(ctx, http.StatusOK, allProducts, count)
}

func (cc ProductController) FilterUserProducts(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	pagination := utils.BuildPagination(ctx)
	products, err := cc.service.FilterUserProducts(int64(id), pagination)
	if err != nil {
		cc.logger.Zap.Error("Error", err)
		return
	}
	responses.SuccessJSON(ctx, http.StatusOK, products)
}
