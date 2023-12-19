package dbstruct
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)
//学生结构
type Student struct{
	Sno int
	Sname string
	Ssex string
	Sage int
	Sdept string
	Scholarship string
}
//学生表的插入操作
func (student Student) insert(db *sql.DB) error{
	query := "INSERT INTO student (Sno, Sname, Ssex, Sage, Sdept, Scholarship) VALUES (?, ?, ?, ?, ?, ?)"

	// 执行插入操作
	_, err := db.Exec(query, student.Sno, student.Sname, student.Ssex, student.Sage, student.Sdept, student.Scholarship)
	if err != nil {
		return err
	}

	fmt.Println("学生表插入成功")
	return nil
}
// //学生表的删除操作
// func (student Student) delete(db *sql.DB) error{
	
// }
// //学生表的修改操作
// func (student Student) update(db *sql.DB) error{
	
// }
// //学生表的搜索操作
// func (student Student) search(db *sql.DB) error{
	
// }
//学生-课程结构
type SC struct{
	Sno int
	Cno int
	Grade int
}
//课程结构
type Course struct{
	Cno int
	Cname string
	Cpno int
	Ccredit int	
}
//获取表中所有的课程信息
func getAllCourses(db *sql.DB) ([]Course, error) {
	// 查询数据库中的课程表
    rows, err := db.Query("SELECT Cno, Cname, Cpno, Ccredit FROM courses")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var courses []Course

    // 遍历查询结果，映射为 Course 结构体
    for rows.Next() {
        var course Course
        if err := rows.Scan(&course.Cno, &course.Cname, &course.Cpno, &course.Ccredit); err != nil {
            return nil, err
        }
        courses = append(courses, course)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return courses, nil
}
//修改课程信息
func updateCourseInformation(db *sql.DB, courseID string, updatedCourse Course) error {
	// 执行数据库更新操作，更新指定 ID 的课程信息
	
	
	// 假设这里执行一个更新的示例操作
	fmt.Printf("Updated course with ID %s\n", courseID)
	return nil
}
//删除课程信息
func deleteCourseFromDB(db *sql.DB, courseID string) error{
	//执行数据库删除操作，删除指定 ID 的课程信息
	
}
