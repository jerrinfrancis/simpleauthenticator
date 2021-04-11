package mongo

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jerrinfrancis/simpleauthenticator/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	client *mongo.Client
	mode   string
}

type userDB struct {
	col *mongo.Collection
}

func (u userDB) Insert(user db.User) (*db.User, error) {
	_, err := u.col.InsertOne(context.TODO(), user)
	return &user, err
}
func (u userDB) InsertForOTP(user db.UserForOTPAuth) (*db.UserForOTPAuth, error) {
	_, err := u.col.InsertOne(context.TODO(), user)
	return &user, err
}

func (u userDB) UpdateVerificationToken(phoneNumber, verificationToken string) (int64, error) {
	var bdoc, set bson.D
	bdoc = append(bdoc, bson.E{Key: "$set", Value: append(set, bson.E{Key: "userlogininfoforotp.verificationtoken", Value: verificationToken})})
	result, error := u.col.UpdateOne(context.Background(), bson.M{"userlogininfoforotp.phonenumber": phoneNumber}, bdoc)
	if error != nil {
		return 0, error
	}
	return result.MatchedCount, nil

}
func (u userDB) Find(userName string) (*db.User, error) {
	var bdoc bson.D

	bdoc = append(bdoc, bson.E{Key: "userlogininfo.username", Value: userName})
	var user db.User
	error := u.col.FindOne(context.Background(), bdoc).Decode(&user)
	if error != nil {
		if error == mongo.ErrNoDocuments {
			log.Println("no usr")
			return nil, nil
		}
		return nil, error
	}
	return &user, nil

}
func (m mongoDB) User() db.UserDB {
	log.Println("m.Mode: ", m.mode)
	if m.mode == "OTP" {
		return userDB{col: m.client.Database("login").Collection("usersotp")}
	}
	return userDB{col: m.client.Database("login").Collection("users")}

}
func (u userDB) FindByPhoneNumber(phoneNumber string) (*db.UserForOTPAuth, error) {
	var bdoc bson.D

	bdoc = append(bdoc, bson.E{Key: "userlogininfoforotp.phonenumber", Value: phoneNumber})
	var user db.UserForOTPAuth
	error := u.col.FindOne(context.Background(), bdoc).Decode(&user)
	if error != nil {
		if error == mongo.ErrNoDocuments {
			log.Println("no usr")
			return nil, nil
		}
		return nil, error
	}
	return &user, nil

}

var client *mongo.Client

func New(mode string) db.DB {
	if client != nil {
		log.Println(mode, "New")
		m := mongoDB{client: client, mode: mode}
		log.Println("set", m.mode)
		return m
	}
	log.Println("creating client :", os.Getenv("USERDB_URL"))
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("USERDB_URL")))
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		cancel()
		log.Fatalln("unable to connect to DB", err)
	}
	defer cancel()

	return mongoDB{client: client, mode: mode}
}
