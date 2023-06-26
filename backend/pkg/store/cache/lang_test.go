package cache

import (
	"fmt"
	"github.com/Code-Hex/go-generics-cache"
	"github.com/macyan13/webdict/backend/pkg/app/domain/lang"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestLangRepo_Create(t *testing.T) {
	type fields struct {
		domainProxy lang.Repository
		cache       *cache.Cache[string, map[string]query.LangView]
	}
	type args struct {
		lang *lang.Lang
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		argsFn   func() args
		wantErr  assert.ErrorAssertionFunc
		wantFn   assert.ValueAssertionFunc
	}{
		{
			"Error on DB request",
			func() fields {
				domainProxy := lang.NewMockRepository(t)
				domainProxy.On("Create", mock.AnythingOfType("*lang.Lang")).Return(fmt.Errorf("testErr"))
				c := cache.New[string, map[string]query.LangView]()
				c.Set("testAuthor", map[string]query.LangView{"lang1": {ID: "lang1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			func() args {
				ln, err := lang.NewLang("en", "testAuthor")
				assert.Nil(t, err)
				return args{lang: ln}
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error())
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, _ := i2[0].(*cache.Cache[string, map[string]query.LangView]).Get("testAuthor")
				_, ok := cacheMap["lang1"]
				assert.True(t, ok, i2)
				return true
			},
		},
		{
			"Cache is set",
			func() fields {
				domainProxy := lang.NewMockRepository(t)
				domainProxy.On("Create", mock.AnythingOfType("*lang.Lang")).Return(nil)
				c := cache.New[string, map[string]query.LangView]()
				c.Set("testAuthor", map[string]query.LangView{"lang1": {ID: "lang1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			func() args {
				ln, err := lang.NewLang("en", "testAuthor")
				assert.Nil(t, err)
				return args{lang: ln}
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				_, ok := i2[0].(*cache.Cache[string, map[string]query.LangView]).Get("testAuthor")
				assert.False(t, ok, i2)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := LangRepo{
				domainProxy: f.domainProxy,
				cache:       f.cache,
				cacheTTL:    time.Minute,
			}
			if tt.wantErr(t, repo.Create(tt.argsFn().lang), fmt.Sprintf("Create(%v)", tt.argsFn().lang)) {
				return
			}
			tt.wantFn(t, tt.argsFn().lang, repo.cache, fmt.Sprintf("Create(%v)", tt.argsFn().lang))
		})
	}
}

func TestLangRepo_Update(t *testing.T) {
	type fields struct {
		domainProxy lang.Repository
		cache       *cache.Cache[string, map[string]query.LangView]
	}
	type args struct {
		lang *lang.Lang
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		argsFn   func() args
		wantErr  assert.ErrorAssertionFunc
		wantFn   assert.ValueAssertionFunc
	}{
		{
			"Error on DB request",
			func() fields {
				domainProxy := lang.NewMockRepository(t)
				domainProxy.On("Update", mock.AnythingOfType("*lang.Lang")).Return(fmt.Errorf("testErr"))
				c := cache.New[string, map[string]query.LangView]()
				c.Set("testAuthor", map[string]query.LangView{"lang1": {ID: "lang1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			func() args {
				ln, err := lang.NewLang("en", "testAuthor")
				assert.Nil(t, err)
				return args{lang: ln}
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error())
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, _ := i2[0].(*cache.Cache[string, map[string]query.LangView]).Get("testAuthor")
				_, ok := cacheMap["lang1"]
				assert.True(t, ok)
				return true
			},
		},
		{
			"Cache is set",
			func() fields {
				domainProxy := lang.NewMockRepository(t)
				domainProxy.On("Update", mock.AnythingOfType("*lang.Lang")).Return(nil)
				c := cache.New[string, map[string]query.LangView]()
				c.Set("testAuthor", map[string]query.LangView{"lang1": {ID: "lang1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			func() args {
				ln, err := lang.NewLang("en", "testAuthor")
				assert.Nil(t, err)
				return args{lang: ln}
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				_, ok := i2[0].(*cache.Cache[string, map[string]query.LangView]).Get("testAuthor")
				assert.False(t, ok)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := LangRepo{
				domainProxy: f.domainProxy,
				cache:       f.cache,
				cacheTTL:    time.Minute,
			}
			if tt.wantErr(t, repo.Update(tt.argsFn().lang), fmt.Sprintf("Update(%v)", tt.argsFn().lang)) {
				return
			}
			tt.wantFn(t, tt.argsFn().lang, repo.cache, fmt.Sprintf("Update(%v)", tt.argsFn().lang))
		})
	}
}

func TestLangRepo_Delete(t *testing.T) {
	type fields struct {
		domainProxy lang.Repository
		cache       *cache.Cache[string, map[string]query.LangView]
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
			"Error on DB request",
			func() fields {
				domainProxy := lang.NewMockRepository(t)
				domainProxy.On("Delete", "lang1", "testAuthor").Return(fmt.Errorf("testErr"))
				c := cache.New[string, map[string]query.LangView]()
				c.Set("testAuthor", map[string]query.LangView{"lang1": {ID: "lang1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			args{
				id:       "lang1",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error())
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, _ := i.(*cache.Cache[string, map[string]query.LangView]).Get("testAuthor")
				_, ok := cacheMap["lang1"]
				assert.True(t, ok, i2)
				return true
			},
		},
		{
			"Cache is set",
			func() fields {
				domainProxy := lang.NewMockRepository(t)
				domainProxy.On("Delete", "lang1", "testAuthor").Return(nil)
				c := cache.New[string, map[string]query.LangView]()
				c.Set("testAuthor", map[string]query.LangView{"lang1": {ID: "lang1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			args{
				id:       "lang1",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				_, ok := i.(*cache.Cache[string, map[string]query.LangView]).Get("testAuthor")
				assert.False(t, ok, i2)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := LangRepo{
				domainProxy: f.domainProxy,
				cache:       f.cache,
				cacheTTL:    time.Minute,
			}
			if tt.wantErr(t, repo.Delete(tt.args.id, tt.args.authorID), fmt.Sprintf("Delete(%v:%v)", tt.args.id, tt.args.authorID)) {
				return
			}
			tt.wantFn(t, repo.cache, fmt.Sprintf("Delete(%v:%v)", tt.args.id, tt.args.authorID))
		})
	}
}

func TestLangRepo_GetAllViews(t *testing.T) {
	type fields struct {
		queryProxy query.LangViewRepository
		cache      *cache.Cache[string, map[string]query.LangView]
	}
	type args struct {
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
			"Cache is missed",
			func() fields {
				queryProxy := query.NewMockLangViewRepository(t)
				queryProxy.On("GetAllViews", "testAuthor").Return([]query.LangView{{ID: "lang1"}, {ID: "lang2"}}, nil)
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.LangView](),
				}
			},
			args{authorID: "testAuthor"},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				result := i.([]query.LangView)
				cacheMap, _ := i2[0].(*cache.Cache[string, map[string]query.LangView]).Get("testAuthor")
				assert.Equal(t, []query.LangView{{ID: "lang1"}, {ID: "lang2"}}, result)

				for r := range result {
					_, ok := cacheMap[result[r].ID]
					assert.True(t, ok)
				}
				return true
			},
		},
		{
			"Error on DB request",
			func() fields {
				queryProxy := query.NewMockLangViewRepository(t)
				queryProxy.On("GetAllViews", "testAuthor").Return(nil, fmt.Errorf("testError"))
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.LangView](),
				}
			},
			args{authorID: "testAuthor"},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testError", err.Error(), i)
				return false
			},
			nil,
		},
		{
			"Cache is present",
			func() fields {
				c := cache.New[string, map[string]query.LangView]()
				c.Set("testAuthor", map[string]query.LangView{"lang1": {ID: "lang1"}, "lang2": {ID: "lang2"}})
				return fields{
					queryProxy: nil,
					cache:      c,
				}
			},
			args{authorID: "testAuthor"},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				result := i.([]query.LangView)
				assert.Equal(t, []query.LangView{{ID: "lang1"}, {ID: "lang2"}}, result)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := LangRepo{
				queryProxy: f.queryProxy,
				cache:      f.cache,
				cacheTTL:   time.Minute,
			}
			got, err := repo.GetAllViews(tt.args.authorID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetAllViews(%v)", tt.args.authorID)) {
				return
			}
			tt.wantFn(t, got, repo.cache)
		})
	}
}

func TestLangRepo_GetView(t *testing.T) {
	type fields struct {
		queryProxy query.LangViewRepository
		cache      *cache.Cache[string, map[string]query.LangView]
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
			"Cache is missed",
			func() fields {
				queryProxy := query.NewMockLangViewRepository(t)
				queryProxy.On("GetAllViews", "testAuthor").Return([]query.LangView{{ID: "lang1"}, {ID: "lang2"}}, nil)
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.LangView](),
				}
			},
			args{
				id:       "lang1",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				result := i.(query.LangView)
				cacheMap, _ := i2[0].(*cache.Cache[string, map[string]query.LangView]).Get("testAuthor")
				assert.Equal(t, query.LangView{ID: "lang1"}, result, i2[1])
				assert.Equal(t, result, cacheMap[result.ID], i2[1])
				return true
			},
		},
		{
			"Error on DB request",
			func() fields {
				queryProxy := query.NewMockLangViewRepository(t)
				queryProxy.On("GetAllViews", "testAuthor").Return(nil, fmt.Errorf("testErr"))
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.LangView](),
				}
			},
			args{
				id:       "lang1",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
			nil,
		},
		{
			"LangID not found",
			func() fields {
				queryProxy := query.NewMockLangViewRepository(t)
				queryProxy.On("GetAllViews", "testAuthor").Return([]query.LangView{{ID: "lang1"}, {ID: "lang2"}}, nil)
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.LangView](),
				}
			},
			args{
				id:       "lang3",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not find lang, userID: testAuthor, langID: lang3", err.Error(), i)
				return true
			},
			nil,
		},
		{
			"Cache is present",
			func() fields {
				c := cache.New[string, map[string]query.LangView]()
				c.Set("testAuthor", map[string]query.LangView{"lang1": {ID: "lang1"}})
				return fields{
					cache: c,
				}
			},
			args{
				id:       "lang1",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				result := i.(query.LangView)
				assert.Equal(t, query.LangView{ID: "lang1"}, result, i2[1])
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			f := tt.fieldsFn()
			repo := LangRepo{
				queryProxy: f.queryProxy,
				cache:      f.cache,
				cacheTTL:   time.Minute,
			}
			got, err := repo.GetView(tt.args.id, tt.args.authorID)
			if tt.wantErr(t1, err, fmt.Sprintf("GetView(%v, %v)", tt.args.id, tt.args.authorID)) {
				return
			}
			tt.wantFn(t, got, repo.cache, "GetView(%v, %v)", tt.args.id, tt.args.authorID)
		})
	}
}

func TestLangRepo_DeleteByAuthorID(t *testing.T) {
	type fields struct {
		domainProxy lang.Repository
		cache       *cache.Cache[string, map[string]query.LangView]
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
			"Error on DB request",
			func() fields {
				domainProxy := lang.NewMockRepository(t)
				domainProxy.On("DeleteByAuthorID", "testAuthor").Return(0, fmt.Errorf("testErr"))
				c := cache.New[string, map[string]query.LangView]()
				c.Set("testAuthor", map[string]query.LangView{"lang1": {ID: "lang1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			args{authorID: "testAuthor"},
			0,
			assert.Error,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, _ := i.(*cache.Cache[string, map[string]query.LangView]).Get("testAuthor")
				_, ok := cacheMap["lang1"]
				assert.True(t, ok, i2)
				return true
			},
		},
		{
			"Cache is cleared",
			func() fields {
				domainProxy := lang.NewMockRepository(t)
				domainProxy.On("DeleteByAuthorID", "testAuthor").Return(5, nil)
				c := cache.New[string, map[string]query.LangView]()
				c.Set("testAuthor", map[string]query.LangView{"lang1": {ID: "lang1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			args{authorID: "testAuthor"},
			5,
			assert.NoError,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				_, ok := i.(*cache.Cache[string, map[string]query.LangView]).Get("testAuthor")
				assert.False(t, ok)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			l := LangRepo{
				domainProxy: f.domainProxy,
				cache:       f.cache,
			}
			got, err := l.DeleteByAuthorID(tt.args.authorID)
			tt.assertCacheFn(t, l.cache, fmt.Sprintf("DeleteByAuthorID(%v)", tt.args.authorID))
			if !tt.wantErr(t, err, fmt.Sprintf("DeleteByAuthorID(%v)", tt.args.authorID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "DeleteByAuthorID(%v)", tt.args.authorID)
		})
	}
}
