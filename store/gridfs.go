package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/lighthouse-and-whale/quark/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"io"
	"time"
)

type FileInfo struct {
	GFS        *mongo.Database
	FileId     primitive.ObjectID `bson:"_id" json:"fileId"`
	Length     float64            `bson:"length" json:"length"`
	ChunkSize  int64              `bson:"chunkSize" json:"chunkSize"`
	UploadDate time.Time          `bson:"uploadDate" json:"uploadDate"`
	Filename   string             `bson:"filename" json:"filename"`
}

func (f *FileInfo) Add(filename string, source io.Reader) (fid primitive.ObjectID, err error) {
	var bucket *gridfs.Bucket
	bucket, err = gridfs.NewBucket(f.GFS)
	if err != nil {
		panic(err)
	}
	fid, err = bucket.UploadFromStream(filename, source)
	if err != nil {
		core.NewError(err, "store.(f *FileInfo) Add().bucket.UploadFromStream()")
		return
	}
	return
}

func (f *FileInfo) Del(fid string) (err error) {
	var bucket *gridfs.Bucket
	bucket, err = gridfs.NewBucket(f.GFS)
	if err != nil {
		panic(err)
	}
	var oid primitive.ObjectID
	if oid, err = primitive.ObjectIDFromHex(fid); err != nil {
		panic(err)
	} else {
		if err = bucket.Delete(oid); err != nil {
			core.NewError(err, "store.(f *FileInfo) Del().bucket.Delete()")
			return
		}
		return
	}
}

func (f *FileInfo) Get(fid string) (stream *gridfs.DownloadStream, info FileInfo, err error) {
	var bucket *gridfs.Bucket
	bucket, err = gridfs.NewBucket(f.GFS)
	if err != nil {
		panic(err)
	}
	var oid primitive.ObjectID
	if oid, err = primitive.ObjectIDFromHex(fid); err != nil {
		panic(err)
	} else {
		var cur *mongo.Cursor
		cur, err = bucket.Find(primitive.M{"_id": oid})
		if err != nil {
			core.NewError(err, "store.(f *FileInfo) Get().bucket.Find()")
			return
		}
		var items = make([]FileInfo, 0)
		if err = cur.All(context.Background(), &items); err != nil {
			core.NewError(err, "store.(f *FileInfo) Get().cur.All()")
			return
		}
		if len(items) == 1 {
			info = items[0]
		} else {
			err = errors.New("GridFS.Get: len(items) != 1")
			return
		}
		stream, err = bucket.OpenDownloadStream(oid)
		if err != nil {
			core.NewError(err, "store.(f *FileInfo) Get().bucket.OpenDownloadStream()")
			return
		}
	}
	return
}

func (f *FileInfo) GetAll(filter interface{}) (output []FileInfo, err error) {
	if filter == nil {
		filter = map[string]interface{}{}
	}
	var bucket *gridfs.Bucket
	bucket, err = gridfs.NewBucket(f.GFS)
	if err != nil {
		panic(err)
	}
	var cur *mongo.Cursor
	cur, err = bucket.Find(filter)
	if err != nil {
		core.NewError(err, "store.(f *FileInfo) GetAll().bucket.Find()")
		return
	}
	if err = cur.All(context.Background(), &output); err != nil {
		core.NewError(err, "store.(f *FileInfo) GetAll().cur.All()")
		return
	}
	return
}

func (f *FileInfo) Len(filter interface{}) (count int64, err error) {
	if filter == nil {
		filter = map[string]interface{}{}
	}
	var bucket *gridfs.Bucket
	bucket, err = gridfs.NewBucket(f.GFS)
	if err != nil {
		panic(err)
	}
	var cur *mongo.Cursor
	cur, err = bucket.Find(filter)
	if err != nil {
		core.NewError(err, "store.(f *FileInfo) Len().bucket.Find()")
		return
	}
	ctx := context.Background()
	for cur.Next(ctx) {
		count++
	}
	return
}

// 文件大小
type FileSize struct {
	Length float64
	Format string
}

func FileSizeFormat(size float64) (format string) {
	switch {
	case size > 1024*1024*1024:
		format = fmt.Sprintf("%.2fGB", size/1024/1024/1024)
	case size > 1024*1024:
		format = fmt.Sprintf("%.2fMB", size/1024/1024)
	case size > 1024:
		format = fmt.Sprintf("%.2fKB", size/1024)
	}
	return
}
