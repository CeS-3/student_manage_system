package main

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"student_manage_system/dbstruct"
	"github.com/gin-gonic/gin"
	"net/http"
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
	//展示学生信息
	r.GET("/students", func(c *gin.Context) {
		// 获取数据库中所有学生信息
		students, err := dbstruct.GetAllStudents(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve students"})
			return
		}

		// 返回学生信息给前端
		c.JSON(http.StatusOK, students)
	})
	//添加新生
	r.POST("/students/add", func(c *gin.Context) {
		// 获取用户提交的新生信息
		var newCourse dbstruct.Student
		if err := c.ShouldBindJSON(&newCourse); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// 执行插入操作
		err := newCourse.AddNewStudent(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert course"})
			return
		}

		// 返回成功信息给前端
		c.JSON(http.StatusOK, gin.H{"message": "Course added successfully"})
	})
	//更新学生信息
	r.POST("/students/:id/edit",func(c *gin.Context) {
		//从前端获取更新的学生信息
		var updateStudent dbstruct.Student
		if err := c.ShouldBind(&updateStudent);err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		}
		err = dbstruct.UpdateStudentInformation(db,updateStudent)
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error": "Failed to update Student"})
		}
	})
	//删除学生信息
	r.DELETE("/students/:id", func(c *gin.Context){
		studentID := c.Param("id")
		// 执行删除操作
		err := dbstruct.DeleteStudentFromDB(db, studentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
			return
		}
		// 返回成功信息给前端
		c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
	})
	//展示所有课程信息
	r.GET("/courses", func(c *gin.Context) {
		// 查询数据库中的所有课程
		courses, err := dbstruct.GetAllCourses(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"错误": "获取课程失败"})
			return
		}
		// 返回课程列表给前端
		c.JSON(http.StatusOK, courses)
	})
	//添加课程
	r.POST("/courses/add", func(c *gin.Context) {
		// 获取用户提交的新增课程信息
		var newCourse dbstruct.Course
		if err := c.ShouldBindJSON(&newCourse); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// 执行插入操作
		err := newCourse.AddNewCourse(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert course"})
			return
		}

		// 返回成功信息给前端
		c.JSON(http.StatusOK, gin.H{"message": "Course added successfully"})
	})
	//修改课程信息,一次只能修改一条
	r.POST("/courses/edit", func(c *gin.Context) {
		// 获取用户提交的修改后的课程信息
		var updatedCourse dbstruct.Course
		if err := c.ShouldBind(&updatedCourse); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// 执行更新操作
		err = dbstruct.UpdateCourseInformation(db, updatedCourse)
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
		err := dbstruct.DeleteCourseFromDB(db, courseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
			return
		}
		// 返回成功信息给前端
		c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
	})
	//成绩展示界面
	r.GET("/grades",func(c *gin.Context) {
		//获取排序后的成绩
		grades,err := dbstruct.GetAllOrderedGrades(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve grades"})
			return
		}
		//传递给前端页面
		c.JSON(http.StatusOK,grades)
	})
	//录入学生成绩
	r.POST("/grades/add", func(c *gin.Context) {
		// 获取用户提交的新增成绩信息
		var newGrade dbstruct.Grade
		if err := c.ShouldBind(&newGrade); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// 执行插入操作
		err := newGrade.AddNewGrade(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert grade"})
			return
		}

		// 返回成功信息给前端
		c.JSON(http.StatusOK, gin.H{"message": "Grade added successfully"})
	})
	//修改学生成绩
	r.POST("/grades/edit", func(c *gin.Context) {
		// 获取用户提交的修改后的成绩信息
		var updatedGrade dbstruct.Grade
		if err := c.ShouldBind(&updatedGrade); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// 执行更新操作
		err := dbstruct.UpdateGradeInformation(db, updatedGrade)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update grade"})
			return
		}
		// 返回成功信息给前端
		c.JSON(http.StatusOK, gin.H{"message": "Grade updated successfully"})
	})
	

}


