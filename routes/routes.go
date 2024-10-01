package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Thiago-Maia/gin-api-rest-alura/controllers"
	"github.com/Thiago-Maia/gin-api-rest-alura/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func HandleRequests() {
	port := os.Getenv("PORT")
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}

	gin.SetMode(ginMode)
	docs.SwaggerInfo.BasePath = "/api/v1"

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "healthy") })
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := r.Group("/api/v1")
	getStudentsRoutes(v1)

	fmt.Printf("Servidor escutando porta: %s", port)
	r.Run(":" + port)
}

func getStudentsRoutes(g *gin.RouterGroup) {
	var studentsController = controllers.NewStudentsController()

	students := g.Group("/students")
	{
		students.GET("/", studentsController.Get)
		students.GET("/:id", studentsController.FindOneById)
		students.GET("/findName/:name", studentsController.FindManyByName)
		students.POST("/", studentsController.Create)
		students.DELETE("/:id", studentsController.Delete)
		students.PUT("/", studentsController.Update)
	}
}
