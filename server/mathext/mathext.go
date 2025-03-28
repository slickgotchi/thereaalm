package mathext

func Abs[T int | int8 | int16 | int32 | int64](x T) T {
    if x < 0 {
        return -x
    }
    return x
}

func Clamp[T int | int8 | int16 | int32 | int64 | float32 | float64](x, min, max T) T {
    if x < min {
        return min
    }
    if x > max {
        return max
    }
    return x
}

func Min[T int | int8 | int16 | int32 | int64 | float32 | float64](a, b T) T {
    if a < b {
        return a
    }
    return b
}

func Max[T int | int8 | int16 | int32 | int64 | float32 | float64](a, b T) T {
    if a > b {
        return a
    }
    return b
}

func Sign[T int | int8 | int16 | int32 | int64 | float32 | float64](x T) int {
    switch {
    case x > 0:
        return 1
    case x < 0:
        return -1
    default:
        return 0
    }
}

func Lerp[T float32 | float64](a, b, t T) T {
    return a + (b-a)*t
}

func Wrap[T int](value, min, max T) T {
    rangeSize := max - min
    for value < min {
        value += rangeSize
    }
    for value >= max {
        value -= rangeSize
    }
    return value
}
