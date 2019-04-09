FROM golang:1.11.6-alpine3.9

RUN apk add --update \
    bash \
    bash-completion \
    build-base \
    git \
    shadow \
    vim

# Some environment variables to influence Terraform behavior
ENV TF_LOG=TRACE TF_LOG_PATH=/var/log/terraform.log

RUN usermod --shell /bin/bash root
ADD docker/bashrc.sh /root/.bashrc

# Install Terraform binary for use executing ad hoc commands
WORKDIR /tmp
RUN wget https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_amd64.zip; unzip terraform_0.11.13_linux_amd64.zip; rm terraform_0.11.13_linux_amd64.zip; mv terraform /usr/local/bin/

# Prevent having to download everything once for every container instantiated
RUN mkdir -p /tmp/terraform-provider-citrixitm/citrixitm/
COPY go.* /tmp/terraform-provider-citrixitm/
COPY citrixitm /tmp/terraform-provider-citrixitm/citrixitm/
WORKDIR /tmp/terraform-provider-citrixitm
RUN go mod download

# Set the default working directory
WORKDIR /terraform-provider-citrixitm
