

# Website Blocker 

## Overview

A Go utility to block and unblock websites by modifying the `hosts` file on Windows, Linux, and macOS.

## Installation

1. **Clone the Repository & Install Dependencies**

   ```bash
   git clone https://github.com/ni5arga/website-blocker
   cd website-blocker
   go mod download
   ```

2. **Build the Applications**

   ```bash
   go build
   ```

   This command installs all dependencies and builds the `website_blocker` and `fallback` executables.

## Configuration

1. **Create `sites.txt`**

   List websites to manage, one per line.

   ```
   www.example.com
   example.com
   instagram.com
   ```

2. **Edit Time Limits**

   Open `website_blocker.go` and modify the `startBlock` and `endBlock` variables to set the blocking time range:

   ```go
   var (
       startBlock = "09:00" // Time to start blocking
       endBlock   = "17:00" // Time to stop blocking
   )
   ```

3. **Adjust Host Path (if needed)**

   The default path to the `hosts` file is set automatically based on the operating system. Adjust it in `website_blocker.go` and `fallback.go` if necessary.

## Usage

### Website Blocker

Run the blocker:

```bash
./website_blocker
```

### Fallback Unblocker

Run the fallback unblocker:

```bash
sudo ./fallback
```

## Emergency Use

Use the `fallback` script to unblock all sites immediately.

## Output

- **Blocked**: `Blocked <website>`
- **Unblocked**: `Unblocked <website>`
- **Errors**: `Error: <message>`

## License

MIT License. See [LICENSE](LICENSE) for details.

