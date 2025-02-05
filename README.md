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

**opencomply** simplifies security and compliance across your entire stack—from containers to cloud—so you can ship faster and worry less.

*   **See everything:** Get a complete view of your infrastructure, data, identities, configurations, and security posture across all your clouds and tools.
*   **Govern anything:** Assess and enforce compliance for configurations, processes, and security across your entire environment.
*   **Adapt easily:** Define any rule—best practices, regulatory, or change requirements—as code, manage them in Git, and integrate with your existing tools and CI/CD pipeline.

**Key Features**

*   **CloudQL:** Based on Steampipe, CloudQL lets you use SQL to explore thousands of different asset types—like that old, unapproved Kubernetes cluster you've been meaning to address.
*   **Policy-as-Query:** Define your configuration checks and other requirements (e.g., "Did X happen before Y?") as simple SQL queries, and then manage those checks in Git.
*   **Customize:** Easily write your own checks—even handle hard-to-tackle policies like Approved Base Docker Images or OpenAI Assistants with governed data sources.
*   **Scalable:** Built with KEDA and OpenSearch, CloudQL scales effortlessly to handle thousands of checks across even very large infrastructures.
*   **Battle-Tested:** We've thoroughly tested CloudQL on over $50M of real-world cloud infrastructure, so you can be confident in its reliability.
*   **Extensible:** Extensive integrations—AWS, Azure, DigitalOcean, Linode, GitHub, Render, Fly.io, Heroku, Cloudflare, and more. Fork one and write your own.

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

- **Download Open-Source:** Access the latest stable release from our GitHub releases.
- **Try Cloud for Free:** Sign up for our hosted Cloud offering (coming soon).
- **Explore Resources:** Check out recommended best practices and compliance guidelines:
  - Reliability & Security Best Practices
  - FedRAMP, HIPAA, CIS, CISA Cyber Essentials
- **Customize:** Configure controls, frameworks, roles, SSO, policies, queries, and tasks to suit your environment.
