package dao

import (
	"context"
	"design/internal/conf"
	"design/internal/model"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Dao struct {
	c           *conf.Config
	mongoClient *mongo.Client
}

func New(c *conf.Config) *Dao {
	d := &Dao{
		c:           c,
		mongoClient: newMongodbClient(c.Mongodb),
	}
	return d
}

func newMongodbClient(c *conf.Mongodb) *mongo.Client {
	url := fmt.Sprintf("mongodb://%s/%s", c.Host, c.DatabaseName)
	clientOptions := options.Client().ApplyURI(url)
	auth := &options.Credential{Username: c.Username, Password: c.Password}
	clientOptions.Auth = auth
	clientOptions.SetMaxPoolSize(c.MaxPoolSize)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.TimeoutInMillisecond))
	defer cancel()

	// 连接到MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

func (d *Dao) GetDatabase() *mongo.Database {
	return d.mongoClient.Database(d.c.Mongodb.DatabaseName)
}

func (d *Dao) GetCollection(collectionName string) *mongo.Collection {
	return d.GetDatabase().Collection(collectionName)
}

func (d *Dao) GetMemberCollection() *mongo.Collection {
	return d.GetCollection("member")
}

func (d *Dao) Register(ctx context.Context, mem model.Member) (bool, error) {
	filter := bson.D{
		{"token", mem.Token},
	}
	count, err := d.GetMemberCollection().CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	}
	ir, err := d.GetMemberCollection().InsertOne(ctx, mem)
	if err != nil {
		return false, err
	}
	return ir.InsertedID != "", nil
}

func (d *Dao) ListMember(ctx context.Context, role, status string) ([]model.Member, error) {
	filter := bson.D{}
	if role != "" {
		filter = append(filter, bson.E{"role", role})
	}
	if status != "" {
		filter = append(filter, bson.E{"status", status})
	}
	result := []model.Member{}
	cur, err := d.GetMemberCollection().Find(ctx, filter)
	if err != nil {
		return result, err
	}
	defer cur.Close(ctx)
	cur.All(ctx, &result)
	return result, nil
}

func (d *Dao) SetEmail(ctx context.Context, token, email string) (bool, error) {
	filter := bson.D{
		{"token", token},
	}
	update := bson.D{
		{"$set", bson.D{
			{"email", email}},
		},
	}
	ur, err := d.GetMemberCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return ur.ModifiedCount > 0, nil
}

func (d *Dao) GetMember(ctx context.Context, token string) (model.Member, error) {
	filter := bson.D{
		{"token", token},
	}
	var result model.Member
	err := d.GetMemberCollection().FindOne(ctx, filter).Decode(&result)
	return result, err
}
