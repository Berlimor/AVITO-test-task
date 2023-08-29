package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserWithSegment struct {
	ID			string		`json:"id"`
	Segments	string		`json:"segments"`
}

func main() {
	env := new(Env)
	var err error
	env.DB, err = DBConnect()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/users", ListUsers)
	r.GET("/users/:id", GetUsersSegments)
	r.POST("/users/", AddSegmentsToUser)

	r.POST("/segments/:segmentName", CreateSegment)
	r.DELETE("/segments/:segmentName", DeleteSegment)

	port := os.Getenv("PORT") 
	if port == "" {
		port = "8080"
	}
	r.Run(port)
}

func ListUsers(c *gin.Context) {
	query := "SELECT * FROM users"

}

func GetUsersSegments(c *gin.Context) {
	var user UserWithSegment
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Could not bind JSON")
		return
	}
}

func CreateSegment(c *gin.Context) {
	var user UserWithSegment
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Could not bind JSON")
		return
	}
	newSegment := user.Segments
}

func DeleteSegment(c *gin.Context) {
	var user UserWithSegment
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Could not bind JSON")
		return
	}
	segmentToDelete := user.Segments
}

func AddSegmentsToUser(c *gin.Context) {
	var user UserWithSegment
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Could not bind JSON")
		return
	}
	segmentsToAdd := strings.Split(user.Segments, ",")
}