package model

import (
	"context"
	"errors"
	"social-api/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollectionName string = "users"

//types here have to implement the  Modeler interface

type UserModel struct {
	Collection *mongo.Collection
}

// simple search when you need to get a entry without any filter options
// will only return single entry
func (um *UserModel) GetEntry(key bson.D) (*types.Users, error) {
	var entry types.Users
	// bson.d is a drivitive of primitive int so cannont be
	// nil, so just check the length of the key
	if len(key) == 0 {
		return nil, errors.New("empty filter given")
	}
	err := um.Collection.FindOne(context.TODO(), key).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (um *UserModel) GetEntryAdvanced(filter bson.D, sort bson.D) ([]*types.Users, error) {
	opts := options.Find().SetSort(sort)
	cur, err := um.Collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	var entrys []*types.Users
	if err = cur.All(context.TODO(), &entrys); err != nil {
		return nil, err
	}
	// gonna return a error if no data return for the given filters
	if len(entrys) == 0 {
		return nil, errors.New("no values found")
	}
	return entrys, nil

}

func (um *UserModel) AddEntry(val bson.D) error {
	if len(val) <= 2 {
		return errors.New("not enough values given to add user")
	}
	if _, err := um.Collection.InsertOne(context.TODO(), val); err != nil {
		return err
	}
	return nil

}
func (um *UserModel) RemoveEntry(val bson.D) error {
	if len(val) == 0 {
		return errors.New("empty val value given")
	}
	if _, err := um.Collection.DeleteOne(context.TODO(), val); err != nil {
		return err
	}
	return nil
}
func (um *UserModel) ModifyEntry(filter bson.D, val bson.D) error {
	if len(filter) == 0 {
		return errors.New("empty filter value given")
	}
	if len(val) == 0 {
		return errors.New("no empty update value given")
	}
	if _, err := um.Collection.UpdateOne(context.TODO(), filter, val); err != nil {
		return err
	}
	return nil
}

func NewUserModel(client *mongo.Database) *UserModel {
	c := client.Collection(userCollectionName)
	return &UserModel{
		Collection: c,
	}
}
