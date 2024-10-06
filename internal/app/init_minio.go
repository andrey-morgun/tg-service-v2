package app

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// TODO: Check bucket
func (a *App) initMinio(ctx context.Context) {
	var err error
	a.minio, err = minio.New(a.config.Minio.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(a.config.Minio.User, a.config.Minio.Password, ""),
	})
	if err != nil {
		a.logger.Fatal(err.Error())
		return
	}

	exists, err := a.minio.BucketExists(ctx, a.config.Minio.Bucket)
	if err != nil {
		a.logger.Fatal(err.Error())
		return
	}
	if !exists {
		err := a.minio.MakeBucket(ctx, a.config.Minio.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			a.logger.Fatal(err.Error())
			return
		}
	}
}
