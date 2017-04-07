#!/usr/bin/env bash

# import.sh updates C and C++ dependencies from their upstream sources and
# applies CockroachDB patches. Usage:
#
#     ./import.sh [DEP...]
#
# If dependency specifications are omitted, all dependencies will be updated.
#
# To add a new dependency, add it to the `deps` array below. To update a
# dependency, update the URL with the new version in the `deps` array. Patch
# files in this directory of the form DEP-*.patch are automatically applied
# after the dependency is download. To perform further transformations, declare
# a niladic function named `mangle_DEP` in this file, which will be invoked
# after successful patching.

set -euo pipefail
shopt -s nullglob

((${BASH_VERSION%%.*} >= 4)) || {
  echo "fatal: bash 4 or later required. You have $BASH_VERSION." >&2
  exit 1
}

declare -A deps
deps=(
    [rocksdb]=https://github.com/facebook/rocksdb/archive/v5.1.4.tar.gz
    [jemalloc]=https://github.com/jemalloc/jemalloc/releases/download/4.5.0/jemalloc-4.5.0.tar.bz2
    [protobuf]=https://github.com/google/protobuf/releases/download/v3.2.0/protobuf-cpp-3.2.0.tar.gz
    [snappy]=https://github.com/google/snappy/releases/download/1.1.3/snappy-1.1.3.tar.gz
)

mangle_rocksdb() {
  # Downcase some windows-only includes for compatibility with mingw64.
  grep -lR '^#include <.*[A-Z].*>' rocksdb | while IFS= read -r source_file
  do
    echo "downcasing headers in $source_file"
    awk '/^#include <.*[A-Z].*>/ { print tolower($0); next; } { print; }' "$source_file" > tmp
    mv tmp "$source_file"
  done

  # Avoid MSVC-only extensions for compatibility with mingw64.
  grep -lRF 'i64;' rocksdb | xargs sed -i~ 's!i64;!LL;!g'

  rm -rf rocksdb/{arcanist_util,docs,java}
}

(($# >= 1)) && goals=("$@") || goals=("${!deps[@]}")

for dep in "${goals[@]}"
do
  [[ "${deps["$dep"]:-}" ]] || {
    echo "unrecognized dep $dep" >&2
    exit 1
  }
done

for dep in "${goals[@]}"
do
  echo "updating $dep"
  url="${deps[$dep]}"
  rm -rf "$dep"
  mkdir "$dep"
  curl -fL "$url" | tar -xz -C "$dep" --strip-components 1
  for patch in $dep-*.patch
  do
    patch -d "$dep" -p1 < "$patch"
  done
  [[ "$(type -t "mangle_$dep")" = function ]] && "mangle_$dep"
done

