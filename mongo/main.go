package mongo

import (
	"context"
	"go-auth-with-chi/ioc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const MONGOURI = "mongodb://localhost:27017/authSampleApp"

type MongoDB struct {
	userCol *mongo.Collection
}

func NewMongoDB(env string) *MongoDB {
	userCol := "users"

	if env == "TEST" {
		userCol += "_test"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := makeMongoClient(ctx)
	if err != nil {
		log.Panicln("mongodb connection failed with err : ", err)
	}
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Panicln("mongodb no response with err : ", err)
	}

	db := mongoClient.Database("authSampleApp")
	return &MongoDB{
		userCol: db.Collection(userCol),
	}
}

func makeMongoClient(ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(MONGOURI).SetAuth(options.Credential{
		Username: "donghyun",
		Password: "donghyun",
	})
	mongoClient, err := mongo.Connect(ctx, clientOptions)

	return mongoClient, err
}

func (conn *MongoDB) RegisterRepos() {
	ioc.Repo.Users = GetMongoUserRepo(conn)
}
