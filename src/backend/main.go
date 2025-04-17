package main

import (
	"defi-backend/config"
	"defi-backend/routes"
	"fmt"
	"log"
	"os"
)

func main() {
	// 加载本地配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建 Nacos 客户端
	nacosClient, err := config.NewNacosClient(cfg.Nacos)
	if err != nil {
		log.Fatalf("Failed to create Nacos client: %v", err)
	}

	// 从 Nacos 获取配置
	content, err := config.GetConfig(nacosClient, cfg.Nacos.Group, cfg.Nacos.DataId)
	if err != nil {
		log.Fatalf("Failed to get config from Nacos: %v", err)
	}

	fmt.Printf("Configuration loaded successfully:\n%s\n", content)

	// 设置路由
	r := routes.SetupRouter()

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务器
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
