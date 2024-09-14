# Load Balancer System in Go - https://blog.lowlevelforest.com/

## Overview

This Go-based load balancer is designed to distribute incoming HTTP and HTTPS requests across multiple backend servers using various load balancing algorithms. The system supports:

- **Round Robin**: Sequential distribution of requests.
- **Least Connections**: Routing requests to the server with the fewest active connections.
- **Weighted Round Robin**: Distribution based on server weights.
- **IP Hash**: (Placeholder for future implementation) Routes requests based on the client's IP address.

The system can handle a large number of connections and automatically reloads the server list every 5 seconds from a configuration file.

## Features

- **Dual Protocol Support**: Handles both HTTP and HTTPS requests.
- **Dynamic Server Management**: Automatically reloads server configurations.
- **Multiple Load Balancing Algorithms**: Choose from Round Robin, Least Connections, Weighted Round Robin, and IP Hash.
- **Scalable**: Designed to handle high traffic and numerous concurrent connections.

## Prerequisites

- **Go 1.18+**: Ensure Go is installed on your machine. Download it from [golang.org](https://golang.org/dl/).
- **SSL Certificates**: For HTTPS support, you'll need SSL certificate files.

## Installation

1. **Clone the Repository**

   ```sh
   git clone https://github.com/coffeecms/coffee_load_balancer.git
   cd coffee_load_balancer
   ```

2. **Build the Application**

   ```sh
   go build -o coffee_load_balancer main.go
   ```

3. **Generate SSL Certificates (for HTTPS)**

   If you don't have SSL certificates, generate self-signed certificates for testing:

   ```sh
   openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt
   ```

   Place `server.crt` and `server.key` in the same directory as the built application.

## Configuration

1. **Server List Configuration**

   Create a `servers.conf` file in the project directory. This file should list the backend servers with optional weights, one per line:

   ```
   <server_address>:<weight>
   ```

   Example `servers.conf`:

   ```
   192.168.1.1:10
   192.168.1.2:5
   192.168.1.3
   ```

   - `<server_address>`: The IP address or hostname of the backend server.
   - `<weight>`: (Optional) The weight of the server for the Weighted Round Robin algorithm.

2. **Load Balancing Algorithm**

   Edit the `algorithm` field in the `main()` function of `main.go` to select the load balancing algorithm:

   - `"round_robin"`
   - `"least_connections"`
   - `"weighted_round_robin"`
   - `"ip_hash"`

   *Note: IP Hashing is a placeholder for future implementation.*

## Usage

1. **Run the Load Balancer**

   Start the load balancer application:

   ```sh
   ./coffee_load_balancer
   ```

   The load balancer will start two servers:

   - **HTTP**: Listens on port `8080`
   - **HTTPS**: Listens on port `8443`

2. **Sending Requests**

   - **HTTP Requests**: Use the HTTP port for regular requests:
   
     ```sh
     curl http://localhost:8080/your-endpoint
     ```

   - **HTTPS Requests**: Use the HTTPS port for secure requests. You may need to bypass SSL certificate validation for self-signed certificates:
   
     ```sh
     curl https://localhost:8443/your-endpoint --insecure
     ```

3. **Automatic Server List Reload**

   The system will automatically reload the server list every 5 seconds. Ensure `servers.conf` is updated with the latest backend servers as needed.

## Notes

- **SSL Certificates**: For production, use certificates from a trusted Certificate Authority (CA). Update the `server.crt` and `server.key` file names in the code if different.
- **Firewall and Network Configuration**: Ensure the ports used (8080 for HTTP and 8443 for HTTPS) are open and accessible.
- **Scaling**: The load balancer is designed for high concurrency, but ensure your backend servers and infrastructure can handle the load.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you have suggestions or improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

For questions or further information, please contact [LowLevelForest](mailto:lowlevelforest@gmail.com).

```

### Key Sections:

- **Overview**: Provides a brief description of the load balancer and its features.
- **Features**: Lists the capabilities of the system.
- **Prerequisites**: Details the requirements to run the system.
- **Installation**: Step-by-step guide for cloning the repo, building the application, and generating SSL certificates.
- **Configuration**: Instructions for setting up the `servers.conf` file and selecting the load balancing algorithm.
- **Usage**: How to run the application and send requests.
- **Notes**: Additional considerations and tips for production use.
- **Contributing**: Encourages community involvement.
- **License**: License details.
- **Contact**: Provides a way to reach out for further information.