package pkg

import (
	"context"

	kit_log "github.com/go-kit/kit/log"

	gstorage "cloud.google.com/go/storage"

	weasel "github.com/revas-hq/weasel/pkg"
)

type storage struct {
	Bucket *gstorage.BucketHandle
	Logger kit_log.Logger
}

func NewStorage(logger kit_log.Logger, client *gstorage.Client, bucketName string) weasel.Storage {
	return &storage{
		Bucket: client.Bucket(bucketName),
		Logger: logger,
	}
}

func setMetadata(attrs *gstorage.ObjectAttrs, object *weasel.Object) {
	object.Metadata = map[string]string{}
	object.ContentType = attrs.ContentType
	object.CacheControl = attrs.CacheControl
	object.ContentDisposition = attrs.ContentDisposition
	object.Etag = attrs.Etag
}

func (s *storage) OpenObject(ctx context.Context, name string, object *weasel.Object) error {
	_ = s.Logger.Log("object", name)
	obj := s.Bucket.Object(name)
	r, err := obj.NewReader(ctx)
	if err != nil {
		_ = s.Logger.Log("object", name, "err", err)
		return err
	}
	as, err := obj.Attrs(ctx)
	if err != nil {
		_ = s.Logger.Log("object", name, "err", err)
		return err
	}
	object.Body = r
	setMetadata(as, object)
	return nil
}
