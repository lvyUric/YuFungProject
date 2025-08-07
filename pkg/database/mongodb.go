package database

import (
	"context"
	"time"

	"YufungProject/configs"
	"YufungProject/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database
var mongoClient *mongo.Client

// InitMongoDB 初始化MongoDB连接
func InitMongoDB(config configs.MongoDBConfig) (*mongo.Database, error) {
	// 解析超时时间
	timeout, err := time.ParseDuration(config.Timeout)
	if err != nil {
		timeout = 10 * time.Second // 默认10秒
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 设置客户端选项
	clientOptions := options.Client().ApplyURI(config.URI)
	clientOptions.SetMaxPoolSize(uint64(config.MaxPoolSize))

	// 连接到MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Errorf("MongoDB连接失败: %v", err)
		return nil, err
	}

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Errorf("MongoDB连接测试失败: %v", err)
		return nil, err
	}

	// 保存客户端引用
	mongoClient = client
	MongoDB = client.Database(config.Database)

	logger.Infof("MongoDB连接成功: %s/%s", config.URI, config.Database)
	return MongoDB, nil
}

// DisconnectMongoDB 断开MongoDB连接
func DisconnectMongoDB() error {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := mongoClient.Disconnect(ctx); err != nil {
			logger.Errorf("MongoDB断开连接失败: %v", err)
			return err
		}
		logger.Info("MongoDB连接已关闭")
	}
	return nil
}
