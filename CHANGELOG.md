# Changelog

## [0.2.0](https://github.com/opentasks/cmd/compare/v0.1.0...v0.2.0) (2025-12-09)


### Features

* **graph:** implement graph builder and visualization for task relationships ([2363d9f](https://github.com/opentasks/cmd/commit/2363d9f1561849fe14402d6f25d04731c6fb88cc))
* **query:** implement in-memory SQLite query acceleration with relationship normalization ([e193768](https://github.com/opentasks/cmd/commit/e193768a3fb61066e9272d3006969b3bdfd46e03))
* remove graph feature entirely ([0815fd9](https://github.com/opentasks/cmd/commit/0815fd9ad80f7892ab7d2cc632ba691c5dfca7c7))

## [0.1.0](https://github.com/opentasks/cmd/compare/v0.0.2-alpha.1...v0.1.0) (2025-12-07)


### Features

* **config:** enable automatic prerelease detection in goreleaser ([f7aed7d](https://github.com/opentasks/cmd/commit/f7aed7dc726e9c61f9483ff4f5845110734023f5))
* implement active project resolution with E2E testing ([#24](https://github.com/opentasks/cmd/issues/24)) ([93c221a](https://github.com/opentasks/cmd/commit/93c221adf0830b153677d752ca965d89e62bf038))


### Bug Fixes

* **brand:** use FontSlant for banner rendering ([2c241a4](https://github.com/opentasks/cmd/commit/2c241a465fb8773aa6987cd1c65a7b6606b568db))
* calculate prerelease tag and remove unsupported GoReleaser flag ([efd3299](https://github.com/opentasks/cmd/commit/efd3299f8711bcb4ed86cf29e37b7f128cafd79c))
* correct prerelease detection and tag handling in publish workflow ([28d3896](https://github.com/opentasks/cmd/commit/28d389642c73e7e4232950ded7dd3e917b0f16c1))
* hk pre commit  ([#9](https://github.com/opentasks/cmd/issues/9)) ([418fc7d](https://github.com/opentasks/cmd/commit/418fc7dcde7040ebf016ac9b863d236f342e56e3))
* **security:** address gosec findings (12 → 0 failing issues) ([#12](https://github.com/opentasks/cmd/issues/12)) ([b8dabcf](https://github.com/opentasks/cmd/commit/b8dabcf762fb3721efb755707644e76b63d3b09d))
* **security:** correct gosec/govulncheck hook configuration ([#11](https://github.com/opentasks/cmd/issues/11)) ([c85c41f](https://github.com/opentasks/cmd/commit/c85c41f7354882f18e3f519c03a720affab20024))
* use toJSON() in release debug step to output raw outputs object ([5a89ef5](https://github.com/opentasks/cmd/commit/5a89ef53e3f59715ec136024743ee43cc76a0422))


### Miscellaneous Chores

* release 0.1.0 ([6b2d9cd](https://github.com/opentasks/cmd/commit/6b2d9cd74ee70c3549a6727bc2741e9afbe974af))

## [0.0.2-alpha.1](https://github.com/opentasks/cmd/compare/v0.0.1-alpha.1...v0.0.2-alpha.1) (2025-11-26)


### Bug Fixes

* **release:** update repository references from zenobi-us/opentask to opentasks/cmd ([#8](https://github.com/opentasks/cmd/issues/8)) ([f1848d1](https://github.com/opentasks/cmd/commit/f1848d17adbccac85f9697eb1804f0dc9a15600f))
