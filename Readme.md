# lm-zap-hook

`lm-zap-hook` is a zapcore.Core implementation for sending the logs from the application to the Logicmonitor Platform.

## Installation

`go get -u github.com/logicmonitor/lm-zap-hook`

## Quick Start

### Authentication:

Set the `LM_ACCESS_ID` and `LM_ACCESS_KEY` for using the LMv1 authentication. The company name or account name must be set to `LM_ACCOUNT` property. All properties can be set using environment variable.

| Environment variable |	Description                                        |
| -------------------- | ------------------------------------------------------|
|   LM_ACCOUNT         | Account name (Company Name) is your organization name |
|   LM_ACCESS_ID       | Access id while using LMv1 authentication.|
|   LM_ACCESS_KEY      | Access key while using LMv1 authentication.|

### Getting Started

Here's an example code snippet for configuring the `lm-zap-hook`:

```go
    // create a new Zap logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	// create resource tags for mapping the log messages to a unique LogicMonitor resource
    resourceTags := map[string]string{"system.displayname": "test-device"}

  	// create a new core that sends zapcore.WarnLevel and above messages to Logicmonitor Platform
	lmCore, err := lmzaphook.NewLMCore(context.Background(),
		lmzaphook.Params{ResourceMapperTags: resourceTags},
		lmzaphook.WithLogLevel(zapcore.WarnLevel),
	)

	// Wrap a NewTee to send log messages to both your main logger and to Logicmonitor
	logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
	  return zapcore.NewTee(core, lmCore)
	}))

	// This message will only go to the main logger
	logger.Info("Test log message for main logger", zap.String("foo", "bar"))

	// This warning will go to both the main logger and to Logicmonitor.
	logger.Warn("Warning message with fields", zap.String("foo", "bar"))
```

Complete source code of the example can be found [here]().

### Options:

Following are the options that can be passed to `NewLMCore()` to configure the `lmCore`:

| Option                                     |   Description                                                                    |             
|--------------------------------------------|----------------------------------------------------------------------------------|
|   WithLogLevel(logLevel)                   | Configures `lmCore` to send the logs having level equal or above the level specified by `logLevel` |
|   WithClientBatchingEnabled(batchInterval) | Enables batching of the log messages for the interval specified by `batchInterval` |
|   WithMetadata(metadata)                   | Metadata to be sent with the every log message.                                    |
|   WithNopLogIngesterClient()               | Configures `lmCore` to use the nopLogIngesterClient which discards the log messages. Can be used for testing.                          |


