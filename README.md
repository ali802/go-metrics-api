# High-Efficiency Infrastructure & Go Metrics API

A high-performance, ultra-lean Go microservice built with the Chi Router that captures system telemetry and exposes real-time runtime metrics. This project serves as a case study in cloud-native optimization, utilizing multi-stage Docker compilation workflows and minimalist Google Distroless runtime layers to drastically reduce the application footprint and minimize attack surfaces.

---

## Core Engineering Achievements

* Image Size Optimization: Reduced final container image footprint from a baseline of 1.27 GB down to an ultra-lean 15 MB (a 98.82% space savings).
* Vulnerability Surface Minimization: Eliminated all operating system shells (bash, sh), package managers (apt), and generic platform tools from the production runtime, establishing a near-zero attack surface.
* Resource Efficiency: Achieved an active runtime idle memory footprint of < 1.00 MB by stripping away operating system background daemons.
* Deterministic Builds: Isolated dependencies via native Go toolchain locking, generating immutable, cryptographically verified go.sum signature maps.

---

## Technical Architecture Details

### 1. Advanced Compilation Tweaks
Within the multi-stage build pipeline, the Go compiler flag arguments ensure that the binary is stripped of bloated debugging overhead:
* CGO_ENABLED=0: Disables dynamic linking into OS C-libraries (such as glibc), producing a fully static standalone binary that runs independently without host OS file dependencies.
* -ldflags="-s -w": Instructs the linker to drop debug tables (-s) and DWARF tracking paths (-w), shaving 30% to 40% off the compiled binary file size.

### 2. Multi-Stage Pipeline Separation
The runtime architecture splits the deployment lifecycle into two isolated environments:
1. Stage 1 (builder): A heavy golang:1.23-bullseye workspace used to pull dependencies, cache modules, and compile the static executable safely.
2. Stage 2 (production): A minimalist gcr.io/distroless/static-debian12:nonroot layer. It contains no operating system utilities and forces the binary to run under a restricted user security context (UID 65532).

---

## Empirical Validation Matrix

| Evaluation Metric | Single-Stage Architecture | Multi-Stage Distroless | System Delta |
| :--- | :--- | :--- | :--- |
| Image Disk Footprint | 1.27 GB | 15 MB | -98.82% |
| Shell Interactivity | Available (bash, sh) | None | Hardened |
| Runtime Security Context | Root User | Restricted (nonroot) | Secured |
| System Package Utilities | Present (apt, dpkg) | Completely Removed | Isolated |
| Idle Memory Footprint | Standard OS Overhead | < 1.00 MB | Efficient |

---

## Getting Started & Deployment

### Local Development
To download internal package dependencies and run the server natively on your machine toolchain:

go mod tidy
go run main.go

### Building the Optimized Image
To trigger the multi-stage container build process locally:

docker build -t go-metrics-api:multi -f Dockerfile.multi .

### Deploying the Microservice
Run the lean production container detached in the background, mapping traffic directly to port 9090:

docker run -d -p 9090:9090 --name live-metrics-app go-metrics-api:multi

---

## Telemetry Verification

Query the running container instance to parse live structural application data:

curl http://localhost:9090/api/v1/metrics

### Expected JSON Response Payload:

{
  "status": "HEALTHY",
  "timestamp": "2026-06-20T15:33:59.075986721Z",
  "go_version": "go1.23.12",
  "num_cpu": 4,
  "goroutines_active": 4,
  "allocated_memory_mb": 0
}
