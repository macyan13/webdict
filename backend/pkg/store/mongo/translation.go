package mongo

import (
	"context"
	"fmt"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"regexp"
	"time"
)

// TranslationRepo Mongo DB implementation for domain translation entity
type TranslationRepo struct {
	collection *mongo.Collection
	tagRepo    query.TagViewRepository
	langRepo   query.LangViewRepository
}

// TranslationModel represents mongo translation document
type TranslationModel struct {
	ID            string    `bson:"_id"`
	AuthorID      string    `bson:"author_id"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updatedAt"`
	Transcription string    `bson:"transcription"`
	Target        string    `bson:"target"`
	Source        string    `bson:"source"`
	Example       string    `bson:"example"`
	TagIDs        []string  `bson:"tag_ids"`
	LangID        string    `bson:"lang_id"`
}

// NewTranslationRepo creates new TranslationRepo
func NewTranslationRepo(db *mongo.Database, tagRepo query.TagViewRepository, langRepo query.LangViewRepository) (*TranslationRepo, error) {
	t := TranslationRepo{collection: db.Collection("translations"), tagRepo: tagRepo, langRepo: langRepo}

	if err := t.initIndexes(); err != nil {
		return nil, err
	}
	return &t, nil
}

// initIndexes creates required for current queries indexes in translation collection
func (r *TranslationRepo) initIndexes() error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "source", Value: 1}, {Key: "lang_id", Value: 1}, {Key: "author_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "tag_ids", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "lang_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "lang_id", Value: 1},
				{Key: "tag_ids", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "lang_id", Value: 1},
				{Key: "source", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "author_id", Value: 1},
				{Key: "lang_id", Value: 1},
				{Key: "target", Value: 1},
				{Key: "created_at", Value: -1},
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
		return replaceOnDuplicateKeyError(err, translation.ErrSourceAlreadyExists)
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

	result, err := r.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: model.ID}}, bson.M{"$set": model})

	if err != nil {
		return replaceOnDuplicateKeyError(err, translation.ErrSourceAlreadyExists)
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

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, translation.ErrNotFound
		}
		return nil, err
	}

	return translation.UnmarshalFromDB(
		record.ID,
		record.Source,
		record.Transcription,
		record.Target,
		record.AuthorID,
		record.Example,
		record.TagIDs,
		record.CreatedAt,
		record.UpdatedAt,
		record.LangID,
	), nil
}

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

func (r *TranslationRepo) ExistByLang(langID, authorID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.D{{Key: "lang_id", Value: langID}, {Key: "author_id", Value: authorID}})

	return count > 0, err
}

func (r *TranslationRepo) ExistByTag(tagID, authorID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.D{{Key: "tag_ids", Value: tagID}, {Key: "author_id", Value: authorID}})

	return count > 0, err
}

func (r *TranslationRepo) DeleteByAuthorID(authorID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()
	result, err := r.collection.DeleteMany(ctx, bson.D{{Key: "author_id", Value: authorID}})
	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}

func (r *TranslationRepo) GetLastViewsByTags(authorID, langID string, pageSize, page int, tagIds []string) (query.LastTranslationViews, error) {
	filter := bson.D{{Key: "author_id", Value: authorID}, {Key: "lang_id", Value: langID}}
	if len(tagIds) != 0 {
		filter = append(filter, bson.E{Key: "tag_ids", Value: bson.D{{Key: "$all", Value: tagIds}}})
	}

	return r.getLastViewsByFilter(filter, pageSize, page)
}

func (r *TranslationRepo) GetLastViewsBySourcePart(authorID, langID, sourcePart string, pageSize, page int) (query.LastTranslationViews, error) {
	patter := primitive.Regex{Pattern: fmt.Sprintf(".*%s.*", regexp.QuoteMeta(sourcePart))}
	filter := bson.D{{Key: "author_id", Value: authorID}, {Key: "lang_id", Value: langID}, {Key: "source", Value: bson.M{"$regex": patter}}}

	return r.getLastViewsByFilter(filter, pageSize, page)
}

