package models

type Exam struct {
	ID          string `json:"id" gorm:"primaryKey;size:64"`
	Department  string `json:"department" gorm:"index:dep"`
	CourseCode  string `json:"course_code" gorm:"index:code"`
	CourseTitle string `json:"course_title" gorm:"index:title"`
	Section     string `json:"section" gorm:"index:section"`
	Teacher     string `json:"teacher" gorm:"index:teacher"`
	ExamDate    string `json:"exam_date"`
	ExamTime    string `json:"exam_time"`
	Room        string `json:"room"`

	TimeCommon
}
