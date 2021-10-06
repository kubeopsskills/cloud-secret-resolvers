[![codecov](https://codecov.io/gh/kubeopsskills/cloud-secret-resolvers/branch/main/graph/badge.svg?token=t65R7COoaz)](https://codecov.io/gh/kubeopsskills/cloud-secret-resolvers)
# Cloud Secret Resolvers (CSR)

Cloud Secret Resolvers is a set of tools to help your applications (on Kubernetes) to retrieve any credentials from cloud managed vaults without the needed to write additional boilerplate code in your applications!.

<!-- TOC -->

- [Cloud Secret Resolvers (CSR)](#cloud-secret-resolvers-csr)
  - [Installation](#installation)
  - [Using on Kubernetes](#using-on-kubernetes)
  - [How it works](#how-it-works)

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
    1. Install az cli with [link](https://docs.microsoft.com/cli/azure/install-azure-cli)
    2. Login with `az login`
    3. Create the service principal with `az ad sp create-for-rbac`
    4. Add Access policies to azure vault secret with step 3. service principal appId for GET permission for authorizing to read key and secret

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
    AZ_CLIENT_ID: "[your service principle appId from step 3 prerequisites]"
    AZ_CLIENT_SECRET: "[your service principle password from step 3 prerequisites]"
    AZ_TENANT_ID: "[your service principle tenant from step 3 prerequisites]"
    ```
  
- Google Cloud
  - Coming Soon!

## How it works
The architecture looks like below.

Internally, the `CSR` find local environment variables in the Kubernetes Pod Container which have Cloud Vault key placeholders for example: export db_username=${db_username}, then the `CSR` will extract db_username as a key and ${db_username} as a value. Finally, the `CSR` will use ${db_username} to match cloud vault key, retrieve cloud vault value, and map the value with db_username local environment.

![Diagram](https://github.com/kubeopsskills/cloud-secret-resolvers/blob/main/assets/diagram.png)