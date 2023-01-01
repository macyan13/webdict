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
	collection *mongo.Collection
	tagRepo    query.TagViewRepository
}

// TranslationModel represents mongo translation document
type TranslationModel struct {
	ID            string    `bson:"_id"`
	AuthorID      string    `bson:"author_id"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updatedAt"`
	Transcription string    `bson:"transcription,omitempty"`
	Translation   string    `bson:"translation"`
	Text          string    `bson:"text"`
	Example       string    `bson:"example,omitempty"`
	TagIDs        []string  `bson:"tag_ids,omitempty"`
}

// NewTranslationRepo creates new TranslationRepo
func NewTranslationRepo(db *mongo.Database, tagRepo query.TagViewRepository) (*TranslationRepo, error) {
	t := TranslationRepo{collection: db.Collection("translations"), tagRepo: tagRepo}

	if err := t.initIndexes(); err != nil {
		return nil, err
	}
	return &t, nil
}

// initIndexes creates required for current queries indexes in translation collection
func (r *TranslationRepo) initIndexes() error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "text", Value: 1},
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

// Create saves new translation to DB
func (r *TranslationRepo) Create(t *translation.Translation) error {
	model, err := r.fromDomainToModel(t)
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

// Update updates already existed translation
func (r *TranslationRepo) Update(t *translation.Translation) error {
	model, err := r.fromDomainToModel(t)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	result, err := r.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: model.ID}}, model)

	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return fmt.Errorf("translation with id %s which must be modified not found", model.ID)
	}

	return nil
}

// Get performs search request based on translation id and author id parameters and returns domain translation entity
func (r *TranslationRepo) Get(id, authorID string) (*translation.Translation, error) {
	var record TranslationModel

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}, {Key: "author_id", Value: authorID}}).Decode(&record)

	if err == mongo.ErrNoDocuments {
		return nil, translation.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return translation.UnmarshalFromDB(
		record.ID,
		record.AuthorID,
		record.CreatedAt,
		record.UpdatedAt,
		record.Transcription,
		record.Translation,
		record.Text,
		record.Example,
		record.TagIDs,
	), nil
}

// Delete removes translation record by passed id and authorId fields
func (r *TranslationRepo) Delete(id, authorID string) error {
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

func (r *TranslationRepo) ExistByText(text, authorID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.D{{Key: "text", Value: text}, {Key: "author_id", Value: authorID}})

	return count > 0, err
}

// GetView perform search request based on translation id and author id parameters and returns translation view representation
func (r *TranslationRepo) GetView(id, authorID string) (query.TranslationView, error) {
	var record TranslationModel

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}, {Key: "author_id", Value: authorID}}).Decode(&record)
	if err != nil {
		return query.TranslationView{}, err
	}

	return r.fromModelToView(record)
}

// GetLastViews provide a limited slice of views ordered in DESC order by created_at field
func (r *TranslationRepo) GetLastViews(authorID string, limit int) ([]query.TranslationView, error) {
	filter := bson.D{{Key: "author_id", Value: authorID}}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	opts.SetLimit(int64(limit))

	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	var models []TranslationModel

	if err = cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	views := make([]query.TranslationView, 0, limit)

	for i := range models {
		view, err := r.fromModelToView(models[i])

		if err != nil {
			return nil, err
		}

		views = append(views, view)
	}

	return views, nil
}

// fromDomainToModel converts domain translation to mongo model
func (r *TranslationRepo) fromDomainToModel(t *translation.Translation) (TranslationModel, error) {
	model := TranslationModel{}
	err := mapstructure.Decode(t.ToMap(), &model)
	return model, err
}

// fromModelToView converts mongo model to translation View performing request for receiving related tag views
func (r *TranslationRepo) fromModelToView(model TranslationModel) (query.TranslationView, error) {
	tagViews, err := r.tagRepo.GetViews(model.TagIDs, model.AuthorID)

	if err != nil {
		return query.TranslationView{}, err
	}

	if len(model.TagIDs) != len(tagViews) {
		return query.TranslationView{}, fmt.Errorf("can not find all translation tags")
	}

	return query.TranslationView{
		ID:            model.ID,
		CreatedAd:     model.CreatedAt,
		Transcription: model.Transcription,
		Translation:   model.Translation,
		Text:          model.Text,
		Example:       model.Example,
		Tags:          tagViews,
	}, nil
}