package mongodb

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	c *mongo.Client
}
type IDRange struct {
	Max string `json:"max"`
	Min string `json:"min"`
}

func NewClient(ctx context.Context, uri string) (*Client, error) {
	loggerOptions := options.
		Logger().
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	// Uses the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.
		Client().
		ApplyURI(uri).
		SetLoggerOptions(loggerOptions)

	// Defines the options for the MongoDB client
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Creates a new client and connects to the server
	client, err := mongo.Connect(ctx, opts, clientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{c: client}, nil
}

func (c *Client) Close() error {
	return c.c.Disconnect(context.Background())
}
func (c *Client) Ping() error {
	var result bson.M
	if err := c.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return err
	}
	return nil
}

func (c *Client) Database(name string) *mongo.Database {
	return c.c.Database(name)
}

func (c *Client) ListDatabaseNames() ([]string, error) {
	return c.c.ListDatabaseNames(context.TODO(), bson.M{})
}

func (c *Client) GetObjectIDRange(ctx context.Context, db, collection, key string, taskNum int64) ([]IDRange, error) {
	var (
		minRes bson.M
		maxRes bson.M
	)
	collect := c.Database(db).Collection(collection)
	if err := collect.
		FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.D{{key, 1}})).
		Decode(&minRes); err != nil {
		return nil, err
	}
	if err := collect.
		FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.D{{key, -1}})).
		Decode(&maxRes); err != nil {
		return nil, err
	}
	minID := minRes["_id"].(primitive.ObjectID).Hex()
	maxID := maxRes["_id"].(primitive.ObjectID).Hex()
	fmt.Printf("minRes: %v, maxRes: %v\n", minRes, maxRes)
	total, err := collect.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	docsPerTask := total / taskNum
	if total%taskNum != 0 {
		docsPerTask += 1
	}
	if docsPerTask == 0 {
		docsPerTask = 1
	}
	// 1. 添加第一个任务（从最小ID开始）
	firstTaskEndID, err := SplitObjectIdRanges(minID, maxID, taskNum)
	if err != nil {
		return nil, err
	}
	return firstTaskEndID, nil
}

// cloneBytes 安全复制字节切片
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// SplitObjectIdRanges 按 ObjectId 值域线性切分
func SplitObjectIdRanges(minOID, maxOID string, taskNum int64) ([]IDRange, error) {
	if taskNum <= 0 {
		return nil, fmt.Errorf("channelCount must be > 0")
	}

	// 1. 将 ObjectId 字符串转为 []byte（12字节）
	minBytes, err := hex.DecodeString(minOID)
	if err != nil || len(minBytes) != 12 {
		return nil, fmt.Errorf("invalid min ObjectId: %s", minOID)
	}
	maxBytes, err := hex.DecodeString(maxOID)
	if err != nil || len(maxBytes) != 12 {
		return nil, fmt.Errorf("invalid max ObjectId: %s", maxOID)
	}

	// 2. 转为大整数便于计算
	minInt := new(big.Int).SetBytes(minBytes)
	maxInt := new(big.Int).SetBytes(maxBytes)

	if minInt.Cmp(maxInt) >= 0 {
		return nil, fmt.Errorf("min ObjectId >= max ObjectId")
	}

	// 3. 计算步长
	diff := new(big.Int).Sub(maxInt, minInt)
	step := new(big.Int).Div(diff, big.NewInt(taskNum))

	ranges := make([]IDRange, taskNum)

	// 4. 生成每个区间
	current := new(big.Int).Set(minInt)
	for i := 0; i < int(taskNum); i++ {
		start := cloneBytes(current.Bytes())
		// 补齐到 12 字节
		for len(start) < 12 {
			start = append([]byte{0}, start...)
		}

		// 计算 end
		next := new(big.Int).Add(current, step)
		// 最后一个区间，end = maxInt + 1（确保覆盖）
		if i == int(taskNum)-1 {
			next = new(big.Int).Add(maxInt, big.NewInt(1))
		}

		endBytes := cloneBytes(next.Bytes())
		for len(endBytes) < 12 {
			endBytes = append([]byte{0}, endBytes...)
		}

		ranges[i] = IDRange{
			Min: hex.EncodeToString(start),
			Max: hex.EncodeToString(endBytes),
		}

		current = next
	}

	return ranges, nil
}

func (c *Client) GetDocByRange(ctx context.Context, db, collection, key string, minID, maxID string) ([]bson.M, error) {
	collect := c.Database(db).Collection(collection)
	fmt.Println("minID:", minID, "maxID:", maxID)
	minObjID, err := primitive.ObjectIDFromHex(minID)
	if err != nil {
		return nil, err

	}
	maxObjID, err := primitive.ObjectIDFromHex(maxID)
	if err != nil {
		return nil, err

	}
	filter := bson.M{
		"_id": bson.M{
			"$gte": minObjID,
			"$lte": maxObjID,
		},
	}
	cur, err := collect.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var docs []bson.M
	for cur.Next(ctx) {
		var doc bson.M
		err := cur.Decode(&doc)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return docs, nil
}