func (r *TranslationRepo) GetLastViewsByTargetPart(authorID, langID, targetPart string, pageSize, page int) (query.LastTranslationViews, error) {
	patter := primitive.Regex{Pattern: fmt.Sprintf(".*%s.*", regexp.QuoteMeta(targetPart))}
	filter := bson.D{{Key: "author_id", Value: authorID}, {Key: "lang_id", Value: langID}, {Key: "target", Value: bson.M{"$regex": patter}}}

	return r.getLastViewsByFilter(filter, pageSize, page)
}

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

func (r *TranslationRepo) GetRandomViews(authorID, langID string, tagIds []string, limit int) (query.RandomViews, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	filter := bson.D{{Key: "author_id", Value: authorID}, {Key: "lang_id", Value: langID}}
	if len(tagIds) != 0 {
		filter = append(filter, bson.E{Key: "tag_ids", Value: bson.D{{Key: "$all", Value: tagIds}}})
	}

	pipeline := []bson.D{
		{{Key: "$match", Value: filter}},
		{{Key: "$sample", Value: bson.M{"size": limit}}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return query.RandomViews{}, err
	}

	var models []TranslationModel

	if err = cursor.All(ctx, &models); err != nil {
		return query.RandomViews{}, err
	}

	err = cursor.Close(ctx)
	if err != nil {
		return query.RandomViews{}, err
	}

	views := make([]query.TranslationView, 0, limit)

	for i := range models {
		view, err := r.fromModelToView(models[i])

		if err != nil {
			return query.RandomViews{}, err
		}

		views = append(views, view)
	}

	return query.RandomViews{Views: views}, nil
}

func (r *TranslationRepo) getLastViewsByFilter(filter bson.D, pageSize, page int) (query.LastTranslationViews, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), queryDefaultTimeoutInSec*time.Second)
	defer cancel()

	totalDocuments, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return query.LastTranslationViews{}, err
	}

	if totalDocuments == 0 && page == 1 {
		return query.LastTranslationViews{}, nil
	}

	totalPages := int(totalDocuments) / pageSize
	if int(totalDocuments)%pageSize != 0 {
		totalPages++
	}

	if totalPages < page {
		return query.LastTranslationViews{}, fmt.Errorf("can not get %d translations page from DB as max page is %d", page, totalPages)
	}

	skip := (page - 1) * pageSize
	sort := bson.M{"created_at": -1}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(sort))
	if err != nil {
		return query.LastTranslationViews{}, err
	}

	var models []TranslationModel

	if err = cursor.All(ctx, &models); err != nil {
		return query.LastTranslationViews{}, err
	}

	err = cursor.Close(ctx)
	if err != nil {
		return query.LastTranslationViews{}, err
	}

	views := make([]query.TranslationView, 0, pageSize)

	for i := range models {
		view, err := r.fromModelToView(models[i])

		if err != nil {
			return query.LastTranslationViews{}, err
		}

		views = append(views, view)
	}

	return query.LastTranslationViews{
		Views:        views,
		TotalRecords: int(totalDocuments),
	}, nil
}

// fromDomainToModel converts domain translation to mongo model
func (r *TranslationRepo) fromDomainToModel(t *translation.Translation) (TranslationModel, error) {
	model := TranslationModel{}
	err := mapstructure.Decode(t.ToMap(), &model)
	return model, err
}

// fromModelToView converts mongo model to translation View performing request for receiving related tag views
func (r *TranslationRepo) fromModelToView(model TranslationModel) (query.TranslationView, error) {
	view := query.TranslationView{
		ID:            model.ID,
		CreatedAd:     model.CreatedAt,
		Transcription: model.Transcription,
		Target:        model.Target,
		Source:        model.Source,
		Example:       model.Example,
	}

	langView, err := r.langRepo.GetView(model.LangID, model.AuthorID)
	if err != nil {
		return query.TranslationView{}, err
	}

	view.Lang = langView

	if model.TagIDs == nil {
		return view, nil
	}
	tagViews, err := r.tagRepo.GetViews(model.TagIDs, model.AuthorID)

	if err != nil {
		return query.TranslationView{}, err
	}

	if len(model.TagIDs) != len(tagViews) {
		return query.TranslationView{}, fmt.Errorf("can not find all translation tags")
	}

	view.Tags = tagViews

	return view, nil
}
