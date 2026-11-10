# Working on this repository (humans and agents)

- **Trunk-based.** `main` is the only long-lived branch. Work on a short-lived branch and open one
  PR into `main`; it is squash-merged once the required checks pass. Never push to `main`.
- **Signed commits.** The ruleset requires a verified signature on every commit of a PR, not only
  on the squash commit: sign with an SSH or GPG key registered on GitHub as a signing key. Commits
  made through the GitHub API or web UI (Renovate, the merge button) are signed by GitHub.
- **Merging is not releasing.** `main` builds the `edge` image; only a published release (a `v*`
  tag, from the draft release) builds `latest`.
- **Shared workflows.** PR conventions, release notes, the image build and registry cleanup call the reusable
  workflows of `thomas-lg-lab/.github` at a version tag; change them there, not here.
- **PR titles follow Conventional Commits** (`type(scope): summary`, commitlint types, `!` for a breaking
  change): checked by the required `conventions / title` check; the label and the release notes come from it.
- **Issues and PR descriptions** follow the organisation's defaults in `thomas-lg-lab/.github`
  (issue forms, PR template, and how to fill them in its README); they are not in this repository.
  Fill every section that applies, delete the others, and say how the change was tested. An issue
  opened through the API or `gh` skips the forms: use the same sections and set the issue type
  (`gh api repos/OWNER/REPO/issues -f title=... -f body=... -f type=Task`, or `Bug`, `Feature`).
- **Dependencies** are updated by Renovate (shared config in `thomas-lg-lab/renovate-config`); do not
  bump them by hand in a feature PR. Security fixes arrive as `security` PRs at any time.
- **Supply chain.** Actions are pinned to commit SHAs and images to digests; keep it that way when
  adding one (Renovate keeps them current).
- **Secrets** never go in the repository, in a commit message or in a log.
- Technical files are written in English.
