package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID         string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string    `bson:"name" json:"name"`
	Data       string    `bson:"data" json:"data"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpadatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:       entry.Name,
		Data:       entry.Data,
		CreatedAt:  time.Now(),
		UpadatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting:", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cusror, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Println("Error finding all logs:", err)
		return nil, err
	}
	defer cusror.Close(ctx)

	var logs []*LogEntry
	for cusror.Next(ctx) {
		var item LogEntry
		err := cusror.Decode(&item)
		if err != nil {
			log.Println("Error decoding:", err)
			return nil, err
		}
		logs = append(logs, &item)
	}
	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting id:", err)
		return nil, err
	}

	var item LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&item)
	if err != nil {
		log.Println("Error finding log:", err)
		return nil, err
	}
	return &item, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	err := collection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping collection:", err)
		return err
	}
	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("Error converting id:", err)
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx, 
		bson.M{"_id": docID}, 
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		log.Println("Error updating:", err)
		return nil, err
	}

	return result, nil
}