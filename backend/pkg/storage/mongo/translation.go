package mongo

import (
	"context"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/macyan13/webdict/backend/pkg/domain/translation"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// TranslationRepo Mongo DB implementation for domain translation entity
type TranslationRepo struct {
	ctx        context.Context
	collection *mongo.Collection
	tagRepo    query.TagViewRepository
}

// TranslationModel represents mongo translation document
type TranslationModel struct {
	Id            string    `bson:"_id"`
	AuthorId      string    `bson:"author_id"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updatedAt"`
	Transcription string    `bson:"transcription,omitempty"`
	Translation   string    `bson:"translation"`
	Text          string    `bson:"text"`
	Example       string    `bson:"example,omitempty"`
	TagIds        []string  `bson:"tag_ids,omitempty"`
}

// NewTranslationRepo creates new TranslationRepo
func NewTranslationRepo(ctx context.Context, db *mongo.Database, tagRepo query.TagViewRepository) (*TranslationRepo, error) {
	t := TranslationRepo{ctx: ctx, collection: db.Collection("translations"), tagRepo: tagRepo}

	if err := t.initIndexes(); err != nil {
		return nil, err
	}
	return &t, nil
}

// initIndexes creates required for current queries indexes in translation collection
func (t *TranslationRepo) initIndexes() error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{"author_id", 1},
				{"created_at", -1},
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	if _, err := t.collection.Indexes().CreateMany(ctx, indexes); err != nil {
		return err
	}
	return nil
}

// Create saves new translation to DB
func (t *TranslationRepo) Create(translation translation.Translation) error {
	model, err := t.fromDomainToModel(translation)
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

// Update updates already existed translation
func (t *TranslationRepo) Update(translation translation.Translation) error {
	model, err := t.fromDomainToModel(translation)
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
		return fmt.Errorf("translation with id %s which must be modified not found", model.Id)
	}

	return nil
}

// Get performs search request based on translation id and author id parameters and returns domain translation entity
func (t *TranslationRepo) Get(id, authorId string) (translation.Translation, error) {
	var record TranslationModel

	ctx, cancel := context.WithTimeout(t.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := t.collection.FindOne(ctx, bson.D{{"_id", id}, {"author_id", authorId}}).Decode(&record)

	if err == mongo.ErrNoDocuments {
		return translation.Translation{}, translation.NotFoundErr
	}

	if err != nil {
		return translation.Translation{}, err
	}

	return translation.UnmarshalFromDB(
		record.Id,
		record.AuthorId,
		record.CreatedAt,
		record.UpdatedAt,
		record.Transcription,
		record.Translation,
		record.Text,
		record.Example,
		record.TagIds,
	), nil
}

// Delete removes translation record by passed id and authorId fields
func (t *TranslationRepo) Delete(id, authorId string) error {
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

// GetView perform search request based on translation id and author id parameters and returns translation view representation
func (t *TranslationRepo) GetView(id, authorId string) (query.TranslationView, error) {
	var record TranslationModel

	ctx, cancel := context.WithTimeout(t.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := t.collection.FindOne(ctx, bson.D{{"_id", id}, {"author_id", authorId}}).Decode(&record)
	if err != nil {
		return query.TranslationView{}, err
	}

	return t.fromModelToView(record)
}

// GetLastViews provide a limited slice of views ordered in DESC order by created_at field
func (t *TranslationRepo) GetLastViews(authorId string, limit int) ([]query.TranslationView, error) {
	filter := bson.D{{"author_id", authorId}}
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	opts.SetLimit(int64(limit))

	ctx, cancel := context.WithTimeout(t.ctx, queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	cursor, err := t.collection.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	var models []TranslationModel
	if err = cursor.All(context.TODO(), &models); err != nil {
		return nil, err
	}

	views := make([]query.TranslationView, 0, limit)
	for _, model := range models {
		view, err := t.fromModelToView(model)

		if err != nil {
			return nil, err
		}

		views = append(views, view)
	}

	return views, nil
}

// fromDomainToModel converts domain translation to mongo model
func (t *TranslationRepo) fromDomainToModel(translation translation.Translation) (TranslationModel, error) {
	model := TranslationModel{}
	err := mapstructure.Decode(translation.ToMap(), &model)
	return model, err
}

// fromModelToView converts mongo model to translation View performing request for receiving related tag views
func (t *TranslationRepo) fromModelToView(model TranslationModel) (query.TranslationView, error) {
	tagViews, err := t.tagRepo.GetViews(model.TagIds, model.AuthorId)

	if err != nil {
		return query.TranslationView{}, err
	}

	if len(model.TagIds) != len(tagViews) {
		return query.TranslationView{}, fmt.Errorf("can not find all translation tags")
	}

	return query.TranslationView{
		Id:            model.Id,
		CreatedAd:     model.CreatedAt,
		Transcription: model.Transcription,
		Translation:   model.Translation,
		Text:          model.Text,
		Example:       model.Example,
		Tags:          tagViews,
	}, nil
}
