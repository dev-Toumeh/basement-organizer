FROM golang:1.22

# User variabels
ARG USERNAME=go-user
ARG GID=1000
ARG UID=1000

WORKDIR /usr/src/app

# install programms
RUN apt-get update && apt-get install -y vim

# to keep the container running
CMD ["tail", "-f", "/dev/null"]

# create new User inside the Conrainer
RUN groupadd -g ${GID} ${USERNAME} \
    && useradd -m -r -u ${UID} -g ${USERNAME} ${USERNAME}