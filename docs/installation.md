# Installation

You can install `docker-ai` using one of the methods below.

## Homebrew (macOS & Linux)

First, tap the repository:

```sh
brew tap Aj7Ay/homebrew-tap
```

Then, install `docker-ai`:

```sh
brew update
brew install docker-ai
```

## Debian/Ubuntu (.deb package)

For Debian-based Linux distributions like Ubuntu, you can download the latest `.deb` package from our GitHub Releases page and install it using `dpkg`.

1.  Go to the [**Releases Page**](https://github.com/Aj7Ay/docker-ai/releases).
2.  Download the `.deb` file from the latest release (e.g., `docker-ai_0.2.5_amd64.deb`).
3.  Install the package:

```bash
sudo dpkg -i docker-ai_*.deb
```

If you encounter any dependency issues, run the following command to fix them:

```bash
sudo apt-get install -f
```

## From Source

You can also build `docker-ai` from source if you have a Go environment set up.

1.  Clone the repository:
    ```sh
    git clone https://github.com/Aj7Ay/docker-ai.git
    cd docker-ai
    ```
2.  Build the binary:
    ```sh
    go build -o docker-ai ./cmd/docker-ai
    ```
3.  Move the binary to a directory in your `$PATH`:
    ```sh
    sudo mv docker-ai /usr/local/bin/
    ``` 