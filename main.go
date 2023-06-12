package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置gin
	router := gin.Default()
	router.Use(Cors())

	// 获取目录下的所有文件
	router.GET("/fileName", func(ctx *gin.Context) {
		// 读取文件
		FileInfo, err := os.ReadDir("./data/")
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"code": http.StatusBadGateway,
				"msg":  fmt.Sprint("读取文件错误,err:", err),
			})
		}

		// 储存文件的名字
		var nameArr []string
		for _, info := range FileInfo {
			nameArr = append(nameArr, info.Name())
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "文件读取成功",
			"data": nameArr,
		})
	})

	// 获取数据
	router.GET("/get", func(ctx *gin.Context) {
		// 获取文件的路径
		filePath := ctx.Query("file")

		// 打开文件
		file, err := os.Open(path.Join("./data", filePath))
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"code": http.StatusBadGateway,
				"msg":  fmt.Sprint("读取文件错误,err:", err),
			})
			return
		}
		defer file.Close()

		// 储存文件中的数据
		var dataArr []interface{} = make([]interface{}, 0)
		// 读取的所有
		index := 1

		// 足行扫描文件
		buf := bufio.NewScanner(file)
		for buf.Scan() {
			// 获取到每一行的内容
			line := buf.Text()

			var d interface{}
			err := json.Unmarshal([]byte(line), &d)
			if err != nil {
				// 存在读取错误的一行
				log.Fatal(fmt.Sprint("第", index, "读取错误"))
				continue
			}

			dataArr = append(dataArr, d)
			index++
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "文件读取成功",
			"data": dataArr,
		})
	})

	// 获取数据，分页的方式  固定使用 data.jsonl的数据
	router.GET("/getPage", func(ctx *gin.Context) {
		// 获取当前页，默认 1 页， 每页10条数据
		currentPage := ctx.DefaultQuery("curPage", "1")
		curPage, _ := strconv.Atoi(currentPage)

		// 获取文件的路径
		filePath := "data.jsonl"

		// 打开文件
		file, err := os.Open(path.Join("./data", filePath))
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"code": http.StatusBadGateway,
				"msg":  fmt.Sprint("读取文件错误,err:", err),
			})
			return
		}
		defer file.Close()

		// 储存文件中的数据
		var dataArr []interface{} = make([]interface{}, 0)
		// 读取的所有
		index := 1

		// 足行扫描文件
		buf := bufio.NewScanner(file)
		for buf.Scan() {
			// 获取到每一行的内容
			line := buf.Text()

			var d interface{}
			err := json.Unmarshal([]byte(line), &d)
			if err != nil {
				// 存在读取错误的一行
				log.Fatal(fmt.Sprint("第", index, "读取错误"))
				continue
			}

			dataArr = append(dataArr, d)
			index++
		}

		// 获得页面的总数
		total := math.Ceil(float64(index) / 10)

		start := (curPage - 1) * 10
		end := curPage * 10

		var data []interface{}
		if end <= index {
			// 获取分页的数据
			data = dataArr[start:end]
		} else {
			data = dataArr[start:]
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "文件读取成功",
			"data": gin.H{
				"data":    data,
				"curPage": curPage,
				"total":   total,
			},
		})
	})

	// 启动gin
	router.Run()
}

// 解决gin跨域问题
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
