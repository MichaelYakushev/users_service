package main

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

// user represents data about a user.
type user struct {
	GithubId    string `json:"github_id"`
	TgId        string `json:"telegramid"`
	Roles       string `json:"roles"`
	Fio         string `json:"fio"`
	GroupNumber string `json:"group_number"`
}

type errorMessage struct {
	Message string `json:"message"`
}

// users slice to seed user data.
var users = []user{
	{GithubId: "11242", TgId: "8837hSh", Roles: "Студент", Fio: "Иван Кононский", GroupNumber: "ИВТ-232"},
	{GithubId: "1242", TgId: "49957hSh", Roles: "Преподаватель", Fio: "Вячеслав Белый", GroupNumber: "ИВТ-232"},
	{GithubId: "11", TgId: "14227hSh", Roles: "Студент,Администратор", Fio: "Иннокентий Васильев", GroupNumber: "ИВТ-232"},
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.GET("/users/:id", getUserByID)
	router.GET("/users/:id/roles", getRolesByID)
	router.POST("/users", postUsers)
	router.DELETE("/users/:id", delUsers)
	router.POST("/users/:id/edit_roles", editRoles)
	router.POST("/users/:id/edit_user", editUsers)

	router.Run("localhost:8080")
}

// getUsers responds with the list of all users as JSON.
func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

// postUsers adds an user from JSON received in the request body.
func postUsers(c *gin.Context) {
	var newUser user

	// Call BindJSON to bind the received JSON to newUser.
	err := c.BindJSON(&newUser)
	if err != nil {
		return
	}

	// Add the new user to the slice.
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func delUsers(c *gin.Context) {
	id := c.Param("id")
	for index, u := range users {
		if (u.TgId == id) || (u.GithubId == id) {
			users = slices.Delete(users, index, index+1)
			c.IndentedJSON(http.StatusOK, u)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})

}

// getUserByID locates the user whose TgID value or GithubId value matches the id
// parameter sent by the client, then returns that user as a response.
func getUserByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of users, looking for
	// an user whose TgID or GithubId value matches the parameter.
	for _, u := range users {
		if (u.TgId == id) || (u.GithubId == id) {
			c.IndentedJSON(http.StatusOK, u)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}
func getRolesByID(c *gin.Context) {
	id := c.Param("id")
	for _, u := range users {
		if (u.TgId == id) || (u.GithubId == id) {
			c.IndentedJSON(http.StatusOK, u.Roles)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "role not found"})
}

// TODO: fix it
func editRoles(c *gin.Context) {
	id := c.Param("id")

	for index, u := range users {
		if (u.TgId == id) || (u.GithubId == id) {
			var newUser user
			err := c.BindJSON(&newUser)
			users[index].Roles = newUser.Roles

			if err != nil {
				return
			}
			c.IndentedJSON(http.StatusOK, newUser.Roles)
			return
		}

	}

	c.IndentedJSON(http.StatusNotFound, errorMessage{Message: "role not found"})

}
func editUsers(c *gin.Context) {
	id := c.Param("id")

	for index, u := range users {
		if (u.TgId == id) || (u.GithubId == id) {
			var newUser user
			err := c.BindJSON(&newUser)
			users[index] = newUser

			if err != nil {
				return
			}
			c.IndentedJSON(http.StatusOK, newUser)
			return
		}

	}

	c.IndentedJSON(http.StatusNotFound, errorMessage{Message: "user not found"})

}
