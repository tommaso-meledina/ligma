# Project Brief

The purpose of this project is to build a simple and fast CLI tool named `ligma`.

The tool is meant to be used by OSS developers when they start a new project and need to copy-paste the appropriate license in the root of their new project.

The CLI will allow to:

- list the available licenses in the terminal (e.g. `ligma ls`)
- display a specific license in the terminal (e.g. `ligma get <license id>`)
- generate a `LICENSE.md` file containing a specific license (e.g. `ligma write <license id> [location]`)
- display a `help` message as usual

The list of licenses and the details of each license shall be taken, in JSON format, from the official SPDX GitHub repo at https://github.com/spdx/license-list-data/tree/main/json.

The CLI shall be written in Golang.