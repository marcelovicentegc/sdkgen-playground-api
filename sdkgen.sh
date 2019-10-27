#!/bin/bash
crystal run ../sdkgen/main.cr -- src/playground.sdkgen -o src/gen/playground.$1 -t $2