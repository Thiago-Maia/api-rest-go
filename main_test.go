package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Thiago-Maia/gin-api-rest-alura/controllers"
	"github.com/Thiago-Maia/gin-api-rest-alura/database"
	"github.com/Thiago-Maia/gin-api-rest-alura/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type response struct {
	Data    models.Student `json:"data"`
	Message string         `json:"message"`
	Success bool           `json:"sucess"`
}

type listResponse struct {
	Data    []models.Student `json:"data"`
	Message string           `json:"message"`
	Success bool             `json:"sucess"`
}

func SetupRoute() *gin.Engine {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	database.Connect()
	gin.SetMode("release")
	rotas := gin.Default()
	return rotas
}

var studentsController = controllers.NewStudentsController()
var studentMock = models.Student{
	Name: "Jhon Doe",
	Cpf:  "00000061653",
	Rg:   "000000004"}

func TestCreateStudend(t *testing.T) {
	r := SetupRoute()
	method := "POST"
	url := "/students"
	r.POST(url, studentsController.Create)

	payload := getPayload(studentMock)
	req, res := getReqAndRes(method, url, payload)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusBadRequest, res.Code)

	studentMock.Cpf = "00000061654"
	payload = getPayload(studentMock)
	req, res = getReqAndRes(method, url, payload)
	r.ServeHTTP(res, req)
	if assert.Equal(t, http.StatusOK, res.Code) {

		jsonResponse := getResponse(*res.Body)

		studentMock.ID = jsonResponse.Data.ID
	}
}

func TestGetStudents(t *testing.T) {
	r := SetupRoute()
	method, url := "GET", "/students"
	r.GET(url, studentsController.Get)
	req, res := getReqAndRes(method, url, nil)

	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetStudentById(t *testing.T) {
	res := getStudentById()

	assert.Equal(t, http.StatusOK, res.Code)
	jsonResponse := getResponse(*res.Body)
	assert.Equal(t, studentMock.Name, jsonResponse.Data.Name)
}

func TestGetStudentByName(t *testing.T) {
	r := SetupRoute()
	method, url := "GET", fmt.Sprintf("/students/findName/%s", studentMock.Name[0:2])
	r.GET("/students/findName/:name", studentsController.FindManyByName)
	req, res := getReqAndRes(method, url, nil)
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	jsonResponse := getListResponse(*res.Body)

	assert.Equal(t, studentMock.Name, jsonResponse.Data[0].Name)
}

func TestUpdateStudent(t *testing.T) {
	r := SetupRoute()
	method, url := "PUT", "/students"
	r.PUT("/students", studentsController.Update)

	studentMock.Cpf = "14695938741"
	payload := getPayload(studentMock)

	req, res := getReqAndRes(method, url, payload)
	r.ServeHTTP(res, req)

	assert.NotEqual(t, http.StatusOK, res.Code)
}

func TestDeleteStudent(t *testing.T) {
	r := SetupRoute()
	method, url := "DELETE", fmt.Sprintf("/students/%d", studentMock.ID)
	r.DELETE("/students/:id", studentsController.Delete)

	req, res := getReqAndRes(method, url, nil)
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	*res = getStudentById()
	assert.Equal(t, http.StatusNotFound, res.Code)
}

// Privates Methods
func getStudentById() httptest.ResponseRecorder {
	r := SetupRoute()
	method, url := "GET", fmt.Sprintf("/students/%d", studentMock.ID)
	r.GET("/students/:id", studentsController.FindOneById)
	req, res := getReqAndRes(method, url, nil)
	r.ServeHTTP(res, req)

	return *res
}

func getReqAndRes(method, url string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(method, url, body)
	res := httptest.NewRecorder()

	if err != nil {
		log.Fatal("Erro ao criar rota!")
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, res
}

func getPayload(jsonBody any) *bytes.Buffer {
	body, err := json.Marshal(jsonBody)

	if err != nil {
		log.Fatal("Erro ao criar rota!")
	}

	return bytes.NewBuffer(body)
}

func getResponse(res bytes.Buffer) response {
	var jsonResponse response
	err := json.Unmarshal(res.Bytes(), &jsonResponse)
	if err != nil {
		log.Fatal(err.Error())
	}
	return jsonResponse
}

func getListResponse(res bytes.Buffer) listResponse {
	var jsonResponse listResponse
	err := json.Unmarshal(res.Bytes(), &jsonResponse)
	if err != nil {
		log.Fatal(err.Error())
	}
	return jsonResponse
}
