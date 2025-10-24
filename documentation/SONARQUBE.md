# SonarQube CI Configuration

The repository ships with a reusable GitHub Actions workflow (`.github/workflows/sonarqube.yaml`) that
runs the SonarQube scanner on pushes to the `main` branch and on pull requests.

## Credentials and project information

The workflow looks for the following values in the repository settings, trying them in the order shown
below. The first non-empty value is used.

| Purpose | Secret | Variable | Notes |
| --- | --- | --- | --- |
| Authentication token | `SONAR_TOKEN` | `SONAR_TOKEN` | Must be defined as a secret to keep the token private. |
| SonarQube server URL | `SONAR_HOST_URL` | `SONAR_HOST_URL` | Example: `https://sonarqube.example.com`. |
| SonarQube project key | `SONAR_PROJECT_KEY` | `SONAR_PROJECT_KEY` | Optional. If omitted the scanner uses the key defined in `sonar-project.properties`. |

If neither the token nor host URL is provided, the workflow will skip the scan instead of failing. This
allows forks to enable their own SonarQube projects without interfering with the upstream configuration.

## Fork-specific setup

1. Create a project for your fork in your SonarQube instance.
2. Generate a project analysis token.
3. In the forked repository, navigate to **Settings → Secrets and variables → Actions** and add:
   - Secret `SONAR_TOKEN` with the project token.
   - Variable or secret `SONAR_HOST_URL` pointing to your SonarQube server.
   - (Optional) Variable or secret `SONAR_PROJECT_KEY` with the project key you created.
4. Push to `main` or open a pull request to trigger the workflow.

With the values above in place the workflow in the fork will run against your SonarQube instance. If you
do not supply a project key override the upstream `sonar-project.properties` file continues to control
the key, ensuring forks do not interfere with the original project configuration.
