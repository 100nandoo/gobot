# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

## Before exploring, read these

- `CONTEXT.md` at the repo root, or
- `CONTEXT-MAP.md` at the repo root if it exists
- `docs/adr/` — read ADRs that touch the area you're about to work in

If any of these files don't exist, proceed silently.

## File structure

This repo is configured as a single-context repo.

## Use the glossary's vocabulary

When naming domain concepts, use the terms defined in `CONTEXT.md`.

## Flag ADR conflicts

If an output contradicts an existing ADR, surface it explicitly rather than silently overriding it.
