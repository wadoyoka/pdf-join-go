# pdf-join

Merge PDF files in a directory into a single PDF.

## Install

```bash
brew tap wadoyoka/tap
brew install pdf-join
```

Or build from source:

```bash
go install github.com/wadoyoka/pdf-join-go@latest
```

## Usage

```bash
# Merge PDFs in the current directory
pdf-join

# Merge PDFs in a specific directory
pdf-join /path/to/pdfs

# Specify output file name
pdf-join -o result.pdf

# Combine options
pdf-join /path/to/pdfs -o result.pdf
```

Files are merged in filename ascending order. The default output file is `merged.pdf`.

## License

This project is licensed under the [MIT License](LICENSE).

### Third-Party Licenses

This tool uses the following open-source libraries:

- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - Apache License 2.0
- [cobra](https://github.com/spf13/cobra) - Apache License 2.0
- [pflag](https://github.com/spf13/pflag) - BSD-3-Clause
- [yaml.v2](https://github.com/go-yaml/yaml) - Apache License 2.0

See the [THIRD_PARTY_LICENSES](THIRD_PARTY_LICENSES/) directory for full license texts.

You can also view license information by running:

```bash
pdf-join --credits
```
