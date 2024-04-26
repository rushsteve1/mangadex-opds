# MangaDex OPDS

This software is a [MangaDex](https://mangadex.org) to [OPDS](https://opds.io)
proxy for reading manga on apps/devices that support OPDS such as e-readers.

It is meant to be self-hosted by the person reading.
This proxy is **NOT** suitable public-internet deployments.
We want to be respectful of MangaDex and their
[API Usage Policy](https://api.mangadex.org/docs/#acceptable-usage-policy)

## Project Goals and Non-Goals

- Easy to self-host
- Easy to understand and hack on
- Easy to configure with env vars
- Minimal dependencies (just Go stdlib and `golang.org/x/`)
- Reasonably performant
- Minimal system requirements
- Implement the OPDS standard closely
- Provide chapters in common formats (EPUB, CBZ)

- **NOT** a standalone manga reading server
- **NOT** a publicly available service
