FROM golang:1.11.6-alpine3.9

RUN apk add --update \
    bash \
    bash-completion \
    build-base \
    git \
    vim

# Some environment variables to influence Terraform behavior
ENV TF_LOG=TRACE TF_LOG_PATH=/var/log/terraform.log

ADD docker/passwd /etc/passwd
ADD docker/bashrc.sh /root/.bashrc

# Prevent having to download everything once for every container instantiated
RUN mkdir -p /tmp/terraform-provider-citrixitm/citrixitm/
COPY go.* /tmp/terraform-provider-citrixitm/
COPY citrixitm /tmp/terraform-provider-citrixitm/citrixitm/
WORKDIR /tmp/terraform-provider-citrixitm
RUN go mod download

# Install Terraform binary for use executing ad hoc commands
WORKDIR /tmp
RUN wget https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_amd64.zip; unzip terraform_0.11.13_linux_amd64.zip; rm terraform_0.11.13_linux_amd64.zip; mv terraform /usr/local/bin/

# Set the default working asdfa directory
WORKDIR /terraform-provider-citrixitm
