package employee

type Result struct {
	Count int
}

func (r Result) Int() int { return r.Count }

type Rows []struct{}

type Stmt interface {
	Close() error
	NumInput() int
	Exec(stmt string, args ...string) (Result, error)
	Query(args []string) (Rows, error)
}

// 返回女性员工总数
func FemaleCount(s Stmt) (int, error) {
	result, err := s.Exec("select count(*) from employee_tab where gender=?", "1")
	if err != nil {
		return 0, err
	}

	return result.Int(), nil
}

type Result2 struct {
	count int
}

func (r Result2) Count() int {
	return r.count
}

type Rows2 []struct{}

type Stmt2 interface {
	Exec(s string, args ...string) (Result2, error)
	Close() error
	NumInput() int
	Query(args []string) (Rows2, error)
}

func MaleCount(s Stmt2) (int, error) {
	r, err := s.Exec("select count(*) from male where gender=?", "1")
	if err != nil {
		return 0, err
	}
	return r.Count(), nil
}
