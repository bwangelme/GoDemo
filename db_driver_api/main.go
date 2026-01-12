// db_driver_api
// 一个实现了 sql/driver 接口的程序，底层是使用 excel 存储文件
package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExcelDriver 实现 driver.Driver 接口
type ExcelDriver struct{}

// Open 打开一个 Excel 文件作为数据库
func (d *ExcelDriver) Open(name string) (driver.Conn, error) {
	return &ExcelConn{filePath: name}, nil
}

// ExcelConn 实现 driver.Conn 接口
type ExcelConn struct {
	filePath string
	file     *excelize.File
}

func (c *ExcelConn) Prepare(query string) (driver.Stmt, error) {
	return &ExcelStmt{conn: c, query: query}, nil
}

func (c *ExcelConn) Close() error {
	if c.file != nil {
		c.file.Close()
	}
	return nil
}

func (c *ExcelConn) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("transactions not supported")
}

// ExcelStmt 实现 driver.Stmt 接口
type ExcelStmt struct {
	conn  *ExcelConn
	query string
}

func (s *ExcelStmt) Close() error {
	return nil
}

func (s *ExcelStmt) NumInput() int {
	return -1 // 不限制参数数量
}

func (s *ExcelStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, fmt.Errorf("exec not supported")
}

func (s *ExcelStmt) Query(args []driver.Value) (driver.Rows, error) {
	return parseAndExecuteQuery(s.conn, s.query, args)
}

// ExcelRows 实现 driver.Rows 接口
type ExcelRows struct {
	columns []string
	data    [][]string
	current int
}

func (r *ExcelRows) Columns() []string {
	return r.columns
}

func (r *ExcelRows) Close() error {
	return nil
}

func (r *ExcelRows) Next(dest []driver.Value) error {
	if r.current >= len(r.data) {
		return io.EOF
	}

	for i, v := range r.data[r.current] {
		if i < len(dest) {
			dest[i] = v
		}
	}
	r.current++
	return nil
}

// 解析并执行查询
func parseAndExecuteQuery(conn *ExcelConn, query string, args []driver.Value) (driver.Rows, error) {
	// 简单解析 SELECT * FROM table
	re := regexp.MustCompile(`SELECT\s+(.+)\s+FROM\s+(\w+)`)
	matches := re.FindStringSubmatch(strings.TrimSpace(query))
	if len(matches) != 3 {
		return nil, fmt.Errorf("unsupported query: %s", query)
	}

	selectFields := strings.TrimSpace(matches[1])
	tableName := strings.TrimSpace(matches[2])

	// 打开 Excel 文件（如果还没有打开的话）
	if conn.file == nil {
		f, err := excelize.OpenFile(conn.filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open Excel file: %v", err)
		}
		conn.file = f
	}

	// 检查工作表是否存在
	allSheets := conn.file.GetSheetMap()
	sheetExists := false
	sheetName := ""

	for sheetNum, name := range allSheets {
		if strings.EqualFold(name, tableName) {
			sheetExists = true
			sheetName = name
			break
		}
		// 也检查数字形式的 sheet name
		if fmt.Sprintf("%d", sheetNum) == tableName {
			sheetExists = true
			sheetName = name
			break
		}
	}

	if !sheetExists {
		return nil, fmt.Errorf("table (sheet) %s not found in Excel file", tableName)
	}

	// 获取工作表的所有行
	rows, err := conn.file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read sheet %s: %v", sheetName, err)
	}

	if len(rows) == 0 {
		return &ExcelRows{columns: []string{}, data: [][]string{}, current: 0}, nil
	}

	// 第一行作为列名
	headers := rows[0]
	selectedColumns := headers

	if selectFields != "*" {
		selectedColumns = strings.Split(selectFields, ",")
		for i := range selectedColumns {
			selectedColumns[i] = strings.TrimSpace(selectedColumns[i])
		}
	}

	// 构建结果数据
	var resultData [][]string
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		selectedRow := make([]string, len(selectedColumns))

		for j, col := range selectedColumns {
			// 找到列索引
			colIndex := -1
			for k, header := range headers {
				if strings.TrimSpace(header) == col {
					colIndex = k
					break
				}
			}

			if colIndex >= 0 && colIndex < len(row) {
				selectedRow[j] = row[colIndex]
			} else {
				selectedRow[j] = ""
			}
		}
		resultData = append(resultData, selectedRow)
	}

	return &ExcelRows{
		columns: selectedColumns,
		data:    resultData,
		current: 0,
	}, nil
}

