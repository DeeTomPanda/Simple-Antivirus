# go-Simple AV

A lightweight antivirus-style file monitoring and scanning tool written in Go.

## Features

### Supported
* Real-time filesystem monitoring
* SHA-256 file hashing
* Known-malware hash matching
* Structured logging
* Configurable watch paths
* On-demand directory scan (-scan flag)
* CLI-based operation
* SQLite hash database

### Planned
* Recursive directory watching
* YARA rule-based scanning
* Quarantine detected files
* Detection history and reporting
* Custom rule and signature management
* Process attribution (Windows)
* Windows Service support

## Detection Pipeline

```text
Filesystem Event
        │
        ▼
    Hash File
        │
        ├── Known Malware Hash Match
        │
        ▼
     YARA Scan
        │
        ▼
     Detection
        │
        ▼
       Log
```

## Roadmap

### Monitoring

* [x] Filesystem watcher
* [x] Recursive directory watching
* [x] Configurable watch targets

### Detection

* [x] SHA-256 hashing
* [x] Hash database loading
* [x] Hash-based detection
* [ ] YARA integration
* [ ] Rule management

### Response

* [x] Detection logging
* [ ] Quarantine support
* [ ] Detection history

### Platform

* [ ] Windows support
* [ ] Process attribution
* [ ] Windows Service mode

## Status

Active development. This project is intended as a learning and portfolio project focused on systems programming, malware detection workflows, and real-time filesystem monitoring in Go.

## License

MIT
