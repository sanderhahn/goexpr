# Goexpr

Expression evaluation using [Shunting-yard algorithm](http://en.wikipedia.org/wiki/Shunting-yard_algorithm).

```
> x = 6
6
> y = x + 2 * 10
26
```

The parsing is done by building a grammar structure.
Semantics are added by surrounding `actions` with a callback.
Backtracking stops after a succesful alternative is found.
