package util

func RemoveFromSlice[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func Reverse[T any](s []T) []T {
    a := make([]T, len(s))
    copy(a, s)

    for i := len(a)/2 - 1; i >= 0; i-- {
        opp := len(a) - 1 - i
        a[i], a[opp] = a[opp], a[i]
    }

    return a
}