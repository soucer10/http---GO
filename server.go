package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

//Createserver -  subir servidor
func Createserver() {
	fmt.Println("Iniciando o Servidor HTTP.......")

	conteudo, _ := ioutil.ReadFile("secret.txt")
	x := strings.Fields(string(conteudo))
	os.Setenv("ACCESS_SECRET", x[0])

	server := echo.New()

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "A minha magrelinha á linda!!!!!!")
	})
	server.GET("/cadastrar", CadastroUser)
	server.GET("/cadastrar1", LogonUser)
	server.Use(valideToken)
	server.GET("/teste", func(c echo.Context) error {
		return c.String(http.StatusOK, "A minha magrelinha não confia em você!!!!!!")
	})

	server.Logger.Fatal(server.Start(":3333"))

}

func valideToken(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		endpoint := c.Path()
		Usertoken := c.Request().Header.Get("token")
		if endpoint == "/" || endpoint == "/cadastrar" || endpoint == "/cadastrar1" {

			return next(c)
		}

		fmt.Println(Usertoken, endpoint)

		token, err := jwt.Parse(Usertoken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Não Passou do Parse")
			}

			return []byte(os.Getenv("ACCESS_SECRET")), nil

		})

		if err != nil {
			return c.JSON(401, Message{Msg: "Deu ruim na verificação"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			User := User{}
			Userjson := claims["user_id"]
			str := fmt.Sprintf("%v", Userjson)

			//byteValueJSON, _ := ioutil.ReadAll(Userjson)
			//json.Unmarshal(byteValueJSON, User)
			bson.Unmarshal([]byte(str), &User)
			fmt.Println(str, User, []byte(str))
		}

		return next(c)

	}

}

//CadastroUser - criação de usuário no Banco
func CadastroUser(c echo.Context) error {

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

	_, _, err1 := User.Include()

	if err1 != nil {
		return c.JSON(401, Message{err1.Error()})
	}

	User.Password = ""
	return c.JSON(200, User)
}

//LogonUser - realiza a autenticação
func LogonUser(c echo.Context) error {

	nome := "Hugo"
	email := "hugobicudo@gmail.com"
	password := "A5gl4hx7"

	user := User{nome, email, password}

	//nome := c.FormValue("nome")
	//email := c.FormValue("email")
	//password := c.FormValue("password")
	logon, aux, err := user.Logon()

	if !logon || err != nil {

		return c.JSON(401, Message{err.Error()})

	}
	id, _ := aux.DecodeBytes()
	token, _ := user.CreateToken(id)
	return c.JSON(200, Login{token, "Login realizado com sucesso"})
}
