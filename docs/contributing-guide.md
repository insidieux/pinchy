# Contributing guide

- Create an issue
- Fork project
- Create your feature branch (git checkout -b issue-id)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin issue-id)
- Create new Pull Request

## Requirements

- Every pull request code **must be** covered by tests
- Every new source/registry **must be** documented at /docs

## Adding new source/provider

- Add new sub directory to `/pkg/core/{source,registry}`
- Put your implementation of interface
- Add new sub directory to `/internal/extension/{source,registry}`
- Put wiring code for dynamic load implementations:
  - Register implementation of `ProviderInterface` at `init` function
  - Add anonymous import to `/cmd/pinchy/modules.go`
