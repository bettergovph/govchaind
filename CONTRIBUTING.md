# Contributing to OpenGovChain

First off, thank you for considering contributing to OpenGovChain. It's people like you that make this project a great tool for transparency.

This document outlines the conventions we follow for contributions. Please read it carefully to ensure a smooth and effective collaboration.

## How Can I Contribute?

You can contribute in several ways:

-   **Reporting Bugs:** If you find a bug, please search for existing issues before opening a new one. If no related issue exists, open a new one with a clear title, a detailed description of the bug, and steps to reproduce it.
-   **Suggesting Enhancements:** If you have an idea for an enhancement, open an issue with the `enhancement` label. Describe the enhancement and its potential benefits.
-   **Pull Requests:** If you want to contribute code, please follow the process described below.

## Contribution Workflow

Our contribution workflow is designed to be simple and effective.

### 1. Fork the Repository

Start by forking the repository and cloning it to your local machine.

### 2. Branching Strategy

All changes should be made in a feature branch. Our branch naming convention is:

```
type/short-description
```

-   **`type`**: Describes the kind of change. We use the following types:
    -   `feat`: A new feature.
    -   `fix`: A bug fix.
    -   `docs`: Documentation only changes.
    -   `style`: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc).
    -   `refactor`: A code change that neither fixes a bug nor adds a feature.
    -   `perf`: A code change that improves performance.
    -   `test`: Adding missing tests or correcting existing tests.
    -   `chore`: Changes to the build process or auxiliary tools and libraries such as documentation generation.
    -   `ci`: Changes to our CI configuration files and scripts.
    -   `build`: Changes that affect the build system or external dependencies.

-   **`short-description`**: A few words describing the change, separated by hyphens.

**Examples:**

```
feat/add-new-dataset-query
fix/incorrect-gas-estimation
docs/update-api-documentation
```

### 3. Development

-   **Code Style:** This project uses `gofmt` for code formatting. Please ensure your code is formatted before committing.
-   **Linting:** We use `golangci-lint` for linting. Run `make lint` to check your code for issues. You can run `make lint-fix` to automatically fix some issues.
-   **Testing:** All new features and bug fixes must be accompanied by unit tests. Run `make test` to run the test suite.

### 4. Commit Message Convention

We follow the **Conventional Commits** specification. Your commit messages should be structured as follows:

```
type(scope): description

[optional body]

[optional footer]
```

-   **`type`**: Must be one of the types listed in the branching strategy section.
-   **`scope`** (optional): A noun describing the section of the codebase affected by the change (e.g., `x/datasets`, `app`, `docker`).
-   **`description`**: A concise, imperative-tense description of the change.

**Examples:**

```
feat(x/datasets): add support for querying datasets by owner
```

```
fix(app): correct genesis file validation
```

```
docs(contributing): add new contribution guidelines
```

### 5. Pull Request (PR) Process

Once your changes are ready, push your branch to your fork and open a pull request to the `main` branch of the main repository.

-   **PR Title:** Your PR title should follow the Conventional Commits format.
-   **PR Description:** Use the following template for your PR description:

    ```markdown
    ## Description

    A clear and concise description of the changes.

    ## Related Issue

    Link to the issue that this PR addresses.

    ## How to Test

    Detailed steps on how to test your changes.
    ```

-   **CI Checks:** Your PR must pass all CI checks before it can be merged.

## Development Setup

Please refer to the main [README.md](./README.md) for instructions on how to set up your development environment. Before you start, run the verification script to ensure your system is ready:

```bash
./scripts/check-build-env.sh
```

Thank you for your contribution!