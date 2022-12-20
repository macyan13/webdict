package mongo

import (
	"context"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/domain/tag"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// TagRepo Mongo DB implementation for domain tag entity
type TagRepo struct {
	ctx        context.Context
	collection *mongo.Collection
}

// TagModel represents mongo tag document
type TagModel struct {
	Id       string `bson:"_id"`
	Tag      string `bson:"tag"`
	AuthorId string `bson:"author_id"`
}

// NewTagRepo creates TagRepo
func NewTagRepo(ctx context.Context, db *mongo.Database) (*TagRepo, error) {
	t := TagRepo{ctx: ctx, collection: db.Collection("tags")}

	if err := t.initIndexes(); err != nil {
		return nil, err
	}
	return &t, nil
}

// initIndexes creates required for current queries indexes in tags collection
func (t *TagRepo) initIndexes() error {
	return nil
}

// Create saves new tag to DB
func (t *TagRepo) Create(tag tag.Tag) error {
	model, err := t.fromDomainToModel(tag)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(t.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	if _, err = t.collection.InsertOne(ctx, model); err != nil {
		return err
	}

	return nil
}

// Update updates already existed tag
func (t *TagRepo) Update(tag tag.Tag) error {
	model, err := t.fromDomainToModel(tag)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(t.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	result, err := t.collection.UpdateOne(ctx, bson.D{{"_id", model.Id}}, model)

	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return fmt.Errorf("tag with id %s which must be modified not found", model.Id)
	}

	return nil
}

// Get searches for tag with id and authorId
func (t *TagRepo) Get(id, authorId string) (tag.Tag, error) {
	var record TagModel

	ctx, cancel := context.WithTimeout(t.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := t.collection.FindOne(ctx, bson.D{{"_id", id}, {"author_id", authorId}}).Decode(&record)

	if err != mongo.ErrNoDocuments {
		return tag.Tag{}, tag.NotFoundErr
	}

	if err != nil {
		return tag.Tag{}, err
	}

	return tag.UnmarshalFromDB(
		record.Id,
		record.Tag,
		record.AuthorId,
	), nil
}

// Delete removes tag with id and authorId
func (t *TagRepo) Delete(id, authorId string) error {
	ctx, cancel := context.WithTimeout(t.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	result, err := t.collection.DeleteOne(ctx, bson.D{{"_id", id}, {"author_id", authorId}})

	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return fmt.Errorf("1 record was supposed to be deleted, %d removed", result.DeletedCount)
	}

	return nil
}

// AllExist checks that all tags exist in DB with passed ids and authorId
func (t *TagRepo) AllExist(ids []string, authorId string) (bool, error) {
	filter := bson.D{{"_id", bson.D{{"$in", ids}}}, {"author_id", authorId}}

	ctx, cancel := context.WithTimeout(t.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	count, err := t.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return int(count) == len(ids), nil
}

// GetAllViews returns all existing tag views for passed authorId
func (t *TagRepo) GetAllViews(authorId string) ([]query.TagView, error) {
	filter := bson.D{{"author_id", authorId}}

	ctx, cancel := context.WithTimeout(t.ctx, 5*time.Second)
	defer cancel()

	cursor, err := t.collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var models []TagModel
	if err = cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	views := make([]query.TagView, 0, len(models))
	for _, model := range models {
		views = append(views, t.fromModelToView(model))
	}

	return views, nil
}

// GetView searches for tag with id and authorId
func (t *TagRepo) GetView(id, authorId string) (query.TagView, error) {
	var record TagModel

	ctx, cancel := context.WithTimeout(t.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := t.collection.FindOne(ctx, bson.D{{"_id", id}, {"author_id", authorId}}).Decode(&record)
	if err != nil {
		return query.TagView{}, err
	}

	return t.fromModelToView(record), nil
}

// GetViews returns tag views for passed ids and authorId
func (t *TagRepo) GetViews(ids []string, authorId string) ([]query.TagView, error) {
	filter := bson.D{{"_id", bson.D{{"$in", ids}}}, {"author_id", authorId}}

	ctx, cancel := context.WithTimeout(t.ctx, 5*time.Second)
	defer cancel()

	cursor, err := t.collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var models []TagModel
	if err = cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	views := make([]query.TagView, 0, len(models))
	for _, model := range models {
		views = append(views, t.fromModelToView(model))
	}

	return views, nil
}

// fromDomainToModel converts domain tag to mongo model
func (t *TagRepo) fromDomainToModel(tag tag.Tag) (TagModel, error) {
	model := TagModel{}
	err := mapstructure.Decode(tag.ToMap(), &model)
	return model, err
}

// fromModelToView converts mongo model to tag View
func (t *TagRepo) fromModelToView(model TagModel) query.TagView {
	return query.TagView{
		Id:  model.Id,
		Tag: model.Tag,
	}
}
