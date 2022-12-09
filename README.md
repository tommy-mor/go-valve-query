# go-valve-query

A [babashka pod](https://github.com/babashka/babashka.pods) for running [a2s](https://developer.valvesoftware.com/wiki/Server_queries) queries and [rcon](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol) commands to valve source engine servers.

Implemented using the Go [go-source-server-query](https://github.com/NewPage-Community/go-source-server-query) and [transit](https://github.com/russolsen/transit) libraries.

This repository was adapted/forked from the [pod-babashka-go-sqlite3](https://github.com/babashka/pod-babashka-go-sqlite3) by @borkdude.

## Usage

Load the pod and `tommy-mor.go-valve-query` namespace:

``` clojure
(ns valve-script
  (:require [babashka.pods :as pods]))

(pods/load-pod 'tommy-mor.go-valve-query "0.1.0")
(require '[tommy-mor.go-valve-query :as steam])
```

The namespace exposes four functions: `ping`, `info`, `players`, and `rcon`. The first argument is always the url of the server you want to query:

``` clojure
(steam/ping "mge.tf:27015")
;;=> {:ms 81}

(steam/info "mge.tf:27015")
;;=> {:game-id 440, :sourcetv-port 0, :max-players 24, :protocol 17, :game "Team Fortress", :folder "tf", :name "Team Fortress", :bots 0, :port 27015, :keywords "", :steam-id 85568392925064336, :id 440, :players 0, :environment 1, :server-type 1, :version "7708610", :sourcetv-name "", :vac 2, :map "mge_chillypunch_final4_fix2", :visibility 1}

(steam/players "mge.tf:27015")
;;=>[{:index 0, :name "JJP", :score 16, :duration 2280.06689453125}]

(steam/rcon "mge.tf:27015" "<password>" "sm help")
;;=> ["SourceMod Menu:" "Usage: sm <command> [arguments]" "    cmds             - List console commands" "    config           - Set core configuration options" "    credits          - Display credits listing" "    cvars            - View convars created by a plugin" "    exts             - Manage extensions" "    plugins          - Manage Plugins" "    prof             - Profiling" "    version          - Display version information"]
```

## Build

### Requirements

- [Go](https://golang.org/dl/) 1.15+ should be installed.
- Clone this repo.
- Run `go build -o go-valve-query` to compile the binary.

## License
This is a fork/adaptation of https://github.com/babashka/pod-babashka-go-sqlite3, which had this license:

Copyright © 2020-2021 Michiel Borkent and Rahul De

License: [BSD 3-Clause](https://opensource.org/licenses/BSD-3-Clause)
