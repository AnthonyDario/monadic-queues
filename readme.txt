Anthony Dario

    This is an exploration of using the Monad pattern to make managing mutable
state in a distributed setting easier for programmers.  In an attempt to drink
from the same functional programming well as mapreduce, the idea is to use
concepts originally developed for functional programming to aid in distributed
systems design, development, and operation.

    Monads are used pervasively in functional programming for a variety of
different purposes.  One of the main uses is as a technique for managing side
effects and mutable state.  Haskell, for example, hides all IO inside of an IO
monad.

    This project will wrap a message broker (rabbitMQ) with a set of monadic
interfaces.  It will then implement a simple pizza ordering website in an
overly complicated distributed system style in order to evaluate the monadic
interface in terms of programming simplicity as well as performance.
