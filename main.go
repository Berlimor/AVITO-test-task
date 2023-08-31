package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type UserWithSegment struct {
	ID			int				`json:"id"`
	Segments	[]string		`json:"segments"`
}

type Segment struct {
	Name		string		`json:"name"`
}

type AddingSegmentsToUser struct {
	SegmentsToAdd		[]string		`json:"add-segments"`
	SegmentsToDelete	[]string		`json:"delete-segments"`
	UserID				int				`json:"user-id"`
}

func main() {
	env := new(Env)
	var err error
	env.DB, err = DBConnect()
	if err != nil {
		panic(err)
	}
	defer env.DB.Close()

	// Create tables
	err = env.CreateUsersTable()
	if err != nil {
		panic(err)
	}
	err = env.CreateSegmentsTable()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/users", env.ListUsers)
	r.GET("/users/:id", env.GetUserByID)
	r.POST("/users", env.AddSegmentsToUser)

	r.POST("/segments", env.CreateSegment)
	r.DELETE("/segments", env.DeleteSegment)


	r.Run(":8000")
}

func (env Env) ListUsers(c *gin.Context) {
	query := "SELECT * FROM users"
	rows, err := env.DB.Query(query)

	switch err{
	case sql.ErrNoRows:
		defer rows.Close()
		c.IndentedJSON(http.StatusInternalServerError, "No rows found in users table")
		return

	case nil:
		result := make([]UserWithSegment, 0)
		for rows.Next() {
			var id int 
			var segments []string
			err = rows.Scan(&id, pq.Array(&segments))
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, "Error scaning rows")
				return
			}
			result = append(result, UserWithSegment{id, segments})
		}
		c.IndentedJSON(http.StatusOK, result)
		return

	default:
		defer rows.Close()
		c.IndentedJSON(http.StatusInternalServerError, "There are no users in the database")
		return
	}
}

func (env Env) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Invalid id parameter")
		return
	}

	query := "SELECT * FROM users WHERE id=$1"
	row := env.DB.QueryRow(query, id)
	var segments []string

	err = row.Scan(&id, pq.Array(&segments))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Error scaning a row")
	}
	c.IndentedJSON(http.StatusOK, UserWithSegment{id, segments})
}

func (env Env) CreateSegment(c *gin.Context) {
	var segment Segment
	if err := c.BindJSON(&segment); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Could not bind JSON")
		return
	}
	newSegment := segment.Name

	query := "INSERT INTO segments(segment) VALUES($1)"
	result, err := env.DB.Exec(query, newSegment)
	if err != nil {
		error := fmt.Sprintf("Error %v", err)
		c.IndentedJSON(http.StatusInternalServerError, error)
		return
	}

	_, err = result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "An error occured while checking the returned result")
		return
	}

	c.IndentedJSON(http.StatusOK, "Segment created successfully")
}

func (env Env) DeleteSegment(c *gin.Context) {
	var segment Segment
	if err := c.BindJSON(&segment); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Could not bind JSON")
		return
	}
	segmentToDelete := segment.Name

	query := "DELETE FROM segments WHERE segment=$1"
	result, err := env.DB.Exec(query, segmentToDelete)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Could not delete specified segment instance")
		return
	}

	_, err = result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "An error occured while checking the returned result")
		return
	}

	c.IndentedJSON(http.StatusOK, "Segment deleted successfully")
}

func (env Env) AddSegmentsToUser(c *gin.Context) {
	var request AddingSegmentsToUser
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Could not bind JSON")
		return
	}
	segmentsToAdd := request.SegmentsToAdd
	segmentsToDelete := request.SegmentsToDelete
	userID := request.UserID

	query := "SELECT * FROM users WHERE id=$1"
	row := env.DB.QueryRow(query, userID)
	var oldID int
	var oldSegments []string
	if err := row.Scan(oldID, pq.Array(&oldSegments)); err == sql.ErrNoRows {
		query = "INSERT INTO users VALUES($1, $2)"
		_, err = env.DB.Exec(query, userID, segmentsToAdd)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Could not insert user and it's segments")
			return
		}
		c.IndentedJSON(http.StatusOK, "Segments successfully added")
		return
	}
	
	uniqueOldSegments := GetMap(oldSegments)

	for _, segment := range segmentsToDelete {
		delete(uniqueOldSegments, segment)
	}

	for _, segment := range segmentsToAdd {
		if _, unique := uniqueOldSegments[segment]; !unique {
			uniqueOldSegments[segment] = true
		}
	}

	var updatedSegments []string
	for key := range uniqueOldSegments {
		// Checking if the segment occures in segments table
		query = "SELECT exists(SELECT * FROM segments WHERE segment=$1)"
		var exists bool
		err := env.DB.QueryRow(query, key).Scan(&exists)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Error validating the segment")
		}
		if exists {
			updatedSegments = append(updatedSegments, key)
		}
	}

	query = "UPDATE users SET segments=$1 WHERE id=$2;"
	result, err := env.DB.Exec(query, pq.Array(updatedSegments), userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Could not update")
		return
	}

	if _, err = result.RowsAffected(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "An error occured while checking the returned result")
	}

	c.IndentedJSON(http.StatusOK, "User segments updated successfully")
}