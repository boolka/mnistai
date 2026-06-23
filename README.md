# mnistai

A small Go-based MNIST neural network project for training, testing, exporting misclassified digit images, and running a tiny web UI for live predictions. The repository uses the `goai` network implementation and embedded MNIST helpers for quick experimentation.

## Repository structure

- `cmd/train` - trains the network and saves weights to a file.
- `cmd/test` - loads a saved network and evaluates it against the MNIST test set.
- `cmd/incorrect` - loads a saved network and exports incorrectly classified test digits as PNG images.
- `cmd/server` - small web server that loads a saved network, serves the static UI, and exposes a `POST /prediction` endpoint.
- `pkg/config` - contains the JSON-backed network configuration model.
- `pkg/server` - server helpers and request/response types for the web API (`PredictionRequest` / `PredictionResponse`).
- `static/` - single-page web UI (draw a digit, send to `/prediction`, see outputs).

## Prerequisites

- Go 1.20+ (project tested with recent Go versions; adjust as needed)
- Network access for module dependencies

## Configuration

The default configuration file is `config.json`.

Example:

```json
{
  "layers": [784, 128, 64, 10],
  "epochs": 10,
  "random_rate": 1.0,
  "skip_rate": 0.1,
  "learning_rate": 0.001
}
```

Configuration fields:

- `layers` - array of layer sizes for the network; the first value should be `784` for MNIST inputs and the last value should be `10` for digit classes.
- `epochs` - number of training epochs.
- `random_rate` - probability of randomly mutating network weights during training or network initialization, depending on the implementation path.
- `skip_rate` - probability of skipping a training sample during each loop iteration.
- `learning_rate` - learning rate used when correcting weights.

## Build

From the repository root build each command you need. Example:

```bash
go build ./cmd/train
go build ./cmd/test
go build ./cmd/incorrect
go build ./cmd/server
```

## Usage

### Train

Train a new network or continue training an existing one.

```bash
go run ./cmd/train -config config.json -net ai.net
```

If `ai.net` does not exist, the command will create a new randomized network and save it after training.

Optional flags:

- `-config`, `-c` - path to the config JSON file (default `config.json`)
- `-net`, `-n` - path to the network file to save or load (default `ai.net`)

### Test

Evaluate a saved network on the MNIST test set.

```bash
go run ./cmd/test -config config.json -net ai.net
```

The command prints total correct, total sureness above threshold, total incorrect, and accuracy.

Optional flags:

- `-config`, `-c` - path to the config JSON file (default `config.json`)
- `-net`, `-n` - path to the saved network file to load (default `ai.net`)

### Export incorrect images

Generate PNGs for misclassified test digits.

```bash
go run ./cmd/incorrect -config config.json -filename ai.net -out incorrect
```

### Web server + UI

Start the web server which loads a saved network, serves the static UI from `static/`, and exposes a prediction API:

```bash
go run ./cmd/server -config config.json -net ai.net
```

The server listens on port `8080` by default. Open http://localhost:8080 in a browser to use the drawing UI.

Client notes:

- The web UI posts JSON to `POST /prediction` with the body `{"inputs": [784 integers 0–255]}`.
- The server responds with `{"outputs": [10 numeric values]}` where the index of the largest value is the predicted digit.
- Static files are served from the `static/` directory relative to the server's working directory.

Optional flags:

- `-config`, `-c` - path to the config JSON file (default `config.json`)
- `-filename`, `-f` - path to the saved network file to load (default `ai.net`)
- `-out`, `-o` - output directory for PNG files (defaults to the current directory)

## Notes

- The project uses `github.com/boolka/goai/pkg/network` for network creation, activation, correction, and serialization.
- MNIST dataset loading uses embedded `github.com/boolka/mnistdb` and `github.com/boolka/mnistidx` packages.
- The training loop applies a simple correction rule using the configured learning rate.

Implementation details worth noting:

- `cmd/server` starts the web server on port `8080` and calls into `pkg/server` to register handlers and serve files from `static/`.
- The server API uses two types defined in `pkg/server`: `PredictionRequest{Inputs []float64}` and `PredictionResponse{Outputs []float64}`.

## Package overview

- `pkg/config` defines `NetworkConfig` and JSON mapping for training parameters.
- `pkg/server` defines a small web API and server helpers used by `cmd/server`.
- `cmd/train` loads config, creates or loads a network, trains it, and saves the network file.
- `cmd/test` loads a saved network and evaluates it on the MNIST test set.
- `cmd/incorrect` exports misclassified MNIST digits as PNG files.
- `cmd/server` runs the web server and serves `static/index.html` as the UI.

## License

MIT License - see [LICENSE](LICENSE) file for details.
