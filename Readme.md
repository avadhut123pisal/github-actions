| Option                                     |   Description                                                                    |             
|--------------------------------------------|----------------------------------------------------------------------------------|
|   WithLogLevel(logLevel)                   | Configures `lmCore` to send the logs having level equal or above the level specified by `logLevel` |
|   WithClientBatchingEnabled(batchInterval) | Enables batching of the log messages for the interval specified by `batchInterval` |
|   WithMetadata(metadata)                   | Metadata to be sent with the every log message.                                    |
|   WithNopLogIngesterClient()               | Configures `lmCore` to use the nopLogIngesterClient which discards the log messages. Can be used for testing.                          |

