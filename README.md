# BigTech.fail

BigTech.fail is an open-source research project aiming to highlight the dangers
of censorship, propaganda and surveillance from today's technology corporations
and governments.

* [Website](./website/)
* [Server](./server/)
* [Documentation](./doc/)
* [Graphics](./graphics/)

## Getting Started

### Development

With [`Hugo 0.83.1`](https://github.com/gohugoio/hugo/releases/tag/v0.83.1)
installed, simply run `hugo server -D` from the `./website` directory. This
will build and serve the website (with live reloading) from
`http://localhost:1313`.

### Production

With both [`Hugo 0.83.1`](https://github.com/gohugoio/hugo/releases/tag/v0.83.1)
and [`Go 1.16`](https://golang.org/doc/go1.16) installed, simply run `make`
from the root of the repository to build the full production binary (available
at `./server/bigtechfaild`) which will contain the static web resources
embedded within it. To build and run production locally on
`http://localhost:8080`, run `make run`.

## Licensing

* Please review the [CREDITS](./CREDITS.md) file for attribution given to work that is not ours.
* Any **source code** in this repository is free and unencumbered software released into the public domain. See the UNLICENSE file or visit [unlicense.org](https://unlicense.org/) for more information.
* Unless otherwise stated, any orignal content (articles, graphics, etc.) in this repo is released under the [Creative Commons Attribution 4.0 International License](https://creativecommons.org/licenses/by/4.0/).
* Please review our [Contributing Guidelines](./CONTRIBUTING.md) before you consider making a contribution.
