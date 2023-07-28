# Project Structure

## API

The root of the project directory contains the basic cross platform source 
code, this includes the specific implementations and the API. Each supported
platform has another folder, following the naming scheme 
`wastebasket_platform`. This folder contains code that compiles on all 
platforms, but will only be usable when executed on a specific platform.
Usage needs to be asserted dynamically at runtime.

## CLI

Additionally to the library functionallity, this library offers a CLI.
There is one root command called `wastebasket` and multiple subcommands, such
as `trash`, `restore` and more. The subcommands can all be compiled 
separately, but the code is shared with the main command.

The actual sub command implementations are all located in `cmd/impl`.
