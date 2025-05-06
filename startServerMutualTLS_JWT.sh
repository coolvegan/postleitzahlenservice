#!/bin/sh
AUTH=JWT JWTSECRET=dies-ist-mein-test-passwort-für-den-jwt  RPCTYPE=mTLS go run cmd/server/main.go