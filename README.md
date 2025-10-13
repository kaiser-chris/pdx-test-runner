# Overview
**pdx-test-runner** is a tool to run scripted tests for games like Victoria 3 in a repeatable and configurable way.

![Title Icon](github_icon_readme.jpg)

## Contents
* [Configuration](#configuration)
* [Special Comments](#special-comments)
* [How To Build](#how-to-build)

## Status
[![Build Binaries](https://github.com/kaiser-chris/pdx-test-runner/actions/workflows/build.yaml/badge.svg)](https://github.com/kaiser-chris/pdx-test-runner/actions/workflows/build.yaml)
[![GitHub Release](https://img.shields.io/github/v/release/kaiser-chris/pdx-test-runner?display_name=release&label=Current%20Version&color=blue)](https://github.com/kaiser-chris/pdx-test-runner/releases)

## Configuration
```json
{
  "game-directory": "X:\\Path\\To\\Game\\Base\\Folder",
  "mod-directories": [
    "X:\\Path\\To\\First\\Mod\\In\\Load\\Order",
    "X:\\Path\\To\\Second\\Mod\\In\\Load\\Order"
  ],
  "ignored-files": [
    "some_test_file.txt",
    "another_test_file.txt"
  ]
}
```

## Special Comments
```
### name = Name of the whole test file
last_date = "1936.1.1"

tests = {
    ### name = Name of specific test
    ### desc = Description of specific test
	some_test = {
        acceptable_fail_rate = 0.0
		success = {
            always = yes
		}
		fail = {
            always = no
		}
	}
}
```

## How To Build
First download and install the Go SDK:
- https://go.dev/doc/install

Next, open the project folder in a terminal (e.g. cmd) and run the following command:
```
go build
```

That is it. There should be an executable in the project folder now.