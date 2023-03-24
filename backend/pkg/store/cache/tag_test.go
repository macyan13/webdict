package cache

import (
	"context"
	"fmt"
	"github.com/Code-Hex/go-generics-cache"
	"github.com/macyan13/webdict/backend/pkg/app/domain/tag"
	"github.com/macyan13/webdict/backend/pkg/app/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestTagRepo_updateCacheWithMisses(t *testing.T) {
	type fields struct {
		queryProxy query.TagViewRepository
	}
	type args struct {
		ids      []string
		authorID string
		cacheMap map[string]query.TagView
	}
	tests := []struct {
		name     string
		fieldsFn func() fields
		args     args
		wantErr  bool
		wantFn   assert.ValueAssertionFunc
	}{
		{
			"Error on receiving tag views from repo",
			func() fields {
				queryProxy := query.NewMockTagViewRepository(t)
				queryProxy.On("GetViews", []string{"ids"}, "testAuthor").Return(nil, fmt.Errorf("test"))
				return fields{
					queryProxy: queryProxy,
				}
			},
			args{
				ids:      []string{"ids"},
				authorID: "testAuthor",
				cacheMap: map[string]query.TagView{},
			},
			true,
			nil,
		},
		{
			"Missed views are filled",
			func() fields {
				queryProxy := query.NewMockTagViewRepository(t)
				queryProxy.On("GetViews", []string{"miss2", "miss3"}, "testAuthor").Return([]query.TagView{{ID: "miss2"}, {ID: "miss3"}}, nil)
				return fields{
					queryProxy: queryProxy,
				}
			},
			args{
				ids:      []string{"miss1", "miss2", "miss3"},
				authorID: "testAuthor",
				cacheMap: map[string]query.TagView{"miss1": {ID: "miss1"}},
			},
			false,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, ok := i.(map[string]query.TagView)
				assert.True(t, ok)
				assert.Equal(t, "miss2", cacheMap["miss2"].ID)
				assert.Equal(t, "miss3", cacheMap["miss3"].ID)
				return true
			},
		},
		{
			"All views are cached",
			func() fields {
				return fields{
					queryProxy: nil,
				}
			},
			args{
				ids:      []string{"view1", "view2"},
				authorID: "testAuthor",
				cacheMap: map[string]query.TagView{"view1": {ID: "view1"}, "view2": {ID: "view2"}},
			},
			false,
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, ok := i.(map[string]query.TagView)
				assert.True(t, ok)
				assert.Equal(t, "view1", cacheMap["view1"].ID)
				assert.Equal(t, "view2", cacheMap["view2"].ID)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := NewTagRepo(context.TODO(), nil, f.queryProxy, Opts{TagCacheTTL: time.Minute})
			err := repo.updateCacheWithMisses(tt.args.ids, tt.args.authorID, tt.args.cacheMap)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("updateCacheWithMisses() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			tt.wantFn(t, tt.args.cacheMap)
		})
	}
}

