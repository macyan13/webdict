package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/Code-Hex/go-generics-cache"
	"github.com/macyan13/webdict/backend/pkg/app/domain/translation"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
	"time"
)

func TestTranslationRepo_Create(t *testing.T) {
	type fields struct {
		domainProxy translation.Repository
		pageCache   *cache.Cache[string, map[string]query.LastTranslationViews]
	}
	type args struct {
		record *translation.Translation
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
		wantFn   assert.ValueAssertionFunc
	}{
		{
			"Error on create",
			func() fields {
				repo := translation.MockRepository{}
				repo.On("Create", mock.AnythingOfType("*translation.Translation")).Return(fmt.Errorf("error"))
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID", map[string]query.LastTranslationViews{"key": {}})
				return fields{
					domainProxy: &repo,
					pageCache:   pageCache,
				}
			},
			args{
				record: createTranslationByAuthorIDAndIDAndLangID("testID", "test", "DE"),
			},
			assert.Error,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				pageCache := i.(*cache.Cache[string, map[string]query.LastTranslationViews])
				entry, ok := pageCache.Get("authorID")
				assert.True(t, ok, i2...)
				_, ok = entry["key"]
				return assert.True(t, ok, i2...)
			},
		},
		{
			"New translation is created and cache is cleared",
			func() fields {
				repo := translation.MockRepository{}
				repo.On("Create", mock.AnythingOfType("*translation.Translation")).Return(nil)
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID-EN", map[string]query.LastTranslationViews{"key": {}})
				pageCache.Set("AuthorID-DE", map[string]query.LastTranslationViews{"key": {}})
				return fields{
					domainProxy: &repo,
					pageCache:   pageCache,
				}
			},
			args{
				record: createTranslationByAuthorIDAndIDAndLangID("authorID", "test", "EN"),
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err, i...)
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				pageCache := i.(*cache.Cache[string, map[string]query.LastTranslationViews])
				_, ok := pageCache.Get("AuthorID-DE")
				assert.True(t, ok, i2...)
				_, ok = pageCache.Get("authorID-EN")
				return assert.False(t, ok, i2...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := TranslationRepo{
				domainProxy:               f.domainProxy,
				cacheTTL:                  time.Minute,
				lastTranslationsPageCache: f.pageCache,
			}
			tt.wantErr(t, repo.Create(tt.args.record), fmt.Sprintf("Create(%v)", tt.args.record))
			tt.wantFn(t, repo.lastTranslationsPageCache, fmt.Sprintf("Create(%v)", tt.args.record))
		})
	}
}

