#!/bin/bash

CWD=$(pwd)
CORE_DIR="$CWD/core"
DEVKIT_RELEASE_DIR="$CWD/devkit-release"

cd $CORE_DIR && \
    make devkit && \
    mkdir -p $DEVKIT_RELEASE_DIR && rm -rf $DEVKIT_RELEASE_DIR/* && \
    cp -r $CWD/devkit/* $DEVKIT_RELEASE_DIR && \
    mkdir $DEVKIT_RELEASE_DIR/core && \
    cp $CORE_DIR/core.so $DEVKIT_RELEASE_DIR/core && \
    cp $CORE_DIR/package.yml $DEVKIT_RELEASE_DIR/core && \
    cp -r $CORE_DIR/resources $DEVKIT_RELEASE_DIR/core/resources && \
    cp -r $CORE_DIR/sdk $DEVKIT_RELEASE_DIR/core/sdk && \
    cp -r $CWD/main $DEVKIT_RELEASE_DIR/main
