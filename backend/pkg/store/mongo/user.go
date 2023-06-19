package mongo

import (
	"context"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/user"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// UserRepo Mongo DB implementation for domain user entity
type UserRepo struct {
	collection    *mongo.Collection
	langRepo      query.LangViewRepository
	roleConverter *query.RoleConverter
}

// UserModel represents mongo user document
type UserModel struct {
	ID            string `bson:"_id"`
	Name          string `bson:"name"`
	Email         string `bson:"email"`
	Password      string `bson:"password"`
	Role          int    `bson:"role"`
	DefaultLangID string `bson:"default_lang_id"`
}

// NewUserRepo creates new UserRepo
func NewUserRepo(db *mongo.Database, langRepo query.LangViewRepository, roleMapper *query.RoleConverter) (*UserRepo, error) {
	u := UserRepo{collection: db.Collection("users"), langRepo: langRepo, roleConverter: roleMapper}

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
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	if _, err := r.collection.Indexes().CreateMany(ctx, indexes); err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) Create(usr *user.User) error {
	model, err := r.fromDomainToModel(usr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	if _, err = r.collection.InsertOne(ctx, model); err != nil {
		return replaceOnDuplicateKeyError(err, user.ErrEmailAlreadyExists)
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

	return r.fromModelToDomain(record), nil
}

func (r *UserRepo) Get(id string) (*user.User, error) {
	var record UserModel

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&record)

	if err == mongo.ErrNoDocuments {
		return nil, user.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return r.fromModelToDomain(record), nil
}

func (r *UserRepo) Update(usr *user.User) error {
	model, err := r.fromDomainToModel(usr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	result, err := r.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: model.ID}}, bson.M{"$set": model})

	if err != nil {
		return replaceOnDuplicateKeyError(err, user.ErrEmailAlreadyExists)
	}

	if result.MatchedCount != 1 {
		return fmt.Errorf("usr with id %s which must be modified not found", model.ID)
	}

	return nil
}

func (r *UserRepo) GetAllViews() ([]query.UserView, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	var models []UserModel
	if err = cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	views := make([]query.UserView, 0, len(models))
	for _, model := range models {
		view, err2 := r.fromModelToView(model)
		if err2 != nil {
			return nil, err2
		}
		views = append(views, view)
	}

	return views, nil
}

func (r *UserRepo) GetView(id string) (query.UserView, error) {
	var record UserModel

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&record)
	if err != nil {
		return query.UserView{}, err
	}

	return r.fromModelToView(record)
}

// fromDomainToModel converts domain user to mongo model
func (r *UserRepo) fromDomainToModel(usr *user.User) (UserModel, error) {
	model := UserModel{}
	err := mapstructure.Decode(usr.ToMap(), &model)
	return model, err
}

// fromModelToView converts mongo model to user View
func (r *UserRepo) fromModelToView(model UserModel) (query.UserView, error) {
	role, err := r.roleConverter.RoleToView(user.Role(model.Role))
	if err != nil {
		return query.UserView{}, err
	}

	view := query.UserView{
		ID:          model.ID,
		Name:        model.Name,
		Email:       model.Email,
		Role:        role,
		DefaultLang: query.LangView{},
	}

	if model.DefaultLangID != "" {
		langView, err := r.langRepo.GetView(model.DefaultLangID, model.ID)

		if err != nil {
			return query.UserView{}, err
		}
		view.DefaultLang = langView
	}

	return view, nil
}

// fromModelToDomain converts mongo model to user entity
func (r *UserRepo) fromModelToDomain(model UserModel) *user.User {
	return user.UnmarshalFromDB(
		model.ID,
		model.Name,
		model.Email,
		model.Password,
		model.Role,
		model.DefaultLangID,
	)
}
