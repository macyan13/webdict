package cache

import (
	"context"
	"fmt"
	"github.com/Code-Hex/go-generics-cache"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"golang.org/x/exp/maps"
	"time"
)

type LangRepo struct {
	domainProxy lang.Repository
	queryProxy  query.LangViewRepository
	cache       *cache.Cache[string, map[string]query.LangView]
	cacheTTL    time.Duration
}

func NewLangRepo(ctx context.Context, domainProxy lang.Repository, queryProxy query.LangViewRepository, opts Opts) *LangRepo {
	return &LangRepo{
		domainProxy: domainProxy,
		queryProxy:  queryProxy,
		cache:       cache.NewContext[string, map[string]query.LangView](ctx),
		cacheTTL:    opts.LangCacheTTL,
	}
}

func (l LangRepo) Create(ln *lang.Lang) error {
	if err := l.domainProxy.Create(ln); err != nil {
		return err
	}

	l.cache.Delete(ln.AuthorID())
	return nil
}

func (l LangRepo) Update(ln *lang.Lang) error {
	if err := l.domainProxy.Update(ln); err != nil {
		return err
	}

	l.cache.Delete(ln.AuthorID())
	return nil
}

func (l LangRepo) Get(id, authorID string) (*lang.Lang, error) {
	return l.domainProxy.Get(id, authorID)
}

func (l LangRepo) Delete(id, authorID string) error {
	if err := l.domainProxy.Delete(id, authorID); err != nil {
		return err
	}

	l.cache.Delete(authorID)
	return nil
}

func (l LangRepo) Exist(id, authorID string) (bool, error) {
	return l.domainProxy.Exist(id, authorID)
}

func (l LangRepo) DeleteByAuthorID(authorID string) (int, error) {
	count, err := l.domainProxy.DeleteByAuthorID(authorID)
	if err == nil {
		l.cache.Delete(authorID)
	}

	return count, err
}

func (l LangRepo) GetAllViews(authorID string) ([]query.LangView, error) {
	if cachedViews, ok := l.cache.Get(authorID); ok {
		return maps.Values(cachedViews), nil
	}

	cachedViews, err := l.initCache(authorID)
	if err != nil {
		return nil, err
	}
	return maps.Values(cachedViews), nil
}

func (l LangRepo) GetView(id, authorID string) (query.LangView, error) {
	cachedViews, ok := l.cache.Get(authorID)
	if !ok {
		refreshedViews, err := l.initCache(authorID)
		if err != nil {
			return query.LangView{}, err
		}
		cachedViews = refreshedViews
	}

	if view, hit := cachedViews[id]; hit {
		return view, nil
	}
	return query.LangView{}, fmt.Errorf("can not find lang, userID: %s, langID: %s", authorID, id)
}

func (l LangRepo) initCache(authorID string) (map[string]query.LangView, error) {
	views, err := l.queryProxy.GetAllViews(authorID)
	if err != nil {
		return nil, err
	}

	cacheMap := make(map[string]query.LangView, len(views))

	for i := range views {
		cacheMap[views[i].ID] = views[i]
	}

	l.cache.Set(authorID, cacheMap, cache.WithExpiration(l.cacheTTL))
	return cacheMap, nil
}
