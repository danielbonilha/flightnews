package articles

import (
	"context"
	"fmt"
	errors "github.com/danielbonilha/flightnews/error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "Cluster0"
const collectionName = "flightnews"

type NoSqlRepository struct {
	Conn *mongo.Client
}

func (r *NoSqlRepository) getArticles(offset int64, limit int64) ([]FlightNews, error) {
	fmt.Printf("Getting articles [offset=%d][limit=%d]\n", offset, limit)
	coll := r.Conn.Database(dbName).Collection(collectionName)
	opts := options.Find().SetSkip(offset).SetLimit(limit)

	cursor, err := coll.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, errors.Message{
			Msg:        err.Error(),
			StatusCode: 500,
		}
	}

	defer func(cursor *mongo.Cursor) {
		_ = cursor.Close(context.TODO())
	}(cursor)

	resultItems := make([]FlightNews, 0)
	for cursor.Next(context.TODO()) {
		var result FlightNews
		if err := cursor.Decode(&result); err != nil {
			return nil, errors.Message{
				Msg:        err.Error(),
				StatusCode: 500,
			}
		}
		resultItems = append(resultItems, result)
	}

	return resultItems, nil
}

func (r *NoSqlRepository) getArticle(id int) (*FlightNews, error) {
	fmt.Printf("Getting article [id=%d]\n", id)
	coll := r.Conn.Database(dbName).Collection(collectionName)
	var result FlightNews

	err := coll.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, errors.Message{
			Msg:        "Item not found",
			StatusCode: 404,
		}
	}
	if err != nil {
		return nil, errors.Message{
			Msg:        err.Error(),
			StatusCode: 500,
		}
	}

	return &result, nil
}

func (r *NoSqlRepository) insertArticle(body *FlightNews) (*FlightNews, error) {
	fmt.Printf("Persisting article [id=%d]\n", body.Id)
	coll := r.Conn.Database(dbName).Collection(collectionName)
	_, err := coll.InsertOne(context.TODO(), body)
	if err != nil {
		return nil, errors.Message{
			Msg:        err.Error(),
			StatusCode: 500,
		}
	}

	return body, nil
}

func (r *NoSqlRepository) updateArticle(id int, body *FlightNews) (*FlightNews, error) {
	fmt.Printf("Updating article [id=%d]\n", id)
	coll := r.Conn.Database(dbName).Collection(collectionName)
	filter := bson.D{{"id", id}}

	_, err := coll.ReplaceOne(context.TODO(), filter, body)
	if err != nil {
		return nil, errors.Message{
			Msg:        err.Error(),
			StatusCode: 500,
		}
	}

	return body, nil
}

func (r *NoSqlRepository) deleteArticle(id int) error {
	fmt.Printf("Deleting article [id=%d]\n", id)
	coll := r.Conn.Database(dbName).Collection(collectionName)
	filter := bson.D{{"id", id}}

	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return errors.Message{
			Msg:        err.Error(),
			StatusCode: 500,
		}
	}
	return nil
}

func (r *NoSqlRepository) countArticles() (int64, error) {
	fmt.Println("Counting local articles")
	coll := r.Conn.Database(dbName).Collection(collectionName)
	filter := bson.D{}

	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, errors.Message{
			Msg:        err.Error(),
			StatusCode: 500,
		}
	}
	return count, nil
}
