# Paxos Simulator

A simple implementation of Basic- and Multi-Paxos.

## How to Use

The project is written in Go, but it should be executed using Docker to render the output as Mermaid-cli sequence diagrams.

To run the project, simply call: `docker-compose up --build`. The program writes output to text files in the `artifacts` directory which are converted into [Mermaid.js](https://mermaid-js.github.io/mermaid/#/) sequence diagrams.

## Examples

### Basic-Paxos
![basic-paxos](./artifacts/basic-paxos-example.png)


### Multi-Paxos
![basic-paxos](./artifacts/multi-paxos-example.png)