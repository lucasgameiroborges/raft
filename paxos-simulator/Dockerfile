FROM golang:latest
LABEL maintainer="Lucas Gameiro <lucasgameirob@gmail.com>"

WORKDIR /

COPY go.mod .
COPY src /src
# Create a new user to run the raft node and
# switch to that user
RUN mkdir node && \
    mkdir -p node/cluster-data/log

# Define the volumes for this container
# that must be supplied from the host
VOLUME [ \
    "node/cluster-data/log" \
]

RUN rm -rf /node/cluster-data/log/*

CMD ["go", "run", "src/multi-paxos/main.go"]