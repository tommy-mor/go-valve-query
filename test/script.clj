#!/usr/bin/env bb

(ns script
  (:require [babashka.pods :as pods]
            [clojure.java.io :as io]
            [clojure.test :as t :refer [deftest is testing]]))

(pods/load-pod "./go-valve-query")
(require '[tommy-mor.go-valve-query :as steam])

(prn (ns-publics 'tommy-mor.go-valve-query))

(def url "mge.tf:27015")

(prn (steam/ping url))

(prn (steam/info url))
(prn (steam/players url))

(->> (steam/rcon url (clojure.string/trim (slurp "password.txt"))
                 "sm")
     clojure.string/split-lines
     (map println))


