# Load Balancer System in Go - https://blog.lowlevelforest.com/

## Overview

This Go-based load balancer distributes incoming HTTP and HTTPS requests across multiple backend servers using various load balancing algorithms. The system supports:

- **Round Robin**: Distributes requests sequentially.
- **Least Connections**: Routes requests to the server with the fewest active connections.
- **Weighted Round Robin**: Distributes requests based on server weights.
- **IP Hash**: (Placeholder for future implementation) Routes requests based on the client's IP address.

The system supports both HTTP and HTTPS protocols and automatically reloads the server list every 5 seconds from a configuration file.

## Features

- **Dual Protocol Support**: Handles both HTTP and HTTPS requests.
- **Dynamic Server Management**: Automatically reloads server configurations.
- **Multiple Load Balancing Algorithms**: Choose from Round Robin, Least Connections, Weighted Round Robin, and IP Hash.
- **Scalable**: Designed to handle high traffic and numerous concurrent connections.

## Prerequisites

- **Go 1.18+**: Ensure Go is installed. Download from [golang.org](https://golang.org/dl/).
- **SSL Certificates**: Required for HTTPS support. Generate or obtain valid SSL certificates.

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

   If you don’t have SSL certificates, generate self-signed certificates for testing:

   ```sh
   openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt
   ```

   Place `server.crt` and `server.key` in the same directory as the built application.

## Configuration

1. **Server List Configuration**

   Create a `servers.conf` file in the project directory with the following format:

   ```
   <IP>:<Port>:<Weight>
   ```

   Example `servers.conf`:

   ```
   10.220.3.1:23252:10
   10.220.3.2:12422:5
   10.220.3.3:54322:15
   ```

   - `<IP>`: The IP address of the backend server.
   - `<Port>`: The port number of the backend server.
   - `<Weight>`: (Optional) The weight of the server for Weighted Round Robin.

2. **Load Balancing Algorithm**

   Update the `algorithm` field in the `main()` function of `main.go` to select the load balancing algorithm:

   - `"round_robin"`
   - `"least_connections"`
   - `"weighted_round_robin"`
   - `"ip_hash"` (Placeholder for future implementation)

   *Example for Round Robin:*

   ```go
   lb := &LoadBalancer{
       algorithm: "round_robin",
   }
   ```

## Usage

1. **Run the Load Balancer**

   Start the load balancer application:

   ```sh
   ./coffee_load_balancer
   ```

   The load balancer will start two servers:

   - **HTTP**: Listens on port `8080`
   - **HTTPS**: Listens on port `443`

2. **Sending Requests**

   - **HTTP Requests**: Use the HTTP port for regular requests:

     ```sh
     curl http://localhost:8080/your-endpoint
     ```

   - **HTTPS Requests**: Use the HTTPS port for secure requests. If using self-signed certificates, bypass SSL verification:

     ```sh
     curl https://localhost:443/your-endpoint --insecure
     ```

3. **Automatic Server List Reload**

   The system will automatically reload the server list every 5 seconds. Ensure `servers.conf` is updated with the latest backend servers as needed.

## Deployment

1. **Configure DNS**

   - Update the DNS settings for `lowlevelforest.com` to point to the IP address of your load balancer server.
   - Ensure both HTTP (port 8080) and HTTPS (port 443) are routed to the load balancer.

2. **Firewall and Network Configuration**

   - Open the required ports on your firewall:
     - Port `8080` for HTTP
     - Port `443` for HTTPS

   - Secure your server by disabling unused services and keeping your system updated.

## Monitoring and Maintenance

1. **Monitor Logs**

   - Regularly check logs for errors or issues. Ensure logging is set up to capture operational metrics.

2. **Update Server List**

   - Modify `servers.conf` as needed to add or remove backend servers. The load balancer will automatically reload the configuration.

3. **Regular Maintenance**

   - Perform regular updates and maintenance on the load balancer and backend servers to ensure stability and security.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you have suggestions or improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

For questions or further information, please contact [LowLevelForest](mailto:lowlevelforest@gmail.com).
```

### Key Sections:

- **Overview**: Describes the load balancer system and its features.
- **Features**: Lists capabilities including dual protocol support and scalability.
- **Prerequisites**: Lists the required tools and SSL certificates.
- **Installation**: Instructions for cloning the repository, building the application, and generating SSL certificates.
- **Configuration**: Details the `servers.conf` file format and how to choose the load balancing algorithm.
- **Usage**: How to run the application, send requests, and manage automatic server list reloading.
- **Deployment**: Instructions for configuring DNS, firewall, and network settings.
- **Monitoring and Maintenance**: Guidelines for log monitoring, server list updates, and regular maintenance.
- **Contributing**: Encourages contributions.
- **License**: Licensing details.
- **Contact**: Contact information for further inquiries.