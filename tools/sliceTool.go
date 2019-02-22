package tools

import "strings"

//数组去重
func RemoveRepByLoop(slc []map[string](map[string]bool) )  []map[string](map[string]bool)  {
	result := []map[string](map[string]bool){}   // 存放结果
	for i := range slc{
		flag := true
		for j := range result{
			if equalMap(slc[i],result[j]){
				flag = false  // 存在重复元素，标识为false
				break
			}
		}
		if flag {  // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}


func equalMap(x, y map[string]map[string]bool) bool {
	if len(x) != len(y) {
		return false
	}
	for k, xv := range x {
		if yv, ok := y[k];!ok || !equelBool(xv,yv){
			return false
		}
	}
	return true
}

func equelBool(x,y map[string]bool) bool{
	if len(x) != len(y) {
		return false
	}
	for k,xv := range x{
		if yv, ok := y[k]; !ok || yv != xv {
			return false
		}
	}
	return true
}

//sample-values/9H200A1700008/#
//sample-values/9H200A1700008/
func CutOutString(s string)string{
	cutIndex := strings.Index(s, "#")
	s = string([]rune(s)[:cutIndex])
	return s
}