func TestTranslationRepo_Update(t *testing.T) {
	type fields struct {
		domainProxy       translation.Repository
		singleRecordCache *cache.Cache[string, query.TranslationView]
		pageCache         *cache.Cache[string, map[string]query.LastTranslationViews]
	}
	type args struct {
		record *translation.Translation
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
		wantFn   assert.ValueAssertionFunc
	}{
		{
			"Error on DB update",
			func() fields {
				repo := translation.MockRepository{}
				repo.On("Update", mock.AnythingOfType("*translation.Translation")).Return(fmt.Errorf("error"))
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID", map[string]query.LastTranslationViews{"key": {}})
				singleCache := cache.NewContext[string, query.TranslationView](context.TODO())
				singleCache.Set("testID", query.TranslationView{})
				return fields{
					domainProxy:       &repo,
					pageCache:         pageCache,
					singleRecordCache: singleCache,
				}
			},
			args{record: createTranslationByAuthorIDAndIDAndLangID("authorID", "testID", "EN")},
			assert.Error,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				repo := i.(TranslationRepo)
				pageCache, ok := repo.lastTranslationsPageCache.Get("authorID")
				assert.True(t, ok, i2...)
				_, ok = pageCache["key"]
				assert.True(t, ok, i2...)
				_, ok = repo.singleRecordCache.Get("testID")
				return assert.True(t, ok, i2...)
			},
		},
		{
			"Translation is updated and cache is cleared",
			func() fields {
				repo := translation.MockRepository{}
				repo.On("Update", mock.AnythingOfType("*translation.Translation")).Return(nil)
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID-EN", map[string]query.LastTranslationViews{"key": {}})
				pageCache.Set("authorID-DE", map[string]query.LastTranslationViews{"key": {}})
				singleCache := cache.NewContext[string, query.TranslationView](context.TODO())
				singleCache.Set("testID", query.TranslationView{})
				singleCache.Set("otherID", query.TranslationView{})
				return fields{
					domainProxy:       &repo,
					pageCache:         pageCache,
					singleRecordCache: singleCache,
				}
			},
			args{record: createTranslationByAuthorIDAndIDAndLangID("authorID", "testID", "EN")},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err, i...)
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				repo := i.(TranslationRepo)
				_, ok := repo.lastTranslationsPageCache.Get("authorID-EN")
				assert.False(t, ok, i2...)
				_, ok = repo.lastTranslationsPageCache.Get("authorID-DE")
				assert.True(t, ok, i2...)
				_, ok = repo.singleRecordCache.Get("otherID")
				assert.True(t, ok, i2...)
				_, ok = repo.singleRecordCache.Get("testID")
				return assert.False(t, ok, i2...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := TranslationRepo{
				domainProxy:               f.domainProxy,
				cacheTTL:                  time.Minute,
				singleRecordCache:         f.singleRecordCache,
				lastTranslationsPageCache: f.pageCache,
			}
			tt.wantErr(t, repo.Update(tt.args.record), fmt.Sprintf("Update(%v)", tt.args.record))
			tt.wantFn(t, repo, fmt.Sprintf("Update(%v)", tt.args.record))
		})
	}
}

