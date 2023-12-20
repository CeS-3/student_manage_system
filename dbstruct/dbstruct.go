package dbstruct
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)
//学生结构
type Student struct{
	Sno string 
	Sname sql.NullString
	Ssex sql.NullString
	Sage sql.NullInt32
	Sdept sql.NullString
	Scholarship sql.NullString
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
func UpdateStudentInformation(db *sql.DB,updateStudent Student) error{
	// 执行数据库更新操作，更新指定 ID 的学生信息
	//构造查询语句
	query := "UPDATE Course SET Sname = ?,Ssex = ?,Sage = ?,Sdept = ?,Scholarship = ? WHERE Sno = ?"
	//执行语句
	_,err := db.Exec(query,updateStudent.Sname,updateStudent.Ssex,updateStudent.Sage,updateStudent.Sdept,updateStudent.Scholarship,updateStudent.Sno)	
	if err != nil{
		return err
	}
	//执行成功输出
	fmt.Printf("Updated Course with ID %s\n", updateStudent.Sno)
	return nil
}
//删除学生信息
func DeleteStudentFromDB(db *sql.DB,studentID string) error{
	//执行数据库删除操作，删除指定 ID 的SC信息
	qurey1 := "DELETE FROM SC WHERE Sno = ?"
	_,err := db.Exec(qurey1,studentID)
	if err != nil{
		return err
	}
	//同时也要删除掉学生表中对应的内容
	qurey2 := "DELETE FROM Student WHERE Sno = ?"
	_,err = db.Exec(qurey2,studentID)
	if err != nil{
		return err
	}
	//执行成功输出
	fmt.Printf("Delete Student with ID %s\n",studentID)
	return nil
}

//成绩
type Grade struct{
	Sname sql.NullString
	Sno string
	Cname sql.NullString
	Cno int
	Ggrade sql.NullInt32
}
//成绩插入操作
func (grade Grade) AddNewGrade(db *sql.DB) error{
	query := "INSERT INTO SC (Sno, Cno, Grade) VALUES(?,?,?)" 
	
	_,err := db.Exec(query,grade.Sno,grade.Cno,grade.Ggrade)
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
	SELECT
	Student.Sname,Student.Sno,Course.Cname,Course.Cno,SC.Grade 
	FROM Student,Course,SC 
	WHERE SC.Sno=Student.Sno 
	AND SC.Cno=Course.Cno `)
	if err != nil{
		return nil,err
	}
	defer rows.Close()
	//用于存储成绩
	var grades []Grade
	for rows.Next(){
        var grade Grade
        if err := rows.Scan(&grade.Sname,&grade.Sno,&grade.Cname,&grade.Cno,&grade.Ggrade); err != nil {
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
func UpdateGradeInformation(db *sql.DB, updatedGrade Grade) error{
	query := "UPDATE SC SET Grade = ? WHERE Sno = ? AND Cno = ?"
	_, err := db.Exec(query,updatedGrade.Ggrade,updatedGrade.Sno,updatedGrade.Cno)
	if err != nil{
		return err
	}
	fmt.Println("成绩修改成功")
	return nil
}
//搜索成绩
func SearchGrade(db *sql.DB, studentID string)([]Grade,error){
	query := `SELECT 
	Student.Sname,Student.Sno,Course.Cname,Course.Cno,SC.Grade 
	FROM Student,Course,SC 
	WHERE SC.Sno=Student.Sno 
	AND SC.Cno=Course.Cno
	AND Student.Sno=? 
	`
	rows,err := db.Query(query,studentID)
	if err != nil{
		return nil,err
	}
	defer rows.Close()
	var SearchResults []Grade
	for rows.Next(){
		var SearchResult Grade
		if err := rows.Scan(&SearchResult.Sname,&SearchResult.Sno,&SearchResult.Cname,&SearchResult.Cno,&SearchResult.Ggrade); err != nil{
			return nil,err
		} 
		SearchResults = append(SearchResults, SearchResult)
	}
	if err := rows.Err(); err != nil{
		return nil,err
	}
	return SearchResults,nil  
}
//成绩特征结构
type GradeAttribution struct{
	Sdept string 
	Avg sql.NullFloat64
	Max sql.NullInt32
	Min sql.NullInt32
	Erate sql.NullFloat64  //优秀率
	Failers sql.NullInt32	 //不及格人数
}
func GetAllGradesAttribution(db *sql.DB) ([]GradeAttribution,error){
	rows, err := db.Query(`SELECT
	Student.Sdept, 
    AVG(SC.Grade), 
    MAX(SC.Grade), 
    MIN(SC.Grade), 
    SUM(CASE WHEN Grade >= 90 THEN 1 ELSE 0 END) / COUNT(*) * 100, 
    SUM(CASE WHEN Grade < 60 THEN 1 ELSE 0 END)
	FROM SC,Student
	WHERE SC.Sno = Student.Sno
	GROUP BY Student.Sdept`)
	if err != nil{
		return nil,err
	}
	defer rows.Close()
	var attributions []GradeAttribution
	for rows.Next(){
		var attribution GradeAttribution
		if err := rows.Scan(&attribution.Sdept,&attribution.Avg,&attribution.Max,&attribution.Min,&attribution.Erate,&attribution.Failers);err != nil{
			return nil,err
		}
		attributions = append(attributions,attribution)
	}
	if err := rows.Err(); err != nil{
		return nil,err
	}
	return attributions,err
}
//排名结构
type Rank struct{
	Rrank sql.NullInt32
	Sdept sql.NullString
	Sno string
	Sname sql.NullString
}
func GetAllRanks(db *sql.DB) ([]Rank,error){
	rows, err := db.Query(`SELECT
    RANK() OVER (PARTITION BY s.Sdept ORDER BY SUM(sc.Grade) DESC) AS Ranking,
    s.Sdept AS Department,
    s.Sno AS StudentID,
    s.Sname AS StudentName
FROM
    Student s, SC sc
WHERE
    s.Sno = sc.Sno
GROUP BY
    s.Sdept, s.Sno, s.Sname
ORDER BY
    Department, Ranking;
`)
	if err != nil{
		return nil,err
	}
	defer rows.Close()
	var ranks []Rank
	for rows.Next(){
		var rank Rank
		if err := rows.Scan(&rank.Rrank,&rank.Sdept,&rank.Sno,&rank.Sname);err != nil{
			return nil,err
		}
		ranks = append(ranks,rank)
	}
	if err := rows.Err(); err != nil{
		return nil,err
	}
	return ranks,err
}
//课程结构
type Course struct{
	Cno int
	Cname sql.NullString
	Cpno sql.NullInt32
	Ccredit sql.NullInt32
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
    rows, err := db.Query("SELECT Cno, Cname, Cpno, Ccredit FROM Course")
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
func UpdateCourseInformation(db *sql.DB, updatedCourse Course) error {
	// 执行数据库更新操作，更新指定 ID 的课程信息
	//构造查询语句
	query := "UPDATE Course SET Cname = ?,Cpno = ?,Ccredit = ? WHERE Cno = ?"
	//执行语句
	_,err := db.Exec(query,updatedCourse.Cname,updatedCourse.Cpno,updatedCourse.Ccredit,updatedCourse.Cno)	
	if err != nil{
		return err
	}
	//执行成功输出
	fmt.Printf("Updated Course with ID %d\n", updatedCourse.Cno)
	return nil
}
//删除课程信息
func DeleteCourseFromDB(db *sql.DB, courseID string) error{
	//执行数据库删除操作，删除指定 ID 的课程信息
	qurey1 := "DELETE FROM SC WHERE Cno = ?"
	_,err := db.Exec(qurey1,courseID)
	if err != nil{
		return err
	}
	//同时也要删除掉SC表中对应的内容
	qurey2 := "DELETE FROM Course WHERE Cno = ?"
	_,err = db.Exec(qurey2,courseID)
	if err != nil{
		return err
	}
	//执行成功输出
	fmt.Printf("Delete Course with ID %s\n",courseID)
	return nil
}

