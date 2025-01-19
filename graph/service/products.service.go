package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/biFebriansyah/goraphql/graph/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductService struct {
	*mongo.Collection
}

func NewProductService(cln *mongo.Collection) *ProductService {
	return &ProductService{cln}
}

type FacetProduct struct {
	Metadata []struct {
		TotalCount int `bson:"totalCount"`
	} `bson:"metadata"`
	Data []*model.Products `bson:"data"`
}

func (product *ProductService) GetAll(page, limit int64, name *string) (*model.ProductDetail, error) {
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
				bson.M{"$skip": (page - 1) * limit},
				bson.M{"$limit": limit}},
		},
	}}

	cursor, err := product.Aggregate(ctx, mongo.Pipeline{matchParams, qurParams})
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	var results []FacetProduct
	for cursor.Next(ctx) {
		var data FacetProduct
		decoder := bson.NewDecoder(bson.NewDocumentReader(bytes.NewReader(cursor.Current)))
		decoder.ObjectIDAsHexString()
		if err := decoder.Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}

		results = append(results, data)
	}

	if len(results[0].Data) <= 0 {
		return nil, fmt.Errorf("failed find data product: %w", errors.New("data not found"))
	}

	totalCount := results[0].Metadata[0].TotalCount
	productMeta := model.ProductMeta{Total: int32(totalCount)}
	if totalCount > 0 {
		cek := int32(math.Ceil(float64(totalCount) / float64(limit)))
		if int32(page) >= cek {
			productMeta.Next = 0
		} else {
			productMeta.Next = int32(page + 1)
		}
	}
	if page > 1 {
		productMeta.Prev = int32(page - 1)
	}

	userData := model.ProductDetail{
		Data: results[0].Data,
		Meta: &productMeta,
	}

	return &userData, nil

}

func (product *ProductService) GetById(productId string) (*model.Products, error) {
	result := new(model.Products)
	obectId, err := bson.ObjectIDFromHex(productId)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID format: %w", err)
	}

	raw, err := product.FindOne(context.TODO(), bson.M{"_id": obectId}).Raw()
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	decoder := bson.NewDecoder(bson.NewDocumentReader(bytes.NewReader(raw)))
	decoder.ObjectIDAsHexString()
	if err := decoder.Decode(result); err != nil {
		return nil, fmt.Errorf("failed to decode product: %w", err)
	}

	return result, nil

}

func (product *ProductService) CreateOne(data model.NewProduct) (*model.Products, error) {
	res, err := product.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, fmt.Errorf("failed to insert product: %w", err)
	}

	datas := model.Products{
		ID:    res.InsertedID.(bson.ObjectID).Hex(),
		Name:  data.Name,
		Price: data.Price,
		Stock: data.Stock,
	}

	return &datas, nil
}

func (product *ProductService) UpdateOne(data model.UpdateProduct) (*model.Products, error) {
	obectId, err := bson.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	updateData := bson.M{}
	marshal, _ := json.Marshal(data)
	err = json.Unmarshal(marshal, &updateData)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal struct: %w", err)
	}

	delete(updateData, "_id")
	res, err := product.FindOneAndUpdate(context.TODO(), bson.M{"_id": obectId}, bson.M{"$set": updateData}).Raw()
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	result := new(model.Products)
	decoder := bson.NewDecoder(bson.NewDocumentReader(bytes.NewReader(res)))
	decoder.ObjectIDAsHexString()
	if err := decoder.Decode(result); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return result, nil

}

func (product *ProductService) DeleteOne(productId string) (string, error) {
	obectId, err := bson.ObjectIDFromHex(productId)
	if err != nil {
		return "nil", fmt.Errorf("invalid user ID format: %w", err)
	}

	res := product.FindOneAndDelete(context.TODO(), bson.M{"_id": obectId})
	if res.Err() != nil {
		return "nil", fmt.Errorf("failed to delete user: %w", err)
	}

	return productId, nil

}
