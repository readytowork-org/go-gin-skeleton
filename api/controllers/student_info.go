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

// StudentInfoController -> struct
type StudentInfoController struct {
	logger             infrastructure.Logger
	StudentInfoService services.StudentInfoService
}

// NewStudentInfoController -> constructor
func NewStudentInfoController(
	logger infrastructure.Logger,
	StudentInfoService services.StudentInfoService,
) StudentInfoController {
	return StudentInfoController{
		logger:             logger,
		StudentInfoService: StudentInfoService,
	}
}

// CreateStudentInfo -> Create StudentInfo
func (cc StudentInfoController) CreateStudentInfo(c *gin.Context) {
	studentInfo := models.StudentInfo{}

	if err := c.ShouldBindJSON(&studentInfo); err != nil {
		cc.logger.Zap.Error("Error [CreateStudentInfo] (ShouldBindJson) : ", err)
		responses.ErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := cc.StudentInfoService.CreateStudentInfo(studentInfo); err != nil {
		cc.logger.Zap.Error("Error [CreateStudentInfo] [db CreateStudentInfo]: ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "StudentInfo Created Sucessfully")
}

// GetAllStudentInfo -> Get All StudentInfo
func (cc StudentInfoController) GetAllStudentInfo(c *gin.Context) {

	pagination := utils.BuildPagination(c)
	pagination.Sort = "created_at desc"
	studentInfos, count, err := cc.StudentInfoService.GetAllStudentInfo(pagination)

	if err != nil {
		cc.logger.Zap.Error("Error finding StudentInfo records", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Find StudentInfo")
		return
	}
	responses.JSONCount(c, http.StatusOK, studentInfos, count)

}

// GetOneStudentInfo -> Get One StudentInfo
func (cc StudentInfoController) GetOneStudentInfo(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	studentInfo, err := cc.StudentInfoService.GetOneStudentInfo(ID)

	if err != nil {
		cc.logger.Zap.Error("Error [GetOneStudentInfo] [db GetOneStudentInfo]: ", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Find StudentInfo")
		return
	}
	responses.JSON(c, http.StatusOK, studentInfo)

}

// UpdateOneStudentInfo -> Update One StudentInfo By Id
func (cc StudentInfoController) UpdateOneStudentInfo(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	studentInfo := models.StudentInfo{}

	if err := c.ShouldBindJSON(&studentInfo); err != nil {
		cc.logger.Zap.Error("Error [UpdateStudentInfo] (ShouldBindJson) : ", err)
		responses.ErrorJSON(c, http.StatusBadRequest, "failed to update studentInfo")
		return
	}
	studentInfo.ID = ID

	if err := cc.StudentInfoService.UpdateOneStudentInfo(studentInfo); err != nil {
		cc.logger.Zap.Error("Error [UpdateStudentInfo] [db UpdateStudentInfo]: ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, "failed to update studentInfo")
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "StudentInfo Updated Sucessfully")
}

// DeleteOneStudentInfo -> Delete One StudentInfo By Id
func (cc StudentInfoController) DeleteOneStudentInfo(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := cc.StudentInfoService.DeleteOneStudentInfo(ID)

	if err != nil {
		cc.logger.Zap.Error("Error [DeleteOneStudentInfo] [db DeleteOneStudentInfo]: ", err.Error())
		responses.ErrorJSON(c, http.StatusBadRequest, "Failed to Delete StudentInfo")
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "StudentInfo Deleted Sucessfully")
}
