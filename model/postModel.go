package model

//types here have to implement the  Modeler interface

import (
	"context"
	"errors"
	"social-api/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const postCollectionName string = "posts"

//types here have to implement the  Modeler interface

type PostModel struct {
	Collection *mongo.Collection
}

// simple search when you need to get a entry without any filter options
// will only return single entry
func (pm *PostModel) GetEntry(key bson.D) (*types.Posts, error) {
	var entry types.Posts
	err := pm.Collection.FindOne(context.TODO(), key).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (pm *PostModel) GetEntryAdvanced(filter bson.D, sort bson.D) ([]*types.Posts, error) {
	opts := options.Find().SetSort(sort)
	cur, err := pm.Collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	var entrys []*types.Posts
	if err = cur.All(context.TODO(), &entrys); err != nil {
		return nil, err
	}
	// gonna return a error if no data return for the given filters
	if len(entrys) == 0 {
		return nil, errors.New("no values found")
	}
	return entrys, nil

}

func (pm *PostModel) AddEntry(val bson.D) error {
	if len(val) < 3 {
		return errors.New("not enough values given to add post")
	}
	if _, err := pm.Collection.InsertOne(context.TODO(), val); err != nil {
		return err
	}
	return nil

}
func (pm *PostModel) RemoveEntry(val bson.D) error {
	if _, err := pm.Collection.DeleteOne(context.TODO(), val); err != nil {
		return err
	}
	return nil
}
func (pm *PostModel) ModifyEntry(filter bson.D, val bson.D) error {
	if _, err := pm.Collection.UpdateOne(context.TODO(), filter, val); err != nil {
		return err
	}
	return nil
}

func NewPostModel(client *mongo.Database) *PostModel {
	c := client.Collection(postCollectionName)
	return &PostModel{
		Collection: c,
	}
}
