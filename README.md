# go-valve-query

A [babashka pod](https://github.com/babashka/babashka.pods) for running [a2s](https://developer.valvesoftware.com/wiki/Server_queries) queries and [rcon](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol) commands to valve source engine servers.

Implemented using the Go [go-source-server-query](https://github.com/NewPage-Community/go-source-server-query) and [transit](https://github.com/russolsen/transit) libraries.

This repository was adapted/forked from the [pod-babashka-go-sqlite3](https://github.com/babashka/pod-babashka-go-sqlite3) by @borkdude.

## Usage

Load the pod and `tommy-mor.go-valve-query` namespace:

``` clojure
(ns valve-script
  (:require [babashka.pods :as pods]))

(pods/load-pod "TODO update when I add to registry")
(require '[tommy-mor.go-valve-query :as steam])
```

TODO update these examples..

The namespace exposes two functions: `execute!` and `query`. Both accept a path
to the sqlite database and a query vector:

``` clojure
(sqlite/execute! "/tmp/foo.db"
  ["create table if not exists foo (the_text TEXT, the_int INTEGER, the_real REAL, the_blob BLOB)"])

;; This pod also supports storing blobs, so lets store a picture.
(def png (java.nio.file.Files/readAllBytes (.toPath (io/file "resources/babashka.png"))))

(sqlite/execute! "/tmp/foo.db"
  ["insert into foo (the_text, the_int, the_real, the_blob) values (?,?,?,?)" "foo" 1 3.14 png])
;;=> {:rows-affected 1, :last-inserted-id 1}

(def results (sqlite/query "/tmp/foo.db" ["select * from foo order by the_int asc"]))
(count results) ;;=> 1

(def row (first results))
(keys row) ;;=> (:the_text :the_int :the_real :the_blob)
(:the_text row) ;;=> "foo"

;; Should be true:
(= (count png) (count (:the_blob row)))
```

Additionally, unparameterised queries are supported if a string is passed
```clojure
(sqlite/query "/tmp/foo.db" "select * from foo")
```

Passing any other kind of data apart from a string or a vector will throw.

See [test/script.clj](test/script.clj) for an example test script.

## Build

### Requirements

- [Go](https://golang.org/dl/) 1.15+ should be installed.
- Clone this repo.
- Run `go build` to compile the binary.

## License

Copyright © 2020-2021 Michiel Borkent and Rahul De TODO? not sure what do to for this...

License: [BSD 3-Clause](https://opensource.org/licenses/BSD-3-Clause)
