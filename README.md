[![codecov](https://codecov.io/gh/kubeopsskills/cloud-secret-resolvers/branch/main/graph/badge.svg?token=t65R7COoaz)](https://codecov.io/gh/kubeopsskills/cloud-secret-resolvers)
# Cloud Secret Resolvers (CSR)

Cloud Secret Resolvers is a set of tools to help your applications (on Kubernetes) to retrieve any credentials from cloud managed vaults without the needed to write additional boilerplate code in your applications!.

<!-- TOC -->

- [Cloud Secret Resolvers (CSR)](#cloud-secret-resolvers-csr)
  - [Installation](#installation)
  - [Using on Kubernetes](#using-on-kubernetes)
  - [How it works](#how-it-works)
  - [Dev tools](#dev-tools)

<!-- /TOC -->

## Installation

Cloud Secret Resolvers is available on Linux, ARM, macOS and Windows platforms.
- Binaries for Linux, ARM, Windows and Mac are available as tarballs in the [release](https://github.com/kubeopsskills/cloud-secret-resolvers/releases) page

## Using on Kubernetes

- AWS
  
  - Prerequisites:
    - Enabled the OIDC provider on your [EKS](https://aws.amazon.com/th/eks/) cluster (https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html)
    - Your application Kubernetes pod has a service account with the following privillege:
       [policy.json](assets/policy.json)
  - Update your application entrypoint as follows:
    ```bash
    #!/bin/bash
    eval $(csr)
    node ... # your application runtime command
    ```
  - Update your application Kubernetes config maps as follows:
    ```yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: [your config map name]
      namespace: [your config map namespace name]
    data:
    ...
    CLOUD_TYPE: "aws"
    AWS_REGION: "[your AWS region name]"
    AWS_SECRET_NAME: "[your AWS secret name]"
    ```

- Azure
  - Prerequisites:

    A. Install az cli with [link](https://docs.microsoft.com/cli/azure/install-azure-cli)
    B. Login with `az login`
    C. Create the service principal with `az ad sp create-for-rbac`
    D. Add Access policies to azure vault secret with step B. service principal app with GET permissions for reading secrets

  - Update your application entrypoint as follows:
    ```bash
    #!/bin/bash
    eval $(csr)
    node ... # your application runtime command
    ```
  - Update your application Kubernetes config maps as follows:
    ```yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: [your config map name]
      namespace: [your config map namespace name]
    data:
    ...
    CLOUD_TYPE: "azure"
    AZ_REGION: "[your Azure region name]"
    AZ_VAULT_URL: "[your azure key vault url like: https://example.vault.azure.net]"
    AZ_CLIENT_ID: "[your service principle appId from step C. prerequisites]"
    AZ_TENANT_ID: "[your service principle tenant from step C. prerequisites]"
    ```
  - Update your application Kubernetes secret as follows:
    ```yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: secret-config
    type: Opaque
    data:
      # AZ_CLIENT_SECRET is service principle password from step C. prerequisites
      AZ_CLIENT_SECRET: "< base64 encoded AZ_CLIENT_SECRET >"
    ```
  
- Google Cloud
  - Coming Soon!

## How it works
The architecture looks like below.

Internally, the `CSR` find local environment variables in the Kubernetes Pod Container which have Cloud Vault key placeholders for example: export db_username=${db_username}, then the `CSR` will extract db_username as a key and ${db_username} as a value. Finally, the `CSR` will use ${db_username} to match cloud vault key, retrieve cloud vault value, and map the value with db_username local environment.

![Diagram](https://github.com/kubeopsskills/cloud-secret-resolvers/blob/main/assets/diagram.png)

## Dev tools

> install make for using command with `brew install make`

- `make run` for running app
- `make test` for testing
- `make test-coverage` for export test coverage
- `make all` for building app for all OS
- `make clean` for cleaning build app
