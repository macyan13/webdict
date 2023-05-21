package cache

import (
	"context"
	"fmt"
	"github.com/Code-Hex/go-generics-cache"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"sort"
	"strings"
	"time"
)

type TranslationRepo struct {
	domainProxy       translation.Repository
	queryProxy        query.TranslationViewRepository
	cacheTTL          time.Duration
	singleRecordCache *cache.Cache[string, query.TranslationView]
	pageCache         *cache.Cache[string, map[string]query.LastViews]
}

func NewTranslationRepo(ctx context.Context, domainProxy translation.Repository, queryProxy query.TranslationViewRepository, cacheTTL time.Duration) *TranslationRepo {
	return &TranslationRepo{
		domainProxy:       domainProxy,
		queryProxy:        queryProxy,
		cacheTTL:          cacheTTL,
		singleRecordCache: cache.NewContext[string, query.TranslationView](ctx),
		pageCache:         cache.NewContext[string, map[string]query.LastViews](ctx),
	}
}

func (t TranslationRepo) Create(record *translation.Translation) error {
	err := t.domainProxy.Create(record)
	if err == nil {
		t.pageCache.Delete(t.authorLangCacheKey(record.AuthorID(), record.LangID()))
	}

	return err
}

func (t TranslationRepo) Update(record *translation.Translation) error {
	err := t.domainProxy.Update(record)
	if err == nil {
		t.singleRecordCache.Delete(record.ID())
		t.pageCache.Delete(t.authorLangCacheKey(record.AuthorID(), record.LangID()))
	}

	return err
}

func (t TranslationRepo) Get(id, authorID string) (*translation.Translation, error) {
	return t.domainProxy.Get(id, authorID)
}

func (t TranslationRepo) ExistByTag(tagID, authorID string) (bool, error) {
	return t.domainProxy.ExistByTag(tagID, authorID)
}

func (t TranslationRepo) ExistByLang(langID, authorID string) (bool, error) {
	return t.domainProxy.ExistByLang(langID, authorID)
}

func (t TranslationRepo) Delete(id, authorID string) error {
	record, err := t.domainProxy.Get(id, authorID)

	if err != nil {
		return err
	}

	err = t.domainProxy.Delete(id, authorID)
	if err == nil {
		t.singleRecordCache.Delete(id)
		t.pageCache.Delete(t.authorLangCacheKey(record.AuthorID(), record.LangID()))
	}

	return err
}

func (t TranslationRepo) GetView(id, authorID string) (query.TranslationView, error) {
	if cachedView, ok := t.singleRecordCache.Get(id); ok {
		return cachedView, nil
	}

	view, err := t.queryProxy.GetView(id, authorID)
	if err == nil {
		t.singleRecordCache.Set(id, view, cache.WithExpiration(t.cacheTTL))
	}

	return view, err
}

func (t TranslationRepo) GetLastViews(authorID, langID string, pageSize, page int, tagIds []string) (query.LastViews, error) {
	pageKey := fmt.Sprintf("%d-%d-%v", pageSize, page, strings.Join(t.sortTagsAlphabetically(tagIds), "-"))
	authorPagesKey := t.authorLangCacheKey(authorID, langID)

	if authorLangPages, ok := t.pageCache.Get(authorPagesKey); ok {
		if cachedViews, ok := authorLangPages[pageKey]; ok {
			return cachedViews, nil
		}

		views, err := t.queryProxy.GetLastViews(authorID, langID, pageSize, page, tagIds)

		if err == nil {
			authorLangPages[pageKey] = views
		}

		return views, err
	}

	views, err := t.queryProxy.GetLastViews(authorID, langID, pageSize, page, tagIds)

	if err == nil {
		cacheMap := map[string]query.LastViews{pageKey: views}
		t.pageCache.Set(authorPagesKey, cacheMap, cache.WithExpiration(t.cacheTTL))
	}

	return views, err
}

func (t TranslationRepo) sortTagsAlphabetically(tagIds []string) []string {
	sort.Slice(tagIds, func(i, j int) bool {
		return tagIds[i] < tagIds[j]
	})

	return tagIds
}

func (t TranslationRepo) authorLangCacheKey(authorID, lang string) string {
	return fmt.Sprintf("%s-%s", authorID, lang)
}
