package models

const (
	// SPRING_24_MID = "24_SPRING_MID"
	// SPRING_24_FIN = "24_SPRING_FINAL"

	SUMMER_24_MID = "24_SUMMER_MID"
)

/*
	Semister Time
	From January 1 to May 31 is Spring
	From June 1 to September 31 is Summer
	From October 1 to December 31 is Fall

	write a go function that takes a date as input and returns the semister in this format: 24_SPRING_MID


*/

type Exam struct {
	ID          string `json:"id" gorm:"primaryKey;size:64"`
	TrimsterID  string `json:"trimster_id" gorm:"index:trimster"`
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
