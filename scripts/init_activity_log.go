package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 连接MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// 选择数据库
	db := client.Database("insurance_db")

	// 删除已存在的集合
	err = db.Collection("activity_logs").Drop(ctx)
	if err != nil {
		log.Printf("删除集合失败: %v", err)
	}

	// 创建集合
	err = db.CreateCollection(ctx, "activity_logs")
	if err != nil {
		log.Printf("创建集合失败: %v", err)
	}

	// 创建索引
	collection := db.Collection("activity_logs")
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "company_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "operation_time", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "operation_type", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "module_name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "company_id", Value: 1}, {Key: "operation_time", Value: -1}},
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("创建索引失败: %v", err)
	}

	// 示例数据
	sampleData := []interface{}{
		bson.M{
			"log_id":         fmt.Sprintf("AL%d001", time.Now().Unix()),
			"user_id":        "admin",
			"username":       "admin",
			"company_id":     "",
			"company_name":   "平台管理",
			"operation_type": "login",
			"module_name":    "认证授权",
			"operation_desc": "用户登录系统",
			"request_url":    "/api/v1/auth/login",
			"request_method": "POST",
			"request_params": bson.M{},
			"ip_address":     "192.168.1.100",
			"user_agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"operation_time": time.Date(2024, 3, 15, 9, 30, 0, 0, time.UTC),
			"execution_time": 150,
			"result_status":  "success",
			"target_id":      "",
			"target_name":    "",
		},
		bson.M{
			"log_id":         fmt.Sprintf("AL%d002", time.Now().Unix()),
			"user_id":        "admin",
			"username":       "admin",
			"company_id":     "",
			"company_name":   "平台管理",
			"operation_type": "create",
			"module_name":    "公司管理",
			"operation_desc": "创建新公司：测试公司A",
			"request_url":    "/api/v1/companies",
			"request_method": "POST",
			"request_params": bson.M{
				"company_name":  "测试公司A",
				"contact_phone": "13800138000",
			},
			"ip_address":     "192.168.1.100",
			"user_agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"operation_time": time.Date(2024, 3, 15, 10, 15, 0, 0, time.UTC),
			"execution_time": 320,
			"result_status":  "success",
			"target_id":      "COMP001",
			"target_name":    "测试公司A",
		},
		bson.M{
			"log_id":         fmt.Sprintf("AL%d003", time.Now().Unix()),
			"user_id":        "user001",
			"username":       "zhangsan",
			"company_id":     "COMP001",
			"company_name":   "测试公司A",
			"operation_type": "create",
			"module_name":    "保单管理",
			"operation_desc": "创建新保单：A000000001",
			"request_url":    "/api/v1/policies",
			"request_method": "POST",
			"request_params": bson.M{
				"proposal_number":  "A000000001",
				"customer_name_cn": "张三",
			},
			"ip_address":     "192.168.1.101",
			"user_agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"operation_time": time.Date(2024, 3, 15, 11, 20, 0, 0, time.UTC),
			"execution_time": 450,
			"result_status":  "success",
			"target_id":      "POL001",
			"target_name":    "A000000001",
		},
		bson.M{
			"log_id":         fmt.Sprintf("AL%d004", time.Now().Unix()),
			"user_id":        "user001",
			"username":       "zhangsan",
			"company_id":     "COMP001",
			"company_name":   "测试公司A",
			"operation_type": "update",
			"module_name":    "保单管理",
			"operation_desc": "更新保单信息：A000000001",
			"request_url":    "/api/v1/policies/POL001",
			"request_method": "PUT",
			"request_params": bson.M{
				"customer_name_cn": "张三",
				"policy_currency":  "USD",
			},
			"ip_address":     "192.168.1.101",
			"user_agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"operation_time": time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC),
			"execution_time": 280,
			"result_status":  "success",
			"target_id":      "POL001",
			"target_name":    "A000000001",
		},
		bson.M{
			"log_id":         fmt.Sprintf("AL%d005", time.Now().Unix()),
			"user_id":        "admin",
			"username":       "admin",
			"company_id":     "",
			"company_name":   "平台管理",
			"operation_type": "create",
			"module_name":    "用户管理",
			"operation_desc": "创建新用户：lisi",
			"request_url":    "/api/v1/users",
			"request_method": "POST",
			"request_params": bson.M{
				"username":     "lisi",
				"display_name": "李四",
			},
			"ip_address":     "192.168.1.100",
			"user_agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"operation_time": time.Date(2024, 3, 15, 16, 45, 0, 0, time.UTC),
			"execution_time": 380,
			"result_status":  "success",
			"target_id":      "USER002",
			"target_name":    "lisi",
		},
		bson.M{
			"log_id":         fmt.Sprintf("AL%d006", time.Now().Unix()),
			"user_id":        "user002",
			"username":       "lisi",
			"company_id":     "COMP001",
			"company_name":   "测试公司A",
			"operation_type": "view",
			"module_name":    "保单管理",
			"operation_desc": "查看保单列表",
			"request_url":    "/api/v1/policies",
			"request_method": "GET",
			"request_params": bson.M{
				"page":      1,
				"page_size": 20,
			},
			"ip_address":     "192.168.1.102",
			"user_agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"operation_time": time.Date(2024, 3, 15, 17, 10, 0, 0, time.UTC),
			"execution_time": 120,
			"result_status":  "success",
			"target_id":      "",
			"target_name":    "",
		},
		bson.M{
			"log_id":         fmt.Sprintf("AL%d007", time.Now().Unix()),
			"user_id":        "user001",
			"username":       "zhangsan",
			"company_id":     "COMP001",
			"company_name":   "测试公司A",
			"operation_type": "export",
			"module_name":    "保单管理",
			"operation_desc": "导出保单数据",
			"request_url":    "/api/v1/policies/export",
			"request_method": "POST",
			"request_params": bson.M{
				"export_type": "xlsx",
			},
			"ip_address":     "192.168.1.101",
			"user_agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"operation_time": time.Date(2024, 3, 15, 18, 20, 0, 0, time.UTC),
			"execution_time": 2500,
			"result_status":  "success",
			"target_id":      "",
			"target_name":    "",
		},
		bson.M{
			"log_id":         fmt.Sprintf("AL%d008", time.Now().Unix()),
			"user_id":        "admin",
			"username":       "admin",
			"company_id":     "",
			"company_name":   "平台管理",
			"operation_type": "delete",
			"module_name":    "公司管理",
			"operation_desc": "删除公司：测试公司B",
			"request_url":    "/api/v1/companies/COMP002",
			"request_method": "DELETE",
			"request_params": bson.M{},
			"ip_address":     "192.168.1.100",
			"user_agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"operation_time": time.Date(2024, 3, 15, 19, 30, 0, 0, time.UTC),
			"execution_time": 180,
			"result_status":  "success",
			"target_id":      "COMP002",
			"target_name":    "测试公司B",
		},
		bson.M{
			"log_id":         fmt.Sprintf("AL%d009", time.Now().Unix()),
			"user_id":        "user001",
			"username":       "zhangsan",
			"company_id":     "COMP001",
			"company_name":   "测试公司A",
			"operation_type": "import",
			"module_name":    "保单管理",
			"operation_desc": "批量导入保单数据",
			"request_url":    "/api/v1/policies/import",
			"request_method": "POST",
			"request_params": bson.M{
				"file_name":   "policies_import.xlsx",
				"total_count": 50,
			},
			"ip_address":     "192.168.1.101",
			"user_agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"operation_time": time.Date(2024, 3, 15, 20, 15, 0, 0, time.UTC),
			"execution_time": 5000,
			"result_status":  "success",
			"target_id":      "",
			"target_name":    "",
		},
		bson.M{
			"log_id":         fmt.Sprintf("AL%d010", time.Now().Unix()),
			"user_id":        "admin",
			"username":       "admin",
			"company_id":     "",
			"company_name":   "平台管理",
			"operation_type": "logout",
			"module_name":    "认证授权",
			"operation_desc": "用户退出登录",
			"request_url":    "/api/v1/auth/logout",
			"request_method": "POST",
			"request_params": bson.M{},
			"ip_address":     "192.168.1.100",
			"user_agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"operation_time": time.Date(2024, 3, 15, 22, 0, 0, 0, time.UTC),
			"execution_time": 80,
			"result_status":  "success",
			"target_id":      "",
			"target_name":    "",
		},
	}

	// 插入数据
	result, err := collection.InsertMany(ctx, sampleData)
	if err != nil {
		log.Fatal("插入数据失败:", err)
	}

	// 验证数据
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("统计文档失败:", err)
	}

	fmt.Printf("活动记录表初始化完成！\n")
	fmt.Printf("插入记录数：%d\n", len(result.InsertedIDs))
	fmt.Printf("总记录数：%d\n", count)
	fmt.Printf("索引创建完成！\n")

	// 显示前几条记录作为验证
	fmt.Printf("\n前3条记录示例：\n")
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetLimit(3))
	if err != nil {
		log.Fatal("查询数据失败:", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Fatal("解码文档失败:", err)
		}
		fmt.Printf("ID: %s, 用户: %s, 操作: %s, 时间: %s\n",
			doc["log_id"], doc["username"], doc["operation_desc"], doc["operation_time"])
	}

	fmt.Printf("\n初始化脚本执行完成！\n")
}
