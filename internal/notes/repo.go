package notes

import (
	"context"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("note not found")

type Repo struct {
	col *mongo.Collection
}

func NewRepo(db *mongo.Database) (*Repo, error) {
	col := db.Collection("notes")
	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: 1}},
		Options: options.Index().SetUnique(true)})
	if err != nil {
		return nil, err
	}
	return &Repo{col: col}, nil
}
func (r *Repo) Create(ctx context.Context, title, content string) (Note, error) {
	now := time.Now()
	n := Note{Title: title, Content: content, CreatedAt: now, UpdatedAt: now}
	res, err := r.col.InsertOne(ctx, n)
	if err != nil {
		return Note{}, err
	}
	n.ID = res.InsertedID.(primitive.ObjectID)
	return n, nil
}
func (r *Repo) ByID(ctx context.Context, idHex string) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return Note{}, ErrNotFound
	}
	var n Note
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&n); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}
	return n, nil
}
func (r *Repo) List(ctx context.Context, q string, limit int64, after string) ([]Note, error) {
	filter := bson.M{}
	if q != "" {
		filter["title"] = bson.M{"$regex": q, "$options": "i"}
	}

	if after != "" {
		afterID, err := primitive.ObjectIDFromHex(after)
		if err != nil {
			return nil, err
		}
		filter["_id"] = bson.M{"$lt": afterID}
	}

	opts := options.Find().SetLimit(limit).SetSort(bson.D{{Key: "_id", Value: -1}})
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []Note
	for cur.Next(ctx) {
		var n Note
		if err := cur.Decode(&n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, cur.Err()
}
func (r *Repo) Update(ctx context.Context, idHex string, title, content *string) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return Note{}, ErrNotFound
	}
	set := bson.M{"updatedAt": time.Now()}
	if title != nil {
		set["title"] = *title
	}
	if content != nil {
		set["content"] = *content
	}
	after := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated Note
	if err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": set}, after).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}
	return updated, nil
}
func (r *Repo) Delete(ctx context.Context, idHex string) error {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return ErrNotFound
	}
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) Stats(ctx context.Context) (map[string]interface{}, error) {
	pipeline := mongo.Pipeline{
		// Группируем все документы
		{{Key: "$group", Value: bson.M{
			"_id":              nil,
			"totalNotes":       bson.M{"$sum": 1},
			"avgContentLength": bson.M{"$avg": bson.M{"$strLenCP": "$content"}},
		}}},
	}

	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []bson.M
	if err = cur.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return map[string]interface{}{
			"totalNotes":       0,
			"avgContentLength": 0,
		}, nil
	}

	stats := results[0]
	ratio := math.Pow(10, float64(2))
	avgContend := math.Round(stats["avgContentLength"].(float64)*ratio) / ratio
	return map[string]any{
		"totalNotes":       stats["totalNotes"],
		"avgContentLength": avgContend,
	}, nil
}
