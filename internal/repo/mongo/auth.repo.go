package mon

import (
	"authservice/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Alright, so we make a lot of these :D
// this is how we interface we with the database from our applications.

// first step is, make a struct with the fields db and collection, these will store the state
// of the database and allow us to use that connection thing we made earlier
// give it a shot

type AuthRepo struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewAuthRepo(client *mongo.Client, dbName string) *AuthRepo {
	db := client.Database(dbName)
	collection := db.Collection("users")

	return &AuthRepo{
		db:         db,
		collection: collection,
	}
}

func (ur *AuthRepo) Store(ctx context.Context, User domain.Auth) error {
	_, err := ur.collection.InsertOne(ctx, User)
	if err != nil {
		return err
	}
	return nil
}

func (ur *AuthRepo) Remove(ctx context.Context, userID domain.UserID) error {
	_, err := ur.collection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return err
	}
	return nil
}

func (ur *AuthRepo) FetchUserByID(ctx context.Context, userId domain.UserID) (*domain.Auth, error) {
	var auth domain.Auth
	err := ur.collection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&auth)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}
