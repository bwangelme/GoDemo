package bdreview

/*
数组 [1,2,3,4,5,6], target = 12

target int，选择3个数的和等于 target

选择不重复的 3 元组
*/

func solution(arr []int, target int) [][]int {
	var (
		res    = make([][]int, 0)
		subset = make([]int, 0)
	)
	helper(arr, 0, subset, target, &res)
	return res
}

/*
add 1

	|--add 2
		|-- add 3
		|-- nadd 3
	|--nadd 2
		|-- add 3
		|-- nadd 3

nadd 1

	|--add 2
		|-- add 3
		|-- nadd 3
	|--nadd 2
		|-- add 3
		|-- nadd 3
*/
func helper(arr []int, idx int, subset []int, target int, result *[][]int) {
	// 符合条件的结果
	if 0 == target && len(subset) == 3 {
		resItem := make([]int, len(subset))
		copy(resItem, subset)
		*result = append(*result, resItem)
	} else if target > 0 && idx < len(arr) && len(subset) < 3 {
		// 不添加当前元素
		helper(arr, idx+1, subset, target, result)

		// 添加当前元素
		subset = append(subset, arr[idx])
		helper(arr, idx+1, subset, target-arr[idx], result)
		subset = subset[0 : len(subset)-1]
	}
}
