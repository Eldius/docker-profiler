FROM ubuntu:latest
RUN apt-get update && apt-get install -y stress
CMD ["stress", "--cpu", "1", "--io", "1", "--vm", "1", "--vm-bytes", "128M", "--timeout", "10s"]
