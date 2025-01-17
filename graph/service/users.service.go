package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/biFebriansyah/goraphql/graph/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserService struct {
	*mongo.Collection
}

func NewUserService(cln *mongo.Collection) *UserService {
	return &UserService{cln}
}

func (user *UserService) GetAll() ([]*model.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := user.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	defer cur.Close(ctx)
	result := []*model.Users{}

	for cur.Next(ctx) {
		var data *model.Users
		decoder := bson.NewDecoder(bson.NewDocumentReader(bytes.NewReader(cur.Current)))
		decoder.ObjectIDAsHexString()
		if err := decoder.Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}

		result = append(result, data)
	}

	return result, nil
}

func (user *UserService) GetById(userId string) (*model.Users, error) {
	result := new(model.Users)
	obectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	raw, err := user.FindOne(context.TODO(), bson.M{"_id": obectId}).Raw()
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	decoder := bson.NewDecoder(bson.NewDocumentReader(bytes.NewReader(raw)))
	decoder.ObjectIDAsHexString()
	if err := decoder.Decode(result); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return result, nil

}

func (user *UserService) CreateOne(data model.NewUser) (*model.Users, error) {
	res, err := user.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	datas := model.Users{
		ID:    res.InsertedID.(bson.ObjectID).Hex(),
		Name:  data.Name,
		Email: data.Email,
	}

	return &datas, nil
}
