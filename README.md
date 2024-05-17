# MangaDex-OPDS

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
- Minimal dependencies
- Reasonably performant
- Minimal system requirements
- Implement the OPDS standard closely
- Provide chapters in common formats (EPUB, CBZ)

- **NOT** a standalone manga reading server
- **NOT** a publicly available service

## Specifications

- [OPDS](https://specs.opds.io/opds-1.2)
- [ATOM](https://validator.w3.org/feed/docs/atom.html)
- [OPDS-PSE](https://anansi-project.github.io/docs/opds-pse/specs/v1.0)
- [EPUB](https://www.w3.org/TR/epub-33/)
- [ComicInfo.xml](hhttps://anansi-project.github.io/docs/comicinfo/documentation)

## Thanks

- The [MangaDex Team](https://mangadex.org/about)
- Dani and the [Panels team](https://panels.app)

## Copyright

MangaDex and the MangaDex Logo are the copyright of [MangaDex](https://mangadex.org/about)

mangadex-opds is licensed under the terms of the [AGPL license](./LICENSE.txt)

## Donations

This project does not take donations, please
[donate to MangaDex instead](https://namicomi.com/en/org/3Hb7HnWG/mangadex/subscriptions)
