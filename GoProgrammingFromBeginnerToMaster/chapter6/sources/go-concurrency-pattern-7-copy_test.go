package main

import (
	"fmt"
	"testing"
	"time"
)

func ShutdownerMaker(tm int) func(timeout time.Duration) error {
	return func(timeout time.Duration) error {
		defer fmt.Printf("ShutdownerMaker done, cost: %ds\n", tm)
		time.Sleep(time.Duration(tm) * time.Second)
		return nil
	}
}

func TestConcurrentShutdown2(t *testing.T) {
	type args struct {
		timeout     time.Duration
		shutdowners []Shutdowner2
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"base case", args{5 * time.Second, []Shutdowner2{ShutdownFunc(ShutdownerMaker(3)), ShutdownFunc(ShutdownerMaker(2))}}, false},
		{"timeout case", args{5 * time.Second, []Shutdowner2{ShutdownFunc(ShutdownerMaker(2)), ShutdownFunc(ShutdownerMaker(6))}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConcurrentShutdown2(tt.args.timeout, tt.args.shutdowners...); (err != nil) != tt.wantErr {
				t.Errorf("ConcurrentShutdown2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSequentialShutdown2(t *testing.T) {
	type args struct {
		timeout    time.Duration
		shutdowner []Shutdowner2
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"base case",
			args{5 * time.Second, []Shutdowner2{ShutdownFunc(ShutdownerMaker(2)), ShutdownFunc(ShutdownerMaker(2))}},
			false},
		{"timeout case",
			args{5 * time.Second, []Shutdowner2{ShutdownFunc(ShutdownerMaker(2)), ShutdownFunc(ShutdownerMaker(4))}},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SequentialShutdown2(tt.args.timeout, tt.args.shutdowner...); (err != nil) != tt.wantErr {
				t.Errorf("SequentialShutdown2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
