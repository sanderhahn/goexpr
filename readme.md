# Goexpr

Expression evaluation using [Shunting-yard algorithm](http://en.wikipedia.org/wiki/Shunting-yard_algorithm).

```
> a=(1+2)*3/4
ast: (= a (/ (* (() (+ 1 2)) 3) 4))
2.25
```

The parsing is done by building a grammar structure.
Semantics are added by surrounding `actions` with a callback.
Backtracking stops after a succesful alternative is found.
The expression parser constructs the ast tree that is visited for evaluation. 