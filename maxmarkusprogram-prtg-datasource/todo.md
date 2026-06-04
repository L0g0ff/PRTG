# Plugin publishing best practices

When publishing a Grafana plugin, follow best practices to ensure a smooth submission, review process, and a higher quality experience for users. Whether you’re fine-tuning your plugin’s functionality or preparing your documentation, by following established guidelines you improve the plugin’s performance, security, and discoverability within the Grafana ecosystem.

The recommendations in this document will help you avoid common pitfalls, streamline the review process, and create a plugin that integrates seamlessly into users' workflows.

caution

The Grafana Plugins team studies each plugin submission individually and approves any given plugin on a case-by-case basis.

Following [best practices](https://grafana.com/developers/plugin-tools/publish-a-plugin/publishing-best-practices) or [providing testing information](https://grafana.com/developers/plugin-tools/publish-a-plugin/provide-test-environment) does not guarantee the approval of a submitted plugin.

Refer to [Plugin submission review](https://grafana.com/developers/plugin-tools/publish-a-plugin/publish-a-plugin#plugin-submission-review) for more details.

## Before you begin

Before you proceed make sure you've done the following:

* [Created your initial plugin](https://grafana.com/developers/plugin-tools/)
* Reviewed the [best practices for plugin development guide](https://grafana.com/developers/plugin-tools/key-concepts/best-practices)
* Familiarized yourself with the [plugin signature levels](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin#public-or-private-plugins)

## Populate your plugin's metadata

Use metadata to make your Grafana plugin discoverable and user-friendly. Properly structuring the [metadata in your `plugin.json` file](https://grafana.com/developers/plugin-tools/reference/plugin-json) not only helps users find your plugin in [Grafana’s plugin catalog](https://grafana.com/grafana/plugins/) but also provides essential details about the plugin’s functionality and compatibility.

Focus on these key elements:

**[Plugin name](https://grafana.com/developers/plugin-tools/reference/plugin-json)**

`name`

The name of your plugin must be clear, concise, and descriptive. It's the first point of interaction for potential users, so avoid overly generic or cryptic names. Aim for a name that reflects the plugin’s primary functionality and makes it easy to understand its purpose at a glance.

**[Description](https://grafana.com/developers/plugin-tools/reference/plugin-json#info)**

`info.description`

The description field summarizes what your plugin does and why users should install it. Limit the description to two sentences, highlighting the core functionality and use cases. A well-written description not only informs users, but also contributes to better search results in the catalog.

**[Keywords](https://grafana.com/developers/plugin-tools/reference/plugin-json#info)**

`info.keywords`

Keywords improve the searchability of your plugin within Grafana’s catalog. Choose terms that accurately describe your plugin’s functionality and data types it supports, such as "JSON", "SQL", or "visualization".

caution

Avoid keyword stuffing. Irrelevant keywords are flagged during the review process, potentially delaying publication.

**[Logos](https://grafana.com/developers/plugin-tools/reference/plugin-json#info)**

`info.logos`

Adding logos improves the overall look and feel of your plugin in the plugin catalog, and adds legitimacy and professionalism to your plugin.

**[Screenshots](https://grafana.com/developers/plugin-tools/reference/plugin-json#info)**

`info.screenshots`

Use the screenshots field to provide an array with one or more screenshot images to be displayed in the plugin catalog. Images give users a visual representation of your plugin, and can help them decide if you plugin solves their problem before installing it.

Make sure that your screenshots:

* Show your plugin in action, highlighting its standout features
* Have a suitable resolution and file type, for example png, jpeg, or gif

**[Sponsorship link](https://grafana.com/developers/plugin-tools/reference/plugin-json#infolinks)**

`info.links`

The sponsorship link provides a way for users to support your work and contribute to its development. It appears in the "Links" section of your plugin's detail page and supports various funding platforms, such as GitHub Sponsors or Patreon.

Example:

```text
{
  info: {
    links: [
      {
        name: "sponsorship",
        url: "https://github.com/sponsors/pluginDeveloper"
      }
    ]
  }
}
```

**[Grafana version compatibility](https://grafana.com/developers/plugin-tools/reference/plugin-json#dependencies)**

`dependencies.grafanaDependency`

Specify the minimum Grafana version your plugin is compatible with so that users running different versions of Grafana know whether your plugin will work for them. Be sure to [run end-to-end tests](https://grafana.com/developers/plugin-tools/e2e-test-a-plugin/) to confirm compatibility with releases you support.

## Create a comprehensive README

Your plugin's README file serves as both a first impression and a detailed guide for your users. It's a combination of a storefront advertisement and an instruction manual—showing what your plugin can do, how to install it, and how users can make the most of it within their Grafana instances.

Use the [README template](https://raw.githubusercontent.com/grafana/plugin-tools/main/packages/create-plugin/templates/common/src/README.md), included as part of the plugin structure generated by the `create-plugin` tool, to provide your users with everything they need to confidently use and contribute to your plugin.

The template covers the essential components, and you can add more specific details to help users understand the value and functionality of your plugin, such as:

* **Screenshots or screen captures:** Include screenshots or even video demonstrations so users can quickly grasp the plugin’s capabilities and setup process, giving them confidence to use it effectively.
* **Dynamic badges:** Badges provide quick information about your plugin, such as the latest release version or whether it has passed security and code checks. Use tools like [shields.io](https://shields.io/) with the Grafana.com API to automatically update these badges whenever you publish a new version, adding transparency and trustworthiness to your plugin.
* **Contribution guidance:** Maintaining a plugin can be demanding, especially for individual developers. Clearly outline how users can provide feedback and report bugs, and direct potential code contributors to your `contributing.md`. This fosters community involvement and makes it easier to maintain and improve your plugin over time.

## Maintain a detailed changelog

A well-maintained changelog is essential for plugin transparency and helps users understand what's changed between versions. Grafana displays your changelog in the plugin details page so users can evaluate whether to update.

info

Use Grafana's automated changelog generation feature to simplify the process of maintaining your changelog. Learn how in the [Automatically Generate Changelogs](https://grafana.com/developers/plugin-tools/publish-a-plugin/build-automation#generate-changelogs-automatically) guide.

### Changelog best practices

Use a dedicated CHANGELOG.md file in your repository root with the following information:

1. Follow semantic versioning (MAJOR.MINOR.PATCH) and organize entries by version
2. Date each release to provide chronological context
3. Group changes by type such as "Features", "Bug Fixes", "Breaking Changes"...
4. Reference pull requests with links to provide additional context
5. Highlight breaking changes prominently to alert users of required actions

### Example

```markdown
### [1.10.0](https://github.com/user/plugin-name/tree/1.10.0) (2025-04-05)

**Implemented enhancements:**

- Add support for dark theme [\#138](https://github.com/user/plugin-name/pull/138) ([username](https://github.com/username))
- Add ability to customize tooltip formats [\#135](https://github.com/user/plugin-name/pull/135) ([username](https://github.com/username))
- Support for PostgreSQL data source [\#129](https://github.com/user/plugin-name/pull/129) ([username](https://github.com/username))

**Fixed bugs:**

- Fix panel crash when switching dashboards [\#139](https://github.com/user/plugin-name/pull/139) ([username](https://github.com/username))
- Fix inconsistent time zone handling [\#134](https://github.com/user/plugin-name/pull/134) ([username](https://github.com/username))

**Closed issues:**

- Documentation needs examples for PostgreSQL queries [\#130](https://github.com/user/plugin-name/issues/130)

**Merged pull requests:**

- Update dependencies to address security vulnerabilities [\#140](https://github.com/user/plugin-name/pull/140) ([username](https://github.com/username))

**Breaking changes:**

- Migrate configuration storage format [\#115](https://github.com/user/plugin-name/pull/115) ([username](https://github.com/username))
```

With this format your changelog becomes a transparent resource that clearly communicates changes, acknowledges contributions, and provides links to more detailed information. It helps users make informed decisions about updating your plugin and demonstrates your commitment to maintaining a high-quality Grafana plugin.

## End-to-end testing

End-to-end (E2E) testing ensures that your Grafana plugin works correctly across various environments and supported Grafana versions. It replicates real-world usage by testing the plugin in an environment similar to the end-user's setup. Implementing E2E tests helps catch issues before submission, saving time during the review process and ensuring a smoother user experience.

**Key points:**

* **Test compatibility across versions:** Ensure your plugin works seamlessly with various versions of Grafana by setting up E2E tests targeting multiple releases.
* **Automate testing:** Integrate E2E testing into your continuous integration (CI) pipeline to catch issues early and frequently, reducing potential problems during review.

For a comprehensive guide on setting up E2E tests, refer to our [E2E test a plugin](https://grafana.com/developers/plugin-tools/e2e-test-a-plugin/) documentation.

## Validate your plugin

Before submitting your plugin for review, use the Plugin Validator to check for potential issues that could prevent your plugin from being accepted, such as security vulnerabilities or structural problems.

**Key points:**

* **Run locally or in CI:** You can run the validator locally or integrate it into your CI workflow to automate the validation process. Note, the validator runs automatically during the default release workflow.
* **Validation reports:** The tool generates a report, highlighting any errors or warnings that need to be addressed before submission.

For more information on using the validator refer to the [Plugin Validator documentation](https://github.com/grafana/plugin-validator).

## Provide a provisioned test environment

Provisioning a test environment for your plugin can significantly reduce the review time and make it easier for others to test and contribute to your plugin. A provisioned environment includes a pre-configured Grafana instance with sample dashboards and data sources that demonstrate your plugin's functionality.

**Key points:**

* **Why provisioning matters:** It ensures that both reviewers and contributors can quickly verify your plugin's behaviour without manual setup, speeding up the review and collaboration process.
* **Automated setup:** You can provision test environments using Docker to create an out-of-the-box experience that replicates a typical Grafana setup.

To learn more about setting up provisioned environments, check out our [provisioning guide](https://grafana.com/developers/plugin-tools/publish-a-plugin/provide-test-environment).

## Automate releases with GitHub Actions

To streamline your plugin development workflow, automate releases using GitHub Actions. This ensures that your plugin is built, signed, and packaged correctly on each release, reducing human error and expediting the publishing process.

**Key points:**

* **Continuous integration (CI):** Use GitHub Actions to automatically build and test your plugin on every commit or pull request, catching issues early.
* **Release workflow:** Automate the signing and packaging of your plugin when you're ready to publish, ensuring it meets the necessary criteria for submission to the Grafana plugin catalog.

For more details refer to Grafana's [Automate packaging and signing with GitHub](https://grafana.com/developers/plugin-tools/publish-a-plugin/build-automation) guide.

## Next steps

Follow these best practices to increase the chances of a successful plugin submission. They are designed to ensure that your plugin not only passes our review process but also delivers an exceptional experience for users. Adopting these practices will streamline your workflow and help create plugins that stand out in the Grafana ecosystem.

Once your plugin is ready to be published follow our guide for [submitting your plugin for review](https://grafana.com/developers/plugin-tools/publish-a-plugin/publish-a-plugin). We look forward to seeing what you create!

# Automate the packaging and signing of your plugin with GitHub CI

Set up your plugin to use the supplied [GitHub workflows](https://grafana.com/developers/plugin-tools/set-up/set-up-github) from [create-plugin](https://grafana.com/developers/plugin-tools/) to ensure that your plugin builds and packages in the correct format. Additionally, you can use the zip file that this workflow produces to test the plugin.

If you include a Grafana Access Policy Token in your [GitHub repository secrets](https://docs.github.com/en/codespaces/managing-codespaces-for-your-organization/managing-development-environment-secrets-for-your-repository-or-organization), the system automatically creates a signed build that you can use to test the plugin locally before submission. For information about how to create this token, refer to the [sign a plugin](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin#generate-an-access-policy-token) documentation.

When you create a release tag, the process becomes automated and results in a zip file that you can submit for publication to the [Grafana plugin catalog](https://grafana.com/plugins).

You can use the links to the archive and zip files from the release page to make your plugin submission.

## Package your plugin with GitHub CI

Follow these steps to package your plugin with GitHub CI.

To package your plugin in a ZIP file manually, refer to [Package a plugin](https://grafana.com/developers/plugin-tools/publish-a-plugin/package-a-plugin).

### Set up the release workflow

Ensure your repository contains a `.github/workflows/release.yml` file with the following contents:

```yaml
# filepath: .github/workflows/release.yml
name: Release

on:
push:
tags:
-'v*'# Run workflow on version tags, e.g. v1.0.0.

jobs:
release:
permissions:
id-token: write
contents: write
attestations: write
runs-on: ubuntu-latest
steps:
-uses: actions/checkout@v4

-uses: grafana/plugin-actions/build-plugin@main
with:
# refer to https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin#generate-an-access-policy-token to generate it
# save the value in your repository secrets
policy_token: ${{ secrets.GRAFANA_ACCESS_POLICY_TOKEN }}
# creates a signed build provenance attestation to verify the authenticity of the plugin build
attestation:true
```

### Trigger the release workflow

To trigger the release workflow, push a tag with the format `vX.X.X` to the repository. Typically, you merge all of your changes into `main`, and you apply the tag to `main`.

#### Create a `vX.X.X` tag with your package manager (recommended)

Use your package manager to create a version tag.

The following examples create a patch version following [Semantic Versioning](https://semver.org/):

With [npm](https://docs.npmjs.com/cli/v7/commands/npm-init):

```sh
npm version patch
```

With [yarn](https://yarnpkg.com/lang/en/docs/cli/version/):

```sh
yarn version patch
```

With [pnpm](https://pnpm.io/):

```sh
pnpm version patch
```

This updates your version in the `package.json` file and creates a new Git tag with the format `vX.X.X`. You can change `patch` to `minor` or `major` to create a new minor or major version.

After you create the tag, push it to the repository:

```sh
git push origin main --follow-tags
```

### Publish your release in GitHub

After you [create and push the tag](https://grafana.com/developers/plugin-tools/publish-a-plugin/build-automation#trigger-the-release-workflow), the release workflow runs and generates a release with all the artifacts you need to submit your plugin to the [Grafana plugin catalog](https://grafana.com/plugins).

The workflow creates a  **draft release** . You can edit the release in GitHub, update the description as needed, and then publish it. For more details about managing repository releases, refer to the [GitHub documentation](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository).

### Use your release assets for your plugin submission

After you publish the draft release, you can use the release assets to submit your plugin to the [Grafana plugin catalog](https://grafana.com/plugins). Copy the links to the archive (zip) file and sha1 sum and use them in the plugin submission form.

### Download the release zip file

Access the release zip file directly from the GitHub repository release path (for example, `https://github.com/org/plugin-id/releases`).

## Sign your plugin automatically

You can sign your plugin releases using GitHub Action.

First, [generate an Access Policy Token](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin#generate-an-access-policy-token) and [save it in your repository secrets](https://docs.github.com/en/actions/security-for-github-actions/security-guides/using-secrets-in-github-actions#creating-secrets-for-a-repository) as `GRAFANA_ACCESS_POLICY_TOKEN`.

By default, create-plugin adds the following `release.yml` to your scaffolded plugin. If this is missing from your plugin repository, copy the following to add the workflow:

```yaml
# filepath: .github/workflows/release.yml
name: Release

on:
push:
tags:
-'v*'# Run workflow on version tags, e.g. v1.0.0.

jobs:
release:
permissions:
id-token: write
contents: write
attestations: write
runs-on: ubuntu-latest
steps:
-uses: actions/checkout@v4

-uses: grafana/plugin-actions/build-plugin@main
with:
# refer to https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin#generate-an-access-policy-token to generate it
# save the value in your repository secrets
policy_token: ${{ secrets.GRAFANA_ACCESS_POLICY_TOKEN }}
attestation:true
use_changelog_generator:true# Enable automatic changelog generation
```

Next, follow the regular process to [trigger](https://grafana.com/developers/plugin-tools/publish-a-plugin/build-automation#trigger-the-release-workflow) the release workflow. Your plugin is signed automatically, and you can use the release assets for your plugin submission.

## Attest provenance for plugin builds

Provenance attestation generates verifiable records of the build's origin and process and enhances the security of your plugin builds. With this feature, users can confirm that the plugin they're installing was created through your official build pipeline.

Currently, this feature is available only with GitHub Actions in public repositories. While we recommend using GitHub Actions with provenance attestation for improved security, you can still build and distribute plugins using other CI/CD platforms or manual methods.

### Enable provenance attestation

To enable provenance attestation in your existing GitHub Actions workflow:

1. Add required permissions to your workflow job:

```yaml
permissions:
id-token: write
contents: write
attestations: write
```

2. Enable attestation in the `build-plugin` action:

```yaml
-uses: grafana/plugin-actions/build-plugin@main
with:
policy_token: ${{ secrets.GRAFANA_ACCESS_POLICY_TOKEN }}
attestation:true
```

The workflow generates attestations automatically when it builds your plugin zip file.

### Troubleshoot provenance attestation

If you encounter errors in the plugin validator or your plugin submission like these:

* "No provenance attestation. This plugin was built without build verification."
* "Cannot verify plugin build."

Follow the steps in the [Enable provenance attestation](https://grafana.com/developers/plugin-tools/publish-a-plugin/build-automation#enable-provenance-attestation) section to enable provenance attestation in your GitHub Actions workflow.

## Generate changelogs automatically

Maintaining a detailed changelog is essential for communicating updates to your users and displays prominently in the Grafana plugin details page. To simplify this process, our plugin build workflow supports automatic changelog generation.

### Use the GitHub Actions workflow to generate changelog

The build-plugin GitHub Action can automatically generate and maintain your plugin's changelog using the [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator) tool. This feature:

* Creates a comprehensive `CHANGELOG.md` file organized by release.
* Groups changes by type (features, bug fixes, and more).
* Includes links to pull requests and issues.
* Acknowledges contributors.
* Commits the updated changelog to your repository.

To enable automatic changelog generation in your workflow, add the `use_changelog_generator: true` parameter to your build-plugin action:

```yaml
-uses: grafana/plugin-actions/build-plugin@main
with:
policy_token: ${{ secrets.GRAFANA_ACCESS_POLICY_TOKEN }}
attestation:true
use_changelog_generator:true# Enable automatic changelog generation
```

### Requirements

To use this feature, ensure your workflow has the necessary permissions:

```yaml
permissions:
contents: write
```

The changelog generator requires write access to commit the updated `CHANGELOG.md` file to your repository.

If your target branch is protected, the default `github.token` can't push changes directly, even with write permissions. In this case, you need to:

1. Create a Personal Access Token (PAT) with appropriate permissions.
2. Store it as a repository secret (for example, `CHANGELOG_PAT`).
3. Configure the action to use this token:

```yaml
-name: Build plugin
uses: grafana/plugin-actions/build-plugin@main
with:
use_changelog_generator:true
token: ${{ secrets.CHANGELOG_PAT }}# Replace default github.token
```

### Generated changelog format

The generated changelog follows a standardized format that clearly categorizes changes:

```markdown
## [1.2.0](https://github.com/user/plugin-name/tree/1.2.0) (2025-03-15)

**Implemented enhancements:**

- Add dark theme support [\#138](https://github.com/user/plugin-name/pull/138) ([username](https://github.com/username))
- Add tooltip customization options [\#135](https://github.com/user/plugin-name/pull/135) ([username](https://github.com/username))

**Fixed bugs:**

- Fix panel crash when switching dashboards [\#139](https://github.com/user/plugin-name/pull/139) ([username](https://github.com/username))
- Fix time zone handling inconsistencies [\#134](https://github.com/user/plugin-name/pull/134) ([username](https://github.com/username))

**Closed issues:**

- Documentation needs more examples [\#130](https://github.com/user/plugin-name/issues/130)

**Merged pull requests:**

- Update dependencies for security [\#140](https://github.com/user/plugin-name/pull/140) ([username](https://github.com/username))
```

# Package a plugin

Package a plugin to organize the plugin code and make it ready for use in your organization. Follow these steps to package the plugin in a ZIP file.

1. Build the plugin

   * npm
   * Yarn
   * pnpm

   ```shell
   npminstall
   npm run build
   ```
2. Optional: If your plugin has a backend, build it as well.

   ```text
   mage
   ```

   Make sure that all the binaries are executable and have a `0755` (`-rwxr-xr-x`) permission.
3. Sign the plugin. To learn more, refer to [Sign a plugin](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin).
4. Rename the `dist` directory to match your plugin ID, and then create a ZIP archive.

   ```text
   mv dist/ myorg-simple-panel
   zip myorg-simple-panel-1.0.0.zip myorg-simple-panel -r
   ```
5. Optional: verify that your plugin is packaged correctly using [zipinfo](https://linux.die.net/man/1/zipinfo). It should look like this:

   ```shell
   $ zipinfo grafana-clickhouse-datasource-1.1.2.zip

   Archive:  grafana-clickhouse-datasource-1.1.2.zip
   Zip file size: 34324077 bytes, number of entries: 22
   drwxr-xr-x          0 bx stor 22-Mar-24 23:23 grafana-clickhouse-datasource/
   -rw-r--r--       1654 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/CHANGELOG.md
   -rw-r--r--      11357 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/LICENSE
   -rw-r--r--       2468 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/MANIFEST.txt
   -rw-r--r--       8678 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/README.md
   drwxr-xr-x          0 bx stor 22-Mar-24 23:23 grafana-clickhouse-datasource/dashboards/
   -rw-r--r--      42973 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/dashboards/cluster-analysis.json
   -rw-r--r--      56759 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/dashboards/data-analysis.json
   -rw-r--r--      39406 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/dashboards/query-analysis.json
   -rwxr-xr-x   16469136 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/gpx_clickhouse_darwin_amd64
   -rwxr-xr-x   16397666 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/gpx_clickhouse_darwin_arm64
   -rwxr-xr-x   14942208 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/gpx_clickhouse_linux_amd64
   -rwxr-xr-x   14155776 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/gpx_clickhouse_linux_arm
   -rwxr-xr-x   14548992 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/gpx_clickhouse_linux_arm64
   -rwxr-xr-x   15209472 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/gpx_clickhouse_windows_amd64.exe
   drwxr-xr-x          0 bx stor 22-Mar-24 23:23 grafana-clickhouse-datasource/img/
   -rw-r--r--        304 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/img/logo.png
   -rw-r--r--       1587 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/img/logo.svg
   -rw-r--r--     138400 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/module.js
   -rw-r--r--        808 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/module.js.LICENSE.txt
   -rw-r--r--     487395 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/module.js.map
   -rw-r--r--       1616 bX defN 22-Mar-24 23:23 grafana-clickhouse-datasource/plugin.json
   22 files, 92516655 bytes uncompressed, 34319591 bytes compressed:  62.9%
   ```

## Next steps

After you've packaged your plugin, you can proceed to:

* [Publish your plugin](https://grafana.com/developers/plugin-tools/publish-a-plugin/publish-a-plugin) to share it with the world, or
* [Install a packaged plugin](https://grafana.com/docs/grafana/latest/administration/plugin-management/#install-a-packaged-plugin) by extracting it into your plugin directory.

# Sign a plugin

Grafana Labs signs all Grafana Labs-authored plugins, including Enterprise plugins, so that Grafana can verify their authenticity with signature verification. By default, Grafana requires all plugins to be signed before it loads them.

Refer to [Plugin signatures](https://grafana.com/docs/grafana/latest/administration/plugin-management/plugin-sign/) for more details.

## Before you begin

### Signatures during plugin development

You don't need to sign a plugin during development or when you submit a plugin for review for the first time. The [Docker development environment](https://grafana.com/developers/plugin-tools/set-up/) that `@grafana/create-plugin` scaffolds is configured by default to run in [development mode](https://github.com/grafana/grafana/blob/main/contribute/developer-guide.md#configure-grafana-for-development), which allows you to load the plugin without a signature.

### Generate an Access Policy token

To verify ownership of your plugin, generate an Access Policy token that you use every time you sign a new version of your plugin.

1. [Create a Grafana Cloud account](https://grafana.com/signup).
2. Log in to your account, and then go to  **My Account > Security > Access Policies** .
3. Click  **Create access policy** .
   **Realm:** Set to *`<YOUR_ORG_NAME>`* (all-stacks)
   **Scope:** Set to **plugins:write**
   ![Create access policy.](https://grafana.com/developers/plugin-tools/assets/images/create-access-policy-v2-8b4191d5722032376519b26189fdf158.png)
   Create access policy.
4. Click **Create token** to create a new token.
   **Expiration date** is optional, though you should change tokens periodically for increased security.
   ![Create access policy token.](https://grafana.com/developers/plugin-tools/assets/images/create-access-policy-token-ef1f2131e5a1994707c011bec80e1e2b.png)
   Create access policy token.
5. Click **Create** and then save a copy of the token somewhere secure for future reference.
6. Proceed to signing your [public plugin](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin#sign-a-public-plugin) or [private plugin](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin#sign-a-private-plugin).

## Public or private plugins

Plugins can have different [signature levels](https://grafana.com/legal/plugins/#what-are-the-different-classifications-of-plugins) depending on their author, related technology, and intended use.

A plugin can be either *public* or  *private* :

* **Public plugins:** Grafana signs these as Community or Commercial. Grafana distributes them within the [Grafana plugin catalog](https://grafana.com/plugins) and makes them available for others to install.
* **Private plugins:** These are only available for use within your organization.

Before you sign your plugin, review the [Plugins policy](https://grafana.com/legal/plugins/) to determine the appropriate signature for your plugin.

## Sign a public plugin

The Grafana team needs to review public plugins before you can sign them.

1. Submit your plugin for [review](https://grafana.com/developers/plugin-tools/publish-a-plugin/publish-a-plugin).
2. If the team approves your plugin, you're granted a plugin signature level. You need this signature level to proceed.
3. In your plugin directory, export the Access Policy token as an environment variable using the token you just created:

   ```sh
   exportGRAFANA_ACCESS_POLICY_TOKEN=<YOUR_ACCESS_POLICY_TOKEN>
   ```
4. Sign the plugin. The Grafana sign-plugin tool creates a [MANIFEST.txt](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin#add-a-plugin-manifest-for-verification) file in the `dist` directory of your plugin:

   * npm
   * Yarn
   * pnpm

   ```shell
   npm run sign
   ```

## Sign a private plugin

1. In your plugin directory, export the Access Policy token as an environment variable using the token you created:

   ```sh
   exportGRAFANA_ACCESS_POLICY_TOKEN=<YOUR_ACCESS_POLICY_TOKEN>
   ```
2. Sign the plugin. The Grafana sign-plugin tool creates a [MANIFEST.txt](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin#add-a-plugin-manifest-for-verification) file in the `dist` directory of your plugin. After the `rootUrls` flag, enter a comma-separated list of URLs for the Grafana instances where you intend to install the plugin:

   * npm
   * Yarn
   * pnpm

   ```shell
   npm run sign -- --rootUrls https://example.com/grafana
   ```

## Add a plugin manifest for verification

For Grafana to verify the digital signature of a plugin, the plugin must include a signed manifest file, `MANIFEST.txt`. The signed manifest file contains two sections:

* **Signed message:** Contains plugin metadata and plugin files with their respective checksums (SHA256).
* **Digital signature:** Created by encrypting the signed message using a private key. Grafana has a built-in public key that it uses to verify that the digital signature was encrypted using the expected private key.

**Example:**

```txt
-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

{
  "manifestVersion": "2.0.0",
  "signatureType": "community",
  "signedByOrg": "myorgid",
  "signedByOrgName": "My Org",
  "plugin": "myorgid-simple-panel",
  "version": "1.0.0",
  "time": 1602753404133,
  "keyId": "7e4d0c6a708866e7",
  "files": {
    "LICENSE": "12ab7a0961275f5ce7a428e662279cf49bab887d12b2ff7bfde738346178c28c",
    "module.js.LICENSE.txt": "0d8f66cd4afb566cb5b7e1540c68f43b939d3eba12ace290f18abc4f4cb53ed0",
    "module.js.map": "8a4ede5b5847dec1c6c30008d07bef8a049408d2b1e862841e30357f82e0fa19",
    "plugin.json": "13be5f2fd55bee787c5413b5ba6a1fae2dfe8d2df6c867dadc4657b98f821f90",
    "README.md": "2d90145b28f22348d4f50a81695e888c68ebd4f8baec731fdf2d79c8b187a27f",
    "module.js": "b4b6945bbf3332b08e5e1cb214a5b85c82557b292577eb58c8eb1703bc8e4577"
  }
}
-----BEGIN PGP SIGNATURE-----
Version: OpenPGP.js v4.10.1
Comment: https://openpgpjs.org

wqEEARMKAAYFAl+IE3wACgkQfk0ManCIZudpdwIHTCqjVzfm7DechTa7BTbd
+dNIQtwh8Tv2Q9HksgN6c6M9nbQTP0xNHwxSxHOI8EL3euz/OagzWoiIWulG
7AQo7FYCCQGucaLPPK3tsWaeFqVKy+JtQhrJJui23DAZLSYQYZlKQ+nFqc9x
T6scfmuhWC/TOcm83EVoCzIV3R5dOTKHqkjIUg==
=GdNq
-----END PGP SIGNATURE-----
```

## Troubleshooting

### Why do I get a "Modified signature" error?

In some cases, the system generates an invalid `MANIFEST.txt` because of an issue when signing the plugin on Windows. You can fix this by replacing all double backslashes, `\\`, with a forward slash, `/`, in the `MANIFEST.txt` file. You need to do this every time you sign your plugin.

### Why do I get a "Field is required: `rootUrls`" error for my public plugin?

With a *public* plugin, your plugin doesn't have a plugin signature level assigned to it yet. A Grafana team member assigns a signature level to your plugin after they review and approve it. For more information, refer to the [Sign a public plugin](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin#sign-a-public-plugin) section.

### Why do I get a "Field is required: `rootUrls`" error for my private plugin?

With a *private* plugin, you need to add a `rootUrls` flag to the `plugin:sign` command. The `rootUrls` must match the [`root_url`](https://grafana.com/docs/grafana/latest/setup-grafana/configure-grafana#root_url) configuration. For more information, refer to the [Sign a private plugin](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin#sign-a-private-plugin) section.

If you still get this error, make sure that the Access Policy token was generated by a Grafana Cloud account that matches the first part of the plugin ID.

# Publish or update a plugin

You've just built your plugin and now you want to share it with the world!

Publishing your plugin to the [Grafana plugin catalog](https://grafana.com/plugins) makes it easily discoverable by millions of Grafana users. Read on to learn how to manage the lifecycle of a plugin in the catalog, from publishing and updating to potentially deprecating.

## Before you begin

* [Review our guidelines](https://grafana.com/legal/plugins/#plugin-publishing-and-signing-criteria) - Learn about the Grafana Labs criteria for publishing and signing plugins.
* [Review our publishing best practices](https://grafana.com/developers/plugin-tools/publish-a-plugin/publishing-best-practices) - Ensure your plugin is in the best state it can be before submitting it for review.
* [Package a plugin](https://grafana.com/developers/plugin-tools/publish-a-plugin/package-a-plugin) - Build the plugin and get it ready to share in the form of a ZIP archive.
* Refer to [plugin-examples](https://github.com/grafana/grafana-plugin-examples) to review best practices for building your plugin.

**To speed up the time it takes to review your plugin:**

* Check that your plugin is ready for review using the [plugin validator](https://github.com/grafana/plugin-validator).
* Provide sample dashboards and test data with your repository so that the plugin's functionality can be verified. Use the [provisioning](https://grafana.com/developers/plugin-tools/publish-a-plugin/provide-test-environment) process provided to simplify this step.

## Publish your plugin

Follow these steps to publish your plugin for the first time.

1. [Sign in](https://grafana.com/auth/sign-in) to your Grafana Cloud account. Note that you need to be an administrator for the Grafana Cloud organization being used to publish the plugin.
2. In the left menu, under Org Settings, click  **My Plugins** .
3. Click  **Submit New Plugin** . The Create Plugin Submission dialog box appears.
   ![Submit plugin.](https://grafana.com/developers/plugin-tools/assets/images/plugins-submission-create-080f2482c8361a8c92ef62b79b0319ed.png)
   Submit plugin.
4. Enter the information requested by the form.
   * **OS & Architecture:**
     * Select **Single** if your plugin archive contains binaries for multiple architectures.
     * Select **Multiple** if you'd like to submit separate plugin archives for each architecture. This can lead to faster downloads since users can select the specific architecture on which they want to install the plugin.
   * **URL:** A URL that points to a ZIP archive of your packaged plugin.
   * **Source code URL:** A URL that points to a public Git repository or ZIP archive of your complete plugin source code.
   * **SHA1:** The SHA1 hash of the plugin specified by the  **URL** .
   * **Testing guidance:** An overview covering the installation, configuration, and usage of your plugin.
   * **Provisioning provided for test environment:** Check this box if you have [configured provisioning](https://grafana.com/developers/plugin-tools/publish-a-plugin/provide-test-environment). If you've done this, rest assured it will be identified during the review, and no additional action is needed on your part.
   * The remaining questions help us determine the [signature level](https://grafana.com/legal/plugins/#what-are-the-different-classifications-of-plugins) for your plugin.
5. Click  **Submit** .

### Plugin submission review

After you submit your plugin:

1. The Grafana Plugins team runs an automated validation to make sure it adheres to the Grafana guidelines.
2. Upon the validation, your submission is placed in a review queue.
3. A plugin reviewer performs a manual inspection that consists of:

* **Code review** : For quality and security purposes, we review the source code for the plugin.
* **Tests** : We install your plugin on one of our Grafana instances to test it for basic use.
* We may ask you to assist us in configuring a test environment for your plugin.
* We'll use this test environment whenever you submit a plugin update.

note

Following [best practices](https://grafana.com/developers/plugin-tools/publish-a-plugin/publishing-best-practices) or [providing a test environment](https://grafana.com/developers/plugin-tools/publish-a-plugin/provide-test-environment) does not guarantee the approval of a submitted plugin. The Grafana Plugins team studies each submission individually and decides on a case-by-case basis.

## Update your plugin

To update a plugin, follow the same guidance as for [publish your plugin](https://grafana.com/developers/plugin-tools/publish-a-plugin/publish-a-plugin#publish-your-plugin), except in Step 3 where you can now click **Submit Update** for the plugin you want to update.

All plugin submissions, new or updates, go through the same automated and rigorous manual review process. Because we may have a test environment already setup for an existing plugin, plugin update reviews may go faster.

## Deprecate a plugin

If a plugin is no longer relevant or is unable to be maintained, plugin developers can request that the plugin be deprecated and removed from the catalog. Similarly, Grafana Labs may deprecate and delist a plugin as part of curating the catalog and ensuring plugins meet our standards for security, quality and compatibility.

For more information on plugin deprecation and how to request your plugin to be deprecated, refer to the Grafana Labs [Plugin Deprecation Policy](https://grafana.com/legal/plugin-deprecation/).

# Help us test your plugin

note

Providing the Grafana Plugins team with additional test configurations and environments isn't required as part of the plugin submission process, but will speed up the review process. For more details, refer to [Plugin submission review](https://grafana.com/developers/plugin-tools/publish-a-plugin/publish-a-plugin#plugin-submission-review).

Developers often ask us how long it takes to review a plugin for publishing to the Grafana [plugin catalog](https://grafana.com/plugins). Although we [can&#39;t give you an estimate](https://grafana.com/developers/plugin-tools/publish-a-plugin/publish-faqs#how-long-does-it-take-to-review-a-submission), you can help us reduce cycle times.

The most time-consuming task when reviewing a plugin is creating suitable test configurations and environments so we can verify your plugin's behavior. This step often involves several back-and-forth conversations between you, the plugin developer, and us, the review team.

To help us test your plugin and improve the review time, you can add testing resources to your plugin via [provisioning](https://grafana.com/docs/grafana/latest/administration/provisioning/#provision-grafana).

## Why provision testing information?

Providing testing context offers several benefits:

* **Faster review times.** If you include provisioning information for your plugin before submission, your wait for the review is much shorter.
* **Easier reviews.** An out-of-the-box working example allows us to easily experiment with additions to the plugin.
* **Easier setup for e2e tests.** Provisioned dashboards allow e2e tests to run against specific scenarios across local development and in CI.
* **Improved clarity.** Provisioned plugins make it easier for tech-savvy users to understand how the plugin works.

## What do we need?

Depending on the type and complexity of your plugin, we require the following resources to test your plugin:

* A simple test JSON dashboard for simple panel or data source plugins.
* A docker-compose.yml file with all the configuration necessary to have the plugin running.
* For more complex plugins, we could require access to a test API to thoroughly try out your plugin.

## How to provide test configurations and environments

Provisioning allows you to add resources in a YAML file under a `/provisioning` directory. We can then use those files to test your plugin as you intended it to work, and provide a better and faster review.

Starting in v2.8.0, `create-plugin` generates provisioning capabilities for all plugin types (apps, datasources and panels) and includes a sample dashboard. If you scaffolded your plugin with a previous version of `create-plugin`, you can run a new command to add the missing provisioning files.

### What you need to do

Provision your plugin with the required testing files and setup that can run from scratch without additional manual commands or configuration.

For example:

* Create an example dashboard with your working plugin, export it as a JSON, and put it in your provisioning files. This is required for panel and data source plugins.
  * If necessary, create a provisioning file for the data source associated with the dashboard.
* Prepare a docker-compose.yml file with all the necessary services set up to use the plugin.
  * For example, a data source plugin often requires a service that provides the data and some seed data.
  * Running Docker compose up from a fresh environment should be enough to make the dashboard example work.
* Use and update the sample dashboard to continuously verify behavior as part of your development process. If appropriate, configure your plugin so that it can return data.

caution

The Grafana Plugins team studies each plugin submission individually and approves any given plugin on a case-by-case basis.

Following [best practices](https://grafana.com/developers/plugin-tools/publish-a-plugin/publishing-best-practices) or [providing testing information](https://grafana.com/developers/plugin-tools/publish-a-plugin/provide-test-environment) does not guarantee the approval of a submitted plugin.

Refer to [Plugin submission review](https://grafana.com/developers/plugin-tools/publish-a-plugin/publish-a-plugin#plugin-submission-review) for more details.

# Publish a plugin: Frequently asked questions

## Do I need to submit a private plugin?

* No. Please only submit plugins that you wish to make publicly available for the Grafana community.

## How long does it take to review a submission?

* We're not able to give an estimate because each plugin submission is unique, though we're constantly working to improve the time it takes to review a plugin. Providing a [provisioned](https://grafana.com/developers/plugin-tools/publish-a-plugin/provide-test-environment) test environment can drastically speed up your review.

## Can I decide a date when my plugin will be published?

* No. We cannot guarantee specific publishing dates, as plugins are immediately published after a review based on our internal prioritization.

## Can I see metrics of my plugin installs, downloads or usage?

* No. We don't offer this information at the moment to plugin authors.

## How can I update my plugin's catalog page?

* The plugin's catalog page content is extracted from the plugin README file. To update the plugin's catalog page, submit an updated plugin with the new content included in the README file.

## Can I unlist a plugin?

* In the event of a bug, unlisting the plugin from our catalog may be possible in exceptional cases, such as security concerns. However, we don't have control over the instances where the plugin is installed.
* Also, refer to the Grafana Labs [Plugin Deprecation Policy](https://grafana.com/legal/plugin-deprecation/) to learn more about plugin deprecation.

## Can I distribute my plugin somewhere else other than the Grafana plugin catalog?

* The official method for distributing Grafana plugins is through our catalog. Alternative methods, such as installing private or development plugins on local Grafana instances, are available as per the guidelines provided in [this guide](https://grafana.com/docs/grafana/latest/administration/plugin-management#install-plugin-on-local-grafana).

## Can I still use Angular for a plugin?

* No. We will not accept any new plugin submissions written in Angular. For more information, refer to our [Angular support deprecation documentation](https://grafana.com/docs/grafana/latest/developers/angular_deprecation/).

## Can I submit plugins built with Toolkit?

* The @grafana/toolkit tool is deprecated. Please [migrate to `create-plugin`](https://grafana.com/developers/plugin-tools/migration-guides/migrate-from-toolkit). In the future, we will reject submissions based on @grafana/toolkit as it becomes increasingly out-of-date.

## Do all plugins require signatures?

* All plugins require signatures unless they are in development or being submitted to review for the first time.

## Do plugin signatures expire?

* Plugin signatures do not currently expire.

## What source code URL formats are supported?

* Using a tag or branch: `https://github.com/grafana/clock-panel/tree/v2.1.3`
* Using a tag or branch and the code is in a subdirectory (important for mono repos): `https://github.com/grafana/clock-panel/tree/v2.1.3/plugin/` (here, the plugin contains the plugin code)
* Using the latest main or master branch commit: `https://github.com/grafana/clock-panel/` (not recommended, it's better to pass a tag or branch)