func TestTranslationRepo_Delete(t *testing.T) {
	type fields struct {
		domainProxy       translation.Repository
		singleRecordCache *cache.Cache[string, query.TranslationView]
		pageCache         *cache.Cache[string, map[string]query.LastTranslationViews]
	}
	type args struct {
		id       string
		authorID string
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  assert.ErrorAssertionFunc
		wantFn   assert.ValueAssertionFunc
	}{
		{
			"Error on DB DELETE",
			func() fields {
				repo := translation.MockRepository{}
				repo.On("Get", "testID", "authorID").Return(createTranslationByAuthorIDAndIDAndLangID("authorID", "testID", "EN"), nil)
				repo.On("Delete", "testID", "authorID").Return(fmt.Errorf("error"))
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID-EN", map[string]query.LastTranslationViews{"key": {}})
				singleCache := cache.NewContext[string, query.TranslationView](context.TODO())
				singleCache.Set("testID", query.TranslationView{})
				return fields{
					domainProxy:       &repo,
					pageCache:         pageCache,
					singleRecordCache: singleCache,
				}
			},
			args{authorID: "authorID", id: "testID"},
			assert.Error,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				repo := i.(TranslationRepo)
				pageCache, ok := repo.lastTranslationsPageCache.Get("authorID-EN")
				assert.True(t, ok, i2...)
				_, ok = pageCache["key"]
				assert.True(t, ok, i2...)
				_, ok = repo.singleRecordCache.Get("testID")
				return assert.True(t, ok, i2...)
			},
		},
		{
			"Error on getting record from DB",
			func() fields {
				repo := translation.MockRepository{}
				repo.On("Get", "testID", "authorID").Return(nil, fmt.Errorf("error"))
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID-EN", map[string]query.LastTranslationViews{"key": {}})
				singleCache := cache.NewContext[string, query.TranslationView](context.TODO())
				singleCache.Set("testID", query.TranslationView{})
				return fields{
					domainProxy:       &repo,
					pageCache:         pageCache,
					singleRecordCache: singleCache,
				}
			},
			args{authorID: "authorID", id: "testID"},
			assert.Error,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				repo := i.(TranslationRepo)
				pageCache, ok := repo.lastTranslationsPageCache.Get("authorID-EN")
				assert.True(t, ok, i2...)
				_, ok = pageCache["key"]
				assert.True(t, ok, i2...)
				_, ok = repo.singleRecordCache.Get("testID")
				return assert.True(t, ok, i2...)
			},
		},
		{
			"Translation is deleted and cache is cleared",
			func() fields {
				repo := translation.MockRepository{}
				repo.On("Delete", "testID", "authorID").Return(nil)
				repo.On("Get", "testID", "authorID").Return(createTranslationByAuthorIDAndIDAndLangID("authorID", "testID", "EN"), nil)
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID-EN", map[string]query.LastTranslationViews{"key": {}})
				pageCache.Set("authorID-DE", map[string]query.LastTranslationViews{"key": {}})
				singleCache := cache.NewContext[string, query.TranslationView](context.TODO())
				singleCache.Set("testID", query.TranslationView{})
				singleCache.Set("otherID", query.TranslationView{})
				return fields{
					domainProxy:       &repo,
					pageCache:         pageCache,
					singleRecordCache: singleCache,
				}
			},
			args{authorID: "authorID", id: "testID"},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err, i...)
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				repo := i.(TranslationRepo)
				_, ok := repo.lastTranslationsPageCache.Get("authorID-EN")
				assert.False(t, ok, i2...)
				_, ok = repo.lastTranslationsPageCache.Get("authorID-DE")
				assert.True(t, ok, i2...)
				_, ok = repo.singleRecordCache.Get("otherID")
				assert.True(t, ok, i2...)
				_, ok = repo.singleRecordCache.Get("testID")
				return assert.False(t, ok, i2...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			f := tt.fieldsFn()
			repo := TranslationRepo{
				domainProxy:               f.domainProxy,
				cacheTTL:                  time.Minute,
				singleRecordCache:         f.singleRecordCache,
				lastTranslationsPageCache: f.pageCache,
			}
			tt.wantErr(t1, repo.Delete(tt.args.id, tt.args.authorID), fmt.Sprintf("Delete(%v, %v)", tt.args.id, tt.args.authorID))
			tt.wantFn(t, repo, fmt.Sprintf("Delete(%v, %v)", tt.args.id, tt.args.authorID))
		})
	}
}

func TestTranslationRepo_sortTagsAlphabetically(t *testing.T) {
	type args struct {
		tags []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"Tags are sorted alphabetically",
			args{tags: []string{"a", "c", "b"}},
			[]string{"a", "b", "c"},
		},
		{
			"Tags are sorted alphabetically",
			args{tags: []string{"tagA", "tagC", "tagB", "tagD", "tagE", "tagF"}},
			[]string{"tagA", "tagB", "tagC", "tagD", "tagE", "tagF"},
		},
		{
			"Tags are sorted alphabetically",
			args{tags: []string{"a", "c", "b", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}},
			[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := TranslationRepo{}
			if got := repo.sortTagsAlphabetically(tt.args.tags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortTagsAlphabetically() = %v, want %v", got, tt.want)
			}
		})

	}
}

