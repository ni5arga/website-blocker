

# Website Blocker

## Overview

A Go utility to block and unblock websites by modifying the `hosts` file on Windows, Linux, and macOS.

## Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/ni5arga/website-blocker
   cd website-blocker
   go mod download
   ```

2. **Build the Executable**

   ```bash
   go build -o blocker.exe blocker.go
   ```

   This will create the `blocker.exe` executables.

## Configuration

1. **Create `sites.txt`**

   List the websites to block or unblock, one per line:

   ```
   www.example.com
   example.com
   instagram.com
   ```

2. **Edit Time Limits**

   Modify the `startBlock` and `endBlock` variables in `blocker.go` to set your desired blocking times:

   ```go
   var (
       startBlock = "09:00" // Blocking start time
       endBlock   = "17:00" // Blocking end time
   )
   ```

3. **Adjust Host Path (if needed)**

   The path to the `hosts` file is set automatically based on your operating system. Change it in `blocker.go` if necessary.

## Usage

### Website Blocker

Run the blocker with:

```bash
sudo go run blocker.go
```

## Output

- **Blocked**: `Blocked <website>`
- **Unblocked**: `Unblocked <website>`
- **Errors**: `Error: <message>`

## License

MIT License. See [LICENSE](LICENSE) for details.

