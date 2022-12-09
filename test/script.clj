#!/usr/bin/env bb

(ns script
  (:require [babashka.pods :as pods]
            [clojure.java.io :as io]
            [clojure.test :as t :refer [deftest is testing]]))

(pods/load-pod "./go-valve-query")
(require '[tommy-mor.go-valve-query :as steam])

(prn (ns-publics 'tommy-mor.go-valve-query))

#_(prn (steam/ping "elo2.sappho.io:27215" "query"))

(prn (steam/info "elo2.sappho.io:27215" "query"))


