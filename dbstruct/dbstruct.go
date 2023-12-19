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
func (student Student) AddNewStudent(db *sql.DB) error{
	query := "INSERT INTO Student (Sno, Sname, Ssex, Sage, Sdept, Scholarship) VALUES (?, ?, ?, ?, ?, ?)"

	// 执行插入操作
	_, err := db.Exec(query, student.Sno, student.Sname, student.Ssex, student.Sage, student.Sdept, student.Scholarship)
	if err != nil {
		return err
	}

	fmt.Println("学生表插入成功")
	return nil
}
//查询整张学生表
func GetAllStudents(db *sql.DB) ([]Student,error){
	//进行查询
	rows, err := db.Query("SELECT Sno,Sname,Ssex,Sage,Sdept,Scholarship FROM Student")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    //用于存储查询结果
	var students []Student
	//遍历查询结果
	for rows.Next() {
        var student Student
        if err := rows.Scan(&student.Sno,&student.Sname,&student.Ssex,&student.Sage,&student.Sdept,&student.Scholarship); err != nil {
            return nil, err
        }
        students = append(students, student)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }
		
	return students,nil
}
//更新学生信息
func UpdateStudentInformation(db *sql.DB,studentID string,updateStudent Student) error{
	// 执行数据库更新操作，更新指定 ID 的学生信息
	//构造查询语句
	query := "UPDATE Course SET Sname = ?,Ssex = ?,Sage = ?,Sdept = ?,Scholarship = ? WHERE Sno = ?"
	//执行语句
	_,err := db.Exec(query,updateStudent.Sname,updateStudent.Ssex,updateStudent.Sage,updateStudent.Sdept,updateStudent.Scholarship,studentID)	
	if err != nil{
		return err
	}
	//执行成功输出
	fmt.Printf("Updated Course with ID %s\n", studentID)
	return nil
}
//删除学生信息
func DeleteStudentFromDB(db *sql.DB,studentID string) error{
	//执行数据库删除操作，删除指定 ID 的课程信息
	qurey1 := "DELETE FROM Student WHERE Sno = ?"
	_,err := db.Exec(qurey1,studentID)
	if err != nil{
		return err
	}
	//同时也要删除掉SC表中对应的内容
	qurey2 := "DELETE FROM SC WHERE Sno = ?"
	_,err = db.Exec(qurey2,studentID)
	if err != nil{
		return err
	}
	//执行成功输出
	fmt.Printf("Delete Student with ID %s\n",studentID)
	return nil
}

//学生-课程结构

type Grade struct{
	Rank int
	Sname string
	Sno int
	Cname string
	Cno int
	_grade int
}
type GradeAttribution struct{
	Sdept string
	avg float64
	max int
	min int
	Erate float64  //优秀率
	failers int	 //不及格人数
}
//成绩插入操作
func (grade Grade) AddNewGrade(db *sql.DB) error{
	query := "INSERT INTO SC (Sno, Cno, Grade) VALUES(?,?,?)" 
	
	_,err := db.Exec(query,grade.Sno,grade.Cno,grade._grade)
	if err != nil{
		return err
	}

	fmt.Println("成绩插入成功")
	return nil
}
//查询所有学生的成绩
func GetAllOrderedGrades(db *sql.DB) ([]Grade,error){
	//查询学生的姓名，对应的课程名，对应的成绩
	rows,err := db.Query(`
	SELECT ROW_NUMBER() OVER (ORDER BY SC.Grade DESC) AS Rank, 
	Student.Sname,Student.Sno,Course.Cname,Course.Cno,SC.Grade 
	FROM Student,Course,SC 
	WHERE SC.Sno=Student.Sno 
	AND SC.Cno=Course.Cno 
	ORDER BY SC.Grade DESC`)
	if err != nil{
		return nil,err
	}
	defer rows.Close()
	//用于存储成绩
	var grades []Grade
	for rows.Next(){
        var grade Grade
        if err := rows.Scan(&grade.Rank,&grade.Sname,&grade.Sno,&grade.Cname,&grade.Cno,&grade._grade); err != nil {
            return nil, err
        }
        grades = append(grades, grade)
	}

	if err := rows.Err(); err != nil {
        return nil, err
    }

	return grades,nil
}
//修改成绩
func UpdateGradeInformation(db, studentID, courseID, updatedGrade){

}
//课程结构
type Course struct{
	Cno int
	Cname string
	Cpno int
	Ccredit int	
}
//课程表的插入操作
func (course Course) AddNewCourse(db *sql.DB) error{
	query := "INSERT INTO Course (Cno,Cname,Cpno,Ccredit) VALUES (?,?,?,?)"

	// 执行插入操作
	_, err := db.Exec(query,course.Cno,course.Cname,course.Cpno,course.Ccredit)
	if err != nil {
		return err
	}

	fmt.Println("课程表插入成功")
	return nil
}
//获取表中所有的课程信息
func GetAllCourses(db *sql.DB) ([]Course, error) {
	// 查询数据库中的课程表
    rows, err := db.Query("SELECT Cno, Cname, Cpno, Ccredit FROM Courses")
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
func UpdateCourseInformation(db *sql.DB, courseID string, updatedCourse Course) error {
	// 执行数据库更新操作，更新指定 ID 的课程信息
	//构造查询语句
	query := "UPDATE Course SET Cname = ?,Cpno = ?,Ccredit = ? WHERE Cno = ?"
	//执行语句
	_,err := db.Exec(query,updatedCourse.Cname,updatedCourse.Cpno,updatedCourse.Ccredit,courseID)	
	if err != nil{
		return err
	}
	//执行成功输出
	fmt.Printf("Updated Course with ID %s\n", courseID)
	return nil
}
//删除课程信息
func DeleteCourseFromDB(db *sql.DB, courseID string) error{
	//执行数据库删除操作，删除指定 ID 的课程信息
	qurey1 := "DELETE FROM Course WHERE Cno = ?"
	_,err := db.Exec(qurey1,courseID)
	if err != nil{
		return err
	}
	//同时也要删除掉SC表中对应的内容
	qurey2 := "DELETE FROM SC WHERE Cno = ?"
	_,err = db.Exec(qurey2,courseID)
	if err != nil{
		return err
	}
	//执行成功输出
	fmt.Printf("Delete Course with ID %s\n",courseID)
	return nil
}

