# Contributing

Contributions are more than welcome! Just try to follow the conventions listed
below.

### Conventions

#####File Structure

The declaration order is as follows:

1. Vars
2. Types
    1. Vars
    3. Interfaces
    4. Structs
    5. Functions
3. Exported functions
4. Unexported functions

##### Tests

Code must pass `gofmt` and `gotype`. It should also pass the test suite, unless
it is a breaking change.

Tests should be written with a BDD approach. If something doesn't make sense to
unit test, but you want to test the behaviour in some way, `reflect` is your
best friend.