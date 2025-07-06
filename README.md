# terragrunt-ops

## Overview

The `terragrunt-ops` application is a Go-based utility designed to streamline the process of managing complex infrastructure deployments and ultimately avoid unnececasry cloud costs. In otherwords, I run this locally to tare down and build complex AWS infrastructure to save money. It leverages Terragrunt=>Terraform to handle the provisioning and destruction of AWS resources in a controlled and automated manner.

## Features

- **Automated Infrastructure Provisioning**: Seamlessly apply your Terragrunt configurations to provision the required AWS resources in sequential order, since dependencies are critical.
- Centralized config (steps.json) to define apply/destroy order.
- **Cost Savings**: By automating the ```apply``` and ```destroy``` processes, this application helps reduce cloud costs by minimizing the time resources are running unnecessarily.

## Prerequisites

- Go 1.20+
- Terragrunt & Terraform installed and in PATH
- AWS CLI configured & authenticated.

Project path is hardcoded to:

## Getting Started

1. **Install Dependencies**: Ensure you have Go installed on your system. Additionally, make sure you have Terraform and Terragrunt set up and configured with your AWS credentials.

2. **Clone the Repository**: Clone the `terragrunt-ops` repository to your local machine.

3. **Configure**: Modify the configuration file ```steps.json``` to specify the absolute path to your Terragrunt base and your sequence of steps.

4. **Build and Run**: Build the Go application and run the binary to start the automated infrastructure management process.

## Usage
- Apply: ```./tgops --mode=apply```
- Destory: ```./tgops --mode=destroy```

## Contributing

Contributions to `terragrunt-ops` are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).Golang script to build and destroy large, complex, infrastructure in an effort to save on cloud costs
