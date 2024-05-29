package user

import (
	"context"
	"fmt"

	"github.com/markelca/toggles/pkg/hash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

type UserRepository interface {
	FindAll() ([]*User, error)
	FindByUserName(userName string) (*User, error)
	Create(user User) error
	Update(user *User) error
	Upsert(user User) error
	Authenticate(userName, password string) (*User, error)
}

type UserMongoRepository struct {
	collection *mongo.Collection
}

func NewUserMongoRepository(host string, port uint) (UserRepository, error) {
	hostStr := fmt.Sprintf("mongodb://%v:%v", host, port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(hostStr))
	if err != nil {
		return nil, err
	}
	collection := client.Database("toggles").Collection("users")
	return UserMongoRepository{collection}, nil
}

func (repository UserMongoRepository) Upsert(user User) error {
	pwdHash, err := hash.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = pwdHash
	_, err = repository.collection.ReplaceOne(ctx, bson.D{{Key: "username", Value: user.UserName}}, user, options.Replace().SetUpsert(true))
	return err
}

func (repository UserMongoRepository) FindAll() ([]*User, error) {
	cur, err := repository.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	result := make([]*User, 0)
	for cur.Next(ctx) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		result = append(result, &user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (repository UserMongoRepository) FindByUserName(userName string) (*User, error) {
	x := repository.collection.FindOne(ctx, bson.D{{Key: "username", Value: userName}})
	var user User
	err := x.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repository UserMongoRepository) Create(user User) error {
	pwdHash, err := hash.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = pwdHash
	_, err = repository.collection.InsertOne(ctx, user)
	return err
}

func (repository UserMongoRepository) Update(user *User) error {
	_, err := repository.collection.ReplaceOne(ctx, bson.D{{Key: "username", Value: user.UserName}}, user)
	return err
}

func (repository UserMongoRepository) Authenticate(userName, password string) (*User, error) {
	user, err := repository.FindByUserName(userName)
	if err != nil {
		return nil, err
	}
	if !hash.CheckPasswordHash(password, user.Password) {
		fmt.Println("check password failed")
		return nil, ErrUserAuthenticationFailed
	}
	return user, nil
}