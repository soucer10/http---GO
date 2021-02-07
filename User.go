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

var Client1 mongo.Client

// Conectar_mongo -  Conectar no banco de dados
func Conectar_mongo() {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	Client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://soucer10:soucer10@cluster0.fwaym.mongodb.net/<dbname>?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}

	Client1 = *Client
	fmt.Println("Conectado ao Banco....")
}

// Connect - responsavel pelo client
type Token struct {
	JWT string `json:"token"`
}

type Message struct {
	Msg string `json:"message" xml:"message" form:"message" query:"message"`
}

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
func (u User) Include() (bool, error, *mongo.InsertOneResult) {

	collection := Client1.Database("Golang").Collection("Users")

	var dbaux User
	filter := bson.D{{"email", u.Email}}
	err1 := collection.FindOne(context.TODO(), filter).Decode(&dbaux)

	if err1 == nil {
		return false, errors.New("Usuário já está cadastrado"), nil
	}

	aux, err := collection.InsertOne(context.TODO(), u)

	if err != nil {

		return false, errors.New("Erro...."), nil

	}
	return true, nil, aux
}

func (u User) Logon() (bool, error, *mongo.SingleResult) {

	collection := Client1.Database("Golang").Collection("Users")
	var dbaux User
	filter := bson.D{{"email", u.Email}}

	err := collection.FindOne(context.TODO(), filter).Decode(&dbaux)
	if err != nil {

		return false, errors.New("Usuário não encontrado! \n"), nil
	}

	err1 := bcrypt.CompareHashAndPassword([]byte(dbaux.Password), []byte(u.Password))
	if err1 != nil {

		return false, errors.New("Senha está errada!\n"), nil
	}
	aux := collection.FindOne(context.TODO(), filter)
	return true, nil, aux
}

func (u User) CreateToken(b bson.Raw) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = b
	atClaims["exp"] = time.Now().Add(time.Minute * 24 * 60).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
