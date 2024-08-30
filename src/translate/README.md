# Sango Phrase Translation

This tool accepts a phrase in Sango and provides a word-for-word translation into English.

## Usage

The `translate` command accepts two subcommands:

- `translate cli` _"Sango phrase"_
  - Does not block and returns the translation to stdout.
- `translate www [-port 8000]`
  - Blocks and launches a service [on the specified port, defaulting to 8000].
  - User can interact with the service by loading [translate.html](translate.html) into a browser.
