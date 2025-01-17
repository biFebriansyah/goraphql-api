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

type ProductService struct {
	*mongo.Collection
}

func NewProductService(cln *mongo.Collection) *ProductService {
	return &ProductService{cln}
}

func (product *ProductService) GetAll() ([]*model.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := product.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	defer cur.Close(ctx)
	result := []*model.Products{}

	for cur.Next(ctx) {
		var data *model.Products
		decoder := bson.NewDecoder(bson.NewDocumentReader(bytes.NewReader(cur.Current)))
		decoder.ObjectIDAsHexString()
		if err := decoder.Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}

		result = append(result, data)
	}

	return result, nil
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
