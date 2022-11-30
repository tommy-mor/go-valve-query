#!/usr/bin/env bb

(ns script
  (:require [babashka.pods :as pods]
            [clojure.java.io :as io]
            [clojure.test :as t :refer [deftest is testing]]))

(prn (pods/load-pod "./go-valve-query"))

(require '[tommy-mor.go-valve-query :as steam])

(prn (steam/connect "mge.tf:27015" "query"))