func TestTranslationRepo_GetView(t1 *testing.T) {
	type fields struct {
		queryProxy        query.TranslationViewRepository
		singleRecordCache *cache.Cache[string, query.TranslationView]
	}
	type args struct {
		id       string
		authorID string
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     query.TranslationView
		wantErr  assert.ErrorAssertionFunc
		wantFn   assert.ValueAssertionFunc
	}{
		{
			"Cache is not set, error on DB query",
			func() fields {
				repo := query.MockTranslationViewRepository{}
				repo.On("GetView", "testID", "authorID").Return(query.TranslationView{}, errors.New("error"))
				return fields{
					queryProxy:        &repo,
					singleRecordCache: cache.NewContext[string, query.TranslationView](context.TODO()),
				}
			},
			args{authorID: "authorID", id: "testID"},
			query.TranslationView{},
			assert.Error,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				singleCache := i.(*cache.Cache[string, query.TranslationView])
				_, ok := singleCache.Get("testID")
				return assert.False(t, ok, i2...)
			},
		},
		{
			"Cache is not set",
			func() fields {
				repo := query.MockTranslationViewRepository{}
				repo.On("GetView", "testID", "authorID").Return(query.TranslationView{ID: "testID"}, nil)
				return fields{
					queryProxy:        &repo,
					singleRecordCache: cache.NewContext[string, query.TranslationView](context.TODO()),
				}
			},
			args{authorID: "authorID", id: "testID"},
			query.TranslationView{ID: "testID"},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i...)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				singleCache := i.(*cache.Cache[string, query.TranslationView])
				_, ok := singleCache.Get("testID")
				return assert.True(t, ok, i2...)
			},
		},
		{
			"Cache is set",
			func() fields {
				singleRecordCache := cache.NewContext[string, query.TranslationView](context.TODO())
				singleRecordCache.Set("testID", query.TranslationView{ID: "testID"})
				return fields{
					singleRecordCache: singleRecordCache,
				}
			},
			args{authorID: "authorID", id: "testID"},
			query.TranslationView{ID: "testID"},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i...)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				singleCache := i.(*cache.Cache[string, query.TranslationView])
				_, ok := singleCache.Get("testID")
				return assert.True(t, ok, i2...)
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := TranslationRepo{
				queryProxy:        f.queryProxy,
				cacheTTL:          time.Minute,
				singleRecordCache: f.singleRecordCache,
			}
			got, err := repo.GetView(tt.args.id, tt.args.authorID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetView(%v, %v)", tt.args.id, tt.args.authorID)) {
				assert.Equal(t, tt.want, got)
			}
			tt.wantFn(t, repo.singleRecordCache, fmt.Sprintf("GetView(%v, %v)", tt.args.id, tt.args.authorID))
		})
	}
}

