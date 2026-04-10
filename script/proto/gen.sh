#!/bin/bash

set -e

cd $(dirname $0)/..

echo "Generating RPC code from proto..."

# Generate user service
kitex -module vicomova -service user -gen-path third_party/kitex_gen ./api/rpc/user/user.proto

echo "Done!"
