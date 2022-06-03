package employee

import "testing"

type fakeStmtForMaleCount struct {
	Stmt
}

func (f fakeStmtForMaleCount) Exec(stmt string, args ...string) (Result, error) {
	return Result{Count: 5}, nil
}

func TestEmployeeMaleCount(t *testing.T) {
	f := fakeStmtForMaleCount{}
	c, _ := FemaleCount(f)
	if c != 5 {
		t.Errorf("want: %d, actual: %d", 5, c)
		return
	}
}

type fakeStmt2MaleCount struct {
	Stmt2
}

func (f fakeStmt2MaleCount) Exec(s string, args ...string) (Result2, error) {
	return Result2{count: 5}, nil
}

func TestMaleCount(t *testing.T) {
	f := fakeStmt2MaleCount{}
	r, err := MaleCount(f)
	if err != nil {
		t.Log(err)
		return
	}
	want := 5
	if r != want {
		t.Logf("want: %d, got: %d\n", want, r)
	}
}
