package cache

import (
	"context"
	"fmt"
	"github.com/Code-Hex/go-generics-cache"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"golang.org/x/exp/maps"
	"time"
)

type TagRepo struct {
	domainProxy tag.Repository
	queryProxy  query.TagViewRepository
	cache       *cache.Cache[string, map[string]query.TagView]
	cacheTTL    time.Duration
}

func NewTagRepo(ctx context.Context, domainProxy tag.Repository, queryProxy query.TagViewRepository, opts Opts) *TagRepo {
	return &TagRepo{
		domainProxy: domainProxy,
		queryProxy:  queryProxy,
		cache:       cache.NewContext[string, map[string]query.TagView](ctx),
		cacheTTL:    opts.TagCacheTTL,
	}
}

func (t TagRepo) GetAllViews(authorID string) ([]query.TagView, error) {
	if cachedViews, ok := t.cache.Get(authorID); ok {
		return maps.Values(cachedViews), nil
	}

	cachedViews, err := t.initCache(authorID)
	if err != nil {
		return nil, err
	}
	return maps.Values(cachedViews), nil
}

func (t TagRepo) GetView(id, authorID string) (query.TagView, error) {
	cachedViews, ok := t.cache.Get(authorID)
	if !ok {
		refreshedViews, err := t.initCache(authorID)
		if err != nil {
			return query.TagView{}, err
		}
		cachedViews = refreshedViews
	}

	if view, hit := cachedViews[id]; hit {
		return view, nil
	}
	return query.TagView{}, fmt.Errorf("can not find tag, userID: %s, tagID: %s", authorID, id)
}

func (t TagRepo) GetViews(ids []string, authorID string) ([]query.TagView, error) {
	cachedViews, ok := t.cache.Get(authorID)

	if !ok {
		refreshedViews, err := t.initCache(authorID)
		if err != nil {
			return nil, err
		}

		cachedViews = refreshedViews
	}

	views := make([]query.TagView, 0, len(ids))

	for i := range ids {
		if view, hit := cachedViews[ids[i]]; hit {
			views = append(views, view)
		} else {
			return nil, fmt.Errorf("can not find tag, userID: %s, tagID: %s", authorID, ids[i])
		}
	}

	return views, nil
}

func (t TagRepo) Create(tg *tag.Tag) error {
	if err := t.domainProxy.Create(tg); err != nil {
		return err
	}

	t.cache.Delete(tg.AuthorID())
	return nil
}

func (t TagRepo) Update(tg *tag.Tag) error {
	if err := t.domainProxy.Update(tg); err != nil {
		return err
	}

	t.cache.Delete(tg.AuthorID())
	return nil
}

func (t TagRepo) Get(id, authorID string) (*tag.Tag, error) {
	return t.domainProxy.Get(id, authorID)
}

func (t TagRepo) Delete(id, authorID string) error {
	if err := t.domainProxy.Delete(id, authorID); err != nil {
		return err
	}

	t.cache.Delete(authorID)
	return nil
}

func (t TagRepo) AllExist(ids []string, authorID string) (bool, error) {
	return t.domainProxy.AllExist(ids, authorID)
}

func (t TagRepo) initCache(authorID string) (map[string]query.TagView, error) {
	views, err := t.queryProxy.GetAllViews(authorID)
	if err != nil {
		return nil, err
	}

	cacheMap := make(map[string]query.TagView, len(views))

	for i := range views {
		cacheMap[views[i].ID] = views[i]
	}

	t.cache.Set(authorID, cacheMap, cache.WithExpiration(t.cacheTTL))
	return cacheMap, nil
}
