<p align="center">
  <a href="https://opencomply.io">
    <picture style="width":50%">
          <source media="(prefers-color-scheme: dark)" srcset="https://github.com/opengovern/opencomply/blob/main/assets/logos/logo-dark.png">
          <source media="(prefers-color-scheme: light)" srcset="https://github.com/opengovern/opencomply/blob/main/assets/logos/logo-light.png">
          <img style="width":50%" alt="opencomply" src="https://github.com/opengovern/opencomply/blob/main/assets/logos/logo-light.png" >
    </picture>

  </a>
</p>

Take control of your security and compliance across all clouds, platforms, and regions. OpenComply makes it easy to govern every change, deployment, and asset—so you can build exceptional products without getting bogged down by complex compliance processes.

![App Screenshot](./assets/screenshots/app-%20screenshot%203.png)

## Table of Contents
- [Key Features](#key-features)
- [Example: Query Everything](#example-query-everything)
- [Supported Integrations](#supported-integrations)
- [Getting Started](#getting-started)
- [Next Steps](#next-steps)

---

## Key Features

OpenComply is built by practitioners for practitioners, with a focus on making security and compliance **accessible**, **agile**, and **inclusive**:

- **Accessibility:** Teams of any size can easily adopt compliance best practices.
- **Agility:** Adapt quickly to evolving infrastructure, policies, and frameworks.
- **Inclusivity:** Foster collaboration across DevOps, Security, and Compliance teams.

Here’s how OpenComply delivers on these principles:

- **Full-Stack Security & Compliance:** Get a single source of truth across all your environments—regardless of cloud provider, region, or platform.
- **Policy as Code:** Define and centralize security policies using version-controlled code. Continuously audit changes to maintain compliance and strengthen your security posture.
- **Automated Audits:** Detect risks and promote best practices automatically with real-time scanning and alerts.
- **Highly Customizable:** Leverage flexible frameworks, controls, roles, and integrations to tailor OpenComply to your organization’s unique needs.
- **SQL for Your Entire Tech Stack:** Query your infrastructure and code artifacts as if they were tables in a database—enabling fast insights and powerful automation.

---

## Example: Query Everything

OpenComply supports rich SQL queries across your entire tech stack. For example, find all unique Docker base images and usage counts:

```sql
SELECT image AS "Base Image", COUNT(*) AS "Count"
FROM (
  SELECT DISTINCT sha, jsonb_array_elements_text(images) AS image
  FROM github_artifact_dockerfile
) AS expanded
GROUP BY image
ORDER BY "Count" DESC;
```

### Sample Results

| Base Image                  | Count |
|-----------------------------|-------|
| scratch                     | 14    |
| docker.io/golang:alpine     | 14    |
| cloudql-plugin-base:0.0.161 | —     |
| golang:1.23-alpine          | 4     |
| alpine:latest               | 3     |

OpenComply turns your infrastructure, code, and security data into a searchable database.

## Supported Integrations

Seamlessly integrate OpenComply with your favorite tools and platforms, including but not limited to:

- **Public Clouds:** AWS, Azure, DigitalOcean, Linode (Akamai)
- **AI Services:** OpenAI, CohereAI
- **Identity Providers:** Microsoft Entra ID (Azure AD)
- **DevOps & Code Security:** GitHub
- **Web Security:** Cloudflare

For a full list, check out our Integrations Documentation. You can also easily write your own custom integrations.

## Getting Started

### Installation via Helm

You can install OpenComply on any Kubernetes cluster with Helm:

```bash
helm repo add opencomply https://charts.opencomply.io --force-update
helm install -n opencomply opencomply opencomply/opencomply \
  --create-namespace \
  --timeout=10m
```

Once installation is complete, forward the OpenComply UI to your local machine:

```bash
kubectl port-forward -n opencomply svc/nginx-proxy 8080:80
```

Then browse to [http://localhost:8080](http://localhost:8080) to access the OpenComply dashboard.

We offer quickstart guides for AWS, Azure, DigitalOcean, GKE, and Linode in our docs (link coming soon).

## Next Steps

- **Download Open-Source:** Access the latest stable release from our GitHub releases.
- **Try Cloud for Free:** Sign up for our hosted Cloud offering (coming soon).
- **Explore Resources:** Check out recommended best practices and compliance guidelines:
  - Reliability & Security Best Practices
  - FedRAMP, HIPAA, CIS, CISA Cyber Essentials
- **Customize:** Configure controls, frameworks, roles, SSO, policies, queries, and tasks to suit your environment.
