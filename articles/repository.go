package articles

import (
	"context"
	errors "coodesh/error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const dbName = "Cluster0"
const collectionName = "flightnews"

type NoSqlRepository struct {
	Conn *mongo.Client
}

func (r *NoSqlRepository) getArticles() ([]*FlightNews, error) {
	coll := r.Conn.Database(dbName).Collection(collectionName)
	filter := bson.D{}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, errors.Message{
			Msg:        err.Error(),
			StatusCode: 500,
		}
	}

	defer func(cursor *mongo.Cursor) {
		_ = cursor.Close(context.TODO())
	}(cursor)

	resultItems := make([]*FlightNews, 0)
	for cursor.Next(context.TODO()) {
		var result FlightNews
		if err := cursor.Decode(&result); err != nil {
			return nil, errors.Message{
				Msg:        err.Error(),
				StatusCode: 500,
			}
		}
		resultItems = append(resultItems, &result)
	}

	return resultItems, nil
}

func (r *NoSqlRepository) getArticle(id int) (*FlightNews, error) {
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

func (r *NoSqlRepository) postArticle(body *FlightNews) (*FlightNews, error) {
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

func (r *NoSqlRepository) putArticle(id int, body *FlightNews) (*FlightNews, error) {
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
