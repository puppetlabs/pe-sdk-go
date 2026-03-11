# pe-sdk-go

This project is consumed by pe-cli and if updates are made here,
[the version reference should be updated in pe-cli](https://github.com/puppetlabs/pe-cli/blob/5fb2be380a7f915ef485799244ba9452382faf3e/go.mod#L6)

Note: there is no CI that is versioning/tagging in this repository, so the
reference in pe-cli is a go pseudo-version.

A typical pseudo-version looks like v0.0.0-20180706040059-75cf19382434.

- v0.0.0: The base version.
- 20180706040059: The UTC timestamp (July 6, 2018, 4:00:59 AM UTC).
- 75cf19382434: The 12-character commit hash prefix.

```
git rev-parse --short=12 HEAD
b35f94e4a88e

git show -s --format='%ct %H' b35f94e4a88e

1773269054 b35f94e4a88e74bcccb5ff7bf34f01b3e1ad7403

date -u -r 1773269054 '+%Y%m%d%H%M%S'
20260311224414


v0.0.0-20260311224414-b35f94e4a88e
```