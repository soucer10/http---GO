package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Createserver() {
	fmt.Println("Iniciando o Servidor HTTP.......")
	server := echo.New()
	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "A minha magrelinha á linda!!!!!!")
	})

	server.GET("/cadastrar", Cadastro_User)
	server.GET("/cadastrar1", Logon_User)

	server.Logger.Fatal(server.Start(":3333"))

}

func Cadastro_User(c echo.Context) error {

	nome := "Hugo3"
	email := "hugobicudo3@gmail.com"
	password := "A5gl4hx7"

	//nome := c.FormValue("nome")
	//email := c.FormValue("email")
	//password := c.FormValue("password")

	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		fmt.Println(err)
	}
	User := User{nome, email, string(hash)}

	_, err1, aux := User.Include()

	fmt.Println(aux.InsertedID)
	if err1 != nil {
		return c.JSON(401, Message{err1.Error()})
	}

	User.Password = ""
	return c.JSON(200, User)
}
func Logon_User(c echo.Context) error {

	nome := "Hugo"
	email := "hugobicudo@gmail.com"
	password := "A5gl4hx7"

	user := User{nome, email, password}

	//nome := c.FormValue("nome")
	//email := c.FormValue("email")
	//password := c.FormValue("password")
	logon, err, aux := user.Logon()

	if !logon || err != nil {

		return c.JSON(401, Message{err.Error()})

	}
	id, _ := aux.DecodeBytes()
	token, _ := user.CreateToken(id)
	return c.JSON(200, Login{token, "Login realizado com sucesso"})
}
