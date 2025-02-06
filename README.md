<p align="center">
  <a href="https://opencomply.io">
    <picture>
      <!-- Dark mode logo -->
      <source media="(prefers-color-scheme: dark)" srcset="https://github.com/opengovern/opencomply/blob/main/assets/logos/logo-dark.png">
      <!-- Light mode logo -->
      <source media="(prefers-color-scheme: light)" srcset="https://github.com/opengovern/opencomply/blob/main/assets/logos/logo-light.png">
      <!-- Fallback/logo -->
      <img
        width="70%"
        alt="opencomply"
        src="https://github.com/opengovern/opencomply/blob/main/assets/logos/logo-dark.png"
      >
    </picture>
  </a>
</p>


<p align="center">
  <img 
    src="./assets/screenshots/app-screenshot-1.png"
    alt="App Screenshot"
    width="100%"
  />
</p>


# opencomply

#### Simplify security and compliance across your entire stack—from containers to cloud—so you can ship faster and worry less.


**What *opencomply* Does:**

*   **See Everything:** Complete visibility into your infrastructure, data, identities, configurations, and security across all your clouds and tools.
*   **Govern Anything:** Assess and enforce compliance for configurations, processes, and security.
*   **Adapt Easily:** Define any rule (best practices, regulatory, or internal) as code (SQL policies), manage in Git, and integrate with your tools and CI/CD.

**Key Features:**

*   **Unified Visibility (CloudQL):** Explore all your assets (containers, cloud resources, etc.) with SQL.
*   **Customizable Controls (Policy-as-Code):** Define any compliance check as a SQL policy and manage it in Git.
*   **Flexible Compliance:** Easily create your own checks, even complex ones.
*   **Scalable Audits:** Handles thousands of checks across large infrastructures.
*   **Extensive Integrations:** Connect to AWS, Azure, DigitalOcean, Linode, GitHub, and more.

## Getting Started

**Helm Installation:** 

Install on any Kubernetes clusters with at least 3 nodes (4 vCPUs x 16GB RAM each).

```bash
helm repo add opencomply https://charts.opencomply.io --force-update
helm install -n opencomply opencomply opencomply/opencomply --create-namespace
kubectl port-forward -n opencomply svc/nginx-proxy 8080:80
```
Open http://localhost:8080/ in your browser, sign in with ```admin@opencomply.io``` as the username and ```password``` as the password.

App includes sample data.

## Next Steps

*   **Explore the Documentation:** Visit [docs.opencomply.io](https://docs.opencomply.io) for detailed information and guides.
*   **Try Cloud for Free:** Sign up for our hosted Cloud offering (coming soon).