func TestTagRepo_GetAllViews(t *testing.T) {
	type fields struct {
		queryProxy query.TagViewRepository
		cache      *cache.Cache[string, map[string]query.TagView]
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
				queryProxy := query.NewMockTagViewRepository(t)
				queryProxy.On("GetAllViews", "testAuthor").Return([]query.TagView{{ID: "tag1"}, {ID: "tag2"}}, nil)
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.TagView](),
				}
			},
			args{authorID: "testAuthor"},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				result := i.([]query.TagView)
				cacheMap, _ := i2[0].(*cache.Cache[string, map[string]query.TagView]).Get("testAuthor")
				assert.Equal(t, []query.TagView{{ID: "tag1"}, {ID: "tag2"}}, result)

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
				queryProxy := query.NewMockTagViewRepository(t)
				queryProxy.On("GetAllViews", "testAuthor").Return(nil, fmt.Errorf("testError"))
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.TagView](),
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
				c := cache.New[string, map[string]query.TagView]()
				c.Set("testAuthor", map[string]query.TagView{"tag1": {ID: "tag1"}, "tag2": {ID: "tag2"}})
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
				result := i.([]query.TagView)
				assert.Equal(t, []query.TagView{{ID: "tag1"}, {ID: "tag2"}}, result)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := TagRepo{
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

func TestTagRepo_GetView(t *testing.T) {
	type fields struct {
		queryProxy query.TagViewRepository
		cache      *cache.Cache[string, map[string]query.TagView]
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
				queryProxy := query.NewMockTagViewRepository(t)
				queryProxy.On("GetViews", []string{"tag1"}, "testAuthor").Return([]query.TagView{{ID: "tag1"}}, nil)
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.TagView](),
				}
			},
			args{
				id:       "tag1",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				result := i.(query.TagView)
				cacheMap, _ := i2[0].(*cache.Cache[string, map[string]query.TagView]).Get("testAuthor")
				assert.Equal(t, query.TagView{ID: "tag1"}, result, i2[1])
				assert.Equal(t, result, cacheMap[result.ID], i2[1])
				return true
			},
		},
		{
			"Error on DB request",
			func() fields {
				queryProxy := query.NewMockTagViewRepository(t)
				queryProxy.On("GetViews", []string{"tag1"}, "testAuthor").Return(nil, fmt.Errorf("testErr"))
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.TagView](),
				}
			},
			args{
				id:       "tag1",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
			nil,
		},
		{
			"Tag not found",
			func() fields {
				queryProxy := query.NewMockTagViewRepository(t)
				queryProxy.On("GetViews", []string{"tag1"}, "testAuthor").Return([]query.TagView{}, nil)
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.TagView](),
				}
			},
			args{
				id:       "tag1",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "can not find tag, userID: tag1, tagID: testAuthor", err.Error(), i)
				return true
			},
			nil,
		},
		{
			"Cache is present",
			func() fields {
				c := cache.New[string, map[string]query.TagView]()
				c.Set("testAuthor", map[string]query.TagView{"tag1": {ID: "tag1"}})
				return fields{
					cache: c,
				}
			},
			args{
				id:       "tag1",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				result := i.(query.TagView)
				assert.Equal(t, query.TagView{ID: "tag1"}, result, i2[1])
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			f := tt.fieldsFn()
			repo := TagRepo{
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

func TestTagRepo_GetViews(t *testing.T) {
	type fields struct {
		queryProxy query.TagViewRepository
		cache      *cache.Cache[string, map[string]query.TagView]
	}
	type args struct {
		ids      []string
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
				queryProxy := query.NewMockTagViewRepository(t)
				queryProxy.On("GetViews", []string{"tag1", "tag2"}, "testAuthor").Return([]query.TagView{{ID: "tag1"}, {ID: "tag2"}}, nil)
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.TagView](),
				}
			},
			args{
				ids:      []string{"tag1", "tag2"},
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				result := i.([]query.TagView)
				cacheMap, _ := i2[0].(*cache.Cache[string, map[string]query.TagView]).Get("testAuthor")
				assert.Equal(t, []query.TagView{{ID: "tag1"}, {ID: "tag2"}}, result)

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
				queryProxy := query.NewMockTagViewRepository(t)
				queryProxy.On("GetViews", []string{"tag1"}, "testAuthor").Return(nil, fmt.Errorf("testErr"))
				return fields{
					queryProxy: queryProxy,
					cache:      cache.New[string, map[string]query.TagView](),
				}
			},
			args{
				ids:      []string{"tag1"},
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error(), i)
				return true
			},
			nil,
		},
		{
			"Cache is present",
			func() fields {
				c := cache.New[string, map[string]query.TagView]()
				c.Set("testAuthor", map[string]query.TagView{"tag1": {ID: "tag1"}, "tag2": {ID: "tag2"}})
				return fields{
					cache: c,
				}
			},
			args{
				ids:      []string{"tag1", "tag2"},
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				result := i.([]query.TagView)
				ids := []string{"tag1", "tag2"}
				for _, id := range ids {
					var view query.TagView
					for r := range result {
						if result[r].ID == id {
							view = result[r]
						}
					}
					assert.Equal(t, id, view.ID)
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			f := tt.fieldsFn()
			repo := TagRepo{
				queryProxy: f.queryProxy,
				cache:      f.cache,
				cacheTTL:   time.Minute,
			}
			got, err := repo.GetViews(tt.args.ids, tt.args.authorID)
			if tt.wantErr(t1, err, fmt.Sprintf("GetViews(%v, %v)", tt.args.ids, tt.args.authorID)) {
				return
			}
			tt.wantFn(t, got, repo.cache, "GetViews(%v, %v)", tt.args.ids, tt.args.authorID)
		})
	}
}

