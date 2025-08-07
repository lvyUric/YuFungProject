package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// 选择数据库
	db := client.Database("insurance_db")

	// 1. 添加活动记录菜单
	menuCollection := db.Collection("menus")
	menu := bson.M{
		"menu_id":         "activity_log",
		"parent_id":       "",
		"menu_name":       "活动记录",
		"menu_type":       "menu",
		"route_path":      "/activity-log",
		"component":       "ActivityLog",
		"permission_code": "activity:log:list",
		"icon":            "HistoryOutlined",
		"sort_order":      100,
		"visible":         true,
		"status":          "enable",
		"created_at":      time.Now(),
		"updated_at":      time.Now(),
	}

	_, err = menuCollection.InsertOne(context.Background(), menu)
	if err != nil {
		log.Printf("添加菜单失败: %v", err)
	} else {
		log.Println("活动记录菜单添加成功")
	}

	// 2. 为超级管理员角色添加活动记录权限
	roleCollection := db.Collection("roles")

	// 更新超级管理员角色
	_, err = roleCollection.UpdateOne(
		context.Background(),
		bson.M{"role_id": "super_admin"},
		bson.M{
			"$addToSet": bson.M{"menu_ids": "activity_log"},
			"$set":      bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		log.Printf("更新超级管理员角色失败: %v", err)
	} else {
		log.Println("超级管理员角色权限更新成功")
	}

	// 更新公司管理员角色
	_, err = roleCollection.UpdateOne(
		context.Background(),
		bson.M{"role_id": "company_admin"},
		bson.M{
			"$addToSet": bson.M{"menu_ids": "activity_log"},
			"$set":      bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		log.Printf("更新公司管理员角色失败: %v", err)
	} else {
		log.Println("公司管理员角色权限更新成功")
	}

	log.Println("活动记录菜单配置完成！")
}
