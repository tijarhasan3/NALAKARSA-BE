# Discussions Sort Filter Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans (recommended for this small change) to implement this plan task-by-task.

**Goal:** Make `GET /api/v1/discussions` honor the FE sort contract without changing its endpoint or response shape.

**Architecture:** Keep filtering and pagination in the existing GORM repository query. Map `newest`, `oldest`, and `popular` to SQL ordering and always add `created_at DESC` as the deterministic tie-breaker; unknown or empty values use `newest`.

**Tech Stack:** Go, GORM, PostgreSQL, existing repository tests/tooling.

**Spec:** `C:/Users/ASUS/Downloads/BE_DISCUSSIONS_SORT_FILTER_CONTRACT.md`

## Global Constraints

- Preserve `search`, `role`, `category`, `page`, and `limit` filters.
- Do not add an endpoint.
- `popular` must use the stored actual `upvote_count` field.
- Response and pagination metadata remain unchanged.

## Review Focus

- `sort=popular` orders by `upvote_count DESC`, then `created_at DESC`.
- `sort=oldest` orders by `created_at ASC`, then `created_at DESC` as the required fallback.
- Empty or unknown sort values safely default to newest.
- Existing filters remain combined with sorting before pagination/counting.

### Task 1: Discussion ordering

**Files:**
- Modify: `internal/repository/discussion/discussion_repository.go`
- Test: `internal/repository/discussion/discussion_repository_test.go` (or the repository's existing test location)

- [ ] Add a focused test or query assertion covering `newest`, `oldest`, `popular`, and fallback ordering.
- [ ] Update `pgDiscussionRepository.List` to map the three contract values and append `discussions.created_at DESC` as the stable secondary order.
- [ ] Run `gofmt`, `go test ./...`, and `git diff --check`.
