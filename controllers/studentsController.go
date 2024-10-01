package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Thiago-Maia/gin-api-rest-alura/database"
	"github.com/Thiago-Maia/gin-api-rest-alura/models"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/swag/example/celler/httputil"
)

type StudentsController struct {
}

func NewStudentsController() *StudentsController {
	return &StudentsController{}
}

// Get godoc
//
//	@Summary		Show students
//	@Description	Show all students
//	@Tags			Students Student
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.Student
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/students [get]
func (ctr *StudentsController) Get(c *gin.Context) {
	var students []models.Student

	database.Db.Order("name ASC").Find(&students)

	c.JSON(http.StatusOK, gin.H{
		"sucess":  true,
		"message": "Students",
		"data":    students})
}

// Create godoc
//
//	@Summary		Create student
//	@Description	Create student
//	@Tags			Students Student
//	@Accept			json
//	@Produce		json
//	@Param			student 	body		models.Student	true	"Student Model"
//	@Success		200	{object}	models.Student
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/students [post]
func (ctr *StudentsController) Create(c *gin.Context) {
	var student models.Student

	if err := c.ShouldBind(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sucess":  false,
			"message": err.Error()})
		return
	}

	if err := student.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sucess":  false,
			"message": err.Error()})
		return
	}

	database.Db.Create(&student)
	c.JSON(http.StatusOK, gin.H{
		"sucess":  true,
		"message": "Student created succefuly",
		"data":    student})
}

func (ctr *StudentsController) FindOneById(c *gin.Context) {
	var student models.Student

	id := c.Params.ByName("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"sucess":  false,
			"message": "Id required",
			"data":    nil})
		return
	}

	database.Db.First(&student, id)

	if student.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"sucess":  false,
			"message": "Student not found",
			"data":    nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sucess":  true,
		"message": fmt.Sprintf("Student with id: %s Found succefuly", id),
		"data":    student})

}

func (ctr *StudentsController) FindManyByName(c *gin.Context) {
	var students []models.Student

	name := c.Param("name")

	name = strings.ToLower(name)

	database.Db.Where("LOWER(name) LIKE ?", "%"+name+"%").Find(&students)

	if len(students) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"sucess":  false,
			"message": "Student not found",
			"data":    nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sucess":  true,
		"message": "Student found succefuly",
		"data":    students})

}

func (ctr *StudentsController) Delete(c *gin.Context) {
	id := c.Params.ByName("id")

	database.Db.Delete(&models.Student{}, id)

	c.JSON(http.StatusOK, gin.H{
		"sucess":  true,
		"message": "Student deleted succefuly",
		"data":    nil})
}

func (ctr *StudentsController) Update(c *gin.Context) {
	var newStudent models.Student

	if err := c.ShouldBind(&newStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sucess":  false,
			"message": err.Error()})
		return
	}

	if err := newStudent.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"sucess":  false,
			"message": err.Error()})
		return
	}

	var student models.Student
	database.Db.First(&student, newStudent.ID)

	database.Db.Model(&student).UpdateColumns(newStudent)

	c.JSON(http.StatusOK, gin.H{
		"sucess":  true,
		"message": "Student updated succefuly",
		"data":    student})
}
