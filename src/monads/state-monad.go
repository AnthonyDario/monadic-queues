// The state monad allows mutable state to be modelled.  It will be used for
// the configuration, and perhaps the database.

package main

import (
    "log"
)

// The state monad contains a function that accepts a state, of type S, and
// returns a new state along with a return value, of type T.  
type StateMonad [S any, T any] struct {
    f func(S) (T, S)
}

// The monadic functions
// ---------------------------

// "Return" produces a value without modifying the state.
// Renamed to "produce" to avoid the keyword collision in go.
// This is the "unit" function
func produce [S any, T any] (v T) StateMonad[S, T] {
    return StateMonad[S, T] {
        func (s S) (T, S) {
            return v, s
        },
    }
}

// Bind modifies the monad, applying the given function, f, to its result
// bind : m a -> (a -> m b) -> mb
func bind [S any, T any, U any] (m StateMonad[S, T], f func(T) StateMonad[S, U]) StateMonad[S, U]{
    // The haskell code for this function is below
    // (>>=) :: State s t -> (t -> State s u) -> State s u  
    // m >>= f = \r -> let (x, s) = m r in (f x) s
    
    // This came out pretty cryptically but it is a
    // straight translation of the above haskell code
    return StateMonad[S, U] {
        func (r S) (U, S) {
            x, s := m.f(r)
            return f(x).f(s) 
        },
    }
}

// Get returns the current state in the value
func get [S any] () StateMonad[S, S] {
    return StateMonad[S, S] {
        func (s S) (S, S) {
            return s, s
        },
    }
}

// Put replaces the state, returning nothing
func put [S any, T any] (s S) StateMonad[S, T] {
    return StateMonad[S, T] {
        func (s2 S) (T, S) {
            var zedVal T // Cannot return nil for a type in go, so we return the "zero value"
            return zedVal, s
        },
    }
}

// Modify updates the state, returning nothing
func modify [S any, T any] (f func(S) S) StateMonad[S, T] {
    return StateMonad[S, T] { 
        func (s S) (T, S) {
            var zedVal T // Cannot return nil for a type in go, so we return the "zero value"
            return zedVal, f(s)
        },
    }
}

// RunState applies the monad to an initial state
// runState :: State s a -> s -> (a, s)
func runState [S any, T any] (m StateMonad[S, T]) StateMonad[S, T] {
    return StateMonad[S, T] {
        func (s S) (T, S) {
            return m.f(s)
        },
    }
}

// Runs a stateful computation on a given state and returns the final state
func execState [S any, T any] (m StateMonad[S, T], s S) S {
    _, state := runState(m).f(s)
    return state
}

func main () {
    // A simple example that increments an integer state

    // Our stateful computation increments the state by one
    var sm = modify[int, int](func (s int) (int) {
        return s + 1
    })

    // Now we apply the monad to an initial state (putting the state in the monad)
    finalState := execState(sm, 0)

    log.Print(finalState)
}
