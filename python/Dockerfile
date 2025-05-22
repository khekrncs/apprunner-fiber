FROM alpine:latest

# Install a basic web server
RUN apk add --no-cache python3

# Create a simple script that responds to health checks
RUN echo 'import http.server, socketserver; \
    handler = http.server.SimpleHTTPRequestHandler; \
    httpd = socketserver.TCPServer(("0.0.0.0", 8080), handler); \
    print("Server started at port 8080"); \
    httpd.serve_forever()' > /app.py

# Expose the port
EXPOSE 8080

# Run the server
CMD ["python3", "/app.py"]