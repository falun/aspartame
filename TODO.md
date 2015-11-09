##### TODO
* `ParenExpr`s break parsing of assigned enum values
* Seems dangerous to always overwirte files if it already exists (but necessary to regenerate; dig into what stringer does)
* If the const value aggregates from imported files the generated file isn't going compile -- potentially needs `-includeImports a,b,c` flag
* Not sure if there is a good reason to do so but: what would it take to support non-int enum values?
* There is probably a reason I should not use my home grown `Strings()` think about switching to `stringer` or emulating it
