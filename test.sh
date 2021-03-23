#!/usr/bin/env bash

git pull
go test -v -run TestListZpool zpool.go zfs.go zfs_test.go utils.go utils_notsolaris.go error.go