package structs

type TestSimple struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Category    string `db:"category"`
	DiffLevel   int    `db:"diff_level"`
	Picture     []byte `db:"picturefile"`
}

type TestQuestionDTO struct {
	ID       int    `db:"id"`
	TestID   int    `db:"test_id"`
	Question string `db:"question"`
	IsSong   bool   `db:"is_song"`
	SongFile string `db:"song_file"`
}

type TestFull struct {
	ID          int
	Name        string
	Description string
	Category    string
	DiffLevel   int
	Questions   []TestQuestionFull
}

type TestQuestionFull struct {
	ID       int
	TestID   int
	Question string
	IsSong   bool
	Song     []byte
	Answers  []QuestionAnswer
}

type QuestionAnswer struct {
	ID         int    `db:"id"`
	QuestionID int    `db:"question_id"`
	Answer     string `db:"answer"`
	IsCorrect  bool   `db:"is_correct"`
}

type UserAnswer struct {
	QuestionID int
	AnswerID   int
}

type TestResult struct {
	User      string
	TestID    int
	Total     int
	UserScore int
}

type UserScore struct {
	TestID    int
	UserLogin string
	Score     int
}

type UserScoreDTO struct {
	TestID    int    `db:"test_id"`
	UserLogin string `db:"user_login"`
	Score     int    `db:"score"`
}
