// Unrated Coder t.me/Unrated_Coder

package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client   *mongo.Client
	users    *mongo.Collection
	admins   *mongo.Collection
}

func NewDatabase(uri string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := client.Database("mediainfo_bot")
	return &Database{
		client: client,
		users:  db.Collection("users"),
		admins: db.Collection("admins"),
	}, nil
}

func (db *Database) Ping(ctx context.Context) error {
	return db.client.Ping(ctx, nil)
}

func (db *Database) AddUser(ctx context.Context, userID int64, username string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"username": username,
			"banned":   false,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := db.users.UpdateOne(ctx, filter, update, opts)
	return err
}

func (db *Database) IsBanned(ctx context.Context, userID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	var result struct {
		Banned bool `bson:"banned"`
	}
	err := db.users.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result.Banned, nil
}

func (db *Database) BanUser(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"banned": true,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := db.users.UpdateOne(ctx, filter, update, opts)
	return err
}

func (db *Database) UnbanUser(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"banned": false,
		},
	}
	_, err := db.users.UpdateOne(ctx, filter, update)
	return err
}

func (db *Database) GetAllUsersCount(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return db.users.CountDocuments(ctx, bson.M{})
}

func (db *Database) AddAdmin(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"is_admin": true,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := db.admins.UpdateOne(ctx, filter, update, opts)
	return err
}

func (db *Database) IsAdmin(ctx context.Context, userID int64, ownerID int64) (bool, error) {
	if userID == ownerID {
		return true, nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	var result struct {
		IsAdmin bool `bson:"is_admin"`
	}
	err := db.admins.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return result.IsAdmin, nil
}

func (db *Database) RemoveAdmin(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	_, err := db.admins.DeleteOne(ctx, filter)
	return err
}

func (db *Database) GetAllAdmins(ctx context.Context) ([]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := db.admins.Find(ctx, bson.M{"is_admin": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var admins []int64
	for cursor.Next(ctx) {
		var res struct {
			UserID int64 `bson:"user_id"`
		}
		if err := cursor.Decode(&res); err == nil {
			admins = append(admins, res.UserID)
		}
	}
	return admins, nil
}
