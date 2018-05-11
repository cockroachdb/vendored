// Copyright 2012 The LevelDB-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program creates a leveldb db at /tmp/db.
//
// To build and run:
// g++ make-db.cc -lleveldb && ./a.out

#include <iostream>

#include "leveldb/db.h"

static const char* dbname = "/tmp/db";

// The program consists of up to 4 stages. If stage is in the range [1, 4],
// the program will exit after the stage'th stage.
// 1. create an empty DB.
// 2. add some key/value pairs.
// 3. close and re-open the DB, which forces a compaction.
// 4. add some more key/value pairs.
static const int stage = 4;

int main(int argc, char** argv) {
  leveldb::Status status;
  leveldb::Options o;
  leveldb::WriteOptions wo;
  leveldb::DB* db;

  o.create_if_missing = true;
  o.error_if_exists = true;

  if (stage < 1) {
    return 0;
  }
  std::cout << "Stage 1" << std::endl;

  status = leveldb::DB::Open(o, dbname, &db);
  if (!status.ok()) {
    std::cerr << "DB::Open " << status.ToString() << std::endl;
    return 1;
  }

  if (stage < 2) {
    return 0;
  }
  std::cout << "Stage 2" << std::endl;

  status = db->Put(wo, "foo", "one");
  if (!status.ok()) {
    std::cerr << "DB::Put " << status.ToString() << std::endl;
    return 1;
  }

  status = db->Put(wo, "bar", "two");
  if (!status.ok()) {
    std::cerr << "DB::Put " << status.ToString() << std::endl;
    return 1;
  }

  status = db->Put(wo, "baz", "three");
  if (!status.ok()) {
    std::cerr << "DB::Put " << status.ToString() << std::endl;
    return 1;
  }

  status = db->Put(wo, "foo", "four");
  if (!status.ok()) {
    std::cerr << "DB::Put " << status.ToString() << std::endl;
    return 1;
  }

  status = db->Delete(wo, "bar");
  if (!status.ok()) {
    std::cerr << "DB::Delete " << status.ToString() << std::endl;
    return 1;
  }

  if (stage < 3) {
    return 0;
  }
  std::cout << "Stage 3" << std::endl;

  delete db;
  db = NULL;
  o.create_if_missing = false;
  o.error_if_exists = false;

  status = leveldb::DB::Open(o, dbname, &db);
  if (!status.ok()) {
    std::cerr << "DB::Open " << status.ToString() << std::endl;
    return 1;
  }

  if (stage < 4) {
    return 0;
  }
  std::cout << "Stage 4" << std::endl;

  status = db->Put(wo, "foo", "five");
  if (!status.ok()) {
    std::cerr << "DB::Put " << status.ToString() << std::endl;
    return 1;
  }

  status = db->Put(wo, "quux", "six");
  if (!status.ok()) {
    std::cerr << "DB::Put " << status.ToString() << std::endl;
    return 1;
  }

  status = db->Delete(wo, "baz");
  if (!status.ok()) {
    std::cerr << "DB::Delete " << status.ToString() << std::endl;
    return 1;
  }

  return 0;
}
