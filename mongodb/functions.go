package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os/exec"
	"time"
)

var Client *mongo.Client

const duration = 360 * time.Second
const timeout = 10 * time.Second

var MyDb = MongoDb{
	Url:      "mongodb://localhost:27017",
	DbName:   "mydb",
	Users:    "users",
	Sessions: "sessions",
}

func GenerateToken() string {
	str, _ := exec.Command("uuidgen").Output()
	return string(str[:len(str)-2])
}

func init() {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	clientOptions := options.Client().ApplyURI(MyDb.Url)
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongodb.")
}

func CheckLogin(username string, password string) error {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Users)

	var account User
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&account); err != nil {
		return err
	}

	if account.Password != password {
		return errors.New("bad credentials")
	}

	return nil
}

func InsertSession(username string) (Session, error) {
	token := GenerateToken()
	session := Session{
		Username:  username,
		Token:     token,
		ExpiresAt: time.Now().Add(duration),
	}

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Sessions)

	if _, err := collection.InsertOne(ctx, session); err != nil {
		return Session{}, err
	}

	return session, nil
}

func GetSession(token string) (Session, error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Sessions)

	var session Session
	if err := collection.FindOne(ctx, bson.M{"token": token}).Decode(&session); err != nil {
		return Session{}, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = RemoveSession(token)
		return Session{}, errors.New("expired token")
	}

	return session, nil
}

func RemoveSession(token string) error {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Sessions)

	if _, err := collection.DeleteOne(ctx, bson.M{"token": token}); err != nil {
		return err
	}

	return nil
}

func GetUser(username string) (User, error) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Users)

	var account User
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&account); err != nil {
		return User{}, err
	}

	return account, nil
}

func AddProfileQuote(username string, quote Quote) error {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Users)

	var account User
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&account); err != nil {
		fmt.Println(err)
		return err
	}

	account.Quotes = append(account.Quotes, quote)
	filter := bson.M{"username": username}
	update := bson.M{"$set": account}

	ctx, err := context.WithTimeout(context.Background(), timeout)
	_ = err
	collection = Client.Database(MyDb.DbName).Collection(MyDb.Users)

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func DeleteProfileQuote(username string, data string) error {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Users)

	var account User
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&account); err != nil {
		return err
	}

	var links []Quote
	for _, quote := range account.Quotes {
		if quote.Data == data {
			continue
		}
		links = append(links, quote)
	}

	account.Quotes = links
	filter := bson.M{"username": username}
	update := bson.M{"$set": account}

	ctx, err := context.WithTimeout(context.Background(), timeout)
	_ = err
	collection = Client.Database(MyDb.DbName).Collection(MyDb.Users)

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func UpdateQuotes(username string, data string) error {
	user, _ := GetUser(username)

	var newQuotes []Quote
	for _, el := range user.Quotes {
		el.Special = false
		if el.Data == data {
			el.Special = true
		}
		newQuotes = append(newQuotes, el)
	}

	user.Quotes = newQuotes
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Users)

	filter := bson.M{"username": username}
	update := bson.M{"$set": user}

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
