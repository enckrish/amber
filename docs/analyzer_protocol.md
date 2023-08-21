# Analyzer Protocol

## Setup
Set the following values in .env file:
```env
MONGO_PASSWORD=
AMBER_KAFKA_URL=
```
## Identifying if a logs should be analyzed
Each analyzer should specify a `group_id` and a `self_id`. If `self_id` is not the current leader of the `group_id`, it should not perform any analysis.
If `group_id` is not added as an intended target for the service named in the logs, it should skip the logs.

For accessing current leader and intended targets, the following MongoDB parameters are used:
```bash
ROUTING_DB_NAME = "routing-db"
TARGET_MAP_COLL_NAME = "target-map"
GROUP_ANALYZER_COLL_NAME = "group-leader"
```

Example Python Code:
```python
def is_target(service: str) -> bool:
    client = get_mongo_client()

    db = client.get_database(ROUTING_DB_NAME)
    map_coll = db.get_collection(TARGET_MAP_COLL_NAME)
    tmap = map_coll.find_one({"service": service})
    if tmap is None:
        return False
    associated_groups = tmap.get('targets')
    if GROUP_ID not in associated_groups:
        logger.info("group id mismatch, accepted", associated_groups)
        return False

    leader_coll = db.get_collection(GROUP_ANALYZER_COLL_NAME)
    leader = leader_coll.find_one({"group_id": GROUP_ID})

    if leader is None:
        return False

    group_leader = leader.get('leader')
    if group_leader != SELF_ID:
        logger.info(f"self is not leader, self id: {SELF_ID} leader: {group_leader}")
        return False

    logger.info("self is leader of a accepted group, analysing")
    return True

```

## Request and Response Guidelines

The Analyzer should listen to the `topic.log.requests.analysis.1` topic on Kafka to get the requests.

Results should be published in the `topic.log.analysis.result.1` topic.

Requests are of this type, and is JSON encoded as string:
```
{
  streamId: string
  messageId: string
  logs string[]
}
```
For more details, see `pb/router.proto`.

Results should be of the form:
```
{
    streamId: int
    messageId: int
    rating: int
    review: string
    actions: string[]
    citation: int
}
```

`rating` of 0 is considered to be `ERR`, and is used when the LLM didn't output valid severity rating. `rating` 1 to 5 correspond to:
```
1: none
2: low
3: medium
4: high
5: critical
```

