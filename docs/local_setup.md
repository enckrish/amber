# Local Setup

The steps we will follow are:
1. Start Kafka.
2. Setup MongoDB.
3. Set `MONGO_PASSWORD` and `AMBER_KAFKA_URL` variables in environment.
4. Start the Router.
5. Start analyzers, and register them through the admin APIs in Router.
6. We are good to go now. Now we can use Amber CLI for any file.

## Starting the Router
The router creates the topics needed for Amber to work, and sets up the gRPC endpoints for communication from the CLI.
```
git clone https://www.github.com/enckrish/amberine-router.git
cd amberine-router && go get
go run .
```

## Starting Analyzers
Steps to start the analyzers will vary according to the analyzer code used. For the example analyzer: `aquamarine`, the process will be;
1. Install Pytorch with CUDA support from [here](https://pytorch.org/).
2. Install all required packages using `pip -r requirements.txt`
3. Run `main.py`.

## Starting the CLI
Refer to `docs/quickstart.md`.
