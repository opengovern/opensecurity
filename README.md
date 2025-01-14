<p align="right">
  <a href="https://opengovernance.io">
    <picture>
          <source media="(prefers-color-scheme: dark)" srcset="https://github.com/opengovern/opencomply/blob/main/assets/logos/logo-dark.png">
          <source media="(prefers-color-scheme: light)" srcset="https://github.com/opengovern/opencomply/blob/main/assets/logos/logo-light.png">
          <img alt="OpenGovernance" src="https://github.com/opengovern/opencomply/blob/main/assets/logos/logo-light.png"
    </picture>

  </a>
</p>

<p align="center"> <em>🚀 Full Stack Governance,🚀 Goven across clouds, platforms, and tools 🚀 Maintain policies as Code,🚀 Steampipe Compatible, 🚀 Unify Security, Compliance, and Ops.</em> </p>

OpenGovernance simplifies governance, compliance, security, and operations across clouds, platforms, and on-premises. Steampipe-compatible and Git-managed, it enforces top policies, optimizes costs, boosts efficiency and reliability, and aligns with the Well-Architected Framework.

![App Screenshot](./assets/screenshots/app-%20screenshot%203.png)

Unlike traditional governance tools that are complex to set up and maintain, OpenGovernance is user-friendly and easy to operate. You can have your governance framework up and running in minutes without dealing with intricate configurations.

Additionally, OpenGovernance replaces legacy compliance systems by providing a unified interface, eliminating the need for multiple separate installations. It supports managing standards like SOC2 and HIPAA, ensuring your organization stays compliant with less effort.

By optimizing your compliance and governance processes, OpenGovernance helps reduce operational costs.

## 🌟 Features:
- **Centralized Multi-Cloud Governance**: Manage AWS, Azure, and GCP policies from one platform.
- **Steampipe Compatibility**: Leverage Steampipe Queries, and utilize vendor neutral polices
- **Batteries included**: Over 2,500 unique policies and 50+ benchmarks, including built-in support for NIST, HIPAA, SOC 2, CIS, and more across multiple clouds.
- **Track History & Capture Evidence**: Keep an audit trail and ensure regulatory compliance, over time
- **Customizable Policy Controls**: Use simple SQL to define and enforce your standards.
- **Vendor-Neutral & Open Source**: Flexible integration with existing tools and platforms.
- **Role-Based Access Control (RBAC)**: Secure, fine-grained access management. SSO/OIDC available.
- **User-Friendly Interface**: Intuitive, WebUI, with API

## ⚡️ Quick start on Kubernetes:

### Add the Helm Repository:

```bash
helm repo add opencomply https://charts.opencomply.io --force-update
```

### Install with Helm
```bash
helm install opencomply opencomply/opencomply --namespace opencomply --create-namespace --timeout 10m 
```

### Expose the app

```bash
kubectl port-forward -n opencomply svc/nginx-proxy 8080:80
```
Navigate to http://localhost:8080/ in your browser.
To sign in, use admin@opencomply.io as the username and password as the password.
