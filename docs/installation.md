# Installation

You can install `docker-ai` using one of the methods below.

## Homebrew (macOS & Linux)

If you are on macOS or Linux and use Homebrew, you can install `docker-ai` from our custom tap:

```bash
brew install syedhabeeb/tap/docker-ai
```

## Debian/Ubuntu (.deb package)

For Debian-based Linux distributions like Ubuntu, you can download the latest `.deb` package from our GitHub Releases page and install it using `dpkg`.

1.  Go to the [**Releases Page**](https://github.com/syedhabeeb/docker-aj/releases).
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

```bash
git clone https://github.com/syedhabeeb/docker-aj.git
cd docker-aj
go build ./cmd/docker-ai
```
This will create a `docker-ai` binary in the project directory. 