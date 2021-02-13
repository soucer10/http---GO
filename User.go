package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

//Client1 - Ponteiro para conexão do banco
var Client1 mongo.Client

// Conectarmongo -  Conectar no banco de dados
func Conectarmongo() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://soucer10:soucer10@cluster0.fwaym.mongodb.net/<dbname>?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}

	Client1 = *Client

	fmt.Println("Conectado ao Banco....")
}

// Token - responsavel pela autênticação
type Token struct {
	JWT string `json:"token"`
}

//Message =  envio de mensagens
type Message struct {
	Msg string `json:"message" xml:"message" form:"message" query:"message"`
}

//Login - Model user
type Login struct {
	Token string `json:"token" xml:"token" form:"token" query:"token"`
	Msg   string `json:"message" xml:"message" form:"message" query:"message"`
}

//User model
type User struct {
	Nome     string `json:"nome" xml:"nome" form:"nome" query:"nome"`
	Email    string `json:"email" xml:"email" form:"email" query:"email"`
	Password string `json:"password" xml:"password" form:"password" query:"password"`
}

//Include - Incluir usuário
func (u User) Include() (bool, *mongo.InsertOneResult, error) {

	collection := Client1.Database("Golang").Collection("Users")

	var dbaux User
	filter := bson.D{{"email", u.Email}}
	err1 := collection.FindOne(context.TODO(), filter).Decode(&dbaux)

	if err1 == nil {
		return false, nil, errors.New("Usuário já está cadastrado")
	}

	aux1, err := collection.InsertOne(context.TODO(), u)

	if err != nil {

		return false, nil, errors.New("Erro")

	}
	return true, aux1, nil
}

//Logon realiza o login
func (u User) Logon() (bool, *mongo.SingleResult, error) {

	collection1 := Client1.Database("Golang").Collection("Users")
	var dbaux User
	filter := bson.D{{"email", u.Email}}

	err := collection1.FindOne(context.TODO(), filter).Decode(&dbaux)
	if err != nil {

		return false, nil, errors.New("Usuário não encontrado")
	}

	err1 := bcrypt.CompareHashAndPassword([]byte(dbaux.Password), []byte(u.Password))
	if err1 != nil {

		return false, nil, errors.New("Senha está errada")
	}
	aux := collection1.FindOne(context.TODO(), filter)
	return true, aux, nil
}

//CreateToken - Criação do Token
func (u User) CreateToken(b bson.Raw) (string, error) {

	//Creating Access Token

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = b.String()
	atClaims["exp"] = time.Now().Add(time.Minute * 24 * 60).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
