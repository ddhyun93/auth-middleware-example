package mongo

import (
	"context"
	"go-auth-with-chi/domain"
	"go-auth-with-chi/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type userRepo struct {
	col *mongo.Collection
}

func (repo *userRepo) Get(ID string) (*domain.UserDAO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{"_id": objectID}
	var user *domain.UserDAO
	if err := repo.col.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepo) GetByEmail(email string) (*domain.UserDAO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"email": email}
	var user *domain.UserDAO
	if err := repo.col.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepo) Upsert(user *domain.UserDAO) (*domain.UserDAO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Update().SetUpsert(true)

	var updatedUserID primitive.ObjectID
	if user.ID != primitive.NilObjectID {
		user.UpdatedAt = time.Now()
		// update
		filter := bson.M{"_id": user.ID}
		if _, err := repo.col.UpdateOne(ctx, filter, bson.M{"$set": &user}, opts); err != nil {
			return nil, err
		}
		updatedUserID = user.ID
	} else {
		// insert
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.ID = primitive.NewObjectID()
		updatedUserID = user.ID
		_, err := repo.col.InsertOne(ctx, *user)
		if err != nil {
			return nil, err
		}
	}

	var upsertedUser *domain.UserDAO
	filter := bson.M{"_id": updatedUserID}
	if err := repo.col.FindOne(ctx, filter).Decode(&upsertedUser); err != nil {
		return nil, err
	}

	return upsertedUser, nil
}

func (repo *userRepo) Destroy() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := repo.col.DeleteMany(ctx, bson.M{})
	return err == nil
}

func GetMongoUserRepo(conn *MongoDB) repository.UserRepository {
	return &userRepo{
		col: conn.userCol,
	}
}
