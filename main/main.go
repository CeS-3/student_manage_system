package main

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"student_manage_system/dbstruct"
)

//连接主机
var sql_host = os.Getenv("sql_host")
var sql_password = os.Getenv("sql_password")
//连接池对象
var db *sql.DB
//进行数据库的初始化
func db_init()(err error){
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/S_T_U202212057",sql_host,sql_password)
	//连接数据集
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("开启 %s 时发生错误:%v\n", dsn, err)
		return err
	}
	err = db.Ping() //尝试连接数据库
	if err != nil {
		fmt.Printf("开启 %s 时发生错误:%v\n", dsn, err)
		return err
	}
	return nil
}

func main(){
	//初始化数据库连接池
	err := db_init()
	if err != nil{
		fmt.Printf("数据库初始化时发生错误:%v\n",err)
		return
	}
	defer db.Close()
	//初始化gin引擎
	r := gin.Default()

	//展示所有课程信息
	r.GET("/courses", func(c *gin.Context) {
		// 查询数据库中的所有课程
		courses, err := dbstruct.getAllCourses(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"错误": "获取课程失败"})
			return
		}
		// 返回课程列表给前端
		c.JSON(http.StatusOK, courses)
	})
	//修改课程信息,一次只能修改一条
	r.POST("/courses/:id/edit", func(c *gin.Context) {
		// 从路由参数中获取课程 ID
		courseID := c.Param("id")
		// 获取用户提交的修改后的课程信息
		var updatedCourse dbstruct.Course
		if err := c.ShouldBind(&updatedCourse); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// 执行更新操作
		err := dbstruct.updateCourseInformation(db, courseID, updatedCourse)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
			return
		}
	})	
	//删除课程
	r.DELETE("/courses/:id", func(c *gin.Context) {
		// 从路由参数中获取课程 ID
		courseID := c.Param("id")

		// 执行删除操作
		err := dbstruct.deleteCourseFromDB(db, courseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
			return
		}
		// 返回成功信息给前端
		c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
	})
}
func ()

