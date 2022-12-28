package mongo

import (
	"context"
	"github.com/macyan13/webdict/backend/pkg/domain/user"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// UserRepo Mongo DB implementation for domain user entity
type UserRepo struct {
	ctx        context.Context
	collection *mongo.Collection
}

// UserModel represents mongo user document
type UserModel struct {
	Id       string `bson:"_id"`
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
	Role     int    `bson:"role"`
}

// NewUserRepo creates new UserRepo
func NewUserRepo(ctx context.Context, db *mongo.Database) (*UserRepo, error) {
	u := UserRepo{ctx: ctx, collection: db.Collection("users")}

	if err := u.initIndexes(); err != nil {
		return nil, err
	}
	return &u, nil
}

// initIndexes creates required for current queries indexes in user collection
func (u *UserRepo) initIndexes() error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{"email", 1},
			},
		},
	}

	ctx, cancel := context.WithTimeout(u.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	if _, err := u.collection.Indexes().CreateMany(ctx, indexes); err != nil {
		return err
	}
	return nil
}

// Exist checks if user with the email exists
func (u *UserRepo) Exist(email string) (bool, error) {
	filter := bson.D{{Key: "email", Value: email}}

	ctx, cancel := context.WithTimeout(u.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	count, err := u.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (u *UserRepo) Create(user user.User) error {
	model, err := u.fromDomainToModel(user)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(u.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	if _, err = u.collection.InsertOne(ctx, model); err != nil {
		return err
	}

	return nil
}

// GetByEmail returns User for email
func (u *UserRepo) GetByEmail(email string) (user.User, error) {
	var record UserModel

	ctx, cancel := context.WithTimeout(u.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := u.collection.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&record)

	if err == mongo.ErrNoDocuments {
		return user.User{}, user.NotFoundErr
	}
	if err != nil {
		return user.User{}, err
	}

	return user.UnmarshalFromDB(
		record.Id,
		record.Name,
		record.Email,
		record.Password,
		record.Role,
	), nil
}

// fromDomainToModel converts domain user to mongo model
func (u *UserRepo) fromDomainToModel(user user.User) (UserModel, error) {
	model := UserModel{}
	err := mapstructure.Decode(user.ToMap(), &model)
	return model, err
}
