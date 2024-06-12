# QA Challenge Go Application

---

## Table of Contents

- [QA Challenge Go Application](#qa-challenge-go-application)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
    - [Fixed issues](#fixed-issues)
  - [Quick Launch](#quick-launch)
    - [Prerequisites](#prerequisites)
    - [Run locally, using `make` commands](#run-locally-using-make-commands)
    - [Run locally, using bash script](#run-locally-using-bash-script)
    - [Run in Docker container](#run-in-docker-container)
  - [Repository structure](#repository-structure)
  - [Recommendations and possible improvements](#recommendations-and-possible-improvements)

---

## Overview

This repository contains a Go application designed to interact with Ethereum-based smart contracts on the evm networks/chains. The application provides functionalities to deploy, call, read GetterSetter smart contract and write its values to a JSON file. This guide will help you set up, build, and run the application both locally and within a Docker container.

### Fixed issues

- Add `memory` keyword to the `setBytes` and `requestedBytes` parameters signature in `GetterSetter.sol` (i.e. `function setBytes(bytes memory _value)`).

---

## Quick Launch

### Prerequisites

Before you begin, ensure you have the following installed on your machine:

- [Git](https://git-scm.com/downloads)
- [Golang](https://golang.org/dl/) (version 1.22.3 or later)
- [Metamask](https://metamask.io/) - EOA (externally owned account) and wallet's private key
- RPC URL - [Infura](https://www.infura.io/), [Alchemy](https://www.alchemy.com/), [Chainstack](https://chainstack.com/)
- [Make](https://www.gnu.org/software/make/manual/make.html)
- [Docker](https://www.docker.com/products/docker-desktop)

### Run locally, using `make` commands

1. Clone the repository and navigate to its root:

  ```sh
  git clone https://github.com/yourusername/qa-challenge-application.git && cd qa-challenge-application
  ```

2. Create and configure a `config.toml` file in the root directory based on the [config-example.toml](config-example.toml).

The `Contract.Mode` defines the way the application will interact with the contract:

- `demo` - will deploy a new contract and call its setter and getter methods, depending on the values set under the the `Contract.Values` section
- `deploy-contract` - will deploy a new contract. Specify its address in the `Contract.Address` section to reuse it in `call-contract` mode.
- `call-contract` - Requires `Contract.Address` to be set. This mode will call setters and getters on the specified contract, depending on the values set in the `Contract.Values` section.
- `read-only-contract` - will read/fetch the values of the specified contract.  
  
The `Client` section:

- `GasLimit` - the gas limit used for the transaction. If not specified, defaults to 3000000.
- `WaitingTimeout` - the timeout used to wait for the transaction to be mined. The transaction is not reverted/cancelled if it is not mined within the timeout. If it is not specified, defaults to 300 seconds.

3. Build the application (bin):

  ```sh
  make build
  ```

4. Run the application:

  ```sh
  make run
  ```

The JSON file will be created in the `./output` directory.

### Run locally, using bash script

1. Repeat steps 1-3 in the previous section

2. Make `./run_app.sh` executable:
  
  ```sh
  chmod +x ./run_app.sh
  ```

3. Run the bash script:

  ```sh
  ./run_app.sh
  ```

### Run in Docker container

1. Create Docker image by running the following command from the root directory:

  ```sh
  docker build -t go-assessment:latest .
  ```

2. Create and run container from the previously created image:

  ```sh
  docker run -it --name go-assessment_v1 go-assessment:latest
  ```

---

## Repository structure

```txt
qa-assessment/
├── config-example.toml # Example configuration file
├── go.mod              # dependencies
├── main.go             # Application entrypoint
├── Makefile            # Tool for managing and maintaining the project
├── run_app.sh          # Bash script to run the application
└── src/
    ├── config/    # Go mappings for config.toml
    ├── contracts/ # Source code for the GetterSetter smart contract, with go bindings
    ├── evm/
    │   └── clients/
    │       └── geth/
    │           ├── Runner.go             # Geth client Runner.
    │           ├── dto/                  # Data Transfer Objects for GetterSetter smart contract. Implements Builder pattern for flexibility.
    │           ├── types/                # Data model for JSON output
    │           ├── transactions/         # Re-usable logic to handle transactions.
    │           ├── account/              # API to manage accounts-related data (private key, EOA, balance, etc).
    │           └── client/               # Geth client wrapper.
    └──utils/      # Utility functions to handle data type conversions, application configurations and JSON responses
```

**Noteworthy:**

1. The output JSON is created in the `./output` directory. The path is hardcoded in the [json_writer.go](./src/utils/json_writer.go) function. Update it if you want to write the output JSON to a different location.

---

## Recommendations and possible improvements

- Extend app with CLI API to make the application more interactive
- Add linting and code formatters to follow the recommended
- Add (unit) testing
- Add GitHub actions to run tests in CI/CD
- Add possibility to override the `config.toml` values with user-defined values (for `docker run` and CLI)
- Convert raw bytes values returned in the response to the human-readable format
