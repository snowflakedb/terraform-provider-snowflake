# Our roadmap

## (2025-05-22) Roadmap after GA

Recent efforts concentrated on stabilization, SnowflakeDB migration, and the Terraform provider's GA release, which delayed [feature parity](#feature-gap-closing---the-current-approach) and [grants migration](#grants-migration). With the provider now GA and officially supported by Snowflake, we can prioritize further improvements.

In the upcoming months, our focus will be on the following topics:

* [Migration help and solving reported issues for the 2.0.0 versions of the provider](#migration-help)
* [Security-related improvements (PAT support)](#pat-support)
* [Organization accounts adaptation](#organization-accounts)
* Feature gap:
    * [SPCS](#spcs)
    * [Listings](#listings)
    * [Account](#account)
    * [Missing functionalities in stable resources](#missing-functionalities-in-stable-resources)
* [Research and PoC of Plugin Framework](#plugin-framework)
* [Research and PoC of Snowflake REST APIs](#snowflake-rest-apis)

Further topics (these are the currently expected next topics, but they are still subject to change):

* Feature gap:
    * Tables stabilization
    * Iceberg tables
    * Function and procedure stabilization
* Improved error handling
* Improved logging

### Migration help

The v1.0.0 release and later the v2.0.0 GA release were essential steps in maturing the provider. It’s really important for us to move the majority of our users to these versions. The previous [entry already covered multiple reasons](#reasons-to-migrate-to-v1). Also, it’s connected with the [old namespace deprecation](#snowflakedb-migration).

We are committed to resolving your migration challenges. Please contact us if you encounter any migration difficulties. Before starting, be sure to consult the [migration guide](./MIGRATION_GUIDE.md).

### PAT support

Programmatic access tokens are one of the newest authentication methods supported by the Snowflake ecosystem ([2025-04-30-programmatic-access-tokens](https://docs.snowflake.com/en/release-notes/2025/other/2025-04-30-programmatic-access-tokens)). As Snowflake is pushing users away from weak user password authentication, we want to have a wide choice of stronger options inside the provider. You can read more in [Using programmatic access tokens for authentication](https://docs.snowflake.com/en/user-guide/programmatic-access-tokens).

### Organization accounts

Organization accounts were introduced some time ago ([2024-10-01-organization-account](https://docs.snowflake.com/en/release-notes/2024/other/2024-10-01-organization-account)). We believe that the provider's widest use case is account-level work. [Organization accounts](https://docs.snowflake.com/en/user-guide/organization-accounts) offer a new way of managing multi-account setups. We want to allow the provider's users to benefit from their current and future potential.

### SPCS

SPCS is a widely requested area for the Snowflake Terraform Provider (e.g. [#2478](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2478)). We want to provide the following object support (first in preview, later as stable):

* Image repository;
* Compute pool;
* Service.

### Listings

Listings are the object that is also highly anticipated (e.g. [#2379](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2379)). We will share later if we will support all variants right away (e.g. [create-organization-listing](https://docs.snowflake.com/en/sql-reference/sql/create-organization-listing) vs [create-listing](https://docs.snowflake.com/en/sql-reference/sql/create-listing)).

### Account

The [snowflake_account](https://registry.terraform.io/providers/snowflakedb/snowflake/2.1.0/docs/resources/account) resource was reworked [before the V1 release](./MIGRATION_GUIDE.md#snowflake_account-resource-changes). It offers only the creation capabilities and not certain attribute modifications (as explained in [this comment](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3425#issuecomment-2690179698)). We would like to fill this gap.

### Missing functionalities in stable resources

The stable resources have new features (e.g., `resource_constraint` in warehouses [#3258](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3258)). While this roadmap entry won't detail every planned feature, we will focus on ensuring our stable resources offer the same functionality as their SQL counterparts.

### Plugin Framework

Our provider was written using [Terraform SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2). There is a new, [recommended](https://developer.hashicorp.com/terraform/plugin/framework-benefits) way of writing plugins: the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework), which opens more possibilities for plugin writers. This can improve the usability and reliability of our provider (check e.g. [#3327](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3327#issuecomment-2597989866) or [#3366](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3366#issuecomment-2622223398)). We also hope to tackle some limitations listed in the [Decisions before V1](./v1-preparations/CHANGES_BEFORE_V1.md). SDKv2 is not yet deprecated by HashiCorp, but the new features for the Terraform ecosystem are added only to the new framework (e.g. [moved block support](https://github.com/hashicorp/terraform-plugin-sdk/issues/1411#issuecomment-2599203301)).

Transitioning to the new framework will require significant time. Initially, our focus will be on research and developing a PoC for a subset of resources. This will enable us to create a more informed implementation plan.

### Snowflake REST APIs

The current provider was developed due to the initial lack of alternatives for managing Snowflake objects beyond executing SQL commands via a dedicated driver. However, Snowflake now offers [REST APIs](https://docs.snowflake.com/en/developer-guide/snowflake-rest-api/snowflake-rest-api), which present significant advantages for the provider's resources and data sources. Utilizing REST endpoints could lead to improved generation capabilities, faster development cycles, and quicker adoption of new Snowflake features, potentially accelerating the closure of existing feature gaps.

Adopting this new provider represents a significant change. We intend to proceed cautiously, initially focusing on validating its benefits and limitations through a proof of concept involving a limited number of resources. This approach will allow us to assess its viability without impacting the stability of the existing provider, as we anticipate running duplicate resources during this phase. We may also seek user feedback on these new resources to inform our decision-making process.

## (2025-04-23) v2.0.0 release - GA
<a name="23042025-v200-release---ga"></a>

As announced previously, we have released today the v2.0.0 version of the provider, which is an officially supported Snowflake product!

Some essential information:
- We updated the project state in the registry docs and the project's README.
- The [official docs](https://docs.snowflake.com/en/user-guide/terraform) will be refreshed soon.
- Before upgrading, please read the [migration guide](MIGRATION_GUIDE.md#v121--v200) carefully, as there are some breaking changes. We wanted to make them as smooth as possible, so for some of them, there is the "opt-out" option (e.g. skipping the file permission validation can still be set in the provider config through [skip_toml_file_permission_verification](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#skip_toml_file_permission_verification-1) attribute), which will be removed in future major releases.
- For more details on "what is GA?" check the [previous entry](#10042025-ga-announcement). 
- In case of any migration issues, please contact us.
- We will share our next milestones in early May (as the next roadmap entry). One of our main priorities is to help with migrations, so that everyone can benefit from the official support.

As always, we encourage you to contact us with any questions or issues. Happy Terraforming!

## (2025-04-10) GA announcement
<a name="10042025-ga-announcement"></a>

### GA release

If there are no significant obstacles or critical issues, we aim to release GA on **April 23rd**.

#### Migration to the GA version

With the v1.0.0 release, we are no longer making breaking changes in the stable resources without bumping the major version of the provider (with the exception of Snowflake BCRs and preview features).
As we need to introduce some small breaking changes to stable parts of the provider, we will release the 2.0.0 version soon. We target it to be the GA version.

The breaking changes will be added to the [migration guide](MIGRATION_GUIDE.md#migration-guide) as usual. There won’t be a need to migrate the resources. The expected changes are:

* changing a default value in the provider configuration;
* marking some more attributes as sensitive in resource schemas - we will share the full list with recommendations on what to do if the current value is used as input value in another resource.

We do not expect any more breaking changes to be part of the GA release. Please contact us with any questions or problems regarding this migration.

#### Scope

Please check the ["What is GA?](#what-is-ga)" and ["What GA is not?](#what-ga-is-not)" sections of one of the previous announcements.

The main misunderstanding we observe is the stable and preview features distinction. Let us reiterate it:

* resources and data sources that were stabilized during our road to V1 are stable; they will receive the official Snowflake support after GA;
* resources and data sources that were not reworked were marked as preview features. A few clarifications:
    * They were NOT downgraded - they were part of 0.x.x experimental versions of the provider before this classification.
    * We did not want to remove them from the provider entirely as it could block adoption to v1+ versions. That’s why we decided to leave them, even if they may be unstable.
    * If you used any resource or data source in 0.x.x versions, you acknowledged the risks of using experimental features either way.
    * In v1+ versions, you can continue using them, but you need to explicitly allow these features in the provider (as explained [here](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#preview_features_enabled-1)).
    * Remember, Snowflake does not offer official support for these preview resources and data sources. Please consider these features as preview features in the provider, regardless of their state in Snowflake.
    * Each preview feature will be reworked and marked as a stable feature in future releases. Please expect that preview features might introduce breaking changes, even when the provider’s major version number does not change.
* Resources and data sources marked as preview features or not yet supported will be prioritized based on demand and worked on after GA. We will share the current priorities in **early May**.

#### Versioning

Versioning will be clarified in [Snowflake official docs](https://docs.snowflake.com/en/user-guide/terraform) as part of the GA release. Short summary:

* Snowflake Terraform Provider follows semantic versioning.
* Minor version releases may sometimes include unexpected changes, depending on the configuration or environment. We balance the occasional one-time inconvenience for some users against the overall benefits these updates bring to the community.
* Starting with the GA release, Snowflake will offer official support only for the latest version of the provider. When a new version is released, it immediately becomes the officially supported version.
* Although the provider's latest version is the only officially supported version, we make a best effort to resolve issues with earlier versions. After assessing the issue, Snowflake Support may, at its discretion, require an update to the latest version to support the troubleshooting process.
* Generally, we introduce new features and fixes in the latest minor version. This is due to the resource-intensive development process and the need for extensive regression testing.
* If we discover a security vulnerability, we consider backporting critical fixes to earlier versions on a case-by-case basis.
* We assess BCRs introduced by underlying Snowflake features for impacts to the provider. The migration guide provides information about how to manage potential breaking changes We prioritize BCR fixes in each latest version release of the provider and recommend updating your version of the provider regularly.

### snowflakedb migration

The migration to the snowflakedb GitHub organization was successfully concluded. It includes creating the new namespace in the official Terraform registry: [https://registry.terraform.io/providers/snowflakedb/snowflake](https://registry.terraform.io/providers/snowflakedb/snowflake). The [old provider](https://registry.terraform.io/providers/Snowflake-Labs/snowflake) still exists, but as [announced](#24032025-snowflakedb-migration), it will be deprecated soon. The migration is straightforward. You can find details on how to migrate [here](SNOWFLAKEDB_MIGRATION.md).

The old namespace will be deleted on **July 31st**. We will soon add the deprecation notice to the old namespace. Please contact us if you encounter any problems switching to the provider in the new namespace.

## (2025-03-24) snowflakedb migration
<a name="24032025-snowflakedb-migration"></a>

As part of GA preparation, we must migrate out of the [Snowflake-Labs](https://github.com/Snowflake-Labs) GitHub organization to [snowflakedb](https://github.com/snowflakedb), as noted in the [previous Roadmap entry](#what-is-ga). It should happen within a few days. A few important notes on that:
- **The [old terraform registry](http://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs) will still be accessible so no immediate action is needed** (at least for now; we will have separate communication regarding its deprecation). However, don’t wait with the migration, as we plan to deprecate and remove the old registry sooner rather than later.
- There will be a new entry in the registry; it should contain the same releases as the [old one](http://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest/docs) (however, we want to ensure that it includes the releases starting from v0.60.0; if you are on the older version, please migrate to the newer one soon).
- The old GitHub repository ([link](https://github.com/Snowflake-Labs/terraform-provider-snowflake)) will redirect to the same repository located in the snowflakedb organization. All issues, pull requests, and discussions will be persisted.
- When the repository is migrated, we will first verify if the new registry is working correctly and provide detailed instructions on migrating (it should be a straightforward process similar to the description in [#upgrading-from-czi-provider](./CZI_UPGRADE.md#upgrading-from-czi-provider)).

If you have any questions about the migration topic, please contact us through GitHub or your account managers.

We will announce the GA date soon, so stay tuned!

## (2025-02-07) GA scope and roadmap
<a name="07022025-ga-scope-and-roadmap"></a>

### Current focus and goals

Since the last update, we have focused on the following:

* supporting the V1 migration;
* assessing the scope and timeline for the GA.

The biggest migration challenge (ATM) is transitioning from old grants to new ones
(e.g. [#3335](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/3335)).
Check the [grants migration](#grants-migration) section for more details.

### GA

We took a long road to stabilize the provider and recently got it to [V1](#13122024-v1-release-update).
The next essential step in making the provider official and supported by Snowflake is bringing it to GA.

#### What is GA?

Because of the project’s long history, we were asked multiple times about the difference between GA and V1.

The GA of the Snowflake Terraform Provider will mean:

* having official Snowflake support (ability to submit official Support Cases for the Provider);
* migrating the project to the [snowflakedb](https://github.com/snowflakedb) GitHub organization
(we are still in [Snowflake-Labs](https://github.com/Snowflake-Labs), reserved for unofficial/experimental projects).

The above will mean changes in the support process and the provider setup
(most probably the change in the [registry](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs)).
We will share the details in the upcoming weeks.

Because of the recent V1 and upcoming GA, we will also clarify our versioning policies (e.g. how long the given version is supported).

**Important**: it will apply to the v1.x.x+ versions of the provider, so it’s essential to upgrade to the v1.0.0 version as soon as possible
(read more in the [migration](#will-migration-be-needed) section below).

#### What GA is not?

There is a common misconception of what will be supported in the provider’s GA.

The GA of the Snowflake Terraform Provider does NOT mean that:

* all features are stable;
* all Snowflake GA objects are supported.

Functionally, the provider will offer almost the same set of objects as the [recently released V1](#13122024-v1-release-update).
Read more about the [feature gap closing](#feature-gap-closing---the-current-approach) below.

#### Timeline

We aim to reach GA by the end of May 2025. We will update the timeline in mid-March.

#### Will migration be needed?

There will be the following migrations involved:

1. Getting to v1.0.0.

    It should already be an ongoing process. The V1 version offers stability and will be the basis for official support.

    Remember to follow our [migration guide](./MIGRATION_GUIDE.md#migration-guide) closely, as there were many breaking changes between the 0.x.x versions.
    Reach out to us if you have any problems with it.

2. Getting to v1.x.x.

    As mentioned in the [What is GA](#what-is-ga) section, official support will start with one of the 1.x.x versions (we will announce the precise version later).
    This migration should be easy because we don’t plan to introduce breaking changes between the 1.0.0 and 1.x.x versions.

    Remember that enabling [preview features](https://registry.terraform.io/providers/snowflakedb/snowflake/1.0.0/docs#preview_features_enabled-3)
    in the provider’s configuration may result in manual migration as these features do not offer stable schemas.

3. Changes in the terraform config files.

    Because we have to migrate the project from [Snowflake-Labs](https://github.com/Snowflake-Labs) to [snowflakedb](https://github.com/snowflakedb),
    we will also most probably create a new registry instead of the [existing one](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs).
    This migration should be painless (basically changing the required provider block and running the `terraform state replace-provider`,
    similar to the [#upgrading-from-czi-provider](./CZI_UPGRADE.md#upgrading-from-czi-provider)).

    We will share official instructions closer to the GA release date.

### Grants migration

We considered adding small migration helpers before going to GA (check [this discussion](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/3335#discussioncomment-11799443));
however, we don’t have the resources to work on them in parallel with GA preparations.
We might be able to publish more helpful examples and pointers instead.
We encourage [contributing](./CONTRIBUTING.md) to the project.

Our idea would be to create simple scripts that can be reused, or that can at least serve as a base/inspiration for the user-side migration automation.
They may contain:

* printing the desired config for new grants based on the old grant resource input;
* printing the desired config for new grants based on the output from the Snowflake query;
* generating proper import statements (specifically generating correct identifiers).

We will treat this topic as a high-priority nice-to-have before the GA, and an essential topic right after reaching GA.

### Reasons to migrate to v1+

While we can’t make anyone migrate to the newer versions of the provider, we would like to point out a few things:

* Snowflake is not officially supporting the Snowflake Terraform Provider project yet.

    We put the disclaimers everywhere but have learned that it’s not always enough. **It won’t change for the 0.x.x versions after reaching GA**.
* The old versions (0.x.x) will not be back-fixed (our policy before v1.0.0 was that we were always introducing fixes only in the newest 0.x.x versions; examples: [comment](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2982#issuecomment-2296211672) and [comment](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2739#issuecomment-2071555398)).

    It means that **pre-GA versions can break entirely** when a breaking change is introduced on the Snowflake side
    (the provider works on SQL statements; if the syntax changes through an official BCR, we won’t provide patches to unsupported versions, basically making them inoperable).
* New features will only be introduced in the newest versions.
* Migrating to v1.0.0 may be challenging, but there won’t be any breaking changes in stable resources, and no resource removals are planned until v2.0.0, which is not planned to be released anytime soon.
* The engineering team handles the current support directly on a best-effort basis. The GA versions, which will be officially supported by Snowflake, will enable quicker triage and response.

### Feature gap closing - the current approach

As part of the V1 release, we have introduced a distinction between stable and preview resources (check [the previous update](#13122024-v1-release-update)).
In addition to the preview resources that need to be stabilized, some objects have not yet been created in the provider
(e.g. [iceberg tables](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2249)
or [listings](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2379)).
Additionally, the stable resource can also be subject to changes when new attributes are added on the Snowflake side to the object.

We are not backing off the strategy to ultimately support all Snowflake objects.
However, reaching GA is our highest priority now, and all the feature-related work will be postponed until after GA.
We will discuss feature priorities after reaching GA.
Please contact us through your account managers if any feature is critical.

The same applies to the non-critical issues where a workaround exists.
We will still fix the critical issues as part of our best-effort support.

## (2024-12-13) V1 release update
<a name="13122024-v1-release-update"></a>

We have released a long-awaited [v1.0.0](https://github.com/snowflakedb/terraform-provider-snowflake/releases/tag/v1.0.0). A few things to know now:
- Together with v1.0.0 we have also released "the last" 0.x.x version - 0.100.0. v1.0.0 is built on top of that; it removed the [deprecated resources](https://github.com/snowflakedb/terraform-provider-snowflake/blob/ab015e8cf6f4db762b4043e8bfce2a010b623602/v1-preparations/LIST_OF_REMOVED_RESOURCES_FOR_V1.md) and attributes mostly, so if you are using one of the latest 0.x versions, you should be really close to v1.
- Check the migration guides for [v1.0.0](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#v01000--v100) and [v0.100.0](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#v0990--v01000).
- The provider entered a stable version from the engineering point of view. It will prohibit us from introducing breaking changes in stable resources without bumping the major version.
- Resources and data sources in our provider now have two states, [stable](https://github.com/snowflakedb/terraform-provider-snowflake/blob/ab015e8cf6f4db762b4043e8bfce2a010b623602/v1-preparations/LIST_OF_STABLE_RESOURCES_FOR_V1.md) and [preview](https://github.com/snowflakedb/terraform-provider-snowflake/blob/ab015e8cf6f4db762b4043e8bfce2a010b623602/v1-preparations/LIST_OF_PREVIEW_FEATURES_FOR_V1.md). To allow the given preview feature you have to explicitly set it in [the provider config](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#preview_features_enabled-1). Please familiarize yourselves with the limitations of the preview feature before enabling it (most notably, preview features may require migrations between minor versions).
- Our current main goal is to help with migration and address all the incoming v1 issues.
- Keep in mind that V1 does not mean we have an official Snowflake support (check our new disclaimer in [README](https://github.com/snowflakedb/terraform-provider-snowflake?tab=readme-ov-file#snowflake-terraform-provider)).
- Our next milestone is reaching GA, which requires mostly procedural steps. Before that, no big changes are planned for the provider.
- Besides the GA, we want to focus mostly on stabilizing the preview resources. We will share their current prioritization in January. The main ones for now are functions, procedures, and tables.

## (2024-10-25) Project state overview
<a name="25102024-project-state-overview"></a>

### Goals

Since the last update we have focused on:

* [Reducing the feature gap](#reducing-the-feature-gap) (focusing on the Snowflake essential GA resources)
* Redesigning identifiers (check [\#3045](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/3045) and [identifiers_rework_design_decisions](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/docs/guides/identifiers_rework_design_decisions.md))
* Reworking the provider's configuration (the doc/discussion will be shared when ready)
* Researching the object renaming in our provider (the doc will be shared when ready)

These steps were all needed to get us closer to the first stable version of the provider which... is really close. In the next 1-2 months we want to:

* [Wrap up the functional scope](#wrap-up-the-functional-scope) (not all the objects will be declared stable, more details below)
* [Prepare for the V1 release](#prepare-for-the-v1-release)
* [Prepare some basic performance benchmarks](#prepare-some-basic-performance-benchmarks) (especially, after a few major changes to the resources logic)
* [Improve/update the documentation](#improveupdate-the-documentation)
* [Run a closed early adopter program](#run-a-closed-early-adopter-program) to verify the readiness of the provider to enter a stable V1

If there won't be any major obstacles or critical issues we aim to release V1 on **December 9th**. To better understand its scope, please check the ["What is V1?"](#what-is-v1) section.

#### Reducing the feature gap

During the last six months, we have been tackling objects from the [essential](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/v1-preparations/ESSENTIAL_GA_OBJECTS.MD) and [remaining](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/v1-preparations/REMAINING_GA_OBJECTS.MD) object lists. We’ve been aligning the implementation, adding missing attributes, and fixing known issues of the chosen objects (full list below). We had to make design decisions that sometimes were not only dictated by our engineering assessments but also by the limitations of Terraform and the underlying [SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2). The main decisions are listed inside the repository in the [Design decisions before v1](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md#design-decisions-before-v1) (we will validate if all essential ones are present there before releasing V1).

#### Wrap up the functional scope

It’s about finishing the redesign of objects we want to declare stable. This mainly affects tables and accounts, but it also involves small alterations in other objects (which will be listed in the migration guide as usual).

As shown [below](#which-resources-will-be-declared-stable), all but one of the [essential](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/v1-preparations/ESSENTIAL_GA_OBJECTS.MD) objects and a few of the [remaining](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/v1-preparations/REMAINING_GA_OBJECTS.MD) objects made the cut.

#### Prepare for the V1 release

This is mainly cleaning up the repository but also activities around the release:

* removing deprecated resources
* marking the resources as [preview features](#preview-resourcesdatasources)
* removing deprecated attributes
* potentially renaming some configuration options
* summarizing migration guidelines between v0.x.x and v1.0.0

#### Prepare some basic performance benchmarks

During the resources redesign we introduced multiple changes that may affect the performance. Namely:

* more SQL statements are run (`SHOW`, `DESCRIBE`, and `SHOW PARAMETERS` when needed)
* the state we save is bigger because of the `show_output`, `describe_output`, and `parameters`.

We observed that our customers tend to have lots of objects in single terraform deployments. This leads to longer planning and execution times. To be able to guide “what is too much”, we need to perform tests with more objects on our end.

#### Improve/update the documentation

We greatly improved the docs and the transparency of the project. However, there are still topics that need our attention (e.g. adding a migration guide directly to the [registry](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs), adding missing design decisions like granting ownership, or adding more guides \- similar to [identifiers rework](https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs/guides/identifiers) \- like importing existing infrastructure into Snowflake).

#### Run a closed early adopter program

We planned V1 to be as close as possible to the latest 0.x.x version before the V1 release. However, some changes/migrations are still expected. To improve confidence, we have decided to provide early V1 binaries to early adopters. We are still actively recruiting customers; please reach out to your Snowflake Account Manager at the earliest if you would like to participate. The program runs from mid-November to mid-December.

#### What is V1?

The first major version, V1, marks the first step in getting to GA by providing stable versions to customers who use the provider. We hope to have all our current customers migrate to V1. The provider's Product and Engineering teams will be available for migration or any other questions, as we believe this migration is key in preparing our customers for seamless GA adoption.

From the engineering point of view, the provider will be in the stable version, but it will still stay in the Snowflake-Labs GitHub organization. We plan to change that and move it to the official snowflakedb org so that it gets the official Snowflake support. This will be a necessary step to reach the GA.

#### Which resources will be declared stable

Check [this list](v1-preparations/LIST_OF_STABLE_RESOURCES_FOR_V1.md) for details.

#### Preview resources/datasources

On our road to V1, we went through the resources, starting with the most used ones. We did not cover all of them (as described above). Because of that, in the newest [v0.97.0](https://registry.terraform.io/providers/snowflakedb/snowflake/0.97.0/docs) version of the provider, we have multiple resources that were not redesigned/fixed.

We discussed two main options: removing them from 1.0.0 or marking them as preview features. We were mostly worried that removing resources would prevent the majority of our users from migrating to the stable version. On the other hand, we know they are not ready so we don’t want to declare them as stable.

After consideration, we decided to leave them as preview features that need to be **explicitly enabled by the user**. This way, we are not reducing the provider's functionality between v0.x.x and v1.0.0 and leave the possibility to use them while accepting the limitations they have. However, these resources will be subject to change after V1. They should be treated as [Snowflake Preview Features](https://docs.snowflake.com/en/release-notes/preview-features) so changes to their schemas (breaking changes included\!) may be introduced even without bumping the major version of the provider.

#### “Attachment” resources clarification

During our road to V1 we tried to limit the number of resources needed to be configured in order to manage the given Snowflake object correctly. Because of that, we moved [Snowflake parameters](https://docs.snowflake.com/en/sql-reference/parameters) handling directly to the given object’s resource (check [this](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/v1-preparations/CHANGES_BEFORE_V1.md#snowflake-parameters)). We did that to other types of properties too (e.g. we changed the logic for public keys handling in the [snowflake_user](https://registry.terraform.io/providers/snowflakedb/snowflake/0.95.0/docs/resources/user#rsa_public_key) resource, so that [snowflake_user_public_keys](https://registry.terraform.io/providers/snowflakedb/snowflake/0.94.1/docs/resources/user_public_keys) is no longer compatible with it).

Still, these “attachment” objects serve a specific use case (i.e. the main object is not managed by Terraform but part of the object may be). It opened a question for the future not only because of the aforementioned use case but also because of a wider perspective on the default resource behavior. For example, a resource monitor can be attached to a warehouse only by a user with an ACCOUNTADMIN role (check [\#3019](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3019)). Some of our users would like to provision warehouses separately from assigning resource monitors but the caveat here is that leaving the resource monitor empty in the resource config will currently remove any assigned resources. Handling this would require adding a separate attachment resource and allowing a conditional change in behavior for empty assignments in the main object.

The topic is wide. For the V1, we decided to keep most of the attachment resources as [preview features](#preview-resourcesdatasources) and we will discuss the need for handling the use cases described in this section as a separate topic after V1.

#### Which resources will be left as preview features

Check [this list](v1-preparations/LIST_OF_PREVIEW_FEATURES_FOR_V1.md) for details.

#### Which resources will be removed

Check [this list](v1-preparations/LIST_OF_REMOVED_RESOURCES_FOR_V1.md) for details.

#### Roadmap short after V1

Right after V1, we would like to focus on helping all of you with the migration. We will prioritize it so we encourage you to approach us with any issues you might have.

In the meantime, if we have enough time, we want to prioritize redesigning the object marked as preview features. Currently, stages and shares open the list.

#### Next year priorities

This is only a general overview of the next year and may be subject to change:

* Graduate out of Snowflake-Labs into the official snowflakedb organization
* GA of the Snowflake Terraform Provider
* Research performance improvements (optimize Snowflake invocations)
* Grants improvements
* Redesign remaining GA objects
* Design transition to the [plugin framework](https://developer.hashicorp.com/terraform/plugin/framework)
* Introduce Terraform modules

## (2024-05-05) Roadmap Overview
<a name="05052024-roadmap-overview"></a>

### Goals
Since the last update we have focused on:
- [Finishing the SDK rewrite](#finishing-sdk-rewrite).
- [Redesigning grants](#redesigning-grants) (check announcements: [discussions/1890#discussioncomment-9071073](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/1890#discussioncomment-9071073), [discussions/2235](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/2235), and [discussions/2736](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/2736)).
- Improving the provider’s stability (by [stabilizing the tests](#tests-stabilization), solving new incoming issues on a daily basis, and [introducing repository-wide fixes to multiple objects](#resolving-existing-issues)).
- Preparing the scope for the V1 (more below). Part of [supporting-all-snowflake-ga-features](#supporting-all-snowflake-ga-features).
- Raising the transparency of the project (this roadmap, [contribution guidelines](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/CONTRIBUTING.md), [old issues cleanup](https://github.com/snowflakedb/terraform-provider-snowflake/discussions/2755), and [FAQ](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/FAQ.md)).

The primary goals we are working on currently are:
- Introducing support for the fundamental GA features and improving the existing objects (resolving existing provider issues included). Continuation of [supporting-all-snowflake-ga-features](#supporting-all-snowflake-ga-features).
- Reworking identifiers.

The more concrete topics we are currently dealing with are presented in the following three sections: current, upcoming, and next.

|                      Current                       |                             Upcoming                             |                  Next                   |
|:--------------------------------------------------:|:----------------------------------------------------------------:|:---------------------------------------:|
| Preparing essential GA objects for the provider V1 | Preparing rest of the fundamental GA objects for the provider V1 |    Official Snowflake public preview    |
|                 Identifiers rework                 |                 Support object renaming properly                 |        Enable Snowflake support         |
|                                                    |                 Provider’s configuration rework                  | Support for the public preview features |
|                                                    |                      Prepare the V1 release                      |                                         |

#### Current (expected mid/late this year)
##### Preparing essential GA objects for the provider V1
As we stated in the [previous entry](#supporting-all-snowflake-ga-features) we want to inspect all the existing objects to find missing parameters and flaws in their designs. We gathered a list of objects we believe are most crucial, and we will address them first. The list is available [here](v1-preparations/ESSENTIAL_GA_OBJECTS.MD). After them, we will address the ones described in the [following entry](#preparing-the-rest-of-the-fundamental-ga-objects-for-the-provider-v1).

##### Identifiers rework
([previous entry](#identifiers-rework-1)) Identifiers were recently the second, next to the Grants, most common error source in users’ configurations. We want to make interaction with them easier (at least to the extent we have control of).

#### Upcoming (expected likely in late Fall)
##### Preparing the rest of the fundamental GA objects for the provider V1
This will be the continuation of [Preparing essential GA objects for the provider V1](#preparing-essential-ga-objects-for-the-provider-v1). It will address objects listed [here](v1-preparations/REMAINING_GA_OBJECTS.MD).

##### Support object renaming properly
Object renaming is a topic that arises in different contexts like renaming a database, column, or schema object to name a few. The renaming topic was brought up a long time ago, e.g. in [#420](https://github.com/snowflakedb/terraform-provider-snowflake/issues/420), [#753](https://github.com/snowflakedb/terraform-provider-snowflake/issues/753), and [this forum entry](https://community.snowflake.com/s/question/0D5Do00000KWFhqKAH/how-to-rename-snowflake-database-on-terraform-with-the-snowflake-provider). We want to address the renaming in general before the stable V1.

##### Provider’s configuration rework
([previous entry](#providers-configuration-rework-1)) It is one of the last moments before going V1 to make incompatible changes in the provider. The current configuration contains many deprecated parameters, inconsistencies with the documentation, and other design flaws. We want to address it.

##### Prepare the V1 release
This will be the moment to validate our V1 efforts by checking if everything was implemented and making the migration for all of you as smooth as possible. This includes:
- Listing of all breaking changes
- Summarizing the migration notes
- Communicating the V1 release in detail
- Describing the new release cycle post-V1
- And many more…

**Important** We plan to introduce the changes before the V1 to allow you to migrate most of the objects before the official release. Because we are still not providing the backward bugfixes, it's always best to bump the provider version with the new releases (following the [migration guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/MIGRATION_GUIDE.md#migration-guide)).

#### Next
- Official Snowflake public preview
- Enable Snowflake support
- Support for the public preview features

## (2024-01-19) Roadmap Overview
<a name="19012024-roadmap-overview"></a>
### Goals
The primary goals we are working on currently are:
- Adding missing and updating existing functionalities (resources and data sources);
- Resolving existing provider issues;
- Improving provider’s stability.

We believe fulfilling these goals will help us reach V1 with a stable, reliable, and functional provider. The more concrete topics we are currently dealing with are presented in the following three sections: current, upcoming, and next.

|                                                                                     Current                                                                                      |                                                                                              Upcoming                                                                                               |                                                                         Next                                                                          |
|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-----------------------------------------------------------------------------------------------------------------------------------------------------:|
|                                                                                 Redesign Grants                                                                                  | Design proper resources for the majority of Snowflake objects. Support all Snowflake GA features, starting with the most critical resources like databases, schemas, tables, tasks, and warehouses. | Rework provider’s configuration. Covers current configurations deprecated parameters, inconsistencies with the documentation, and other design flaws. |
| Finish SDK rewrite. Migrate existing resources and data sources to the new SDK, aiding in safer and more extendable generation of SQL statements executed against Snowflake API. |                                                                                         Rework identifiers                                                                                          |                                     Stabilization of tests, ensuring quicker development and stability assurance.                                     |
|                                                     Address open issues in the repo repository, focusing on critical issues.                                                     |                                                              Address open issues in the repo repository, focusing on critical issues.                                                               |                                                                                                                                                       |

#### Current
##### Redesigning Grants
Grants proved to be one of the most common pain points for the provider’s users. We have been focusing on designing the proper resources for the past few weeks. The development is in progress, but more topics still need our attention (like granting ownership, and imported privileges, to name a few).

##### Finishing SDK rewrite
Last year, we changed the approach to generating the SQL statements executed against Snowflake API. The previous, old implementation was error-prone and hard to maintain. We are concluding migrating existing resources and data sources to the new SDK we are developing. It has already proved to be safer and more extendable.

##### Resolving existing issues
Having the ~470 open issues in the repository is not fun. We want to reduce that number drastically. We have recently taken multiple different steps to achieve it:
- We respond to most of the incoming issues faster.
- We classified and prioritized the existing issues. We picked the resources that were causing the most trouble for our users. We will focus first on resource monitors, databases, and tasks. At the same time, we introduce improvements in reporting errors and handle common pitfalls globally.
- We plan to close the issues regarding ancient provider versions. There will be a separate announcement about it.

#### Upcoming
##### Supporting all Snowflake GA features
Eventually, we want to support all Snowflake features. We first want to support all the GA ones. It does not only mean that we will add the missing resources; we will also carefully inspect the existing ones to find missing parameters and flaws in their designs. We will start with the most critical resources like databases, schemas, tables, tasks, and warehouses.

##### Identifiers rework
Identifiers were recently the second, next to the Grants, most common error source in users’ configurations. We want to make interaction with them easier (at least to the extent we have control of).

##### Increasing transparency and involving the community in discussions
We are actively being asked about the state of the development, plans for introducing new resources, and design decisions. This roadmap is one of the many steps we are willing to take to be more transparent to our users.

#### Next
##### Provider’s configuration rework
It is one of the last moments before going V1 to make incompatible changes in the provider. The current configuration contains many deprecated parameters, inconsistencies with the documentation, and other design flaws. We want to address it.

##### Tests stabilization
We are extensively testing our provider. We rely on our tests when introducing new features. Unfortunately, historically, testing was not the biggest concern in the project; many tests are missing, and existing ones are not always correct. Having reliable test sets is essential for quicker development and stability assurance.
