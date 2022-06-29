package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// StudentInfoRoutes -> struct
type StudentInfoRoutes struct {
	logger                    infrastructure.Logger
	router                    infrastructure.Router
	studentInfoController controllers.StudentInfoController
	middleware                middlewares.FirebaseAuthMiddleware
}

// NewStudentInfoRoutes -> creates new StudentInfo controller
func NewStudentInfoRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	studentInfoController controllers.StudentInfoController,
	middleware middlewares.FirebaseAuthMiddleware,
) StudentInfoRoutes {
	return StudentInfoRoutes{
		router:                    router,
		logger:                    logger,
		studentInfoController: studentInfoController,
		middleware:                middleware,
	}
}

// Setup studentInfo routes
func (c StudentInfoRoutes) Setup() {
	c.logger.Zap.Info(" Setting up StudentInfo routes")
	studentInfo := c.router.Gin.Group("/student_info")
	{
		studentInfo.POST("", c.studentInfoController.CreateStudentInfo)
		studentInfo.GET("", c.studentInfoController.GetAllStudentInfo)
		studentInfo.GET("/:id", c.studentInfoController.GetOneStudentInfo)
		studentInfo.PUT("/:id", c.studentInfoController.UpdateOneStudentInfo)
		studentInfo.DELETE("/:id", c.studentInfoController.DeleteOneStudentInfo)
	}
}
