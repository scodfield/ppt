package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"ppt/dao"
	"ppt/logger"
)

func FilterUsersByBrandID(client *mongo.Client, users []uint64, brandID int32) ([]uint64, error) {
	var result []uint64
	usersT := client.Database(dao.MongoDB).Collection(dao.MongoUsers)
	opts := options.Find()
	opts.SetProjection(bson.M{"user_id": 1})
	filter := bson.M{"brand_id": brandID}

	batchSize := 1000
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}
		batchUsers := users[i:end]
		filter["user_id"] = bson.M{"$in": batchUsers}
		cursor, err := usersT.Find(dao.Ctx, filter, opts)
		if err != nil {
			logger.Error("FilterUsersByBrandID find error", zap.Error(err))
			return nil, err
		}
		var tmp []map[string]interface{}
		if err = cursor.All(dao.Ctx, &tmp); err != nil {
			logger.Error("FilterUsersByBrandID cursor.All error", zap.Error(err))
			return nil, err
		}
		for _, u := range tmp {
			result = append(result, u["user_id"].(uint64))
		}
		cursor.Close(dao.Ctx)
	}
	return result, nil
}
