# Required Project Dependencies

- [Go](https://go.dev)
  - [ogen](https://ogen.dev) requires Go v1.18+

- [Task](https://taskfile.dev)
  - Like _Makefiles_ but made easier and for Go projects!

- [GoFmt](https://pkg.go.dev/cmd/gofmt)
  - Opinionated go formatter

- [Bun](https://bun.sh)
  - Fast Typescript/Javascript runtime.
  - Used along with _fluid-oas_ for generated OpenApiV3 specifications.

Typical development workflow:

1. Make changes in `api`, generate the openapi specification by using `task api:generate`

2. Regenerate handler layer in backend by using `task challenge:generate`

3. Implement handlers, business logic, etc.

## Helpful commands:

You can grab all existing commands

```bash
task --list-all
```

Generate OpenApi Spec
```bash
# Generate openapi specification
task api:generate
# Generate server code from openapi specification
task challenge:generate
```

```bash
# Run the server locally (Not in docker)
task challenge:run
```

```bash
# Run all tests including integration tests (Might take a while)
task challenge:test
```