func TestTranslationRepo_GetLastViews(t *testing.T) {
	type fields struct {
		queryProxy query.TranslationViewRepository
		pageCache  *cache.Cache[string, map[string]query.LastTranslationViews]
	}
	type args struct {
		authorID string
		lang     string
		pageSize int
		page     int
		tagIds   []string
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		want     query.LastTranslationViews
		wantErr  assert.ErrorAssertionFunc
		wantFn   assert.ValueAssertionFunc
	}{
		{
			"Cache is not set, error on DB query",
			func() fields {
				repo := query.MockTranslationViewRepository{}
				repo.On("GetLastViewsByTags", "authorID", "EN", 10, 1, []string{"tag1", "tag2"}).Return(query.LastTranslationViews{}, errors.New("error"))
				return fields{
					queryProxy: &repo,
					pageCache:  cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO()),
				}
			},
			args{authorID: "authorID", lang: "EN", pageSize: 10, page: 1, tagIds: []string{"tag1", "tag2"}},
			query.LastTranslationViews{},
			assert.Error,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				pageCache := i.(*cache.Cache[string, map[string]query.LastTranslationViews])
				_, ok := pageCache.Get("authorID")
				return assert.False(t, ok, i2...)
			},
		},
		{
			"Author cache is set, requested page is missed",
			func() fields {
				repo := query.MockTranslationViewRepository{}
				repo.On("GetLastViewsByTags", "authorID", "EN", 10, 2, []string{"tag1", "tag2"}).Return(query.LastTranslationViews{TotalRecords: 5}, nil)
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID-EN", map[string]query.LastTranslationViews{"1": {}})
				return fields{
					queryProxy: &repo,
					pageCache:  pageCache,
				}
			},
			args{authorID: "authorID", lang: "EN", pageSize: 10, page: 2, tagIds: []string{"tag1", "tag2"}},
			query.LastTranslationViews{TotalRecords: 5},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i...)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				pageCache := i.(*cache.Cache[string, map[string]query.LastTranslationViews])
				authorCache, ok := pageCache.Get("authorID-EN")
				assert.True(t, ok, i2...)
				_, ok = authorCache["10-2-tag1-tag2"]
				return assert.True(t, ok, i2...)
			},
		},
		{
			"Cache is not set",
			func() fields {
				repo := query.MockTranslationViewRepository{}
				repo.On("GetLastViewsByTags", "authorID", "EN", 10, 2, []string{"tag1", "tag2"}).Return(query.LastTranslationViews{TotalRecords: 5}, nil)
				return fields{
					queryProxy: &repo,
					pageCache:  cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO()),
				}
			},
			args{authorID: "authorID", lang: "EN", pageSize: 10, page: 2, tagIds: []string{"tag1", "tag2"}},
			query.LastTranslationViews{TotalRecords: 5},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i...)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				pageCache := i.(*cache.Cache[string, map[string]query.LastTranslationViews])
				authorCache, ok := pageCache.Get("authorID-EN")
				assert.True(t, ok, i2...)
				_, ok = authorCache["10-2-tag1-tag2"]
				return assert.True(t, ok, i2...)
			},
		},
		{
			"Author cache is set, page cache is not set, error on getting requested page",
			func() fields {
				repo := query.MockTranslationViewRepository{}
				repo.On("GetLastViewsByTags", "authorID", "EN", 10, 2, []string{"tag1", "tag2"}).Return(query.LastTranslationViews{}, errors.New("error"))
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID-EN", map[string]query.LastTranslationViews{"1": {}})
				return fields{
					queryProxy: &repo,
					pageCache:  pageCache,
				}
			},
			args{authorID: "authorID", lang: "EN", pageSize: 10, page: 2, tagIds: []string{"tag1", "tag2"}},
			query.LastTranslationViews{},
			assert.Error,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				pageCache := i.(*cache.Cache[string, map[string]query.LastTranslationViews])
				authorCache, ok := pageCache.Get("authorID-EN")
				assert.True(t, ok, i2...)
				_, ok = authorCache["10-2-tag1-tag2"]
				return assert.False(t, ok, i2...)
			},
		},
		{
			"Page is set",
			func() fields {
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID-EN", map[string]query.LastTranslationViews{"10-2-tag1-tag2": {TotalRecords: 5}})
				return fields{
					pageCache: pageCache,
				}
			},
			args{authorID: "authorID", lang: "EN", pageSize: 10, page: 2, tagIds: []string{"tag2", "tag1"}},
			query.LastTranslationViews{TotalRecords: 5},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err, i...)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				pageCache := i.(*cache.Cache[string, map[string]query.LastTranslationViews])
				authorCache, ok := pageCache.Get("authorID-EN")
				assert.True(t, ok, i2...)
				_, ok = authorCache["10-2-tag1-tag2"]
				return assert.True(t, ok, i2...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := TranslationRepo{
				queryProxy:                f.queryProxy,
				lastTranslationsPageCache: f.pageCache,
				cacheTTL:                  time.Minute,
			}
			got, err := repo.GetLastViewsByTags(tt.args.authorID, tt.args.lang, tt.args.pageSize, tt.args.page, tt.args.tagIds)
			if !tt.wantErr(t, err, fmt.Sprintf("GetLastViewsByTags(%v, %v, %v, %v, %v)", tt.args.authorID, tt.args.lang, tt.args.pageSize, tt.args.page, tt.args.tagIds)) {
				assert.Equalf(t, tt.want, got, "GetLastViewsByTags(%v, %v, %v, %v, %v)", tt.args.authorID, tt.args.lang, tt.args.pageSize, tt.args.page, tt.args.tagIds)
			}
			tt.wantFn(t, repo.lastTranslationsPageCache, fmt.Sprintf("GetLastViewsByTags(%v, %v, %v, %v, %v)", tt.args.authorID, tt.args.lang, tt.args.pageSize, tt.args.page, tt.args.tagIds))
		})
	}
}

