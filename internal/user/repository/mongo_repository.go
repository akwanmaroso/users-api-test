package repository

import (
	"context"

	"github.com/akwanmaroso/users-api/internal/models"
	"github.com/akwanmaroso/users-api/internal/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepositoryImpl struct {
	db *mongo.Collection
}

func NewMongoRepository(db *mongo.Collection) user.Repository {
	return &MongoRepositoryImpl{
		db: db,
	}
}

func (m *MongoRepositoryImpl) Create(ctx context.Context, user *models.User) (*models.User, error) {
	// var result models.User
	res, err := m.db.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	if lastObjectID, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = lastObjectID
	}

	return user, nil
}

func (m *MongoRepositoryImpl) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var result models.User

	filter := bson.M{"username": username}
	docs := m.db.FindOne(ctx, filter)

	if err := docs.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *MongoRepositoryImpl) GetByID(ctx context.Context, id string) (*models.User, error) {
	var result models.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	docs := m.db.FindOne(ctx, filter)

	if err = docs.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *MongoRepositoryImpl) List(ctx context.Context) ([]*models.User, error) {
	result := make([]*models.User, 0)
	cursor, err := m.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (m *MongoRepositoryImpl) Edit(ctx context.Context, user *models.User) (*models.User, error) {
	var result models.User

	after := options.After
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	objID, err := primitive.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"username":   user.Username,
			"role":       user.Role,
			"updated_at": user.UpdatedAt,
		},
	}

	err = m.db.FindOneAndUpdate(ctx, filter, update, &opts).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *MongoRepositoryImpl) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	_, err = m.db.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
