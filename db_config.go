package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"db_name"`
}

type AppConfig struct {
	ServerPort string         `json:"server_port"`
	Database   DatabaseConfig `json:"database"`
	LogLevel   string         `json:"log_level"`
}

// 全局变量，用于存储加载的配置
var (
	db     *gorm.DB
	config AppConfig
)

func loadConfig(path string) error {
	// 打开配置文件
	configFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("无法打开配置文件 %s: %v", path, err)
	}
	defer configFile.Close()

	// 读取文件内容
	bytes, err := io.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("无法读取配置文件: %v", err)
	}

	// 使用 'encoding/json' 将文件内容解析(Unmarshal)到 AppConfig 结构体中
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return fmt.Errorf("无法解析 JSON 配置: %v", err)
	}

	log.Println("配置加载成功。")
	return nil
}

func initDb() {
	// 从 'config' 变量动态构建 DSN，而不是硬编码
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
	)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("无法从 GORM 获取 sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("成功连接到 MySQL 数据库 (%s@%s)!", config.Database.User, config.Database.Host)
}
