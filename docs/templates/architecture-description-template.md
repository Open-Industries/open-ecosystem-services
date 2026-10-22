# Architecture Description Template

> **How to use:** Copy this file to `docs/architecture.md` (or the repository's established architecture path) and replace all bracketed prompts. Mark a section `Not applicable` with a short reason when it does not apply. Link to versioned diagrams, ADRs, API contracts, IaC, runbooks, and test evidence instead of duplicating their contents.
>
> This is an adaptable template informed by ISO/IEC/IEEE 42010:2022, TOGAF Standard 10th Edition, and IEEE 1016-2009 for software design detail. It is not a claim of certification or strict compliance. Tailor it to the system and stakeholders.

## Document control

| Field | Value |
|---|---|
| System / product / service | [Name] |
| Document owner / lead architect | [Name or role] |
| Version / status | [0.1 / Draft] |
| Date / change | [YYYY-MM-DD / Change or issue ID] |
| Reviewers / approvers | [Roles or names] |
| Classification | [Public / Internal / Confidential / other] |
| Related records | [Requirements, ADRs, tickets, source repository, model links] |

### Change history

| Version | Date | Author | Summary | Approval / review |
|---|---|---|---|---|
| 0.1 | [YYYY-MM-DD] | [Name] | Initial draft | [Pending] |

## 1. Executive summary

- **Problem / opportunity:** [What outcome is required?]
- **Proposed solution:** [Describe the architecture in 3–5 sentences.]
- **Expected value:** [Business or operational outcomes.]
- **Key trade-offs:** [Most important trade-offs.]
- **Decision / approval needed:** [Owner and due date, or None.]
- **Top risks:** [Link to risk records or summarize.] 

## 2. Purpose, scope, and context

- **Purpose and audience:** [Why this description exists and who uses it.]
- **System of interest:** [System, application, platform, service, or change being described.]
- **In scope:** [Capabilities, teams, environments, and lifecycle stages.]
- **Out of scope:** [Explicit exclusions.]
- **Current / target state:** [Baseline and intended end state; link to diagrams if applicable.]
- **Assumptions:** [Assumption, owner, validation date.]
- **Constraints:** [Regulatory, organizational, technology, hosting, budget, or schedule constraints.]
- **Dependencies:** [External teams, systems, vendors, services, and required dates.]

### Terms and abbreviations

| Term | Meaning |
|---|---|
| [Term] | [Definition] |

## 3. Stakeholders and concerns

Record the people or roles that have an interest in the system and the questions the architecture must address.

| Stakeholder / role | Concern / question | Priority | Evidence / view that addresses it |
|---|---|---|---|
| [Service owner] | [Which capabilities and service levels are delivered?] | [High] | [Context, capability, SLO] |
| [Security / privacy] | [Where are trust boundaries and sensitive-data flows?] | [High] | [Security, identity, data-flow] |
| [Platform / operations] | [How is it deployed, monitored, backed up, and recovered?] | [High] | [Deployment and operations] |
| [Development / integration] | [What are the components, contracts, and failure behaviors?] | [High] | [Logical, interface, sequence] |

## 4. Drivers, requirements, and quality attributes

Link to authoritative requirements where possible. State measurable acceptance criteria rather than adjectives such as “fast,” “secure,” or “highly available.”

| ID | Driver / requirement | Target or acceptance criterion | Verification evidence / owner |
|---|---|---|---|
| [NFR-01] | [Availability] | [e.g., monthly API availability ≥ 99.95%, measured at the public endpoint] | [SLO report / owner] |
| [NFR-02] | [Performance] | [e.g., p95 latency ≤ 300 ms at stated load] | [Load test / owner] |

## 5. Principles and design constraints

| ID | Principle / constraint | Rationale and implications | Type / owner |
|---|---|---|---|
| [P-01] | [Example: externalize environment-specific configuration] | [Reason and effect on design] | [Mandatory or preferred / role] |

## 6. Architecture overview

[Insert a high-level context diagram. Show actors, external systems, system boundary, trust boundaries, and major interfaces. Include a legend and source link.]

### Building blocks and responsibilities

| Element | Responsibility | Owner | Critical dependencies |
|---|---|---|---|
| [API gateway] | [Authentication, routing, quotas] | [Platform team] | [Identity provider] |

### End-to-end flow

1. [Describe the primary user or system flow.]
2. [Describe validation, authorization, and data transformations.]
3. [Describe processing, events, and response.]
4. [Describe failure handling and audit/telemetry.]

## 7. Viewpoints and architecture views

A **viewpoint** defines how to build and interpret a view to address stakeholder concerns. For each view, name the viewpoint, intended stakeholders, concerns, notation, and validation rules. Include diagrams as Mermaid, PlantUML, C4, UML, ArchiMate, or a linked model as appropriate.

### 7.1 Context / ecosystem view

- **Viewpoint and stakeholders:** [Name roles]
- **Concerns addressed:** [Boundary, actors, external dependencies]
- **Diagram / model:** [Insert or link]
- **Assumptions / validation:** [Explain]

### 7.2 Functional / logical view

- **Viewpoint and stakeholders:** [Name roles]
- **Concerns addressed:** [Responsibilities, component boundaries, interactions]
- **Diagram / model:** [Insert or link]
- **Assumptions / validation:** [Explain]

### 7.3 Information / data view

- **Viewpoint and stakeholders:** [Name roles]
- **Concerns addressed:** [Data ownership, classification, flow, retention, residency, lifecycle]
- **Diagram / model:** [Insert or link]
- **Assumptions / validation:** [Explain]

### 7.4 Interface / integration view

- **Viewpoint and stakeholders:** [Name roles]
- **Concerns addressed:** [APIs, events, protocols, schemas, versions, errors]
- **Contracts:** [OpenAPI / AsyncAPI / schema links]
- **Assumptions / validation:** [Explain]

### 7.5 Deployment / technology view

- **Viewpoint and stakeholders:** [Name roles]
- **Concerns addressed:** [Runtime, clusters, networks, zones, placement, capacity]
- **Diagram / model:** [Insert or link]
- **Environment differences:** [Development, test, staging, production]

### 7.6 Security / privacy view

- **Viewpoint and stakeholders:** [Name roles]
- **Concerns addressed:** [Identity, authorization, trust zones, sensitive data, threats]
- **Diagram / model:** [Insert data-flow or threat-model link]
- **Controls and residual risks:** [Link to security review and accepted exceptions]

### 7.7 Operations / resilience view

- **Viewpoint and stakeholders:** [Name roles]
- **Concerns addressed:** [SLOs, scaling, failure domains, observability, backup, DR]
- **Diagram / model:** [Insert or link]
- **Runbooks / evidence:** [Links]

### 7.8 Migration / transition view

- **Viewpoint and stakeholders:** [Name roles]
- **Concerns addressed:** [Migration increments, coexistence, cutover, rollback, retirement]
- **Roadmap / diagram:** [Insert or link]

Delete or add view sections based on actual stakeholder concerns. Do not keep empty views simply to complete a checklist.

## 8. Interfaces, contracts, and dependencies

| Interface / dependency | Provider → consumer | Contract / protocol | Identity and security | Timeout / retry / failure behavior | Owner |
|---|---|---|---|---|---|
| [Name] | [A → B] | [OpenAPI, AsyncAPI, HTTPS, event] | [AuthN/Z, TLS, data class] | [Behavior and SLO] | [Team] |

## 9. Architecture decisions and rationale

Record significant choices in separate, version-controlled ADRs. Link each decision to affected requirements and views.

| ADR | Decision | Status | Consequences / revisit trigger |
|---|---|---|---|
| [docs/adrs/0001-example.md] | [Decision summary] | [Proposed / Accepted / Superseded] | [Trade-offs and trigger] |

Minimum ADR content: context and drivers; decision; alternatives considered; positive and negative consequences; security, operations, cost, and compliance impact; owner and date.

## 10. Security, privacy, and compliance

- **Data classification and flows:** [Link to data view and owner.]
- **Identity and access:** [Human, workload, and external identities; authentication and authorization.]
- **Trust boundaries and network policy:** [Zones, ingress/egress, service-to-service controls.]
- **Encryption and secrets:** [In transit, at rest, key and certificate lifecycle.]
- **Threats and mitigations:** [Threat model, residual risks, security owner.]
- **Audit and privacy:** [Events, retention, consent, data minimization, residency.]
- **Applicable controls / exceptions:** [Policy, regulation, standard, evidence, exception owner.]

## 11. Deployment and operations

- **Runtime and deployment units:** [Kubernetes/OpenShift resources, VMs, managed services, etc.]
- **Capacity assumptions:** [CPU, memory, storage, throughput, limits, evidence.]
- **Availability and failure handling:** [Redundancy, zones, timeouts, retries, circuit breakers, degradation.]
- **Observability:** [Metrics, logs, traces, alerts, dashboards, correlation identifiers.]
- **Backup and disaster recovery:** [RPO, RTO, restore procedure, test frequency.]
- **Release and rollback:** [Versioning, promotion, compatibility, rollback conditions.]
- **Support and ownership:** [Service owner, on-call, runbooks, escalation.]
- **Operational acceptance criteria:** [Objective evidence required before release.]

## 12. Implementation and transition plan

| Increment / work package | Scope | Dependencies | Acceptance evidence | Owner / target |
|---|---|---|---|---|
| [1] | [Foundation] | [Identity, network] | [Tests / review] | [Role / date] |

Describe migration sequencing, data migration, coexistence, cutover, rollback, and decommissioning. Link implementation code, IaC, GitOps manifests, and release plan.

## 13. Risks, issues, assumptions, and dependencies

| ID | Type | Description | Impact / likelihood | Mitigation / action | Owner / status |
|---|---|---|---|---|---|
| [R-01] | [Risk] | [Description] | [High / Medium] | [Action] | [Role / Open] |

## 14. Traceability and consistency

| Concern / requirement | View(s) / ADR(s) | Validation evidence | Status / accepted gap |
|---|---|---|---|
| [Availability] | [Operations view, ADR-02] | [Failover test link] | [Met / gap owner] |

Check that each high-priority concern maps to a view or explicit decision; views use defined viewpoints; diagrams agree on element names and relationships; interface contracts match the described flows; and each key requirement maps to validation evidence or an accepted gap.

## 15. Review and approval

| Review role | Name | Outcome / conditions | Date |
|---|---|---|---|
| Architecture | [Name] | [Approved / changes requested] | [YYYY-MM-DD] |
| Security / privacy | [Name] | [Outcome] | [YYYY-MM-DD] |
| Service owner / operations | [Name] | [Outcome] | [YYYY-MM-DD] |

## References

- [ISO/IEC/IEEE 42010:2022 — Architecture description](https://www.iso.org/standard/74393.html)
- [ISO/IEC/IEEE 42020:2019 — Architecture processes](https://www.iso.org/standard/68982.html)
- [ISO/IEC/IEEE 15288:2023 — System life cycle processes](https://www.iso.org/standard/81702.html)
- [TOGAF Standard, 10th Edition](https://www.opengroup.org/togaf)
- [TOGAF example deliverable templates](https://help.opengroup.org/hc/en-us/articles/21726647171730-Are-There-Any-Template-Deliverables-for-the-TOGAF-Standard)
- [IEEE 1016-2009 — Software Design Descriptions](https://standards.ieee.org/standard/1016-2009.html)

## Tailoring and TOGAF mapping

This document can serve as the architecture description / architecture definition work product for a solution. Tailor content to the engagement; it does not replace all ADM artifacts.

| TOGAF ADM area | Sections to use |
|---|---|
| Preliminary / capability | Document control, principles, constraints, governance |
| Phase A — Architecture Vision | Executive summary, purpose, scope, stakeholders, drivers |
| Phases B–D — Business, Information Systems, Technology | Requirements, principles, views, building blocks, interfaces, decisions |
| Phase E — Opportunities & Solutions | Implementation increments, risks, transition architecture |
| Phase F — Migration Planning | Transition plan and dependencies |
| Phase G — Implementation Governance | Validation, approvals, exceptions, traceability |
| Phase H — Architecture Change Management | Change history, impact assessment, review triggers |
| Requirements Management | Traceability and consistency |
