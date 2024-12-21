package configs

var TablesToResourceTypes = map[string]string{
	 "Github/Actions/Artifact": "github_actions_artifact",
  "Github/Actions/Repository/Runner": "github_actions_runner",
  "Github/Actions/Repository/Secret": "github_actions_secret",
  "Github/Actions/Repository/Workflow_run": "github_actions_workflow_run",
  "Github/Blob": "github_blob",
  "Github/Branch": "github_branch",
  "Github/Branch/Protection": "github_branch_protection",
  "Github/Commit": "github_commit",
  "Github/Issue": "github_issue",
  "Github/License": "github_license",
  "Github/Organization": "github_organization",
  "Github/Organization/Collaborator": "github_organization_collaborator",
  "Github/Organization/Dependabot/Alert": "github_organization_dependabot_alert",
  "Github/Organization/External/Identity": "github_organization_external_identity",
  "Github/Organization/Member": "github_organization_member",
  "Github/PullRequest": "github_pull_request",
  "Github/Release": "github_release",
  "Github/Repository": "github_repository",
  "Github/Repository/Collaborator": "github_repository_collaborator",
  "Github/Repository/DependabotAlert": "github_repository_dependabot_alert",
  "Github/Repository/Deployment": "github_repository_deployment",
  "Github/Repository/Environment": "github_repository_environment",
  "Github/Repository/Ruleset": "github_repository_ruleset",
  "Github/Repository/SBOM": "github_repository_sbom",
  "Github/Repository/VulnerabilityAlert": "github_repository_vulnerability_alert",
  "Github/Tag": "github_tag",
  "Github/Team": "github_team",
  "Github/Team/Member": "github_team_member",
  "Github/Tree": "github_tree",
  "Github/User": "github_user",
  "Github/Workflow": "github_workflow",
  "Github/Container/Package": "github_container_package",
  "Github/Package/Maven": "github_maven_package",
  "Github/NPM/Package": "github_npm_package",
  "Github/Nuget/Package": "github_nuget_package",
  "Github/Artifact/DockerFile": "github_artifact_dockerfile",
}

var ResourceTypesList = []string{
  "Github/Actions/Artifact",
  "Github/Actions/Repository/Runner",
  "Github/Actions/Repository/Secret",
  "Github/Actions/Repository/Workflow_run",
  "Github/Blob",
  "Github/Branch",
  "Github/Branch/Protection",
  "Github/Commit",
  "Github/Issue",
  "Github/License",
  "Github/Organization",
  "Github/Organization/Collaborator",
  "Github/Organization/Dependabot/Alert",
  "Github/Organization/External/Identity",
  "Github/Organization/Member",
  "Github/PullRequest",
  "Github/Release",
  "Github/Repository",
  "Github/Repository/Collaborator",
  "Github/Repository/DependabotAlert",
  "Github/Repository/Deployment",
  "Github/Repository/Environment",
  "Github/Repository/Ruleset",
  "Github/Repository/SBOM",
  "Github/Repository/VulnerabilityAlert",
  "Github/Tag",
  "Github/Team",
  "Github/Team/Member",
  "Github/Tree",
  "Github/User",
  "Github/Workflow",
  "Github/Container/Package",
  "Github/Package/Maven",
  "Github/NPM/Package",
  "Github/Nuget/Package",
  "Github/Artifact/DockerFile",
}