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

	views, err := t.queryProxy.GetAllViews(authorID)
	if err != nil {
		return nil, err
	}

	cacheMap := make(map[string]query.TagView, len(views))

	for i := range views {
		cacheMap[views[i].ID] = views[i]
	}

	t.cache.Set(authorID, cacheMap, cache.WithExpiration(t.cacheTTL))
	return views, nil
}

func (t TagRepo) GetView(id, authorID string) (query.TagView, error) {
	cachedViews, ok := t.cache.Get(authorID)

	if !ok {
		cachedViews = map[string]query.TagView{}
	}

	if err := t.updateCacheWithMisses([]string{id}, authorID, cachedViews); err != nil {
		return query.TagView{}, err
	}

	view, ok := cachedViews[id]
	if ok {
		t.cache.Set(authorID, cachedViews, cache.WithExpiration(t.cacheTTL))
		return view, nil
	}

	return query.TagView{}, fmt.Errorf("can not find tag, userID: %s, tagID: %s", id, authorID)
}

func (t TagRepo) GetViews(ids []string, authorID string) ([]query.TagView, error) {
	cachedViews, ok := t.cache.Get(authorID)
	if !ok {
		cachedViews = map[string]query.TagView{}
	}

	if err := t.updateCacheWithMisses(ids, authorID, cachedViews); err != nil {
		return nil, err
	}

	views := make([]query.TagView, 0, len(ids))

	for i := range ids {
		if view, hit := cachedViews[ids[i]]; hit {
			views = append(views, view)
		}
	}

	t.cache.Set(authorID, cachedViews, cache.WithExpiration(t.cacheTTL))
	return views, nil
}

func (t TagRepo) Create(tg *tag.Tag) error {
	err := t.domainProxy.Create(tg)

	if err == nil {
		t.cache.Delete(tg.AuthorID())
	}

	return err
}

func (t TagRepo) Update(tg *tag.Tag) error {
	if err := t.domainProxy.Update(tg); err != nil {
		return err
	}

	if cachedViews, ok := t.cache.Get(tg.AuthorID()); ok {
		delete(cachedViews, tg.ID())
	}

	return nil
}

func (t TagRepo) Get(id, authorID string) (*tag.Tag, error) {
	return t.domainProxy.Get(id, authorID)
}

func (t TagRepo) ExistByTag(tg, authorID string) (bool, error) {
	return t.domainProxy.ExistByTag(tg, authorID)
}

func (t TagRepo) Delete(id, authorID string) error {
	if err := t.domainProxy.Delete(id, authorID); err != nil {
		return err
	}

	if cachedViews, ok := t.cache.Get(authorID); ok {
		delete(cachedViews, id)
	}

	return nil
}

func (t TagRepo) AllExist(ids []string, authorID string) (bool, error) {
	return t.domainProxy.AllExist(ids, authorID)
}

func (t TagRepo) updateCacheWithMisses(ids []string, authorID string, cacheMap map[string]query.TagView) error {
	cacheMiss := make([]string, 0, len(ids))

	for i := range ids {
		if _, cacheHit := cacheMap[ids[i]]; !cacheHit {
			cacheMiss = append(cacheMiss, ids[i])
		}
	}

	if len(cacheMiss) == 0 {
		return nil
	}

	views, err := t.queryProxy.GetViews(cacheMiss, authorID)

	if err != nil {
		return err
	}

	for i := range views {
		cacheMap[views[i].ID] = views[i]
	}

	return nil
}