func main() {
	// 注册驱动
	sql.Register("excel", &ExcelDriver{})

	// 创建示例 Excel 文件
	createSampleExcel()

	// 连接数据库（实际上是 Excel 文件）
	db, err := sql.Open("excel", "./sample.xlsx")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// 执行查询 - 从 Users 工作表查询
	fmt.Println("=== Querying Users sheet ===")
	rows, err := db.Query("SELECT name, age FROM Users")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	// 获取列名
	columns, _ := rows.Columns()
	fmt.Println("Columns:", columns)

	// 遍历结果
	for rows.Next() {
		var name, age string
		err := rows.Scan(&name, &age)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		fmt.Printf("Name: %s, Age: %s\n", name, age)
	}

	// 查询 Products 工作表
	fmt.Println("\n=== Querying Products sheet ===")
	rows2, err := db.Query("SELECT product_name, price FROM Products")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows2.Close()

	// 获取列名
	columns2, _ := rows2.Columns()
	fmt.Println("Columns:", columns2)

	// 遍历结果
	for rows2.Next() {
		var productName, price string
		err := rows2.Scan(&productName, &price)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		fmt.Printf("Product: %s, Price: %s\n", productName, price)
	}

	// 查询所有列
	fmt.Println("\n=== Querying all columns from Users ===")
	rows3, err := db.Query("SELECT * FROM Users")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows3.Close()

	// 获取列名
	columns3, _ := rows3.Columns()
	fmt.Println("Columns:", columns3)

	// 遍历结果
	for rows3.Next() {
		values := make([]interface{}, len(columns3))
		valuePtrs := make([]interface{}, len(columns3))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows3.Scan(valuePtrs...)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		for i, v := range values {
			fmt.Printf("%s: %v ", columns3[i], v)
		}
		fmt.Println()
	}
}

// 创建示例 Excel 文件
func createSampleExcel() {
	f := excelize.NewFile()

	// 删除默认工作表
	f.DeleteSheet("Sheet1")

	// 创建 Users 工作表
	usersSheet := "Users"
	f.NewSheet(usersSheet)

	// 添加表头
	f.SetCellValue(usersSheet, "A1", "name")
	f.SetCellValue(usersSheet, "B1", "age")
	f.SetCellValue(usersSheet, "C1", "city")

	// 添加数据
	f.SetCellValue(usersSheet, "A2", "毛一一")
	f.SetCellValue(usersSheet, "B2", "25")
	f.SetCellValue(usersSheet, "C2", "江西九江")

	f.SetCellValue(usersSheet, "A3", "孙二二")
	f.SetCellValue(usersSheet, "B3", "30")
	f.SetCellValue(usersSheet, "C3", "北京")

	f.SetCellValue(usersSheet, "A4", "周三三")
	f.SetCellValue(usersSheet, "B4", "35")
	f.SetCellValue(usersSheet, "C4", "山东烟台")

	// 创建 Products 工作表
	productsSheet := "Products"
	f.NewSheet(productsSheet)

	// 添加表头
	f.SetCellValue(productsSheet, "A1", "product_name")
	f.SetCellValue(productsSheet, "B1", "price")
	f.SetCellValue(productsSheet, "C1", "category")

	// 添加数据
	f.SetCellValue(productsSheet, "A2", "平板")
	f.SetCellValue(productsSheet, "B2", "999.99")
	f.SetCellValue(productsSheet, "C2", "电子产品")

	f.SetCellValue(productsSheet, "A3", "书藉")
	f.SetCellValue(productsSheet, "B3", "19.99")
	f.SetCellValue(productsSheet, "C3", "学习资料")

	f.SetCellValue(productsSheet, "A4", "手机")
	f.SetCellValue(productsSheet, "B4", "699.00")
	f.SetCellValue(productsSheet, "C4", "电子产品")

	// 保存文件
	f.SaveAs("./sample.xlsx")
}
