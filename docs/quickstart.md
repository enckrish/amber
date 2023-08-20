# Quickstart

To run Amber CLI on any file, run through the following steps:

1. Define `AMBER_KAFKA_URL` as an environment variable. It should contain the URL contained in `KAFKA_ADVERTISED_LISTENER`. Amber doesn't support multiple brokers for now.
2. Amber supports multiple flags to customize the behavior. To see available options: `amber --help`.
3. To tail a log file for analysis, use the following command:
```
amber -p <absolute path to file> -t <service name e.g. Docker, nginx etc.> -bs <number of lines to send at once> -hs <number of analysis to retain for historical context>
```

