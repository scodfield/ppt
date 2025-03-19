package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ppt/dao"
)

func UpdateUserBalance(mongoClient *mongo.Client, userID uint64, amount int64) (int64, error) {
	userInfo := mongoClient.Database("ppt").Collection("user_info")
	filter := bson.M{"user_id": userID}
	update := bson.M{"$inc": bson.M{"balance": amount}}
	updateOpts := options.FindOneAndUpdate()
	updateOpts.SetProjection(bson.M{"user_id": 1, "balance": 1})
	updateOpts.SetReturnDocument(options.After)
	var result map[string]interface{}
	if err := userInfo.FindOneAndUpdate(dao.Ctx, filter, update, updateOpts).Decode(&result); err != nil {
		return 0, err
	}
	return result["balance"].(int64), nil
}
