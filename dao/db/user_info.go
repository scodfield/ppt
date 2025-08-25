package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"ppt/dao"
	"ppt/log"
)

func UpdateUserBalance(mongoClient *mongo.Client, userID uint64, amount int64) (int64, error) {
	userCredit := mongoClient.Database(dao.MongoDBPPT).Collection(dao.MongoCollUserCredit)
	filter := bson.M{"user_id": userID}
	update := bson.M{"$inc": bson.M{"balance": amount}}
	updateOpts := options.FindOneAndUpdate()
	updateOpts.SetProjection(bson.M{"user_id": 1, "balance": 1})
	updateOpts.SetReturnDocument(options.After)
	var result map[string]interface{}
	if err := userCredit.FindOneAndUpdate(dao.Ctx, filter, update, updateOpts).Decode(&result); err != nil {
		return 0, err
	}
	return result["balance"].(int64), nil
}

func UpdateUserLogin(mongoClient *mongo.Client, userID uint64, loginTime int64, loginIP string) error {
	userLogin := mongoClient.Database(dao.MongoDBPPT).Collection(dao.MongoCollUserLogin)
	res, err := userLogin.InsertOne(dao.Ctx, map[string]interface{}{
		"user_id":    userID,
		"login_time": loginTime,
		"login_ip":   loginIP,
	})
	if err != nil {
		log.Error("UpdateUserLogin insert login log error", zap.Uint64("user_id", userID), zap.Int64("login_time", loginTime), zap.String("login_ip", loginIP), zap.Error(err))
		return err
	}
	log.Info("UpdateUserLogin success to insert one login log", zap.Uint64("user_id", userID), zap.Any("inserted_id", res.InsertedID))
	return nil
}
