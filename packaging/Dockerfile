FROM artifactory.ci1.us-west-2.aws-dev.app.snowflake.com/docker-remote/library/ubuntu:18.04 AS base

ENV GO_VERSION 1.23.6
ENV TERRAFORM_VERSION 1.8.4

# Update packages and install necessary tools
RUN apt-get update && apt-get install -y wget make git gnupg software-properties-common unzip

# Download, unzip, and install Terraform to a location in the PATH
RUN wget -O terraform.zip "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip" \
    && unzip terraform.zip -d /usr/local/bin \
    && rm terraform.zip

# Install Go
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz

# Add Go to the PATH
ENV PATH "/usr/local/go/bin:${PATH}"
