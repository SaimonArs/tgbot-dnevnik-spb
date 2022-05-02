package utilslice

func ReversInt(a []int) []int {
    for j, i := 0, len(a) - 1; i > j; j, i = j+1, i-1 {
        a[i], a[j] = a[j], a[i]
    }
    return a
}
