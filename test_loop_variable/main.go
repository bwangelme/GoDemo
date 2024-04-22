package main

import "fmt"

func MultiMatchUserProfileAds(ads []int, uid uint64, userProfileTags []int32) []int {
	result := make([]int, 0)
	for _, ad := range ads {
		if userProfileTags == nil {
			fmt.Println("get userProfileTags")
			userProfileTags = []int32{}
			result = append(result, ad)
		}
	}

	return result
}

func main() {
	MultiMatchUserProfileAds([]int{2, 3, 4, 5, 7}, 23, nil)
}
