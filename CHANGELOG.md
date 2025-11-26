# Changelog

## 1.0.0 (2025-11-26)


### ⚠ BREAKING CHANGES

* ID format updated
    - Old: s-42, e-1, p-3a (type-specific sequential with collision suffix)
    - New: 42, 1, 3 (global sequential integer)
    - Type now stored in frontmatter only

### Features

* Add config hierarchical discovery, merging, init and view commands ([b16d947](https://github.com/opentasks/cmd/commit/b16d947136663204bd08988c43d6caf3abe614cd))
* Add project context support for ergonomic project selection ([0cbea5a](https://github.com/opentasks/cmd/commit/0cbea5a67f5d92cc1577e83c94530aff2401a8af))
* Add tree view of resolved config files ([9bb6a8f](https://github.com/opentasks/cmd/commit/9bb6a8f64982017efb18cdf66a65929cf62d0271))
* **cli:** integrate new config resolution system into CLI commands ([f825893](https://github.com/opentasks/cmd/commit/f8258936b8e288d5c1268951947f95e592724784))
* **config:** implement new schema types and merging logic for global and project configs ([462b563](https://github.com/opentasks/cmd/commit/462b56314d1e9c668c0ed1b32dc57d57d77d8485))
* implement core opentasks MVP ([306109e](https://github.com/opentasks/cmd/commit/306109e7da06a4b6782face46db9ed9a41ac435c))
* Implement project context CLI commands ([6297991](https://github.com/opentasks/cmd/commit/6297991611b78017d5bb13effa267f9884dac564))
* Implement project context CLI commands ([752fc2a](https://github.com/opentasks/cmd/commit/752fc2aaa0bda8a1f50e250bbd50e8f8b7fbc4ea))
* show virtual defaults layer in config view verbose output ([5a86d6a](https://github.com/opentasks/cmd/commit/5a86d6a528435b978e36000071aa936424793129))


### Bug Fixes

* build to bin ([3767132](https://github.com/opentasks/cmd/commit/3767132dfbc92cdf3e0fe8bae88d9cf9650df753))
* **ci:** simplify publish workflow and fix tag resolution for all triggers ([d7fd66d](https://github.com/opentasks/cmd/commit/d7fd66dc6057d16affd2a811fe36f18a6dd01d97))
* **cli:** improve config file path display in config view ([a188b77](https://github.com/opentasks/cmd/commit/a188b7701fed424ff299c01638e04b9f7978de31))
* **config:** eliminate duplicate global config in discovered files ([30f7c14](https://github.com/opentasks/cmd/commit/30f7c143fd6b3b3e8dd7bdc3cf1a49b29470b252))
* **config:** merge storage path from matching global project ([92074b9](https://github.com/opentasks/cmd/commit/92074b9d341ee1c73313ac29ff4d28fddf7db81c))
* **config:** merge templates from matching global project ([081f992](https://github.com/opentasks/cmd/commit/081f99251abe4c44c71980e6f32c732beab811d5))
* display proper directory tree structure in config view file listing ([e1ed08f](https://github.com/opentasks/cmd/commit/e1ed08f6f6659ad87f724a8ebaa93a0bb70f071a))
* improve config file tree display with proper hierarchy and defaults ([1a325bf](https://github.com/opentasks/cmd/commit/1a325bff93a216a1207b9dbad21bd9f2821b1551))
* move config show as alias to existing cmd ([3f45e4e](https://github.com/opentasks/cmd/commit/3f45e4eacc29660673e394f4b743b589b0b03b57))
* Prevent panic in config view when no config files found ([a5b6802](https://github.com/opentasks/cmd/commit/a5b6802cc1a534310016c87daffd7dbc6632e5ab))
* restore proper YAML frontmatter to design decision file ([1b79728](https://github.com/opentasks/cmd/commit/1b797280bc05c0b10e035a75eaee8c5051fd98a3))
* **test:** account for user global config in discovery tests ([96926bb](https://github.com/opentasks/cmd/commit/96926bb1db47a947885bd6c8c0dfde1fef8ded9b))
* use hierarchical config discovery in CLI instead of looking for single config.toml ([23b9c7d](https://github.com/opentasks/cmd/commit/23b9c7d596b5b7dbcbfc9d50e564c4372dc043a4))


### Code Refactoring

* update ID system from per-type semantic to global sequential ([bccffc1](https://github.com/opentasks/cmd/commit/bccffc18c116072d612fdb260c03c27c28681e1e))
