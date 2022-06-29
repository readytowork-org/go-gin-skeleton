package repository

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/utils"
)

// StudentInfoRepository database structure
type StudentInfoRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// NewStudentInfoRepository creates a new StudentInfo repository
func NewStudentInfoRepository(db infrastructure.Database, logger infrastructure.Logger) StudentInfoRepository {
	return StudentInfoRepository{
		db:     db,
		logger: logger,
	}
}

// Create StudentInfo
func (c StudentInfoRepository) Create(StudentInfo models.StudentInfo) (models.StudentInfo, error) {
	return StudentInfo, c.db.DB.Create(&StudentInfo).Error
}

// GetAllStudentInfo -> Get All studentInfos
func (c StudentInfoRepository) GetAllStudentInfo(pagination utils.Pagination) ([]models.StudentInfo, int64, error) {
	var studentInfos []models.StudentInfo
	var totalRows int64 = 0
	queryBuider := c.db.DB.Model(&models.StudentInfo{}).Offset(pagination.Offset).Order(pagination.Sort)

	if !pagination.All {
		queryBuider = queryBuider.Limit(pagination.PageSize)
	}

	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuider.Where(c.db.DB.Where("`student_infos`.`name` LIKE ?", searchQuery))
	}

	err := queryBuider.
		Find(&studentInfos).
		Offset(-1).
		Limit(-1).
		Count(&totalRows).Error
	return studentInfos, totalRows, err
}

// GetOneStudentInfo -> Get One StudentInfo By Id
func (c StudentInfoRepository) GetOneStudentInfo(ID int64) (models.StudentInfo, error) {
	StudentInfo := models.StudentInfo{}
	return StudentInfo, c.db.DB.
		Where("id = ?", ID).First(&StudentInfo).Error
}

// UpdateOneStudentInfo -> Update One StudentInfo By Id
func (c StudentInfoRepository) UpdateOneStudentInfo(StudentInfo models.StudentInfo) error {
	return c.db.DB.Model(&models.StudentInfo{}).
		Where("id = ?", StudentInfo.ID).
		Updates(map[string]interface{}{
			"name": StudentInfo.Name,
		}).Error
}

// DeleteOneStudentInfo -> Delete One StudentInfo By Id
func (c StudentInfoRepository) DeleteOneStudentInfo(ID int64) error {
	return c.db.DB.
		Where("id = ?", ID).
		Delete(&models.StudentInfo{}).
		Error
}
