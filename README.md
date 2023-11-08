# dynamic
Dual view of types, static and dynamic, in go

Tired of:
* having to juggle between nested `map[string]any` and your own incomplete types?
* keep your static types up to date with source implementation in an API you are consuming?

You probably don't have to.

Meet `dynamic.T[...]`.

## Usage

replace:
```
var result MyTyp
err := json.Unmarshall(bytes, &result)
```

with:
```
var result dynamic.T[MyTyp]
err := json.Unmarshall(bytes, &result)
```

And you're done!

Now you have access to the static fields by `.Static` and `.Extra`.

Of course, you can also use `dynamic.T[..]` types nested inside other `dynamic.T[..]` types.

kthx, bye
