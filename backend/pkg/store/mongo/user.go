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
	collection *mongo.Collection
}

// UserModel represents mongo user document
type UserModel struct {
	ID       string `bson:"_id"`
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
	Role     int    `bson:"role"`
}

// NewUserRepo creates new UserRepo
func NewUserRepo(db *mongo.Database) (*UserRepo, error) {
	u := UserRepo{collection: db.Collection("users")}

	if err := u.initIndexes(); err != nil {
		return nil, err
	}
	return &u, nil
}

// initIndexes creates required for current queries indexes in user collection
func (r *UserRepo) initIndexes() error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	if _, err := r.collection.Indexes().CreateMany(ctx, indexes); err != nil {
		return err
	}
	return nil
}

// Exist checks if user with the email exists
func (r *UserRepo) Exist(email string) (bool, error) {
	filter := bson.D{{Key: "email", Value: email}}

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (r *UserRepo) Create(usr *user.User) error {
	model, err := r.fromDomainToModel(usr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	if _, err = r.collection.InsertOne(ctx, model); err != nil {
		return err
	}

	return nil
}

// GetByEmail returns User for email
func (r *UserRepo) GetByEmail(email string) (*user.User, error) {
	var record UserModel

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&record)

	if err == mongo.ErrNoDocuments {
		return nil, user.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user.UnmarshalFromDB(
		record.ID,
		record.Name,
		record.Email,
		record.Password,
		record.Role,
	), nil
}

// fromDomainToModel converts domain user to mongo model
func (r *UserRepo) fromDomainToModel(usr *user.User) (UserModel, error) {
	model := UserModel{}
	err := mapstructure.Decode(usr.ToMap(), &model)
	return model, err
}
