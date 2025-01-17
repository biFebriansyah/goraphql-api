package service

import (
	"context"
	"fmt"
	"log"
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

	cur, err := user.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	result := []*model.Users{}

	for cur.Next(ctx) {
		var data model.Users
		if err := cur.Decode(&data); err != nil {
			log.Fatal(err)
		}

		result = append(result, &data)
	}

	return result, nil
}

func (user *UserService) GetById(userId string) (*model.Users, error) {
	result := new(model.Users)
	obectId, _ := bson.ObjectIDFromHex(userId)
	if err := user.FindOne(context.TODO(), bson.M{"_id": obectId}).Decode(result); err != nil {
		return nil, err
	}

	return result, nil

}

func (user *UserService) CreateOne(data model.Users) (*model.Users, error) {
	res, err := user.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}

	fmt.Println(res)

	return &data, nil
}
