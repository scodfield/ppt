package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"ppt/dao"
)

func FilterUsersByBrandID(client *mongo.Client, users []uint64, brandID int32) ([]uint64, error) {
	var result []uint64
	usersT := client.Database(dao.MongoDBPTT).Collection(dao.MongoCollUsers)
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
			log.Error("FilterUsersByBrandID find error", zap.Error(err))
			return nil, err
		}
		var tmp []map[string]interface{}
		if err = cursor.All(dao.Ctx, &tmp); err != nil {
			log.Error("FilterUsersByBrandID cursor.All error", zap.Error(err))
			return nil, err
		}
		for _, u := range tmp {
			result = append(result, u["user_id"].(uint64))
		}
		cursor.Close(dao.Ctx)
	}
	return result, nil
}

const FriendVisitBulkWriteSize = 1000

func UpdateUserFriendVisits(client *mongo.Client, userID uint64, visits []uint64) error {
	friendVisit := client.Database(dao.MongoDBPTT).Collection(dao.MongoCollFriendVisit)
	batchSize := 0
	var operations []mongo.WriteModel
	for i := 0; i < len(visits); i++ {
		if batchSize > FriendVisitBulkWriteSize {
			res, err := friendVisit.BulkWrite(dao.Ctx, operations)
			if err != nil {
				log.Error("UpdateUserFriendVisits bulk write error", zap.Error(err))
				return err
			}
			log.Info("UpdateUserFriendVisits bulk write result", zap.Any("result", res))
			batchSize = 0
			operations = operations[:0]
		}
		batchSize++
		updateModel := mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": userID}).SetUpsert(true).SetUpdate(bson.M{"$inc": bson.M{"today_visit": 1, "total_visit": 1}, "$push": bson.M{"today_visit_friends": visits[i]}})
		operations = append(operations, updateModel)
	}
	if len(operations) > 0 {
		res, err := friendVisit.BulkWrite(dao.Ctx, operations)
		if err != nil {
			log.Error("UpdateUserFriendVisits bulk write error", zap.Error(err))
			return err
		}
		log.Info("UpdateUserFriendVisits bulk write result", zap.Any("result", res))
		operations = nil
	}
	return nil
}
