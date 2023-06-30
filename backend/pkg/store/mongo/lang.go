package mongo

import (
	"context"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LangRepo struct {
	collection *mongo.Collection
}

type LangModel struct {
	ID       string `bson:"_id"`
	Name     string `bson:"name"`
	AuthorID string `bson:"author_id"`
}

func NewLangRepo(db *mongo.Database) (*LangRepo, error) {
	t := LangRepo{collection: db.Collection("langs")}

	if err := t.initIndexes(); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *LangRepo) initIndexes() error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
				{Key: "author_id", Value: 1},
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

func (r *LangRepo) Create(l *lang.Lang) error {
	model, err := r.fromDomainToModel(l)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	if _, err = r.collection.InsertOne(ctx, model); err != nil {
		return replaceOnDuplicateKeyError(err, lang.ErrLangAlreadyExists)
	}

	return nil
}

func (r *LangRepo) Update(l *lang.Lang) error {
	model, err := r.fromDomainToModel(l)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	result, err := r.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: model.ID}}, bson.M{"$set": model})

	if err != nil {
		return replaceOnDuplicateKeyError(err, lang.ErrLangAlreadyExists)
	}

	if result.MatchedCount != 1 {
		return fmt.Errorf("lang with id %s which must be modified not found", model.ID)
	}

	return nil
}

func (r *LangRepo) Get(id, authorID string) (*lang.Lang, error) {
	var record LangModel

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}, {Key: "author_id", Value: authorID}}).Decode(&record)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, lang.ErrNotFound
		}

		return nil, err
	}

	return lang.UnmarshalFromDB(
		record.ID,
		record.Name,
		record.AuthorID,
	), nil
}

func (r *LangRepo) Delete(id, authorID string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}, {Key: "author_id", Value: authorID}})

	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return fmt.Errorf("1 record was supposed to be deleted, %d removed", result.DeletedCount)
	}

	return nil
}

func (r *LangRepo) DeleteByAuthorID(authorID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()
	result, err := r.collection.DeleteMany(ctx, bson.D{{Key: "author_id", Value: authorID}})
	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}

func (r *LangRepo) Exist(id, authorID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.D{{Key: "_id", Value: id}, {Key: "author_id", Value: authorID}})

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *LangRepo) GetAllViews(authorID string) ([]query.LangView, error) {
	filter := bson.D{{Key: "author_id", Value: authorID}}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, filter, &options.FindOptions{Sort: bson.M{"created_at": -1}})

	if err != nil {
		return nil, err
	}

	var models []LangModel
	if err = cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	views := make([]query.LangView, 0, len(models))
	for _, model := range models {
		views = append(views, r.fromModelToView(model))
	}

	return views, nil
}

func (r *LangRepo) GetView(id, authorID string) (query.LangView, error) {
	var record LangModel

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}, {Key: "author_id", Value: authorID}}).Decode(&record)
	if err != nil {
		return query.LangView{}, err
	}

	return r.fromModelToView(record), nil
}

func (r *LangRepo) fromDomainToModel(l *lang.Lang) (LangModel, error) {
	model := LangModel{}
	err := mapstructure.Decode(l.ToMap(), &model)
	return model, err
}

func (r *LangRepo) fromModelToView(model LangModel) query.LangView {
	return query.LangView{
		ID:   model.ID,
		Name: model.Name,
	}
}
