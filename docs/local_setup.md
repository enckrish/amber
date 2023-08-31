# Local Setup

The steps we will follow are:
1. Start Kafka.
2. Setup MongoDB.
3. Set `MONGO_PASSWORD` and `AMBER_KAFKA_URL` variables in environment.
4. Start the Router.
5. Start analyzers, and register them through the admin APIs in Router.
6. We are good to go now. Now we can use Amber CLI for any file.

## [Optional] Starting Kafka Docker Instance
These instructions are written for the purpose of easy testing, but not from the standpoint of security or performance, so please keep that in mind.

We will be accessing the Kafka service from outside the host, so we will expose the Kafka port (9092) externally. I will be using `ngrok` for this purpose. 
```bash
ngrok tcp localhost:9092
```

Next we will start the Kafka instance using Docker using the `bashj79/kafka-kraft` image. Any image may be chosen. The reason behind choosing this was easy configuration for my simple needs and because of not needing Zookeeper.
```bash
docker run -p 9092:9092 -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://tcp.in.ngrok.io:XXXXX bashj79/kafka-kraft
```
Replace `XXXXX` with the port number that ngrok allocated to you.

## Starting the Router
The router creates the topics needed for Amber to work, and sets up the gRPC endpoints for communication from the CLI.
Remember to add `AMBER_KAFKA_URL` set to the Kafka endpoint, and the `MONGO_PASSWORD` in the environment.
```
git clone https://www.github.com/enckrish/amberine-router.git
cd amberine-router && go get
go run .
```

## Starting Analyzers
Remember to add `AMBER_KAFKA_URL` set to the Kafka endpoint, and the `MONGO_PASSWORD` in the environment.
Steps to start the analyzers will vary according to the analyzer code used. For the example analyzer: `aquamarine`, the process will be:
1. Install Pytorch with CUDA support from [here](https://pytorch.org/).
2. Install all required packages using `pip -r requirements.txt`
3. Run `main.py`.

## Starting the CLI
Refer to `docs/quickstart.md`.
