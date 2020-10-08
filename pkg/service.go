package pkg

import (
	"context"
	"io"
	"path"
	"path/filepath"
	"strings"

	kit_log "github.com/go-kit/kit/log"
)

type CacheControlOption string

const (
	DisableCache = "no-cache"
	PublicCache  = "max-age=86400"
)

type Object struct {
	Metadata           map[string]string `json:"metadata"`
	Body               io.Reader         `json:"-"`
	ContentType        string            `json:"contentType"`
	CacheControl       string            `json:"cacheControl"`
	ContentDisposition string            `json:"contentDisposition"`
	Etag               string            `json:"etag"`
}

type Service interface {
	GetObject(ctx context.Context, host string, path string, object *Object) error
}

type Storage interface {
	OpenObject(ctx context.Context, name string, object *Object) error
}

type service struct {
	Storage Storage
	Logger  kit_log.Logger
	Index   string
}

func NewService(logger kit_log.Logger, Storage Storage) Service {
	return &service{
		Storage: Storage,
		Logger:  logger,
		Index:   "index.html",
	}
}

func (s *service) configureObject(name string, object *Object) {
	object.CacheControl = PublicCache
	if object.ContentType == "text/html" {
		object.CacheControl = DisableCache
	}
	_ = s.Logger.Log("object", name, "type", object.ContentType, "cache", object.CacheControl)
}

func (s *service) GetObject(ctx context.Context, host string, p string, object *Object) error {
	isDir := !strings.HasSuffix(p, s.Index) && filepath.Ext(p) == ""
	if isDir {
		ps := path.Join(p, s.Index)
		err := s.Storage.OpenObject(ctx, path.Join(host, ps), object)
		if err == nil {
			s.configureObject(path.Join(host, ps), object)
			return nil
		}
	}
	err := s.Storage.OpenObject(ctx, path.Join(host, p), object)
	if err != nil {
		return err
	}
	s.configureObject(path.Join(host, p), object)
	return nil
}
