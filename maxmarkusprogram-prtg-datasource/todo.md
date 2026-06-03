# "

## First of all: thank you!

This plugin is exactly what the PRTG + Grafana community has been missing. The work you've put into building a proper backend plugin with Go, caching, and all the different query types is really impressive. Thank you for making this available!

---

## Feature request

The pre-built binaries currently committed to the repository appear to be out of sync with the source code. Users who skip the build step run into issues that don't occur when building from the latest source.

It would be great if a GitHub Actions pipeline could automatically build the plugin (frontend + backend) and publish the output as release assets on every push to `main`. That way users can grab a ready-to-use zip from the Releases page without needing Go, Node.js, npm or Mage installed locally.

Thanks again for maintaining this!" 

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
