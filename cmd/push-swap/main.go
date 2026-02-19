package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// --- 1. THE BRAIN (TURK ALGORITHM LOGIC) ---

type MovePlan struct {
	Shared       int
	RemainingA   int
	RemainingB   int
	MoveA, MoveB string
	SharedType   string
	TotalCost    int
}

func GetBestPlan(stackA, stackB []int) MovePlan {
	var bestPlan MovePlan
	bestPlan.TotalCost = math.MaxInt
	lenA, lenB := len(stackA), len(stackB)

	for i, valA := range stackA {
		targetIdxB := findTargetIdxB(valA, stackB)
		distA, isATop := getDist(i, lenA)
		distB, isBTop := getDist(targetIdxB, lenB)

		plan := calculateStrategy(distA, distB, isATop, isBTop)
		if plan.TotalCost < bestPlan.TotalCost {
			bestPlan = plan
		}
	}
	return bestPlan
}

func calculateStrategy(distA, distB int, isATop, isBTop bool) MovePlan {
	p := MovePlan{}
	if isATop && isBTop {
		p.Shared = min(distA, distB)
		p.SharedType = "rr"
		p.RemainingA, p.RemainingB = distA-p.Shared, distB-p.Shared
		p.MoveA, p.MoveB = "ra", "rb"
	} else if !isATop && !isBTop {
		p.Shared = min(distA, distB)
		p.SharedType = "rrr"
		p.RemainingA, p.RemainingB = distA-p.Shared, distB-p.Shared
		p.MoveA, p.MoveB = "rra", "rrb"
	} else {
		p.RemainingA, p.RemainingB = distA, distB
		if isATop { p.MoveA = "ra" } else { p.MoveA = "rra" }
		if isBTop { p.MoveB = "rb" } else { p.MoveB = "rrb" }
	}
	p.TotalCost = p.Shared + p.RemainingA + p.RemainingB
	return p
}

func findTargetIdxB(valA int, stackB []int) int {
	targetIdx := -1
	closestSmaller := math.MinInt64
	for i, valB := range stackB {
		if valB < valA && valB > closestSmaller {
			closestSmaller = valB
			targetIdx = i
		}
	}
	if targetIdx == -1 {
		maxVal := math.MinInt64
		for i, valB := range stackB {
			if int64(valB) > int64(maxVal) {
				maxVal = valB
				targetIdx = i
			}
		}
	}
	return targetIdx
}

func findTargetIdxA(valB int, stackA []int) int {
	targetIdx := -1
	closestBigger := math.MaxInt64
	for i, valA := range stackA {
		if valA > valB && valA < closestBigger {
			closestBigger = valA
			targetIdx = i
		}
	}
	if targetIdx == -1 {
		minVal := math.MaxInt64
		for i, valA := range stackA {
			if valA < minVal {
				minVal = valA
				targetIdx = i
			}
		}
	}
	return targetIdx
}

// --- 2. THE OPERATIONS (THE BODY) ---

func pa(a *[]int, b *[]int) {
	if len(*b) == 0 { return }
	val := (*b)[0]
	*b = (*b)[1:]
	*a = append([]int{val}, (*a)...)
	fmt.Println("pa")
}

func pb(a *[]int, b *[]int) {
	if len(*a) == 0 { return }
	val := (*a)[0]
	*a = (*a)[1:]
	*b = append([]int{val}, (*b)...)
	fmt.Println("pb")
}

func ra(s []int, label string) {
	if len(s) < 2 { return }
	first := s[0]
	copy(s, s[1:])
	s[len(s)-1] = first
	fmt.Println(label)
}

func rra(s []int, label string) {
	if len(s) < 2 { return }
	last := s[len(s)-1]
	copy(s[1:], s)
	s[0] = last
	fmt.Println(label)
}

func rotateSilent(s []int) {
	if len(s) < 2 { return }
	f := s[0]; copy(s, s[1:]); s[len(s)-1] = f
}

func revRotateSilent(s []int) {
	if len(s) < 2 { return }
	l := s[len(s)-1]; copy(s[1:], s); s[0] = l
}

// --- 3. THE EXECUTORS ---

func ExecuteMove(p MovePlan, a, b *[]int) {
	for i := 0; i < p.Shared; i++ {
		if p.SharedType == "rr" {
			rotateSilent(*a); rotateSilent(*b); fmt.Println("rr")
		} else {
			revRotateSilent(*a); revRotateSilent(*b); fmt.Println("rrr")
		}
	}
	for i := 0; i < p.RemainingA; i++ {
		if p.MoveA == "ra" { ra(*a, "ra") } else { rra(*a, "rra") }
	}
	for i := 0; i < p.RemainingB; i++ {
		if p.MoveB == "rb" { ra(*b, "rb") } else { rra(*b, "rrb") }
	}
	pb(a, b)
}

func sortThree(a []int) {
	max := a[0]
	for _, v := range a { if v > max { max = v } }
	if a[0] == max { ra(a, "ra") } else if a[1] == max { rra(a, "rra") }
	if a[0] > a[1] {
		a[0], a[1] = a[1], a[0]
		fmt.Println("sa")
	}
}

func finalizeA(a []int) {
	minIdx := 0
	minVal := a[0]
	for i, val := range a {
		if val < minVal { minVal = val; minIdx = i }
	}
	dist, isTop := getDist(minIdx, len(a))
	for i := 0; i < dist; i++ {
		if isTop { ra(a, "ra") } else { rra(a, "rra") }
	}
}

// --- 4. MAIN FLOW ---

func main() {
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

	// Phase 1: Push A to B until 3 left
	if len(stackA) > 3 { pb(&stackA, &stackB) }
	if len(stackA) > 3 { pb(&stackA, &stackB) }
	for len(stackA) > 3 {
		plan := GetBestPlan(stackA, stackB)
		ExecuteMove(plan, &stackA, &stackB)
	}

	// Phase 2: Sort Three
	sortThree(stackA)

	// Phase 3: Push B back to A
	for len(stackB) > 0 {
		targetIdxA := findTargetIdxA(stackB[0], stackA)
		dist, isTop := getDist(targetIdxA, len(stackA))
		for i := 0; i < dist; i++ {
			if isTop { ra(stackA, "ra") } else { rra(stackA, "rra") }
		}
		pa(&stackA, &stackB)
	}

	// Phase 4: Align A
	finalizeA(stackA)
}

func min(a, b int) int { if a < b { return a }; return b }
func getDist(idx, length int) (int, bool) {
	if idx <= length/2 { return idx, true }
	return length - idx, false
}