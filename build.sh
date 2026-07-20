#!/bin/bash

go build -ldflags "\
    -X 'darkroom/pkg/config.DarkroomSecret=$DARKROOM_SECRET' \
    -X 'darkroom/pkg/config.EncryptionKey=$ENCRYPTION_KEY'" -o darkroom
