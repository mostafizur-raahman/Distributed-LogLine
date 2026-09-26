It is 2 AM. Your payment service is dropping transactions. The on-call engineer SSHs into three different machines, tailing three different log files, grepping for a transaction ID that may or may not appear in any of them. By the time they find the right file on the right host, it has been 40 minutes and the incident channel has 200 messages.

This is what happens when logs live on the machines that produce them. Every service writes to its own file, in its own format, on its own disk. Searching across services means remembering which host runs which service, what log format it uses, and where the files live. Correlating a single user's request across five services means opening five terminals.

Services like Datadog Logs, Papertrail, and Grafana Loki exist because every team eventually hits this wall. The idea is straightforward: instead of writing logs to local files, applications push structured log entries to a central service over HTTP. The service stores them, indexes them, and lets you search across everything from one place. It is infrastructure that every team needs and most teams pay for.

Ingestion. Any application sends a POST request to /ingest with a JSON payload the log level, a message, the service name, a timestamp, and an optional bag of structured metadata. logline validates the payload, checks the API key, and stores it. A single logline instance handles thousands of entries per second using PostgreSQL batch inserts and connection pooling.

```json
{
    "level": "error",
    "message": "connection refused",
    "service": "payment-api",
    "timestamp": "2026-02-27T14:30:00Z",
    "data": {
        "transaction_id": "txn_8f2a9c",
        "retry_count": 3
    }
}
```
