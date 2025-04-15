# Directory Reader

A powerful Go utility that scans directories and generates detailed markdown documentation of their structure and file contents. Perfect for codebase documentation, project handoffs, or preparing code for LLM analysis.

## Features

- 📂 Generates a clean Markdown representation of your directory structure
- 📄 Includes file contents in properly formatted code blocks
- 🔍 Intelligently ranks and orders files based on importance
- 🚫 Skips binary files, ensuring only readable content is included
- 🧠 Automatically detects programming languages in your codebase
- 🙈 Respects `.gitignore` and provides additional ignore capabilities
- 🔢 Counts tokens for LLM compatibility

## Installation

### Using Go

```bash
go install github.com/jvanmelckebeke/directory-reader/cmd/directory-reader@latest
```

### From Source

```bash
git clone https://github.com/jvanmelckebeke/directory-reader.git
cd directory-reader
go build -o bin/directory-reader cmd/directory-reader/main.go
```

## Usage

```bash
directory-reader [--ignore=lang1,lang2] <target_directory>
```

### Options

- `--ignore`: Comma-separated list of programming languages to ignore (e.g., 'go,python')

### Example

```bash
directory-reader --ignore=java,node ./my-project
```

This will create a `directory_content.md` file in `./my-project` containing:
1. A visual representation of the directory structure
2. The contents of each non-binary file, properly formatted

## Configuration

### Ignore Files

The tool respects the following ignore configurations:

1. `.gitignore` files in the target directory
2. `.readerignore` files in the target directory
3. Language-specific ignore patterns from [gitignore.io](https://www.toptal.com/developers/gitignore)
4. Default ignore patterns for common build artifacts

### File Ranking

Files are ordered in the output based on importance:
- README files appear last (most important for context)
- Markdown and documentation files are prioritized
- Configuration files follow
- Implementation files come next
- Test files appear first (typically less important)

## Use Cases

- **LLM Integration**: ask your favorite LLM to generate a readme, provide codereviews, ...

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request