func TestTagRepo_Update(t *testing.T) {
	type fields struct {
		domainProxy tag.Repository
		cache       *cache.Cache[string, map[string]query.TagView]
	}
	type args struct {
		tag *tag.Tag
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
				domainProxy := tag.NewMockRepository(t)
				domainProxy.On("Update", mock.AnythingOfType("*tag.Tag")).Return(fmt.Errorf("testErr"))
				c := cache.New[string, map[string]query.TagView]()
				c.Set("testAuthor", map[string]query.TagView{"tag1": {ID: "tag1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			func() args {
				tg, err := tag.NewTag("tag1", "testAuthor")
				assert.Nil(t, err)
				return args{tag: tg}
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error())
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, _ := i2[0].(*cache.Cache[string, map[string]query.TagView]).Get("testAuthor")
				_, ok := cacheMap["tag1"]
				assert.True(t, ok)
				return true
			},
		},
		{
			"Cache is set",
			func() fields {
				domainProxy := tag.NewMockRepository(t)
				domainProxy.On("Update", mock.AnythingOfType("*tag.Tag")).Return(nil)
				c := cache.New[string, map[string]query.TagView]()
				c.Set("testAuthor", map[string]query.TagView{"tag1": {ID: "tag1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			func() args {
				tg, err := tag.NewTag("tag1", "testAuthor")
				assert.Nil(t, err)
				return args{tag: tg}
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, _ := i2[0].(*cache.Cache[string, map[string]query.TagView]).Get("testAuthor")
				_, ok := cacheMap[i.(*tag.Tag).ID()]
				assert.False(t, ok)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := TagRepo{
				domainProxy: f.domainProxy,
				cache:       f.cache,
				cacheTTL:    time.Minute,
			}
			if tt.wantErr(t, repo.Update(tt.argsFn().tag), fmt.Sprintf("Update(%v)", tt.argsFn().tag)) {
				return
			}
			tt.wantFn(t, tt.argsFn().tag, repo.cache)
		})
	}
}

func TestTagRepo_Delete(t *testing.T) {
	type fields struct {
		domainProxy tag.Repository
		cache       *cache.Cache[string, map[string]query.TagView]
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
				domainProxy := tag.NewMockRepository(t)
				domainProxy.On("Delete", "tag1", "testAuthor").Return(fmt.Errorf("testErr"))
				c := cache.New[string, map[string]query.TagView]()
				c.Set("testAuthor", map[string]query.TagView{"tag1": {ID: "tag1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			args{
				id:       "tag1",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error())
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, _ := i.(*cache.Cache[string, map[string]query.TagView]).Get("testAuthor")
				_, ok := cacheMap["tag1"]
				assert.True(t, ok)
				return true
			},
		},
		{
			"Cache is set",
			func() fields {
				domainProxy := tag.NewMockRepository(t)
				domainProxy.On("Delete", "tag1", "testAuthor").Return(nil)
				c := cache.New[string, map[string]query.TagView]()
				c.Set("testAuthor", map[string]query.TagView{"tag1": {ID: "tag1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			args{
				id:       "tag1",
				authorID: "testAuthor",
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, _ := i.(*cache.Cache[string, map[string]query.TagView]).Get("testAuthor")
				_, ok := cacheMap["tag1"]
				assert.False(t, ok)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := TagRepo{
				domainProxy: f.domainProxy,
				cache:       f.cache,
				cacheTTL:    time.Minute,
			}
			if tt.wantErr(t, repo.Delete(tt.args.id, tt.args.authorID), fmt.Sprintf("Delete(%v:%v)", tt.args.id, tt.args.authorID)) {
				return
			}
			tt.wantFn(t, repo.cache)
		})
	}
}

func TestTagRepo_Create(t *testing.T) {
	type fields struct {
		domainProxy tag.Repository
		cache       *cache.Cache[string, map[string]query.TagView]
	}
	type args struct {
		tag *tag.Tag
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
				domainProxy := tag.NewMockRepository(t)
				domainProxy.On("Create", mock.AnythingOfType("*tag.Tag")).Return(fmt.Errorf("testErr"))
				c := cache.New[string, map[string]query.TagView]()
				c.Set("testAuthor", map[string]query.TagView{"tag1": {ID: "tag1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			func() args {
				tg, err := tag.NewTag("tag1", "testAuthor")
				assert.Nil(t, err)
				return args{tag: tg}
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, "testErr", err.Error())
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, _ := i2[0].(*cache.Cache[string, map[string]query.TagView]).Get("testAuthor")
				_, ok := cacheMap["tag1"]
				assert.True(t, ok)
				return true
			},
		},
		{
			"Cache is set",
			func() fields {
				domainProxy := tag.NewMockRepository(t)
				domainProxy.On("Create", mock.AnythingOfType("*tag.Tag")).Return(nil)
				c := cache.New[string, map[string]query.TagView]()
				c.Set("testAuthor", map[string]query.TagView{"tag1": {ID: "tag1"}})
				return fields{
					domainProxy: domainProxy,
					cache:       c,
				}
			},
			func() args {
				tg, err := tag.NewTag("tag1", "testAuthor")
				assert.Nil(t, err)
				return args{tag: tg}
			},
			func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return false
			},
			func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				cacheMap, _ := i2[0].(*cache.Cache[string, map[string]query.TagView]).Get("testAuthor")
				assert.Zero(t, cacheMap)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fieldsFn()
			repo := TagRepo{
				domainProxy: f.domainProxy,
				cache:       f.cache,
				cacheTTL:    time.Minute,
			}
			if tt.wantErr(t, repo.Create(tt.argsFn().tag), fmt.Sprintf("Update(%v)", tt.argsFn().tag)) {
				return
			}
			tt.wantFn(t, tt.argsFn().tag, repo.cache)
		})
	}
}
