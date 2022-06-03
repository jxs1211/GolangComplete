package main

import "fmt"

type FinishedHouse struct {
	style                  int    // 0: Chinese, 1: American, 2: European
	centralAirConditioning bool   // true or false
	floorMaterial          string // "ground-tile" or ”wood"
	wallMaterial           string // "latex" or "paper" or "diatom-mud"
}

type Option func(*FinishedHouse)

func NewFinishedHouse(options ...Option) *FinishedHouse {
	h := &FinishedHouse{
		// default options
		style:                  0,
		centralAirConditioning: true,
		floorMaterial:          "wood",
		wallMaterial:           "paper",
	}

	for _, option := range options {
		option(h)
	}

	return h
}

func WithStyle(style int) Option {
	return func(h *FinishedHouse) {
		h.style = style
	}
}

func WithFloorMaterial(material string) Option {
	return func(h *FinishedHouse) {
		h.floorMaterial = material
	}
}

func WithWallMaterial(material string) Option {
	return func(h *FinishedHouse) {
		h.wallMaterial = material
	}
}

func WithCentralAirConditioning(centralAirConditioning bool) Option {
	return func(h *FinishedHouse) {
		h.centralAirConditioning = centralAirConditioning
	}
}

func showFunctionalOption() {
	fmt.Printf("%+v\n", NewFinishedHouse()) // use default options
	fmt.Printf("%+v\n", NewFinishedHouse(WithStyle(1),
		WithFloorMaterial("ground-tile"),
		WithCentralAirConditioning(false)))

}

type Room struct {
	f1 bool
	f2 int
	f3 string
	f4 string
}

type Option2 func(*Room)

func NewRoom(option ...Option2) *Room {
	r := &Room{
		f1: true,
		f2: 1,
		f3: "f3 default",
		f4: "f4 default",
	}

	for _, op := range option {
		op(r)
	}
	return r
}

func WithF1(b bool) Option2 {
	return func(r *Room) {
		r.f1 = b
	}
}

func WithF2(i int) Option2 {
	return func(r *Room) {
		r.f2 = i
	}
}

func WithF3(s string) Option2 {
	return func(r *Room) {
		r.f3 = s
	}
}

func WithF4(s string) Option2 {
	return func(r *Room) {
		r.f4 = s
	}
}

func showFunctionalOption2() {
	r := NewRoom()
	fmt.Printf("%+v\n", r)
	r = NewRoom(
		WithF1(false),
		WithF2(2),
	)
	fmt.Printf("%+v\n", r)
	r = NewRoom(
		WithF1(false),
		WithF2(2),
		WithF3("f3 manual"),
		WithF4("f4 manual"),
	)
	fmt.Printf("%+v\n", r)
}

func main() {
	showFunctionalOption2()
}
