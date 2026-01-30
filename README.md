# Minecraft Docker Server Manager

This is a simple web-based manager for Minecraft servers that run in Docker containers. Written in Go using the Gin framework, it allows you to view your Minecraft server containers, see their resource usage (CPU and Memory), state adn status, port, and start or stop them directly from the web interface.

## Features

-   **List Containers**: View all containers on the host machine.
-   **Container Stats**: See CPU and Memory usage for each container.
-   **Toggle Containers**: Start and stop containers with the click of a button.
-   **Embedded UI**: The frontend is embedded into the Go binary, making it a single, easy-to-distribute file.

## Project Structure

-   `main.go`: The main application logic.
-   `index.html`: The frontend HTML file (embedded in the binary).
-   `go.mod`, `go.sum`: Go module files.
-   `makefile`: Contains build commands.
-   `Dockerfile`: An example showing how a minecraft server image can be made. This image is published on docker hub.
-   `server1/`: This directory is an example of how a managed server could be set up. It includes a `docker-compose.yml` file to demonstrate a potential use case.

## Getting Started

### Prerequisites

-   Go (version 1.16+ for `embed` support)
-   Docker
-   Make (optional, for using the makefile)

### Build and Run

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/Rithwik0604/minecraft-server-manager.git
    cd minecraft-server-manager
    ```

2.  **Build the binary:**
    You can use the provided makefile (builds for linux):
    ```bash
    make build
    ```
    Or build it manually:
    ```bash
    go build -o mcserver-manager
    ```

3.  **Run the application:**
    ```bash
    ./mcserver-manager
    ```

4.  Open your web browser and navigate to `http://localhost:8080`.
