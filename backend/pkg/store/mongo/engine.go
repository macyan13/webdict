package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type Opts struct {
	Database string
	Host     string
	Port     int
	Username string
	Passwd   string
}

const queryDefaultTimeoutInSec = 3

func InitDatabase(ctx context.Context, opts Opts) (*mongo.Database, error) {
	clientOpts := options.Client()

	// Set the write concern
	wc := writeconcern.New(writeconcern.WMajority())
	clientOpts.ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", opts.Username, opts.Passwd, opts.Host, opts.Port)).SetWriteConcern(wc)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client.Database(opts.Database), nil
}

func replaceOnDuplicateKeyError(original, replace error) error {
	if mongo.IsDuplicateKeyError(original) {
		return replace
	}

	return original
}
