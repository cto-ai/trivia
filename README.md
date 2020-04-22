![](https://cto.ai/static/oss-banner.png)

# trivia

Trivia for everyone! Trivia questions from the Open Trivia Database.

![](https://raw.githubusercontent.com/cto-ai/trivia/blob/master/assets/screenshot_cli.png)

## Requirements

To run this or any other Op, install the [Ops Platform](https://cto.ai/platform).

Find information about how to run and build Ops via the [Ops Platform Documentation](https://cto.ai/docs/overview).

## Usage

To initiate the interactive trivia CLI prompt run:

```bash
ops run cto.ai/trivia
```

## Local Development / Running from Source

**1. Clone the repo:**

```bash
git clone <git url>
```

**2. Navigate into the directory and install dependencies:**

```bash
cd trivia && go get ./...
```

**1. Run the Op from your current working directory with:**

```bash
ops run .
```

## Debugging Issues

Use the `DEBUG` flag in the terminal to see verbose Op output like so:

```bash
DEBUG=trivia:* ops run trivia
```

When submitting issues or requesting help, be sure to also include the version information. To get your ops version run:

```bash
ops -v
```

## Resources

### trivia Docs

- [Open Trivia DB](https://opentdb.com/api_config.php)
## Contributing

See the [Contributing Docs](CONTRIBUTING.md) for more information.

## Contributors

<table>
  <tr>
    <td align="center"><a href="https://github.com/aschereT"><img src="https://github.com/aschereT.png" width="100px;" alt="aschereT's face here"/><br /><sub><b>Vincent Tan</b></sub></a><br/></td>
  </tr>
</table>

## LICENSE

[MIT](LICENSE)
