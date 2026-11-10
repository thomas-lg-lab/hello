# repo-template

Template for new `thomas-lg-lab` repositories: exactly the files every repository of the organisation carries,
ready from the first commit. Create a repository from it with **Use this template** on GitHub, or
`gh repo create thomas-lg-lab/<name> --template thomas-lg-lab/repo-template --public --clone`.

The organisation drives the rest from [`thomas-lg-lab/.github`](https://github.com/thomas-lg-lab/.github): issue
forms, PR template, SECURITY.md and code of conduct defaults, the shared workflows, the Release Drafter
configuration and the settings script.

## What you get

| File | Role |
| --- | --- |
| `.github/workflows/pull-request.yml` | Conventional Commit PR title (required check) and label, shared workflow |
| `.github/workflows/ci.yml` | The project's checks, then the shared image build: Trivy gate on PRs, `edge` and `sha-*` from main, `X.Y.Z` and `latest` from a release |
| `.github/workflows/release-drafter.yml` | Draft release from the merged PRs; a manual run publishes it and builds the release image |
| `.github/workflows/registry-cleanup.yml` | Weekly pruning of old `sha-*` images |
| `.github/workflows/scorecard.yml` | OpenSSF Scorecard (Security tab and badge) |
| `.github/rulesets/` | `Protect main` (squash only, linear history, signed commits, required checks, CodeQL results, no bypass) and `Code owner review` (changes under `.github/`) |
| `.github/CODEOWNERS`, `renovate.json`, `LICENSE`, `SECURITY.md`, `AGENTS.md`, `.editorconfig`, `.gitattributes` | Per repository by nature; `renovate.json` extends the shared preset of `thomas-lg-lab/renovate-config` |

## After creating a repository from it

1. Replace the `checks` job of `ci.yml` with the project's build, tests and lint; list their job names in `needs`
   of `image` and in `required_status_checks` of `.github/rulesets/main.json`. Without an image, delete the
   `image` job, `image / build` from the ruleset and `registry-cleanup.yml`.
2. From the root of the new clone, with the organisation's `.github` repository cloned next to it:
   `../.github/scripts/apply-standards.sh thomas-lg-lab/<name>` (settings, security features, CodeQL default
   setup, labels, rulesets). Run it again after editing a ruleset file.
3. Write the local check commands in a CONTRIBUTING.md (the organisation's PR template points there), and the
   project's specifics in `AGENTS.md`. Renovate needs nothing: the organisation's App covers every repository.
4. Replace this README, keep the badge:

```markdown
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/thomas-lg-lab/<name>/badge)](https://scorecard.dev/viewer/?uri=github.com/thomas-lg-lab/<name>)
```
