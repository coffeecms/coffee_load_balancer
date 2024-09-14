```markdown
# Load Balancer System in Go - https://blog.lowlevelforest.com/

## Overview

This is a simple yet robust load balancer system implemented in Go. The load balancer supports several request distribution algorithms including:

- **Round Robin**: Distributes requests sequentially to each server.
- **Least Connections**: Routes requests to the server with the fewest connections.
- **Weighted Round Robin**: Distributes requests based on the server's weight.
- **IP Hash**: Routes requests based on the client's IP address (currently a placeholder).

The system automatically reloads the server list from a configuration file every 5 seconds and supports handling various types of HTTP requests (GET, POST, PUT, DELETE, etc.).

## Features

- **Dynamic Server List**: Automatically reloads server configurations from `servers.conf`.
- **Load Balancing Algorithms**: Includes Round Robin, Least Connections, Weighted Round Robin, and IP Hash.
- **Scalable**: Designed to handle a large number of concurrent connections efficiently.

## Prerequisites

- Go 1.18 or higher
- Basic knowledge of Go and HTTP servers

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/coffeecms/coffee_load_balancer.git
   cd coffee_load_balancer
   ```

2. **Build the application:**

   ```sh
   go build -o coffee_load_balancer main.go
   ```

## Configuration

1. **Create a `servers.conf` file**:

   The `servers.conf` file should contain a list of servers with optional weights. Each line should have the format:

   ```
   <server_address>:<weight>
   ```

   Example `servers.conf`:

   ```
   192.168.1.1:10
   192.168.1.2:5
   192.168.1.3
   ```

   - `<server_address>`: The IP address or domain of the server.
   - `<weight>`: (Optional) The weight of the server for weighted round-robin algorithm.

2. **Choose a load balancing algorithm**:

   Update the `algorithm` field in the `main()` function of `main.go` to one of the following:

   - `"round_robin"`
   - `"least_connections"`
   - `"weighted_round_robin"`
   - `"ip_hash"`

## Usage

1. **Run the load balancer:**

   ```sh
   ./coffee_load_balancer
   ```

2. **Send HTTP requests to the load balancer**:

   The load balancer will forward requests to the appropriate backend server based on the selected algorithm. You can test it using `curl` or any HTTP client:

   ```sh
   curl http://localhost:8080/your-endpoint
   ```

## Notes

- Ensure that your backend servers are reachable and properly configured to handle requests forwarded by the load balancer.
- The IP Hashing algorithm is a placeholder and not yet implemented. You can extend the code to support it if needed.
- The load balancer is designed to handle a large number of connections, but make sure your hardware and network infrastructure are scaled accordingly.

## Contributing

If you have suggestions or improvements, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

For any questions or further information, please contact [Your Name](mailto:your.email@example.com).

```

### Key Points:

- **Overview**: Brief introduction to what the system does.
- **Features**: Highlights the key features of the load balancer.
- **Prerequisites**: Lists the requirements to run the system.
- **Installation**: Step-by-step guide to clone, build, and set up the system.
- **Configuration**: Instructions on how to configure the `servers.conf` file and select the load balancing algorithm.
- **Usage**: How to run the application and test it.
- **Notes**: Additional information about the system's capabilities and limitations.
- **Contributing**: Encourages contributions and provides a way to get in touch.
- **License**: Specifies the license for the project.
- **Contact**: Provides contact information for further queries.

Feel free to adjust any sections based on your specific needs or additional features.