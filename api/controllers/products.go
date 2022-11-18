package controllers

import (
	"boilerplate-api/api/responses"
	"boilerplate-api/api/services"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	logger         infrastructure.Logger
	productService services.ProductService
}

func NewProductController(
	logger infrastructure.Logger,
	productService services.ProductService,
	env infrastructure.Env,
) ProductController {
	return ProductController{
		logger:         logger,
		productService: productService,
	}
}

func (cc ProductController) CreateProduct(c *gin.Context) {
	product := models.Products{}
	fmt.Println(cc, "context")
	if err := c.ShouldBindJSON(&product); err != nil {
		cc.logger.Zap.Error("Error [CreateProduct] (ShouldBindJson) : ", err.Error())
		err := errors.BadRequest.Wrap(err, "Failed to bind product body data")
		responses.HandleError(c, err)
		return
	}
	if err := cc.productService.CreateProduct(product); err != nil {
		cc.logger.Zap.Error("Error [CreateProduct] (db CreateProduct) : ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to bind create product")
		responses.HandleError(c, err)
		return
	}
	responses.SuccessJSON(c, http.StatusOK, product)
}

func (cc ProductController) GetAllProducts(c *gin.Context) {
	pagination := utils.BuildPagination(c)
	products, count, err := cc.productService.GetAllProducts(pagination)
	if err != nil {
		cc.logger.Zap.Error("Error finding products", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users data")
		responses.HandleError(c, err)
		return
	}
	responses.JSONCount(c, http.StatusOK, products, count)
}

func (cc ProductController) GetProduct(c *gin.Context) {
	product := models.Products{}

	product.ID = 22
	products, err := cc.productService.GetProduct(product)

	if err != nil {
		cc.logger.Zap.Error("Error finding product")
	}

	responses.SuccessJSON(c, http.StatusOK, products)
}
