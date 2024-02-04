<div align="center">

<img src="./assets/sync-stream@1x.png" alt="Sync-Stream" width="100">

# Sync-Stream

Automatically process files uploaded to an **S3 Bucket**, indexing their details in a **DynamoDB** table leveraging **Lambda** for event-driven execution.

[Getting Started](#-getting-started) •
[Components](#components) •
[Deployment](#-deployment) •
[Usage](#-usage) •
[Development](#-development) •
[Testing](#-testing)

</div>

## 🚀 Getting Started

### Installation

1. Clone this repository:
    ```bash
    git clone https://github.com/tomweston/sync-stream.git
    ```

2. Navigate to the `lambda` directory and install dependencies:
    ```bash
    cd sync-stream/lambda
    npm install
    ```

3. Deploy the infrastructure using Pulumi:
    ```bash
    cd ..
    pulumi up
    ```

## Components

- **`Pulumi.yaml` and `Pulumi.main.yaml`**: Define the Pulumi project and stack configuration.
- **`main.go`**: Contains the Go code defining the infrastructure resources using Pulumi's AWS SDK.
- **`lambda/index.js`**: The Node.js code executed by the Lambda function.
- **`lambda/index.test.js`**: Contains unit tests for the Lambda function's logic.

## 🚀 Deployment

The project utilizes GitHub Actions for continuous integration and deployment. Changes to the main branch trigger automated deployment to AWS via Pulumi.

## 🔧 Usage

1. **File Upload**: Upload files to the specified AWS S3 bucket.
2. **Event Processing**: Lambda processes each upload, extracting metadata.
3. **Data Indexing**: Metadata is stored in DynamoDB, providing a searchable index.

## 🤝 Development

To contribute or modify:

1. Clone the repository and install dependencies as described in [Installation](#installation).
2. Make necessary changes to the infrastructure (`main.go`) or Lambda function (`index.js`).
3. Deploy updates with `pulumi up`.

## ✅ Testing

Run `npm test` in the `lambda` directory to execute unit tests, ensuring the Lambda function processes file events as expected.
