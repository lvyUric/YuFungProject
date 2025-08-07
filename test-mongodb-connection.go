package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB连接配置
	uri := "mongodb://admin:yf2025@106.52.172.124:27017/insurance_db?authSource=admin"

	fmt.Println("🔍 开始测试MongoDB连接...")
	fmt.Printf("📍 连接地址: %s\n", uri)

	// 设置连接选项
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建客户端
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("❌ 创建MongoDB客户端失败: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("⚠️ 断开连接时出错: %v", err)
		}
	}()

	// 测试连接
	fmt.Println("🔄 正在测试连接...")
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ MongoDB连接失败: %v", err)
	}

	fmt.Println("✅ MongoDB连接成功！")

	// 获取数据库列表
	fmt.Println("📋 获取数据库列表...")
	databases, err := client.ListDatabaseNames(ctx, nil)
	if err != nil {
		log.Printf("⚠️ 获取数据库列表失败: %v", err)
	} else {
		fmt.Println("📚 可用数据库:")
		for _, db := range databases {
			fmt.Printf("  - %s\n", db)
		}
	}

	// 测试访问insurance_db数据库
	fmt.Println("🔍 测试访问insurance_db数据库...")
	database := client.Database("insurance_db")
	collections, err := database.ListCollectionNames(ctx, nil)
	if err != nil {
		log.Printf("⚠️ 获取集合列表失败: %v", err)
	} else {
		fmt.Println("📄 insurance_db数据库中的集合:")
		for _, collection := range collections {
			fmt.Printf("  - %s\n", collection)
		}
	}

	fmt.Println("�� MongoDB连接测试完成！")
}
