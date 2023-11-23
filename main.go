package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"slices"
)

// user представляет данные о пользователе.
type User struct {
	GithubId    string `json:"github_id" db:"github_id"`
	TgId        string `json:"telegram_id" db:"telegram_id"`
	Roles       string `json:"roles" db:"roles"`
	Fio         string `json:"fio" db:"fio"`
	GroupNumber string `json:"group_number" db:"group_number"`
}

type errorMessage struct {
	Message string `json:"message"`
}

// слайс для заполнения данных о пользователях.
var users = []User{
	{GithubId: "11242", TgId: "8837hSh", Roles: "Студент", Fio: "Иван Кононский", GroupNumber: "ИВТ-232"},
	{GithubId: "1242", TgId: "49957hSh", Roles: "Преподаватель", Fio: "Вячеслав Белый", GroupNumber: "ИВТ-232"},
	{GithubId: "11", TgId: "14227hSh", Roles: "Студент,Администратор", Fio: "Иннокентий Васильев", GroupNumber: "ИВТ-232"},
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

var config = Config{
	Host:     "localhost",
	Port:     "5432",
	Username: "postgres",
	Password: "postgres1G5",
	DBName:   "postgres",
	SSLMode:  "disable",
}

func main() {
	db, err := ConnectDB(config)

	if err != nil {
		fmt.Println("Connection to DB fail", err)
		log.Fatal(err)
	}

	user := User{}
	rows, err := db.Queryx("SELECT * FROM users")
	for rows.Next() {
		err := rows.StructScan(&user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%#v\n", user)
	}

	//router := gin.Default()
	//router.GET("/users", getUsers)
	//router.GET("/users/:id", getUserByID)
	//router.GET("/users/:id/roles", getRolesByID)
	//router.POST("/users", postUsers)
	//router.DELETE("/users/:id", delUsers)
	//router.POST("/users/:id/roles", editRoles)
	//router.POST("/users/:id", editUsers)
	//
	//router.Run("localhost:8080")
}

// getUsers выдает список всех пользователей в формате JSON.
func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

// postUsers добавляет пользователя из JSON, полученного в теле запроса.
func postUsers(c *gin.Context) {
	var newUser User

	// Вызываем BindJSON, чтобы привязать полученный JSON к новому пользователю.
	err := c.BindJSON(&newUser)
	if err != nil {
		return
	}

	// Добавляем нового пользователя в слайс.
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

// delUsers удаляет пользователя из слайса.
func delUsers(c *gin.Context) {
	id := c.Param("id")
	for index, u := range users {
		if (u.TgId == id) || (u.GithubId == id) {
			users = slices.Delete(users, index, index+1)
			c.IndentedJSON(http.StatusOK, u)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, errorMessage{Message: "user not found"})
}

// getUserByID проверяет существование пользователя.
func getUserByID(c *gin.Context) {
	id := c.Param("id")
	for _, u := range users {
		if (u.TgId == id) || (u.GithubId == id) {
			c.IndentedJSON(http.StatusOK, u)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, errorMessage{Message: "user not found"})
}

// getRolesByID проверяет роль пользователя.
func getRolesByID(c *gin.Context) {
	id := c.Param("id")
	for _, u := range users {
		if (u.TgId == id) || (u.GithubId == id) {
			c.IndentedJSON(http.StatusOK, u.Roles)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, errorMessage{Message: "user not found"})
}

// editRoles изменяет роль пользователя на роль из JSON, полученного в теле запроса.
func editRoles(c *gin.Context) {
	id := c.Param("id")
	for index, u := range users {
		if (u.TgId == id) || (u.GithubId == id) {
			var newUser User
			err := c.BindJSON(&newUser)
			users[index].Roles = newUser.Roles
			if err != nil {
				return
			}
			c.IndentedJSON(http.StatusOK, newUser.Roles)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, errorMessage{Message: "user not found"})
}

// editUsers изменяет данные пользователя на данные из JSON, полученного в теле запроса.
func editUsers(c *gin.Context) {
	id := c.Param("id")
	for index, u := range users {
		if (u.TgId == id) || (u.GithubId == id) {
			var newUser User
			err := c.BindJSON(&newUser)
			users[index].Fio = newUser.Fio
			users[index].GroupNumber = newUser.GroupNumber
			if err != nil {
				return
			}
			c.IndentedJSON(http.StatusOK, users[index])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, errorMessage{Message: "user not found"})
}

func ConnectDB(cfg Config) (*sqlx.DB, error) {

	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port,
		cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
