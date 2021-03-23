#!/usr/bin/env bash

git pull
go test -v -run TestListZpool pkg/zfs/zpool.go pkg/zfs/zfs.go pkg/zfs/zfs_test.go pkg/zfs/utils.go pkg/zfs/utils_notsolaris.go pkg/zfs/error.go