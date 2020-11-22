package mongo

import (
	"atlant/service/dto"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductsActions struct {
	client     *mongo.Client
	timeout    time.Duration
	collection *mongo.Collection
}

func (p ProductsActions) page(start uint64, leng int64, opts *options.FindOptions) ([]*dto.Product, error) {
	opts.SetSkip(int64(start))
	if leng > 0 {
		opts.SetLimit(leng)
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	cursor, err := p.collection.Find(ctx, bson.M{}, opts)
	defer cursor.Close(context.Background())
	if err != nil {
		log.Printf("Can't get data from DB, %v", err)
		return []*dto.Product{}, nil
	}
	var results []*dto.Product
	for cursor.Next(ctx) {

		// create a value into which the single document can be decoded
		var elem dto.Product
		err := cursor.Decode(&elem)
		if err != nil {
			log.Printf("Can't deserialized data: %v", err)
			return nil, err
		}
		results = append(results, &elem)
	}

	return results, nil
}

func (p ProductsActions) Save(product string, price float64) error {
	opts := options.Update().SetUpsert(true)

	filter := bson.D{{"product", product}}
	update := bson.D{
		{"$inc", bson.D{
			{"change_count", 1},
		}},
		{"$currentDate", bson.D{
			{"change_date", bson.D{
				{"$type", "timestamp"},
			}},
		}},
		{"$set", bson.D{
			{"price", price},
		}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	result, err := p.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error during inserting document: %v", err)
		return err
	}
	log.Printf("Upserted %d, Modified %d", result.UpsertedCount, result.ModifiedCount)
	return nil
}

func (p ProductsActions) LoadByProduct(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
	opts := options.Find()
	direction := -1
	if upSort {
		direction = 1
	}
	opts.SetSort(bson.M{
		"product": direction,
	})
	return p.page(start, leng, opts)
}
func (p ProductsActions) LoadByPrice(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
	opts := options.Find()
	direction := -1
	if upSort {
		direction = 1
	}
	opts.SetSort(bson.M{
		"price": direction,
	})
	return p.page(start, leng, opts)
}
func (p ProductsActions) LoadByChangeCount(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
	opts := options.Find()
	direction := -1
	if upSort {
		direction = 1
	}
	opts.SetSort(bson.M{
		"change_count": direction,
	})
	return p.page(start, leng, opts)
}
func (p ProductsActions) LoadByDate(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
	opts := options.Find()
	direction := -1
	if upSort {
		direction = 1
	}
	opts.SetSort(bson.M{
		"change_date": direction,
	})
	return p.page(start, leng, opts)
}

func NewProductsActions(connection string, timeout time.Duration) ProductsActions {
	client, err := mongo.NewClient(options.Client().ApplyURI(connection))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("atlant").Collection("products")

	log.Printf("Connected to MongoDB: %v", connection)
	return ProductsActions{
		client:     client,
		timeout:    timeout,
		collection: collection,
	}
}
