# Sango Phrase Translation

This tool accepts a phrase in Sango and provides a word-for-word translation into English.

## Installation

Currently, the (toy) service running <code>translate www</code> interacts with a
(toy) HTMX example client when navigating the browser to http://localhost:8000.

To make this demo functional, you must first have downloaded and installed a Postgres
database running locally and have restored the sample `dvdrental` database following
[these instructions](https://www.postgresqltutorial.com/postgresql-getting-started/postgresql-sample-database/).

## Usage

The `translate` command accepts two subcommands:

- `translate cli` _"Sango phrase"_ **[TODO]**
  - Does not block and returns the translation to stdout.
- `translate www` ~~[-port 8000]~~
  - Blocks and launches a service [~~on the specified port, defaulting~~hardwired to 8000].
  - User can interact with this via the HTMX client in [translate.html](translate.html).
