# WhisprFlow — Comprehensive Build Plan

> A local-first speech processing pipeline: **Audio → Text → Embeddings → Insights**, all running on your Ryzen 5 3500U / 14GB RAM machine.

---

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌───────────────┐     ┌──────────────┐     ┌─────────────┐
│  Audio In   │ ──► │  VAD Filter  │ ──► │  STT Engine   │ ──► │  Post-Process │ ──► │   Output    │
│ (mic/file)  │     │ (silero-vad) │     │ (whisper.cpp) │     │ (embeddings)  │     │ (text/json) │
└─────────────┘     └──────────────┘     └───────────────┘     └──────────────┘     └─────────────┘
                                                                                          │
                                                                                          ▼
                                                                                   ┌──────────────┐
                                                                                   │   Storage    │
                                                                                   │  (libSQL)    │
                                                                                   └──────────────┘
                                                                                          │
                                                                                          ▼
                                                                                   ┌──────────────┐
                                                                                   │   Search     │
                                                                                   │  (vector)    │
                                                                                   └──────────────┘
```

---

## Tech Stack

| Layer | Choice | Rationale |
|-------|--------|-----------|
| **Language** | Go | Your existing stack, great for CLI, CGO support |
| **STT** | whisper.cpp (official Go bindings) | Best CPU performance, native Go bindings exist |
| **VAD** | silero-vad via ONNX Runtime (CGO) | ~2MB model, <1ms per chunk, no Python dependency |
| **Database** | libSQL (Turso's SQLite fork) | Native vector search, no extensions needed |
| **Go DB Driver** | `go-libsql` | Official SDK, embedded mode, CGO-based |
| **Embeddings** | `all-MiniLM-L6-v2` (384-dim) | ~80MB, fast on CPU, good quality |
| **LLM (optional)** | Ollama (external API call) | Don't embed — just call localhost:11434 |
| **Audio capture** | `portaudio` via CGO or `miniaudio` | Cross-platform, mature |
| **Audio conversion** | FFmpeg (external) | Convert any input to 16kHz 16-bit WAV |

---

## Hardware-Aware Optimizations (Ryzen 5 3500U / 14GB RAM / Vega 8 iGPU)

| Constraint | Strategy |
|------------|----------|
| **No NVIDIA GPU** | whisper.cpp uses BLAS/OpenCL — no CUDA needed |
| **8 threads** | Use 6-7 threads for STT, leave headroom |
| **14GB RAM** | STT ~500MB + VAD ~200MB + libSQL ~50MB + Go ~100MB = well within budget |
| **Integrated GPU** | Vega 8 can help via GGML Vulkan backend (optional) |
| **Default model** | `base.en` (~140MB) for speed; `small` (~460MB) for accuracy |
| **VAD first** | Skip 30-60% silence = massive compute savings |
| **Chunk size** | 30s max segments (Whisper's context window) |

---

## Project Structure

```
whisprflow/
├── cmd/
│   └── whispr/
│       └── main.go              # CLI entry point (Cobra CLI)
├── internal/
│   ├── stt/
│   │   ├── whisper.go           # whisper.cpp Go bindings wrapper
│   │   ├── model.go             # Model management (download, load, select)
│   │   └── transcribe.go        # Core transcription logic
│   ├── vad/
│   │   ├── vad.go               # VAD interface
│   │   ├── onnx.go              # ONNX Runtime implementation
│   │   └── timestamps.go        # Speech segment extraction
│   ├── audio/
│   │   ├── capture.go           # Microphone input (portaudio/miniaudio)
│   │   ├── convert.go           # FFmpeg wrapper (any format → 16kHz WAV)
│   │   └── chunk.go             # Segment audio into 30s chunks
│   ├── embed/
│   │   └── embed.go             # Embedding generation (call local model or API)
│   ├── output/
│   │   ├── text.go              # Plain text formatter
│   │   ├── json.go              # JSON with timestamps, confidence
│   │   └── srt.go               # SubRip subtitle output
│   ├── store/
│   │   ├── db.go                # libSQL connection (embedded mode)
│   │   ├── schema.go            # Table creation, migrations
│   │   └── transcript.go        # CRUD for transcripts + vectors
│   ├── search/
│   │   └── search.go            # Semantic search via vector_top_k
│   └── config/
│       └── config.go            # CLI config, env vars, defaults
├── models/                      # Downloaded model files (gitignored)
│   ├── whisper/                 # ggml-base.en.bin, etc.
│   ├── vad/                     # silero_vad.onnx
│   └── embed/                   # all-MiniLM-L6-v2 ONNX
├── scripts/
│   ├── download-models.sh       # Model download helper
│   └── build-whisper.sh         # Build whisper.cpp from source
├── go.mod
├── go.sum
├── Makefile                     # build, test, setup, clean
└── PLAN.md                      # This file
```

---

## Phase 1: Core STT Engine

**Goal:** Transcribe an audio file using whisper.cpp.

### Tasks
1. Clone `whisper.cpp` and build from source (`cmake -B build && cmake --build build`)
2. Download `base.en` model via `models/download-ggml-model.sh`
3. Use official Go bindings: `github.com/ggerganov/whisper.cpp/bindings/go`
4. Create `internal/stt/` package:
   - `whisper.go` — init context, load model
   - `transcribe.go` — `Transcribe(audioPath string) (*Transcript, error)`
   - `Transcript` struct: `Segments []Segment`, each with `Text`, `Start`, `End`, `Confidence`
5. Handle audio format: input must be **16-bit WAV, 16kHz, mono**
6. Create `internal/audio/convert.go` — shell out to FFmpeg for conversion
7. Write a simple CLI command: `whispr transcribe test.wav`
8. Benchmark on your hardware — target: ~0.5x real-time for `base.en`

### Deliverables
- [ ] `whispr transcribe audio.wav` works end-to-end
- [ ] Output as plain text
- [ ] Benchmark numbers documented

---

## Phase 2: Voice Activity Detection

**Goal:** Skip silence to reduce STT compute by 30-60%.

### Tasks
1. Download `silero_vad.onnx` (~2MB) from the silero-vad repo
2. Integrate ONNX Runtime for Go: `github.com/yalue/onnxruntime_go`
3. Create `internal/vad/` package:
   - `vad.go` — load ONNX model, init session
   - `timestamps.go` — `DetectSpeech(audioPath string) ([]SpeechSegment, error)`
   - `SpeechSegment`: `Start`, `End` (in seconds)
4. Pipeline: raw audio → VAD → extract speech segments → feed each to STT
5. Handle edge cases: segments too short, overlapping, gaps
6. Add `--vad` flag to enable/disable VAD pre-filtering

### Deliverables
- [ ] `whispr transcribe audio.wav --vad` skips silence
- [ ] Measurable speedup vs. non-VAD run
- [ ] VAD timestamps included in JSON output

---

## Phase 3: Audio Input Layer

**Goal:** Support both file input and live microphone capture.

### Tasks
1. **File mode** (already partially done in Phase 1):
   - Support `.wav`, `.mp3`, `.m4a`, `.ogg`, `.flac`
   - Auto-convert to 16kHz 16-bit mono WAV via FFmpeg
   - Clean up temp files after processing

2. **Live mode** (new):
   - Integrate `portaudio` or `miniaudio` via CGO
   - Capture from default microphone
   - Buffer audio in 30s chunks
   - Stream chunks through VAD → STT pipeline
   - Print transcripts in real-time as they complete
   - Handle graceful shutdown (Ctrl+C)

3. **Chunking logic** (`internal/audio/chunk.go`):
   - Split long audio into ≤30s segments
   - Overlap segments by 100ms to avoid cutting words
   - Track global timestamps across chunks

### Deliverables
- [ ] `whispr transcribe audio.mp3` (any format)
- [ ] `whispr transcribe --live` (real-time mic)
- [ ] Real-time streaming output

---

## Phase 4: Output Formats & Post-Processing

**Goal:** Rich output formats and optional LLM summarization.

### Tasks
1. **Output formatters** (`internal/output/`):
   - `text.go` — plain text (default)
   - `json.go` — structured JSON with segments, timestamps, confidence
   - `srt.go` — SubRip subtitle format
   - `vtt.go` — WebVTT format (optional)
   - `--output` flag to select format

2. **Embedding generation** (`internal/embed/`):
   - Use `all-MiniLM-L6-v2` ONNX model (~80MB)
   - Generate 384-dim embeddings for each transcript segment
   - Store alongside transcript in libSQL

3. **LLM summarization** (optional, requires Ollama running):
   - `--summarize` flag
   - Send transcript to `localhost:11434/api/generate` (Phi-3, Mistral, etc.)
   - Print summary after transcription

### Deliverables
- [ ] `whispr transcribe audio.wav --output json`
- [ ] `whispr transcribe audio.wav --output srt`
- [ ] `whispr transcribe audio.wav --summarize`
- [ ] Embeddings generated and stored

---

## Phase 5: Storage with libSQL (Turso)

**Goal:** Persist transcripts with vector embeddings for future search.

### Tasks
1. Set up libSQL in embedded mode:
   ```go
   import "github.com/tursodatabase/go-libsql"
   connector := libsql.NewEmbeddedConnector("./whisprflow.db")
   db := sql.OpenDB(connector)
   ```

2. Create schema (`internal/store/schema.go`):
   ```sql
   CREATE TABLE IF NOT EXISTS transcripts (
       id         INTEGER PRIMARY KEY,
       title      TEXT NOT NULL,
       source     TEXT,              -- file path or "live"
       created_at TEXT NOT NULL,
       duration   REAL,
       model      TEXT,              -- whisper model used
       full_text  TEXT,              -- concatenated transcript
       embedding  F32_BLOB(384)      -- all-MiniLM-L6-v2 embedding
   );

   CREATE TABLE IF NOT EXISTS segments (
       id             INTEGER PRIMARY KEY,
       transcript_id  INTEGER REFERENCES transcripts(id),
       start          REAL,
       "end"          REAL,
       text           TEXT,
       confidence     REAL,
       embedding      F32_BLOB(384)
   );

   -- Vector indexes
   CREATE INDEX IF NOT EXISTS transcripts_vec_idx
       ON transcripts (libsql_vector_idx(embedding));

   CREATE INDEX IF NOT EXISTS segments_vec_idx
       ON segments (libsql_vector_idx(embedding));
   ```

3. CRUD operations (`internal/store/transcript.go`):
   - `SaveTranscript(t *Transcript) error`
   - `GetTranscript(id int) (*Transcript, error)`
   - `ListTranscripts() ([]TranscriptMeta, error)`
   - `DeleteTranscript(id int) error`

4. Auto-save on every transcription
5. `whispr list` and `whispr show <id>` commands

### Deliverables
- [ ] `whisprflow.db` created automatically on first run
- [ ] Transcripts auto-saved with embeddings
- [ ] `whispr list` shows history
- [ ] `whispr show <id>` shows full transcript

---

## Phase 6: Semantic Search

**Goal:** Search past transcripts by meaning, not just keywords.

### Tasks
1. Create vector index on embeddings (done in Phase 5 schema)
2. Implement semantic search (`internal/search/search.go`):
   ```sql
   SELECT t.title, t.full_text,
          vector_distance_cos(t.embedding, vector32(?)) AS distance
   FROM transcripts t
   ORDER BY distance ASC
   LIMIT ?;
   ```
3. Query embedding for search text:
   - Generate embedding for the search query using same model
   - Run `vector_top_k` or cosine distance query
   - Return top N matching transcripts

4. CLI command:
   ```
   whispr search "what did I say about the deployment issue"
   ```

5. Optional: keyword search fallback (FTS5 on `full_text` column)

### Deliverables
- [ ] `whispr search "query"` returns relevant past transcripts
- [ ] Results ranked by similarity
- [ ] Fallback keyword search (optional)

---

## Phase 7: CLI Polish & DX

**Goal:** Professional CLI experience.

### Tasks
1. Use [Cobra](https://github.com/spf13/cobra) for CLI structure:
   ```
   whispr transcribe <file>     # Transcribe a file
   whispr transcribe --live     # Live mic input
   whispr list                  # List past transcripts
   whispr show <id>             # Show a specific transcript
   whispr search "<query>"      # Semantic search
   whispr models                # List/download models
   whispr config                # Show/set configuration
   ```

2. Global flags:
   - `--model <base|small|medium>` — whisper model size
   - `--output <text|json|srt>` — output format
   - `--vad` — enable VAD pre-filtering
   - `--threads <n>` — CPU threads (default: 6)
   - `--language <en|auto>` — language detection
   - `--summarize` — run through Ollama LLM
   - `--no-save` — skip database storage

3. Config file (`~/.config/whisprflow/config.yaml`):
   - Default model
   - Default output format
   - Thread count
   - Ollama endpoint
   - Database path

4. Progress indicators:
   - Spinner during transcription
   - Real-time segment output for live mode
   - ETA for long files

5. Error handling:
   - Friendly error messages
   - Suggest fixes (e.g., "FFmpeg not found, install with: brew install ffmpeg")

### Deliverables
- [ ] All commands work with `--help`
- [ ] Config file support
- [ ] Progress indicators
- [ ] Clean error messages

---

## Phase 8: Testing & Benchmarking

**Goal:** Reliable, well-tested codebase.

### Tasks
1. Unit tests for each package:
   - `stt/` — mock whisper, test parsing
   - `vad/` — test with known audio files (silence vs. speech)
   - `audio/` — test format conversion
   - `output/` — test format correctness
   - `store/` — test CRUD with temp database
   - `search/` — test vector queries

2. Integration tests:
   - Full pipeline: audio file → transcript → store → search
   - Use a sample audio file from `testdata/`

3. Benchmarks:
   - Transcription speed (real-time factor) for each model
   - VAD overhead
   - Embedding generation time
   - Search latency

4. CI/CD (optional):
   - GitHub Actions: build, test, lint
   - Release binaries for Linux/macOS

### Deliverables
- [ ] `go test ./...` passes
- [ ] `go test -bench=. ./...` produces benchmark numbers
- [ ] Coverage > 70%

---

## Build Order & Milestones

| Phase | What | Estimated Effort |
|-------|------|-----------------|
| **1** | Core STT — transcribe a file | First milestone |
| **2** | VAD — skip silence | Speedup visible |
| **3** | Live mic input | Real-time demo |
| **4** | Output formats + embeddings | Rich output |
| **5** | libSQL storage | Persistence |
| **6** | Semantic search | AI-powered search |
| **7** | CLI polish | Professional DX |
| **8** | Tests & benchmarks | Quality gate |

---

## Key Dependencies

```go
// go.mod
require (
    github.com/ggerganov/whisper.cpp/bindings/go  // STT
    github.com/tursodatabase/go-libsql             // Database + vectors
    github.com/yalue/onnxruntime_go               // VAD (ONNX)
    github.com/spf13/cobra                        // CLI
    github.com/spf13/viper                        // Config
    github.com/schollz/progressbar/v3             // Progress bars
)
```

### External Dependencies
- **FFmpeg** — audio format conversion
- **CMake + C compiler** — build whisper.cpp
- **ONNX Runtime** — for silero-vad (download .so/.dll)
- **Ollama** (optional) — local LLM for summarization

---

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| whisper.cpp CGO build fails | Medium | High | Shell out to CLI as fallback |
| ONNX Runtime CGO complexity | Medium | Medium | Python subprocess as fallback |
| Real-time latency on CPU | Low | High | Use `base` model, VAD first |
| AMD GPU = no CUDA | N/A | Low | whisper.cpp uses BLAS, not CUDA |
| libSQL Go driver maturity | Low | Medium | Fallback to `modernc.org/sqlite` + sqlite-vec |
| Portaudio build issues | Medium | Medium | Use miniaudio (single-header, simpler CGO) |
| Memory pressure with large files | Low | Medium | Process in chunks, GC between segments |

---

## Model Reference

### Whisper Models
| Model | Size | RAM | Speed (relative) | Best For |
|-------|------|-----|------------------|----------|
| `tiny.en` | ~75MB | ~300MB | 1x | Quick tests |
| `base.en` | ~140MB | ~500MB | 0.5x | **Default for your hardware** |
| `small.en` | ~460MB | ~1GB | 0.2x | Better accuracy |
| `medium.en` | ~1.5GB | ~3GB | 0.1x | High accuracy, slow |

### VAD Model
- `silero_vad.onnx` — ~2MB, ~200MB RAM during inference

### Embedding Model
- `all-MiniLM-L6-v2` (ONNX) — ~80MB, 384-dim vectors

---

## Quick Start Commands (Target UX)

```bash
# Setup
whispr models download base.en
whispr models download vad
whispr models download embed

# Basic transcription
whispr transcribe meeting.mp3

# With VAD, JSON output
whispr transcribe meeting.mp3 --vad --output json

# Live microphone
whispr transcribe --live --model base.en

# With summarization
whispr transcribe meeting.mp3 --summarize

# List past transcripts
whispr list

# Show specific transcript
whispr show 3

# Semantic search
whispr search "what did we decide about the API design?"

# Config
whispr config set model small.en
whispr config set threads 6
whispr config set output json
```

---

## Next Steps

1. **Approve this plan** — review and adjust as needed
2. **Start Phase 1** — scaffold project, build whisper.cpp, transcribe first file
3. **Iterate** — each phase builds on the previous one
