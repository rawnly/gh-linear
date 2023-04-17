# GH Linear

> Create a new branch from a linear issue

## Installation

```
  gh extension install rawnly/gh-linear
```

Upgrade:

```
  gh extension upgrade rawnly/gh-linear
```

## Usage

```
gh-linear is a tool to help you create new branches from Linear issues

Usage:
  gh-linear [flags]

Examples:
$ gh linear --issue <IDENTIFIER>
$ gh linear


Flags:
  -h, --help           help for gh-linear
  -i, --issue string   The issue identifier
```

## Setup

Right now in order to run `gh linear` you must have the environment variable `LINEAR_API_KEY` correctly configured.

### Suggestion

Create a `.env` file with the correct `LINEAR_API_KEY` in your projects and run `source .env` before running `gh linear`