func createTranslationByAuthorIDAndIDAndLangID(authorID, id, langID string) *translation.Translation {
	return translation.UnmarshalFromDB(
		id,
		"test",
		"test",
		"test",
		authorID,
		"test",
		[]string{},
		time.Now(),
		time.Now(),
		langID,
	)
}

func TestTranslationRepo_authorLangCacheKey(t *testing.T) {
	authorID := "authorID"
	lang := "EN"
	repo := TranslationRepo{}
	assert.Equal(t, "authorID-EN", repo.authorLangCacheKey(authorID, lang))
}

func TestTranslationRepo_DeleteByAuthorID(t *testing.T) {
	type fields struct {
		domainProxy       translation.Repository
		singleRecordCache *cache.Cache[string, query.TranslationView]
		pageCache         *cache.Cache[string, map[string]query.LastTranslationViews]
	}
	type args struct {
		authorID string
	}
	tests := []struct {
		name          string
		fieldsFn      func() fields
		args          args
		want          int
		wantErr       assert.ErrorAssertionFunc
		assertCacheFn assert.ValueAssertionFunc
	}{
		{
			"Error on DB DELETE",
			func() fields {
				repo := translation.MockRepository{}
				repo.On("DeleteByAuthorID", "authorID").Return(0, fmt.Errorf("error"))
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID-EN", map[string]query.LastTranslationViews{"key": {}})
				singleCache := cache.NewContext[string, query.TranslationView](context.TODO())
				singleCache.Set("testID", query.TranslationView{})
				return fields{
					domainProxy:       &repo,
					pageCache:         pageCache,
					singleRecordCache: singleCache,
				}
			},
			args{authorID: "authorID"},
			0,
			assert.Error,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				repo := i.(TranslationRepo)
				pageCache, ok := repo.lastTranslationsPageCache.Get("authorID-EN")
				assert.True(t, ok, i2...)
				_, ok = pageCache["key"]
				assert.True(t, ok, i2...)
				_, ok = repo.singleRecordCache.Get("testID")
				return assert.True(t, ok, i2...)
			},
		},
		{
			"Cache is cleared",
			func() fields {
				repo := translation.MockRepository{}
				repo.On("DeleteByAuthorID", "authorID").Return(5, nil)
				pageCache := cache.NewContext[string, map[string]query.LastTranslationViews](context.TODO())
				pageCache.Set("authorID-EN", map[string]query.LastTranslationViews{"key": {}})
				singleCache := cache.NewContext[string, query.TranslationView](context.TODO())
				singleCache.Set("testID", query.TranslationView{})
				return fields{
					domainProxy:       &repo,
					pageCache:         pageCache,
					singleRecordCache: singleCache,
				}
			},
			args{authorID: "authorID"},
			5,
			assert.NoError,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				repo := i.(TranslationRepo)
				_, ok := repo.lastTranslationsPageCache.Get("authorID-EN")
				assert.False(t, ok, i2...)
				_, ok = repo.singleRecordCache.Get("testID")
				return assert.False(t, ok, i2...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			f := tt.fieldsFn()
			repo := TranslationRepo{
				domainProxy:               f.domainProxy,
				singleRecordCache:         f.singleRecordCache,
				lastTranslationsPageCache: f.pageCache,
			}
			got, err := repo.DeleteByAuthorID(tt.args.authorID)
			tt.assertCacheFn(t, repo, fmt.Sprintf("DeleteByAuthorID(%v)", tt.args.authorID))
			if !tt.wantErr(t1, err, fmt.Sprintf("DeleteByAuthorID(%v)", tt.args.authorID)) {
				return
			}
			assert.Equalf(t1, tt.want, got, "DeleteByAuthorID(%v)", tt.args.authorID)
		})
	}
}
