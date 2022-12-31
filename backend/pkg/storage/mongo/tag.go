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
	ID       string `bson:"_id"`
	Tag      string `bson:"tag"`
	AuthorID string `bson:"author_id"`
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
func (r *TagRepo) initIndexes() error {
	return nil
}

// Create saves new tag to DB
func (r *TagRepo) Create(t *tag.Tag) error {
	model, err := r.fromDomainToModel(t)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(r.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	if _, err = r.collection.InsertOne(ctx, model); err != nil {
		return err
	}

	return nil
}

// Update updates already existed tag
func (r *TagRepo) Update(t *tag.Tag) error {
	model, err := r.fromDomainToModel(t)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(r.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	result, err := r.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: model.ID}}, model)

	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return fmt.Errorf("tag with id %s which must be modified not found", model.ID)
	}

	return nil
}

// Get searches for tag with id and authorId
func (r *TagRepo) Get(id, authorID string) (*tag.Tag, error) {
	var record TagModel

	ctx, cancel := context.WithTimeout(r.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}, {Key: "author_id", Value: authorID}}).Decode(&record)

	if err != mongo.ErrNoDocuments {
		return nil, tag.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return tag.UnmarshalFromDB(
		record.ID,
		record.Tag,
		record.AuthorID,
	), nil
}

// Delete removes tag with id and authorId
func (r *TagRepo) Delete(id, authorID string) error {
	ctx, cancel := context.WithTimeout(r.ctx, queryDefaultTimeoutInSec*time.Second)
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

// AllExist checks that all tags exist in DB with passed ids and authorId
func (r *TagRepo) AllExist(ids []string, authorID string) (bool, error) {
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}, {Key: "author_id", Value: authorID}}

	ctx, cancel := context.WithTimeout(r.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return int(count) == len(ids), nil
}

// GetAllViews returns all existing tag views for passed authorId
func (r *TagRepo) GetAllViews(authorID string) ([]query.TagView, error) {
	filter := bson.D{{Key: "author_id", Value: authorID}}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var models []TagModel
	if err = cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	views := make([]query.TagView, 0, len(models))
	for _, model := range models {
		views = append(views, r.fromModelToView(model))
	}

	return views, nil
}

// GetView searches for tag with id and authorId
func (r *TagRepo) GetView(id, authorID string) (query.TagView, error) {
	var record TagModel

	ctx, cancel := context.WithTimeout(r.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}, {Key: "author_id", Value: authorID}}).Decode(&record)
	if err != nil {
		return query.TagView{}, err
	}

	return r.fromModelToView(record), nil
}

// GetViews returns tag views for passed ids and authorId
func (r *TagRepo) GetViews(ids []string, authorID string) ([]query.TagView, error) {
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}, {Key: "author_id", Value: authorID}}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var models []TagModel
	if err = cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	views := make([]query.TagView, 0, len(models))
	for _, model := range models {
		views = append(views, r.fromModelToView(model))
	}

	return views, nil
}

// fromDomainToModel converts domain tag to mongo model
func (r *TagRepo) fromDomainToModel(t *tag.Tag) (TagModel, error) {
	model := TagModel{}
	err := mapstructure.Decode(t.ToMap(), &model)
	return model, err
}

// fromModelToView converts mongo model to tag View
func (r *TagRepo) fromModelToView(model TagModel) query.TagView {
	return query.TagView{
		ID:  model.ID,
		Tag: model.Tag,
	}
}
