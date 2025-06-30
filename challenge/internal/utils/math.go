package utils

func PowerSet[T comparable](array []T) [][]T {
	return powerSetHelper(array, []T{})
}

func powerSetHelper[T comparable](array, power []T) (powerSt [][]T) {
	if len(power) == 0 {
		powerSt = append(powerSt, []T{})
	} else {
		powerSt = append(powerSt, power)
	}

	for i := range array {
		newPower := make([]T, 0)
		newPower = append(newPower, array[i])
		newPower = append(newPower, power...)

		powerSt = append(powerSt, powerSetHelper(array[:i], newPower)...)
	}
	return
}
