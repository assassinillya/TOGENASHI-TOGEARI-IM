package set

// Union 求并集
func Union[T uint | int | string | uint32](slice1, slice2 []T) []T {
	m := make(map[T]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// Intersect 求交集
func Intersect[T uint | int | string | uint32](slice1, slice2 []T) []T {
	m := make(map[T]int)
	nn := make([]T, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// Difference 求差集 slice1-并集
func Difference[T uint | int | string | uint32](slice1, slice2 []T) []T {
	m := make(map[T]int)
	nn := make([]T, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}
