package main

import (
	//"fmt"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MovePlan struct {
	Shared       int
	RemainingA   int
	RemainingB   int
	MoveA, MoveB string
	SharedType   string
	TotalCost    int
}


func pa(a *[]int, b *[]int) {
	if len(*b) == 0 { return }
	val := (*b)[0]
	*b = (*b)[1:]
	*a = append([]int{val}, (*a)...)
	//fmt.Println("pa")
}

func pb(a *[]int, b *[]int) {
	if len(*a) == 0 { return }
	val := (*a)[0]
	*a = (*a)[1:]
	*b = append([]int{val}, (*b)...)
	//fmt.Println("pb")
}

func ra(s []int, label string) {
	if len(s) < 2 { return }
	first := s[0]
	copy(s, s[1:])
	s[len(s)-1] = first
	//fmt.Println(label)
}

func rra(s []int, label string) {
	if len(s) < 2 { return }
	last := s[len(s)-1]
	copy(s[1:], s)
	s[0] = last
	//fmt.Println(label)
}

func rotateSilent(s []int) {
	if len(s) < 2 { return }
	f := s[0]; copy(s, s[1:]); s[len(s)-1] = f
}

func revRotateSilent(s []int) {
	if len(s) < 2 { return }
	l := s[len(s)-1]; copy(s[1:], s); s[0] = l
}


func main() {
    inst := []string{"pb", "pb", "rb", "pb", "sa", "pa", "pa", "pa"}

	if len(os.Args) < 2 { return }
	

	// Parsing Input
	var stackA []int
	input := os.Args[1:]
	if len(input) == 1 { input = strings.Fields(input[0]) }
	for _, arg := range input {
		n, _ := strconv.Atoi(arg)
		stackA = append(stackA, n)
	}
	stackB := []int{}
	for _ , ins := range inst {
		switch ins {
		case"pb":
			pb(&stackA,&stackB)
		case"ra":
			ra(stackA,"ra")
		case"sa":
			stackA[0], stackA[1] = stackA[1], stackA[0]
		case"pa":
			pa(&stackA,&stackB)
		case"rb":
			ra(stackB, "rb")

		}
	}
	fmt.Println(stackA, stackB)
	

}