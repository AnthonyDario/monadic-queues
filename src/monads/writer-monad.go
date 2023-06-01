package main

/*
import (
    "fmt"
    "log"
)
*/
// we need to have a writer queue monad

// Need a unit : a -> Writer a
// Need a bind : Writer a -> (a -> Writer b) -> Writer b

type Writer [T any] struct {
    Value T
    Log string
}

func unit [T any] (a T) Writer[T] {
    return Writer[T] {a, ""}
}

func bind [T any, U any] (w Writer[T], f func(T) Writer [U]) Writer[U] {
    var w2 = f(w.Value)
    return Writer[U] {w2.Value, w.Log + "\n" + w2.Log}
}

/*
func main () {
    // Build our writer with unit
    var w = unit(1) 

    // Our function from int -> Writer[bool]
    var f = func (i int) Writer[bool] {
        var isEven = i % 2 == 0
        var log string
        if isEven {
            log = fmt.Sprintf("%d is even", i)
        } else {
            log = fmt.Sprintf("%d is odd", i)
        }
        
        return Writer[bool] {i % 2 == 0, log}
    }

    var g = func (i int) Writer[int] {
        return Writer[int] {i + 1, "incremented i"}
    }
    
    var w2 = bind(w, g)
    var w3 = bind(w2, g)
    var w4 = bind(w3, g)
    var w5 = bind(w4, f)

    log.Print(w5.Log)
}
*/
