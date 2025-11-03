---
id: 26
title: Write epic folder structure with computed ID and summarized name
type: task
status: todo
tags:
    - testing
relationships:
    - type: parent
      taskID: 12
createdAt: "2025-11-02T08:41:05Z"
updatedAt: "2025-11-02T20:35:00Z"
---

# Epic folder structure with computed ID and summarized name


## Instruction
When a user creates a new epic task, the system should automatically organize it into a folder structure using the computed epic ID and a summarized version of the title.

## User Story
As a user, when I create a new epic, I should see the new epic in a folder that matches the computed ID and summarized name:
`<id>-<summarized-title>/<e>-<id>-<summarized-title>.md`

## Example
Creating an epic with title "Implement user authentication system" might result in:
- Folder: `13-implement-user-auth/`
- File: `e-13-implement-user-auth.md`

## Requirements
1. **ID Computation**: Generate sequential epic ID (13, 14, 15, etc.)
2. **Title Summarization**: Convert long titles to slug format (lowercase, hyphens, max ~20 chars)
3. **Folder Creation**: Create `<id>-<summarized-title>/` directory automatically
4. **File Placement**: Save epic as `<folder>/e-<id>-<summarized-title>.md`
5. **Storage Integration**: Update MarkdownFileStorage to handle epic folder creation
6. **ID Sequence Management**: Track highest epic ID to compute next ID correctly

## Tasks
- [ ] Analyze current epic ID computation in codebase
- [ ] Implement title summarization/slugification logic
- [ ] Update MarkdownFileStorage to create epic folders
- [ ] Update ID sequencing to handle epic-specific IDs
- [ ] Write unit tests for title summarization
- [ ] Write integration tests for epic folder creation
- [ ] Verify existing epics can be migrated if needed
- [ ] Test edge cases (very long titles, special chars, unicode)

## Deliverable
- Epic creation flow that automatically creates folder structure
- Updated MarkdownFileStorage with epic folder handling
- Title summarization utility function
- Unit tests for summarization logic (90%+ coverage)
- Integration tests for full epic creation flow
- Documentation of epic folder structure format

## Log
