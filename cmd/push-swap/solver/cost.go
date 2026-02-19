package solver

import "math"

type MovePlan struct {
    AIdx, BIdx int
    Shared     int    // Number of rr or rrr
    RemainingA int    // Remaining ra/rra
    RemainingB int    // Remaining rb/rrb
    MoveA      string // "ra" or "rra"
    MoveB      string // "rb" or "rrb"
    SharedType string // "rr" or "rrr"
    TotalCost  int
}

func GetBestPlan(stackA, stackB []int) MovePlan {
    var bestPlan MovePlan
    bestPlan.TotalCost = math.MaxInt
    lenA, lenB := len(stackA), len(stackB)

    for i, valA := range stackA {
        targetIdxB := findTargetIdxB(valA, stackB)
        currentPlan := calculatePlan(i, targetIdxB, lenA, lenB)

        if currentPlan.TotalCost < bestPlan.TotalCost {
            bestPlan = currentPlan
        }
    }
    return bestPlan
}

func calculatePlan(idxA, idxB, lenA, lenB int) MovePlan {
    plan := MovePlan{AIdx: idxA, BIdx: idxB}
    
    // Determine directions
    isATop := idxA <= lenA/2
    isBTop := idxB <= lenB/2

    // Calculate individual raw distances
    distA := idxA
    if !isATop { distA = lenA - idxA }
    distB := idxB
    if !isBTop { distB = lenB - idxB }

    if isATop && isBTop {
        // Shared "Rotate"
        plan.Shared = min(distA, distB)
        plan.SharedType = "rr"
        plan.RemainingA = distA - plan.Shared
        plan.RemainingB = distB - plan.Shared
        plan.MoveA, plan.MoveB = "ra", "rb"
        plan.TotalCost = plan.Shared + plan.RemainingA + plan.RemainingB
    } else if !isATop && !isBTop {
        // Shared "Reverse Rotate"
        plan.Shared = min(distA, distB)
        plan.SharedType = "rrr"
        plan.RemainingA = distA - plan.Shared
        plan.RemainingB = distB - plan.Shared
        plan.MoveA, plan.MoveB = "rra", "rrb"
        plan.TotalCost = plan.Shared + plan.RemainingA + plan.RemainingB
    } else {
        // Mixed - No shared moves
        plan.RemainingA, plan.RemainingB = distA, distB
        plan.TotalCost = distA + distB
        if isATop { plan.MoveA = "ra" } else { plan.MoveA = "rra" }
        if isBTop { plan.MoveB = "rb" } else { plan.MoveB = "rrb" }
    }
    return plan
}

func findTargetIdxB(valA int, stackB []int) int {
	targetIdx := -1
	closestSmaller := math.MinInt64 // Use 64-bit to be safe

	// 1. Try to find the closest value that is smaller than valA
	for i, valB := range stackB {
		if valB < valA && valB > int(closestSmaller) {
			closestSmaller = valB
			targetIdx = i
		}
	}

	// 2. Edge Case: If no smaller number exists, target is the Max value in B
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