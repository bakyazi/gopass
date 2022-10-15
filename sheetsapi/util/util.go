package util

func Map[T, K any](s []T, fn func(T, int) K) []K {
	res := make([]K, len(s))
	for i, t := range s {
		res[i] = fn(t, i)
	}
	return res
}
