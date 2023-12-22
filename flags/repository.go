package flags

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var ctx = context.Background()

type FlagRepository interface {
    Get(key string) (bool, error)
    Keys() ([]string, error)
    Exists(name string) (bool, error)
    Set(key string, value interface{}) error
}

type FlagMongoRepository struct {
    collection *mongo.Collection
}

func NewFlagMongoRepository(host string, port int) (FlagRepository, error) {
    hostStr := fmt.Sprintf("mongodb://%v:%v",host,port)
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(hostStr))
    if err != nil {
        return nil,err
    }
    collection := client.Database("toggles").Collection("flags")
    return FlagMongoRepository{collection},nil
}

func(r FlagMongoRepository) Keys() ([]string, error) {
    return nil,nil
}

func(r FlagMongoRepository) Get(key string) (bool, error) {
    x := r.collection.FindOne(ctx,bson.D{{Key:"name",Value:key}})
    var f Flag
    err := x.Decode(&f)
    if err == mongo.ErrNoDocuments {
        return false,FlagNotFoundError
    } else if err != nil {
        return false,err
    }
    return f.Value,nil
}

func (r FlagMongoRepository) Set(key string, value interface{}) error {
    filter := bson.D{{"name",key}}
    update := bson.D{{"$set", bson.D{{"value",value}}}}
    opts := options.Update().SetUpsert(true)
    _, err := r.collection.UpdateOne(context.TODO(), filter, update, opts)

    if err != nil {
        return err
    } else {
        return nil
    }
}

func (r FlagMongoRepository) Exists(name string) (bool,error) {
    _,err := r.Get(name)
    log.Printf("aaaa ver: %v",err)
    if err == nil {
        return true,nil
    } else if err == FlagNotFoundError {
        return false,nil
    } else {
        return false,err
    }
}
