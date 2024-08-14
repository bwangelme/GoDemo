package main

import (
	"fmt"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func (*Product) TableName() string {
	return "product"
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "zyb:passwd@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Create
	//db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 3) // 根据整型主键查找
	fmt.Println(product.ID)
	product = Product{}
	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	fmt.Println(product.ID)

	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	fmt.Println(product.ID, product.UpdatedAt)

	// Delete - 删除 product
	db.Delete(&product, 1)
}
