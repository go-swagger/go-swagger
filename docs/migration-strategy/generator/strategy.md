Path to complete implementation
- Alpha (MVP slice, deliver developer-preview)
  - Complete: S0.2, S1.1, S1.4, S1.5, S2.1, S2.4, S3.1, S3.3, S4.3, S7.1.
  - Output: Compiling server scaffolding for the MVP spec; parameter binding for form/simple styles; JSON requestBody/response; v3 templates; CLI flagging.
  - Quality: Basic unit tests and one golden fixture; v2 regression job (smoke).

- Beta (broaden feature coverage and client parity)
  - Add: S1.2 (deref with anchors), S2.2–S2.3 (null/union/composition), S2.6–S2.7 (enum/format policy), S3.4–S3.5 (response content, negotiation), S4.1–S4.2 (servers + security), S5.1–S5.2 (client surface/readers), S7.2–S7.4 (helpers/CLI/naming), S8.3 (lints), S9.1 (goldens), S9.5 (dual-lane CI).
  - Output: Functional client SDK; multiple media types in responses; security hooks; server URL selection; richer validations; robust templates and naming.
  - Quality: Golden tests across varied fixtures; dual-lane CI green.

- GA (completeness, performance, and docs)
  - Add: S2.5 (runtime validator integration), S3.2 (client encoders), S5.3 (client auth), S1.3 (dialect handling), S8.2 (migration guide), S9.2–S9.4 (schema suite, fuzzing, perf), S6.x (callbacks/webhooks) if in-scope.
  - Output: Production-ready generator for OAS3.1; documented migration and known limitations; performance baselines.
  - Quality: JSON Schema 2020-12 subset passing; fuzz-stable; benchmarks tracked.

Notes on estimates and priorities
- P0: MVP-critical or infra that unblocks many items.
- P1: Beta-critical; rounds out server+client parity.
- P2: Important but can follow GA.
- P3: Nice-to-have/advanced.
