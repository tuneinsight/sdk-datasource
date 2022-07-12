# SDK

## Data source plugin for the Tune Insight Note

A TI Note data source plugin extends the data manager to allow clients to query external data sources.
It exposes operations that can have arbitrary parameters and results, and can potentially output data objects.

### How to develop a data source plugin

A TI Note data source plugin is a [Go plugin](https://pkg.go.dev/plugin) that exposes two variables:
- `DataSourceType` a string (of the type `sdk.DataSourceType`) that identifies uniquely the type of data source.
- `DataSourcePluginFactory` a function (of the type `sdk.DataSourcePluginFactory`) that can be invoked by the TI Note to create a new instance of the data source.

The factory function `sdk.DataSourcePluginFactory` takes as parameters:
- `logger` that allows the data source plugin logging to be integrated with the one of the TI Note;
- `config` a map of arbitrary config keys to allow the data source to be configured by the TI Note.

It returns a data source, which is a struct that implements the interface `sdk.DataSourcePlugin` with the function `Query()`.
This function takes as arguments:
- `userID` the unique identifier of the user invoking the query;
- `operation` the operation requested for the query;
- `jsonParameters` the parameters of the operation, that can be anything in the JSON format;
- `outputDataObjectsSharedIDs` the shared IDs defined by the client the potential output data objects of the operation should have.  

The `Query()` function returns:
- `jsonResults` the results of the operation, that can be anything in the JSON format;
- `outputDataObjects` the potential output data objects of the operation;
- `err` the potential error of the operation.

If the operation outputs data objects, those must be of the type `sdk.DataObject`, with the shared ID and output name provided by the client.

### Use the plugin in Tune Insight Note

The plugin compiled into a `.so` file must be placed under the `plugins/datasources` directory of the TI Note.
The plugin is then dynamically loaded at TI Note startup, after which it will be possible to instantiate,
query, and manage data sources of the newly defined type throug the TI Note API.
