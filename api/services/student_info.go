package services

import (
	"boilerplate-api/api/repository"
	"boilerplate-api/models"
	"boilerplate-api/utils"

)

// StudentInfoService -> struct
type StudentInfoService struct {
	repository repository.StudentInfoRepository
}

// NewStudentInfoService  -> creates a new StudentInfoservice
func NewStudentInfoService(repository repository.StudentInfoRepository) StudentInfoService {
	return StudentInfoService{
		repository: repository,
	}
}

// CreateStudentInfo -> call to create the StudentInfo
func (c StudentInfoService) CreateStudentInfo(studentInfo models.StudentInfo) (models.StudentInfo, error) {
	return c.repository.Create(studentInfo)
}

// GetAllStudentInfo -> call to create the StudentInfo
func (c StudentInfoService) GetAllStudentInfo(pagination utils.Pagination) ([]models.StudentInfo, int64, error) {
	return c.repository.GetAllStudentInfo(pagination)
}

// GetOneStudentInfo -> Get One StudentInfo By Id
func (c StudentInfoService) GetOneStudentInfo(ID int64) (models.StudentInfo, error) {
	return c.repository.GetOneStudentInfo(ID)
}

// UpdateOneStudentInfo -> Update One StudentInfo By Id
func (c StudentInfoService) UpdateOneStudentInfo(studentInfo models.StudentInfo) error {
	return c.repository.UpdateOneStudentInfo(studentInfo)
}

// DeleteOneStudentInfo -> Delete One StudentInfo By Id
func (c StudentInfoService) DeleteOneStudentInfo(ID int64) error {
	return c.repository.DeleteOneStudentInfo(ID)

}