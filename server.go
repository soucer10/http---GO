package main

import (
	f "fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

//Pessoa - estrutura do usuário
type Pessoa struct {
	Nome string `json:"nome" xml:"nome" form:"nome" query:"nome"`

	Email string `josn:"email" xml:"email" form:"email" query:"email"`

	Senha string `json:"senha" xml:"senha" form:"senha" query:"senha"`
}

//Createserver responsavel pela criação do Servidor
func Createserver() {

	server := echo.New()
	server.Debug = true
	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "A minha magrelinha á linda!!!!!!")
	})
	server.Use(Process)
	server.GET("/user", getuser)
	server.GET("/saudacao/:nome", saudacao)

	server.Logger.Fatal(server.Start(":3333"))

}

func getuser(c echo.Context) error {

	pessoa := Pessoa{"Hugo", "hugobicudo@gmail.com", "senha"}
	f.Println(pessoa.Email)
	return c.JSON(http.StatusOK, pessoa)
}

func saudacao(c echo.Context) error {
	nome := c.Param("nome")
	mensage := "Seja bem Vindo " + nome
	return c.JSON(http.StatusOK, "mensage:"+mensage)
}

//Process - Validar existecia do token Middleware
func Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		endpoint := c.Path()

		if endpoint == "/" {
			return next(c)
		}

		f.Println(endpoint)

		token := c.Request().Header.Get("token")

		if token == "" || len(token) == 0 {
			return c.String(401, "Você é safado!!!!!!!!!")
		}
		return next(c)

	}
}
