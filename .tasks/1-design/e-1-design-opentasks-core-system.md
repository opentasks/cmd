---
id: 1
title: Design OpenTasks Core System
type: epic
status: done
tags: [design, architecture]
relationships:
  - type: blocks
    taskID: 5
  - type: blocks
    taskID: 6
  - type: blocks
    taskID: 7
  - type: blocks
    taskID: 12

createdAt: 2025-11-02T08:00:00Z
updatedAt: 2025-11-02T21:00:00Z
---

# Design OpenTasks Core System

Complete architectural design for the OpenTasks markdown-based task management system. This epic encompasses the core data model, storage interface, configuration system, and query capabilities.

## Overview

OpenTasks is a markdown-based task manager written in Go with pluggable storage backends, customizable workflows, and semantic task IDs. The design focuses on simplicity, extensibility, and dog-fooding the system itself for internal task management.

## Scope

- Task data model and relationships
- Semantic ID system with collision detection
- Configuration system (config.toml structure)
- Storage plugin interface (BaseStorage)
- Query engine for filtering and retrieval
- CLI architecture (viper/cobra)
- MCP integration for AI agents

## Success Criteria

- [x] All design documents completed and reviewed
- [x] BaseStorage interface fully specified
- [x] config.toml schema documented
- [x] ID generation strategy validated
- [x] CLI architecture designed and validated
- [x] Ready for implementation phase (Phase 2)

## Status: Complete ✅

All core design tasks completed:
- s-5: Task Data Model and Relationships ✅
- s-6: Semantic ID System ✅
- s-7: BaseStorage Interface ✅
- s-8: Configuration System ✅
- s-9: Query Engine and Filtering ✅
- s-10: CLI Architecture (reviewed and validated) ✅
- d-0: Final Design Decisions ✅
- r-2,3,4: Research documents ✅
- p-11: Design Roadmap ✅

## Next Phase: Testing & Polish (Phase 2)

Phase 2 epic (e-12) now begins implementation of:
- Unit tests (s-13)
- Integration tests (s-14)
- Feature implementations (s-15 through s-25)
- Bug fixes and enhancements (b-21)
- End-to-end testing (t-19)
