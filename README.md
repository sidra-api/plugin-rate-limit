# Rate Limit Plugin

This is a rate limit plugin for the Sidra Api, written in Go. It limits the number of requests a client can make within a minute.

## Features

- Limits the number of requests per minute per client IP.
- Configurable rate limit through environment variables.
- Logs rate limit status for each request.

## Installation

To build and run the plugin using Docker, follow these steps:

1. Clone the repository:
    ```sh
    git clone <repository-url>
    cd <repository-directory>
    ```

2. Build the Docker image:
    ```sh
    docker build -t plugin-rate-limit .
    ```

3. Run the Docker container:
    ```sh
    docker run -e PLUGIN_NAME=rate-limit -e RATE_LIMIT=5 -p 8080:8080 plugin-rate-limit
    ```

## Configuration

The plugin can be configured using the following environment variables:

- `PLUGIN_NAME`: The name of the plugin (default: `rate-limit`).
- `RATE_LIMIT`: The maximum number of requests allowed per minute (default: `5`).

## Usage

The plugin processes incoming requests and applies rate limiting based on the client's IP address. If the rate limit is exceeded, it returns a `429 Too Many Requests` response. Otherwise, it allows the request and returns a `200 OK` response.

## Development

To develop and test the plugin locally:

1. Ensure you have Go installed on your machine.
2. Clone the repository and navigate to the project directory.
3. Build and run the plugin:
    ```sh
    go build -o rate-limit main.go
    ./rate-limit
    ```
## License

This project is licensed under the MIT License.