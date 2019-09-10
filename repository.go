package file

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	CreateFile(ctx context.Context, name string, size int64, path string) string
	DeleteFile(ctx context.Context, id string)
	GetFilePath(ctx context.Context, id string) *string
	GetFiles(ctx context.Context,fileIds []string)[]Exported
}

type mongoExported struct {
	ID primitive.ObjectID `bson:"_id"`
	Size int64 `bson:"size"`
	Path string `bson:"path"`
	Name string `bson:"name"`
	Created int64 `bson:"created"`
}

type mongoExportedArray []mongoExported

func (m mongoExportedArray)asDomain()[]Exported{
	var result []Exported
	for _, item := range m{
		result = append(result, item.asDomain())
	}
	return result
}

func (m *mongoExported)asDomain()Exported{
	return Exported{
		ID:   m.ID.Hex(),
		Size: m.Size,
		Name: m.Name,
		Path: m.Path,
	}
}

const DBName = "file"
const CollectionName = "entity"

type MongoRepository struct {
	Client *mongo.Client
}

func (m *MongoRepository) GetFiles(ctx context.Context, fileIds []string) []Exported {
	var objIdList []primitive.ObjectID
	for _, item := range fileIds{
		objID, err := primitive.ObjectIDFromHex(item)
		if err != nil{
			panic(err)
		}
		objIdList = append(objIdList, objID)
	}
	var result mongoExportedArray
	cursor, err := m.collection().Find(ctx, bson.M{"_id": bson.M{"$in": objIdList}})
	if err != nil{
		panic(err)
	}
	err = cursor.All(ctx, &result)
	if err != nil{
		panic(err)
	}
	return result.asDomain()
}

func (m *MongoRepository)collection() *mongo.Collection{
	return m.Client.Database(DBName).Collection(CollectionName)
}

func (m *MongoRepository) DeleteFile(ctx context.Context, id string) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		panic(err)
	}
	_, err = m.collection().DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil{
		panic(err)
	}
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{Client: client}
}

func (m *MongoRepository) CreateFile(ctx context.Context, name string, size int64, path string) string {
	result, err := m.collection().InsertOne(ctx, bson.M{
		"size": size,
		"path": path,
		"name": name,
		"created": time.Now().Unix(),
	})
	if err != nil{
		panic(err)
	}
	return result.InsertedID.(primitive.ObjectID).Hex()
}

func (m *MongoRepository) GetFilePath(ctx context.Context, id string) *string {
	var result map[string]interface{}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		panic(err)
	}
	err = m.collection().FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil{
		return nil
	}
	if len(result) == 0{
		return nil
	}
	path, ok := result["path"].(string)
	if !ok{
		return nil
	}
	return &path
}




