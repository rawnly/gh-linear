# GH Linear

> Create a new branch from a linear issue

## Installation

```
  gh extensions install rawnly/gh-linear
```

Upgrade:

```
  gh extensions upgrade rawnly/gh-linear
```

## Usage

```
gh-linear is a tool to help you create new branches from Linear issues

Usage:
  gh-linear [flags]
  gh-linear [command]

Examples:
$ gh linear --issue <IDENTIFIER>
$ gh linear


Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  project     Manage projects

Flags:
  -h, --help           help for gh-linear
  -i, --issue string   The issue identifier

Use "gh-linear [command] --help" for more information about a command.
```

## Setup

Right now in order to run `gh linear` you must have the environment variable `LINEAR_API_KEY` correctly configured.

### Suggestion

Create a `.env` file with the correct `LINEAR_API_KEY` in your projects and run `source .env` before running `gh linear`
