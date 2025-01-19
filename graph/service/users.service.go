package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/biFebriansyah/goraphql/graph/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserService struct {
	*mongo.Collection
}

type FacetResult struct {
	Metadata []struct {
		TotalCount int `bson:"totalCount"`
	} `bson:"metadata"`
	Data []*model.Users `bson:"data"`
}

func NewUserService(cln *mongo.Collection) *UserService {
	return &UserService{cln}
}

func (user *UserService) GetAll(page, limit int64, name *string) (*model.UsersDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	serachName := ""
	if name != nil {
		serachName = *name
	}

	matchParams := bson.D{{Key: "$match", Value: bson.D{{Key: "name", Value: bson.M{"$regex": serachName, "$options": "i"}}}}}
	qurParams := bson.D{{
		Key: "$facet",
		Value: bson.M{
			"metadata": bson.A{bson.M{"$count": "totalCount"}},
			"data": bson.A{
				bson.M{"$unset": bson.A{"password"}},
				bson.M{"$skip": (page - 1) * limit},
				bson.M{"$limit": limit}},
		},
	}}

	cursor, err := user.Aggregate(ctx, mongo.Pipeline{matchParams, qurParams})
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	var results []FacetResult
	for cursor.Next(ctx) {
		var data FacetResult
		decoder := bson.NewDecoder(bson.NewDocumentReader(bytes.NewReader(cursor.Current)))
		decoder.ObjectIDAsHexString()
		if err := decoder.Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}

		results = append(results, data)
	}

	if len(results[0].Data) <= 0 {
		return nil, fmt.Errorf("failed find data user: %w", errors.New("data not found"))
	}

	totalCount := results[0].Metadata[0].TotalCount
	userMeta := model.UserMeta{Total: int32(totalCount)}
	if totalCount > 0 {
		cek := int32(math.Ceil(float64(totalCount) / float64(limit)))
		if int32(page) >= cek {
			userMeta.Next = 0
		} else {
			userMeta.Next = int32(page + 1)
		}
	}
	if page > 1 {
		userMeta.Prev = int32(page - 1)
	}

	userData := model.UsersDetail{
		Data: results[0].Data,
		Meta: &userMeta,
	}

	return &userData, nil
}

func (user *UserService) GetAllold(page, limit int64) ([]*model.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	qurOption := options.Find().SetLimit(limit).SetSkip(page)
	cur, err := user.Find(ctx, bson.M{}, qurOption)
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

func (user *UserService) GetByEmail(email string) (*model.Users, error) {
	result := new(model.Users)
	raw, err := user.FindOne(context.TODO(), bson.M{"email": email}).Raw()
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

func (user *UserService) CreateOne(data model.SignupInput) (*model.Users, error) {
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

func (user *UserService) UpdateOne(data model.UpdateInput) (*model.Users, error) {
	obectId, err := bson.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	updateData := model.SignupInput{Name: data.Name, Email: *data.Email, Password: *data.Password}
	res, err := user.FindOneAndUpdate(context.TODO(), bson.M{"_id": obectId}, bson.M{"$set": updateData}).Raw()
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	result := new(model.Users)
	decoder := bson.NewDecoder(bson.NewDocumentReader(bytes.NewReader(res)))
	decoder.ObjectIDAsHexString()
	if err := decoder.Decode(result); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return result, nil

}

func (user *UserService) DeleteOne(userId string) (string, error) {
	obectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return "nil", fmt.Errorf("invalid user ID format: %w", err)
	}

	res := user.FindOneAndDelete(context.TODO(), bson.M{"_id": obectId})
	if res.Err() != nil {
		return "nil", fmt.Errorf("failed to delete user: %w", err)
	}

	return userId, nil

}
