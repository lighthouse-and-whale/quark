package store

import (
	"context"
	"fmt"
	"github.com/lighthouse-and-whale/quark/core"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// 数据写入模式
// 1.  "w=0"          - 发起写操作，不关心是否成功
// 2.  "w=2"          - 发起写操作，数据被写入指定数量的节点才算成功
// 3.  "w=majority"   - 发起写操作，数据被写入大多数节点才算成功
// 4.  "journal=true" - 发起写操作，落地到journal日志文件中才算成功
const writeConcern = "w=majority"

// 数据读取模式
// 1.  "available"    - 读取所有可用的数据
// 1.  "local"        - 读取当前分片所有可用的数据
// 1.  "majority"     - 读取在大多数节点上落地的数据
// 1.  "linearizable" - 线性化读取数据
// 1.  "snapshot"     - 读取最近快照中的数据
const readConcern = "readConcern=majority"

// 数据读取节点选择
// 1.  "primary"            - 只选择主节点
// 1.  "primaryPreferred"   - 优先择主节点，如果不可用则选择从节点
// 1.  "secondary"          - 只选择从节点
// 1.  "secondaryPreferred" - 优先择从节点，如果不可用则选择主节点
const readPreference = "readPreference=primaryPreferred"

type AggregateLookup struct {
	From         string
	LocalField   string
	ForeignField string
	As           string
	Project      bson.M
}

func NewStorageDatabase(auth bool, nodes string, name, user, pwd, replicaSet string) *mongo.Database {
	var client *mongo.Client
	var err error
	var urlString string
	if auth {
		format := "mongodb://%s:%s@%s/%s?replicaSet=%s&%s&%s&%s"
		urlString = fmt.Sprintf(format, user, pwd, nodes, name, replicaSet, writeConcern, readPreference, readConcern)
	} else {
		urlString = fmt.Sprintf("mongodb://%s", nodes)
	}
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(urlString))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		core.NewError(err, "store.NewStorageDatabase().client.Ping()")
		os.Exit(1)
	}
	return client.Database(name)
}

func MongoIndexesCreateMany(coll *mongo.Collection, keys []string) {
	models := make([]mongo.IndexModel, 0)
	for _, item := range keys {
		models = append(models, mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: item, Value: bsonx.Int64(1)}},
			Options: options.Index().SetBackground(true),
		})
	}
	if _, err := coll.Indexes().CreateMany(context.Background(), models); err != nil {
		log.Printf("[  error  ] %s store (*MongoDriver) IndexesCreateMany(): %s\n", coll.Name(), err)
	}
	return
}

func MongoIndexesCreateUnique(coll *mongo.Collection, key string) {
	models := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: key, Value: bsonx.String("")}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := coll.Indexes().CreateOne(context.Background(), models); err != nil {
		log.Printf("[  error  ] %s store (*MongoDriver) IndexesCreateUnique(): %s\n", coll.Name(), err)
	}
	return
}
