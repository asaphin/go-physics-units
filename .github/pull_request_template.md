<!-- What changes and why. Link the ROADMAP step if there is one. -->

## Checklist

- [ ] `make check` passes (tidy, format, vet, golangci-lint, race tests, allocation guards).
- [ ] `CHANGELOG.md` has an entry under `## Unreleased` for every user-visible change.
- [ ] Exported API changes are deliberate: the "API changes" CI job lists them, and the CHANGELOG mentions each one.
- [ ] Hot paths: `make bench-compare` shows no unexplained regression (calculations matter most; constructors may get slower).
- [ ] Tests and examples cover the change; `AGENTS.md`, `README.md` and `ROADMAP.md` are in step with it.
