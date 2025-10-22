# Docker for Lethean

This document explains how to use the multi-stage `Dockerfile` to build the Lethean blockchain software.

## Overview

The `Dockerfile` is designed to provide a consistent and reproducible build environment. It uses a multi-stage build to create different artifacts: a build environment, build cache, compiled binaries, and a final service image.

You can target a specific stage of the build using the `--target` flag with the `docker build` command.

For example:
```bash
docker build --target=builder .
```

## Build Stages

There are four distinct build stages:

### 1. `builder`

This stage creates a development image that contains all the necessary dependencies and tools to compile the software. It can be used for development by mounting your local source code.

**To build the builder image:**
```bash
docker build --target=builder -t lethean/builder .
```

You can then use it for interactive development:
```bash
docker run -it -v .:/code lethean/builder /bin/bash
```

### 2. `build-cache`

This stage produces an image containing only the Conan package cache. This can be useful for speeding up subsequent builds by pre-populating the cache.

**To create the cache image:**
```bash
docker build --target=build-cache -t lethean/build-cache .
```

### 3. `build-artifacts`

This stage produces an image containing only the compiled binaries.

**To create the artifacts image:**
```bash
docker build --target=build-artifacts -t lethean/build-artifacts .
```

You can extract the binaries from this image:
```bash
docker create --name artifacts lethean/build-artifacts
docker cp artifacts:/ /path/to/local/binaries
docker rm artifacts
```

### 4. `chain-service`

This is the final stage, which produces a lean image containing the Lethean blockchain node and its runtime dependencies.

**To build the service image:**
```bash
docker build --target=chain-service -t lethean/chain-service .
```
or simply:
```bash
docker build -t lethean/chain-service .
```

## Build Arguments

The build process can be customized using `docker build --build-arg` flags.

| Argument        | Default Value                                  | Description                                                                                                     |
|-----------------|------------------------------------------------|-----------------------------------------------------------------------------------------------------------------|
| `THREADS`       | `1`                                            | Number of parallel threads to use for compilation.                                                              |
| `BUILD_BRANCH`  | `dev-configure-testnet-1`                      | The git branch to clone and build if `BUILD_LOCAL=0`.                                                           |
| `BUILD_LOCAL`   | `1`                                            | If set to `1`, it builds from the local source code in the Docker context. If `0`, it clones from `BUILD_REPO`. |
| `BUILD_REPO`    | `https://github.com/letheanVPN/blockchain.git` | The git repository to clone when `BUILD_LOCAL=0`.                                                               |
| `BUILD_TARGET`  | `gcc-linux-armv8`                              | The Conan build profile target. Profiles are located in `cmake/profiles/`.                                      |
| `BUILD_FOLDER`  | `build/release`                                | The output folder for the build.                                                                                |
| `BUILD_TYPE`    | `Release`                                      | The CMake build type (e.g., `Release`, `Debug`).                                                                |
| `BUILD_TESTNET` | `1`                                            | If set to `1`, it builds the testnet version and creates symlinks for binaries without the `-testnet` suffix.   |

### Example: Building for a different target

To build for a different architecture, you can change the `BUILD_TARGET`. For example, for `gcc-linux-x86-64`:

```bash
docker build --build-arg BUILD_TARGET=gcc-linux-x86-64 -t lethean/chain-service .
```

### Example: Building from a git branch

To build a specific branch from the git repository instead of local files:

```bash
docker build --build-arg BUILD_LOCAL=0 --build-arg BUILD_BRANCH=main -t lethean/chain-service .
```
