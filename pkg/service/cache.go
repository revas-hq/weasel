package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/allegro/bigcache"
	kit_log "github.com/go-kit/kit/log"

	weasel "github.com/revas-hq/weasel/pkg"
)

type cache struct {
	Logger kit_log.Logger
	Cache  *bigcache.BigCache
	Next   weasel.Service
}

func NewCache(Logger kit_log.Logger, Cache *bigcache.BigCache, Next weasel.Service) weasel.Service {
	return &cache{
		Logger: Logger,
		Cache:  Cache,
		Next:   Next,
	}
}

func (s *cache) getCache(ctx context.Context, name string, object *weasel.Object) error {
	meta, err := s.Cache.Get(name + ":meta")
	if err != nil {
		return err
	}
	body, err := s.Cache.Get(name)
	if err != nil {
		return err
	}

	err = json.Unmarshal(meta, &object)
	if err != nil {
		return err
	}
	object.Body = bytes.NewReader(body)
	return nil
}

func (s *cache) setCache(ctx context.Context, name string, object *weasel.Object) error {
	if object.CacheControl == weasel.DisableCache {
		return nil
	}
	meta, err := json.Marshal(object)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(object.Body)
	if err != nil {
		return err
	}

	err = s.Cache.Set(name+":meta", meta)
	if err != nil {
		return err
	}
	err = s.Cache.Set(name, body)
	if err != nil {
		return err
	}
	object.Body = bytes.NewReader(body)
	return nil
}

func (s *cache) GetObject(ctx context.Context, host string, p string, object *weasel.Object) error {
	name := path.Join(host, p)
	err := s.getCache(ctx, name, object)
	if err == nil {
		_ = s.Logger.Log("object", name, "cache", "hit")
		return nil
	}
	_ = s.Logger.Log("object", name, "cache", "miss")
	err = s.Next.GetObject(ctx, host, p, object)
	if err != nil {
		_ = s.Logger.Log("object", name, "err", err)
		return err
	}
	err = s.setCache(ctx, name, object)
	if err != nil {
		_ = s.Logger.Log("object", name, "err", err)
		return err
	}
	return nil
}
