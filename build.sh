#!/bin/bash

go build -ldflags "-X 'darkroom/pkg/config.DarkroomSecret=$DARKROOM_SECRET'" -o darkroom
