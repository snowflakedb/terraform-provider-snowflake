# Changelog

## [2.7.0](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v2.6.0...v2.7.0) (2025-09-11)


### 🎉 **What's new:**

* Add generation field to warehouses ([#4002](https://github.com/snowflakedb/terraform-provider-snowflake/issues/4002)) ([f8c0036](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f8c003679bd82f741468c8848890a74dab2f5afd))
* Add resource_constraint to the warehouse resource ([#3880](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3880)) ([9a88420](https://github.com/snowflakedb/terraform-provider-snowflake/commit/9a88420a38212dbdf1bc36b013b2b863f4ee90f0))


### 🔧 **Misc**

* Add 2025_05 bundle migration guide ([#3974](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3974)) ([46ddcef](https://github.com/snowflakedb/terraform-provider-snowflake/commit/46ddcef751719b57105bbb40d73ba04b28aa6e3c))
* Add a note about dropping the primary account role ([#4017](https://github.com/snowflakedb/terraform-provider-snowflake/issues/4017)) ([8074b8b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8074b8bb59092812aee3be1cf2d3abd1d66e95bc))
* Add a note about the auth policy change in recent bcrs ([#4021](https://github.com/snowflakedb/terraform-provider-snowflake/issues/4021)) ([2d56c48](https://github.com/snowflakedb/terraform-provider-snowflake/commit/2d56c48041d0f8ded128b649a48a7d3578c8194e))
* Add gen2 warehouses feature to the migration guide ([#4022](https://github.com/snowflakedb/terraform-provider-snowflake/issues/4022)) ([82744c8](https://github.com/snowflakedb/terraform-provider-snowflake/commit/82744c83906072b39aff15efdc537bca9f0caf38))
* Add migration script ([#3948](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3948)) ([bd967af](https://github.com/snowflakedb/terraform-provider-snowflake/commit/bd967afb91f19fac8629f5bd4b0af34c6be69f17))
* Adjust existing and add missing sweepers ([#4005](https://github.com/snowflakedb/terraform-provider-snowflake/issues/4005)) ([ec83457](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ec83457b98e629e8bdebe14b1ce108d4b184f469))
* Extract account level tests for several objects ([#3996](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3996)) ([66c0b78](https://github.com/snowflakedb/terraform-provider-snowflake/commit/66c0b787c7006737477861dc3af705b5d11ead23))
* Fix warehouse long-running query test ([#4004](https://github.com/snowflakedb/terraform-provider-snowflake/issues/4004)) ([7f23a5f](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7f23a5fd426e24cee47bc111e21dcff409ba21af))
* Prove session parameters behavior, task test problems, and add PoC for better provider testing cache ([#3997](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3997)) ([8249662](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8249662cc22979ce0dd6c16285de5389c6704d01))
* Stabilize service integration and acceptance tests ([#4000](https://github.com/snowflakedb/terraform-provider-snowflake/issues/4000)) ([ddd7b86](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ddd7b864ba501fc795febe47c013948d12521542))
* Test account parameter mapping ([#3989](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3989)) ([6e7f2b3](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6e7f2b376afe7d97f46c3ea09942ecf27ad8560c))
* Unskip or remove BCR/functionality missing related tests ([#3988](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3988)) ([a581f17](https://github.com/snowflakedb/terraform-provider-snowflake/commit/a581f170c6af710a56e57ba22404d7c0a7178044))


### 🐛 **Bug fixes:**

* Fix diff suppression for identifier set ([#4007](https://github.com/snowflakedb/terraform-provider-snowflake/issues/4007)) ([7b2de9d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7b2de9d7d114d5707e07fbf42e5f1e5f214dfe2e))
* Fix integration and acceptance tests ([#4016](https://github.com/snowflakedb/terraform-provider-snowflake/issues/4016)) ([77f03ea](https://github.com/snowflakedb/terraform-provider-snowflake/commit/77f03ea5ea011f1e9451de5e2cf1c0a03aadcb0c))
* Fix privileges validation in `grant_privileges_to_account_role` resource ([#3994](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3994)) ([0a3816f](https://github.com/snowflakedb/terraform-provider-snowflake/commit/0a3816f451ae163d9c6d13b0e93a53d7c7131c77))

## [2.6.0](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v2.5.0...v2.6.0) (2025-08-29)


### 🎉 **What's new:**

* Add semantic view definition ([#3925](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3925)) ([0b1b68d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/0b1b68d466aebbc27bc459a235d3c1dca8b3d8c6))


### 🔧 **Misc**

* Add manual test for consumption billing entity ([#3918](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3918)) ([321852d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/321852d7b4860fcd199583d28f0c733ebb679af7))
* Add new database to the protected databases in sweepers ([#3927](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3927)) ([50fc2dd](https://github.com/snowflakedb/terraform-provider-snowflake/commit/50fc2dddb8a72af640501cd76bc57ea2015a86d7))
* Add new roadmap entry ([#3935](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3935)) ([6bfb929](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6bfb92999ba9ece400991feede9c7ad41f0ea899))
* Address skipped tests (SNOW-2010844) ([#3972](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3972)) ([dbad401](https://github.com/snowflakedb/terraform-provider-snowflake/commit/dbad4013fe2bb254dbeb3ddc2a03588df6912502))
* Adjust architests after moving acceptance tests to a separate package ([#3970](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3970)) ([1379449](https://github.com/snowflakedb/terraform-provider-snowflake/commit/13794492985c1e4a855abb0f654219f4d72f3c6f))
* Adjust PAT alter after fixes in snowflake ([#3922](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3922)) ([b76a0ea](https://github.com/snowflakedb/terraform-provider-snowflake/commit/b76a0eadbc1c1d77d0604fb6466baa8d99cd4590))
* Adjust security integration tests ([#3975](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3975)) ([6067d97](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6067d97721869d156693b7c587159a78a68c7de7))
* Change convert signature in all other resources ([#3940](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3940)) ([962e0a4](https://github.com/snowflakedb/terraform-provider-snowflake/commit/962e0a4176c8bfadb5b5e4468b29f78b527b8c35))
* Clean up all TODOs created during Terraform Plugin Framework and RestAPI research ([#3963](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3963)) ([4ea4e8b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/4ea4e8b87225d983ac6b7e1e3c5630b0256c0ef9))
* Clean up TODO comments left with "next PR" without issue number ([#3964](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3964)) ([3d7b7bf](https://github.com/snowflakedb/terraform-provider-snowflake/commit/3d7b7bf6419c0f245d3e86b4901473f44e054a36))
* Convert function adjustments follow-up ([#3955](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3955)) ([ce144b1](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ce144b12c724cec725cc3fd84e62baf365e62383))
* Find versions of provider in gh organization ([#3932](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3932)) ([c3e6051](https://github.com/snowflakedb/terraform-provider-snowflake/commit/c3e6051a4dacfb75a544089a0030af1c5f0fc1a9))
* Fix check destroy logic ([#3967](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3967)) ([2ca61f0](https://github.com/snowflakedb/terraform-provider-snowflake/commit/2ca61f0f24b2730f679e646521882be7c0d32e87))
* Fix unsafe query integration test ([#3968](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3968)) ([07e3924](https://github.com/snowflakedb/terraform-provider-snowflake/commit/07e392441e8fce1149df8ce2edfbff910263f642))
* Move more task tests to account level ([#3971](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3971)) ([74f781e](https://github.com/snowflakedb/terraform-provider-snowflake/commit/74f781e1bc1994da60fd22d2aeedb40ef7592beb))
* Nuke roles, databases, warehouses, shares, and network policies ([#3976](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3976)) ([678de64](https://github.com/snowflakedb/terraform-provider-snowflake/commit/678de64fa96c50170737208a4712668c0c2c88c4))
* Remove unnecessary TestAccPreCheck ([#3966](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3966)) ([279b622](https://github.com/snowflakedb/terraform-provider-snowflake/commit/279b622d27d3740557fab6dbc2ebf317f988de36))
* Run functional tests after setting up Terraform ([#3962](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3962)) ([3e186e7](https://github.com/snowflakedb/terraform-provider-snowflake/commit/3e186e7dbd7903852492d0078620c5e06d2bc363))
* Small failover group fixes ([#3947](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3947)) ([cebfca6](https://github.com/snowflakedb/terraform-provider-snowflake/commit/cebfca6f4170901bc8259aa0efd0dc8466b400e0))
* Sweep security integrations ([#3950](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3950)) ([62f79db](https://github.com/snowflakedb/terraform-provider-snowflake/commit/62f79dbe33f2d3ae28b2728982dfaf975c72e129))
* Update failover group docs ([#3949](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3949)) ([12dc244](https://github.com/snowflakedb/terraform-provider-snowflake/commit/12dc24414c96b9f65d22d65e7f16a2ed9be41b4b))
* Upgrade go Snowflake driver ([#3980](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3980)) ([0ed3c60](https://github.com/snowflakedb/terraform-provider-snowflake/commit/0ed3c607af36d19dfb3801b39285b4aedfd50a77))
* Use DropSafely in test setup helpers ([#3969](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3969)) ([7f5d4ca](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7f5d4ca7ba600d1f015d026b1cc5ccb34b205d5b))


### 🐛 **Bug fixes:**

* Change convert signature for a single object ([#3926](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3926)) ([edec510](https://github.com/snowflakedb/terraform-provider-snowflake/commit/edec5107cb87ece817e4ba9076db4aaa2455a085))
* Connection fixes follow-up ([#3957](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3957)) ([5647304](https://github.com/snowflakedb/terraform-provider-snowflake/commit/5647304b65d31c8bcd699ec2fcf4b8157ed35aa6))
* Enhance connection management and logging for secondary connections ([#3937](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3937)) ([8fbcc85](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8fbcc856c0843a7d79142f790e441aa96fdb5000))
* Fix updating privileges in grant_privileges_to_account_role ([#3953](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3953)) ([70136ea](https://github.com/snowflakedb/terraform-provider-snowflake/commit/70136ea7bfb6b9251d9b93620edeadc29f218b3b))
* Update migration guide, import docs, and task tests ([#3933](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3933)) ([ca1568d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ca1568d5aa88f429a438307151cbd29b215d9ef2))

## [2.5.0](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v2.4.0...v2.5.0) (2025-08-07)


### 🎉 **What's new:**

* Add listing resource implementation ([#3893](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3893)) ([75ba6a7](https://github.com/snowflakedb/terraform-provider-snowflake/commit/75ba6a757d33d83b32918952163ead774ad1acad))
* Support setting `storage_aws_external_id` in storage integration resource ([#3659](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3659)) ([62d561b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/62d561b85487d8329f05a1bfa53d4fc75cec5c0c))


### 🔧 **Misc**

* Account test assert follow up ([#3815](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3815)) ([64d3a07](https://github.com/snowflakedb/terraform-provider-snowflake/commit/64d3a076e6eb624228a7bde7c90667c1a88f504d))
* Add 2025_04 BCR migration guide ([#3855](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3855)) ([24aaa44](https://github.com/snowflakedb/terraform-provider-snowflake/commit/24aaa44ad451411a2906ae6b964958afc080df32))
* Add a functional test for nested sensitives in the framework ([#3886](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3886)) ([f1eb135](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f1eb135183ae64063ad996b8c3a64f8b9e80b349))
* Add note that TOML file is distinct from CLI and SnowSQL ones ([#3775](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3775)) ([9e45bc5](https://github.com/snowflakedb/terraform-provider-snowflake/commit/9e45bc5a2ba25f0d563263e85ed394db56c4b592))
* Add open api spec for warehouse ([#3892](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3892)) ([df8c26c](https://github.com/snowflakedb/terraform-provider-snowflake/commit/df8c26c41adbe0483dff1764ffbe8a666f6f063a))
* Add REST API PoC for warehouse resource ([#3891](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3891)) ([901cc4a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/901cc4aa83da6eaec9008c3498905f7703d70655))
* Adjust grant privileges to share test ([#3905](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3905)) ([c603137](https://github.com/snowflakedb/terraform-provider-snowflake/commit/c603137986623b88028009f34befe32bfa84f6ea))
* Adjust or remove tests using enabled by default bcr bundles ([#3909](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3909)) ([e6bf3ac](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e6bf3ac0c9a6b6a1e8270cbe626c48812ae2c9ee))
* Checkout to the commit SHA from the slash command ([#3897](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3897)) ([42f6465](https://github.com/snowflakedb/terraform-provider-snowflake/commit/42f6465cdaca70e56406758807c777230a6f9bcf))
* Documentation improvements ([#3901](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3901)) ([94cc28a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/94cc28a16584c3e9284ca7abcedf67aab14ea5aa))
* Fix migration guide ([#3841](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3841)) ([9d34d89](https://github.com/snowflakedb/terraform-provider-snowflake/commit/9d34d89e138425f28c9628e21819df1c8af42688))
* Generate warehouse resource from OpenAPI with HashiCorp tools ([#3895](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3895)) ([d08db34](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d08db34310d29171a8ebb9f19fde727d4556ebd5))
* Implement PoC for managing warehouse using terraform plugin framework ([#3888](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3888)) ([69f9c96](https://github.com/snowflakedb/terraform-provider-snowflake/commit/69f9c96bd1fe005572a6ec62016df615164bd2ce))
* Listing basic resource with autogenerated content ([#3885](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3885)) ([c11882a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/c11882a1590e43ee6dd213203a7269d87ed79dd6))
* Listings integration tests ([#3873](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3873)) ([640cf8c](https://github.com/snowflakedb/terraform-provider-snowflake/commit/640cf8cc9b5d07ebaef0280946ec3cc48121aa96))
* Remove old framework config ([#3889](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3889)) ([7a8d92f](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7a8d92f1af410e819bfb9ee14e20d89635de01a3))
* Simplify listing manifests ([#3902](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3902)) ([3e9cdf9](https://github.com/snowflakedb/terraform-provider-snowflake/commit/3e9cdf97f139615ce4676f8a3d67dbace0953765))
* Sweepers and Test workflow changes ([#3881](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3881)) ([302dc88](https://github.com/snowflakedb/terraform-provider-snowflake/commit/302dc884594da25749ff28d517a4d71e8a0d87dd))


### 🐛 **Bug fixes:**

* Fix acceptance tests identifier setup ([#3911](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3911)) ([93d4d9c](https://github.com/snowflakedb/terraform-provider-snowflake/commit/93d4d9cf5b0b25d85cd59b6cc9fd8d3a7af0401b))
* Fix BCR tests ([#3916](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3916)) ([53725bd](https://github.com/snowflakedb/terraform-provider-snowflake/commit/53725bde81a8e1912d0b7d854062f3869f4e627a))
* Fix secondary connection read ([#3914](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3914)) ([99eb947](https://github.com/snowflakedb/terraform-provider-snowflake/commit/99eb947271b03c05df24af595ff63b8d77ea3d35))
* Fix the default authenticator when the token is set ([#3899](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3899)) ([4122d73](https://github.com/snowflakedb/terraform-provider-snowflake/commit/4122d730b64be10cd4fa2935d619afabdb8b6abd))
* Fix using lowercased network policies in security integrations ([#3867](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3867)) ([941c918](https://github.com/snowflakedb/terraform-provider-snowflake/commit/941c9185c9a4b965fc8e75806a1450e393f7253d))
* Listing tests fix ([#3912](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3912)) ([d00f5f8](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d00f5f86cad3a0d829b0a9b49af3e9d1a08ff138))
* permadiff when granting privileges to shares on non-standard tables ([#3859](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3859)) ([d8d311b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d8d311be7b0027411ff37ceb8910df353069c902))
* schema default ddl collation handling ([#3907](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3907)) ([ea2c10d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ea2c10ddae79f0deaa22d6ce8c3c3e495416b947))

## [2.4.0](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v2.3.0...v2.4.0) (2025-07-24)


### 🎉 **What's new:**

* Add a user_programmatic_access_token resource ([#3835](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3835)) ([9bbeb4b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/9bbeb4bac067499f23955100e5b94d5183b4c139))
* Add a user_programmatic_access_tokens data source ([#3838](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3838)) ([de41e0c](https://github.com/snowflakedb/terraform-provider-snowflake/commit/de41e0c1c1c48198aa6551de0ee5a0d70253c49f))
* Add Current Organization Account resource and related assertions ([#3846](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3846)) ([e59f929](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e59f9297d2577c319650b4c80156a55146147e3a))
* Current organization account follow-up ([#3866](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3866)) ([4c11272](https://github.com/snowflakedb/terraform-provider-snowflake/commit/4c1127252ec718e6f08c1b6e82afac47af9d3f24))
* Implement Current Organization Account ([#3847](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3847)) ([4bbccea](https://github.com/snowflakedb/terraform-provider-snowflake/commit/4bbccea192b6a96d797daf88a341fd377ff65896))
* Implement rotating pats ([#3843](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3843)) ([c99f84d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/c99f84da0d56abae37a63767da19c579852c14fb))


### 🔧 **Misc**

* Add acceptance tests for rotating pats with an external resource ([#3871](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3871)) ([b8289c8](https://github.com/snowflakedb/terraform-provider-snowflake/commit/b8289c86b67aa44510ef3e091732aa8ea18c0374))
* Add an example of recreating the pipe on changing stage attributes ([#3729](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3729)) ([8fd4e7b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8fd4e7b7ce683caae5e2a0135ee8376af2e96586))
* Add basic listings to sdk ([#3800](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3800)) ([500a19a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/500a19ac46f4c363eef1131f8b8c55024faf70e1))
* Add builders for legacy config DTO ([#3842](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3842)) ([f3d4a9a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f3d4a9a2250c1ec539a217ef2d92c645d1aa2126))
* Add current organization account resource sdk part ([#3818](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3818)) ([56ea900](https://github.com/snowflakedb/terraform-provider-snowflake/commit/56ea9007ccc7c11c888865c3739be43eccd9e32a))
* Add PAT integration tests ([#3830](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3830)) ([6f688d6](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6f688d65cfc13cd10ef60bac3082a76dc3910a90))
* Address pre-push errors ([#3849](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3849)) ([39274dd](https://github.com/snowflakedb/terraform-provider-snowflake/commit/39274dd795fda0d7279967c735a196f46a57d753))
* Adjust sweepers and test workflows ([#3876](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3876)) ([3a4bd36](https://github.com/snowflakedb/terraform-provider-snowflake/commit/3a4bd366945faf7bfe129a0f2653aa4e67a93829))
* Fix tests and update documentation after reverting changes from BCR 2025_03 ([#3851](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3851)) ([04071aa](https://github.com/snowflakedb/terraform-provider-snowflake/commit/04071aac31d6c5328d6d6443946fa9ec855cc68a))
* Mark DecodeSnowflakeID function as legacy ([#3848](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3848)) ([d4a775f](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d4a775f4023f761e4c591e19d559f0f4e1a9c9f1))
* Test custom type with metadata ([#3856](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3856)) ([a168a24](https://github.com/snowflakedb/terraform-provider-snowflake/commit/a168a24a92c101acd85776b7b3d9ba52152f169a))
* Test enum handling (suppresion and validation) ([#3861](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3861)) ([0f051b8](https://github.com/snowflakedb/terraform-provider-snowflake/commit/0f051b82e855eefb9360233119fe87c112758b0b))
* Test optional computed in plugin framework ([#3865](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3865)) ([6674c12](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6674c125b70b4f22474bd3976fbc1a370d2bd1b6))
* Test optional with backing field ([#3857](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3857)) ([78e641a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/78e641af9b15a648b73b0e36bb2f405764405607))
* Test parameters handling (multiple variants) ([#3858](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3858)) ([5b88c06](https://github.com/snowflakedb/terraform-provider-snowflake/commit/5b88c06c4e2c00db1e4b24ab0e1e27951760f39e))
* Test zero values in plugin framework ([#3850](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3850)) ([d7816fa](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d7816fa7df164a0f6140ccb105e019a9e85e5bca))
* Update labels in the repository ([#3863](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3863)) ([92669bb](https://github.com/snowflakedb/terraform-provider-snowflake/commit/92669bb7df91a9f6f0179e0c3c001befb9605e82))
* Use existing schema descriptions instead of string equivalents ([#3868](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3868)) ([e9d0448](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e9d044867d716cb412add8f1abbd6cafe07508ea))

## [2.3.0](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v2.2.0...v2.3.0) (2025-07-03)


### 🎉 **What's new:**

* Add programmatic access token support to SDK ([#3819](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3819)) ([6b8ae55](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6b8ae5504c4a3e5c768e818f02d927ad06fa338e))
* Add support for PROGRAMMATIC_ACCESS_TOKEN authenticator ([#3805](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3805)) ([2a4379b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/2a4379baca95a08f7d9ee291fd9a7dca4ad17286))


### 🔧 **Misc**

* Account modificaiton test assertion ([#3759](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3759)) ([6151757](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6151757bca49f02836f57c506e41404aca95c345))
* Add Snowflake BCR migration guide ([#3829](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3829)) ([305853d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/305853d4929f7051eb7a1e63be79ca8302a8efc5))
* Configure plugin framework in functional tests ([#3824](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3824)) ([0cfced0](https://github.com/snowflakedb/terraform-provider-snowflake/commit/0cfced0ce1bf8d82b66a2dcbfe725e6795f3a8ab))
* Do not build the whole project after the changelog entry ([#3817](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3817)) ([d40cc45](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d40cc450b0d177ed672b813e715f9f938e0065d3))
* Enable testifylint and fix reported issues ([#3793](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3793)) ([605f0cf](https://github.com/snowflakedb/terraform-provider-snowflake/commit/605f0cfdc3475170e222275ae47edb9a1b5e512e))
* Set up muxing in tests ([#3804](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3804)) ([ab73524](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ab73524df4e2d84634352b504b1119576e351673))
* Small account adjustments ([#3786](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3786)) ([b561ce5](https://github.com/snowflakedb/terraform-provider-snowflake/commit/b561ce53d577c821ccc7def0ccbf4f410f6ef076))


### 🐛 **Bug fixes:**

* Fix data types parsing for functions and procedures with 2025_03 Bundle ([#3827](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3827)) ([35bdf1c](https://github.com/snowflakedb/terraform-provider-snowflake/commit/35bdf1c1bf029abb9b54347ff5971bb6d86e4bb1))
* Introduce a new function and procedure parsing function ([#3825](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3825)) ([75faf30](https://github.com/snowflakedb/terraform-provider-snowflake/commit/75faf302e724e0cd1b5f041fef209e179cf920b2))
* Remove unused conversion functions interfering with other tests ([#3820](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3820)) ([878a4e7](https://github.com/snowflakedb/terraform-provider-snowflake/commit/878a4e7006e04220dbe62846dce3423adf275c99))

## [2.2.0](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v2.1.0...v2.2.0) (2025-06-25)


### 🎉 **What's new:**

* Add Compute Pool resource ([#3689](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3689)) ([c6f0429](https://github.com/snowflakedb/terraform-provider-snowflake/commit/c6f04298cb43eeb7e1c75ed3fb518ae32e533f0d))
* Add Compute Pools data source ([#3707](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3707)) ([44a8639](https://github.com/snowflakedb/terraform-provider-snowflake/commit/44a8639f1338e5c781ebdad23bdd729ab336ba4e))
* Add compute pools integration tests ([#3680](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3680)) ([c249f55](https://github.com/snowflakedb/terraform-provider-snowflake/commit/c249f557cbafe6f03215a8181aa70eef6c87d1fb))
* Add compute pools to SDK ([#3678](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3678)) ([2baf141](https://github.com/snowflakedb/terraform-provider-snowflake/commit/2baf141c2018f2bf464a39131bd5a6d15ada42b7))
* Add current account resource ([#3745](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3745)) ([9914a9b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/9914a9b7dd12d4db0f958dcd0a8befb7c2a0ea87))
* Add embedding model support to Cortex Search Service ([#3711](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3711)) ([169b942](https://github.com/snowflakedb/terraform-provider-snowflake/commit/169b94298e0df9f1f2c3aef4e7e3925e63af9bf1))
* Add EXECUTE JOB support to SDK ([#3738](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3738)) ([92c4fda](https://github.com/snowflakedb/terraform-provider-snowflake/commit/92c4fda387267d64e7b212169e1ef77722fea980))
* Add git repository datasource ([#3746](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3746)) ([73835ff](https://github.com/snowflakedb/terraform-provider-snowflake/commit/73835ffc542206a16f99c9dbe16f5831c77675da))
* Add git repository definition ([#3705](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3705)) ([cdd600d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/cdd600d4568584ca9e7cb9a4536587c992b98a8a))
* Add git repository resource ([#3739](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3739)) ([a3d333d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/a3d333d12f8932fd21f3693d4e7472003198843e))
* Add Image Repository data source ([#3661](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3661)) ([957691f](https://github.com/snowflakedb/terraform-provider-snowflake/commit/957691f2b336f194e34daa42a54043e626737c2c))
* Add Image Repository resource ([#3660](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3660)) ([adda1c9](https://github.com/snowflakedb/terraform-provider-snowflake/commit/adda1c98320e3422264e4e77be5d88728f69b147))
* Add job service resource ([#3747](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3747)) ([3efbbc2](https://github.com/snowflakedb/terraform-provider-snowflake/commit/3efbbc25ab426dd439be6b20f425ffce0ffe7e1e))
* Add managing tags and integration tests for git repositories ([#3723](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3723)) ([2f1172d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/2f1172d05521f2d27f0dd07f03cf848ea13942e8))
* Add service resource ([#3732](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3732)) ([c3bfef9](https://github.com/snowflakedb/terraform-provider-snowflake/commit/c3bfef967ab039ec1a76c263b91e38b33ad1b5f4))
* Add services data source ([#3752](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3752)) ([9f7f5ac](https://github.com/snowflakedb/terraform-provider-snowflake/commit/9f7f5ac1c52fe7c36b4b4186e358c8f4361c8b3d))
* Add services to SDK part 1 ([#3706](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3706)) ([9e1e31b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/9e1e31b8524f7907f8e93eaa33321511aaf8c718))
* Add tag association support for services ([#3726](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3726)) ([30cf1b8](https://github.com/snowflakedb/terraform-provider-snowflake/commit/30cf1b8a27080929b486e6adca66685b9213c5b1))
* Add templates support to services and job services ([#3763](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3763)) ([7558b90](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7558b900286a6b546448e9b0131229b14e88916a))
* Allow granting privileges on future cortex search services ([#3760](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3760)) ([e96d097](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e96d097930f319002287f3eba77763468e8b8c37))
* get datasource tables on par with views ([#3249](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3249)) ([d5b8dca](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d5b8dca48490900589c1902258a476b4927e16d5))
* Git repository code refactor ([#3753](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3753)) ([daeb8cf](https://github.com/snowflakedb/terraform-provider-snowflake/commit/daeb8cf7cfd3a4a0cd9fe646143dc5a597fe8ab8))
* Recreate service objects on externally changed type ([#3748](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3748)) ([8597e35](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8597e3593d0876bd96840ef46cffb3c4b53c527a))
* Support managing tags for image repositories ([#3662](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3662)) ([8c4aed2](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8c4aed2aad527fdcc0334abf0fb9029b251999ad))


### 🔧 **Misc**

* Account resource parameters ([#3736](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3736)) ([1e3ff03](https://github.com/snowflakedb/terraform-provider-snowflake/commit/1e3ff034870ca0ddfef3dc8d24d9c87af78adcf7))
* Add auto label assignment workflow ([#3712](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3712)) ([20394d5](https://github.com/snowflakedb/terraform-provider-snowflake/commit/20394d56560434827d4e14498522d061b6f4077b))
* Add current account resource policies ([#3768](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3768)) ([92cfa5a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/92cfa5a1e0cc59a13daad59e840351cbcec46437))
* Add location identifier ([#3721](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3721)) ([39c79db](https://github.com/snowflakedb/terraform-provider-snowflake/commit/39c79dba78f201e32b25b90670ef9b86d8339901))
* Add new roadmap entry ([#3702](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3702)) ([5bdcd12](https://github.com/snowflakedb/terraform-provider-snowflake/commit/5bdcd127d9288212b10ea7b138bebc0cb770c5b9))
* Add PAT authentication method to the documentation ([#3749](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3749)) ([44cb288](https://github.com/snowflakedb/terraform-provider-snowflake/commit/44cb2887e43087da856467ec67e3fd2987d2cd81))
* Add service integration tests ([#3725](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3725)) ([8cd7e78](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8cd7e78cc52a5f2114cad8c0c6cd3dc11eacfb63))
* Add tag propagation limitation to the tag resource documentation ([#3770](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3770)) ([4f38421](https://github.com/snowflakedb/terraform-provider-snowflake/commit/4f38421d4cbc02e3f4d87e6bf6e71553b4455450))
* Add TOML config generator ([#3619](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3619)) ([7685d5a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7685d5ada21cce838ba6205a6a84156e010ddcc0))
* Add variable model builder ([#3663](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3663)) ([112587b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/112587bda4fff948c187643d38057530b2099843))
* Adjust resource assertions ([#3715](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3715)) ([ab74816](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ab7481609246c5ccc95705c6f07f788dce13c30b))
* Adjust resource assertions follow up ([#3716](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3716)) ([b1e4faf](https://github.com/snowflakedb/terraform-provider-snowflake/commit/b1e4fafbb5b49308d64c0e102aa03c6c764e4690))
* Adjust trace level possible values ([#3737](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3737)) ([7d5447b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7d5447b0cb9d4220726755a39aa27c50d367134d))
* Apply changes from previous releases ([#3802](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3802)) ([d5a40df](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d5a40dfb46fe9c709b3ce3703d77f6c237c76a45))
* Clarify SDK generator future ([#3771](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3771)) ([6cfa554](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6cfa554b7415edcdcced4f0430ae186f9e8cb31f)), closes [#3671](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3671)
* Clean old acceptance test methods ([#3727](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3727)) ([f94e773](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f94e7734a8a699377c5cc7d0799ffdeb26c263e7))
* enhance privilege-granting resource documentation ([#3700](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3700)) ([0168e70](https://github.com/snowflakedb/terraform-provider-snowflake/commit/0168e70a194baf21e4e47e7f7d55eac1e3acdc3e))
* Extract a separate directory for the acceptance tests ([#3679](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3679)) ([6b75033](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6b750336a6c15c9b9854943e9e038d92fdb759c0))
* Fix tests that run on secondary account profile ([#3666](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3666)) ([976b36c](https://github.com/snowflakedb/terraform-provider-snowflake/commit/976b36c5702b5566eef7a1db04142f7920b7c188))
* Implement double dollar quotes in SDK builders ([#3720](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3720)) ([e2fa6c5](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e2fa6c57f6de57c612df2657890fa63f00aa02d2))
* Introduce account parameters ([#3734](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3734)) ([4b6ec2f](https://github.com/snowflakedb/terraform-provider-snowflake/commit/4b6ec2f871d4a160fdffdc83715d32161b2b4901))
* Regenerate resource assertions ([#3713](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3713)) ([53cb291](https://github.com/snowflakedb/terraform-provider-snowflake/commit/53cb2918f950a8e0754ed0dec6c4d447d1f87cf0))
* Rename acceptance tests files ([#3701](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3701)) ([4ae3eac](https://github.com/snowflakedb/terraform-provider-snowflake/commit/4ae3eac836c472ad4dde3b49a18f023aca62e596))
* Rename account parameters field name ([#3708](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3708)) ([bb83144](https://github.com/snowflakedb/terraform-provider-snowflake/commit/bb83144cd0a5cae64fdf97861fab2be12b61d650))
* Rename account unset parameters ([#3728](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3728)) ([57176e3](https://github.com/snowflakedb/terraform-provider-snowflake/commit/57176e36280a084359a2cc45a25189699f0f689e))
* Show output assertion adjustments ([#3718](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3718)) ([f25db1d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f25db1d169745a335c6bb57bffc77da2087ab5a9))
* Show output assertion adjustments follow up ([#3719](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3719)) ([842526b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/842526b9c85ecfc70d6bd2e7b52a5f7d804dd21a))
* Simplify constructor and builder method generation ([#3640](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3640)) ([697bacf](https://github.com/snowflakedb/terraform-provider-snowflake/commit/697bacfa0406f67d7f75382768ab46aa265eaf7c))
* Small improvements ([#3730](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3730)) ([5abb2f9](https://github.com/snowflakedb/terraform-provider-snowflake/commit/5abb2f960279577b29683693d2fe3fe251e7400f))
* Stabilize object assertions generation ([#3722](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3722)) ([cf907df](https://github.com/snowflakedb/terraform-provider-snowflake/commit/cf907df794f7569bc1a7f03365a8433248526a65))
* Test fixes ([#3740](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3740)) ([2d0c915](https://github.com/snowflakedb/terraform-provider-snowflake/commit/2d0c915ac520ee2a82dfaf0341510e5d29b2c5a7))
* Test provider config tri-value booleans ([#3765](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3765)) ([33b50c2](https://github.com/snowflakedb/terraform-provider-snowflake/commit/33b50c2c3791eb634790cf6cfc011251e72502f7))
* Update docs ([#3758](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3758)) ([f3e43a0](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f3e43a025f97e168199d0d387a88ca53446bb981))
* Upgrade gosnowflake to v1.14.1 ([#3792](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3792)) ([a239157](https://github.com/snowflakedb/terraform-provider-snowflake/commit/a239157df0fdbdbab42f90c975a3505b0ef1e42b))
* Use common setup and cleanup for acceptance tests - data sources ([#3691](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3691)) ([2226369](https://github.com/snowflakedb/terraform-provider-snowflake/commit/2226369ec746e13d2b4f0b7195b30589306028b4))
* Use common setup and cleanup for acceptance tests - provider ([#3724](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3724)) ([403368e](https://github.com/snowflakedb/terraform-provider-snowflake/commit/403368e6d426efbba2bd52c696a5abda02d64230))
* Use common setup and cleanup for acceptance tests - resources ([#3692](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3692)) ([5e2fd24](https://github.com/snowflakedb/terraform-provider-snowflake/commit/5e2fd24c01401cda3334202faa1ef3118dc9465e))


### 🐛 **Bug fixes:**

* Adjust function and procedure tests ([#3693](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3693)) ([702d6df](https://github.com/snowflakedb/terraform-provider-snowflake/commit/702d6dfad0cb56f9a916aed105d324f475beb76f))
* Correct error handling in CreateContextExternalVolume function ([#3636](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3636)) ([2be1b43](https://github.com/snowflakedb/terraform-provider-snowflake/commit/2be1b433a9a21e3a84bc714f3b7c8e5727a0804f))
* Empty privileges ([#3695](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3695)) ([ae1a621](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ae1a62165156823fbb85159092fbbc7f92c5d003)), closes [#3690](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3690)
* Fix account parameter mapping ([#3799](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3799)) ([882173c](https://github.com/snowflakedb/terraform-provider-snowflake/commit/882173c5eb40c671d68069bd27805462afa2240b))
* Fix external oauth integration issuer tests ([#3769](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3769)) ([74e23cf](https://github.com/snowflakedb/terraform-provider-snowflake/commit/74e23cf76cafe88c61e578122c1de2dbf6545f0c))
* Fix external oauth test ([#3773](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3773)) ([8cbd052](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8cbd0526feb550fe44195aec4287ad9a2e81515e))
* Fix image repositories cleanup ([#3694](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3694)) ([e602615](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e602615a9a979f1656e2e41ddb8c1e8c952438e3))
* Fix parsing username in handling grants ([#3688](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3688)) ([d52b98b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d52b98bb962fc59d10914789035c8baa3f479bf3))
* Grant ownership on serverless tasks ([#3772](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3772)) ([0a457bc](https://github.com/snowflakedb/terraform-provider-snowflake/commit/0a457bc0dfc144da6c98b1da0114e4db49d7483b))
* Table test ([#3696](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3696)) ([420a673](https://github.com/snowflakedb/terraform-provider-snowflake/commit/420a673284956b697a942a5ccc4aebea49cd91b5))
* Update streamlit docs and fix view rename with modification ([#3766](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3766)) ([0e89dfd](https://github.com/snowflakedb/terraform-provider-snowflake/commit/0e89dfd317bdfdc3e1374def962776e76cf6c7a8)), closes [#3676](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3676)
* user authentication policy attachment ([#3697](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3697)) ([5b2412f](https://github.com/snowflakedb/terraform-provider-snowflake/commit/5b2412f915d22e8578419ca55817b36815c591ad)), closes [#3672](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3672)

## [2.1.1](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v2.1.0...v2.1.1) (2025-06-24)


### 🔧 **Misc**

* prepare v2.1.1 release ([f6561a4](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f6561a4c1ed019c47b0168d65fc74158ac74e6ee))


### 🐛 **Bug fixes:**

* account parameter mapping ([e52be9b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e52be9b857408f4c3918a9b32c54e13881ed5cf5))

## [2.1.0](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v2.0.0...v2.1.0) (2025-05-12)


### 🎉 **What's new:**

* Add image repositories to SDK ([#3628](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3628)) ([ffdb1ee](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ffdb1ee7b6add7d4a3bef7d03623b3dd4149ffc7))


### 🔧 **Misc**

* Add a holiday notice ([#3635](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3635)) ([f202704](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f202704766ccf68fa214a60a6aa1566b96671ab3))
* Change date format DD.MM.YYYY to YYYY-MM-DD ([#3621](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3621)) ([d0fa668](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d0fa6688bc0cf52a12cf2587b9f2c3574ea3821d))
* Document deprecated grant resource mapping ([#3643](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3643)) ([7090b9f](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7090b9f35b4fb31ea53f52a9f382fb9e09807ca5))
* Ensure no changes to generated model builders in pre push ([#3649](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3649)) ([25a9e17](https://github.com/snowflakedb/terraform-provider-snowflake/commit/25a9e17f7131f5fd1c29e14164b7970d6638fe3e))
* Fix ordering in all model builders ([#3618](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3618)) ([d5218ef](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d5218efed65ad2e0be7d1c6a9ea2188a609dc36f))
* Fix ordering in resource model builder constructors ([#3615](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3615)) ([5724340](https://github.com/snowflakedb/terraform-provider-snowflake/commit/5724340e7dba887c3e05ab42b64ed617ab5e5193))
* Fix tests that will fail once 02_2025 is enabled ([#3650](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3650)) ([b39c8fd](https://github.com/snowflakedb/terraform-provider-snowflake/commit/b39c8fd4374f4c7fefb41752cab819a1d06671a9))
* Generate stable and preview resources in docs ([#3616](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3616)) ([3005841](https://github.com/snowflakedb/terraform-provider-snowflake/commit/30058419f3c7b376baf48c1aa74eac4d4c8fcc0f))
* Handle forbidden attribute names in resource model builder generation ([#3613](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3613)) ([87caec3](https://github.com/snowflakedb/terraform-provider-snowflake/commit/87caec3575f2cd9ad5d25a8d724298256e788c7c))
* Handle slices in resource model builder generation ([#3614](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3614)) ([b480fbf](https://github.com/snowflakedb/terraform-provider-snowflake/commit/b480fbfa04af484c405fb098f3f5955ef8d8c3ad))
* Injecting provider version ([#3612](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3612)) ([d14d8ab](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d14d8ab5f162e2f51ee00a117a234d0e0a02d583))
* Manually imported gpg key ([#3653](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3653)) ([f732ec5](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f732ec5b4036dd124c83f86e54e21c70d7a4ab67))
* Masking random values ([#3583](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3583)) ([3216df2](https://github.com/snowflakedb/terraform-provider-snowflake/commit/3216df2f3cf0236a765d038bcfda355efac90022))
* Minor fixes after v2 release ([#3606](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3606)) ([99a28b4](https://github.com/snowflakedb/terraform-provider-snowflake/commit/99a28b45e76c4393dbca00573b988a6442f55425))
* Minor fixes and documentation adjustment ([#3607](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3607)) ([1daabaf](https://github.com/snowflakedb/terraform-provider-snowflake/commit/1daabaf693aa5138693f11b385e665dae0b64137))
* remove note regarding future Streamlits ([#3611](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3611)) ([8de216c](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8de216c55e8ef8b15faf65fd549fe092f5f01e82)), closes [#3423](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3423) [#3569](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3569)
* Remove the holiday note ([#3644](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3644)) ([11b3848](https://github.com/snowflakedb/terraform-provider-snowflake/commit/11b384852862d75cd5a4d68520a44b058eaafe18))
* Update migration guide ([#3665](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3665)) ([c1426f2](https://github.com/snowflakedb/terraform-provider-snowflake/commit/c1426f22c849c92a1fc186a1e6d5e17e49cc7e4c))
* Upgrade GoSnowflake driver ([#3634](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3634)) ([aaebce5](https://github.com/snowflakedb/terraform-provider-snowflake/commit/aaebce52f07677235f3c1e570e076d0d4095c16d))
* Use proper envs for providing config files in pipelines ([#3625](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3625)) ([26a9a06](https://github.com/snowflakedb/terraform-provider-snowflake/commit/26a9a062bfb878fb424cbcb945203caaec773c27))
* Use the new TOML format in pipelines ([#3608](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3608)) ([4f56f31](https://github.com/snowflakedb/terraform-provider-snowflake/commit/4f56f316cdf438488c47cb353aa9b25d5ca1235e))


### 🐛 **Bug fixes:**

* 3622 issue ([#3626](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3626)) ([4b2f8c6](https://github.com/snowflakedb/terraform-provider-snowflake/commit/4b2f8c6a31771d273f5acb56cdfefe8631d45035))
* Add missing DISABLE_USER_PRIVILEGE_GRANTS parameter ([#3652](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3652)) ([eb7ef76](https://github.com/snowflakedb/terraform-provider-snowflake/commit/eb7ef76f64834dacdfa94a9bd77d0f4f2888b62b)), closes [#3639](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3639)
* Adjust tag association model builder ([#3642](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3642)) ([e80f4c9](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e80f4c9c5ed1ab2f151ad39022c1e1ba66c00288))
* Bring back account retry functionality ([#3645](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3645)) ([8a63774](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8a6377444964b7384052327ff07245b2d8427570))
* Generate dynamic blocks for all resource model builders ([#3648](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3648)) ([6bb868a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6bb868aa8e735e7f5442e3ae57569ab6e70be3b7))
* issue 3629 ([#3633](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3633)) ([010946e](https://github.com/snowflakedb/terraform-provider-snowflake/commit/010946e807854dff8e6e205fb4f7a69f57076967))
* object tracking ([#3627](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3627)) ([e73cd77](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e73cd772599d23a7424fe309ab9f68cadf50affd))
* Prove import diff suppression problem in procedure resources ([#3656](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3656)) ([8e5f55e](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8e5f55ec959b80ed240f3dd76ff4a3d8e499cdf3)), closes [#3401](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3401)
* remove tautological err == nil check in network policy attachment test ([#3631](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3631)) ([24de0b7](https://github.com/snowflakedb/terraform-provider-snowflake/commit/24de0b7a1657f98517d6f401fc48876479b68dd3))
* Set value into correct field for field transformers ([#3646](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3646)) ([1da734f](https://github.com/snowflakedb/terraform-provider-snowflake/commit/1da734f7c298a85c14c4e1b1ec889e83817d3ccf))
* Suppress diff for network policy in user resources ([#3657](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3657)) ([b4900c6](https://github.com/snowflakedb/terraform-provider-snowflake/commit/b4900c6665e86208564762b4a5adc1765ad202d5)), closes [#3655](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3655)

## [2.0.1](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v2.0.0...v2.0.1) (2025-06-24)


### 🔧 **Misc**

* prepare v2.0.1 release ([d17bcec](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d17bcec7ed9f79b8345c9ca0dff50f8a6e0c6ee0))


### 🐛 **Bug fixes:**

* Fix account parameter mapping ([#3797](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3797)) ([837bf9e](https://github.com/snowflakedb/terraform-provider-snowflake/commit/837bf9ecf01b8320855695db1c9a0997c8e06267))

## [2.0.0](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v1.2.1...v2.0.0) (2025-04-23)


### 🎉 **What's new:**

* Diff suppression handling follow up part 2 ([#3592](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3592)) ([39cf134](https://github.com/snowflakedb/terraform-provider-snowflake/commit/39cf1348582a44d3333c9188660a098ddef930f2))
* Use the new TOML format by default ([#3585](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3585)) ([4fd8a0d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/4fd8a0d921a29a860494cc5f751b014aa1e0a82b))
* Verify file permissions by default ([#3519](https://github.com/snowflakedb/terraform-provider-snowflake/pull/3519)) ([551dbfa](https://github.com/snowflakedb/terraform-provider-snowflake/commit/551dbfab36c95d528b7a1721953dfcc32d70810e))


### 🔧 **Misc**

* Add missing migration guide ([#3597](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3597)) ([5af6aec](https://github.com/snowflakedb/terraform-provider-snowflake/commit/5af6aece96e8bf2b8cdc32bcf4637bc59419ae3c))
* Add new roadmap entry ([#3598](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3598)) ([e2a56e9](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e2a56e92c843a14d8e785d58aaf0003446fb4007))
* Documentation adjustments ([#3599](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3599)) ([94417fe](https://github.com/snowflakedb/terraform-provider-snowflake/commit/94417feb24766194eab26fa0b3ab3731da409c91))
* GH documentation updates ([#3596](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3596)) ([97bbd7e](https://github.com/snowflakedb/terraform-provider-snowflake/commit/97bbd7ed42d070524a89b9e5f0a083fad76a1911))
* Prepare 2.0.0 release ([#3602](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3602)) ([91c7673](https://github.com/snowflakedb/terraform-provider-snowflake/commit/91c767382fcc06ecc64b4131feb70008566bc032))
* Sensitive values and documentation changes ([#3595](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3595)) ([70a4dc5](https://github.com/snowflakedb/terraform-provider-snowflake/commit/70a4dc5d19763242e0dbcb3605c518250c6b0e05))
* Adjust sensitive values ([#3572](https://github.com/snowflakedb/terraform-provider-snowflake/pull/3572)) ([abef31c](https://github.com/snowflakedb/terraform-provider-snowflake/commit/abef31c4fbed205b171a1d59a60948aefc02d42a))


### 🐛 **Bug fixes:**

* Add single quotes to CSV_TIMESTAMP_FORMAT parameter value ([#3582](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3582)) ([bad52e3](https://github.com/snowflakedb/terraform-provider-snowflake/commit/bad52e378b1c823240756c9c6444f55842d9189b))

## [1.2.2](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v1.2.1...v1.2.2) (2025-06-24)


### 🔧 **Misc**

* prepare v1.2.2 release ([f8445d8](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f8445d89799f856e839b3a7c0fb8035fe71bbf6b))


### 🐛 **Bug fixes:**

* Fix account parameter mapping ([#3796](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3796)) ([44ed482](https://github.com/snowflakedb/terraform-provider-snowflake/commit/44ed482ed10afe577fa07668ca8921cdecff57fb))

## [1.2.1](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v1.2.0...v1.2.1) (2025-04-22)


### 🔧 **Misc**

* Bump versioning ([#3589](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3589)) ([3fa737a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/3fa737a45ca1d58d9d984b0703f03da5cab2e3d1))

## [1.2.0](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v1.1.0...v1.2.0) (2025-04-22)


### 🎉 **What's new:**

* Add a new TOML file format ([#3577](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3577)) ([e92de36](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e92de36244b6acfced4b59f3cef1c897aaa9a263))
* Add streamlits as valid future grant target ([#3569](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3569)) ([7c870aa](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7c870aa3e10343d2af2b10ec8cc21cdd9c4ed227))
* Handle operations in resources using safe functions ([#3578](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3578)) ([ba39d91](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ba39d9174514206b1eb7bf36cbf8a123fd824a5f))


### 🔧 **Misc**

* Adjust codeql analysis ([#3563](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3563)) ([c881ab3](https://github.com/snowflakedb/terraform-provider-snowflake/commit/c881ab39dc8af5a5212fcb53ed8b99a61cf43ab9))
* Apply safe functions ([#3576](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3576)) ([ae45777](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ae457770179480ba2cc22e14500b96c0c33a4d0c))
* Masking random values ([#3545](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3545)) ([244e12b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/244e12bdf94dd88b0422969a13062104f87d8472))
* Rename ConfigDTO to LegacyConfigDTO ([#3573](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3573)) ([515d685](https://github.com/snowflakedb/terraform-provider-snowflake/commit/515d685800fabcc243499736028169fc037dbea2))
* Test proper diff suppression for data types ([#3579](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3579)) ([059d205](https://github.com/snowflakedb/terraform-provider-snowflake/commit/059d20591dac8f8a999b03db445d46ffaca9c1de))
* Update repo_meta.yaml ([#3566](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3566)) ([6629dbd](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6629dbdf8deb9221215d29aca0494e6a2df5d0c8))
* Use opts for reading TOML config ([#3575](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3575)) ([9fe0d14](https://github.com/snowflakedb/terraform-provider-snowflake/commit/9fe0d144adae85bdfa3d5f2da1c385e0f8bc7f34))


### 🐛 **Bug fixes:**

* Fix handling references to computed fields in context of `show_output` ([#3564](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3564)) ([d464752](https://github.com/snowflakedb/terraform-provider-snowflake/commit/d46475226a1f13bbbaed2b8e5fdc3b74ddd76c24)), closes [#3522](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3522)
* Fix identifier tests ([#3561](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3561)) ([fd295a5](https://github.com/snowflakedb/terraform-provider-snowflake/commit/fd295a56b7f9c6a5d719858db062f7fbe9bcfde0))
* Handle drop operation safely ([#3570](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3570)) ([7bd15fd](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7bd15fd1ad13048e1d075310ad52356ac3affaad))
* Split acceptance tests ([#3562](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3562)) ([73ad5d9](https://github.com/snowflakedb/terraform-provider-snowflake/commit/73ad5d9f8563993730899ecf4c6d0aff324da445))

## [1.1.1](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v1.1.0...v1.1.1) (2025-06-24)

### 🔧 **Misc**

* prepare v1.1.1 release ([96454d7](https://github.com/snowflakedb/terraform-provider-snowflake/commit/96454d7b1e6e1d3c4d162d8487d85b9a20d0497f))
* update gosnowflake driver version ([0daa7b9](https://github.com/snowflakedb/terraform-provider-snowflake/commit/0daa7b9cfa087450dae4af28672c67f689ab68d4))
* update gosnowflake driver version ([a0289d4](https://github.com/snowflakedb/terraform-provider-snowflake/commit/a0289d48a53db2b4bbd3d1947fe3df61f6f11402))

### 🐛 **Bug fixes:**

* Fix account parameter mapping ([#3795](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3795)) ([e953bad](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e953bad901bdd97985a0bacb43173103743b62c6))

## [1.1.0](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v1.0.5...v1.1.0) (2025-04-10)


### 🎉 **What's new:**

* Add timeouts to resources ([#3511](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3511)) ([4bd364d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/4bd364d0d116404ffbe8041b9d2e7b367b622d4f))


### 🔧 **Misc**

* Add a note about default values before every schema markdown  ([#3483](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3483)) ([730590d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/730590d773f87ee49b6ea9cc46ccbe0ba2df7b7e))
* Add auto-generated sensitive value mask ([#3527](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3527)) ([6ca2629](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6ca2629786c18c03e9cc7a0596066da1285f139e))
* Add information about migrating to snowflakedb namespace ([#3513](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3513)) ([8a8b446](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8a8b44669e7bd6a079355f65c84a73e89483f350))
* Add missing default value descriptions ([#3482](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3482)) ([18c219b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/18c219b4ec0a069dc0448f31376ec720e2cae7f2))
* Add PoC for dynamic blocks in model builders ([#3492](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3492)) ([7a5ee0f](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7a5ee0f1861355ba4cb87eb57096ea4e6d364a7b))
* Add repo_meta.yaml files ([#3484](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3484)) ([5581e88](https://github.com/snowflakedb/terraform-provider-snowflake/commit/5581e889dc9f0d67fa2aee8276fbdefb6803cadd))
* Add Snyk Scanning ([#3520](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3520)) ([c5d1df0](https://github.com/snowflakedb/terraform-provider-snowflake/commit/c5d1df07be1e4043847e6f5c7770370831df3f0c))
* Adjust issue templates ([#3486](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3486)) ([53727a0](https://github.com/snowflakedb/terraform-provider-snowflake/commit/53727a0a4372cd024deca68c08b4570b7e6be7fe))
* Adjust links and namespaces part 1 ([#3529](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3529)) ([a1538e0](https://github.com/snowflakedb/terraform-provider-snowflake/commit/a1538e051b7c2d55eb0cec9227dc8423c6bf8218))
* Adjust logs ([#3523](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3523)) ([37e70ac](https://github.com/snowflakedb/terraform-provider-snowflake/commit/37e70acd1b855927fceb13a1da7b3ebb05203b82))
* Adjust ok-to-test ([#3502](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3502)) ([bc12db8](https://github.com/snowflakedb/terraform-provider-snowflake/commit/bc12db81e6ff0ff2faa1190b18b631874e03a328))
* Adjust the repository documentation. ([#3454](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3454)) ([1e8c59a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/1e8c59a4c4e29f1feb897ff17cf14b42bec665ce))
* Allow only one concurrent job for tests on the given branch ([#3533](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3533)) ([5839c93](https://github.com/snowflakedb/terraform-provider-snowflake/commit/5839c935086ffd6a1ec452033b21373d465f1b95))
* Bump provider version ([#3537](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3537)) ([757eb60](https://github.com/snowflakedb/terraform-provider-snowflake/commit/757eb60eb57fcfce9a6563258efa3dc30a92fc7d))
* Change privileges in ok-to-test ([#3500](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3500)) ([dcc884a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/dcc884a26682113a0d49392e46d2b6bd8b4608d3))
* Fix changelog after v1.0.5 ([#3485](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3485)) ([e80804b](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e80804b5c64e3f3d400a59786dbbdfc910e6be6f))
* Fix paths ignore ([#3541](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3541)) ([b8a2158](https://github.com/snowflakedb/terraform-provider-snowflake/commit/b8a21587de4694b818f83998093f5275ec4ea5c1))
* Handle read operation safely ([#3491](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3491)) ([baaab78](https://github.com/snowflakedb/terraform-provider-snowflake/commit/baaab78aab2ee3d15035b82533d0b44518ec92b9))
* Improve error messages ([#3503](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3503)) ([e2ed124](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e2ed124f494ec33e447e9e0fd538a5acc14fa9cb))
* New roadmap entry GA announcement ([#3534](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3534)) ([f21a62d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f21a62d86357667fe68049c45c6b1495ebaf8799))
* Public key downloader ([#3499](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3499)) ([ef51a13](https://github.com/snowflakedb/terraform-provider-snowflake/commit/ef51a13d35dfddf4bceabfef6c22ca1ce058a6a0))
* Release workflow adjustments ([#3531](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3531)) ([2309228](https://github.com/snowflakedb/terraform-provider-snowflake/commit/2309228744af0bb98c506473f6678fcdd304e0b3))
* Replace registry namespace ([#3536](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3536)) ([dd45bef](https://github.com/snowflakedb/terraform-provider-snowflake/commit/dd45bef35d7cf72fd7e2d295c794d56699618db8))
* Set a fixed golangci-lint version in GH actions ([#3470](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3470)) ([e8e4b2d](https://github.com/snowflakedb/terraform-provider-snowflake/commit/e8e4b2d3592cc68e88507ba8f897e45f63bc0e34))
* Set parameters for database precreated for integration tests ([#3524](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3524)) ([3c1f46e](https://github.com/snowflakedb/terraform-provider-snowflake/commit/3c1f46e079dac2636f62ec6731063476aded9498))
* Skip tests workflow for docs only ([#3535](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3535)) ([fac40f8](https://github.com/snowflakedb/terraform-provider-snowflake/commit/fac40f8dfb5260ed56153884dbe9e8c979bfab8a))
* Update codeowners ([#3481](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3481)) ([cffba96](https://github.com/snowflakedb/terraform-provider-snowflake/commit/cffba969cb58fd170e146f701a28f444506ae63f))
* Update codeql scanner ([#3517](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3517)) ([629f307](https://github.com/snowflakedb/terraform-provider-snowflake/commit/629f307d1631e0fd8be34346bd6d2c50c402de33))
* Update documentation to include default values ([#3452](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3452)) ([72eaab6](https://github.com/snowflakedb/terraform-provider-snowflake/commit/72eaab6cfcf337a0ce2585749614f4140b3f2795))
* Update labels ([#3457](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3457)) ([7dbe180](https://github.com/snowflakedb/terraform-provider-snowflake/commit/7dbe180e36221ce541d9260db1b4936b4ccaf892))
* Upgrade dependencies ([#3521](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3521)) ([22b4091](https://github.com/snowflakedb/terraform-provider-snowflake/commit/22b409177962e15491b2795aab08d06742aba520))
* Upgrade deps ([#3487](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3487)) ([98737bc](https://github.com/snowflakedb/terraform-provider-snowflake/commit/98737bcdffc1a0b33b5f804290891e66763f31de))


### 🐛 **Bug fixes:**

* Change setup for precreated objects in acceptance tests ([#3515](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3515)) ([0b76012](https://github.com/snowflakedb/terraform-provider-snowflake/commit/0b76012a0c41de74a67a944bbb841c3510c885f6))
* Fix acceptance tests setup ([#3516](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3516)) ([b5f5676](https://github.com/snowflakedb/terraform-provider-snowflake/commit/b5f56768345347b4c64907c6e474443428382319))
* Fix dispatch event checkout ([#3407](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3407)) ([b0be814](https://github.com/snowflakedb/terraform-provider-snowflake/commit/b0be814323043bcc0deac774883304dbec8bc5f0))
* Fix resource acceptance tests setup part1 ([#3450](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3450)) ([3ff4ba5](https://github.com/snowflakedb/terraform-provider-snowflake/commit/3ff4ba5ee9715785403f1fd89bb4c826af514c9f))
* Fix resource acceptance tests setup part2 ([#3480](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3480)) ([f61051a](https://github.com/snowflakedb/terraform-provider-snowflake/commit/f61051a685007856c07da62610eaecb9c019d4d9))
* Fix resource acceptance tests setup part3 ([#3489](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3489)) ([12f6ca0](https://github.com/snowflakedb/terraform-provider-snowflake/commit/12f6ca09a1e14f0a07b929629472512bb42fdc0c))
* Fix resource acceptance tests setup part4 ([#3504](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3504)) ([6a3b69c](https://github.com/snowflakedb/terraform-provider-snowflake/commit/6a3b69cce4201166c52ad9172ade9af09bcf4483))
* Fix resource acceptance tests setup part5 ([#3514](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3514)) ([45abb78](https://github.com/snowflakedb/terraform-provider-snowflake/commit/45abb78c1efe28697ff902cdf3e35ec5b42180ec))
* Increase stale objects timeout ([#3532](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3532)) ([efe15a1](https://github.com/snowflakedb/terraform-provider-snowflake/commit/efe15a1cb1354e4108a2e8c303b61d8d7a601b29))
* ok-to-test permission ([#3475](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3475)) ([0630eb9](https://github.com/snowflakedb/terraform-provider-snowflake/commit/0630eb97903ab5f3197c28cd7029e87f065d0709))
* Remove the outdated section from README's ToC ([#3471](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3471)) ([8398a44](https://github.com/snowflakedb/terraform-provider-snowflake/commit/8398a44acbb875d29c7b285fd5e8746911d5a1b6))
* Return a warning when grant_privileges_to_account_role is used with all_privileges ([#3512](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3512)) ([a2981b1](https://github.com/snowflakedb/terraform-provider-snowflake/commit/a2981b10fd5a9dcbd9aa2eee5bd275f4e2bcbc77))

## [1.0.6](https://github.com/snowflakedb/terraform-provider-snowflake/compare/v1.0.5...v1.0.6) (2025-06-24)


### 🔧 **Misc**

* prepare v1.0.6 release ([dc2f79e](https://github.com/snowflakedb/terraform-provider-snowflake/commit/dc2f79e141c9a1914219d39068be81d8ec63543f))


### 🐛 **Bug fixes:**

* Fix account parameter mapping ([#3794](https://github.com/snowflakedb/terraform-provider-snowflake/issues/3794)) ([1d13b91](https://github.com/snowflakedb/terraform-provider-snowflake/commit/1d13b9183c95848d1526777e59d041923a4ec185))

## [1.0.5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v1.0.4...v1.0.5) (2025-03-25)


### 🔧 **Misc**

* Bump tracking version to v1.0.5 ([#3460](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3460)) ([3aaf555](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3aaf555e72772da8e1088cb08bfd4a66f0bbb926))
* Document godebug flag usage ([#3429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3429)) ([2dce172](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2dce1724570c40e77dbaefe283c371491f682d91))
* fix failing tests ([#3428](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3428)) ([245fd86](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/245fd866ec1b79f680f5b9ff363439a6aaeb7f22))
* Fix migration guide ([#3468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3468)) ([6540841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6540841542d258e943d3f901397eb0d708029fc0))
* GitHub actions cleanup ([#3431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3431)) ([a6e7429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6e74292643db3ad27ad2296d6ba1f03215660f0))
* Protect pentesting user and role from sweepers ([#3426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3426)) ([6354dba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6354dba3a6bb7c20b9666436764c13d3251baa1a))
* Remove driver instrumentation ([#3439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3439)) ([e9fc6ed](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9fc6ed2020b561dc32e8276fd1b1dcab66c1c90))
* Remove SF_TF_ADDITIONAL_DEBUG_LOGGING ([#3441](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3441)) ([b57db81](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57db81966c80ecf5f27d3fe8f0878af86152ffc))
* Update GitHub actions ([#3455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3455)) ([368f090](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368f090cdf3b7acf3c2f06f61603207edabaa59a))
* Use environmental secrets in GitHub Actions ([#3440](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3440)) ([8e5a262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8e5a2629b38be9e8bb31ada43f545bf778a8b0a0))


### 🐛 **Bug fixes:**

* Add boolean env validations and unit tests for TOML config validation ([#3453](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3453)) ([408db4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/408db4eefbf49bf6494ad3afd46fe95dd0315cc1))
* Adjust docs ([#3451](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3451)) ([29b0e4d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/29b0e4d901fcb79e0f8837b62bf17166cd6b6cc9))
* Apply new assertions setup ([#3409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3409)) ([af83da0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/af83da0058655fddbc3469a2fc63c5fe24b69e77))
* Correct a typo in a tag resource example ([#3446](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3446)) ([afd19a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afd19a208be8b1336c756d8d7fa4a51596c74984))
* Fix datasource acceptance tests setup part1 ([#3443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3443)) ([5046654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5046654ecd503c531dda8787edde96bc7268402a))
* Fix datasource acceptance tests setup part2 ([#3445](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3445)) ([4ba2bb6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ba2bb6835f89068242fe7d851b2e2fe8a0bb51d))
* Fix datasource acceptance tests setup part3 ([#3448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3448)) ([f6fdc17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f6fdc17cbc8ffba56b581ccbc92db9de0b6b6594))
* Fix datasource acceptance tests setup part4 ([#3449](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3449)) ([935d969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/935d9698eb2645a4532c101b23e141c0b3968478))
* Fix tasks tests ([#3434](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3434)) ([b5457f6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5457f61efd654b0c9d126ed7d57ed2ed4ded830))
* Handle TOML file permissions ([#3444](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3444)) ([e4ed171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e4ed171aabd3ce2173ce429df01abb029bf96e6a))
* Limit TOML file size ([#3432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3432)) ([2e12981](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2e12981b0155de46d1839b4dff2714fbcbf82327))
* Quick fix assertions ([#3438](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3438)) ([a430ba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a430ba79a9826d823b8617b53cfdd1c7174f01b5))
* Skip file permission verification by default ([#3476](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3476)) ([562173d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/562173d38643bb4ab3ec388b64226e37362bcea9))

## [1.0.4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v1.0.3...v1.0.4) (2025-02-26)


### 🔧 **Misc**

* Add GA scope to the roadmap ([#3385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3385)) ([9be2196](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9be21966ca5b1187d5c055026dd5f563900a1aa9))
* Adjust saml2_integration documentation ([#3415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3415)) ([b8c127d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8c127d3e19f27a541010c08a94c8289166e00a5))
* Bump Go version ([#3408](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3408)) ([902670e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/902670e2eb7fa38033e5348cb403bc57dc6c4862))
* Bump provider version ([#3419](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3419)) ([552e44b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/552e44b3505a3db343bee1bbe240d6b30e61c3b9))
* Remove hardcoded values from documentation ([#3381](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3381)) ([cb1d6e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cb1d6e22af5584ec6d25e3c440dcc2275fab6544))
* Update user docs ([#3416](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3416)) ([9d5224f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d5224f01762104de67187e3f47662af528fd365)), closes [#3404](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3404)
* Upgrade deps ([#3389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3389)) ([086bf15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/086bf152e97c1873b9f1c55a1d1f2e223d05a56e))


### 🐛 **Bug fixes:**

* account for null comment ([#3417](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3417)) ([2f7f80f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f7f80f7d35c02e0e3ad56b3bb77faee385601c4))
* external function varchar return type validation ([#3400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3400)) ([abf5883](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/abf58839374281cbe66ee4430b73a501de01cc63))
* Fix managed account and adjust account ticket numbers ([#3395](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3395)) ([a7cb5cb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7cb5cb42a522ce3787129b3923ca6a83d34368a))
* Propose changes to assertions setup with varying test client ([#3406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3406)) ([fb27f6a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb27f6a9d42be8e2f148f3f78b401cfa09012fcc))

## [1.0.3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v1.0.2...v1.0.3) (2025-02-05)


### 🔧 **Misc**

* Add basic performance tests ([#3349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3349)) ([c57346b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c57346bc9fdcc263ed0c6d8dd34b51ff0a753968))
* Add notes about managing grants on hybrid table ([#3368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3368)) ([e5e98bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e5e98bd59a25eae994f389c4e4c940d3aa1958ca))
* Add preview and stable categories to resources and data sources in docs ([#3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3370)) ([b061434](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b061434f9f67e54bdd32c783d01ec4d4b74e3bf8))
* Fix docs and update the protected users list ([#3365](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3365)) ([f23e8cb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f23e8cbcc5326ca60bb8f759c6970ecf4b0d64d6))
* Grant ownership common use cases ([#3356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3356)) ([97813b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/97813b6308651071f05d8d4e8697f3a34e0f531a))
* Metric level fix with object renaming guide ([#3376](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3376)) ([629ff92](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/629ff92623a3d7a74eea970134291a84a9fec6a3))
* Move technical guides to the guide directory ([#3371](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3371)) ([2b95809](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b95809bc7619cb589eb08fd4e3a9dcbf1fcfbc4))
* Prepare authentication methods guide ([#3364](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3364)) ([c4a1c5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4a1c5ff882686377682cb6828549868b4efc909))
* Publish Performance tests summary ([#3359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3359)) ([d30d002](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d30d0024668e8825e4cb85a393d085a7a5423734))
* Resolve issues and document account resource ([#3360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3360)) ([46b7a9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/46b7a9d75426734bdc14edc3161070e84744f702))
* Upgrade tfplugindocs ([#3361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3361)) ([89f2b0a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/89f2b0abdacf86dcb4f2e47e4119dea993e749f6))


### 🐛 **Bug fixes:**

* change EnforceNetworkRulesForInternalStages ddl to parameter ([#3343](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3343)) ([beb01c7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/beb01c7bbedeb888d8c3af1460669ad3c8194d32))

## [1.0.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v1.0.1...v1.0.2) (2025-01-20)


### 🔧 **Misc**

* Generate ID and ObjectType Show Object Methods ([#3292](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3292)) ([f7ff70a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f7ff70ab0fe571d1fc616a6b202501533d132d43))
* Generate multiline by default ([#3315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3315)) ([520bb4c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/520bb4c57b8ec26d6a6284e219e45133ffb6be9f))
* Overview of the grant ownership resource ([#3342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3342)) ([9f891d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f891d64692fbe357ba96369fbed2629c665c402))
* Update docs and migration guide ([#3313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3313)) ([59e333a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/59e333aeaf188760e35c1e14a40433af9b79d8ca))


### 🐛 **Bug fixes:**

* Fixes in account parameters ([#3310](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3310)) ([96a1f17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/96a1f17ce839288566b9c4b58d168e2cfe1fac89))
* Handle optional account fields in the state upgrader correctly ([#3330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3330)) ([4eae4c8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4eae4c8a3f27e81286479e0522513c745c83d41c)), closes [#3332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3332)
* Small fixes ([#3337](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3337)) ([8807839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8807839f14b0757071ac62dfff50774cb144b4ac))

## [1.0.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v1.0.0...v1.0.1) (2024-12-20)


### 🔧 **Misc**

* Add external removal tests to functions and procedures ([#3305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3305)) ([e2d0705](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d07055015803b86f3c165851a63350beba350f))
* Add holidays note, change the disclaimer, and fix an example ([#3288](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3288)) ([39a0cb7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39a0cb76a1e64f2c4b5ee5923223d6e0df0ddb46))
* Adjust GitHub templates to the new development branch ([#3286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3286)) ([fc10672](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc10672318af90ab6d879dfa63b8efaf93589e62))
* Fix tests after v1 release ([#3306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3306)) ([32983fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/32983fe0d3ff502117b9c81d0d2c49b9d18958ba))
* Generate marshal json for each model ([#3307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3307)) ([7ebbe36](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ebbe364c5829546a70843d249629f8fcc3cdcf6))
* remove generating integration tests placeholder file ([#3256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3256)) ([e0ed15b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0ed15b09017ddfc89c727a6f5a8e4bdaf7aa34d))
* Update docs and migration guide ([#3313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3313)) ([c98e449](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c98e44922a090495feb7b5c31c23b01d07f31bfc))


### 🐛 **Bug fixes:**

* Fixes in account parameters ([#3310](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3310)) ([0de5733](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0de57338fe32200a44bdd89730eb8bccaeeac3d5))

## [1.0.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.100.0...v1.0.0) (2024-12-12)


### ⚠ BREAKING CHANGES

* True v1 release ([#3283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3283))
* Release v1

### 🎉 **What's new:**

* True v1 release ([#3283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3283)) ([112c852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/112c85244191a65bf552ffe27285299d6cc3831f))


### 🔧 **Misc**

* release 1.0.0 ([b8ddbf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8ddbf91bafdc8eef21a89da1ea86e6a455e2b96))
* Release v1 ([#3281](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3281)) ([82f240e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82f240eeed9ec24d6afc82cf5d0106544bec5838))


### 🐛 **Bug fixes:**

* Add missing preview features ([#3284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3284)) ([3aac502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3aac502beb404449674c54d153ca5aab6ef3cbaa))

## [0.100.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.99.0...v0.100.0) (2024-12-12)


### 🎉 **What's new:**

* Account v1 readiness ([#3236](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3236)) ([5df33a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5df33a8c1abe38c29124bac1e03727202c556347))
* Account v1 readiness generated files ([#3242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3242)) ([3df59dd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3df59dd51b53acae9155b732811cfda56d7f20b8))
* Account v1 readiness resource ([#3252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3252)) ([8f5698d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5698dbce3325461d572c4029ef2dbc364e819b))
* Add a new account roles data source ([#3257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3257)) ([b3d6b9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b3d6b9e5b4f327b186161f50dc9ac732d199fb19))
* Add account data source ([#3261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3261)) ([6087fc9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6087fc9fdb2467e022ec7489137e7f5a5fe1ff25))
* Add all other functions and procedures implementations ([#3275](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3275)) ([7a6f68d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6f68df2fb0a0a4696a5442569344039a839c27))
* Basic functions implementation ([#3269](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3269)) ([6d4a103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6d4a10364276e92fa791eaa022c3bd7bce16228d))
* Basic procedures implementation ([#3271](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3271)) ([933335f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/933335f56d1e53bf3e95d1f552672f35425b4878))
* Docs, test, and missing parameter ([#3280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3280)) ([10517f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/10517f337c6b22d5f7f2a4f6c747b6fd2d2f47e9))
* Functions and procedures schemas and generated sources ([#3262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3262)) ([9b70f87](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9b70f872ca799126bc2051b4ed70160f868ac267))
* Functions sdk update ([#3254](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3254)) ([fc1eace](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc1eace306e8d919c3349d56480fa3386ca664af))
* Handle missing fields in function and procedure ([#3273](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3273)) ([53e7a0a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53e7a0aea3350e9e03a804d67e7df796f15bff3a))
* Procedures schemas and generated sources ([#3263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3263)) ([211ad46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/211ad46223f1bdf03b20cc7a06110bfce18a967e))
* Procedures sdk update ([#3255](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3255)) ([682606a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/682606adea5e40befa7e599ced5aa7dc8570f80a))
* Rework account parameter resource ([#3264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3264)) ([15aa9c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/15aa9c2c94d80ae1d299a333b8035e38de6a6dfc))
* Rework data types ([#3244](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3244)) ([05ada91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/05ada917414ea7c574be3974c7de4f09535961fd))
* support table data type ([#3274](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3274)) ([13401d5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/13401d5fff320eedcf40eed7c0831154cc6cc13a))
* Tag association v1 readiness ([#3210](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3210)) ([04f6d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/04f6d54a83cf4e9ea4b292087eefa056114eb5b5))
* Test imports and small fixes ([#3276](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3276)) ([a712195](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7121952892847f61e24e7a7a4fe78c38a450985))
* Unsafe execute v1 readiness ([#3266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3266)) ([c4f1e8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4f1e8fd55150e40d8a556580016ff83fe65bdaf))
* Use new data types in sql builder for functions and procedures ([#3247](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3247)) ([69f677a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69f677a6f86faa79cdece4d422eb61284c1599a6))


### 🔧 **Misc**

* Add ShowByID filtering generation ([#3227](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3227)) ([548ec42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/548ec42ae7bcb8daa038de4cb2f81ced9c028f2d))
* Adress small task-related todos ([#3243](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3243)) ([40de9ae](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/40de9ae93796afbc3091aa2fbb2c5dfba71f911c))
* Apply masking ([#3234](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3234)) ([c209a8a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c209a8ae6c15fa9515e933d18add962070b60257))
* fix missing references in toOpts and changes with newlines ([#3240](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3240)) ([246547f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/246547f8eb13118a325881630f33433b3f5d8f0a))
* function tests ([#3279](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3279)) ([5af6efb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5af6efb08c479edbaea54f87f79672c802edcc86))
* Improve config builders ([#3207](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3207)) ([425787c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/425787c5938e88895af1157f505889611bdef398))
* Revert to proper env ([#3238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3238)) ([5d4ed3b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d4ed3bc233a77196f01351d9c972bb56730298e))
* Use service user for ci ([#3228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3228)) ([2fb50d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2fb50d7b5beb1f361d3c761b344bab3216f6ea59))


### 🐛 **Bug fixes:**

* Make blocked_roles_field optional in OAuth security integrations ([#3267](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3267)) ([7197b57](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7197b57c5dd75be34fc77eb82aabbd091074b809))
* Minor fixes ([#3231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3231)) ([1863bf6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1863bf697f05177f27c351c0687c4bee24fe2c1b))
* Minor fixes 2 ([#3230](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3230)) ([73b7e74](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/73b7e74bf44b1ae6ddc78cac752f2b7febb836cd))

## [0.99.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.98.0...v0.99.0) (2024-11-26)


### 🎉 **What's new:**

* Add tags data source ([#3211](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3211)) ([8907d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8907d9dfea69d6b8ac26fc0a9e249676f332f8b3))
* Tag resource v1 ([#3197](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3197)) ([77b3bf0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/77b3bf0c9998c05a30951730439f8b03a2e418ac))
* Tasks v1 readiness ([#3222](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3222)) ([e2284d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2284d98d23586031514934d7bc7c67139f5e272))


### 🔧 **Misc**

* Add support for usage tracking to data sources ([#3224](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3224)) ([8210bb8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8210bb84b69fe91e0fff22ac836feb79d6e9a402))
* Add usage tracking for the rest of the resources and fix views ([#3223](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3223)) ([231f653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/231f65323611f110564117a325062355e7ed7cf6))
* Basic object tracking ([#3205](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3205)) ([1f0dc94](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1f0dc94e6ac95940ac5fd0e0b5f62152b8f821a5))
* basic object tracking part 2 ([#3214](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3214)) ([e44f2e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e44f2e1938807285ed4d521b56d2efeab7b927bb))
* Improve tags integration tests ([#3193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3193)) ([7736e0a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7736e0a5fa6a97f9e5551507cea955fb62dd1e90))
* parser and secret tests  ([#3192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3192)) ([5ec9c86](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5ec9c86fdc3450f6f07820a4a5fe7f74779c7c41))
* Storage integration with custom protocol ([#3213](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3213)) ([a3a44ae](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a3a44ae5a6eca2a9623369499d8cac4516a87004))
* Unskip auth config tests ([#3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3180)) ([46ab142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/46ab142ad74e5fdc5deb6cc6edc409f487434862))


### 🐛 **Bug fixes:**

* Small fixes and adjustments ([#3226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3226)) ([9f67457](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f6745743daba831422627b5171df404373e9650))

## [0.98.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.97.0...v0.98.0) (2024-11-08)


### 🎉 **What's new:**

* Add authentication policy resource ([#3098](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3098)) ([ddea819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddea819cb9d0dd7918e0f4bbdaa0a2204da7b8b5))
* Add Resource for External Volumes ([#3106](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3106)) ([64ba674](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/64ba6747f9b2364a1f84a43242472af3c4ebeca7))
* Add stream on directory table ([#3129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3129)) ([4391473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/439147392436b9427e31ea578cc5bb971189a932))
* Add stream on view ([#3150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3150)) ([494af6d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/494af6dd19d1098d76e7ea0451d1e144138c7a29))
* Connection datasource ([#3173](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3173)) ([4127b3f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4127b3f782660772b77a72f9b38dafc728254de3))
* Connection resource ([#3162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3162)) ([5aef117](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5aef117f415f238a0e786b9a063d44fadeb879e5))
* Rework config hierarchy ([#3166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3166)) ([04cd9f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/04cd9f04f1713b47f50c1931fd8955665ea8cbcc))
* Rework provider configuration fields ([#3152](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3152)) ([fd6af43](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd6af43d5e50e83d247333e4a0edf85008538a9e))
* Rework streams data source ([#3151](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3151)) ([b18bf30](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b18bf30eee1b86e1f85c52e36122e9fa052053bd))
* SDK Connection ([#3155](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3155)) ([bd11e0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd11e0f5eab8fe420f8af6644f4c1eb90910e69a))
* Secret resource ([#3110](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3110)) ([16a812d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16a812dae20d59f31457790dcd99db03db697051))
* Secrets datasource ([#3131](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3131)) ([8110138](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/811013887cdf1a6624c93481f18e12a183865463))
* Upgrade tag SDK ([#3126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3126)) ([893b288](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/893b288f91ee3d31c875f2e38c7346fe362632e6))


### 🔧 **Misc**

* Add a company name field to the issue templates ([#3182](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3182)) ([0d3248a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d3248a2cb5323f94d8a984aaabe60857041fc3b))
* Add object renaming research summary ([#3172](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3172)) ([721ee40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/721ee40d6954254fc1f45eec50a7098b45d1afc3))
* Add tests to 3117 and bump build time ([#3133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3133)) ([ca90fde](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca90fdefc18f6c77627e602df46650e17cb54eaa))
* Detect changes in lists and sets ([#3147](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3147)) ([c3edb79](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c3edb79758e152a60370fcaa251fdbbdda9bcc56))
* Exclude methods from test function checks in architest ([#3174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3174)) ([edc46cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/edc46cc29fa9d0f19f5f62f04de71c5a56e4b902))
* join object renaming tests into parameterized ones ([#3154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3154)) ([be13502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be13502d7fc67c896ee1387297bddad904f515b0))
* New roadmap entry ([#3158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3158)) ([d83cdde](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83cdde59a5fb98fea5d31e5017869f841c678cb))
* Test more authentication methods ([#3178](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3178)) ([d345cd2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d345cd291170db27e7ba1b83d69e613938fd9538))
* Test support for object renaming ([#3130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3130)) ([d665419](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6654195c099b81969bb8043e98771cf199ee0f4))


### 🐛 **Bug fixes:**

* Apply various fixes ([#3176](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3176)) ([55591da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55591da967db58a3ef11e3e454101e27ca7abb42))
* Connection and secret-datasource tests ([#3177](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3177)) ([167de4b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/167de4be15530d481b6e46953a7984d1b2777899))
* Fix grant import docs ([#3183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3183)) ([94ac910](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/94ac910f490a4cc3d9975e073e7503781bab2ffd))
* Fix main ([#3157](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3157)) ([89b9705](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/89b9705d2af43aad3f0e098e3b8af7ced0b3d406))
* Fix main ([#3160](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3160)) ([5b7412f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b7412fe00a3620d1428d61ce30f51232eccf213))
* Fix main ([#3186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3186)) ([59a0a26](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/59a0a2699cbbb0a18517b607e736a04b62f6c3ba))
* Fix user resource import ([#3181](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3181)) ([34bbbc1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/34bbbc18b61ce1224c4f8607a12a9f6dfd95c958))
* handle external change of secret type ([#3141](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3141)) ([649b839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/649b8397987d13fa53c67729a81dc7fceef218a7))
* Handle external type changes in stream resources ([#3164](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3164)) ([9fd8f88](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9fd8f887401c0931bbbd18d3b3c4d770b8410fd4))
* merge diffs on test clients ([#3149](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3149)) ([0f06b4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f06b4a063c9a9f9922aa4d83c47935d09757571))
* Skip connection data source acc test ([#3184](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3184)) ([2942374](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2942374f668b4bdc48f81c7580262fc1c6473179))

## [0.97.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.96.0...v0.97.0) (2024-10-10)


### 🎉 **What's new:**

* Add secret to sdk ([#3091](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3091)) ([7430aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7430aeeee8bfd59f3e48e689d37c3969ea7e367e))
* Add service user and legacy service user resources ([#3119](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3119)) ([0e88e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0e88e082282adf35f605c323569908a99bd406f9))
* Handle all Task parameters in SDK ([#3103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3103)) ([08ae072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/08ae0729acea25293d4143e9f70f77fe4bfbd29e))
* Stream on external table resource ([#3122](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3122)) ([d837341](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d837341c2d18b6fbb4657ad3a1837190a8ee77d8))
* Stream on table resource ([#3109](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3109)) ([97fa9b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/97fa9b4485cf26c5bf93dedfe3b88f688a71b3e6))
* Tasks v1 readiness - SDK part ([#3086](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3086)) ([0a77383](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a77383ff51a13be40db72f4d0462fb92d0a4fbd))
* Upgrade stream sdk ([#3105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3105)) ([ad5fa11](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad5fa113da84a2a1aeaa275643cd586126e6f603))


### 🔧 **Misc**

* Add pre check to each datasource ([#3065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3065)) ([560ab6b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/560ab6b444ac461684b0d25c68bcb2d55201aba8))
* Bump golang-ci lint to 1.61 ([#3112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3112)) ([f23e085](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f23e085536febfddec941b49ac0033f93e8662cc))
* Secret Validation change ([#3111](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3111)) ([666630e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/666630e4a865bebafaa95e1afb149d27fa90a935))


### 🐛 **Bug fixes:**

* Fix parsing text in view, check parenthesis in ParseSchemaObjectIdentifierWithArguments ([#3102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3102)) ([b0a67e6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0a67e638314e795ffd7075c302db51e6d0c73d9))
* Try to reproduce 2679 and 3117 ([#3124](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3124)) ([ccdbc30](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccdbc306f2a074cbabf0ba96d67e522eb82d36e7))

## [0.96.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.95.0...v0.96.0) (2024-09-18)


### 🎉 **What's new:**

* Add authentication policies to SDK ([#2937](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2937)) ([7b49685](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b496859d07ee650a52735ffc77b398ca6aa01f9))
* Add SDK for External Volumes ([#3033](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3033)) ([2844a30](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2844a301b4d9e303c7d47f1badc2fe5c3acd34e3))
* Masking policy data source v1 ([#3083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3083)) ([2553fbf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2553fbf75f551ac736896155d787384e3b205ef7))
* Masking policy resource v1 ([#3078](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3078)) ([fa4ce56](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fa4ce56c1d4071722ee8dd039f7090c07afc0057))
* Resource monitor v1 readiness ([#3052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3052)) ([ac493f4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ac493f4dab2e6efe2fef05e39fd7cdc622a729ca))
* Resource monitor v1 readiness part 2 ([#3064](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3064)) ([01b66cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/01b66cf0f2977071285bf5988c790cb257b61b3e))
* Row access policy data source v1 ([#3066](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3066)) ([138a288](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/138a288c54971d62af32f74635f083c198f3774c))
* Row access policy resource v1 ([#3063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3063)) ([13cd694](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/13cd6948d98a1e40647407ec2e8c5bd0f0635917))


### 🔧 **Misc**

* Add Filip legacy to protected users ([#3054](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3054)) ([5f54ff9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f54ff98d17f789872912ad9ef6e5ee21fb53b58))
* Add Filip to protected users ([#3053](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3053)) ([74b8ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/74b8ea8dc209b42a395e575f77bb1cb96373747c))
* Adjust authentication policy ([#3068](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3068)) ([f39a68c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f39a68c7af69a22a656aa74fbe1d7eac36f245ed))
* Clean up old test object helpers ([#3049](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3049)) ([78ecb62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78ecb62935c1a40611418e59ed5f0d644cf85f03))
* Example of granting role to multiple objects ([#3047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3047)) ([86a7902](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/86a7902562d81c1eeb47be4530587e87ed8cfd44))
* Update readme and objects rework state ([#3046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3046)) ([5802dca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5802dcae57ab7a0ed24b02d2c492394cf2f53e70))


### 🐛 **Bug fixes:**

* Add fix for model grants ([#3070](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3070)) ([cbc3d74](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cbc3d746c7153025598b4c9d5cccfaae3a957c0e))
* authentication policies ([#3061](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3061)) ([a65e564](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65e564d49bffc7c2e215adb7bd4138f398c00a8))
* Fix database show by and resource logic ([#3055](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3055)) ([1887e55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1887e55908da3f689724f729d012740841995905))
* Fix default secondary roles option import ([#3041](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3041)) ([0a3ce34](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a3ce34f5f78db6438eceb6c58fb1dc21d16e6a5))
* Fix sweepers for warehouse and database ([#3057](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3057)) ([6dd586b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6dd586bbe8a5945c0324ee2f0b5c16535859544e))
* Fix views permadiff ([#3079](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3079)) ([c2d3e5c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c2d3e5c8104dcd43e81dcd3027a7fd3184f3a578))
* Update v0.95.0 migration guide ([#3062](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3062)) ([d5dac60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5dac6005bf2d44f730e51bdab452fd1b4f1aaff))

## [0.95.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.94.1...v0.95.0) (2024-09-04)


### 🎉 **What's new:**

* Add change_tracking, row access policy and aggregation policy to views ([#2988](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2988)) ([1f88bb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1f88bb1aa11397fedd0653c44a56b9e3d294c362))
* Add fully_qualified_name to all resources ([#2990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2990)) ([1b0462f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b0462fe4ffdb190211f3188fb0c388ad07a7e1a))
* Add identifier parsers ([#2957](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2957)) ([824ec52](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/824ec52228afcdec282672018fc41c35a42a429d))
* Add identifier with arguments ([#2979](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2979)) ([00ae1c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/00ae1c5a3384e31697830d84e9c53f9dce2cb6e6))
* Add timeouts block to cortex ([#3004](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3004)) ([34d764b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/34d764b98988431134417346c7d4aa9f3b25e3ea))
* Add user parameters to resource ([#2968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2968)) ([f4ae380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f4ae38086c9a17259aff698f80a21e8d14f748e7))
* Conclude user rework ([#3036](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3036)) ([23e4625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23e462542f6295c8d9e7e2ca13b28aad1a3a2d40))
* database role v1 readiness ([#3014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3014)) ([c4db255](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4db255c3e3b892f8d4d7da628c0344515a185b4))
* Identifier with arguments for procedure and external function ([#2987](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2987)) ([f13cc5c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f13cc5c2c52c84c6e8d417a22e0c8508e921f5bc))
* Rework user resource ([#3026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3026)) ([bde2638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde2638c1040bba5ffb54627badabd18d2e0b26e)), closes [#1572](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1572)
* Rework users datasource ([#3030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3030)) ([751239b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/751239b7d2fee4757471db6c03b952d4728ee099)), closes [#2902](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2902)
* Upgrade view sdk ([#2969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2969)) ([ef2d50a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ef2d50a5d69b296eaa32aa39a089a956379e3a6c))
* View rework part 2 ([#3021](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3021)) ([e05377d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e05377d251fbec831214fc1c481b67d786dc829b))
* View rework part 3 ([#3023](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3023)) ([195b41c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/195b41c927d8813f488ecea65c820f56f7e2853e))


### 🔧 **Misc**

* Add annotation about fully_qualified_name and fix handling granteeName ([#3009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3009)) ([94e6345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/94e6345a6b3c3a09d11ec5cc46c0e5fadb5193f9))
* Apply identifier conventions ([#2996](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2996)) ([5cbea84](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5cbea8487a364011428575f4705a05993fe43724))
* apply identifier conventions to grants ([#3008](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3008)) ([d7780ae](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7780aed31dcdfbe4dc252b06fb74eb5d9483c57))
* Clean collection utils ([#3028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3028)) ([426ddb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/426ddb11723ab7d583b1bc4f74f2a88f8454caa0))
* Clean old assertions ([#3029](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3029)) ([ad657eb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad657eb1e04d56b22001828fa9f22f160e7c6494))
* Conclude identifiers rework ([#3011](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3011)) ([c1b53f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1b53f3c6472467e675cd137671c48f27432646c))
* Improve user test and add manual test for user default database and role ([#3035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3035)) ([6cb0b4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6cb0b4e9422e22c8fe0554248aec0b8fd7b2849a))
* Use new identifier with arguments in function, external function and procedure grants ([#3002](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3002)) ([5053f8b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5053f8b9f9657a50ccce407f29a1c017596563b8))
* User improvements ([#3034](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3034)) ([65b64d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/65b64d78f94fdc444c482bc52c7e84c5db5b4d2e))


### 🐛 **Bug fixes:**

* database tests and introduce a new parameter ([#2981](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2981)) ([3bae7f6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3bae7f634f86c938bf7f00170f25ee1f102d9abb))
* Fix custom diffs for fields with diff supression ([#3032](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3032)) ([2499602](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2499602b4e18a9fea5627053bd1aced748f5ec46))
* Fix default secondary roles after BCR 2024_07 ([#3040](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3040)) ([2ca465a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2ca465adcc6d28f4e10f42bae15462075a6010de)), closes [#3038](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3038)
* Fix issues 2972 and 3007 ([#3020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3020)) ([1772387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1772387a710ff92886e3fe7f30732595f092bb6a))
* Fix known user resource issues ([#3013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3013)) ([a5dfeac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a5dfeac5a825eed6cf0f847887e288e5a648a3f9))
* identifier issues ([#2998](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2998)) ([6fb76b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6fb76b78c02f0d37e851902c9caf05014a905e58))
* minor issues ([#3027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3027)) ([467b06e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/467b06effca5eba2b64b4e49d17c5c8e0de2da23)), closes [#3015](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3015) [#2807](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2807) [#3025](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3025)
* Nuke users ([#2971](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2971)) ([0d90cc9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d90cc948eccc2e281860d341281d13c27805a0b))

## [0.94.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.94.0...v0.94.1) (2024-08-02)


### 🐛 **Bug fixes:**

* Use ALTER for managing PUBLIC schemas that exist ([#2973](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2973)) ([567e9be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/567e9be5efb0be731fa7ee56143b8ca4326bd037))

## [0.94.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.93.0...v0.94.0) (2024-07-26)


### 🎉 **What's new:**

* Add missing session parameters ([#2936](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2936)) ([4ce662d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ce662d7360b2882c414074fc5aa021e1bd70433))
* Adjust user SDK ([#2947](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2947)) ([1127bb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1127bb3fbd5e21e3176acf3ffd012e18440bac1c))
* Better tests poc ([#2917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2917)) ([ef496c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ef496c23a28dde8fc6dff59d47ec79eb8f4e89d4))
* Introduce assertions generators part1 ([#2952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2952)) ([1582a9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1582a9fe73d18bf89529bf12e391a2b300e42eaa))
* Introduce assertions generators part2 ([#2956](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2956)) ([f715e8a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f715e8a990b72b62f87e526f7b30a7b6e9db5f19))
* network policy v1 readiness ([#2914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2914)) ([3408c3f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3408c3f704da9490e518a4206ed52396a2c2315c))
* Rework schema datasource ([#2954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2954)) ([f70e40e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f70e40e82a7f22e86276ab55ddf2584882229ccd))
* Rework schema resource ([#2955](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2955)) ([400a5c8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/400a5c8f02af5134e7f5c5d1bb7f73cf89bf8627))
* Role v1 readiness ([#2916](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2916)) ([32c7690](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/32c76906c2eb8571e4a995df20f6450ddd6f94b4))
* Schema SDK upgrade ([#2945](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2945)) ([bca0836](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bca083607978df37513937cf1b7b2b2dfbe55ef8))
* Streamlit v1 readiness ([#2930](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2930)) ([aa42260](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aa422601b23bf73891b0f50bd4df7944cde56670))


### 🔧 **Misc**

* Remove deprecation from unsafe execute ([#2941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2941)) ([ed712d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed712d7e468269a01a66b9bfe77b5d8018783a30))
* Update documentation ([#2931](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2931)) ([da98bc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da98bc3f5f37f721fed623d7ab110e75ff9557e2))


### 🐛 **Bug fixes:**

* ATTRIBUTE set(string) parsing for cortex search service ([#2953](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2953)) ([70a1c9a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70a1c9a9f4a9982f678eba8cd8e795a02a414966))
* external function header parsing and add missing privileges ([#2961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2961)) ([9d882fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d882fe27332d4e49e1596721cba54ca0eb33b07))
* Fix sync_password field for Azure scim clients ([#2950](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2950)) ([6781133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/678113395f120305e002efc8d969307b0bbf244d))
* Fix tests and relax warehouse validations ([#2959](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2959)) ([dd01ce9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dd01ce9f5cc4aedef9091a2238bc9f8c09905b84)), closes [#2948](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2948)

## [0.93.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.92.0...v0.93.0) (2024-07-10)


### 🎉 **What's new:**

* Add OAUTH integration for custom clients ([#2908](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2908)) ([d9b557f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d9b557f7de6c20c16e9d39be9becc92a5a74fb40))
* Add oauth integration for partner applications ([#2912](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2912)) ([91788e5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/91788e5942c581c59fa8138e775e14c53c749956))
* Add support for cortex search service ([#2860](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2860)) ([43aa89f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/43aa89f16419943ccffdd10f63e2aa9034891d6a))
* API Authentication integration v1 readiness ([#2898](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2898)) ([91931da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/91931dacbc6c92bef7ef4299b7a5ddfaf5faf171))
* External Oauth integration v1 readiness ([#2907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2907)) ([ed237c3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed237c36e82291d75d38e2741bdcfd53e598a447))
* Generate show outputs with mappers ([#2886](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2886)) ([1cada88](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1cada880bd22bc9b0904b859730bf66e0c56639a))
* Introduce security integrations datasource ([#2892](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2892)) ([7f6c657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f6c65777966f9555a6bb8013667dd14cb44b431))
* SAML2 integration v1 readiness ([#2868](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2868)) ([d0c136d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d0c136d1e669ee97b76199ce1e906b469b279842))
* SCIM integration v1 readiness ([#2846](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2846)) ([269df6b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/269df6b11f14358c19dd37e2d6eeb75d33d633ee))
* Security integrations datasource v1 readiness ([#2913](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2913)) ([d10474a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d10474ababa7489dec3bf6c220b419b3e6a97c8b))
* standard database v1 readiness ([#2842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2842)) ([3c11953](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3c11953bf2800ac8297da14b2394d49e3e11dd39))
* Warehouse redesign final touches ([#2900](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2900)) ([0eab636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0eab636f9e16ba4ff60f20d81be82ebcd4e642d3))
* Warehouse redesign part1 ([#2864](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2864)) ([6664457](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/66644576f2e6ba77e0601090231bf42e5d08fcdf))
* Warehouse redesign part2 ([#2887](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2887)) ([1aaf417](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1aaf417de916a0aa6f9ba29eb79c9aea1b9c4bf6))
* Warehouse redesign part3 ([#2890](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2890)) ([873a1ed](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/873a1edc3a96b7d9f8eb304b43747f3843484a0f))
* Warehouse redesign part4 ([#2893](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2893)) ([d525fd9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d525fd952517afb540543ede22a6d2e818146ce5))


### 🔧 **Misc**

* Add documentation on unset and defaults ([#2882](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2882)) ([85a7836](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/85a78365490c2a1b8ec08b1b7883350cb41c1d76))
* apply minor database changes ([#2872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2872)) ([6ccac59](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ccac596d2115cb0ecf9475a3c10c2e0794259b0))
* Apply new resource conventions to scim integration ([#2891](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2891)) ([e11e608](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e11e6085c2896981c73c78640b5fc1743e29142b))
* Improve generator template organization ([#2820](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2820)) ([5035e2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5035e2f0632c9dfca26ab61b5f3af406af743f50))
* Nuke stale objects ([#2869](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2869)) ([9c4a117](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c4a1172da745af8f5059291e58b7cf8c8ec60c5))
* Show a possible solution for [#2877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2877) ([#2878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2878)) ([6fb437b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6fb437bfad221fe9322d24f5e3167ca161edc9c8))
* Validations cleanup and old grants removal ([#2884](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2884)) ([05b7eee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/05b7eeec0e29a633c30e58625ec4fb08bdaeeb1b))


### 🐛 **Bug fixes:**

* Add disclaimers and fix tests ([#2905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2905)) ([1deaedc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1deaedc95df9b82229256055a4803387098c58d9))
* Fix cortex search service ([#2904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2904)) ([763d06c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/763d06c9124bcc43aa918035f67a9ad0680e5635))
* use suppressQuoting to fix stage file_format permadiff ([#2885](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2885)) ([fd70f6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd70f6ef21e90e991763211dbacc0c3870864985))

## [0.92.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.91.0...v0.92.0) (2024-06-06)


### 🎉 **What's new:**

* Add Api Authentication security integration to sdk ([#2840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2840)) ([57a07ee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/57a07ee792fdf61299cf7ad51f06dd97408a2e52))
* Add External Oauth security integration to sdk ([#2835](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2835)) ([82d1c09](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82d1c09962d71d2b2a4bad23f8b11ec05cba1a7a))
* add network rules ([#2746](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2746)) ([c79fa29](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c79fa2991700f4e4ec523c0f102a946534a9597f))
* Add SCIM and SAML2 security integrations to sdk ([#2799](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2799)) ([1312ff1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1312ff14b33aee828604d0d64f6e7a6cfe1e2ebf))
* Add Snowflake Oauth security integration to sdk ([#2830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2830)) ([b576f29](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b576f297f3831cf60e9cc849bfe9e48d5febe4b0))
* Database resource v1 readiness ([#2834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2834)) ([30fe136](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30fe1364be6ed92c4ddb072167bf2150a60d4d10))
* Database SDK upgrade ([#2814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2814)) ([750fe37](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/750fe372a1b4698bde51bdf537cce91dd1605582))


### 🔧 **Misc**

* accept non-pointer values in the generated builder methods ([#2816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2816)) ([c29fbf1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c29fbf13674dd8441db0da4ff5d309fa0f6e925f))
* Add a script for creating labels ([#2778](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2778)) ([ce0fbad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce0fbad38fa926a2b9a6b8b9ffc53c2c1bf1f9f0))
* Adjust before 0.92.0 ([#2857](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2857)) ([0598656](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0598656d1f53fd31f11fc979555f0dea00fa107b))
* Continue random ids rework ([#2819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2819)) ([f20940c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f20940cc05029a3f3de5c51b21c8bac4c790cc3e))
* Random ids rework part3 ([#2833](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2833)) ([36ead85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36ead8576d0ea54f9cc361d2baf9593006dc801f))
* Random ids rework part4 ([#2837](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2837)) ([64518a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/64518a3a8305236353e5bf2de25a9c8a9723e10b))
* Update issue for table and warehouse redesign state ([#2845](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2845)) ([149e55e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/149e55ede78b01e8621ce7dda5c33df3284b584b))


### 🐛 **Bug fixes:**

* Fix failing integration tests ([#2832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2832)) ([2e2ca6c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2e2ca6cd32f88a4ab8453bae05554e4656d3a975))
* Fix QUOTED_IDENTIFIERS_IGNORE_CASE parameter test ([#2841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2841)) ([92ad1d3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92ad1d309c3fb5daf07fe71ea138256ca7e29c5a))

## [0.91.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.90.0...v0.91.0) (2024-05-16)


### 🎉 **What's new:**

* add snowflake_grant_application_role resource ([#2690](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2690)) ([838d241](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/838d241ddde33fba8d4e85ff3b2cd52d8c256aaa))
* datasource database role ([#2731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2731)) ([319ddc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/319ddc3df2f0a73c1f6258632de0e42c04532551))


### 🔧 **Misc**

* Add few documentation adjustments ([#2789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2789)) ([6db8bf3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6db8bf38d3a03b05bdbf53162697da66c8c22b77))
* Bump dependencies ([#2802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2802)) ([54ea6bc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/54ea6bcb42487db8a15782cd4b53fc5fba3223d2))
* Replace parsing function for saving granted object names ([#2813](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2813)) ([175cfc7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/175cfc746a8c9070234e8a0ee2a8d0b04c1f97cc))


### 🐛 **Bug fixes:**

* Fix sweepers ([#2800](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2800)) ([a01115e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a01115ea54fc54c4600a83b877765f3643856080))
* Invoke SetId immediately after alert creation ([#2786](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2786)) ([181b4d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/181b4d0b45699b5076f31a194d36c29735fd3bf7))
* Update the tests after snowflake bugfix ([#2806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2806)) ([6843c1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6843c1a35f201a16c5932faaff390e62bd7dc611))

## [0.90.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.89.0...v0.90.0) (2024-05-08)


### 🎉 **What's new:**

* Adjust owner_role_type and schema_evolution_record columns ([#2740](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2740)) ([424e393](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/424e393813b1c0d668250937bf04487fb323b9f0))


### 🔧 **Misc**

* Add a guide on creating issues ([#2758](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2758)) ([2b006aa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b006aa12d0d0b0fae72ae295cb82404262288be))
* Add missing documentation ([#2781](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2781)) ([cc0a6a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc0a6a7e605735d83e3247c27b75679f45142659))
* Add scripts to close issues from Pre Snowflake bucket ([#2762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2762)) ([44c0c37](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44c0c37047ee15cd2856056270a1fb708473c4d6))
* Add small adjustments ([#2783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2783)) ([e5b0b4b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e5b0b4b95bde9e0d29c0066278393a572dccd825))
* Adjust issue templates ([#2748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2748)) ([64ab76d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/64ab76d82662b8e21679e887d5b3df8d15cfaf3a))
* Cleanup helpers part 3 ([#2730](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2730)) ([eb7bee4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/eb7bee400f6a5ec94b6a142eaeed78dbc349fc00))
* Cleanup helpers part 5 ([#2744](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2744)) ([1f165bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1f165bf72e29d1a9135d2e5e8bf3f46a0ffdddb6))
* Cleanup helpers part4 ([#2741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2741)) ([9475e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9475e35dfca2c15876ee14ace0a6f17c7225d615))
* Cleanup helpers part6 ([#2745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2745)) ([eba3029](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/eba3029bd2b6def28901c47df7e142a417d6ba3a))
* Cleanup helpers poc ([#2724](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2724)) ([70b99fb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70b99fb19be76bff673e4ba7740d07b9a131ca4f))
* Helpers cleanup continuation ([#2726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2726)) ([a70d1af](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a70d1afc607c5e9501762448a96a11c22efca3d4))
* Prepare new roadmap entry ([#2773](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2773)) ([b0bf28f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0bf28f4e08f36a5aee74bba21ec2cf2d1d1179b))
* Prepare parallel builds ([#2737](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2737)) ([6974e25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6974e255ebdd772f18b164943c82aab2f19708d6))
* Update the contribution guidelines (and small fixes) ([#2753](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2753)) ([aafdc72](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aafdc72f21314f809f1ab1f7e4f00d5b2df5e0f5))


### 🐛 **Bug fixes:**

* Fix issue templates ([#2760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2760)) ([d0d5048](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d0d504868b0dda15e3a779e94969d7742253160c))
* Fix setup for two tests ([#2757](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2757)) ([df025b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/df025b0c4dce30554eaf15b29031357aababe192))
* Fix Test (wrong order of arguments) ([#2780](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2780)) ([02f467e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02f467ebe230e077ce0f67ae1d1025d9219f65ed))
* Fix/prove issues [#2733](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2733) [#2735](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2735) [#2763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2763) [#2767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2767) ([#2777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2777)) ([7b1c67e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b1c67e115f3dddeed796ae00f160efeeb64f85b))
* Prove problems with procedure resource data types ([#2782](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2782)) ([68d0318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68d03184047ebe6891fbcebb1524377e10fb82fd))
* read after create ([#2718](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2718)) ([2e9b68f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2e9b68f9da09150c5976d175a225ff368f37f2d2))
* UNSET and empty SET in network policies ([#2759](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2759)) ([3eacb0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3eacb0bd869bc073be756dca9337b69e07bb52f4))
* unset network policy should not have single quotes ([#2584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2584)) ([8f2a363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f2a363db7ea0bcc747ef7c3576d24255f4bf208))
* Update failover group allowed integration types ([#2776](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2776)) ([efde48d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/efde48db07a9eb47c34a5febf1b4de6e4f0671f9))

## [0.89.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.88.0...v0.89.0) (2024-04-18)


### 🎉 **What's new:**

* Update target objects for privilege-granting resources ([#2688](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2688)) ([74e2b6b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/74e2b6bafb6e5a347b6bfc9b1b3e7283a5ad06b4))


### 🔧 **Misc**

* Handle generic check destroy in acceptance tests ([#2716](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2716)) ([63a5324](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63a53247fb3a52fbc66206813900f6006205322c))
* Initialize the SDK client in fewer places ([#2710](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2710)) ([382bfc1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/382bfc1f805851b3b944437128cdf9da74e8d6f2))


### 🐛 **Bug fixes:**

* Adds case statement for ObjectTypeUser in ShowObjectParameter. ([#2675](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2675)) ([23a3341](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23a3341ab416e291fe617888ef09e94b9e69247d))
* diffs always occurring when multiple columns exist ([#2686](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2686)) ([3275ad4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3275ad47747ad556faa5043eebf1295f11ae5bcc))
* Fix issues [#2679](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2679) [#2721](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2721) ([#2723](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2723)) ([b0c9dd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0c9dd42221e50e931324fbc85ad69f890b2d874))
* Fix several small issues ([#2697](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2697)) ([e3f6a15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e3f6a154fc92b1bfd47fb68ae5ec6478aa971437))
* granting ownership on database roles ([#2703](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2703)) ([88944e7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88944e7a02b1ff930e1a39bc519780c124bdf4b2))
* network policy update ([#2647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2647)) ([8126b28](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8126b2872d2dfec4814dfb06141c0aa3284799e4))
* Rename grants_redesign_desgin_decisions.md to grants_redesign_design_… ([#2691](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2691)) ([5000b2b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5000b2b5bb845751eb10de890892deadb8df8bdd))
* renames in update operations ([#2702](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2702)) ([793c879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/793c879b31b0a5ac4fc25d34eae495b6e0446daf))
* showbyid method tests ([#2648](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2648)) ([ff5e617](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff5e617bc91fb35a7b95fe59aa1c12972f99790e))

## [0.88.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.87.3-pre...v0.88.0) (2024-04-09)


### 🎉 **What's new:**

* Fix issues 2651 2656 ([#2659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2659)) ([7fa09cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fa09cc8a296003a09c25a1e8b88f997f6307416))
* Grant ownership follow up ([#2628](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2628)) ([d467e5b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d467e5bccfa5d8ee61482608d7e748320b9574e3))
* grant ownership follow-up ([#2658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2658)) ([bfa2317](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bfa231742ae234a04f1b5759f45855554c971700))
* grant ownership on tasks ([#2684](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2684)) ([2ba7889](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2ba7889206aa73014def4673266967cb2a5be6cc))
* Introduce shared function to suspend task roots ([#2650](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2650)) ([d684b5d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d684b5d7b41aefef8a707d4983b9e1a5bcd01aef))
* Redesign snowflake_grants datasource ([#2667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2667)) ([918873d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/918873d19fff32f88acdce8e62e78fcbf24b244f))
* user password policy attachment ([#2627](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2627)) ([382e49d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/382e49d7254e2422d64d00c51e341644ce614d70))


### 🔧 **Misc**

* release 0.88.0 ([02d60e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02d60e07dc928a8d7398ca54f0b91a259255fd99))
* Update grant examples in all resources ([#2660](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2660)) ([b542b69](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b542b69b198dae06e71d16f8fbf3a13a30f72752))


### 🐛 **Bug fixes:**

* Adjust dynamic tables after BCR-1543 ([#2664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2664)) ([cf32ceb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf32ceb9c6780ec39c14022b983d73afafbd1c84))
* Fix handling secure views read and import ([#2655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2655)) ([3c3ede6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3c3ede6b1a0e1c752c8435457c92db6200f0e7b6))
* Fix issues 2606 2642 2652 2653 ([#2654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2654)) ([4a73721](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a73721ae890929032d70d3d831bb829711494ee))
* Fix release workflow ([#2639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2639)) ([dfd07e9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dfd07e934c5c3d31f1aab211368f85c80068dad6))
* Generate docs for dynamic table ([#2666](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2666)) ([16c75b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16c75b0053da0e52c779765665eb9361f4f32ef3))
* grant read operation ([#2665](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2665)) ([0b947a5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b947a54a7bd285a4eaa337b6efa69da66d42a17))
* migration guide for databases created from share ([#2637](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2637)) ([f9651bc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f9651bcc9c0ec1afd6ec2e1ee258de1791a76b80))
* with_grant_option diff drift ([#2608](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2608)) ([f0018c6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0018c6c0d89a537b269696a59fe255a3c1a6f28))

## [0.87.3-pre](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.87.2...v0.87.3-pre) (2024-03-18)


### 🎉 **What's new:**

* Add snowflake grant ownership resource ([#2604](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2604)) ([bfadd24](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bfadd243a70f1df99426d99c265a12befa93f5e6)), closes [#2549](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2549) [#2199](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2199) [#2084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2084) [#1942](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1942) [#1875](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1875)


### 🔧 **Misc**

* Fix env variables for tests ([#2603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2603)) ([8bc2437](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8bc24378c6fb233151be1f414ecf9ae768c6a255))
* release 0.87.3-pre ([a2be7b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a2be7b90d288be52e3ee1556d7127d320cf13f15))


### 🐛 **Bug fixes:**

* alter table column data type ([#2607](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2607)) ([538b6dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/538b6dc1048a5ccf6d08c1ad8eca520d1e4477fd))
* cgo goreleaser alt solution ([#2613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2613)) ([5d31856](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d31856a25fb5b7520a390b98202d741c75600a0))

## [0.87.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.87.1...v0.87.2) (2024-03-07)


### 🐛 **Bug fixes:**

* Add missing path ([#2597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2597)) ([41865d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/41865d95a0936be351eed7c024f38f3803581cb8))
* Fix the deprecation message for snowflake_role_grants ([#2602](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2602)) ([5d9c1a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d9c1a8324af388602cfddf78338e902a71e9d8e))

## [0.87.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.87.0...v0.87.1) (2024-03-06)


### 🔧 **Misc**

* Add deprecated resources/datasources to the main documentation page ([#2581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2581)) ([68bbf4f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68bbf4f4f89ed86d1a1deaf7145c9c13b0c7dff3))
* Add deprecation message to docs ([#2578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2578)) ([3675d6d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3675d6d6ad0bfa6f901d98686165a2297e6256ba))
* Add Terraform setup to the testing pipeline ([#2579](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2579)) ([216e35a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/216e35a8364a66f85e2c77abae0c3a587a493687))
* Chore use get or skip in other tests ([#2570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2570)) ([2829b90](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2829b9053e4123805b8f8e58b3223c6603d0f616))
* modify provider context to pass snowflake client instead of db connection ([#2577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2577)) ([e7fd4ef](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7fd4efde121d82eb278f32938c0aa612c13e55b))
* Speed up tests ([#2580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2580)) ([f003715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f003715b32c73e24b9c62c79c5436f058b127e9c))


### 🐛 **Bug fixes:**

* file_format validate() for field_optionally_enclosed_by ([#2575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2575)) ([5da5d93](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5da5d9369e79809ae9dc94830ea2c87b574d687c))
* Fix import for account parameter ([#2594](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2594)) ([cac884c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cac884c9f694aeb21476ed5541e579e48d226aea)), closes [#2573](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2573)
* support snowflake application as database ([#2596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2596)) ([b9a4a19](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9a4a19d59866a7f82d7ff50a884d4ac1de8ce04))
* tag identifiers problems ([#2534](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2534)) ([3c300e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3c300e1f7a0ce7159eb3722f4f04415fdf051757))

## [0.87.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.86.0...v0.87.0) (2024-02-28)


### 🎉 **What's new:**

* Add network rule to the sdk ([#2526](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2526)) ([b379565](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b379565de39a00d3863f22351806f87cdd920066))
* supports collation of table column ([#2496](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2496)) ([56771f1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56771f188841a37202827eee73ab412a633ef3e4))


### 🔧 **Misc**

* Clean up environment variables in tests and on CI ([#2543](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2543)) ([9a10cb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a10cb1211940dfd7339d842a91d3a8f790e947f))
* replace warning in new grant resources with info log ([#2521](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2521)) ([c3014b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c3014b93c24146f3abbd9583a1eca5c94081b80f))


### 🐛 **Bug fixes:**

* data retention days follow up ([#2566](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2566)) ([7aba384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aba384168ad1168a5e0ed32e4c91280d69d9a22))
* data retention time parameters ([#2502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2502)) ([76abf21](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/76abf217c519ca116ad8ab10d9947eaca3eb015b))
* data retention time parameters follow-up ([#2530](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2530)) ([5544544](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5544544dff8c25c83f23f7759341992a9b282fda))
* Demote warning to info and set volatility for procedures and functions ([#2567](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2567)) ([abaad7c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/abaad7cdccd058eab689e097bf0b5faa16c8f820)), closes [#2536](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2536)
* Fix ACCOUNT PARAMETERS option failover group resource ([#2522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2522)) ([61883f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/61883f3ee23f5aaeba065a75d0375fd8c69e2b4a)), closes [#2517](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2517)
* Fix failover group alter syntax and suppression for pipe statement ([#2562](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2562)) ([24d76c3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/24d76c383136e9b25d3d969107f2a770aad1d253))
* Fix few tests ([#2515](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2515)) ([a523a6b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a523a6b30f73165a2b51f7c88b8d6701ea19f0f9))
* Fix provider config hierarchy ([#2551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2551)) ([677a12b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/677a12b6587ce72fc3eb18de44a8f99a403d9502))
* Fix query_results in unsafe_execute resource ([#2512](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2512)) ([94ca158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/94ca158966785f4eb50d224a4231231f5aa8f24c)), closes [#2491](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2491)
* Fix replication for database resource ([#2524](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2524)) ([767fbce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/767fbce50fa3412a7b964fefae114a18d9fdaae2)), closes [#2021](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2021)
* Fix show by id for external functions ([#2531](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2531)) ([d910a84](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d910a84a7d829d14e414757a1732c934937e559e)), closes [#2528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2528)
* Fix various small problems ([#2552](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2552)) ([f558ce6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f558ce65e501e77fa95b26c56a1e261e2811e4c4))
* Granting database roles ([#2511](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2511)) ([dc27d64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dc27d644c15a0d2321d9f374f73a9df239c6a972)), closes [#2402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2402)
* grants on external volumes ([#2538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2538)) ([1de9a29](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1de9a296e43a2e0262a3a77346d8d5c6d1d9500f))
* Handle old reference for table_id in table constraint resource ([#2558](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2558)) ([d1e8912](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d1e891267c055b0a00246d5098787f3d21e807cb)), closes [#2535](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2535)
* loosen identifier field validation for account object identifiers ([#2564](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2564)) ([a5ed8cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a5ed8cdafab7cbdf58a27e155cf6b49db286f944))

## [0.86.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.85.0...v0.86.0) (2024-02-15)


### 🎉 **What's new:**

* add refresh_mode and initialize to dynamic tables ([#2437](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2437)) ([d301b20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d301b200b45398a1bf1e6e4127983912ba1a52a8))
* add resource snowflake_user_password_policy_attachment ([#2162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2162)) ([#2307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2307)) ([93af462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/93af462d66e4d8e1606c350a3a9c9a7caea4eb8e))
* create a workaround for granting privileges on all pipes ([#2477](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2477)) ([64f2346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/64f2346674086c01f22e20e0f8d96572e364ff54))
* Handle IMPORTED PRIVILEGES privileges in privilege-to-role granting resources ([#2471](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2471)) ([eb20051](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/eb20051ee1d2b6f23e851fc7b9b87c3aac9ef589))
* use external functions ([#2454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2454)) ([417d473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/417d473774e9abb8fda033f5082c80c917442ba5))
* use funcs from sdk ([#2462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2462)) ([a5f969c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a5f969caafb606c6ee56442d3f191e042cbf6651))
* use sdk for procedures ([#2450](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2450)) ([94ac78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/94ac78a4272f35c715349174d78801e691eed566))
* Use sdk in table constraint resource ([#2466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2466)) ([d685603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6856036ffc854631b1f44ee851dac56dc81ae25))
* Use tables from SDK ([#2453](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2453)) ([fdb4f88](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdb4f88ade09c9a1d63029b1937f1ef87528db8d))


### 🔧 **Misc**

* Add migration notes to the docs and change jira integration ([#2497](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2497)) ([b17f1af](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b17f1af92aaaeda9f471b680cd4d2df1f4c1dc55))
* Change email and issue reporter  ([#2470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2470)) ([5865655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58656556a3001259318aac39d5484d8ff7ed4f0e))
* Grants migration guide ([#2455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2455)) ([62c70fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62c70fd8d1f1ec97a50806cdd2e13af6518fe28e))
* Remove unused old implementation from snowflake pkg ([#2458](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2458)) ([2d0e508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d0e508412b9f3c429f2014afa67225624aec02d))
* update password policy attachment ([#2485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2485)) ([6ec9ff7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ec9ff7128d7b4bb2cf87f6cdbc57216c966b0aa))


### 🐛 **Bug fixes:**

* allow DT warehouse to be updated in-place ([#2439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2439)) ([d565af1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d565af19ca3829539e296320a7a9037cb2547ebc))
* correct test dependencies ([#2493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2493)) ([dfb247f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dfb247f1b2fc91b3ff49c7bca8e9f909711fdf14))
* FileFormat not detecting changes correctly ([#2436](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2436)) ([018bb74](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018bb74429d95a2554b3418061fdf3eacb3c735a))
* Fix few smaller issues ([#2507](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2507)) ([a836871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a836871cdd23a48804d59c424845dd54e853e081))
* Fix functions and small other fixes ([#2503](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2503)) ([0d4aba4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d4aba46a45e148e1914fc93730504bbb27f1600)), closes [#2490](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2490)
* Fix tag tests in view and in materialized view ([#2457](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2457)) ([2de942a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2de942acb308505256b2a7830c2af48106c18c4f))
* Fix task related issues ([#2479](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2479)) ([0385650](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0385650f9ea7d1708b65139f2eeeacc373d5c9f8))
* Fix tests that base on default data retention ([#2465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2465)) ([682e28c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/682e28c9dcac9023339c792912b5aa7cd83f42eb))
* grant privileges to share test terraform dependencies ([#2473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2473)) ([ede8d95](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ede8d9547ed5d4e6f2b39d8e41f1b09981d7d9bf))
* parameter issues ([#2463](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2463)) ([7ee4986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ee49869b799e339a672ca675ffcd704a96d2c67))
* parse dynamic table query from DDL ([#2438](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2438)) ([d76815c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d76815cf330399b7da303cf8502501f62537c801))
* Remove title and body temporarily from jira integration ([#2499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2499)) ([672c97d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/672c97db56f56ecf578cb26b53b30dc106d799a8))
* SHOW GRANTS mapping for share data type ([#2508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2508)) ([feb4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/feb4d445d27579bf88a8ee521ad380e0feac7c1f))
* user error handling ([#2486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2486)) ([dfa52b2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dfa52b2c46e5596184969da96aeb30471c74b492))

## [0.85.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.84.1...v0.85.0) (2024-02-01)


### 🎉 **What's new:**

* Add API integration to the SDK ([#2409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2409)) ([23acda5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23acda5dba9c8378f3b5631446d380a27cf1732c))
* add application to sdk ([#2350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2350)) ([de97ad8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de97ad84db925b62ab10046e0893a5c285a26d67))
* add external funcs to sdk ([#2440](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2440)) ([c8cf09b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c8cf09b2af605dc373a138e6ca6863b5546303d5))
* Add grant privileges to share resource ([#2447](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2447)) ([d8241a5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d8241a5cc76ea7b929abdada81cf6929b5f6ad9e))
* Add materialized view to the SDK ([#2403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2403)) ([a5ce699](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a5ce69920328cce899260249d319ff7726ae3911))
* Add notification integration to the SDK ([#2412](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2412)) ([d84240c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d84240cda369ed9106c7cb3e3eedf85b8d1fa944))
* add sequences to sdk ([#2351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2351)) ([d2e5ffd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2e5ffd5405f10ff30c5ad9f7cd58bd54a5cc028))
* add snowflake grant privileges to account role ([#2365](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2365)) ([e3d086e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e3d086eddc05e0d4963234f82e09e174a018bb08))
* add streamlits to sdk ([#2400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2400)) ([129d24c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/129d24c00fa244d1401cb2169b5b7fb0ba6c465c))
* add-call-with to sdk ([#2337](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2337)) ([ebcd1bc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ebcd1bc40d554abe6863b67d2ab76f2d992dfb32))
* stages migration follow-up ([#2372](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2372)) ([3939dbe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3939dbe2f9189968c087a883ed97dd3b7350787f))
* Use API integration from SDK ([#2429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2429)) ([1ccc864](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1ccc8641106a3ceb4de813ce7c0e5077ead5272e))
* Use managed account from the SDK ([#2420](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2420)) ([3aaa080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3aaa08071a14f820e08751cc7b1e8bef5db16e30))
* Use materialized views and views from SDK ([#2448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2448)) ([dc66d02](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dc66d02304a99a7cb152e91a8e942587cab7e60f))
* Use notification integration from sdk ([#2445](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2445)) ([e8915cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e8915ccb99eeec1f0ac5777fe80be7ef443d8f5c))
* use roles from the SDK ([#2405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2405)) ([c645b4d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c645b4d0e2036d932766480e9c1e0334ef79c16e))
* Use row access policy from SDK ([#2428](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2428)) ([119af5e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/119af5ea74bb219ae822962096e6220ed00f5910))
* Use SDK in the storage integration ([#2380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2380)) ([ce0741c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce0741ce226be9464407b549e90cb179b0fe5880))
* use sequence from sdk and add ordering attr ([#2419](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2419)) ([973b8f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/973b8f76a8ed1540bfd948ba8cb57c212c0d4abc)), closes [#2387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2387)
* Use stage from sdk ([#2427](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2427)) ([c17effd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c17effd16ccd77ba4c5d45f43dcc53a9f11601c6))


### 🔧 **Misc**

* add missing deprecation message ([#2451](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2451)) ([77de569](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/77de5694f73e5ad1443bb99407d2e8aec9a87320))


### 🐛 **Bug fixes:**

* account role test ([#2422](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2422)) ([c1b47d1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1b47d1ade4b198b5bf14dc32162d34797a3b344))
* Adjust tests after Snowflake behavior change ([#2404](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2404)) ([8c03ffb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c03ffb0430445c903168da9706e1ce2630675da))
* app-pkg unset ([#2399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2399)) ([fedb1df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fedb1df2a731d139d68d2284bf3be47fcc4d0115))
* Fix some bugs ([#2421](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2421)) ([dec7cd9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dec7cd9e199ac8658f5c939f811686ba9f5e2e21)), closes [#2358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2358) [#2369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2369) [#2329](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2329)
* snowflake_grant_privileges_to_role read ([#2424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2424)) ([5385cec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5385cec3e5c03d2dbff762b63523bdddee8632d3))

## [0.85.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.84.1...v0.85.0) (2024-01-22)


### ⚠ BREAKING CHANGES

* mark snowflake_user.login_name as non sensitive

### 🎉 **What's new:**

* account password policy attachment ([#1824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1824)) ([f408828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f408828fd023c2207ec41f702cec7dae524b1e93))
* Add "CREATE DYNAMIC TABLE" to schema_grant ([#2144](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2144)) ([6f026f6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f026f64e6e24638df2b9d4110362836a9071011))
* add app packages ([#2323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2323)) ([ca030fc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca030fce209ccc1e5e924b3bb1fd3c680c930471))
* add authenticator ([#2109](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2109)) ([4f3a551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4f3a5519484b0aab91ff5fa08f37a8cf512d1ec0))
* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* Add COPY GRANTS arg to views resources. ([#1668](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1668)) ([7225d93](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7225d93ef3e50c6810c0dd57cfd7079e882d443f))
* Add create streamlit privilege to the SDK ([#2303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2303)) ([be01d5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be01d5fdab4f2d31db9c4c849b349e657c0352c8))
* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))
* Add database role to SDK ([#2009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2009)) ([f5efc09](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f5efc09ea60bd2d66c65c9e07cb84321f95531f0))
* add dynamic tables to sdk ([#2074](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2074)) ([d1dfb05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d1dfb05fbb3bcc59cc2622b6b2d02ebadf1cf33f))
* add event tables to sdk ([#2215](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2215)) ([66cc80a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/66cc80a8cac24f4b7967986032b7da9e20bd4eab))
* Add external functions translators ([#1735](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1735)) ([1f67286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1f672862adb29a658f5e81e940f9afb994347f2f))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* add failover groups to sdk v2, and add data source for failover groups ([#1825](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1825)) ([44e8c06](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44e8c06ba4c665c81f0b909dbd3df90c4925e179))
* add functions to sdk ([#2205](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2205)) ([e542b67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e542b67f761850e812114ce593fa5f6deca941cb))
* Add gcp_pubsub_topic_name parameter for gcp notification integration ([#1687](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1687)) ([a30d0cb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a30d0cb756a2281a4d880af9e32651c04409028e))
* add GRANT ... ON ALL TABLES IN ... ([#1626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1626)) ([505a5f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/505a5f35d9ea8388ca33e5117c545408243298ae))
* Add grant ownership to SDK ([#2064](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2064)) ([f85ec8b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f85ec8baa2f9aaeead4f619dccfa3d38880a16d7))
* Add grant privileges to database role to SDK ([#2023](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2023)) ([717289f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/717289f71fd4a08f44d4f20f6e16dc4dded77803))
* add grant role resource ([#2304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2304)) ([ba91e25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba91e25be15197067e934eb867fb6837b4354e4a))
* add grant_database_role resource ([#2301](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2301)) ([2e7651f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2e7651f35b625bc11f17f42443d062118b6354ba))
* add in more functionality for UpdateResourceMonitor  ([#1456](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1456)) ([2df570f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2df570f0c3271534a37b0cb61b7f4b491081baf7))
* Add Manage Warehouses Account Grant ([#2017](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2017)) ([89c7148](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/89c7148c11378af9e42ea32bdad3e5a5c465d39c))
* Add managed account to the SDK ([#2357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2357)) ([f968db1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f968db11a9b41ffa594170415c4edd9376f574e7))
* add mfa auth ([#2077](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2077)) ([922358a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/922358a43c5383ee5840bf2971ecd27d96f86573))
* Add missing database role operations to SDK ([#2014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2014)) ([d2ea67d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2ea67d7fee00b15e1222fe37efe8e7a1cddecb5))
* Add missing parameters to password policy resource ([#2231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2231)) ([c189782](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c189782bb7d117e27979c21dfea60ba733e996df))
* add missing PrivateLink URLs to datasource ([#1603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1603)) ([78782b1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78782b1b471b7fbd434de1803cd687f6866cada7))
* add new account resource ([#1492](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1492)) ([b1473ba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1473ba158946d81bc4eac095c40c8d0446cf2ed))
* add ON STAGE support for Stream resource ([#1413](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1413)) ([447febf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/447febfef46ef89570108d3447998d6d379b7be7))
* add on_account to session and object params ([#1685](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1685)) ([1329430](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/13294304c7626c9d428682986669d2e97ab2c23b))
* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* add parse_header option to file format resource ([#2132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2132)) ([1e6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e6e54f828efa60edd258b316709fc4dfd370f93))
* Add Pipes to SDK ([#1968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1968)) ([69a543f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69a543fd729b64bdd8964dc34626dee83b3f96a7))
* add PREVENT_LOAD_FROM_INLINE_URL ([#1612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1612)) ([4945a3a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4945a3ae62d887dae6332742edcde715751459b5))
* add procedures to sdkv2 ([#2202](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2202)) ([6b563ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b563acef0f702fc6a366219afac602fa106129c))
* Add row access policy to the SDK ([#2363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2363)) ([14a3e5b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14a3e5b07a97cb18ba99f314990191241d4b15c2))
* Add session policy to SDK ([#2088](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2088)) ([038241c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/038241c00a158b389df8034864a52a252fcc41bf))
* add shares to sdk v2 ([#1813](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1813)) ([a814841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a814841011f08857d8d37691fa5ff01cd9412176))
* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* Add support for `packages`, `imports`, `handler` and `runtimeVersion` to `snowflake_procedure` resource ([#1516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1516)) ([a88f3ad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a88f3ada75dad18b7b4b9024f664de8d687f54e0))
* Add support for is_secure to snowflake_function resource ([#1575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1575)) ([c41b6a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c41b6a35271f12c97f5a4da947eb433013f2aaf2))
* Add tables to the SDK ([#2042](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2042)) ([c1700de](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1700de2a062852da7cb5e3cf3277cc19f6466d6))
* Add task clone to sdk ([#2105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2105)) ([acddb2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/acddb2cd6bfb1a7ffaf6dbb3c8349f7bc550c124))
* Add task to SDK ([#2099](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2099)) ([d52f334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d52f3347f091f0edff5e6daded1120542f1e9bd1))
* Add unsafe_execute resource ([#2225](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2225)) ([196134c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/196134cbf91996eabc50bdc586a657fe7ac71900))
* Add views to the SDK ([#2171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2171)) ([ed079d3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed079d3d06dc3af083da04ca18314c8e7b07308e))
* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* added on_all grants for view, stage, schema and materialized view resource ([#1686](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1686)) ([f27a9e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f27a9e454102709e1134d796f1594b4932670ae8))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* adding a subset of missing parameters ([#1828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1828)) ([6d1a572](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6d1a572623438dd96139edbd698ce5974fa3df61))
* adding CREATE STREAMLIT as privalege ([#1771](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1771)) ([7b112cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b112cc94822807d3c25c78ec812ec2c2c66ea16))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* addition of the on_all option for all remaining resources ([#1742](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1742)) ([e3a7710](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e3a771000fa527fc8420b892450031ba39394fcf))
* alert resource ([#1653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1653)) ([4a84eb9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a84eb91af6201330a2ec60de306db5ed3b5ab38))
* alerts datasource ([#1727](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1727)) ([a1d0d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1d0d2eeb9a8a03d4c347b9f153182e66c237a4f))
* alerts sdk v2 ([#1871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1871)) ([fcf0380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fcf03800cf92e1acff7f89e59968914f1bcad4f4))
* allow all priv on grants ([#1786](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1786)) ([aa3f873](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aa3f873b4669c5942c10681151d6afd4ea03623f))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* Convert alert datasource to new sdk ([#2020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2020)) ([2d0eaeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d0eaebe8ea1605251dd24f86a2f83e6e3be67c8))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))
* data source shares ([#1651](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1651)) ([5dd447f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dd447fc79c153f025e420052a23b90751cdaa9f))
* database role resource ([#1654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1654)) ([9f4a516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f4a5165d6b7e1fa63ad06450ffc55689be66356))
* database roles datasource ([#1726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1726)) ([f1cb7f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1cb7f3b368c0a7fa6a0a08784d53b0e92efda81))
* databases sdk v2 ([#1840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1840)) ([410b493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/410b493ff75188b70cee91db56c60538d6823aac))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* dynamic table resource + data source ([#2104](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2104)) ([59c5de6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/59c5de6becf4b16542da4db59e839f7ca5625e50))
* email notification integration ([#1720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1720)) ([5d21fe1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d21fe19f0da9ed4939e9bd6dc66237685c5a327))
* external_oauth_integration.scope_mapping_attribute ([#1722](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1722)) ([66d88bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/66d88bd29b6aa78dbc2acd11cf729fa9380b011d))
* Generate constructor and builder methods for DTOs ([#2001](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2001)) ([79d9c9c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79d9c9ce77e85dc9bd0889da30155cef0a74c293))
* Go Generator + Network Policy migration to the new SDK ([#2061](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2061)) ([231b081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/231b081a34462392954c18b0f3452b1fed7d016f))
* grant privileges to database role resource ([#2306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2306)) ([0311cf8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0311cf8554f0fbd202b489a54c9428f55c52a490))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))
* identifier validation ([#2269](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2269)) ([9687972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9687972fc673bb7b768f9b4100a6d2b67fc46e48))
* Initialize SDK generator ([#2033](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2033)) ([96b47e5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/96b47e5fdbb6175a7a24d0518b975279927c22de))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))
* **integration:** add google api integration ([#1589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1589)) ([56909cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56909cdc18245f38b0f58bceaf2aa9cbc013d212))
* Introduce simple arch tests ([#2210](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2210)) ([c60db80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c60db80f44d949258f0a692baafdc22b886c3010))
* mark PII as sensitive in account and user resources. ([#1678](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1678)) ([caa461f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/caa461fea56df4f33f73e586b537802c01f1eb5d))
* mark snowflake_user.login_name as non sensitive ([71c0a9c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/71c0a9c2897e07b2251680d619cbbcaeba5779b2))
* masking policy application resource ([#1739](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1739)) ([ce80f57](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce80f57225d6019765a50da6eb74619e401db622))
* masking policy in v2 sdk ([#1777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1777)) ([6978c42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6978c42b304bd8c5429c88bcd7d6ed20ac3fd98a))
* Migrate application role to new sdk ([#2149](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2149)) ([7abb4db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7abb4dbb645f0eb0fcb1d39414b1ed0c322916c9))
* Migrate external tables to new sdk ([#2006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2006)) ([5af17cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5af17cfbe4b6570c61b4ac77ab7564fd1f0529b7))
* Migrate roles to the new sdk ([#2007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2007)) ([5b996e6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b996e6f4b3ea28bae0cd13e75bd112d9af09c68))
* Migrate stage to the new sdk ([#2163](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2163)) ([1d08c46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d08c46c470cda2d6b7c7ed3d599d2c1a57b837b))
* Migrate storage integration to the new sdk ([#2339](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2339)) ([d970a56](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d970a567741edf1ce72e566f9e76b3858ece6288))
* Migrate streams to the new sdk ([#2113](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2113)) ([521fde5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/521fde58b7836a21614d2ae31a7060352a68b465))
* Move integration tests to separate package ([#2111](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2111)) ([2755589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2755589bc8c7ecbda3115dc574339c40bfd5b096))
* new grant resource "snowflake_grant_privileges_to_role"  ([#1929](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1929)) ([e241c22](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e241c223b0c5d159efc803b9fcd4671c8325ae00))
* Optional switch for instrumentedsql ([#2261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2261)) ([9934a59](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9934a59c5ee31dc292961806440427000ea41ea7))
* Parameters sdk v2 ([#1914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1914)) ([1d15355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1535572e1b82f3cb2308e3c19e632eb34fe24b))
* Password policy resource ([#1702](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1702)) ([7ee293b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ee293b52262170a1d53811f9beb6fc3772ae913))
* Poc custom error type ([#2052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2052)) ([b86c4c3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b86c4c34d05f8b982fb6218a3a3a7500a23abf72))
* Provider config refactor ([#2126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2126)) ([9858fea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9858feaa2fb382b57e3c89aae044f09ce2dcc5a3))
* provider config support insecureMode ([#1703](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1703)) ([e269925](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e269925e0f54ca9bbd742e5c60eaebad5f9b1f04))
* resource monitor sdk v2 ([#1892](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1892)) ([707d723](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/707d7233bd867d60e13c8785d1c0adef558a5853))
* roles support numbers ([#1585](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1585)) ([d72dee8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d72dee82d0e0a4d8b484e5b204e156a13117cb76))
* Schemas sdk v2 ([#1975](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1975)) ([289ad8a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/289ad8a8d61cf217d4548a65cf4843e416baa1da))
* set query tag for terraform sessions  ([#1826](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1826)) ([6629583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6629583bc288e5c23cf02d8a53597adbadf55fce))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))
* support custom ownership transfer for grants ([#1743](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1743)) ([eaa6e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/eaa6e01820cb04ffa3c647ecd893e479af8e35a1))
* support drop for accounts and fix read race condition ([#2065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2065)) ([5412252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/541225258ff19de864c8d25d05a8c1b4a7941813))
* support is_org_admin column in SHOW ORGANIZATION ACCOUNTS ([#1673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1673)) ([263c521](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/263c5215dedf70195c0c2c7d8e1505e4b9c0828c))
* support object parameters on account level ([#1583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1583)) ([fb24802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb2480214c8ac4e61fffe3a8e3448597462ad9a1))
* tags sdk refactoring ([#2079](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2079)) ([7013f83](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7013f835864139ef40c2521e6af87589e909e201))
* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))
* Use database role from SDK ([#2012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2012)) ([294075a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/294075a4d145dc08071c11f067afaecb78fe8ef7))
* Use external tables from SDK ([#2228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2228)) ([6941023](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6941023f95230a14f92dd099c1a6375129ee4bfb))
* Use network policy from sdk ([#2087](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2087)) ([50f2935](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/50f2935dc12fc2479c2919a5666d5573a8b737ff))
* Use pipes from SDK ([#2003](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2003)) ([079d47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/079d47d22af22edba6c6c499409b264c4c5f5945))
* use SDK in schema resource and datasource ([#2082](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2082)) ([f7d0d97](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f7d0d9716c4da9b00e5171a368a03a9cd52c31ae))
* Use streams from the new SDK in resource / datasource ([#2129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2129)) ([5c633be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5c633be461fd373d412b02b108e64b6cfc4eb856))
* Use task from SDK in resource and data source ([#2140](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2140)) ([de23f2b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de23f2ba939eb368d9734217e1bb2d4ebc75eef4))
* Use tasks from the SDK followup ([#2153](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2153)) ([82c3c13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82c3c13b6166168e470d7cb9b2982a8979275f17))
* Users sdk v2 ([#1945](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1945)) ([d644b63](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d644b63e673223b67f7092563a0bfc8c826fc38b))


### 🔧 **Misc**

* 0.84 release ([#2374](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2374)) ([18638fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18638fea4c498158c8840a3ccfe03bdaac942923))
* Account SDK ([#1822](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1822)) ([4c22b89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c22b89e82a608375fd78993134cf3a0b45b985b))
* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))
* Add migration guide ([#2142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2142)) ([ee0f6af](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee0f6af54dbd269f8bfa2c3d73a396d98d10a6ab))
* Add missing issues to existing TODOs and add missing ones ([#2354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2354)) ([c3bc66c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c3bc66cffa4b87c49e027cc685df1abd8a6ebfe6))
* Add short scripts used to fetch all currently opened issues ([#2288](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2288)) ([0b5ce4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b5ce4eff21562406a473ba8d796804efb1bd94f))
* Adjust integration tests after moving to separate package ([#2115](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2115)) ([3f528a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f528a87f4c0b3bc95a0dfb35d93d22251b5112e))
* Allow setting gosnowflake logging level from environment variable ([#2285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2285)) ([843e8fc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/843e8fc8f8879fdb5b1148155e7e64174571e975))
* Bump dependencies ([#2252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2252)) ([581d75c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/581d75c566ea758c5f888dedb101e6318d7c0cc3))
* Bump dependencies and fix linter complaints ([#2300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2300)) ([124e862](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/124e862b406526e08d865e0f70414926e50f4d4d))
* Bump go to 1.21 ([#2267](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2267)) ([6b852c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b852c2ba6b1348ba83ec2e2190e8175371524be))
* Change return type of show ([#2045](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2045)) ([21f069a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/21f069a3232a0039792f30445e25b5e88f09ab25))
* Cleanup comments ([#2092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2092)) ([3a06a66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3a06a66f27d9b03d037affdda5b7dbb3bd81d5d7))
* **deps:** bump actions/setup-go from 3 to 4 ([#1634](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1634)) ([3f128c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f128c1ba887c377b7bd5f3d508d7b81676fdf90))
* **deps:** bump github.com/hashicorp/terraform-plugin-docs from 0.14.1 to 0.16.0 ([#1931](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1931)) ([07c4f14](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/07c4f149bf08ea9fecd28aa12a47a9fae3e0523b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1639)) ([330777e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/330777eb0ad93acede6ffef9d7571c0989540657))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#2090](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2090)) ([24b6313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/24b6313a7f28731ad4c10dd80c32e7b9cf397c5d))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1638)) ([107bb4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/107bb4abfb5de896acc1f224afae77b8100ffc02))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump github.com/stretchr/testify from 1.8.2 to 1.8.4 ([#1842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1842)) ([9b0825d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9b0825db520c244cea962830d2139fb50186f23c))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/crypto from 0.13.0 to 0.14.0 ([#2100](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2100)) ([060750f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/060750f668b986d690619cd2d915a7a945350ec5))
* **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1454)) ([ed6bfe0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed6bfe07930e5703036ada816845176d46f5623c))
* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#1528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1528)) ([8a011e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a011e0b1920833c77eb7832f821a4bd52176657))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#1551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1551)) ([35de62f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/35de62f5b722c1dc6eaf2f39f6699935f67557cd))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* **deps:** bump golang.org/x/tools from 0.13.0 to 0.14.0 ([#2101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2101)) ([45a14d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/45a14d96d5fd6b051a73ced2fb03cab8659ded84))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump golang.org/x/tools from 0.4.0 to 0.5.0 ([#1455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1455)) ([ff01970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff019702fdc1ef810bb94533489b89a956f09ef4))
* **deps:** bump golang.org/x/tools from 0.7.0 to 0.13.0 ([#2089](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2089)) ([0ace968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ace9689c6b29382b1d93d9a3542eacc0d7c4270))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))
* **deps:** bump goreleaser/goreleaser-action from 4 to 5 ([#2108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2108)) ([c4b3e82](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b3e82a9ebf8f036282a96f320cc08761e229b2))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))
* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))
* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))
* Extract reviewdog job from integration workflow ([#2027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2027)) ([243fc28](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/243fc285811de787328a7f3ca0d3a84193976b85))
* fix parameters datasource acceptance test ([#1749](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1749)) ([338a19d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/338a19de373525087372cc319aa065f2cddbb724))
* Integration tests setup db once ([#2118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2118)) ([f533368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f5333684d5e46877ec4da6dc708ce2b09b5694af))
* **main:** release 0.34.0 ([#1332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1332)) ([7037952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7037952180309441ac865eed0bc2a44a714b484d))
* **main:** release 0.34.0 ([#1436](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1436)) ([7358fdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7358fdd94a3b106a13dd7000b3c6a8f1272cf233))
* **main:** release 0.34.0 ([#1662](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1662)) ([129e4dd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/129e4ddbc7424306d75298486c1084a27f2a1807))
* **main:** release 0.47.0 ([#1259](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1259)) ([887297f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/887297fc5670b180f3d7158d3092ad035fb473e9))
* **main:** release 0.48.0 ([#1284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1284)) ([cf6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf6e54f720dd852c1663a4e9ff8a74054f51325b))
* **main:** release 0.49.0 ([#1303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1303)) ([fb90556](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb90556c324ffc14b6e90adbdf9a06705af8e7e9))
* **main:** release 0.49.1 ([#1319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1319)) ([431b8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/431b8b7818cd7eccb3dafb11612f72ce8e73b58f))
* **main:** release 0.49.2 ([#1323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1323)) ([c19f307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c19f3070b8aa063c96e1e30a1e6d754b7070d296))
* **main:** release 0.49.3 ([#1327](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1327)) ([102ed1d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/102ed1de7f4fca659869fc0485b42843b394d7e9))
* **main:** release 0.50.0 ([#1344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1344)) ([a860a76](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a860a7623b9e22433ece8cede537c187a45b4bc2))
* **main:** release 0.51.0 ([#1348](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1348)) ([2b273f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b273f7e3baaf855ed6e02a7779022f38ade6745))
* **main:** release 0.52.0 ([#1363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1363)) ([e122715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1227159be50bf26841acead8730dad516a96ebc))
* **main:** release 0.53.0 ([#1401](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1401)) ([80488da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/80488dae4e16f5c55f913449fc729fbd6e1fd6d2))
* **main:** release 0.53.1 ([#1406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1406)) ([8f5ac41](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5ac41265bc08256630b2d95fa8845249098310))
* **main:** release 0.54.0 ([#1431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1431)) ([6b6b55d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b6b55d88a875f30395f2bd3250a2af1b99f9205))
* **main:** release 0.55.0 ([#1449](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1449)) ([1a00052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1a0005296689ad3ae45e5fd92b06e25ed16232de))
* **main:** release 0.55.1 ([#1469](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1469)) ([509ce3f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/509ce3f168d977de71758518e99ce0e38ab9f875))
* **main:** release 0.56.0 ([#1493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1493)) ([9a5fc2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a5fc2c0fdf993285bae42efb83b3384085540a0))
* **main:** release 0.56.1 ([#1504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1504)) ([00fc00c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/00fc00c46f22984e02ed10acdc8041cfc79b507d))
* **main:** release 0.56.2 ([#1505](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1505)) ([f950dac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f950dac5d13516075c416f6abc6d7667474a36a8))
* **main:** release 0.56.3 ([#1511](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1511)) ([9c69643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c69643a31d91d0f3d249f7aea3beeefc53880ae))
* **main:** release 0.56.4 ([#1519](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1519)) ([d0384b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d0384b6d3bfc1bc358f39e58f136c1acef452456))
* **main:** release 0.56.5 ([#1555](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1555)) ([41663ee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/41663ee5900206a03c62e046bfb9659092197bd5))
* **main:** release 0.57.0 ([#1570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1570)) ([44b96cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44b96cf67813f45feb67da4367936748bc04391f))
* **main:** release 0.58.0 ([#1587](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1587)) ([6b20b8d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b20b8d848620a7e9796ae230f6f87300e3fc50c))
* **main:** release 0.58.1 ([#1616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1616)) ([4780ba0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4780ba08b1bdf15785be63ec8dd488a03ddfe378))
* **main:** release 0.58.2 ([#1621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1621)) ([1c34ac1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c34ac157bc064d5d6fe5297231ce87eccbcc298))
* **main:** release 0.59.0 ([#1622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1622)) ([afb18aa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afb18aa8ed3c3f80630bc2f824bb756ddb5eda86))
* **main:** release 0.60.0 ([#1641](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1641)) ([ab4d49f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab4d49f259db99c2c0c6131143c55ca11d2a6610))
* **main:** release 0.60.0 ([#1665](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1665)) ([ea23020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea23020801ea4d43f055f2b443400385d96a135b))
* **main:** release 0.60.0 ([#1667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1667)) ([9d3e40f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d3e40fab05023aff16795266ec8a30761560c26))
* **main:** release 0.60.0 ([#1670](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1670)) ([c61ed63](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61ed6379895f83cebdc8e1c2a82ec0bf97f42d6))
* **main:** release 0.60.1 ([#1649](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1649)) ([56a9b2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56a9b2e5747bffb2456ad2a556e226e8450c242e))
* **main:** release 0.61.0 ([#1655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1655)) ([2fbe15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2fbe15a65a64adb8604d301e9a6d11632b6e3a44))
* **main:** release 0.61.0 ([#1676](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1676)) ([bce459e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bce459eee4e2c247c7167ce48447adf4de5a20df))
* **main:** release 0.62.0 ([#1689](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1689)) ([736ce1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/736ce1e7c6cf331bbbc1e260e71f99f96e8fc5b5))
* **main:** release 0.63.0 ([#1750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1750)) ([f5c65ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f5c65ac1009b88723009171fe9da416c5595c292))
* **main:** release 0.64.0 ([#1766](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1766)) ([1037f98](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1037f9833131ca6d5a0e75adcafe9b99f5f86173))
* **main:** release 0.65.0 ([#1789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1789)) ([883ee1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/883ee1cc34a728dea01c96b2630bed739d2da06c))
* **main:** release 0.66.0 ([#1839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1839)) ([6e5d643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e5d6434f7e576754c2f67f69a4af432ad885486))
* **main:** release 0.66.1 ([#1854](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1854)) ([30d0017](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30d0017608f112b959b142e56bcf35d4ee23ff8b))
* **main:** release 0.66.2 ([#1876](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1876)) ([6aa8fa1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6aa8fa1b3c0ffb1680d4e7a69dd7a219632106c9))
* **main:** release 0.67.0 ([#1880](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1880)) ([43eda8b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/43eda8bea61d44fb872d6ccb2ec48107e0ff2d3c))
* **main:** release 0.68.0 ([#1912](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1912)) ([6d252fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6d252fa4e666b146976bdf47fb6865fe4afc8130))
* **main:** release 0.68.1 ([#1937](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1937)) ([067f3a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/067f3a2ede98ce449524f9a039ebf96581a54832))
* **main:** release 0.68.2 ([#1958](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1958)) ([605ad92](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/605ad92281e34dc93c86e6ba66a285e0d4d8a0f4))
* **main:** release 0.69.0 ([#1960](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1960)) ([c3edbeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c3edbebeda6552a4691081c0c962e3b15cd4535e))
* **main:** release 0.70.0 ([#1991](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1991)) ([6d1703c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6d1703cfb2cdb03c2ecae8464ac0542df1d0a2e2))
* **main:** release 0.70.1 ([#2026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2026)) ([43277a4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/43277a40e791a80db935f1a5241317eafb661b65))
* **main:** release 0.71.0 ([#2049](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2049)) ([538efa5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/538efa5ee9b0370e3090cd7fb867f3df4bbfd0ed))
* **main:** release 0.72.0 ([#2071](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2071)) ([2380e33](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2380e337b7324e95b107ed20519e9aa01fa04afa))
* **main:** release 0.73.0 ([#2103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2103)) ([8311925](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8311925f115dc0a642c2cd46762f8e82c2824dba))
* **main:** release 0.74.0 ([#2114](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2114)) ([16012ae](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16012ae222c611080efc3ffe2b820961bff5fc90))
* **main:** release 0.75.0 ([#2138](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2138)) ([d33a41c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d33a41cd0f477e8375113a546117ccd6e5dbb0eb))
* **main:** release 0.76.0 ([#2161](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2161)) ([a1cb41c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1cb41cc9d0a4c7b2ac009cbd02ee89ae579a36f))
* **main:** release 0.77.0 ([#2204](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2204)) ([8fbc5cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8fbc5cf3005f33c7cfb27815d6e4a1909ffbac80))
* **main:** release 0.78.0 ([#2230](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2230)) ([c7bc9fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c7bc9fa555f7f19296dacf2f169e472d5f5778b0))
* **main:** release 0.79.0 ([#2247](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2247)) ([962b1b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/962b1b6a10185c9e43aa8fd033b6a0f8caef7c3d))
* **main:** release 0.79.1 ([#2251](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2251)) ([411e85d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/411e85defedf12d2a9cc03a89244769365478aea))
* **main:** release 0.80.0 ([#2255](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2255)) ([18b0dcc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18b0dcc5d5523f1a89c87df87a678a6cd26e7ba6))
* **main:** release 0.81.0 ([#2268](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2268)) ([cec02b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cec02b43f087a52307fee4c4058fa6e7d45ca1c8))
* **main:** release 0.82.0 ([#2297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2297)) ([b1a740a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1a740ab29ff8be45de3ab416a01d0cad468d51d))
* **main:** release 0.83.0 ([#2305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2305)) ([ac5dd36](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ac5dd365d29bfff6a62c80564b7efbe8cf48ea05))
* **main:** release 0.83.1 ([#2344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2344)) ([9e6ae30](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e6ae307dcc1eaf31fc003616efeea1a3843da30))
* **main:** release 0.84.0 ([#2355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2355)) ([164f39b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/164f39bec613aa4812476e650493ff8d25825527))
* **main:** release 0.84.1 ([#2386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2386)) ([e2f8626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2f8626d210b425aa6a30373ab763ec79ab9f84b))
* **main:** release 0.84.1 ([#2393](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2393)) ([5194478](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5194478ab77fcc5b9935232f4f835766f4b80751))
* Make options naming consistent ([#2083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2083)) ([df84bd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/df84bd084482ef3f85db2b1f4a3517afe8df4ddb))
* Makefile cleanup ([#1995](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1995)) ([5c6fdbe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5c6fdbe48aadab9d56bfb7cf5d5ea30ca5a2b356))
* Migrate Warehouse resource + datasource to new SDK ([#1792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1792)) ([a14b994](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a14b994d62f73c272cde72651c3c8c18ac4213cf))
* Propose additional debug logging ([#2243](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2243)) ([a0984cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a0984cd355c2004cc245fbd4672967ba59cfb3c7))
* Remove lint ignore comments - test ([#2011](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2011)) ([8f81a42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f81a42544a668b03e0bb5d4d3199e6879d2f937))
* rename "db" symbol to "sql" and make input options more uniform in naming convention ([#1837](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1837)) ([244b5b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/244b5b45e22f66875f14be9bcb9af08ad41d37ca))
* Rename and reuse validatableOpts interface ([#2016](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2016)) ([4b42848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b428480a3f8fb9aa38da7a311f8a04f8d294882))
* Return multiple errors in existing validations ([#2122](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2122)) ([4d4bcdb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d4bcdbe841807da2fa08d534eaf846234934f7c))
* sdk v2 implemented for password policy ([#1752](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1752)) ([0cb1164](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0cb1164414607ce2e4ac7f6fb2da1c563b22da7b))
* Set up a single warehouse for the SDK integration tests ([#2141](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2141)) ([16022ef](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16022ef4171e7dccf2932ae6e8d451b51c93291c))
* Set up schema once for integration tests ([#2121](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2121)) ([5e3ebf4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5e3ebf4e809a4be163a30642f69ea1b9146f21a9))
* skip acc managed account test ([#2376](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2376)) ([dbf645d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dbf645d8ea508582bebd1aab485f54807c818bd2))
* Skip managed account SDK integration tests temporarily ([#2370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2370)) ([9b2f551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9b2f551fc9df34930805b3429fccf88be7e7f851))
* Split existing alter operations ([#2156](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2156)) ([dbb7c91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dbb7c9136c586490a0856cc07ae879be491c8150))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))
* update readme and fmt ([#2373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2373)) ([a648ca7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a648ca7ae1ca46d8bd05a11d7581ffd55d20ffd2))
* Use helper methods in old unit tests ([#2119](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2119)) ([0c44571](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0c44571a92b842ac97d46a38e5bcfc1066367177))


### 🐛 **Bug fixes:**

* 0.54  ([#1435](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1435)) ([4c9dd13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c9dd133574b08d8e67444b6c6b81aa87d9a2acf))
* 0.55 fix ([#1465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1465)) ([8cb3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8cb337048ec5c4a52245feb1b9556fd845d83278))
* 0.59 release fixes ([#1636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1636)) ([0a0256e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a0256ed3f0d56a6c7c22f810419636685094135))
* 0.60 misc bug fixes / linting ([#1643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1643)) ([53da853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53da853c213eec3afbdd2a47a8de3fba897c5d8a))
* 0.65 integration tests and minor linting fixes ([#1835](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1835)) ([8d7d663](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8d7d6636a67d272d5b87e02a178ae799b5b6d777))
* 1969  change default privilege to OPERATE ([#1970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1970)) ([e8721f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e8721f07dd27ce7a903c8c4d2b19fee77ce72cdc))
* add back in 386 build ([#2038](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2038)) ([5522504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55225041d7f4722adbd7b9445990f57c04a4adea))
* Add contributing section to readme ([#1560](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1560)) ([174355d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/174355d740e325ae05435bbbc22b8b335f94fc6f))
* add dynamic tables to grants ([#2059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2059)) ([3767de9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3767de97ec685acebd0e8e105878eb9129a0752c))
* add external table grant to snowflake_grant_privileges_to_role resource ([#1967](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1967)) ([8c84c4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c84c4a5f6eee7cf35b5ec2d77fac1c7fcee05df))
* add grant to sdk v2, better validation ([#1812](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1812)) ([76da3b8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/76da3b84d9d65493ae666fad6ea6d926b81add8f))
* Add missing Google Cloud computed property api_gcp_service_account to the api integration resource ([#1776](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1776)) ([1dd0672](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dd06724aa1add9b504c820583a2a752555c9bdb))
* add missing option ([#1924](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1924)) ([2991a16](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2991a16fe45491d77a8212a3ce0aa9566882de4f))
* Add missing parameters ([#2250](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2250)) ([4f4c4a4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4f4c4a4fcb91725a55726c03429745bbf95fd642))
* add missing reference for email_notification ([#1768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1768)) ([6feeb9c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6feeb9c09f56df0fb222bd0664508defb7c5af6e))
* add nill check for grant_helpers ([#1518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1518)) ([87689bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/87689bb5b60c73bfe3d741c3da6f4f544f16aa45))
* add permissions ([#1464](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1464)) ([e2d249a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d249ae1466e05dad2080f05123e0de66fabcf6))
* Add secondary account and fix tests ([#2324](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2324)) ([da6ca73](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da6ca733c7527c8918d1e1beb86d1641d94062ec))
* add test cases for update repl schedule on failover group ([#1578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1578)) ([ab638f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab638f0b9ba866d22c6f807743eb4eccad2530b8))
* added force_new option to role grant when the role_name has been changed ([#1591](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1591)) ([4ec3613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ec3613de43d70f40a5d29ce5517af53e8ef0a06))
* added privilege for `RESOLVE ALL` ([#1861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1861)) ([18cf7b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18cf7b09c9be927b5e1a625acdadb90abd048ff8))
* Adding double quotes to the column names ([#1729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1729)) ([791dd0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/791dd0b3a613d0ada95d4af7153a7b0a2bbee219))
* all-grants ([#1658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1658)) ([d5d59b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5d59b4e85cd2e97ea0dc42b5ab2955ef35bbb88))
* allow password_policy.max_age_days = 0 ([#1819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1819)) ([eeaadfe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/eeaadfef5b8fc31f84304899c7b205c9a97d8852))
* allow read of really old grant ids and add test cases ([#1615](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1615)) ([cda40ec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cda40ece534cdc3f6849a7d24f2f8acea8476e69))
* backwards compatability for grant helpers id used by procedure and functions ([#1508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1508)) ([3787657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3787657105fbcf18368136813afd558251f92cd1))
* be more cautious when setting db id to empty ([#1725](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1725)) ([e78e0c8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e78e0c816270779a78d839f78debc519fd53f5bf))
* change incorrect select privilege to usage ([#1681](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1681)) ([3d8f5b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d8f5b5ad417d2ebcf171024fc3217440267586b))
* change resource monitor suspend properties to number ([#1545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1545)) ([4bc59e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4bc59e24677260dae94952bdbc5176ad177876dd))
* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* cleanup acc tests ([#2135](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2135)) ([5db751d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5db751d1aa71952b1528e81cf2fdcd05d9d5d0fb))
* cleanup grant logic ([#1522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1522)) ([0502c61](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0502c61e7211253d029a0bec6a8104738624f243))
* cleanup workflows and makefile ([#2150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2150)) ([64335e7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/64335e72e480393437dff9f88122a256a2ac0814))
* database from_share handles organization ([#1711](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1711)) ([9e71eb6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e71eb6b4b69ed7fe1d42febe30fb166ea00812a))
* database grant read ([#2063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2063)) ([d93ddd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d93ddd564cd5cc76b47d49029f99c503d9314f29))
* default_secondary_roles doc ([#1584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1584)) ([23b64fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23b64fa9abcafb59610a77cafbda11a7e2ad648c))
* do not set "query_acceleration_max_scale_factor" if Query Acceleration is not enabled on Warehouse creation ([#1866](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1866)) ([7679e0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7679e0b0a480c6c6867f6e03051e7c3f97709edf))
* do not set query_acceleration_max_scale_factor when enable enable_query_acceleration = false ([#1474](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1474)) ([d62b1b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d62b1b4d6352e7d2dc99e4603370a1f3de3a4624))
* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))
* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))
* documentation for role ownership grant ([#2203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2203)) ([e3d405c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e3d405c91b494413d432e1aef9ff1da1f9ede4a7))
* don't ignore MUST_CHANGE_PASSWORD on account creation ([#1699](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1699)) ([7fcacfa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fcacfa59995766fb8c474b2984b9c9d3092d9ad))
* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* err when reading profile ([#1853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1853)) ([29c4633](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/29c463383da5fbb8de0af5c957aa9f3a34572c6e))
* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* external id for shares ([#2040](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2040)) ([f8357c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8357c06376feda2d9732f5046a408f9b82bc339))
* external tables issues ([#2334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2334)) ([ae41691](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae416917be72ab55c6f1b758dd7e06269831fabc))
* file format now supports reading options ([#1917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1917)) ([2596990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2596990808d97a7b98fa459a74eb7a77edf25cd5))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Fix [#1947](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1947), [#2134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2134), [#2173](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2173), and [#2176](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2176) ([#2192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2192)) ([98d8ccc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/98d8ccc3c52fbd76d9d24c5fc7091ec0afa30e1a))
* Fix database role docs ([#2081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2081)) ([339b1ff](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/339b1ff4fecca3dfe815231c9d3f26dc43a1ed7f))
* Fix email notification integration ([#2292](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2292)) ([70edd3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70edd3e5e010dcd2d2a6aa1bd1dd735a252d22c6))
* Fix encode Snowflake ID for object identifiers ([#2256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2256)) ([1c98a80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c98a80bb1486a04ead57dd0d8abcf65e00ea86c))
* Fix make sweep ([#2025](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2025)) ([beb2f98](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/beb2f9867ee501e630f75ace42dbef577563d593))
* Fix old masking policy implementation ([#2371](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2371)) ([c4837d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4837d660c9595f82ccb95997121aa3e4d49bb64)), closes [#2362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2362)
* Fix small issues ([#2287](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2287)) ([4e2613c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e2613c82edc1e0f91fb0cb56dcaa977d412d785))
* Fix some bugs ([#2234](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2234)) ([774a7db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/774a7db1ae912843f83561c3cd810c31af242958))
* Fix test because of the date ([#2312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2312)) ([9a9ea33](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a9ea3331f090201c2d686ed7e726eb7d9cef926))
* Fix warehouse read and resource monitor empty set ([#2319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2319)) ([05f96c6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/05f96c699e4e36db4b28380d9ac3577d4f50a709)), closes [#2318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2318) [#2316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2316)
* Fix whitespace in goreleaser envs ([#2388](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2388)) ([5e1266d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5e1266d6463219bbe9edddfd327d25c3b2d67f5f))
* Fix workflows ([#2206](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2206)) ([6d7f833](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6d7f8336897dee17c102d69a517e2525c1bb4d91))
* fixed enable_multiple_grants behaviour for role_grants ([#1816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1816)) ([f508129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f50812910a679bdea2d3add65f92d7cebeaa345c))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* fmt / linting fix ([#1694](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1694)) ([d7d910e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d910edc95069f6f01bc350c3319eb36f75c127))
* for 1624 resource monitor timestamps are always considered to have changed ([#2214](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2214)) ([4d5d3ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d5d3ca4cfa3727240b05da9f5010b4fa908f695))
* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* future ownership task grant ([#1954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1954)) ([#1955](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1955)) ([81ac1a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/81ac1a396bba713b0a5ec808d72b880d59e17d47))
* future read on grants ([#1520](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1520)) ([db78f64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/db78f64e56d228f3236b6bdefbe9a9c18c8641e1))
* Generate security integration statements using double quoted names. ([#1897](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1897)) ([a21d44f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a21d44fc9109c589a7e4c78032e2f5a62b6ed514))
* get `on_view` from `table_name` for `stream` ([#1740](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1740)) ([21fc2b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/21fc2b9ad9ddf02a241fbf2d1c3297d98cfa26c9))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* gorelease freebsd ([#1790](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1790)) ([c87f74f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c87f74f0cdffba16108ee822f845d0e961205252))
* goreleaser for mfa token caching ([#2320](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2320)) ([4fef709](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4fef709a376f5905effcc3439b1ad4cb9043ffca))
* grant imports ([#1713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1713)) ([19a4156](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19a415654f198a68e823b9ef62326992feba64df))
* grant read opeartion ([#2364](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2364)) ([5f51f0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f51f0f5c27f2d92409d8cc274917e70dd0c522c))
* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))
* Ignore two flaky table tests and disable comment unset for password policy ([#2086](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2086)) ([ee90014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee90014cc2d0dc5088d356f6c9b0d094c34d2a46))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* integration errors ([#1623](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1623)) ([83a40d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83a40d6361be0685b3864a0f3994298f3991de21))
* interval for failover groups ([#1448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1448)) ([bd1d3cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd1d3cc57f72c7774715f1d92a955536d55fb758))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))
* mfa token caching ([#2360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2360)) ([1404aac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1404aacaf3cdfc670d2b2a9e40c40b47c50f53d8))
* misc linting changes for 0.56.2 ([#1509](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1509)) ([e0d1ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d1ef5c718f9e1e58e80d31bbe2d2f27afec486))
* missing 'kind' field ([#1872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1872)) ([59b7ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/59b7ef5d791ecce366063a5cb8e2d68c68fb3e39))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))
* multiple share grants ([#1510](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1510)) ([d501226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d501226bc2ee8274446efb282c2dfea9599a3c2e))
* oauth integration ([#1315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1315)) ([9087220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9087220af85f08880f7ad453cbe9d13dd3bc11ec))
* openbsd build ([#1647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1647)) ([6895a89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6895a8958775e8e84a1457722f6c282d49458f3d))
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* prevent view drop when copy_grants attribute changes ([#1920](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1920)) ([0cf22d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0cf22d04513e70b0c945cd402b41b964f4d14be0))
* procedure and function grants ([#1502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1502)) ([0d08ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d08ea85541ceff6e591d34a671b44ef778a6611))
* provider config ([#2136](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2136)) ([07b9b4f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/07b9b4fee800fe3f34890783cc463d4fc5904717))
* provider muxing for plugin-framework ([#2130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2130)) ([f3c85c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3c85c0ebe3e54dc91836fe998c2d7c8e3373a52))
* readded imported privileges special case for database grants ([#1597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1597)) ([711ac0c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/711ac0cbc886bf8be6a5a2651234df778452b9df))
* reading file formats with no default schema ([#1936](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1936)) ([c5602f5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5602f5f01e09fc90a82aa2fbb8e116648f04c58))
* reading warehouse with underscore in name ([#1793](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1793)) ([5e184fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5e184fab4cbabcf27f290eb79c01a0cd6cc79282))
* refactor for simplify handling error ([#1472](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1472)) ([3937216](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/393721607c9eee5d73e14c27265eb39f195ccb37))
* refactor handling error to be simple ([#1473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1473)) ([9f37b99](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f37b997de073f01b66c86820237eff8049346ba))
* remove emojis, misc grant id fix ([#1598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1598)) ([fdefbc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdefbc3f1cc5bc7063f1cb1cc922854e8f9914e6))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* resource monitor import test ([#1690](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1690)) ([a4e58bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a4e58bfa0d5d4f90668784980fea7e2cc443f25d))
* resource snowflake_resource_monitor behavior conflict from provider 0.54.0 to 0.55.0 ([#1468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1468)) ([8ce0c53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ce0c533ec5d81273df20be2126b278ca61a59f6))
* resource_monitor suspend triggers conflict ([#1682](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1682)) ([14d1200](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14d120008d6d987dd2c0773ad9815eaac17a2c3c))
* Revert goreleaser mfa token caching ([#2343](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2343)) ([9a98031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a9803107f395ca350bdbb9fb8dc0f17820b8eb6))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* saml integration test ([#1494](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1494)) ([8c31439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c31439253d25aafb54fc09d89e547fa8238258c))
* saml2_sign_request and saml2_force_authn cast type ([#1452](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1452)) ([f8cecd7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8cecd7ca45aabec78fd18d8aa220db7eb34b9e0))
* schema name is optional for future file_format_grant ([#1484](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1484)) ([1450cdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1450cddde6328264f9df37e4dd89a78f5f095b2e))
* schema name is optional for future function_grant ([#1485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1485)) ([dcc550e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dcc550ed5b3df548d5d300cd2b77907ea544bb43))
* schema name is optional for future procedure_grant ([#1486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1486)) ([4cf4561](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cf456151d83cd71a3b9e68abe9c6f29804f2ee2))
* schema name is optional for future sequence_grant ([#1487](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1487)) ([ccf9e78](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccf9e78c9a7884e5beea233dd529a5134c741fb1))
* schema name is optional for future snowflake_stage_grant ([#1466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1466)) ([0b4d814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b4d8146910e8ea31d2ed5ea8b58725449205dcd))
* schema name is optional for future stream_grant ([#1488](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1488)) ([3f7e5d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f7e5d655ed5738107536c873dd11533573bba46))
* schema name is optional for future task_grant ([#1489](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1489)) ([4096fd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4096fd0d8bc65ae23b6d588385e1f81c4f2e7521))
* schema read now checks first if the corresponding database exists ([#1568](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1568)) ([368dc8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368dc8fb3f7e5156d16caed1e03792654d49f3d4))
* schema_name is optional to enable future pipe grant ([#1424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1424)) ([5d966fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d966fd8624fa3208cebae3d7b32c1b59bdcfd4c))
* schemas show by id ([#2238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2238)) ([bd8fc5d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd8fc5d1c9ad049ab76dbf5d689f5beb833a40ca))
* SNOW-59564 remove default data retention in days value ([#2029](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2029)) ([53c20b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53c20b4dc3b1141bcad07976f2dba93198e4d200))
* snowpipe error integration ([#2227](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2227)) ([0b388bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b388bf79346cd4ddedaac99d4651390f1f93358))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* support different tag association queries for COLUMN object types ([#1380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1380)) ([546d0a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/546d0a144e77c759cd6ddb91a193253f27f8fb91))
* support empty comments in oauth integrations ([#1900](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1900)) ([53d46b3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53d46b3a97680a4bf0e402d48c7db15c0d3c6f03))
* support multiple grant_privilege_to_role resources on the same role ([#1953](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1953)) ([dfdd166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dfdd166e20347deb6f464a5bbbe9c345287fcd6b))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* unable to alter stage if url and si have changes  ([#1982](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1982)) ([#1983](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1983)) ([3813aaa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3813aaac52a3ec5186e285066e4db72c2046531a))
* Update .goreleaser.yml ([#2394](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2394)) ([732164f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/732164f13b85ba11c1026a5b3734fac9f586052d))
* update 0.58.2 ([#1620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1620)) ([f1eab04](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1eab04dfdc839144057807953062b3591e6eaf0))
* Update CHANGELOG.md ([#2391](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2391)) ([29e5a82](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/29e5a821b757b1b6c2128513a4384ba89baed628))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))
* update packages ([#1619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1619)) ([79a3acc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79a3acc0e3d6a405593b5adf90a31afef81d700f))
* update read role grants to use new builder ([#1596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1596)) ([e91860a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e91860ae794b034158b71ffb31097e73d8015c51))
* use schema object identifier in external tables ([#2112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2112)) ([f5d4aeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f5d4aebf810697aca96764261b891415f002ee92))
* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))
* view statement update ([#2152](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2152)) ([6de32ae](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6de32ae6ec16ad76fb40afddfcaa7f650322cb67))
* warehouse for sdk v2 ([#1804](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1804)) ([99f7621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/99f7621cad1aafa83f96a4ed52d0fd04d6dfb2ae))
* warehouse size validation ([#1873](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1873)) ([5bbe460](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5bbe4600e2768a66a83df67a16368888b7f1d76b))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))
* wildcards in database name ([#1666](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1666)) ([54bf74c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/54bf74ca3d0119d31668d18bd1599ed029e386c8))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))

## [0.84.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.84.0...v0.84.1) (2024-01-22)


### 🔧 **Misc**

* skip acc managed account test ([#2376](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2376)) ([dbf645d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dbf645d8ea508582bebd1aab485f54807c818bd2))


### 🐛 **Bug fixes:**

* Fix whitespace in goreleaser envs ([#2388](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2388)) ([5e1266d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5e1266d6463219bbe9edddfd327d25c3b2d67f5f))
* grant read opeartion ([#2364](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2364)) ([5f51f0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f51f0f5c27f2d92409d8cc274917e70dd0c522c))

## [0.84.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.83.1...v0.84.0) (2024-01-19)


### 🎉 **What's new:**

* add app packages ([#2323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2323)) ([ca030fc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca030fce209ccc1e5e924b3bb1fd3c680c930471))
* add grant role resource ([#2304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2304)) ([ba91e25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba91e25be15197067e934eb867fb6837b4354e4a))
* add grant_database_role resource ([#2301](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2301)) ([2e7651f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2e7651f35b625bc11f17f42443d062118b6354ba))
* Add managed account to the SDK ([#2357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2357)) ([f968db1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f968db11a9b41ffa594170415c4edd9376f574e7))
* Add row access policy to the SDK ([#2363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2363)) ([14a3e5b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14a3e5b07a97cb18ba99f314990191241d4b15c2))
* Migrate storage integration to the new sdk ([#2339](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2339)) ([d970a56](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d970a567741edf1ce72e566f9e76b3858ece6288))


### 🔧 **Misc**

* 0.84 release ([#2374](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2374)) ([18638fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18638fea4c498158c8840a3ccfe03bdaac942923))
* Add missing issues to existing TODOs and add missing ones ([#2354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2354)) ([c3bc66c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c3bc66cffa4b87c49e027cc685df1abd8a6ebfe6))
* Skip managed account SDK integration tests temporarily ([#2370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2370)) ([9b2f551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9b2f551fc9df34930805b3429fccf88be7e7f851))
* update readme and fmt ([#2373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2373)) ([a648ca7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a648ca7ae1ca46d8bd05a11d7581ffd55d20ffd2))


### 🐛 **Bug fixes:**

* Fix old masking policy implementation ([#2371](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2371)) ([c4837d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4837d660c9595f82ccb95997121aa3e4d49bb64)), closes [#2362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2362)
* mfa token caching ([#2360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2360)) ([1404aac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1404aacaf3cdfc670d2b2a9e40c40b47c50f53d8))

## [0.83.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.83.0...v0.83.1) (2024-01-12)


### 🐛 **Bug fixes:**

* Revert goreleaser mfa token caching ([#2343](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2343)) ([9a98031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a9803107f395ca350bdbb9fb8dc0f17820b8eb6))

## [0.83.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.82.0...v0.83.0) (2024-01-11)


### 🎉 **What's new:**

* Add create streamlit privilege to the SDK ([#2303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2303)) ([be01d5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be01d5fdab4f2d31db9c4c849b349e657c0352c8))
* grant privileges to database role resource ([#2306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2306)) ([0311cf8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0311cf8554f0fbd202b489a54c9428f55c52a490))


### 🐛 **Bug fixes:**

* Add secondary account and fix tests ([#2324](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2324)) ([da6ca73](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da6ca733c7527c8918d1e1beb86d1641d94062ec))
* external tables issues ([#2334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2334)) ([ae41691](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae416917be72ab55c6f1b758dd7e06269831fabc))
* Fix test because of the date ([#2312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2312)) ([9a9ea33](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a9ea3331f090201c2d686ed7e726eb7d9cef926))
* Fix warehouse read and resource monitor empty set ([#2319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2319)) ([05f96c6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/05f96c699e4e36db4b28380d9ac3577d4f50a709)), closes [#2318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2318) [#2316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2316)
* goreleaser for mfa token caching ([#2320](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2320)) ([4fef709](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4fef709a376f5905effcc3439b1ad4cb9043ffca))

## [0.82.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.81.0...v0.82.0) (2023-12-21)


### 🎉 **What's new:**

* add functions to sdk ([#2205](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2205)) ([e542b67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e542b67f761850e812114ce593fa5f6deca941cb))


### 🔧 **Misc**

* Add short scripts used to fetch all currently opened issues ([#2288](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2288)) ([0b5ce4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b5ce4eff21562406a473ba8d796804efb1bd94f))
* Bump dependencies and fix linter complaints ([#2300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2300)) ([124e862](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/124e862b406526e08d865e0f70414926e50f4d4d))


### 🐛 **Bug fixes:**

* Fix email notification integration ([#2292](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2292)) ([70edd3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70edd3e5e010dcd2d2a6aa1bd1dd735a252d22c6))

## [0.81.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.80.0...v0.81.0) (2023-12-20)


### 🎉 **What's new:**

* identifier validation ([#2269](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2269)) ([9687972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9687972fc673bb7b768f9b4100a6d2b67fc46e48))


### 🔧 **Misc**

* Allow setting gosnowflake logging level from environment variable ([#2285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2285)) ([843e8fc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/843e8fc8f8879fdb5b1148155e7e64174571e975))
* Bump go to 1.21 ([#2267](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2267)) ([6b852c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b852c2ba6b1348ba83ec2e2190e8175371524be))


### 🐛 **Bug fixes:**

* Fix small issues ([#2287](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2287)) ([4e2613c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e2613c82edc1e0f91fb0cb56dcaa977d412d785))

## [0.80.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.79.1...v0.80.0) (2023-12-13)


### 🎉 **What's new:**

* Optional switch for instrumentedsql ([#2261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2261)) ([9934a59](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9934a59c5ee31dc292961806440427000ea41ea7))
* Use external tables from SDK ([#2228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2228)) ([6941023](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6941023f95230a14f92dd099c1a6375129ee4bfb))


### 🔧 **Misc**

* Bump dependencies ([#2252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2252)) ([581d75c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/581d75c566ea758c5f888dedb101e6318d7c0cc3))


### 🐛 **Bug fixes:**

* Fix encode Snowflake ID for object identifiers ([#2256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2256)) ([1c98a80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c98a80bb1486a04ead57dd0d8abcf65e00ea86c))
* for 1624 resource monitor timestamps are always considered to have changed ([#2214](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2214)) ([4d5d3ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d5d3ca4cfa3727240b05da9f5010b4fa908f695))

## [0.79.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.79.0...v0.79.1) (2023-12-11)


### 🐛 **Bug fixes:**

* Add missing parameters ([#2250](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2250)) ([4f4c4a4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4f4c4a4fcb91725a55726c03429745bbf95fd642))

## [0.79.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.78.0...v0.79.0) (2023-12-11)


### 🎉 **What's new:**

* add procedures to sdkv2 ([#2202](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2202)) ([6b563ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b563acef0f702fc6a366219afac602fa106129c))


### 🔧 **Misc**

* Propose additional debug logging ([#2243](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2243)) ([a0984cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a0984cd355c2004cc245fbd4672967ba59cfb3c7))

## [0.78.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.77.0...v0.78.0) (2023-12-08)


### 🎉 **What's new:**

* add event tables to sdk ([#2215](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2215)) ([66cc80a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/66cc80a8cac24f4b7967986032b7da9e20bd4eab))
* Add missing parameters to password policy resource ([#2231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2231)) ([c189782](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c189782bb7d117e27979c21dfea60ba733e996df))
* Add tables to the SDK ([#2042](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2042)) ([c1700de](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1700de2a062852da7cb5e3cf3277cc19f6466d6))


### 🐛 **Bug fixes:**

* Fix some bugs ([#2234](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2234)) ([774a7db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/774a7db1ae912843f83561c3cd810c31af242958))
* schemas show by id ([#2238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2238)) ([bd8fc5d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd8fc5d1c9ad049ab76dbf5d689f5beb833a40ca))
* snowpipe error integration ([#2227](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2227)) ([0b388bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b388bf79346cd4ddedaac99d4651390f1f93358))

## [0.77.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.76.0...v0.77.0) (2023-11-30)


### 🎉 **What's new:**

* Add unsafe_execute resource ([#2225](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2225)) ([196134c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/196134cbf91996eabc50bdc586a657fe7ac71900))
* Introduce simple arch tests ([#2210](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2210)) ([c60db80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c60db80f44d949258f0a692baafdc22b886c3010))


### 🐛 **Bug fixes:**

* cleanup workflows and makefile ([#2150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2150)) ([64335e7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/64335e72e480393437dff9f88122a256a2ac0814))
* documentation for role ownership grant ([#2203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2203)) ([e3d405c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e3d405c91b494413d432e1aef9ff1da1f9ede4a7))
* Fix workflows ([#2206](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2206)) ([6d7f833](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6d7f8336897dee17c102d69a517e2525c1bb4d91))

## [0.76.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.75.0...v0.76.0) (2023-11-15)


### Features

* Add "CREATE DYNAMIC TABLE" to schema_grant ([#2144](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2144)) ([6f026f6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f026f64e6e24638df2b9d4110362836a9071011))
* Add views to the SDK ([#2171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2171)) ([ed079d3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed079d3d06dc3af083da04ca18314c8e7b07308e))
* Migrate application role to new sdk ([#2149](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2149)) ([7abb4db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7abb4dbb645f0eb0fcb1d39414b1ed0c322916c9))
* Migrate stage to the new sdk ([#2163](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2163)) ([1d08c46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d08c46c470cda2d6b7c7ed3d599d2c1a57b837b))
* Poc custom error type ([#2052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2052)) ([b86c4c3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b86c4c34d05f8b982fb6218a3a3a7500a23abf72))
* Use tasks from the SDK followup ([#2153](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2153)) ([82c3c13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82c3c13b6166168e470d7cb9b2982a8979275f17))


### Misc

* Add migration guide ([#2142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2142)) ([ee0f6af](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee0f6af54dbd269f8bfa2c3d73a396d98d10a6ab))
* Split existing alter operations ([#2156](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2156)) ([dbb7c91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dbb7c9136c586490a0856cc07ae879be491c8150))


### BugFixes

* Fix [#1947](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1947), [#2134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2134), [#2173](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2173), and [#2176](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2176) ([#2192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2192)) ([98d8ccc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/98d8ccc3c52fbd76d9d24c5fc7091ec0afa30e1a))
* provider muxing for plugin-framework ([#2130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2130)) ([f3c85c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3c85c0ebe3e54dc91836fe998c2d7c8e3373a52))

## [0.75.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.74.0...v0.75.0) (2023-10-26)


### Features

* add parse_header option to file format resource ([#2132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2132)) ([1e6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e6e54f828efa60edd258b316709fc4dfd370f93))
* Use streams from the new SDK in resource / datasource ([#2129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2129)) ([5c633be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5c633be461fd373d412b02b108e64b6cfc4eb856))
* Use task from SDK in resource and data source ([#2140](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2140)) ([de23f2b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de23f2ba939eb368d9734217e1bb2d4ebc75eef4))


### Misc

* Return multiple errors in existing validations ([#2122](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2122)) ([4d4bcdb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d4bcdbe841807da2fa08d534eaf846234934f7c))
* Set up a single warehouse for the SDK integration tests ([#2141](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2141)) ([16022ef](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16022ef4171e7dccf2932ae6e8d451b51c93291c))


### BugFixes

* cleanup acc tests ([#2135](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2135)) ([5db751d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5db751d1aa71952b1528e81cf2fdcd05d9d5d0fb))
* provider config ([#2136](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2136)) ([07b9b4f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/07b9b4fee800fe3f34890783cc463d4fc5904717))
* view statement update ([#2152](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2152)) ([6de32ae](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6de32ae6ec16ad76fb40afddfcaa7f650322cb67))

## [0.74.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.73.0...v0.74.0) (2023-10-18)


### Features

* dynamic table resource + data source ([#2104](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2104)) ([59c5de6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/59c5de6becf4b16542da4db59e839f7ca5625e50))
* Migrate streams to the new sdk ([#2113](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2113)) ([521fde5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/521fde58b7836a21614d2ae31a7060352a68b465))
* Move integration tests to separate package ([#2111](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2111)) ([2755589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2755589bc8c7ecbda3115dc574339c40bfd5b096))
* Provider config refactor ([#2126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2126)) ([9858fea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9858feaa2fb382b57e3c89aae044f09ce2dcc5a3))
* tags sdk refactoring ([#2079](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2079)) ([7013f83](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7013f835864139ef40c2521e6af87589e909e201))
* Use network policy from sdk ([#2087](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2087)) ([50f2935](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/50f2935dc12fc2479c2919a5666d5573a8b737ff))
* use SDK in schema resource and datasource ([#2082](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2082)) ([f7d0d97](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f7d0d9716c4da9b00e5171a368a03a9cd52c31ae))


### Misc

* Adjust integration tests after moving to separate package ([#2115](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2115)) ([3f528a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f528a87f4c0b3bc95a0dfb35d93d22251b5112e))
* Integration tests setup db once ([#2118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2118)) ([f533368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f5333684d5e46877ec4da6dc708ce2b09b5694af))
* Set up schema once for integration tests ([#2121](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2121)) ([5e3ebf4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5e3ebf4e809a4be163a30642f69ea1b9146f21a9))
* Use helper methods in old unit tests ([#2119](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2119)) ([0c44571](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0c44571a92b842ac97d46a38e5bcfc1066367177))


### BugFixes

* use schema object identifier in external tables ([#2112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2112)) ([f5d4aeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f5d4aebf810697aca96764261b891415f002ee92))

## [0.73.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.72.0...v0.73.0) (2023-10-11)


### Features

* add authenticator ([#2109](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2109)) ([4f3a551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4f3a5519484b0aab91ff5fa08f37a8cf512d1ec0))
* Add task clone to sdk ([#2105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2105)) ([acddb2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/acddb2cd6bfb1a7ffaf6dbb3c8349f7bc550c124))
* Add task to SDK ([#2099](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2099)) ([d52f334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d52f3347f091f0edff5e6daded1120542f1e9bd1))


### Misc

* **deps:** bump golang.org/x/crypto from 0.13.0 to 0.14.0 ([#2100](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2100)) ([060750f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/060750f668b986d690619cd2d915a7a945350ec5))
* **deps:** bump golang.org/x/tools from 0.13.0 to 0.14.0 ([#2101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2101)) ([45a14d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/45a14d96d5fd6b051a73ced2fb03cab8659ded84))

## [0.72.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.71.0...v0.72.0) (2023-10-04)


### Features

* add dynamic tables to sdk ([#2074](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2074)) ([d1dfb05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d1dfb05fbb3bcc59cc2622b6b2d02ebadf1cf33f))
* Add grant ownership to SDK ([#2064](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2064)) ([f85ec8b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f85ec8baa2f9aaeead4f619dccfa3d38880a16d7))
* Add Manage Warehouses Account Grant ([#2017](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2017)) ([89c7148](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/89c7148c11378af9e42ea32bdad3e5a5c465d39c))
* add mfa auth ([#2077](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2077)) ([922358a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/922358a43c5383ee5840bf2971ecd27d96f86573))
* Add session policy to SDK ([#2088](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2088)) ([038241c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/038241c00a158b389df8034864a52a252fcc41bf))
* Go Generator + Network Policy migration to the new SDK ([#2061](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2061)) ([231b081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/231b081a34462392954c18b0f3452b1fed7d016f))
* Migrate external tables to new sdk ([#2006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2006)) ([5af17cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5af17cfbe4b6570c61b4ac77ab7564fd1f0529b7))
* Migrate roles to the new sdk ([#2007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2007)) ([5b996e6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b996e6f4b3ea28bae0cd13e75bd112d9af09c68))
* support drop for accounts and fix read race condition ([#2065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2065)) ([5412252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/541225258ff19de864c8d25d05a8c1b4a7941813))


### Misc

* Change return type of show ([#2045](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2045)) ([21f069a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/21f069a3232a0039792f30445e25b5e88f09ab25))
* Cleanup comments ([#2092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2092)) ([3a06a66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3a06a66f27d9b03d037affdda5b7dbb3bd81d5d7))
* **deps:** bump github.com/hashicorp/terraform-plugin-docs from 0.14.1 to 0.16.0 ([#1931](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1931)) ([07c4f14](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/07c4f149bf08ea9fecd28aa12a47a9fae3e0523b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#2090](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2090)) ([24b6313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/24b6313a7f28731ad4c10dd80c32e7b9cf397c5d))
* **deps:** bump github.com/stretchr/testify from 1.8.2 to 1.8.4 ([#1842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1842)) ([9b0825d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9b0825db520c244cea962830d2139fb50186f23c))
* **deps:** bump golang.org/x/tools from 0.7.0 to 0.13.0 ([#2089](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2089)) ([0ace968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ace9689c6b29382b1d93d9a3542eacc0d7c4270))
* Make options naming consistent ([#2083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2083)) ([df84bd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/df84bd084482ef3f85db2b1f4a3517afe8df4ddb))


### BugFixes

* added privilege for `RESOLVE ALL` ([#1861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1861)) ([18cf7b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18cf7b09c9be927b5e1a625acdadb90abd048ff8))
* database grant read ([#2063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2063)) ([d93ddd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d93ddd564cd5cc76b47d49029f99c503d9314f29))
* Fix database role docs ([#2081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2081)) ([339b1ff](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/339b1ff4fecca3dfe815231c9d3f26dc43a1ed7f))
* Ignore two flaky table tests and disable comment unset for password policy ([#2086](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2086)) ([ee90014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee90014cc2d0dc5088d356f6c9b0d094c34d2a46))

## [0.71.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.70.1...v0.71.0) (2023-09-21)


### Features

* Add grant privileges to database role to SDK ([#2023](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2023)) ([717289f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/717289f71fd4a08f44d4f20f6e16dc4dded77803))
* Convert alert datasource to new sdk ([#2020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2020)) ([2d0eaeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d0eaebe8ea1605251dd24f86a2f83e6e3be67c8))
* Initialize SDK generator ([#2033](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2033)) ([96b47e5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/96b47e5fdbb6175a7a24d0518b975279927c22de))
* Parameters sdk v2 ([#1914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1914)) ([1d15355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1535572e1b82f3cb2308e3c19e632eb34fe24b))


### BugFixes

* add dynamic tables to grants ([#2059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2059)) ([3767de9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3767de97ec685acebd0e8e105878eb9129a0752c))
* SNOW-59564 remove default data retention in days value ([#2029](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2029)) ([53c20b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53c20b4dc3b1141bcad07976f2dba93198e4d200))

## [0.70.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.70.0...v0.70.1) (2023-09-01)


### Misc

* Extract reviewdog job from integration workflow ([#2027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2027)) ([243fc28](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/243fc285811de787328a7f3ca0d3a84193976b85))


### BugFixes

* add back in 386 build ([#2038](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2038)) ([5522504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55225041d7f4722adbd7b9445990f57c04a4adea))
* external id for shares ([#2040](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2040)) ([f8357c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8357c06376feda2d9732f5046a408f9b82bc339))
* Fix make sweep ([#2025](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2025)) ([beb2f98](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/beb2f9867ee501e630f75ace42dbef577563d593))

## [0.70.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.69.0...v0.70.0) (2023-08-21)


### Features

* Add database role to SDK ([#2009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2009)) ([f5efc09](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f5efc09ea60bd2d66c65c9e07cb84321f95531f0))
* Add missing database role operations to SDK ([#2014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2014)) ([d2ea67d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2ea67d7fee00b15e1222fe37efe8e7a1cddecb5))
* Add Pipes to SDK ([#1968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1968)) ([69a543f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69a543fd729b64bdd8964dc34626dee83b3f96a7))
* Generate constructor and builder methods for DTOs ([#2001](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2001)) ([79d9c9c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79d9c9ce77e85dc9bd0889da30155cef0a74c293))
* Schemas sdk v2 ([#1975](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1975)) ([289ad8a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/289ad8a8d61cf217d4548a65cf4843e416baa1da))
* Use database role from SDK ([#2012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2012)) ([294075a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/294075a4d145dc08071c11f067afaecb78fe8ef7))
* Use pipes from SDK ([#2003](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2003)) ([079d47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/079d47d22af22edba6c6c499409b264c4c5f5945))
* Users sdk v2 ([#1945](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1945)) ([d644b63](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d644b63e673223b67f7092563a0bfc8c826fc38b))


### Misc

* Makefile cleanup ([#1995](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1995)) ([5c6fdbe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5c6fdbe48aadab9d56bfb7cf5d5ea30ca5a2b356))
* Remove lint ignore comments - test ([#2011](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2011)) ([8f81a42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f81a42544a668b03e0bb5d4d3199e6879d2f937))
* Rename and reuse validatableOpts interface ([#2016](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2016)) ([4b42848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b428480a3f8fb9aa38da7a311f8a04f8d294882))


### BugFixes

* unable to alter stage if url and si have changes  ([#1982](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1982)) ([#1983](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1983)) ([3813aaa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3813aaac52a3ec5186e285066e4db72c2046531a))

## [0.69.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.68.2...v0.69.0) (2023-07-29)


### ⚠ BREAKING CHANGES

* mark snowflake_user.login_name as non sensitive

### Features

* mark snowflake_user.login_name as non sensitive ([71c0a9c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/71c0a9c2897e07b2251680d619cbbcaeba5779b2))
* resource monitor sdk v2 ([#1892](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1892)) ([707d723](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/707d7233bd867d60e13c8785d1c0adef558a5853))


### BugFixes

* 1969  change default privilege to OPERATE ([#1970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1970)) ([e8721f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e8721f07dd27ce7a903c8c4d2b19fee77ce72cdc))
* add external table grant to snowflake_grant_privileges_to_role resource ([#1967](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1967)) ([8c84c4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c84c4a5f6eee7cf35b5ec2d77fac1c7fcee05df))
* future ownership task grant ([#1954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1954)) ([#1955](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1955)) ([81ac1a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/81ac1a396bba713b0a5ec808d72b880d59e17d47))

## [0.68.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.68.1...v0.68.2) (2023-07-17)


### BugFixes

* support multiple grant_privilege_to_role resources on the same role ([#1953](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1953)) ([dfdd166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dfdd166e20347deb6f464a5bbbe9c345287fcd6b))
* supress warehouse size changes ([#1889](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1889) ([daf2cc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daf2cc289bab892e5a0954f852a72a9f99e1b0c1))

## [0.68.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.68.0...v0.68.1) (2023-07-07)


### BugFixes

* reading file formats with no default schema ([#1936](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1936)) ([c5602f5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5602f5f01e09fc90a82aa2fbb8e116648f04c58))

## [0.68.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.67.0...v0.68.0) (2023-07-06)


### Features

* alerts sdk v2 ([#1871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1871)) ([fcf0380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fcf03800cf92e1acff7f89e59968914f1bcad4f4))
* new grant resource "snowflake_grant_privileges_to_role"  ([#1929](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1929)) ([e241c22](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e241c223b0c5d159efc803b9fcd4671c8325ae00))
* set query tag for terraform sessions  ([#1826](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1826)) ([6629583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6629583bc288e5c23cf02d8a53597adbadf55fce))


### BugFixes

* add missing option ([#1924](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1924)) ([2991a16](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2991a16fe45491d77a8212a3ce0aa9566882de4f))
* do not set "query_acceleration_max_scale_factor" if Query Acceleration is not enabled on Warehouse creation ([#1866](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1866)) ([7679e0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7679e0b0a480c6c6867f6e03051e7c3f97709edf))
* file format now supports reading options ([#1917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1917)) ([2596990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2596990808d97a7b98fa459a74eb7a77edf25cd5))
* prevent view drop when copy_grants attribute changes ([#1920](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1920)) ([0cf22d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0cf22d04513e70b0c945cd402b41b964f4d14be0))

## [0.67.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.66.2...v0.67.0) (2023-06-22)


### Features

* account password policy attachment ([#1824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1824)) ([f408828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f408828fd023c2207ec41f702cec7dae524b1e93))


### BugFixes

* Generate security integration statements using double quoted names. ([#1897](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1897)) ([a21d44f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a21d44fc9109c589a7e4c78032e2f5a62b6ed514))
* missing 'kind' field ([#1872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1872)) ([59b7ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/59b7ef5d791ecce366063a5cb8e2d68c68fb3e39))
* support empty comments in oauth integrations ([#1900](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1900)) ([53d46b3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53d46b3a97680a4bf0e402d48c7db15c0d3c6f03))

## [0.66.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.66.1...v0.66.2) (2023-06-13)


### BugFixes

* warehouse size validation ([#1873](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1873)) ([5bbe460](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5bbe4600e2768a66a83df67a16368888b7f1d76b))

## [0.66.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.66.0...v0.66.1) (2023-06-03)


### BugFixes

* err when reading profile ([#1853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1853)) ([29c4633](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/29c463383da5fbb8de0af5c957aa9f3a34572c6e))

## [0.66.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.65.0...v0.66.0) (2023-06-02)


### Features

* databases sdk v2 ([#1840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1840)) ([410b493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/410b493ff75188b70cee91db56c60538d6823aac))


### Misc

* rename "db" symbol to "sql" and make input options more uniform in naming convention ([#1837](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1837)) ([244b5b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/244b5b45e22f66875f14be9bcb9af08ad41d37ca))

## [0.65.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.64.0...v0.65.0) (2023-05-30)


### Features

* add failover groups to sdk v2, and add data source for failover groups ([#1825](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1825)) ([44e8c06](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44e8c06ba4c665c81f0b909dbd3df90c4925e179))
* add shares to sdk v2 ([#1813](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1813)) ([a814841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a814841011f08857d8d37691fa5ff01cd9412176))
* adding a subset of missing parameters ([#1828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1828)) ([6d1a572](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6d1a572623438dd96139edbd698ce5974fa3df61))
* adding CREATE STREAMLIT as privalege ([#1771](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1771)) ([7b112cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b112cc94822807d3c25c78ec812ec2c2c66ea16))
* alerts datasource ([#1727](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1727)) ([a1d0d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1d0d2eeb9a8a03d4c347b9f153182e66c237a4f))
* allow all priv on grants ([#1786](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1786)) ([aa3f873](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aa3f873b4669c5942c10681151d6afd4ea03623f))
* support custom ownership transfer for grants ([#1743](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1743)) ([eaa6e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/eaa6e01820cb04ffa3c647ecd893e479af8e35a1))


### Misc

* Account SDK ([#1822](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1822)) ([4c22b89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c22b89e82a608375fd78993134cf3a0b45b985b))
* Migrate Warehouse resource + datasource to new SDK ([#1792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1792)) ([a14b994](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a14b994d62f73c272cde72651c3c8c18ac4213cf))


### BugFixes

* 0.65 integration tests and minor linting fixes ([#1835](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1835)) ([8d7d663](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8d7d6636a67d272d5b87e02a178ae799b5b6d777))
* add grant to sdk v2, better validation ([#1812](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1812)) ([76da3b8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/76da3b84d9d65493ae666fad6ea6d926b81add8f))
* Add missing Google Cloud computed property api_gcp_service_account to the api integration resource ([#1776](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1776)) ([1dd0672](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dd06724aa1add9b504c820583a2a752555c9bdb))
* allow password_policy.max_age_days = 0 ([#1819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1819)) ([eeaadfe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/eeaadfef5b8fc31f84304899c7b205c9a97d8852))
* fixed enable_multiple_grants behaviour for role_grants ([#1816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1816)) ([f508129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f50812910a679bdea2d3add65f92d7cebeaa345c))
* gorelease freebsd ([#1790](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1790)) ([c87f74f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c87f74f0cdffba16108ee822f845d0e961205252))
* reading warehouse with underscore in name ([#1793](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1793)) ([5e184fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5e184fab4cbabcf27f290eb79c01a0cd6cc79282))
* warehouse for sdk v2 ([#1804](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1804)) ([99f7621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/99f7621cad1aafa83f96a4ed52d0fd04d6dfb2ae))

## [0.64.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.63.0...v0.64.0) (2023-05-09)


### Features

* Add external functions translators ([#1735](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1735)) ([1f67286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1f672862adb29a658f5e81e940f9afb994347f2f))
* addition of the on_all option for all remaining resources ([#1742](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1742)) ([e3a7710](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e3a771000fa527fc8420b892450031ba39394fcf))
* database roles datasource ([#1726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1726)) ([f1cb7f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1cb7f3b368c0a7fa6a0a08784d53b0e92efda81))
* email notification integration ([#1720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1720)) ([5d21fe1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d21fe19f0da9ed4939e9bd6dc66237685c5a327))
* masking policy in v2 sdk ([#1777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1777)) ([6978c42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6978c42b304bd8c5429c88bcd7d6ed20ac3fd98a))


### Misc

* sdk v2 implemented for password policy ([#1752](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1752)) ([0cb1164](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0cb1164414607ce2e4ac7f6fb2da1c563b22da7b))


### BugFixes

* add missing reference for email_notification ([#1768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1768)) ([6feeb9c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6feeb9c09f56df0fb222bd0664508defb7c5af6e))
* database from_share handles organization ([#1711](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1711)) ([9e71eb6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e71eb6b4b69ed7fe1d42febe30fb166ea00812a))
* get `on_view` from `table_name` for `stream` ([#1740](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1740)) ([21fc2b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/21fc2b9ad9ddf02a241fbf2d1c3297d98cfa26c9))

## [0.63.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.62.0...v0.63.0) (2023-04-25)


### Features

* masking policy application resource ([#1739](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1739)) ([ce80f57](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce80f57225d6019765a50da6eb74619e401db622))


### Misc

* fix parameters datasource acceptance test ([#1749](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1749)) ([338a19d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/338a19de373525087372cc319aa065f2cddbb724))

## [0.62.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.61.0...v0.62.0) (2023-04-19)


### Features

* Add gcp_pubsub_topic_name parameter for gcp notification integration ([#1687](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1687)) ([a30d0cb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a30d0cb756a2281a4d880af9e32651c04409028e))
* added on_all grants for view, stage, schema and materialized view resource ([#1686](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1686)) ([f27a9e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f27a9e454102709e1134d796f1594b4932670ae8))
* alert resource ([#1653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1653)) ([4a84eb9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a84eb91af6201330a2ec60de306db5ed3b5ab38))
* data source shares ([#1651](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1651)) ([5dd447f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dd447fc79c153f025e420052a23b90751cdaa9f))
* database role resource ([#1654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1654)) ([9f4a516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f4a5165d6b7e1fa63ad06450ffc55689be66356))
* external_oauth_integration.scope_mapping_attribute ([#1722](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1722)) ([66d88bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/66d88bd29b6aa78dbc2acd11cf729fa9380b011d))
* Password policy resource ([#1702](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1702)) ([7ee293b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ee293b52262170a1d53811f9beb6fc3772ae913))
* provider config support insecureMode ([#1703](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1703)) ([e269925](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e269925e0f54ca9bbd742e5c60eaebad5f9b1f04))


### BugFixes

* Adding double quotes to the column names ([#1729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1729)) ([791dd0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/791dd0b3a613d0ada95d4af7153a7b0a2bbee219))
* be more cautious when setting db id to empty ([#1725](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1725)) ([e78e0c8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e78e0c816270779a78d839f78debc519fd53f5bf))
* change incorrect select privilege to usage ([#1681](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1681)) ([3d8f5b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d8f5b5ad417d2ebcf171024fc3217440267586b))
* don't ignore MUST_CHANGE_PASSWORD on account creation ([#1699](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1699)) ([7fcacfa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fcacfa59995766fb8c474b2984b9c9d3092d9ad))
* fmt / linting fix ([#1694](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1694)) ([d7d910e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d910edc95069f6f01bc350c3319eb36f75c127))
* grant imports ([#1713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1713)) ([19a4156](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19a415654f198a68e823b9ef62326992feba64df))
* resource monitor import test ([#1690](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1690)) ([a4e58bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a4e58bfa0d5d4f90668784980fea7e2cc443f25d))
* resource_monitor suspend triggers conflict ([#1682](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1682)) ([14d1200](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14d120008d6d987dd2c0773ad9815eaac17a2c3c))

## [0.61.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.0...v0.61.0) (2023-04-03)


### Features

* Add COPY GRANTS arg to views resources. ([#1668](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1668)) ([7225d93](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7225d93ef3e50c6810c0dd57cfd7079e882d443f))
* add on_account to session and object params ([#1685](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1685)) ([1329430](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/13294304c7626c9d428682986669d2e97ab2c23b))
* mark PII as sensitive in account and user resources. ([#1678](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1678)) ([caa461f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/caa461fea56df4f33f73e586b537802c01f1eb5d))
* support is_org_admin column in SHOW ORGANIZATION ACCOUNTS ([#1673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1673)) ([263c521](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/263c5215dedf70195c0c2c7d8e1505e4b9c0828c))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.0...v0.34.0) (2023-03-28)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* add GRANT ... ON ALL TABLES IN ... ([#1626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1626)) ([505a5f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/505a5f35d9ea8388ca33e5117c545408243298ae))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* add in more functionality for UpdateResourceMonitor  ([#1456](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1456)) ([2df570f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2df570f0c3271534a37b0cb61b7f4b491081baf7))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* add missing PrivateLink URLs to datasource ([#1603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1603)) ([78782b1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78782b1b471b7fbd434de1803cd687f6866cada7))
* add new account resource ([#1492](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1492)) ([b1473ba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1473ba158946d81bc4eac095c40c8d0446cf2ed))
* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* add ON STAGE support for Stream resource ([#1413](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1413)) ([447febf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/447febfef46ef89570108d3447998d6d379b7be7))
* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))
* add PREVENT_LOAD_FROM_INLINE_URL ([#1612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1612)) ([4945a3a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4945a3ae62d887dae6332742edcde715751459b5))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* Add support for `packages`, `imports`, `handler` and `runtimeVersion` to `snowflake_procedure` resource ([#1516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1516)) ([a88f3ad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a88f3ada75dad18b7b4b9024f664de8d687f54e0))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for is_secure to snowflake_function resource ([#1575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1575)) ([c41b6a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c41b6a35271f12c97f5a4da947eb433013f2aaf2))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding slack bot for PRs and Issues ([#1106](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1106)) ([95c255e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/95c255e5ca65b619b35692671848877793cac29e))
* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))
* **integration:** add google api integration ([#1589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1589)) ([56909cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56909cdc18245f38b0f58bceaf2aa9cbc013d212))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* roles support numbers ([#1585](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1585)) ([d72dee8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d72dee82d0e0a4d8b484e5b204e156a13117cb76))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* support object parameters on account level ([#1583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1583)) ([fb24802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb2480214c8ac4e61fffe3a8e3448597462ad9a1))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))
* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### BugFixes

* 0.54  ([#1435](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1435)) ([4c9dd13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c9dd133574b08d8e67444b6c6b81aa87d9a2acf))
* 0.55 fix ([#1465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1465)) ([8cb3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8cb337048ec5c4a52245feb1b9556fd845d83278))
* 0.59 release fixes ([#1636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1636)) ([0a0256e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a0256ed3f0d56a6c7c22f810419636685094135))
* 0.60 misc bug fixes / linting ([#1643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1643)) ([53da853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53da853c213eec3afbdd2a47a8de3fba897c5d8a))
* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add contributing section to readme ([#1560](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1560)) ([174355d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/174355d740e325ae05435bbbc22b8b335f94fc6f))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* add nill check for grant_helpers ([#1518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1518)) ([87689bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/87689bb5b60c73bfe3d741c3da6f4f544f16aa45))
* add permissions ([#1464](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1464)) ([e2d249a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d249ae1466e05dad2080f05123e0de66fabcf6))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* add test cases for update repl schedule on failover group ([#1578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1578)) ([ab638f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab638f0b9ba866d22c6f807743eb4eccad2530b8))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* added force_new option to role grant when the role_name has been changed ([#1591](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1591)) ([4ec3613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ec3613de43d70f40a5d29ce5517af53e8ef0a06))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* all-grants ([#1658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1658)) ([d5d59b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5d59b4e85cd2e97ea0dc42b5ab2955ef35bbb88))
* Allow creation of database-wide future external table grants ([#1041](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1041)) ([5dff645](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dff645291885cd437e341148c0629fe7ab7383f))
* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* allow read of really old grant ids and add test cases ([#1615](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1615)) ([cda40ec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cda40ece534cdc3f6849a7d24f2f8acea8476e69))
* backwards compatability for grant helpers id used by procedure and functions ([#1508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1508)) ([3787657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3787657105fbcf18368136813afd558251f92cd1))
* change resource monitor suspend properties to number ([#1545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1545)) ([4bc59e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4bc59e24677260dae94952bdbc5176ad177876dd))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* cleanup grant logic ([#1522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1522)) ([0502c61](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0502c61e7211253d029a0bec6a8104738624f243))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* default_secondary_roles doc ([#1584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1584)) ([23b64fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23b64fa9abcafb59610a77cafbda11a7e2ad648c))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* do not set query_acceleration_max_scale_factor when enable enable_query_acceleration = false ([#1474](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1474)) ([d62b1b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d62b1b4d6352e7d2dc99e4603370a1f3de3a4624))
* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))
* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))
* future read on grants ([#1520](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1520)) ([db78f64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/db78f64e56d228f3236b6bdefbe9a9c18c8641e1))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* integration errors ([#1623](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1623)) ([83a40d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83a40d6361be0685b3864a0f3994298f3991de21))
* interval for failover groups ([#1448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1448)) ([bd1d3cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd1d3cc57f72c7774715f1d92a955536d55fb758))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))
* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* misc linting changes for 0.56.2 ([#1509](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1509)) ([e0d1ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d1ef5c718f9e1e58e80d31bbe2d2f27afec486))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))
* multiple share grants ([#1510](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1510)) ([d501226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d501226bc2ee8274446efb282c2dfea9599a3c2e))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* oauth integration ([#1315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1315)) ([9087220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9087220af85f08880f7ad453cbe9d13dd3bc11ec))
* openbsd build ([#1647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1647)) ([6895a89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6895a8958775e8e84a1457722f6c282d49458f3d))
* OSCP -&gt; OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* procedure and function grants ([#1502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1502)) ([0d08ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d08ea85541ceff6e591d34a671b44ef778a6611))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))
* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))
* read Database and Schema name during Stream import ([#732](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/732)) ([9f747b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f747b571b2fcf5b0663696efd75c55a6f8b6a89))
* read Name, Database and Schema during Procedure import ([#819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/819)) ([d17656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d17656fdd2803516b6ee067a6bd5a457bf31d905))
* readded imported privileges special case for database grants ([#1597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1597)) ([711ac0c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/711ac0cbc886bf8be6a5a2651234df778452b9df))
* Recreate notification integration when type changes ([#792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/792)) ([e9768bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9768bf059268fb933ad74f2b459e91e2c0563e0))
* refactor for simplify handling error ([#1472](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1472)) ([3937216](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/393721607c9eee5d73e14c27265eb39f195ccb37))
* refactor handling error to be simple ([#1473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1473)) ([9f37b99](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f37b997de073f01b66c86820237eff8049346ba))
* refactor ReadWarehouse function to correctly read object parameters ([#745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/745)) ([d83c499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83c499910c0f2b6348191a93f917e450b9e69b2))
* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))
* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))
* remove emojis, misc grant id fix ([#1598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1598)) ([fdefbc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdefbc3f1cc5bc7063f1cb1cc922854e8f9914e6))
* Remove force_new from masking_expression ([#588](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/588)) ([fc3e78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3e78acbdda5346f32a004711d31ad6f68529ed))
* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* Removing force new and adding update for data base replication config ([#1105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1105)) ([f34f012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f34f012195d0b9718904ffa7a3a529f58167a74e))
* resource snowflake_resource_monitor behavior conflict from provider 0.54.0 to 0.55.0 ([#1468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1468)) ([8ce0c53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ce0c533ec5d81273df20be2126b278ca61a59f6))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* saml integration test ([#1494](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1494)) ([8c31439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c31439253d25aafb54fc09d89e547fa8238258c))
* saml2_sign_request and saml2_force_authn cast type ([#1452](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1452)) ([f8cecd7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8cecd7ca45aabec78fd18d8aa220db7eb34b9e0))
* schema name is optional for future file_format_grant ([#1484](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1484)) ([1450cdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1450cddde6328264f9df37e4dd89a78f5f095b2e))
* schema name is optional for future function_grant ([#1485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1485)) ([dcc550e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dcc550ed5b3df548d5d300cd2b77907ea544bb43))
* schema name is optional for future procedure_grant ([#1486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1486)) ([4cf4561](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cf456151d83cd71a3b9e68abe9c6f29804f2ee2))
* schema name is optional for future sequence_grant ([#1487](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1487)) ([ccf9e78](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccf9e78c9a7884e5beea233dd529a5134c741fb1))
* schema name is optional for future snowflake_stage_grant ([#1466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1466)) ([0b4d814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b4d8146910e8ea31d2ed5ea8b58725449205dcd))
* schema name is optional for future stream_grant ([#1488](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1488)) ([3f7e5d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f7e5d655ed5738107536c873dd11533573bba46))
* schema name is optional for future task_grant ([#1489](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1489)) ([4096fd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4096fd0d8bc65ae23b6d588385e1f81c4f2e7521))
* schema read now checks first if the corresponding database exists ([#1568](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1568)) ([368dc8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368dc8fb3f7e5156d16caed1e03792654d49f3d4))
* schema_name is optional to enable future pipe grant ([#1424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1424)) ([5d966fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d966fd8624fa3208cebae3d7b32c1b59bdcfd4c))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* support different tag association queries for COLUMN object types ([#1380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1380)) ([546d0a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/546d0a144e77c759cd6ddb91a193253f27f8fb91))
* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))
* update 0.58.2 ([#1620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1620)) ([f1eab04](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1eab04dfdc839144057807953062b3591e6eaf0))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))
* update packages ([#1619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1619)) ([79a3acc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79a3acc0e3d6a405593b5adf90a31afef81d700f))
* update read role grants to use new builder ([#1596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1596)) ([e91860a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e91860ae794b034158b71ffb31097e73d8015c51))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))
* Updating shares to disallow account locators ([#1102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1102)) ([4079080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4079080dd0b9e3caf4b5d360000bd216906cb81e))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))
* wildcards in database name ([#1666](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1666)) ([54bf74c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/54bf74ca3d0119d31668d18bd1599ed029e386c8))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))
* **deps:** bump actions/setup-go from 3 to 4 ([#1634](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1634)) ([3f128c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f128c1ba887c377b7bd5f3d508d7b81676fdf90))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1639)) ([330777e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/330777eb0ad93acede6ffef9d7571c0989540657))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1638)) ([107bb4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/107bb4abfb5de896acc1f224afae77b8100ffc02))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1454)) ([ed6bfe0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed6bfe07930e5703036ada816845176d46f5623c))
* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#1528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1528)) ([8a011e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a011e0b1920833c77eb7832f821a4bd52176657))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#1551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1551)) ([35de62f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/35de62f5b722c1dc6eaf2f39f6699935f67557cd))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump golang.org/x/tools from 0.4.0 to 0.5.0 ([#1455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1455)) ([ff01970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff019702fdc1ef810bb94533489b89a956f09ef4))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))
* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))
* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))
* **main:** release 0.26.0 ([#841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/841)) ([4a6d659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a6d659c7069c1d2d64c43ce3262d3a7a840b7c5))
* **main:** release 0.26.1 ([#849](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/849)) ([a2607e5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a2607e5e15d6dfd66e756e381c0aeccf8106bbd4))
* **main:** release 0.26.2 ([#851](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/851)) ([016e02d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/016e02d3cc51360ecae43df6a931ada2b398e424))
* **main:** release 0.26.3 ([#854](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/854)) ([63f2b85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f2b8507930b9429ebf7c8ee8f65334ef16931e))
* **main:** release 0.27.0 ([#870](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/870)) ([5178aa6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5178aa6a977fe296bc4b5abeae6e55ca27291dca))
* **main:** release 0.28.0 ([#886](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/886)) ([b523f7e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b523f7e263f988a839528bb19b04405890013879))
* **main:** release 0.28.1 ([#895](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/895)) ([c92c518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c92c5184eab43141116d760f9e336714eb535fd7))
* **main:** release 0.28.2 ([#902](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/902)) ([e1c228b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1c228b9bbba60d297082b665159ca9160607e62))
* **main:** release 0.28.3 ([#905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/905)) ([b01a21b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b01a21bc7fa2055bbaf77af8e263e69fbb1bfd54))
* **main:** release 0.28.4 ([#913](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/913)) ([4fa19e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4fa19e37edda8c3958232c647d35bf99a7d00319))
* **main:** release 0.28.5 ([#915](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/915)) ([d9a7f01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d9a7f0165cc440cc7ed6086d033ab7252e56c5e2))
* **main:** release 0.28.6 ([#920](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/920)) ([3a17e34](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3a17e34a9e762ee24d8ff12144fe6c6ac4b68e2e))
* **main:** release 0.28.7 ([#921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/921)) ([adbb52e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/adbb52e3f33c86519ed20f490bddd84980437e3f))
* **main:** release 0.28.8 ([#928](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/928)) ([96152d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/96152d7db14c397bff9f9cc79ba0d85f6f2706b4))
* **main:** release 0.29.0 ([#943](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/943)) ([f1d0af9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1d0af96bbe8e57558bd3a57808298d8d49fe349))
* **main:** release 0.30.0 ([#954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/954)) ([bfd3108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bfd31080b96af02f908e93ff0c20b8cb24840431))
* **main:** release 0.31.0 ([#968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/968)) ([1e21100](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e2110065c23d851e04cd2de1727b683a38168f1))
* **main:** release 0.32.0 ([#974](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/974)) ([e947417](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e947417c459e424829a9b9e4cbb96f04ba7db3cd))
* **main:** release 0.33.0 ([#988](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/988)) ([bf3482e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bf3482e4de81e96b31aec192a15f5bee33d34e78))
* **main:** release 0.33.1 ([#991](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/991)) ([1c5af87](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c5af874164d8b40e7cae54e9206ec6b46c2e75b))
* **main:** release 0.34.0 ([#1014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1014)) ([f1c651e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1c651e17d1697f37be43857318573cb56812f5d))
* **main:** release 0.34.0 ([#1019](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1019)) ([83db3a4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83db3a4c14ec6f1539fbef55c72ae36b22e76906))
* **main:** release 0.34.0 ([#1020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1020)) ([7116025](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7116025e3523cc6d385752f7e71bff1b5fded68b))
* **main:** release 0.34.0 ([#1022](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1022)) ([d06c91f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d06c91fdacbc223cac709743a0fbe9d2c340da83))
* **main:** release 0.34.0 ([#1332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1332)) ([7037952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7037952180309441ac865eed0bc2a44a714b484d))
* **main:** release 0.34.0 ([#1436](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1436)) ([7358fdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7358fdd94a3b106a13dd7000b3c6a8f1272cf233))
* **main:** release 0.34.0 ([#1662](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1662)) ([129e4dd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/129e4ddbc7424306d75298486c1084a27f2a1807))
* **main:** release 0.35.0 ([#1026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1026)) ([f9036e8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f9036e8914b9c139eb6798276124c5544a083eb8))
* **main:** release 0.36.0 ([#1056](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1056)) ([d055d4c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d055d4c57f9a48855382506a313a4f6386da2e3e))
* **main:** release 0.37.0 ([#1065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1065)) ([6aecc46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6aecc46ddc0804a3a8b90422dfeb4c3bfbf093e5))
* **main:** release 0.37.1 ([#1096](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1096)) ([1de53b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1de53b5ee9247216b547398c29c22956247c0563))
* **main:** release 0.38.0 ([#1103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1103)) ([aee8431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aee8431ea64f085de0f4e9cfd46f2b82d16f09e2))
* **main:** release 0.39.0 ([#1130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1130)) ([82616e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82616e325890613d4b2eca5ef6ffa2e3b50a0352))
* **main:** release 0.40.0 ([#1132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1132)) ([f3f1f3b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3f1f3b517963c544da1a64d8d778c118a502b29))
* **main:** release 0.41.0 ([#1157](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1157)) ([5b9b47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b9b47d6fa2da7cd6d4b0bfe1722794003a5fce5))
* **main:** release 0.42.0 ([#1179](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1179)) ([ba45fc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba45fc27b7e3d2eda70966a857ebcd37964a5813))
* **main:** release 0.42.1 ([#1191](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1191)) ([7f9a3c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f9a3c2dd172fa93d1d2648f13b77b1f8f7981f0))
* **main:** release 0.43.0 ([#1196](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1196)) ([3ac84ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3ac84ab0834d3ab875d078489a2d2b7a45cfad28))
* **main:** release 0.43.1 ([#1207](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1207)) ([e61c15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e61c15aab3991e9740da365ec739f0c03fbbbf65))
* **main:** release 0.44.0 ([#1222](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1222)) ([1852308](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/185230847b7179079c718078780d240a9c29bbb0))
* **main:** release 0.45.0 ([#1232](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1232)) ([da886d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da886d4e05f7bb9443168f0fa04b8b397a1db5c7))
* **main:** release 0.46.0 ([#1244](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1244)) ([b9bf009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9bf009a11a7af0413c8f182927731f55379dff4))
* **main:** release 0.47.0 ([#1259](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1259)) ([887297f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/887297fc5670b180f3d7158d3092ad035fb473e9))
* **main:** release 0.48.0 ([#1284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1284)) ([cf6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf6e54f720dd852c1663a4e9ff8a74054f51325b))
* **main:** release 0.49.0 ([#1303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1303)) ([fb90556](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb90556c324ffc14b6e90adbdf9a06705af8e7e9))
* **main:** release 0.49.1 ([#1319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1319)) ([431b8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/431b8b7818cd7eccb3dafb11612f72ce8e73b58f))
* **main:** release 0.49.2 ([#1323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1323)) ([c19f307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c19f3070b8aa063c96e1e30a1e6d754b7070d296))
* **main:** release 0.49.3 ([#1327](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1327)) ([102ed1d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/102ed1de7f4fca659869fc0485b42843b394d7e9))
* **main:** release 0.50.0 ([#1344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1344)) ([a860a76](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a860a7623b9e22433ece8cede537c187a45b4bc2))
* **main:** release 0.51.0 ([#1348](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1348)) ([2b273f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b273f7e3baaf855ed6e02a7779022f38ade6745))
* **main:** release 0.52.0 ([#1363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1363)) ([e122715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1227159be50bf26841acead8730dad516a96ebc))
* **main:** release 0.53.0 ([#1401](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1401)) ([80488da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/80488dae4e16f5c55f913449fc729fbd6e1fd6d2))
* **main:** release 0.53.1 ([#1406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1406)) ([8f5ac41](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5ac41265bc08256630b2d95fa8845249098310))
* **main:** release 0.54.0 ([#1431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1431)) ([6b6b55d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b6b55d88a875f30395f2bd3250a2af1b99f9205))
* **main:** release 0.55.0 ([#1449](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1449)) ([1a00052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1a0005296689ad3ae45e5fd92b06e25ed16232de))
* **main:** release 0.55.1 ([#1469](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1469)) ([509ce3f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/509ce3f168d977de71758518e99ce0e38ab9f875))
* **main:** release 0.56.0 ([#1493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1493)) ([9a5fc2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a5fc2c0fdf993285bae42efb83b3384085540a0))
* **main:** release 0.56.1 ([#1504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1504)) ([00fc00c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/00fc00c46f22984e02ed10acdc8041cfc79b507d))
* **main:** release 0.56.2 ([#1505](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1505)) ([f950dac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f950dac5d13516075c416f6abc6d7667474a36a8))
* **main:** release 0.56.3 ([#1511](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1511)) ([9c69643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c69643a31d91d0f3d249f7aea3beeefc53880ae))
* **main:** release 0.56.4 ([#1519](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1519)) ([d0384b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d0384b6d3bfc1bc358f39e58f136c1acef452456))
* **main:** release 0.56.5 ([#1555](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1555)) ([41663ee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/41663ee5900206a03c62e046bfb9659092197bd5))
* **main:** release 0.57.0 ([#1570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1570)) ([44b96cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44b96cf67813f45feb67da4367936748bc04391f))
* **main:** release 0.58.0 ([#1587](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1587)) ([6b20b8d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b20b8d848620a7e9796ae230f6f87300e3fc50c))
* **main:** release 0.58.1 ([#1616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1616)) ([4780ba0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4780ba08b1bdf15785be63ec8dd488a03ddfe378))
* **main:** release 0.58.2 ([#1621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1621)) ([1c34ac1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c34ac157bc064d5d6fe5297231ce87eccbcc298))
* **main:** release 0.59.0 ([#1622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1622)) ([afb18aa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afb18aa8ed3c3f80630bc2f824bb756ddb5eda86))
* **main:** release 0.60.0 ([#1641](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1641)) ([ab4d49f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab4d49f259db99c2c0c6131143c55ca11d2a6610))
* **main:** release 0.60.0 ([#1665](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1665)) ([ea23020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea23020801ea4d43f055f2b443400385d96a135b))
* **main:** release 0.60.0 ([#1667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1667)) ([9d3e40f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d3e40fab05023aff16795266ec8a30761560c26))
* **main:** release 0.60.1 ([#1649](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1649)) ([56a9b2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56a9b2e5747bffb2456ad2a556e226e8450c242e))
* **main:** release 0.61.0 ([#1655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1655)) ([2fbe15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2fbe15a65a64adb8604d301e9a6d11632b6e3a44))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.0...v0.34.0) (2023-03-28)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* add GRANT ... ON ALL TABLES IN ... ([#1626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1626)) ([505a5f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/505a5f35d9ea8388ca33e5117c545408243298ae))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* add in more functionality for UpdateResourceMonitor  ([#1456](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1456)) ([2df570f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2df570f0c3271534a37b0cb61b7f4b491081baf7))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* add missing PrivateLink URLs to datasource ([#1603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1603)) ([78782b1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78782b1b471b7fbd434de1803cd687f6866cada7))
* add new account resource ([#1492](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1492)) ([b1473ba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1473ba158946d81bc4eac095c40c8d0446cf2ed))
* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* add ON STAGE support for Stream resource ([#1413](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1413)) ([447febf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/447febfef46ef89570108d3447998d6d379b7be7))
* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))
* add PREVENT_LOAD_FROM_INLINE_URL ([#1612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1612)) ([4945a3a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4945a3ae62d887dae6332742edcde715751459b5))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* Add support for `packages`, `imports`, `handler` and `runtimeVersion` to `snowflake_procedure` resource ([#1516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1516)) ([a88f3ad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a88f3ada75dad18b7b4b9024f664de8d687f54e0))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for is_secure to snowflake_function resource ([#1575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1575)) ([c41b6a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c41b6a35271f12c97f5a4da947eb433013f2aaf2))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Add title lint ([#570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/570)) ([d2142fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2142fd408f158a68230f0188c35c7b322c70ab7))
* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding slack bot for PRs and Issues ([#1106](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1106)) ([95c255e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/95c255e5ca65b619b35692671848877793cac29e))
* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))
* **integration:** add google api integration ([#1589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1589)) ([56909cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56909cdc18245f38b0f58bceaf2aa9cbc013d212))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* roles support numbers ([#1585](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1585)) ([d72dee8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d72dee82d0e0a4d8b484e5b204e156a13117cb76))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* support object parameters on account level ([#1583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1583)) ([fb24802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb2480214c8ac4e61fffe3a8e3448597462ad9a1))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))
* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))
* **deps:** bump actions/setup-go from 3 to 4 ([#1634](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1634)) ([3f128c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f128c1ba887c377b7bd5f3d508d7b81676fdf90))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1639)) ([330777e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/330777eb0ad93acede6ffef9d7571c0989540657))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1638)) ([107bb4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/107bb4abfb5de896acc1f224afae77b8100ffc02))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1454)) ([ed6bfe0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed6bfe07930e5703036ada816845176d46f5623c))
* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#1528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1528)) ([8a011e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a011e0b1920833c77eb7832f821a4bd52176657))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#1551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1551)) ([35de62f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/35de62f5b722c1dc6eaf2f39f6699935f67557cd))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump golang.org/x/tools from 0.4.0 to 0.5.0 ([#1455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1455)) ([ff01970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff019702fdc1ef810bb94533489b89a956f09ef4))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))
* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))
* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))
* **main:** release 0.26.0 ([#841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/841)) ([4a6d659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a6d659c7069c1d2d64c43ce3262d3a7a840b7c5))
* **main:** release 0.26.1 ([#849](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/849)) ([a2607e5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a2607e5e15d6dfd66e756e381c0aeccf8106bbd4))
* **main:** release 0.26.2 ([#851](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/851)) ([016e02d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/016e02d3cc51360ecae43df6a931ada2b398e424))
* **main:** release 0.26.3 ([#854](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/854)) ([63f2b85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f2b8507930b9429ebf7c8ee8f65334ef16931e))
* **main:** release 0.27.0 ([#870](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/870)) ([5178aa6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5178aa6a977fe296bc4b5abeae6e55ca27291dca))
* **main:** release 0.28.0 ([#886](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/886)) ([b523f7e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b523f7e263f988a839528bb19b04405890013879))
* **main:** release 0.28.1 ([#895](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/895)) ([c92c518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c92c5184eab43141116d760f9e336714eb535fd7))
* **main:** release 0.28.2 ([#902](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/902)) ([e1c228b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1c228b9bbba60d297082b665159ca9160607e62))
* **main:** release 0.28.3 ([#905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/905)) ([b01a21b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b01a21bc7fa2055bbaf77af8e263e69fbb1bfd54))
* **main:** release 0.28.4 ([#913](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/913)) ([4fa19e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4fa19e37edda8c3958232c647d35bf99a7d00319))
* **main:** release 0.28.5 ([#915](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/915)) ([d9a7f01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d9a7f0165cc440cc7ed6086d033ab7252e56c5e2))
* **main:** release 0.28.6 ([#920](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/920)) ([3a17e34](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3a17e34a9e762ee24d8ff12144fe6c6ac4b68e2e))
* **main:** release 0.28.7 ([#921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/921)) ([adbb52e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/adbb52e3f33c86519ed20f490bddd84980437e3f))
* **main:** release 0.28.8 ([#928](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/928)) ([96152d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/96152d7db14c397bff9f9cc79ba0d85f6f2706b4))
* **main:** release 0.29.0 ([#943](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/943)) ([f1d0af9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1d0af96bbe8e57558bd3a57808298d8d49fe349))
* **main:** release 0.30.0 ([#954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/954)) ([bfd3108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bfd31080b96af02f908e93ff0c20b8cb24840431))
* **main:** release 0.31.0 ([#968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/968)) ([1e21100](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e2110065c23d851e04cd2de1727b683a38168f1))
* **main:** release 0.32.0 ([#974](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/974)) ([e947417](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e947417c459e424829a9b9e4cbb96f04ba7db3cd))
* **main:** release 0.33.0 ([#988](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/988)) ([bf3482e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bf3482e4de81e96b31aec192a15f5bee33d34e78))
* **main:** release 0.33.1 ([#991](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/991)) ([1c5af87](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c5af874164d8b40e7cae54e9206ec6b46c2e75b))
* **main:** release 0.34.0 ([#1014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1014)) ([f1c651e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1c651e17d1697f37be43857318573cb56812f5d))
* **main:** release 0.34.0 ([#1019](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1019)) ([83db3a4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83db3a4c14ec6f1539fbef55c72ae36b22e76906))
* **main:** release 0.34.0 ([#1020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1020)) ([7116025](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7116025e3523cc6d385752f7e71bff1b5fded68b))
* **main:** release 0.34.0 ([#1022](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1022)) ([d06c91f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d06c91fdacbc223cac709743a0fbe9d2c340da83))
* **main:** release 0.34.0 ([#1332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1332)) ([7037952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7037952180309441ac865eed0bc2a44a714b484d))
* **main:** release 0.34.0 ([#1436](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1436)) ([7358fdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7358fdd94a3b106a13dd7000b3c6a8f1272cf233))
* **main:** release 0.34.0 ([#1662](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1662)) ([129e4dd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/129e4ddbc7424306d75298486c1084a27f2a1807))
* **main:** release 0.35.0 ([#1026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1026)) ([f9036e8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f9036e8914b9c139eb6798276124c5544a083eb8))
* **main:** release 0.36.0 ([#1056](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1056)) ([d055d4c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d055d4c57f9a48855382506a313a4f6386da2e3e))
* **main:** release 0.37.0 ([#1065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1065)) ([6aecc46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6aecc46ddc0804a3a8b90422dfeb4c3bfbf093e5))
* **main:** release 0.37.1 ([#1096](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1096)) ([1de53b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1de53b5ee9247216b547398c29c22956247c0563))
* **main:** release 0.38.0 ([#1103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1103)) ([aee8431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aee8431ea64f085de0f4e9cfd46f2b82d16f09e2))
* **main:** release 0.39.0 ([#1130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1130)) ([82616e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82616e325890613d4b2eca5ef6ffa2e3b50a0352))
* **main:** release 0.40.0 ([#1132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1132)) ([f3f1f3b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3f1f3b517963c544da1a64d8d778c118a502b29))
* **main:** release 0.41.0 ([#1157](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1157)) ([5b9b47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b9b47d6fa2da7cd6d4b0bfe1722794003a5fce5))
* **main:** release 0.42.0 ([#1179](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1179)) ([ba45fc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba45fc27b7e3d2eda70966a857ebcd37964a5813))
* **main:** release 0.42.1 ([#1191](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1191)) ([7f9a3c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f9a3c2dd172fa93d1d2648f13b77b1f8f7981f0))
* **main:** release 0.43.0 ([#1196](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1196)) ([3ac84ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3ac84ab0834d3ab875d078489a2d2b7a45cfad28))
* **main:** release 0.43.1 ([#1207](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1207)) ([e61c15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e61c15aab3991e9740da365ec739f0c03fbbbf65))
* **main:** release 0.44.0 ([#1222](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1222)) ([1852308](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/185230847b7179079c718078780d240a9c29bbb0))
* **main:** release 0.45.0 ([#1232](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1232)) ([da886d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da886d4e05f7bb9443168f0fa04b8b397a1db5c7))
* **main:** release 0.46.0 ([#1244](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1244)) ([b9bf009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9bf009a11a7af0413c8f182927731f55379dff4))
* **main:** release 0.47.0 ([#1259](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1259)) ([887297f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/887297fc5670b180f3d7158d3092ad035fb473e9))
* **main:** release 0.48.0 ([#1284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1284)) ([cf6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf6e54f720dd852c1663a4e9ff8a74054f51325b))
* **main:** release 0.49.0 ([#1303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1303)) ([fb90556](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb90556c324ffc14b6e90adbdf9a06705af8e7e9))
* **main:** release 0.49.1 ([#1319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1319)) ([431b8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/431b8b7818cd7eccb3dafb11612f72ce8e73b58f))
* **main:** release 0.49.2 ([#1323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1323)) ([c19f307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c19f3070b8aa063c96e1e30a1e6d754b7070d296))
* **main:** release 0.49.3 ([#1327](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1327)) ([102ed1d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/102ed1de7f4fca659869fc0485b42843b394d7e9))
* **main:** release 0.50.0 ([#1344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1344)) ([a860a76](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a860a7623b9e22433ece8cede537c187a45b4bc2))
* **main:** release 0.51.0 ([#1348](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1348)) ([2b273f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b273f7e3baaf855ed6e02a7779022f38ade6745))
* **main:** release 0.52.0 ([#1363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1363)) ([e122715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1227159be50bf26841acead8730dad516a96ebc))
* **main:** release 0.53.0 ([#1401](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1401)) ([80488da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/80488dae4e16f5c55f913449fc729fbd6e1fd6d2))
* **main:** release 0.53.1 ([#1406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1406)) ([8f5ac41](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5ac41265bc08256630b2d95fa8845249098310))
* **main:** release 0.54.0 ([#1431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1431)) ([6b6b55d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b6b55d88a875f30395f2bd3250a2af1b99f9205))
* **main:** release 0.55.0 ([#1449](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1449)) ([1a00052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1a0005296689ad3ae45e5fd92b06e25ed16232de))
* **main:** release 0.55.1 ([#1469](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1469)) ([509ce3f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/509ce3f168d977de71758518e99ce0e38ab9f875))
* **main:** release 0.56.0 ([#1493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1493)) ([9a5fc2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a5fc2c0fdf993285bae42efb83b3384085540a0))
* **main:** release 0.56.1 ([#1504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1504)) ([00fc00c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/00fc00c46f22984e02ed10acdc8041cfc79b507d))
* **main:** release 0.56.2 ([#1505](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1505)) ([f950dac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f950dac5d13516075c416f6abc6d7667474a36a8))
* **main:** release 0.56.3 ([#1511](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1511)) ([9c69643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c69643a31d91d0f3d249f7aea3beeefc53880ae))
* **main:** release 0.56.4 ([#1519](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1519)) ([d0384b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d0384b6d3bfc1bc358f39e58f136c1acef452456))
* **main:** release 0.56.5 ([#1555](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1555)) ([41663ee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/41663ee5900206a03c62e046bfb9659092197bd5))
* **main:** release 0.57.0 ([#1570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1570)) ([44b96cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44b96cf67813f45feb67da4367936748bc04391f))
* **main:** release 0.58.0 ([#1587](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1587)) ([6b20b8d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b20b8d848620a7e9796ae230f6f87300e3fc50c))
* **main:** release 0.58.1 ([#1616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1616)) ([4780ba0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4780ba08b1bdf15785be63ec8dd488a03ddfe378))
* **main:** release 0.58.2 ([#1621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1621)) ([1c34ac1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c34ac157bc064d5d6fe5297231ce87eccbcc298))
* **main:** release 0.59.0 ([#1622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1622)) ([afb18aa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afb18aa8ed3c3f80630bc2f824bb756ddb5eda86))
* **main:** release 0.60.0 ([#1641](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1641)) ([ab4d49f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab4d49f259db99c2c0c6131143c55ca11d2a6610))
* **main:** release 0.60.0 ([#1665](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1665)) ([ea23020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea23020801ea4d43f055f2b443400385d96a135b))
* **main:** release 0.60.1 ([#1649](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1649)) ([56a9b2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56a9b2e5747bffb2456ad2a556e226e8450c242e))
* **main:** release 0.61.0 ([#1655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1655)) ([2fbe15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2fbe15a65a64adb8604d301e9a6d11632b6e3a44))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))


### BugFixes

* 0.54  ([#1435](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1435)) ([4c9dd13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c9dd133574b08d8e67444b6c6b81aa87d9a2acf))
* 0.55 fix ([#1465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1465)) ([8cb3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8cb337048ec5c4a52245feb1b9556fd845d83278))
* 0.59 release fixes ([#1636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1636)) ([0a0256e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a0256ed3f0d56a6c7c22f810419636685094135))
* 0.60 misc bug fixes / linting ([#1643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1643)) ([53da853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53da853c213eec3afbdd2a47a8de3fba897c5d8a))
* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add contributing section to readme ([#1560](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1560)) ([174355d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/174355d740e325ae05435bbbc22b8b335f94fc6f))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* add nill check for grant_helpers ([#1518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1518)) ([87689bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/87689bb5b60c73bfe3d741c3da6f4f544f16aa45))
* add permissions ([#1464](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1464)) ([e2d249a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d249ae1466e05dad2080f05123e0de66fabcf6))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* add test cases for update repl schedule on failover group ([#1578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1578)) ([ab638f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab638f0b9ba866d22c6f807743eb4eccad2530b8))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* added force_new option to role grant when the role_name has been changed ([#1591](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1591)) ([4ec3613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ec3613de43d70f40a5d29ce5517af53e8ef0a06))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* all-grants ([#1658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1658)) ([d5d59b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5d59b4e85cd2e97ea0dc42b5ab2955ef35bbb88))
* Allow creation of database-wide future external table grants ([#1041](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1041)) ([5dff645](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dff645291885cd437e341148c0629fe7ab7383f))
* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* allow read of really old grant ids and add test cases ([#1615](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1615)) ([cda40ec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cda40ece534cdc3f6849a7d24f2f8acea8476e69))
* backwards compatability for grant helpers id used by procedure and functions ([#1508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1508)) ([3787657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3787657105fbcf18368136813afd558251f92cd1))
* change resource monitor suspend properties to number ([#1545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1545)) ([4bc59e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4bc59e24677260dae94952bdbc5176ad177876dd))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* cleanup grant logic ([#1522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1522)) ([0502c61](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0502c61e7211253d029a0bec6a8104738624f243))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* default_secondary_roles doc ([#1584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1584)) ([23b64fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23b64fa9abcafb59610a77cafbda11a7e2ad648c))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* do not set query_acceleration_max_scale_factor when enable enable_query_acceleration = false ([#1474](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1474)) ([d62b1b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d62b1b4d6352e7d2dc99e4603370a1f3de3a4624))
* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))
* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))
* future read on grants ([#1520](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1520)) ([db78f64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/db78f64e56d228f3236b6bdefbe9a9c18c8641e1))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* integration errors ([#1623](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1623)) ([83a40d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83a40d6361be0685b3864a0f3994298f3991de21))
* interval for failover groups ([#1448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1448)) ([bd1d3cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd1d3cc57f72c7774715f1d92a955536d55fb758))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))
* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* misc linting changes for 0.56.2 ([#1509](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1509)) ([e0d1ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d1ef5c718f9e1e58e80d31bbe2d2f27afec486))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))
* multiple share grants ([#1510](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1510)) ([d501226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d501226bc2ee8274446efb282c2dfea9599a3c2e))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* oauth integration ([#1315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1315)) ([9087220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9087220af85f08880f7ad453cbe9d13dd3bc11ec))
* openbsd build ([#1647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1647)) ([6895a89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6895a8958775e8e84a1457722f6c282d49458f3d))
* OSCP -&gt; OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* procedure and function grants ([#1502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1502)) ([0d08ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d08ea85541ceff6e591d34a671b44ef778a6611))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))
* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))
* read Database and Schema name during Stream import ([#732](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/732)) ([9f747b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f747b571b2fcf5b0663696efd75c55a6f8b6a89))
* read Name, Database and Schema during Procedure import ([#819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/819)) ([d17656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d17656fdd2803516b6ee067a6bd5a457bf31d905))
* readded imported privileges special case for database grants ([#1597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1597)) ([711ac0c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/711ac0cbc886bf8be6a5a2651234df778452b9df))
* Recreate notification integration when type changes ([#792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/792)) ([e9768bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9768bf059268fb933ad74f2b459e91e2c0563e0))
* refactor for simplify handling error ([#1472](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1472)) ([3937216](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/393721607c9eee5d73e14c27265eb39f195ccb37))
* refactor handling error to be simple ([#1473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1473)) ([9f37b99](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f37b997de073f01b66c86820237eff8049346ba))
* refactor ReadWarehouse function to correctly read object parameters ([#745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/745)) ([d83c499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83c499910c0f2b6348191a93f917e450b9e69b2))
* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))
* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))
* remove emojis, misc grant id fix ([#1598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1598)) ([fdefbc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdefbc3f1cc5bc7063f1cb1cc922854e8f9914e6))
* Remove force_new from masking_expression ([#588](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/588)) ([fc3e78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3e78acbdda5346f32a004711d31ad6f68529ed))
* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* Removing force new and adding update for data base replication config ([#1105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1105)) ([f34f012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f34f012195d0b9718904ffa7a3a529f58167a74e))
* resource snowflake_resource_monitor behavior conflict from provider 0.54.0 to 0.55.0 ([#1468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1468)) ([8ce0c53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ce0c533ec5d81273df20be2126b278ca61a59f6))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* saml integration test ([#1494](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1494)) ([8c31439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c31439253d25aafb54fc09d89e547fa8238258c))
* saml2_sign_request and saml2_force_authn cast type ([#1452](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1452)) ([f8cecd7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8cecd7ca45aabec78fd18d8aa220db7eb34b9e0))
* schema name is optional for future file_format_grant ([#1484](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1484)) ([1450cdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1450cddde6328264f9df37e4dd89a78f5f095b2e))
* schema name is optional for future function_grant ([#1485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1485)) ([dcc550e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dcc550ed5b3df548d5d300cd2b77907ea544bb43))
* schema name is optional for future procedure_grant ([#1486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1486)) ([4cf4561](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cf456151d83cd71a3b9e68abe9c6f29804f2ee2))
* schema name is optional for future sequence_grant ([#1487](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1487)) ([ccf9e78](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccf9e78c9a7884e5beea233dd529a5134c741fb1))
* schema name is optional for future snowflake_stage_grant ([#1466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1466)) ([0b4d814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b4d8146910e8ea31d2ed5ea8b58725449205dcd))
* schema name is optional for future stream_grant ([#1488](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1488)) ([3f7e5d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f7e5d655ed5738107536c873dd11533573bba46))
* schema name is optional for future task_grant ([#1489](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1489)) ([4096fd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4096fd0d8bc65ae23b6d588385e1f81c4f2e7521))
* schema read now checks first if the corresponding database exists ([#1568](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1568)) ([368dc8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368dc8fb3f7e5156d16caed1e03792654d49f3d4))
* schema_name is optional to enable future pipe grant ([#1424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1424)) ([5d966fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d966fd8624fa3208cebae3d7b32c1b59bdcfd4c))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* support different tag association queries for COLUMN object types ([#1380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1380)) ([546d0a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/546d0a144e77c759cd6ddb91a193253f27f8fb91))
* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))
* update 0.58.2 ([#1620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1620)) ([f1eab04](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1eab04dfdc839144057807953062b3591e6eaf0))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))
* update packages ([#1619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1619)) ([79a3acc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79a3acc0e3d6a405593b5adf90a31afef81d700f))
* update read role grants to use new builder ([#1596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1596)) ([e91860a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e91860ae794b034158b71ffb31097e73d8015c51))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))
* Updating shares to disallow account locators ([#1102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1102)) ([4079080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4079080dd0b9e3caf4b5d360000bd216906cb81e))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))
* wildcards in database name ([#1666](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1666)) ([54bf74c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/54bf74ca3d0119d31668d18bd1599ed029e386c8))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.0...v0.34.0) (2023-03-28)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* add GRANT ... ON ALL TABLES IN ... ([#1626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1626)) ([505a5f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/505a5f35d9ea8388ca33e5117c545408243298ae))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* add in more functionality for UpdateResourceMonitor  ([#1456](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1456)) ([2df570f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2df570f0c3271534a37b0cb61b7f4b491081baf7))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* add missing PrivateLink URLs to datasource ([#1603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1603)) ([78782b1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78782b1b471b7fbd434de1803cd687f6866cada7))
* add new account resource ([#1492](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1492)) ([b1473ba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1473ba158946d81bc4eac095c40c8d0446cf2ed))
* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* add ON STAGE support for Stream resource ([#1413](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1413)) ([447febf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/447febfef46ef89570108d3447998d6d379b7be7))
* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))
* add PREVENT_LOAD_FROM_INLINE_URL ([#1612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1612)) ([4945a3a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4945a3ae62d887dae6332742edcde715751459b5))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* Add support for `packages`, `imports`, `handler` and `runtimeVersion` to `snowflake_procedure` resource ([#1516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1516)) ([a88f3ad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a88f3ada75dad18b7b4b9024f664de8d687f54e0))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for is_secure to snowflake_function resource ([#1575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1575)) ([c41b6a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c41b6a35271f12c97f5a4da947eb433013f2aaf2))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Add title lint ([#570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/570)) ([d2142fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2142fd408f158a68230f0188c35c7b322c70ab7))
* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding slack bot for PRs and Issues ([#1106](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1106)) ([95c255e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/95c255e5ca65b619b35692671848877793cac29e))
* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))
* **integration:** add google api integration ([#1589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1589)) ([56909cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56909cdc18245f38b0f58bceaf2aa9cbc013d212))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* roles support numbers ([#1585](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1585)) ([d72dee8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d72dee82d0e0a4d8b484e5b204e156a13117cb76))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* support object parameters on account level ([#1583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1583)) ([fb24802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb2480214c8ac4e61fffe3a8e3448597462ad9a1))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))
* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### BugFixes

* 0.54  ([#1435](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1435)) ([4c9dd13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c9dd133574b08d8e67444b6c6b81aa87d9a2acf))
* 0.55 fix ([#1465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1465)) ([8cb3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8cb337048ec5c4a52245feb1b9556fd845d83278))
* 0.59 release fixes ([#1636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1636)) ([0a0256e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a0256ed3f0d56a6c7c22f810419636685094135))
* 0.60 misc bug fixes / linting ([#1643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1643)) ([53da853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53da853c213eec3afbdd2a47a8de3fba897c5d8a))
* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add contributing section to readme ([#1560](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1560)) ([174355d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/174355d740e325ae05435bbbc22b8b335f94fc6f))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* add nill check for grant_helpers ([#1518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1518)) ([87689bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/87689bb5b60c73bfe3d741c3da6f4f544f16aa45))
* add permissions ([#1464](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1464)) ([e2d249a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d249ae1466e05dad2080f05123e0de66fabcf6))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* add test cases for update repl schedule on failover group ([#1578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1578)) ([ab638f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab638f0b9ba866d22c6f807743eb4eccad2530b8))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* added force_new option to role grant when the role_name has been changed ([#1591](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1591)) ([4ec3613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ec3613de43d70f40a5d29ce5517af53e8ef0a06))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* all-grants ([#1658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1658)) ([d5d59b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5d59b4e85cd2e97ea0dc42b5ab2955ef35bbb88))
* Allow creation of database-wide future external table grants ([#1041](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1041)) ([5dff645](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dff645291885cd437e341148c0629fe7ab7383f))
* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* allow read of really old grant ids and add test cases ([#1615](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1615)) ([cda40ec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cda40ece534cdc3f6849a7d24f2f8acea8476e69))
* backwards compatability for grant helpers id used by procedure and functions ([#1508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1508)) ([3787657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3787657105fbcf18368136813afd558251f92cd1))
* change resource monitor suspend properties to number ([#1545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1545)) ([4bc59e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4bc59e24677260dae94952bdbc5176ad177876dd))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* cleanup grant logic ([#1522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1522)) ([0502c61](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0502c61e7211253d029a0bec6a8104738624f243))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* default_secondary_roles doc ([#1584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1584)) ([23b64fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23b64fa9abcafb59610a77cafbda11a7e2ad648c))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* do not set query_acceleration_max_scale_factor when enable enable_query_acceleration = false ([#1474](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1474)) ([d62b1b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d62b1b4d6352e7d2dc99e4603370a1f3de3a4624))
* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))
* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))
* future read on grants ([#1520](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1520)) ([db78f64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/db78f64e56d228f3236b6bdefbe9a9c18c8641e1))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* integration errors ([#1623](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1623)) ([83a40d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83a40d6361be0685b3864a0f3994298f3991de21))
* interval for failover groups ([#1448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1448)) ([bd1d3cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd1d3cc57f72c7774715f1d92a955536d55fb758))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))
* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* misc linting changes for 0.56.2 ([#1509](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1509)) ([e0d1ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d1ef5c718f9e1e58e80d31bbe2d2f27afec486))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))
* multiple share grants ([#1510](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1510)) ([d501226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d501226bc2ee8274446efb282c2dfea9599a3c2e))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* oauth integration ([#1315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1315)) ([9087220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9087220af85f08880f7ad453cbe9d13dd3bc11ec))
* openbsd build ([#1647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1647)) ([6895a89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6895a8958775e8e84a1457722f6c282d49458f3d))
* OSCP -&gt; OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* procedure and function grants ([#1502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1502)) ([0d08ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d08ea85541ceff6e591d34a671b44ef778a6611))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))
* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))
* read Database and Schema name during Stream import ([#732](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/732)) ([9f747b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f747b571b2fcf5b0663696efd75c55a6f8b6a89))
* read Name, Database and Schema during Procedure import ([#819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/819)) ([d17656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d17656fdd2803516b6ee067a6bd5a457bf31d905))
* readded imported privileges special case for database grants ([#1597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1597)) ([711ac0c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/711ac0cbc886bf8be6a5a2651234df778452b9df))
* Recreate notification integration when type changes ([#792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/792)) ([e9768bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9768bf059268fb933ad74f2b459e91e2c0563e0))
* refactor for simplify handling error ([#1472](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1472)) ([3937216](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/393721607c9eee5d73e14c27265eb39f195ccb37))
* refactor handling error to be simple ([#1473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1473)) ([9f37b99](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f37b997de073f01b66c86820237eff8049346ba))
* refactor ReadWarehouse function to correctly read object parameters ([#745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/745)) ([d83c499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83c499910c0f2b6348191a93f917e450b9e69b2))
* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))
* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))
* remove emojis, misc grant id fix ([#1598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1598)) ([fdefbc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdefbc3f1cc5bc7063f1cb1cc922854e8f9914e6))
* Remove force_new from masking_expression ([#588](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/588)) ([fc3e78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3e78acbdda5346f32a004711d31ad6f68529ed))
* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* Removing force new and adding update for data base replication config ([#1105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1105)) ([f34f012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f34f012195d0b9718904ffa7a3a529f58167a74e))
* resource snowflake_resource_monitor behavior conflict from provider 0.54.0 to 0.55.0 ([#1468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1468)) ([8ce0c53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ce0c533ec5d81273df20be2126b278ca61a59f6))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* saml integration test ([#1494](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1494)) ([8c31439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c31439253d25aafb54fc09d89e547fa8238258c))
* saml2_sign_request and saml2_force_authn cast type ([#1452](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1452)) ([f8cecd7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8cecd7ca45aabec78fd18d8aa220db7eb34b9e0))
* schema name is optional for future file_format_grant ([#1484](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1484)) ([1450cdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1450cddde6328264f9df37e4dd89a78f5f095b2e))
* schema name is optional for future function_grant ([#1485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1485)) ([dcc550e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dcc550ed5b3df548d5d300cd2b77907ea544bb43))
* schema name is optional for future procedure_grant ([#1486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1486)) ([4cf4561](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cf456151d83cd71a3b9e68abe9c6f29804f2ee2))
* schema name is optional for future sequence_grant ([#1487](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1487)) ([ccf9e78](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccf9e78c9a7884e5beea233dd529a5134c741fb1))
* schema name is optional for future snowflake_stage_grant ([#1466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1466)) ([0b4d814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b4d8146910e8ea31d2ed5ea8b58725449205dcd))
* schema name is optional for future stream_grant ([#1488](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1488)) ([3f7e5d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f7e5d655ed5738107536c873dd11533573bba46))
* schema name is optional for future task_grant ([#1489](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1489)) ([4096fd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4096fd0d8bc65ae23b6d588385e1f81c4f2e7521))
* schema read now checks first if the corresponding database exists ([#1568](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1568)) ([368dc8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368dc8fb3f7e5156d16caed1e03792654d49f3d4))
* schema_name is optional to enable future pipe grant ([#1424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1424)) ([5d966fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d966fd8624fa3208cebae3d7b32c1b59bdcfd4c))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* support different tag association queries for COLUMN object types ([#1380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1380)) ([546d0a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/546d0a144e77c759cd6ddb91a193253f27f8fb91))
* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))
* update 0.58.2 ([#1620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1620)) ([f1eab04](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1eab04dfdc839144057807953062b3591e6eaf0))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))
* update packages ([#1619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1619)) ([79a3acc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79a3acc0e3d6a405593b5adf90a31afef81d700f))
* update read role grants to use new builder ([#1596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1596)) ([e91860a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e91860ae794b034158b71ffb31097e73d8015c51))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))
* Updating shares to disallow account locators ([#1102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1102)) ([4079080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4079080dd0b9e3caf4b5d360000bd216906cb81e))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))
* **deps:** bump actions/setup-go from 3 to 4 ([#1634](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1634)) ([3f128c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f128c1ba887c377b7bd5f3d508d7b81676fdf90))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1639)) ([330777e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/330777eb0ad93acede6ffef9d7571c0989540657))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1638)) ([107bb4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/107bb4abfb5de896acc1f224afae77b8100ffc02))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1454)) ([ed6bfe0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed6bfe07930e5703036ada816845176d46f5623c))
* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#1528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1528)) ([8a011e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a011e0b1920833c77eb7832f821a4bd52176657))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#1551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1551)) ([35de62f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/35de62f5b722c1dc6eaf2f39f6699935f67557cd))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump golang.org/x/tools from 0.4.0 to 0.5.0 ([#1455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1455)) ([ff01970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff019702fdc1ef810bb94533489b89a956f09ef4))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))
* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))
* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))
* **main:** release 0.26.0 ([#841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/841)) ([4a6d659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a6d659c7069c1d2d64c43ce3262d3a7a840b7c5))
* **main:** release 0.26.1 ([#849](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/849)) ([a2607e5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a2607e5e15d6dfd66e756e381c0aeccf8106bbd4))
* **main:** release 0.26.2 ([#851](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/851)) ([016e02d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/016e02d3cc51360ecae43df6a931ada2b398e424))
* **main:** release 0.26.3 ([#854](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/854)) ([63f2b85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f2b8507930b9429ebf7c8ee8f65334ef16931e))
* **main:** release 0.27.0 ([#870](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/870)) ([5178aa6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5178aa6a977fe296bc4b5abeae6e55ca27291dca))
* **main:** release 0.28.0 ([#886](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/886)) ([b523f7e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b523f7e263f988a839528bb19b04405890013879))
* **main:** release 0.28.1 ([#895](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/895)) ([c92c518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c92c5184eab43141116d760f9e336714eb535fd7))
* **main:** release 0.28.2 ([#902](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/902)) ([e1c228b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1c228b9bbba60d297082b665159ca9160607e62))
* **main:** release 0.28.3 ([#905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/905)) ([b01a21b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b01a21bc7fa2055bbaf77af8e263e69fbb1bfd54))
* **main:** release 0.28.4 ([#913](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/913)) ([4fa19e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4fa19e37edda8c3958232c647d35bf99a7d00319))
* **main:** release 0.28.5 ([#915](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/915)) ([d9a7f01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d9a7f0165cc440cc7ed6086d033ab7252e56c5e2))
* **main:** release 0.28.6 ([#920](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/920)) ([3a17e34](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3a17e34a9e762ee24d8ff12144fe6c6ac4b68e2e))
* **main:** release 0.28.7 ([#921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/921)) ([adbb52e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/adbb52e3f33c86519ed20f490bddd84980437e3f))
* **main:** release 0.28.8 ([#928](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/928)) ([96152d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/96152d7db14c397bff9f9cc79ba0d85f6f2706b4))
* **main:** release 0.29.0 ([#943](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/943)) ([f1d0af9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1d0af96bbe8e57558bd3a57808298d8d49fe349))
* **main:** release 0.30.0 ([#954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/954)) ([bfd3108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bfd31080b96af02f908e93ff0c20b8cb24840431))
* **main:** release 0.31.0 ([#968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/968)) ([1e21100](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e2110065c23d851e04cd2de1727b683a38168f1))
* **main:** release 0.32.0 ([#974](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/974)) ([e947417](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e947417c459e424829a9b9e4cbb96f04ba7db3cd))
* **main:** release 0.33.0 ([#988](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/988)) ([bf3482e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bf3482e4de81e96b31aec192a15f5bee33d34e78))
* **main:** release 0.33.1 ([#991](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/991)) ([1c5af87](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c5af874164d8b40e7cae54e9206ec6b46c2e75b))
* **main:** release 0.34.0 ([#1014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1014)) ([f1c651e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1c651e17d1697f37be43857318573cb56812f5d))
* **main:** release 0.34.0 ([#1019](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1019)) ([83db3a4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83db3a4c14ec6f1539fbef55c72ae36b22e76906))
* **main:** release 0.34.0 ([#1020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1020)) ([7116025](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7116025e3523cc6d385752f7e71bff1b5fded68b))
* **main:** release 0.34.0 ([#1022](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1022)) ([d06c91f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d06c91fdacbc223cac709743a0fbe9d2c340da83))
* **main:** release 0.34.0 ([#1332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1332)) ([7037952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7037952180309441ac865eed0bc2a44a714b484d))
* **main:** release 0.34.0 ([#1436](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1436)) ([7358fdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7358fdd94a3b106a13dd7000b3c6a8f1272cf233))
* **main:** release 0.34.0 ([#1662](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1662)) ([129e4dd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/129e4ddbc7424306d75298486c1084a27f2a1807))
* **main:** release 0.35.0 ([#1026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1026)) ([f9036e8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f9036e8914b9c139eb6798276124c5544a083eb8))
* **main:** release 0.36.0 ([#1056](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1056)) ([d055d4c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d055d4c57f9a48855382506a313a4f6386da2e3e))
* **main:** release 0.37.0 ([#1065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1065)) ([6aecc46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6aecc46ddc0804a3a8b90422dfeb4c3bfbf093e5))
* **main:** release 0.37.1 ([#1096](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1096)) ([1de53b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1de53b5ee9247216b547398c29c22956247c0563))
* **main:** release 0.38.0 ([#1103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1103)) ([aee8431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aee8431ea64f085de0f4e9cfd46f2b82d16f09e2))
* **main:** release 0.39.0 ([#1130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1130)) ([82616e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82616e325890613d4b2eca5ef6ffa2e3b50a0352))
* **main:** release 0.40.0 ([#1132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1132)) ([f3f1f3b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3f1f3b517963c544da1a64d8d778c118a502b29))
* **main:** release 0.41.0 ([#1157](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1157)) ([5b9b47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b9b47d6fa2da7cd6d4b0bfe1722794003a5fce5))
* **main:** release 0.42.0 ([#1179](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1179)) ([ba45fc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba45fc27b7e3d2eda70966a857ebcd37964a5813))
* **main:** release 0.42.1 ([#1191](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1191)) ([7f9a3c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f9a3c2dd172fa93d1d2648f13b77b1f8f7981f0))
* **main:** release 0.43.0 ([#1196](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1196)) ([3ac84ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3ac84ab0834d3ab875d078489a2d2b7a45cfad28))
* **main:** release 0.43.1 ([#1207](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1207)) ([e61c15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e61c15aab3991e9740da365ec739f0c03fbbbf65))
* **main:** release 0.44.0 ([#1222](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1222)) ([1852308](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/185230847b7179079c718078780d240a9c29bbb0))
* **main:** release 0.45.0 ([#1232](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1232)) ([da886d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da886d4e05f7bb9443168f0fa04b8b397a1db5c7))
* **main:** release 0.46.0 ([#1244](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1244)) ([b9bf009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9bf009a11a7af0413c8f182927731f55379dff4))
* **main:** release 0.47.0 ([#1259](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1259)) ([887297f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/887297fc5670b180f3d7158d3092ad035fb473e9))
* **main:** release 0.48.0 ([#1284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1284)) ([cf6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf6e54f720dd852c1663a4e9ff8a74054f51325b))
* **main:** release 0.49.0 ([#1303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1303)) ([fb90556](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb90556c324ffc14b6e90adbdf9a06705af8e7e9))
* **main:** release 0.49.1 ([#1319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1319)) ([431b8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/431b8b7818cd7eccb3dafb11612f72ce8e73b58f))
* **main:** release 0.49.2 ([#1323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1323)) ([c19f307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c19f3070b8aa063c96e1e30a1e6d754b7070d296))
* **main:** release 0.49.3 ([#1327](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1327)) ([102ed1d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/102ed1de7f4fca659869fc0485b42843b394d7e9))
* **main:** release 0.50.0 ([#1344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1344)) ([a860a76](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a860a7623b9e22433ece8cede537c187a45b4bc2))
* **main:** release 0.51.0 ([#1348](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1348)) ([2b273f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b273f7e3baaf855ed6e02a7779022f38ade6745))
* **main:** release 0.52.0 ([#1363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1363)) ([e122715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1227159be50bf26841acead8730dad516a96ebc))
* **main:** release 0.53.0 ([#1401](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1401)) ([80488da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/80488dae4e16f5c55f913449fc729fbd6e1fd6d2))
* **main:** release 0.53.1 ([#1406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1406)) ([8f5ac41](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5ac41265bc08256630b2d95fa8845249098310))
* **main:** release 0.54.0 ([#1431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1431)) ([6b6b55d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b6b55d88a875f30395f2bd3250a2af1b99f9205))
* **main:** release 0.55.0 ([#1449](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1449)) ([1a00052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1a0005296689ad3ae45e5fd92b06e25ed16232de))
* **main:** release 0.55.1 ([#1469](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1469)) ([509ce3f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/509ce3f168d977de71758518e99ce0e38ab9f875))
* **main:** release 0.56.0 ([#1493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1493)) ([9a5fc2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a5fc2c0fdf993285bae42efb83b3384085540a0))
* **main:** release 0.56.1 ([#1504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1504)) ([00fc00c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/00fc00c46f22984e02ed10acdc8041cfc79b507d))
* **main:** release 0.56.2 ([#1505](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1505)) ([f950dac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f950dac5d13516075c416f6abc6d7667474a36a8))
* **main:** release 0.56.3 ([#1511](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1511)) ([9c69643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c69643a31d91d0f3d249f7aea3beeefc53880ae))
* **main:** release 0.56.4 ([#1519](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1519)) ([d0384b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d0384b6d3bfc1bc358f39e58f136c1acef452456))
* **main:** release 0.56.5 ([#1555](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1555)) ([41663ee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/41663ee5900206a03c62e046bfb9659092197bd5))
* **main:** release 0.57.0 ([#1570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1570)) ([44b96cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44b96cf67813f45feb67da4367936748bc04391f))
* **main:** release 0.58.0 ([#1587](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1587)) ([6b20b8d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b20b8d848620a7e9796ae230f6f87300e3fc50c))
* **main:** release 0.58.1 ([#1616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1616)) ([4780ba0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4780ba08b1bdf15785be63ec8dd488a03ddfe378))
* **main:** release 0.58.2 ([#1621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1621)) ([1c34ac1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c34ac157bc064d5d6fe5297231ce87eccbcc298))
* **main:** release 0.59.0 ([#1622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1622)) ([afb18aa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afb18aa8ed3c3f80630bc2f824bb756ddb5eda86))
* **main:** release 0.60.0 ([#1641](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1641)) ([ab4d49f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab4d49f259db99c2c0c6131143c55ca11d2a6610))
* **main:** release 0.60.1 ([#1649](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1649)) ([56a9b2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56a9b2e5747bffb2456ad2a556e226e8450c242e))
* **main:** release 0.61.0 ([#1655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1655)) ([2fbe15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2fbe15a65a64adb8604d301e9a6d11632b6e3a44))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.61.0...v0.34.0) (2023-03-28)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* add GRANT ... ON ALL TABLES IN ... ([#1626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1626)) ([505a5f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/505a5f35d9ea8388ca33e5117c545408243298ae))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* add in more functionality for UpdateResourceMonitor  ([#1456](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1456)) ([2df570f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2df570f0c3271534a37b0cb61b7f4b491081baf7))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* add missing PrivateLink URLs to datasource ([#1603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1603)) ([78782b1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78782b1b471b7fbd434de1803cd687f6866cada7))
* add new account resource ([#1492](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1492)) ([b1473ba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1473ba158946d81bc4eac095c40c8d0446cf2ed))
* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* add ON STAGE support for Stream resource ([#1413](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1413)) ([447febf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/447febfef46ef89570108d3447998d6d379b7be7))
* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))
* add PREVENT_LOAD_FROM_INLINE_URL ([#1612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1612)) ([4945a3a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4945a3ae62d887dae6332742edcde715751459b5))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* Add support for `packages`, `imports`, `handler` and `runtimeVersion` to `snowflake_procedure` resource ([#1516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1516)) ([a88f3ad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a88f3ada75dad18b7b4b9024f664de8d687f54e0))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for is_secure to snowflake_function resource ([#1575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1575)) ([c41b6a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c41b6a35271f12c97f5a4da947eb433013f2aaf2))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Add title lint ([#570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/570)) ([d2142fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2142fd408f158a68230f0188c35c7b322c70ab7))
* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding slack bot for PRs and Issues ([#1106](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1106)) ([95c255e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/95c255e5ca65b619b35692671848877793cac29e))
* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))
* **integration:** add google api integration ([#1589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1589)) ([56909cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56909cdc18245f38b0f58bceaf2aa9cbc013d212))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* roles support numbers ([#1585](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1585)) ([d72dee8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d72dee82d0e0a4d8b484e5b204e156a13117cb76))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* support object parameters on account level ([#1583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1583)) ([fb24802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb2480214c8ac4e61fffe3a8e3448597462ad9a1))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))
* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### BugFixes

* 0.54  ([#1435](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1435)) ([4c9dd13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c9dd133574b08d8e67444b6c6b81aa87d9a2acf))
* 0.55 fix ([#1465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1465)) ([8cb3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8cb337048ec5c4a52245feb1b9556fd845d83278))
* 0.59 release fixes ([#1636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1636)) ([0a0256e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a0256ed3f0d56a6c7c22f810419636685094135))
* 0.60 misc bug fixes / linting ([#1643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1643)) ([53da853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53da853c213eec3afbdd2a47a8de3fba897c5d8a))
* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add contributing section to readme ([#1560](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1560)) ([174355d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/174355d740e325ae05435bbbc22b8b335f94fc6f))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* add nill check for grant_helpers ([#1518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1518)) ([87689bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/87689bb5b60c73bfe3d741c3da6f4f544f16aa45))
* add permissions ([#1464](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1464)) ([e2d249a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d249ae1466e05dad2080f05123e0de66fabcf6))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* add test cases for update repl schedule on failover group ([#1578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1578)) ([ab638f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab638f0b9ba866d22c6f807743eb4eccad2530b8))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* added force_new option to role grant when the role_name has been changed ([#1591](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1591)) ([4ec3613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ec3613de43d70f40a5d29ce5517af53e8ef0a06))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* all-grants ([#1658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1658)) ([d5d59b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5d59b4e85cd2e97ea0dc42b5ab2955ef35bbb88))
* Allow creation of database-wide future external table grants ([#1041](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1041)) ([5dff645](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dff645291885cd437e341148c0629fe7ab7383f))
* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* allow read of really old grant ids and add test cases ([#1615](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1615)) ([cda40ec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cda40ece534cdc3f6849a7d24f2f8acea8476e69))
* backwards compatability for grant helpers id used by procedure and functions ([#1508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1508)) ([3787657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3787657105fbcf18368136813afd558251f92cd1))
* change resource monitor suspend properties to number ([#1545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1545)) ([4bc59e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4bc59e24677260dae94952bdbc5176ad177876dd))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* cleanup grant logic ([#1522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1522)) ([0502c61](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0502c61e7211253d029a0bec6a8104738624f243))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* default_secondary_roles doc ([#1584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1584)) ([23b64fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23b64fa9abcafb59610a77cafbda11a7e2ad648c))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* do not set query_acceleration_max_scale_factor when enable enable_query_acceleration = false ([#1474](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1474)) ([d62b1b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d62b1b4d6352e7d2dc99e4603370a1f3de3a4624))
* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))
* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))
* future read on grants ([#1520](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1520)) ([db78f64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/db78f64e56d228f3236b6bdefbe9a9c18c8641e1))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* integration errors ([#1623](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1623)) ([83a40d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83a40d6361be0685b3864a0f3994298f3991de21))
* interval for failover groups ([#1448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1448)) ([bd1d3cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd1d3cc57f72c7774715f1d92a955536d55fb758))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))
* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* misc linting changes for 0.56.2 ([#1509](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1509)) ([e0d1ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d1ef5c718f9e1e58e80d31bbe2d2f27afec486))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))
* multiple share grants ([#1510](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1510)) ([d501226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d501226bc2ee8274446efb282c2dfea9599a3c2e))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* oauth integration ([#1315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1315)) ([9087220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9087220af85f08880f7ad453cbe9d13dd3bc11ec))
* openbsd build ([#1647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1647)) ([6895a89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6895a8958775e8e84a1457722f6c282d49458f3d))
* OSCP -&gt; OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* procedure and function grants ([#1502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1502)) ([0d08ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d08ea85541ceff6e591d34a671b44ef778a6611))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))
* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))
* read Database and Schema name during Stream import ([#732](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/732)) ([9f747b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f747b571b2fcf5b0663696efd75c55a6f8b6a89))
* read Name, Database and Schema during Procedure import ([#819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/819)) ([d17656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d17656fdd2803516b6ee067a6bd5a457bf31d905))
* readded imported privileges special case for database grants ([#1597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1597)) ([711ac0c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/711ac0cbc886bf8be6a5a2651234df778452b9df))
* Recreate notification integration when type changes ([#792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/792)) ([e9768bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9768bf059268fb933ad74f2b459e91e2c0563e0))
* refactor for simplify handling error ([#1472](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1472)) ([3937216](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/393721607c9eee5d73e14c27265eb39f195ccb37))
* refactor handling error to be simple ([#1473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1473)) ([9f37b99](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f37b997de073f01b66c86820237eff8049346ba))
* refactor ReadWarehouse function to correctly read object parameters ([#745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/745)) ([d83c499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83c499910c0f2b6348191a93f917e450b9e69b2))
* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))
* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))
* remove emojis, misc grant id fix ([#1598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1598)) ([fdefbc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdefbc3f1cc5bc7063f1cb1cc922854e8f9914e6))
* Remove force_new from masking_expression ([#588](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/588)) ([fc3e78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3e78acbdda5346f32a004711d31ad6f68529ed))
* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* Removing force new and adding update for data base replication config ([#1105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1105)) ([f34f012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f34f012195d0b9718904ffa7a3a529f58167a74e))
* resource snowflake_resource_monitor behavior conflict from provider 0.54.0 to 0.55.0 ([#1468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1468)) ([8ce0c53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ce0c533ec5d81273df20be2126b278ca61a59f6))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* saml integration test ([#1494](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1494)) ([8c31439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c31439253d25aafb54fc09d89e547fa8238258c))
* saml2_sign_request and saml2_force_authn cast type ([#1452](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1452)) ([f8cecd7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8cecd7ca45aabec78fd18d8aa220db7eb34b9e0))
* schema name is optional for future file_format_grant ([#1484](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1484)) ([1450cdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1450cddde6328264f9df37e4dd89a78f5f095b2e))
* schema name is optional for future function_grant ([#1485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1485)) ([dcc550e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dcc550ed5b3df548d5d300cd2b77907ea544bb43))
* schema name is optional for future procedure_grant ([#1486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1486)) ([4cf4561](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cf456151d83cd71a3b9e68abe9c6f29804f2ee2))
* schema name is optional for future sequence_grant ([#1487](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1487)) ([ccf9e78](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccf9e78c9a7884e5beea233dd529a5134c741fb1))
* schema name is optional for future snowflake_stage_grant ([#1466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1466)) ([0b4d814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b4d8146910e8ea31d2ed5ea8b58725449205dcd))
* schema name is optional for future stream_grant ([#1488](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1488)) ([3f7e5d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f7e5d655ed5738107536c873dd11533573bba46))
* schema name is optional for future task_grant ([#1489](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1489)) ([4096fd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4096fd0d8bc65ae23b6d588385e1f81c4f2e7521))
* schema read now checks first if the corresponding database exists ([#1568](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1568)) ([368dc8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368dc8fb3f7e5156d16caed1e03792654d49f3d4))
* schema_name is optional to enable future pipe grant ([#1424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1424)) ([5d966fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d966fd8624fa3208cebae3d7b32c1b59bdcfd4c))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* support different tag association queries for COLUMN object types ([#1380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1380)) ([546d0a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/546d0a144e77c759cd6ddb91a193253f27f8fb91))
* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))
* update 0.58.2 ([#1620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1620)) ([f1eab04](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1eab04dfdc839144057807953062b3591e6eaf0))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))
* update packages ([#1619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1619)) ([79a3acc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79a3acc0e3d6a405593b5adf90a31afef81d700f))
* update read role grants to use new builder ([#1596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1596)) ([e91860a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e91860ae794b034158b71ffb31097e73d8015c51))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))
* Updating shares to disallow account locators ([#1102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1102)) ([4079080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4079080dd0b9e3caf4b5d360000bd216906cb81e))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))
* **deps:** bump actions/setup-go from 3 to 4 ([#1634](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1634)) ([3f128c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f128c1ba887c377b7bd5f3d508d7b81676fdf90))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1639)) ([330777e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/330777eb0ad93acede6ffef9d7571c0989540657))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1638)) ([107bb4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/107bb4abfb5de896acc1f224afae77b8100ffc02))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1454)) ([ed6bfe0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed6bfe07930e5703036ada816845176d46f5623c))
* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#1528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1528)) ([8a011e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a011e0b1920833c77eb7832f821a4bd52176657))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#1551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1551)) ([35de62f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/35de62f5b722c1dc6eaf2f39f6699935f67557cd))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump golang.org/x/tools from 0.4.0 to 0.5.0 ([#1455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1455)) ([ff01970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff019702fdc1ef810bb94533489b89a956f09ef4))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))
* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))
* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))
* **main:** release 0.26.0 ([#841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/841)) ([4a6d659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a6d659c7069c1d2d64c43ce3262d3a7a840b7c5))
* **main:** release 0.26.1 ([#849](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/849)) ([a2607e5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a2607e5e15d6dfd66e756e381c0aeccf8106bbd4))
* **main:** release 0.26.2 ([#851](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/851)) ([016e02d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/016e02d3cc51360ecae43df6a931ada2b398e424))
* **main:** release 0.26.3 ([#854](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/854)) ([63f2b85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f2b8507930b9429ebf7c8ee8f65334ef16931e))
* **main:** release 0.27.0 ([#870](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/870)) ([5178aa6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5178aa6a977fe296bc4b5abeae6e55ca27291dca))
* **main:** release 0.28.0 ([#886](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/886)) ([b523f7e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b523f7e263f988a839528bb19b04405890013879))
* **main:** release 0.28.1 ([#895](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/895)) ([c92c518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c92c5184eab43141116d760f9e336714eb535fd7))
* **main:** release 0.28.2 ([#902](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/902)) ([e1c228b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1c228b9bbba60d297082b665159ca9160607e62))
* **main:** release 0.28.3 ([#905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/905)) ([b01a21b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b01a21bc7fa2055bbaf77af8e263e69fbb1bfd54))
* **main:** release 0.28.4 ([#913](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/913)) ([4fa19e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4fa19e37edda8c3958232c647d35bf99a7d00319))
* **main:** release 0.28.5 ([#915](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/915)) ([d9a7f01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d9a7f0165cc440cc7ed6086d033ab7252e56c5e2))
* **main:** release 0.28.6 ([#920](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/920)) ([3a17e34](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3a17e34a9e762ee24d8ff12144fe6c6ac4b68e2e))
* **main:** release 0.28.7 ([#921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/921)) ([adbb52e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/adbb52e3f33c86519ed20f490bddd84980437e3f))
* **main:** release 0.28.8 ([#928](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/928)) ([96152d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/96152d7db14c397bff9f9cc79ba0d85f6f2706b4))
* **main:** release 0.29.0 ([#943](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/943)) ([f1d0af9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1d0af96bbe8e57558bd3a57808298d8d49fe349))
* **main:** release 0.30.0 ([#954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/954)) ([bfd3108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bfd31080b96af02f908e93ff0c20b8cb24840431))
* **main:** release 0.31.0 ([#968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/968)) ([1e21100](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e2110065c23d851e04cd2de1727b683a38168f1))
* **main:** release 0.32.0 ([#974](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/974)) ([e947417](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e947417c459e424829a9b9e4cbb96f04ba7db3cd))
* **main:** release 0.33.0 ([#988](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/988)) ([bf3482e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bf3482e4de81e96b31aec192a15f5bee33d34e78))
* **main:** release 0.33.1 ([#991](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/991)) ([1c5af87](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c5af874164d8b40e7cae54e9206ec6b46c2e75b))
* **main:** release 0.34.0 ([#1014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1014)) ([f1c651e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1c651e17d1697f37be43857318573cb56812f5d))
* **main:** release 0.34.0 ([#1019](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1019)) ([83db3a4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83db3a4c14ec6f1539fbef55c72ae36b22e76906))
* **main:** release 0.34.0 ([#1020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1020)) ([7116025](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7116025e3523cc6d385752f7e71bff1b5fded68b))
* **main:** release 0.34.0 ([#1022](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1022)) ([d06c91f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d06c91fdacbc223cac709743a0fbe9d2c340da83))
* **main:** release 0.34.0 ([#1332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1332)) ([7037952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7037952180309441ac865eed0bc2a44a714b484d))
* **main:** release 0.34.0 ([#1436](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1436)) ([7358fdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7358fdd94a3b106a13dd7000b3c6a8f1272cf233))
* **main:** release 0.35.0 ([#1026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1026)) ([f9036e8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f9036e8914b9c139eb6798276124c5544a083eb8))
* **main:** release 0.36.0 ([#1056](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1056)) ([d055d4c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d055d4c57f9a48855382506a313a4f6386da2e3e))
* **main:** release 0.37.0 ([#1065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1065)) ([6aecc46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6aecc46ddc0804a3a8b90422dfeb4c3bfbf093e5))
* **main:** release 0.37.1 ([#1096](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1096)) ([1de53b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1de53b5ee9247216b547398c29c22956247c0563))
* **main:** release 0.38.0 ([#1103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1103)) ([aee8431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aee8431ea64f085de0f4e9cfd46f2b82d16f09e2))
* **main:** release 0.39.0 ([#1130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1130)) ([82616e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82616e325890613d4b2eca5ef6ffa2e3b50a0352))
* **main:** release 0.40.0 ([#1132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1132)) ([f3f1f3b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3f1f3b517963c544da1a64d8d778c118a502b29))
* **main:** release 0.41.0 ([#1157](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1157)) ([5b9b47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b9b47d6fa2da7cd6d4b0bfe1722794003a5fce5))
* **main:** release 0.42.0 ([#1179](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1179)) ([ba45fc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba45fc27b7e3d2eda70966a857ebcd37964a5813))
* **main:** release 0.42.1 ([#1191](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1191)) ([7f9a3c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f9a3c2dd172fa93d1d2648f13b77b1f8f7981f0))
* **main:** release 0.43.0 ([#1196](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1196)) ([3ac84ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3ac84ab0834d3ab875d078489a2d2b7a45cfad28))
* **main:** release 0.43.1 ([#1207](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1207)) ([e61c15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e61c15aab3991e9740da365ec739f0c03fbbbf65))
* **main:** release 0.44.0 ([#1222](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1222)) ([1852308](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/185230847b7179079c718078780d240a9c29bbb0))
* **main:** release 0.45.0 ([#1232](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1232)) ([da886d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da886d4e05f7bb9443168f0fa04b8b397a1db5c7))
* **main:** release 0.46.0 ([#1244](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1244)) ([b9bf009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9bf009a11a7af0413c8f182927731f55379dff4))
* **main:** release 0.47.0 ([#1259](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1259)) ([887297f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/887297fc5670b180f3d7158d3092ad035fb473e9))
* **main:** release 0.48.0 ([#1284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1284)) ([cf6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf6e54f720dd852c1663a4e9ff8a74054f51325b))
* **main:** release 0.49.0 ([#1303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1303)) ([fb90556](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb90556c324ffc14b6e90adbdf9a06705af8e7e9))
* **main:** release 0.49.1 ([#1319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1319)) ([431b8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/431b8b7818cd7eccb3dafb11612f72ce8e73b58f))
* **main:** release 0.49.2 ([#1323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1323)) ([c19f307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c19f3070b8aa063c96e1e30a1e6d754b7070d296))
* **main:** release 0.49.3 ([#1327](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1327)) ([102ed1d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/102ed1de7f4fca659869fc0485b42843b394d7e9))
* **main:** release 0.50.0 ([#1344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1344)) ([a860a76](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a860a7623b9e22433ece8cede537c187a45b4bc2))
* **main:** release 0.51.0 ([#1348](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1348)) ([2b273f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b273f7e3baaf855ed6e02a7779022f38ade6745))
* **main:** release 0.52.0 ([#1363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1363)) ([e122715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1227159be50bf26841acead8730dad516a96ebc))
* **main:** release 0.53.0 ([#1401](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1401)) ([80488da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/80488dae4e16f5c55f913449fc729fbd6e1fd6d2))
* **main:** release 0.53.1 ([#1406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1406)) ([8f5ac41](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5ac41265bc08256630b2d95fa8845249098310))
* **main:** release 0.54.0 ([#1431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1431)) ([6b6b55d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b6b55d88a875f30395f2bd3250a2af1b99f9205))
* **main:** release 0.55.0 ([#1449](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1449)) ([1a00052](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1a0005296689ad3ae45e5fd92b06e25ed16232de))
* **main:** release 0.55.1 ([#1469](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1469)) ([509ce3f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/509ce3f168d977de71758518e99ce0e38ab9f875))
* **main:** release 0.56.0 ([#1493](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1493)) ([9a5fc2c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9a5fc2c0fdf993285bae42efb83b3384085540a0))
* **main:** release 0.56.1 ([#1504](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1504)) ([00fc00c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/00fc00c46f22984e02ed10acdc8041cfc79b507d))
* **main:** release 0.56.2 ([#1505](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1505)) ([f950dac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f950dac5d13516075c416f6abc6d7667474a36a8))
* **main:** release 0.56.3 ([#1511](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1511)) ([9c69643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c69643a31d91d0f3d249f7aea3beeefc53880ae))
* **main:** release 0.56.4 ([#1519](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1519)) ([d0384b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d0384b6d3bfc1bc358f39e58f136c1acef452456))
* **main:** release 0.56.5 ([#1555](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1555)) ([41663ee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/41663ee5900206a03c62e046bfb9659092197bd5))
* **main:** release 0.57.0 ([#1570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1570)) ([44b96cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/44b96cf67813f45feb67da4367936748bc04391f))
* **main:** release 0.58.0 ([#1587](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1587)) ([6b20b8d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b20b8d848620a7e9796ae230f6f87300e3fc50c))
* **main:** release 0.58.1 ([#1616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1616)) ([4780ba0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4780ba08b1bdf15785be63ec8dd488a03ddfe378))
* **main:** release 0.58.2 ([#1621](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1621)) ([1c34ac1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c34ac157bc064d5d6fe5297231ce87eccbcc298))
* **main:** release 0.59.0 ([#1622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1622)) ([afb18aa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afb18aa8ed3c3f80630bc2f824bb756ddb5eda86))
* **main:** release 0.60.0 ([#1641](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1641)) ([ab4d49f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab4d49f259db99c2c0c6131143c55ca11d2a6610))
* **main:** release 0.60.1 ([#1649](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1649)) ([56a9b2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56a9b2e5747bffb2456ad2a556e226e8450c242e))
* **main:** release 0.61.0 ([#1655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1655)) ([2fbe15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2fbe15a65a64adb8604d301e9a6d11632b6e3a44))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))

## [0.61.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.1...v0.61.0) (2023-03-27)


### Features

* add GRANT ... ON ALL TABLES IN ... ([#1626](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1626)) ([505a5f3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/505a5f35d9ea8388ca33e5117c545408243298ae))


### BugFixes

* all-grants ([#1658](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1658)) ([d5d59b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d5d59b4e85cd2e97ea0dc42b5ab2955ef35bbb88))

## [0.60.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.60.0...v0.60.1) (2023-03-23)


### BugFixes

* openbsd build ([#1647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1647)) ([6895a89](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6895a8958775e8e84a1457722f6c282d49458f3d))

## [0.60.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.59.0...v0.60.0) (2023-03-23)


### Features

* add missing PrivateLink URLs to datasource ([#1603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1603)) ([78782b1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/78782b1b471b7fbd434de1803cd687f6866cada7))
* add PREVENT_LOAD_FROM_INLINE_URL ([#1612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1612)) ([4945a3a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4945a3ae62d887dae6332742edcde715751459b5))
* Add support for `packages`, `imports`, `handler` and `runtimeVersion` to `snowflake_procedure` resource ([#1516](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1516)) ([a88f3ad](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a88f3ada75dad18b7b4b9024f664de8d687f54e0))


### BugFixes

* 0.60 misc bug fixes / linting ([#1643](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1643)) ([53da853](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53da853c213eec3afbdd2a47a8de3fba897c5d8a))
* change resource monitor suspend properties to number ([#1545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1545)) ([4bc59e2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4bc59e24677260dae94952bdbc5176ad177876dd))


### Misc

* **deps:** bump actions/setup-go from 3 to 4 ([#1634](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1634)) ([3f128c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f128c1ba887c377b7bd5f3d508d7b81676fdf90))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1639)) ([330777e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/330777eb0ad93acede6ffef9d7571c0989540657))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1638)) ([107bb4a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/107bb4abfb5de896acc1f224afae77b8100ffc02))

## [0.59.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.58.2...v0.59.0) (2023-03-21)


### Features

* add ON STAGE support for Stream resource ([#1413](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1413)) ([447febf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/447febfef46ef89570108d3447998d6d379b7be7))
* Add support for is_secure to snowflake_function resource ([#1575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1575)) ([c41b6a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c41b6a35271f12c97f5a4da947eb433013f2aaf2))


### BugFixes

* 0.59 release fixes ([#1636](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1636)) ([0a0256e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a0256ed3f0d56a6c7c22f810419636685094135))
* integration errors ([#1623](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1623)) ([83a40d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83a40d6361be0685b3864a0f3994298f3991de21))
* oauth integration ([#1315](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1315)) ([9087220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9087220af85f08880f7ad453cbe9d13dd3bc11ec))
* readded imported privileges special case for database grants ([#1597](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1597)) ([711ac0c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/711ac0cbc886bf8be6a5a2651234df778452b9df))

## [0.58.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.58.1...v0.58.2) (2023-03-16)


### BugFixes

* update 0.58.2 ([#1620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1620)) ([f1eab04](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1eab04dfdc839144057807953062b3591e6eaf0))

## [0.58.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.58.0...v0.58.1) (2023-03-16)


### BugFixes

* allow read of really old grant ids and add test cases ([#1615](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1615)) ([cda40ec](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cda40ece534cdc3f6849a7d24f2f8acea8476e69))
* update packages ([#1619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1619)) ([79a3acc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79a3acc0e3d6a405593b5adf90a31afef81d700f))

## [0.58.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.57.0...v0.58.0) (2023-03-03)


### Features

* **integration:** add google api integration ([#1589](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1589)) ([56909cd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/56909cdc18245f38b0f58bceaf2aa9cbc013d212))
* roles support numbers ([#1585](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1585)) ([d72dee8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d72dee82d0e0a4d8b484e5b204e156a13117cb76))


### BugFixes

* added force_new option to role grant when the role_name has been changed ([#1591](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1591)) ([4ec3613](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4ec3613de43d70f40a5d29ce5517af53e8ef0a06))
* default_secondary_roles doc ([#1584](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1584)) ([23b64fa](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/23b64fa9abcafb59610a77cafbda11a7e2ad648c))
* remove emojis, misc grant id fix ([#1598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1598)) ([fdefbc3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fdefbc3f1cc5bc7063f1cb1cc922854e8f9914e6))
* update read role grants to use new builder ([#1596](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1596)) ([e91860a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e91860ae794b034158b71ffb31097e73d8015c51))

## [0.57.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.5...v0.57.0) (2023-02-28)


### Features

* support object parameters on account level ([#1583](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1583)) ([fb24802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb2480214c8ac4e61fffe3a8e3448597462ad9a1))


### BugFixes

* Add contributing section to readme ([#1560](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1560)) ([174355d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/174355d740e325ae05435bbbc22b8b335f94fc6f))
* add test cases for update repl schedule on failover group ([#1578](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1578)) ([ab638f0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab638f0b9ba866d22c6f807743eb4eccad2530b8))
* schema read now checks first if the corresponding database exists ([#1568](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1568)) ([368dc8f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/368dc8fb3f7e5156d16caed1e03792654d49f3d4))

## [0.56.5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.4...v0.56.5) (2023-02-21)


### Misc

* **deps:** bump golang.org/x/crypto from 0.5.0 to 0.6.0 ([#1528](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1528)) ([8a011e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a011e0b1920833c77eb7832f821a4bd52176657))
* **deps:** bump golang.org/x/net from 0.5.0 to 0.7.0 ([#1551](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1551)) ([35de62f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/35de62f5b722c1dc6eaf2f39f6699935f67557cd))

## [0.56.4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.3...v0.56.4) (2023-02-17)


### BugFixes

* add nill check for grant_helpers ([#1518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1518)) ([87689bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/87689bb5b60c73bfe3d741c3da6f4f544f16aa45))
* cleanup grant logic ([#1522](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1522)) ([0502c61](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0502c61e7211253d029a0bec6a8104738624f243))
* future read on grants ([#1520](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1520)) ([db78f64](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/db78f64e56d228f3236b6bdefbe9a9c18c8641e1))

## [0.56.3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.2...v0.56.3) (2023-02-02)


### BugFixes

* multiple share grants ([#1510](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1510)) ([d501226](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d501226bc2ee8274446efb282c2dfea9599a3c2e))

## [0.56.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.1...v0.56.2) (2023-02-01)


### BugFixes

* backwards compatability for grant helpers id used by procedure and functions ([#1508](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1508)) ([3787657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3787657105fbcf18368136813afd558251f92cd1))
* misc linting changes for 0.56.2 ([#1509](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1509)) ([e0d1ef5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d1ef5c718f9e1e58e80d31bbe2d2f27afec486))
* support different tag association queries for COLUMN object types ([#1380](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1380)) ([546d0a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/546d0a144e77c759cd6ddb91a193253f27f8fb91))

## [0.56.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.56.0...v0.56.1) (2023-01-31)


### BugFixes

* procedure and function grants ([#1502](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1502)) ([0d08ea8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0d08ea85541ceff6e591d34a671b44ef778a6611))

## [0.56.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.55.1...v0.56.0) (2023-01-27)


### Features

* add new account resource ([#1492](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1492)) ([b1473ba](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b1473ba158946d81bc4eac095c40c8d0446cf2ed))


### BugFixes

* do not set query_acceleration_max_scale_factor when enable enable_query_acceleration = false ([#1474](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1474)) ([d62b1b4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d62b1b4d6352e7d2dc99e4603370a1f3de3a4624))
* refactor for simplify handling error ([#1472](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1472)) ([3937216](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/393721607c9eee5d73e14c27265eb39f195ccb37))
* refactor handling error to be simple ([#1473](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1473)) ([9f37b99](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f37b997de073f01b66c86820237eff8049346ba))
* saml integration test ([#1494](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1494)) ([8c31439](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8c31439253d25aafb54fc09d89e547fa8238258c))
* schema name is optional for future file_format_grant ([#1484](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1484)) ([1450cdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1450cddde6328264f9df37e4dd89a78f5f095b2e))
* schema name is optional for future function_grant ([#1485](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1485)) ([dcc550e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/dcc550ed5b3df548d5d300cd2b77907ea544bb43))
* schema name is optional for future procedure_grant ([#1486](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1486)) ([4cf4561](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cf456151d83cd71a3b9e68abe9c6f29804f2ee2))
* schema name is optional for future sequence_grant ([#1487](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1487)) ([ccf9e78](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ccf9e78c9a7884e5beea233dd529a5134c741fb1))
* schema name is optional for future snowflake_stage_grant ([#1466](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1466)) ([0b4d814](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b4d8146910e8ea31d2ed5ea8b58725449205dcd))
* schema name is optional for future stream_grant ([#1488](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1488)) ([3f7e5d6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3f7e5d655ed5738107536c873dd11533573bba46))
* schema name is optional for future task_grant ([#1489](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1489)) ([4096fd0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4096fd0d8bc65ae23b6d588385e1f81c4f2e7521))

## [0.55.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.55.0...v0.55.1) (2023-01-11)


### BugFixes

* resource snowflake_resource_monitor behavior conflict from provider 0.54.0 to 0.55.0 ([#1468](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1468)) ([8ce0c53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ce0c533ec5d81273df20be2126b278ca61a59f6))

## [0.55.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.54.0...v0.55.0) (2023-01-10)


### Features

* add in more functionality for UpdateResourceMonitor  ([#1456](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1456)) ([2df570f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2df570f0c3271534a37b0cb61b7f4b491081baf7))


### Misc

* **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1454)) ([ed6bfe0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ed6bfe07930e5703036ada816845176d46f5623c))
* **deps:** bump golang.org/x/tools from 0.4.0 to 0.5.0 ([#1455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1455)) ([ff01970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ff019702fdc1ef810bb94533489b89a956f09ef4))


### BugFixes

* 0.55 fix ([#1465](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1465)) ([8cb3370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8cb337048ec5c4a52245feb1b9556fd845d83278))
* add permissions ([#1464](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1464)) ([e2d249a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e2d249ae1466e05dad2080f05123e0de66fabcf6))
* interval for failover groups ([#1448](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1448)) ([bd1d3cc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bd1d3cc57f72c7774715f1d92a955536d55fb758))
* saml2_sign_request and saml2_force_authn cast type ([#1452](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1452)) ([f8cecd7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f8cecd7ca45aabec78fd18d8aa220db7eb34b9e0))
* schema_name is optional to enable future pipe grant ([#1424](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1424)) ([5d966fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5d966fd8624fa3208cebae3d7b32c1b59bdcfd4c))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.54.0...v0.34.0) (2022-12-23)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Add title lint ([#570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/570)) ([d2142fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2142fd408f158a68230f0188c35c7b322c70ab7))
* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding slack bot for PRs and Issues ([#1106](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1106)) ([95c255e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/95c255e5ca65b619b35692671848877793cac29e))
* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* Resource to manage a user's public keys ([#540](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/540)) ([590c22e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/590c22ec40ed28c7d280192ed66c4d93623e32fd))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))
* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))
* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))
* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))
* **main:** release 0.26.0 ([#841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/841)) ([4a6d659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a6d659c7069c1d2d64c43ce3262d3a7a840b7c5))
* **main:** release 0.26.1 ([#849](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/849)) ([a2607e5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a2607e5e15d6dfd66e756e381c0aeccf8106bbd4))
* **main:** release 0.26.2 ([#851](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/851)) ([016e02d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/016e02d3cc51360ecae43df6a931ada2b398e424))
* **main:** release 0.26.3 ([#854](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/854)) ([63f2b85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f2b8507930b9429ebf7c8ee8f65334ef16931e))
* **main:** release 0.27.0 ([#870](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/870)) ([5178aa6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5178aa6a977fe296bc4b5abeae6e55ca27291dca))
* **main:** release 0.28.0 ([#886](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/886)) ([b523f7e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b523f7e263f988a839528bb19b04405890013879))
* **main:** release 0.28.1 ([#895](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/895)) ([c92c518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c92c5184eab43141116d760f9e336714eb535fd7))
* **main:** release 0.28.2 ([#902](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/902)) ([e1c228b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1c228b9bbba60d297082b665159ca9160607e62))
* **main:** release 0.28.3 ([#905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/905)) ([b01a21b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b01a21bc7fa2055bbaf77af8e263e69fbb1bfd54))
* **main:** release 0.28.4 ([#913](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/913)) ([4fa19e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4fa19e37edda8c3958232c647d35bf99a7d00319))
* **main:** release 0.28.5 ([#915](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/915)) ([d9a7f01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d9a7f0165cc440cc7ed6086d033ab7252e56c5e2))
* **main:** release 0.28.6 ([#920](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/920)) ([3a17e34](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3a17e34a9e762ee24d8ff12144fe6c6ac4b68e2e))
* **main:** release 0.28.7 ([#921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/921)) ([adbb52e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/adbb52e3f33c86519ed20f490bddd84980437e3f))
* **main:** release 0.28.8 ([#928](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/928)) ([96152d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/96152d7db14c397bff9f9cc79ba0d85f6f2706b4))
* **main:** release 0.29.0 ([#943](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/943)) ([f1d0af9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1d0af96bbe8e57558bd3a57808298d8d49fe349))
* **main:** release 0.30.0 ([#954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/954)) ([bfd3108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bfd31080b96af02f908e93ff0c20b8cb24840431))
* **main:** release 0.31.0 ([#968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/968)) ([1e21100](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e2110065c23d851e04cd2de1727b683a38168f1))
* **main:** release 0.32.0 ([#974](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/974)) ([e947417](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e947417c459e424829a9b9e4cbb96f04ba7db3cd))
* **main:** release 0.33.0 ([#988](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/988)) ([bf3482e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bf3482e4de81e96b31aec192a15f5bee33d34e78))
* **main:** release 0.33.1 ([#991](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/991)) ([1c5af87](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c5af874164d8b40e7cae54e9206ec6b46c2e75b))
* **main:** release 0.34.0 ([#1014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1014)) ([f1c651e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1c651e17d1697f37be43857318573cb56812f5d))
* **main:** release 0.34.0 ([#1019](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1019)) ([83db3a4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83db3a4c14ec6f1539fbef55c72ae36b22e76906))
* **main:** release 0.34.0 ([#1020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1020)) ([7116025](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7116025e3523cc6d385752f7e71bff1b5fded68b))
* **main:** release 0.34.0 ([#1022](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1022)) ([d06c91f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d06c91fdacbc223cac709743a0fbe9d2c340da83))
* **main:** release 0.34.0 ([#1332](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1332)) ([7037952](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7037952180309441ac865eed0bc2a44a714b484d))
* **main:** release 0.35.0 ([#1026](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1026)) ([f9036e8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f9036e8914b9c139eb6798276124c5544a083eb8))
* **main:** release 0.36.0 ([#1056](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1056)) ([d055d4c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d055d4c57f9a48855382506a313a4f6386da2e3e))
* **main:** release 0.37.0 ([#1065](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1065)) ([6aecc46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6aecc46ddc0804a3a8b90422dfeb4c3bfbf093e5))
* **main:** release 0.37.1 ([#1096](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1096)) ([1de53b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1de53b5ee9247216b547398c29c22956247c0563))
* **main:** release 0.38.0 ([#1103](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1103)) ([aee8431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/aee8431ea64f085de0f4e9cfd46f2b82d16f09e2))
* **main:** release 0.39.0 ([#1130](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1130)) ([82616e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82616e325890613d4b2eca5ef6ffa2e3b50a0352))
* **main:** release 0.40.0 ([#1132](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1132)) ([f3f1f3b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f3f1f3b517963c544da1a64d8d778c118a502b29))
* **main:** release 0.41.0 ([#1157](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1157)) ([5b9b47d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5b9b47d6fa2da7cd6d4b0bfe1722794003a5fce5))
* **main:** release 0.42.0 ([#1179](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1179)) ([ba45fc2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ba45fc27b7e3d2eda70966a857ebcd37964a5813))
* **main:** release 0.42.1 ([#1191](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1191)) ([7f9a3c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f9a3c2dd172fa93d1d2648f13b77b1f8f7981f0))
* **main:** release 0.43.0 ([#1196](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1196)) ([3ac84ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3ac84ab0834d3ab875d078489a2d2b7a45cfad28))
* **main:** release 0.43.1 ([#1207](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1207)) ([e61c15a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e61c15aab3991e9740da365ec739f0c03fbbbf65))
* **main:** release 0.44.0 ([#1222](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1222)) ([1852308](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/185230847b7179079c718078780d240a9c29bbb0))
* **main:** release 0.45.0 ([#1232](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1232)) ([da886d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/da886d4e05f7bb9443168f0fa04b8b397a1db5c7))
* **main:** release 0.46.0 ([#1244](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1244)) ([b9bf009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9bf009a11a7af0413c8f182927731f55379dff4))
* **main:** release 0.47.0 ([#1259](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1259)) ([887297f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/887297fc5670b180f3d7158d3092ad035fb473e9))
* **main:** release 0.48.0 ([#1284](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1284)) ([cf6e54f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cf6e54f720dd852c1663a4e9ff8a74054f51325b))
* **main:** release 0.49.0 ([#1303](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1303)) ([fb90556](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb90556c324ffc14b6e90adbdf9a06705af8e7e9))
* **main:** release 0.49.1 ([#1319](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1319)) ([431b8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/431b8b7818cd7eccb3dafb11612f72ce8e73b58f))
* **main:** release 0.49.2 ([#1323](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1323)) ([c19f307](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c19f3070b8aa063c96e1e30a1e6d754b7070d296))
* **main:** release 0.49.3 ([#1327](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1327)) ([102ed1d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/102ed1de7f4fca659869fc0485b42843b394d7e9))
* **main:** release 0.50.0 ([#1344](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1344)) ([a860a76](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a860a7623b9e22433ece8cede537c187a45b4bc2))
* **main:** release 0.51.0 ([#1348](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1348)) ([2b273f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2b273f7e3baaf855ed6e02a7779022f38ade6745))
* **main:** release 0.52.0 ([#1363](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1363)) ([e122715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1227159be50bf26841acead8730dad516a96ebc))
* **main:** release 0.53.0 ([#1401](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1401)) ([80488da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/80488dae4e16f5c55f913449fc729fbd6e1fd6d2))
* **main:** release 0.53.1 ([#1406](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1406)) ([8f5ac41](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8f5ac41265bc08256630b2d95fa8845249098310))
* **main:** release 0.54.0 ([#1431](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1431)) ([6b6b55d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b6b55d88a875f30395f2bd3250a2af1b99f9205))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))


### BugFixes

* 0.54  ([#1435](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1435)) ([4c9dd13](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4c9dd133574b08d8e67444b6c6b81aa87d9a2acf))
* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* Allow creation of database-wide future external table grants ([#1041](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1041)) ([5dff645](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5dff645291885cd437e341148c0629fe7ab7383f))
* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))
* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))
* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* OSCP -&gt; OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))
* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))
* read Database and Schema name during Stream import ([#732](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/732)) ([9f747b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f747b571b2fcf5b0663696efd75c55a6f8b6a89))
* read Name, Database and Schema during Procedure import ([#819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/819)) ([d17656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d17656fdd2803516b6ee067a6bd5a457bf31d905))
* Recreate notification integration when type changes ([#792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/792)) ([e9768bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9768bf059268fb933ad74f2b459e91e2c0563e0))
* refactor ReadWarehouse function to correctly read object parameters ([#745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/745)) ([d83c499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83c499910c0f2b6348191a93f917e450b9e69b2))
* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))
* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))
* Remove force_new from masking_expression ([#588](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/588)) ([fc3e78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3e78acbdda5346f32a004711d31ad6f68529ed))
* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* Removing force new and adding update for data base replication config ([#1105](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1105)) ([f34f012](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f34f012195d0b9718904ffa7a3a529f58167a74e))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))
* Updating shares to disallow account locators ([#1102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1102)) ([4079080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4079080dd0b9e3caf4b5d360000bd216906cb81e))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))

## [0.54.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.53.1...v0.54.0) (2022-12-23)


### Features

* add parameters resources + ds ([#1429](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1429)) ([be81aea](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be81aea070d47acf11e2daed4a0c33cd120ab21c))
* Current role data source ([#1415](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1415)) ([8152aee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8152aee136e279832b59a6ec1b165390e27a1e0e))


### Misc

* **deps:** bump github.com/snowflakedb/gosnowflake ([#1423](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1423)) ([84c9389](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/84c9389c7e945c0b616cacf23b8252c35ff307b3))
* **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1426](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1426)) ([409bcb1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/409bcb19ce17a1babd685ddebbea32f2552d29bd))


### BugFixes

* Don't throw an error on unhandled Role Grants ([#1414](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1414)) ([be7e78b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be7e78b31cc460e562de47613a0a095ec623a0ae))
* go syntax ([#1410](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1410)) ([c5f6b9f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c5f6b9f6a4ccd7c96ad5fb67a10161cdd71da833))
* Go syntax to add revive ([#1411](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1411)) ([b484bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b484bc8a70ab90eb3882d1d49e3020464dd654ec))
* linting errors ([#1432](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1432)) ([665c944](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/665c94480be82831ec33650175d905c048174f7c))

## [0.53.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.53.0...v0.53.1) (2022-12-08)


### Misc

* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1373](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1373)) ([b22a2bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b22a2bdc5c2ec3031fb116323f9802945efddcc2))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1375](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1375)) ([e1891b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1891b61ef5eeabc49276099594d9c1726ca5373))
* **deps:** bump github/codeql-action from 1 to 2 ([#1353](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1353)) ([9d7bc15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d7bc15790eca62d893a2bec3535d468e34710c2))
* **deps:** bump golang.org/x/crypto from 0.1.0 to 0.4.0 ([#1407](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1407)) ([fc96d62](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc96d62119bdd985eca8b7c6b09031592a4a7f65))
* **deps:** bump golang.org/x/tools from 0.2.0 to 0.4.0 ([#1400](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1400)) ([58ca9d8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/58ca9d895254574bc54fadf0ca202a0ab99992fb))
* **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1354](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1354)) ([9ad93a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9ad93a85a72e54d4b93339a3078ab1d4ca85a764))
* **deps:** bump peter-evans/create-or-update-comment from 1 to 2 ([#1350](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1350)) ([d4d340e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d4d340e85369fa1727014d3f51f752b85687994c))
* **deps:** bump peter-evans/find-comment from 1 to 2 ([#1352](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1352)) ([ce13a8e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ce13a8e6655f9cbe03bb2e1c91b9f5746fd5d5f7))
* **deps:** bump peter-evans/slash-command-dispatch from 2 to 3 ([#1351](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1351)) ([9d17ead](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9d17ead0156979a5001f95bbc5636221b232fb17))


### BugFixes

* docs ([#1409](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1409)) ([fb68c25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb68c25d9c1145fa9bbe38395ce1594d9d127139))

## [0.53.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.52.0...v0.53.0) (2022-12-07)


### Features

* Added (missing) API Key to API Integration ([#1386](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1386)) ([500d6cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/500d6cf21e983515a95b142d2745594684df33a0))
* Adding warehouse type for snowpark optimized warehouses ([#1369](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1369)) ([b5bedf9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b5bedf90720fcc64cf3e01add659b077b34e5ae7))
* **ci:** add depguard ([#1368](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1368)) ([1b29f05](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1b29f05d67a1d2fb7938f2c1c0b27071d47f13ab))
* **ci:** add goimports and makezero ([#1378](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1378)) ([b0e6580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0e6580d1086cc9cc2000b201425aa049e684502))


### BugFixes

* errors package with new linter rules ([#1360](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1360)) ([b8df2d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b8df2d737239d7c7b472fb3e031cccdeef832c2d))
* File Format Update Grants ([#1397](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1397)) ([19933c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/19933c04d7e9c10a08b5a06fe70a2f31fdd6c52e))
* Go Expression Fix [#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384) ([#1403](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1403)) ([8936e1a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8936e1a0defc2b6d11812a88f486903a3ced31ac))
* insecure go expression ([#1384](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1384)) ([a6c8e75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6c8e75e142f28ad6e2e9ef3ff4b2b877c101c90))
* preallocate slice ([#1385](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1385)) ([9e972c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9e972c06f7840d1b516766068bb92f7cb458c428))
* remove shares from snowflake_stage_grant [#1285](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1285) ([#1361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1361)) ([3167d9d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3167d9d402960cb2535a036aa373ad9e62d3ef18))
* Removed Read for API_KEY ([#1402](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1402)) ([ddd00c5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ddd00c5b7e1862e2328dbdf599d157a443dce134))
* stop file format failure when does not exist ([#1399](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1399)) ([3611ff5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3611ff5afe3e44c63cdec6ff8b191d0d88849426))
* warehouses update issue ([#1405](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1405)) ([1c57462](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c57462a78f6836ed67678a88b6529a4d75f6b9e))
* wrong usage of testify Equal() function ([#1379](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1379)) ([476b330](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/476b330e69735a285322506d0656b7ea96e359bd))

## [0.52.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.51.0...v0.52.0) (2022-11-17)


### Features

* **ci:** golangci lint adding thelper, wastedassign and whitespace ([#1356](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1356)) ([0079bee](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0079bee139f9cbaaa4b26c2a92a56c37a9366d68))
* grants datasource ([#1377](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1377)) ([0daafa0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0daafa09cb0c53e9a51e42a9574533ebd81135b4))


### Misc

* **docs:** terraform fmt ([#1358](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1358)) ([0a2fe08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0a2fe089fd777fc44583ee3616a726840a13d984))


### BugFixes

* **ci:** remove unnecessary type conversions ([#1357](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1357)) ([1d2b455](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d2b4550902767baad67f88df42d773b76b952b8))
* Improve table constraint docs ([#1355](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1355)) ([7c650bd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7c650bd601662ed71aa06f5f71eddbf9dedb95bd))
* test modules in acceptance test for warehouse ([#1359](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1359)) ([2d8f2b6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2d8f2b6ec0564bbbf577f8efaf9b2d8103198b22))
* update id serialization ([#1362](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1362)) ([4d08a8c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4d08a8cd4058df12d536739965efed776ec7f364))

## [0.51.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.50.0...v0.51.0) (2022-11-07)


### Features

* add support for `notify_users` to `snowflake_resource_monitor` resource ([#1340](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1340)) ([7094f15](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7094f15133cd768bd4aa4431adc66802a7f955c0))
* **ci:** add some linters and fix codes to pass lint ([#1345](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1345)) ([75557d4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75557d49bd03b21fa3cca903c1207b01cf6fcead))
* Delete ownership grant updates ([#1334](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1334)) ([4e6aba7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4e6aba780edf81624b0b12c171d24802c9a2911b))
* show roles data source ([#1309](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1309)) ([b2e5ecf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b2e5ecf050711a9562857bd5e0eee383a6ed497c))


### Misc

* **docs:** update documentation adding double quotes ([#1346](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1346)) ([c4af174](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4af1741347dc080211c726dd1c80116b5e121ef))


### BugFixes

* format for go ci ([#1349](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1349)) ([75d7fd5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/75d7fd54c2758783f448626165062bc8f1c8ebf1))
* tag on schema fix ([#1313](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1313)) ([62bf8b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/62bf8b77e841cf58b622e77d7f2b3cb53d7361c5))
* workflow warnings ([#1316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1316)) ([6f513c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f513c27810ed62d49f0e10895cefc219e9d9226))

## [0.50.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.49.0...v0.50.0) (2022-11-04)


### Features

* task after dag support ([#1342](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1342)) ([a117802](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a117802016c7e47ef539522c7308966c9f1c613a))

## [0.49.3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.49.2...v0.49.3) (2022-11-01)


### BugFixes

* doc ([#1326](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1326)) ([d7d5e08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d7d5e08159b2e199e344048c4ab40f3d756e670a))

## [0.49.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.49.1...v0.49.2) (2022-11-01)


### BugFixes

* validate identifier ([#1312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1312)) ([295bc0f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/295bc0fd852ff417c740d19fab4c7705537321d5))

## [0.49.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.49.0...v0.49.1) (2022-10-31)


### BugFixes

* ie-proxy for go build ([#1318](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1318)) ([c55c101](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c55c10178520a9d668ee7b64145a4855a40d9db5))

## [0.49.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.48.0...v0.49.0) (2022-10-31)


### Features

* add column masking policy specification ([#796](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/796)) ([c1e763c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1e763c953ba52292a0473341cdc0c03b6ff83ed))
* add failover groups ([#1302](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1302)) ([687742c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/687742cc3bd81f1d94de3c28f272becf893e365e))
* Added Query Acceleration for Warehouses ([#1239](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1239)) ([ad4ce91](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ad4ce919b81a8f4e93835244be0a98cb3e20204b))
* task with allow_overlapping_execution option ([#1291](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1291)) ([8393763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/839376316478ab7903e9e4352e3f17665b84cf60))


### Misc

* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1280](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1280)) ([657a180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/657a1800f9394c5d03cc356cf92ed13d36e9f25b))
* **deps:** bump github.com/snowflakedb/gosnowflake ([#1304](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1304)) ([fb61921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb61921f0f28b0745279063402feb5ff95d8cca4))
* **deps:** bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#1300](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1300)) ([2f3c612](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f3c61237d21bc3affadf1f0e08234f5c404dde6))
* **deps:** bump golang.org/x/tools from 0.1.12 to 0.2.0 ([#1295](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1295)) ([5de7a51](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5de7a5188089e7bf55b6af679ebff43f98474f78))
* update docs ([#1297](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1297)) ([495558c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/495558c57ed2158fd5f1ea26edd111de902fd607))


### BugFixes

* golangci.yml to keep quality of codes ([#1296](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1296)) ([792665f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/792665f7fea6cbe3c5df4906ba298efd2f6727a1))
* run check docs ([#1306](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1306)) ([53698c9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/53698c9e7d020f1711e42d024139132ecee1c09f))
* update doc ([#1305](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1305)) ([4a82c67](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a82c67baf7ef95129e76042ff46d8870081f6d1))
* weird formatting ([526b852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/526b852cf3b2d40a71f0f8fad359b21970c2946e))

## [0.48.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.47.0...v0.48.0) (2022-10-24)


### Features

* add custom oauth int ([#1286](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1286)) ([d6397f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d6397f9d331e2e4f658e62f17892630c7993606f))


### BugFixes

* clean up tag association read ([#1261](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1261)) ([de5dc85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/de5dc852dff2d3b9cfb2cf6d20dea2867f1e605a))
* Fixed Grants Resource Update With Futures ([#1289](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1289)) ([132373c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/132373cbe944899e0b5b0043bfdcb85e8913704b))
* Pin Jira actions versions ([#1283](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1283)) ([ca25f25](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ca25f256e52cd70248d0fcb33e60a7741041a268))
* tag association name convention ([#1294](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1294)) ([472f712](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/472f712f1db1c4fabd70b4f98188b157d8fb00f5))

## [0.47.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.46.0...v0.47.0) (2022-10-11)


### Features

* add new table constraint resource ([#1252](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1252)) ([fb1f145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fb1f145900dc27479e3769042b5b303d1dcef047))
* integer return type for procedure ([#1266](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1266)) ([c1cf881](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c1cf881c0faa8634a375de80a8aa921fdfe090bf))


### Misc

* add godot to golangci ([#1263](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1263)) ([3323470](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3323470a7be1988d0d3d11deef3191078872c06c))


### BugFixes

* Database tags UNSET ([#1256](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1256)) ([#1257](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1257)) ([3d5dcac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3d5dcac99c7fa859a811c72ce3dcd1f217c4f7d7))
* missing t.Helper for thelper function ([#1264](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1264)) ([17bd501](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17bd5014282201023572348a5ab51a3bf849ce86))
* misspelling ([#1262](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1262)) ([e9595f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9595f27d0f181a32e77116c950cf141708221f5))

## [0.46.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.45.0...v0.46.0) (2022-09-29)


### Features

* Added Missing Grant Updates + Removed ForceNew ([#1228](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1228)) ([1e9332d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e9332d522beed99d80ecc2d0fc40fedc41cbd12))


### BugFixes

* Table Tags Acceptance Test ([#1245](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1245)) ([ab34763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ab347635d2b1a1cb349a3762c0869ef71ab0bacf))
* Update 'user_ownership_grant' schema validation ([#1242](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1242)) ([061a28a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/061a28a9a88717c0b37b18a564f55f88cbed56ea))

## [0.45.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.44.0...v0.45.0) (2022-09-22)


### Features

* add connection param for snowhouse ([#1231](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1231)) ([050c0a2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/050c0a213033f6f83b5937c0f34a027347bfbb2a))
* add port and protocol to provider config ([#1238](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1238)) ([7a6d312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a6d312e0becbb562707face1b0d87b705692687))

## [0.44.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.43.1...v0.44.0) (2022-09-20)


### Features

* Create a snowflake_user_grant resource. ([#1193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1193)) ([37500ac](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/37500ac88a3980ea180d7b0992bedfbc4b8a4a1e))


### BugFixes

* function not exist and integration grant ([#1154](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1154)) ([ea01e66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea01e66797703e53c58e29d3bdb36557b22dbf79))

## [0.43.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.43.0...v0.43.1) (2022-09-20)


### BugFixes

* add sweepers ([#1203](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1203)) ([6c004a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6c004a31d7d5192f4136126db3b936a4be26ff2c))
* Pass file_format values as-is in external table configuration ([#1183](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1183)) ([d3ad8a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d3ad8a8019ffff65e644e347e21b8b1512be65c4)), closes [#1046](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1046)
* remove stage from statefile if not found ([#1220](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1220)) ([b570217](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b57021705f5b554499b00289e7219ee6dabb70a1))

## [0.43.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.42.1...v0.43.0) (2022-08-31)


### Features

* tag based masking policy ([#1143](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1143)) ([e388545](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e388545cae20da8c011e644ac7ecaf2724f1e374))


### BugFixes

* log fmt ([#1192](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1192)) ([0f2e2db](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0f2e2db2343237620aceb416eb8603b8e42e11ec))

## [0.42.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.42.0...v0.42.1) (2022-08-24)


### Misc

* update-license ([#1190](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1190)) ([e9cfc3e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9cfc3e7d07ee5d60f55d842c13f2d8fc20e7ba6))

## [0.42.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.41.0...v0.42.0) (2022-08-24)


### Features

* tag association resource ([#1187](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1187)) ([123fd2f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/123fd2f88a18242dbb3b1f20920c869fd3f26651))
* transient database ([#1165](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1165)) ([f65a0b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f65a0b501ee7823575c73071115f96973834b07c))


### BugFixes

* Database from share comment on create and docs ([#1167](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1167)) ([fc3a8c2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3a8c289fa8466e0ad8fa9454e31c27d75de563))
* doc of resource_monitor_grant ([#1188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1188)) ([03a6cb3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a6cb3c58f6ce5860b70f62a08befa7c9905df8))
* Fix snowflake_share resource not unsetting accounts ([#1186](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1186)) ([03a225f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/03a225f94a8e641dc2a08fdd3247cc5bd64708e1))

## [0.41.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.40.0...v0.41.0) (2022-08-10)


### Features

* Adding support for debugger-based debugging. ([#1145](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1145)) ([5509899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5509899df90be7e01826261d2f626239f121437c))
* tag grants ([#1127](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1127)) ([018e7ab](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/018e7ababa73a579c79f3939b83a9010fe0b2774))


### BugFixes

* adding in issue link to slackbot ([#1158](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1158)) ([6f8510b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f8510b8e8b7c6b415ef6258a7c1a2f9e1b547c4))
* Deleting a snowflake_user and their associated snowlfake_role_grant causes an error ([#1142](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1142)) ([5f6725a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5f6725a8d0df2f5924c6d6dc2f62ebeff77c8e14))
* doc pipe ([#1171](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1171)) ([c94c2f9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c94c2f913bc47c69edfda2f6e0ef4ff34f52da63))
* expand allowed special characters in role names ([#1162](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1162)) ([30a59e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30a59e0657183aee670018decf89e1c2ef876310))
* Remove validate_utf8 parameter from file_format ([#1166](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1166)) ([6595eeb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6595eeb52ef817981bfa44602a211c5c8b8de29a))

## [0.40.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.39.0...v0.40.0) (2022-07-14)


### Features

* add AWS GOV support in api_integration ([#1118](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1118)) ([2705970](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/270597086e3c9ec2af5b5c2161a09a5a2e3f7511))
* S3GOV support to storage_integration ([#1133](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1133)) ([92a5e35](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/92a5e35726be737df49f2c416359d1c591418ea2))
* Streams on views ([#1112](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1112)) ([7a27b40](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7a27b40cff5cc75fe9743e1ba775254888291662))


### BugFixes

* update team slack bot configurations ([#1134](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1134)) ([b83a461](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b83a461771c150b53f566ad4563a32bea9d3d6d7))

## [0.39.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.38.0...v0.39.0) (2022-07-13)


### Features

* deleting gpg agent before importing key ([#1123](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1123)) ([e895642](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e895642db51988807aa7cb3fc3d787aee37963f1))


### BugFixes

* changing tool to ghaction-import for importing gpg keys ([#1129](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1129)) ([5fadf08](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5fadf08de5cba1a34988b10e12eec392842777b5))
* Delete gpg change ([#1126](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1126)) ([ea27084](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ea27084cda350684025a2a58055ea4bec7427ef5))

## [0.37.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.37.0...v0.37.1) (2022-07-01)


### BugFixes

* Allow creation of stage with storage integration including special characters ([#1081](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1081)) ([7b5bf00](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b5bf00183595a5412f0a5f19c0c3df79502a711)), closes [#1080](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1080)
* warehouse import when auto_suspend is set to null ([#1092](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1092)) ([9dc748f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9dc748f2b7ff98909bf285685a21175940b8e0d8))

## [0.37.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.36.0...v0.37.0) (2022-06-28)


### Features

* add python language support for functions ([#1063](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1063)) ([ee4c2c1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ee4c2c1b3b2fecf7319a5d58d17ae87ff4bcf685))
* Python support for functions ([#1069](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1069)) ([bab729a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bab729a802a2ae568943a89ebad53727afb86e13))


### BugFixes

* Handling 2022_03 breaking changes ([#1072](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1072)) ([88f4d44](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/88f4d44a7f33abc234b3f67aa372230095c841bb))
* Temporarily disabling acceptance tests for release ([#1083](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1083)) ([8eeb4b7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8eeb4b7ff62ef442c45f0b8e3105cd5dc1ff7ccb))

## [0.36.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.35.0...v0.36.0) (2022-06-16)


### Features

* Add support for default_secondary_roles ([#1030](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1030)) ([ae8f3da](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/ae8f3dac67e8bf5c4cd73fb08101d378be32e39f))


### BugFixes

* allow custom characters to be ignored from validation ([#1059](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1059)) ([b65d692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b65d692c83202d3e23628d727d71abf1f603d32a))
* Correctly read INSERT_ONLY mode for streams ([#1047](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1047)) ([9c034fe](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9c034fef3f6ac1e51f6a6aae3460221d642a2bc8))
* handling not exist gracefully ([#1031](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1031)) ([101267d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/101267dd26a03cb8bc6147e06bd467fe895e3b5e))

## [0.35.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.34.0...v0.35.0) (2022-06-07)


### Features

* add allowed values ([#1006](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1006)) ([e7dcfd4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7dcfd4c1f9b77b4d03bfb9c13a8753000f281e2))
* Add allowed values ([#1028](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1028)) ([e756867](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e7568674807af2899a2d1579aec53c706598bccf))
* Add support for creation of streams on external tables ([#999](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/999)) ([0ee1d55](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ee1d556abcf6aaa14ff041155c57ff635c5cf94))
* Support for selecting language in snowflake_procedure ([#1010](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1010)) ([3161827](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/31618278866604e8bfd7d2fa984ec9502c0b7bbb))


### BugFixes

* makefile remove outdated version reference ([#1027](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1027)) ([d066d0b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d066d0b7b7b1604e157d70cc14e5babae2b3ef6b))
* provider upgrade doc ([#1039](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1039)) ([e1e23b9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1e23b94c536f40e1e2418d8af6aa727dfec0d52))


### Misc

* **deps:** bump github.com/hashicorp/terraform-plugin-sdk/v2 ([#1035](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1035)) ([f885f1c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f885f1c0325c019eb3bb6c0d27e58a0aedcd1b53))

## [0.34.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.34.0...v0.34.0) (2022-05-25)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))
* Add a resource to manage sequences ([#582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/582)) ([7fab82f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7fab82f6e9e7452b726ccffc7e935b6b47c53df4))
* Add CREATE ROW ACCESS POLICY to schema grant priv list ([#581](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/581)) ([b9d0e9e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b9d0e9e5b3076eaeec1e47b9d3c9ca14902e5b79))
* Add file format resource ([#577](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/577)) ([6b95dcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6b95dcb0236a7064dd99418de90fc0086f548a78))
* Add importer to integration grant ([#574](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/574)) ([3739d53](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3739d53a72cf0103e7dbfb42260cb7ab98b94f92))
* Add INSERT_ONLY option to streams ([#655](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/655)) ([c906e01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c906e01181d8ffce332e61cf82c57d3bf0b4e3b1))
* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* Add private key passphrase support ([#639](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/639)) ([a1c4067](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a1c406774728e25c51e4da23896b8f40a7090453))
* Add REBUILD table grant ([#638](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/638)) ([0b21c66](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0b21c6694a0e9f7cf6a1dbf28f07a7d0f9f875e9))
* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Add SHOW_INITIAL_ROWS to stream resource ([#575](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/575)) ([3963193](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39631932d6e90e4707a73cca9c5f1237cf3c3a1c))
* add STORAGE_AWS_OBJECT_ACL support to storage integration ([#755](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/755)) ([e136b1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e136b1e0fddebec6874d37bec43e45c9cdab134d))
* Add support for error notifications for Snowpipe ([#595](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/595)) ([90af4cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90af4cf9ed17d06d303a17126190d5b5ea953bc6))
* Add support for GCP notification integration ([#603](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/603)) ([8a08ee6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8a08ee621fea310b627f5be349019ff8638e491b))
* Add support for table column comments and to control a tables data retention and change tracking settings ([#614](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/614)) ([daa46a0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/daa46a072aa2d8d7fe8ac45250c8a93769687f81))
* add the param "pattern" for snowflake_external_table ([#657](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/657)) ([4b5aef6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4b5aef6afd4fed147604c1658b69f3a80bccebab))
* Add title lint ([#570](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/570)) ([d2142fd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d2142fd408f158a68230f0188c35c7b322c70ab7))
* Added Functions (UDF) Resource & Datasource ([#647](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/647)) ([f28c7dc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f28c7dc7cab3ac27df6201954c535c266c6564db))
* Added Procedures Datasource ([#646](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/646)) ([633f2bb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/633f2bb6db84576f07ad3800808dbfe1a84633c4))
* Added Row Access Policy Resources ([#624](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/624)) ([fd97816](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fd97816411189956b63fafbfcdeed12810c91080))
* Added Several Datasources Part 2 ([#622](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/622)) ([2a99ea9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a99ea97972e2bbf9e8a27c9e6ecefc990145f8b))
* Adding Database Replication ([#1007](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1007)) ([26aa08e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/26aa08e767be2ee4ed8a588b796845f76a75c02d))
* adding in tag support ([#713](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/713)) ([f75cd6e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75cd6e5f727b149f9c04f672c985a214a0ceb7c))
* Adding users datasource ([#1013](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1013)) ([4cd86e4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4cd86e4abd58292ebf6fdfa0c5f250e7e9de9fcb))
* Allow creation of saml2 integrations ([#616](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/616)) ([#805](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/805)) ([c07d582](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c07d5820bea7ac3d8a5037b0486c405fdf58420e))
* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))
* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* Allow setting resource monitor on account ([#768](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/768)) ([2613aa3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2613aa31da958e3557849e0615067c649c704110))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))
* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* handle serverless tasks ([#736](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/736)) ([bde252e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bde252ef2b225b128728e2cd4f2dcab62a65ba58))
* handle-account-grant-managed-task ([#751](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/751)) ([8952382](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8952382ca701cb5be19276b82bb740b997c0033a))
* Identity Column Support ([#726](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/726)) ([4da8014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4da801445d0523ce287c00600d1c1fd3f5af330f)), closes [#538](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/538)
* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))
* OAuth security integration for partner applications ([#763](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/763)) ([0ec5f4e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ec5f4ed993a4fa96b144924ddc34caa936819b0))
* Pipe and Task Grant resources ([#620](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/620)) ([90b9f80](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/90b9f80ea7fba568c595b87813324eef5bfa9d26))
* Procedures ([#619](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/619)) ([869ff75](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/869ff759eaaa50b364b41956af11e21fd130a4e8))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* Resource to manage a user's public keys ([#540](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/540)) ([590c22e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/590c22ec40ed28c7d280192ed66c4d93623e32fd))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))
* Support create function with Java language ([#798](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/798)) ([7f077f2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7f077f22c53b23cbed62c9b9284220a8f417f5c8))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))
* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))
* Table Column Defaults ([#631](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/631)) ([bcda1d9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bcda1d9cd3678647c056b5d79c7e7dd49a6380f9))
* table constraints ([#599](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/599)) ([b0417a8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b0417a80440f44833769e666fcf872a9dbd4ea74))
* Task error integration ([#830](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/830)) ([8acfd5f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8acfd5f0f3bcb383ddb74ea05636f84b5b215dbe))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))


### BugFixes

* Add AWS_SNS notification_provider support for error notifications for Snowpipe. ([#777](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/777)) ([02a97e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/02a97e051c804938a6a5137a34c0ff6c4fdc531f))
* Add AWS_SQS_IAM_USER_ARN to notification integration ([#610](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/610)) ([82a340a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82a340a356b7e762ea0beae3625fd6663b31ce33))
* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))
* Add importer to account grant ([#576](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/576)) ([a6d7f6f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a6d7f6fcf6b0e362f2f98f1fcde8b26221bf0cb7))
* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))
* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))
* Add valid property AWS_SNS_TOPIC_ARN to AWS_SNS notification provider  ([#783](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/783)) ([8224954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/82249541b1fb01cf686b7e0ff88e24f1b82e6ec0))
* add warehouses attribute to resource monitor ([#831](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/831)) ([b041e46](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b041e46c21c05597e600ac3cff316dac712442fe))
* Added Missing Account Privileges ([#635](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/635)) ([c9cc806](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c9cc80693c0884e120b62a7f097154dcf1d3490f))
* Allow empty result when looking for storage integration on refresh ([#692](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/692)) ([16363cf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/16363cfe9ea565e94b1cdc5814e31e95e1aa93b5))
* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
* build: Add trimpath Go flag to build ([#316](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/316)) ([420a84b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/420a84b61cf342e8210f440ccfaca5dcaa589ede))
* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* Dependabot configuration to make it easier to work with ([a7c60f7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a7c60f734fc3826b2a1444c3c7d17fdf8b6742c1))
* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))
* Escape String for AS in external table ([#580](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/580)) ([3954741](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3954741ed5ef6928bcb238dd8249fc072259db3f))
* **external_function:** Allow Read external_function where return_type is VARIANT ([#720](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/720)) ([1873108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/18731085333bfc83a1d729e9089c357873b9230c))
* external_table headers order doesn't matter ([#731](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/731)) ([e0d74be](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e0d74be5029f6bf73915dee07cadd03ac52bf135))
* Handling of task error_integration nulls ([#834](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/834)) ([3b27905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3b279055b66cd62f43da05559506f1afa282aa16))
* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))
* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))
* make platform info compatible with quoted identifiers ([#729](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/729)) ([30bb7d0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/30bb7d0214c58382b72b55f0685c3b0e9f5bb7d0))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* make saml2_enable_sp_initiated bool throughout ([#828](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/828)) ([b79988e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b79988e06ebc2faff5ad4667867df46fdbb89309))
* materialized view grant incorrectly requires schema_name ([#654](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/654)) ([faf0767](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/faf076756ec9fa348418fd938517c70578b1db11))
* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))
* OSCP -> OCSP misspelling ([#664](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/664)) ([cc8eb58](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/cc8eb58fceae64348d9e51bcc9258e011788484c))
* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))
* read Database and Schema name during Stream import ([#732](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/732)) ([9f747b5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9f747b571b2fcf5b0663696efd75c55a6f8b6a89))
* read Name, Database and Schema during Procedure import ([#819](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/819)) ([d17656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d17656fdd2803516b6ee067a6bd5a457bf31d905))
* Recreate notification integration when type changes ([#792](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/792)) ([e9768bf](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9768bf059268fb933ad74f2b459e91e2c0563e0))
* refactor ReadWarehouse function to correctly read object parameters ([#745](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/745)) ([d83c499](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d83c499910c0f2b6348191a93f917e450b9e69b2))
* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))
* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))
* Remove force_new from masking_expression ([#588](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/588)) ([fc3e78a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/fc3e78acbdda5346f32a004711d31ad6f68529ed))
* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))
* remove table where is_external is Y ([#667](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/667)) ([14b17b0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/14b17b00d47de1b971bf8967605ae38b348531f8))
* SCIM access token compatible identifiers ([#750](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/750)) ([afc92a3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/afc92a35eedc4ab054d67b75a93aeb03ef86cefd))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))
* Share example ([#673](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/673)) ([e9126a9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e9126a9757a7cf5c0578ea0d274ec489440132ca))
* Share resource to use REFERENCE_USAGE instead of USAGE ([#762](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/762)) ([6906760](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/69067600ac846930e06e857964b8a0cd2d28556d))
* Shares can't be updated on table_grant resource ([#789](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/789)) ([6884748](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/68847481e7094b00ab639f41dc665de85ed117de))
* **snowflake_share:** Can't be renamed, ForceNew on name changes ([#659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/659)) ([754a9df](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/754a9dfb7be5b64196f3c3015d32a5d675726ca9))
* Stream append only ([#653](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/653)) ([807c6ce](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/807c6ce566b08ba1fe3b13eb84e1ae0cf9cf69a8))
* table: Properly read and import table state ([#314](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/314)) ([df1f6bc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/df1f6bcabfca27058c860a7db815d263457afd6c))
* tagging for db, external_table, schema ([#795](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/795)) ([7aff6a1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7aff6a1e04358790a3890e8534ea4ffbc414024b))
* Update go and docs package ([#1009](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1009)) ([72c3180](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72c318052ad6c29866cfee01e9a50a1aaed8f6d0))
* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))
* update ReadTask to correctly set USER_TASK_TIMEOUT_MS ([#761](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/761)) ([7b388ca](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7b388ca4957880e7204a15536e2c6447df43919a))
* Upgrade go ([#715](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/715)) ([f0e59c0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f0e59c055d32d5d152b4c2c384b18745b8e9ef0a))
* Upgrade tf for testing ([#625](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/625)) ([c03656f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c03656f8e97df3f8ba93cd73fcecc9702614e1a0))
* use "DESCRIBE USER" in ReadUser, UserExists ([#769](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/769)) ([36a4f2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/36a4f2e3423fb3c8591d8e96f7a5e1f863e7fea8))
* Warehouse create and alter properties ([#598](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/598)) ([632fd42](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/632fd421f8acbc358d4dfd5ae30935512532ba64))


### Misc

* **main:** release 0.26.0 ([#841](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/841)) ([4a6d659](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4a6d659c7069c1d2d64c43ce3262d3a7a840b7c5))
* **main:** release 0.26.1 ([#849](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/849)) ([a2607e5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a2607e5e15d6dfd66e756e381c0aeccf8106bbd4))
* **main:** release 0.26.2 ([#851](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/851)) ([016e02d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/016e02d3cc51360ecae43df6a931ada2b398e424))
* **main:** release 0.26.3 ([#854](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/854)) ([63f2b85](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f2b8507930b9429ebf7c8ee8f65334ef16931e))
* **main:** release 0.27.0 ([#870](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/870)) ([5178aa6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5178aa6a977fe296bc4b5abeae6e55ca27291dca))
* **main:** release 0.28.0 ([#886](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/886)) ([b523f7e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b523f7e263f988a839528bb19b04405890013879))
* **main:** release 0.28.1 ([#895](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/895)) ([c92c518](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c92c5184eab43141116d760f9e336714eb535fd7))
* **main:** release 0.28.2 ([#902](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/902)) ([e1c228b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e1c228b9bbba60d297082b665159ca9160607e62))
* **main:** release 0.28.3 ([#905](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/905)) ([b01a21b](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b01a21bc7fa2055bbaf77af8e263e69fbb1bfd54))
* **main:** release 0.28.4 ([#913](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/913)) ([4fa19e3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/4fa19e37edda8c3958232c647d35bf99a7d00319))
* **main:** release 0.28.5 ([#915](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/915)) ([d9a7f01](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/d9a7f0165cc440cc7ed6086d033ab7252e56c5e2))
* **main:** release 0.28.6 ([#920](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/920)) ([3a17e34](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/3a17e34a9e762ee24d8ff12144fe6c6ac4b68e2e))
* **main:** release 0.28.7 ([#921](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/921)) ([adbb52e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/adbb52e3f33c86519ed20f490bddd84980437e3f))
* **main:** release 0.28.8 ([#928](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/928)) ([96152d7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/96152d7db14c397bff9f9cc79ba0d85f6f2706b4))
* **main:** release 0.29.0 ([#943](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/943)) ([f1d0af9](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1d0af96bbe8e57558bd3a57808298d8d49fe349))
* **main:** release 0.30.0 ([#954](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/954)) ([bfd3108](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bfd31080b96af02f908e93ff0c20b8cb24840431))
* **main:** release 0.31.0 ([#968](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/968)) ([1e21100](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1e2110065c23d851e04cd2de1727b683a38168f1))
* **main:** release 0.32.0 ([#974](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/974)) ([e947417](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e947417c459e424829a9b9e4cbb96f04ba7db3cd))
* **main:** release 0.33.0 ([#988](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/988)) ([bf3482e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/bf3482e4de81e96b31aec192a15f5bee33d34e78))
* **main:** release 0.33.1 ([#991](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/991)) ([1c5af87](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1c5af874164d8b40e7cae54e9206ec6b46c2e75b))
* **main:** release 0.34.0 ([#1014](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1014)) ([f1c651e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f1c651e17d1697f37be43857318573cb56812f5d))
* **main:** release 0.34.0 ([#1019](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1019)) ([83db3a4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83db3a4c14ec6f1539fbef55c72ae36b22e76906))
* **main:** release 0.34.0 ([#1020](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/1020)) ([7116025](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7116025e3523cc6d385752f7e71bff1b5fded68b))
* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))
* release 0.34.0 ([836dfcb](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/836dfcb28020519a5c4dee820f61581c65b4f3f2))
* Update go files ([#839](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/839)) ([5515443](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/55154432dd5424b6d37b04163613b6db94fda70e))
* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))

### [0.33.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.33.0...v0.33.1) (2022-05-03)


### BugFixes

* Network Attachment (Set For Account) ([#990](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/990)) ([1dde150](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1dde150fdc74937b67d6e94d0be3a1163ac9ebc7))

## [0.33.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.32.0...v0.33.0) (2022-04-28)


### Features

* Add 'snowflake_role' datasource ([#986](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/986)) ([6983d17](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6983d17a47d0155c82faf95a948ebf02f44ef157))

## [0.32.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.31.0...v0.32.0) (2022-04-14)


### Features

* allow in-place renaming of Snowflake schemas ([#972](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/972)) ([2a18b96](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2a18b967b92f716bfc0ae788be36ce762b8ab2f4))

## [0.31.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.30.0...v0.31.0) (2022-04-11)


### Features

* Add manage support cases account grants ([#961](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/961)) ([1d1084d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/1d1084de4d3cef4f76df681812656dd87afb64df))
* snowflake_user_ownership_grant resource ([#969](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/969)) ([6f3f09d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f3f09d37bad59b21aacf7c5d59de020ed47ecf2))

## [0.30.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.29.0...v0.30.0) (2022-03-29)


### Features

* support host option to pass down to driver ([#939](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/939)) ([f75f102](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f75f102f04d8587a393a6c304ea34ae8d999882d))

## [0.29.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.8...v0.29.0) (2022-03-23)


### Features

* Allow in-place renaming of Snowflake tables ([#904](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/904)) ([6ac5188](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6ac5188f62be3dcaf5a420b0e4277bd161d4d71f))
* create snowflake_role_ownership_grant resource ([#917](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/917)) ([17de20f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/17de20f5d5103ccc518ce07cb58a1e9b7cea2865))


### BugFixes

* Legacy role grantID to work with new grant functionality ([#941](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/941)) ([5182361](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/5182361c48463325e7ad830702ad58a9617064df))

### [0.28.8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.7...v0.28.8) (2022-03-18)


### BugFixes

* change the function_grant documentation example privilege to usage ([#901](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/901)) ([70d9550](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/70d9550a7bd05959e709cfbc440d3c72844457ac))
* remove share feature from stage because it isn't supported ([#918](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/918)) ([7229387](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7229387e760eab4ba4316448273b000be514704e))

### [0.28.7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.6...v0.28.7) (2022-03-15)


### BugFixes

* Allow legacy version of GrantIDs to be used with new grant functionality ([#923](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/923)) ([b640a60](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/b640a6011a1f2761f857d024d700d4363a0dc927))
* Make ReadWarehouse compatible with quoted resource identifiers ([#907](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/907)) ([72cedc4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/72cedc4853042ff2fbc4e89a6c8ee6f4adb35c74))
* sequence import ([#775](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/775)) ([e728d2e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e728d2e70d25de76ddbf274bcd2c3fc989c7c449))

### [0.28.6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.5...v0.28.6) (2022-03-14)


### BugFixes

* Add release step in goreleaser ([#919](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/919)) ([63f221e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/63f221e6c2db8ceec85b7bca71b4953f67331e79))

### [0.28.5](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.4...v0.28.5) (2022-03-12)


### BugFixes

* Add manifest json ([#914](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/914)) ([c61fcdd](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c61fcddef12e9e2fa248d5da8df5038cdcd99b3b))

### [0.28.4](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.3...v0.28.4) (2022-03-11)


### BugFixes

* Add gpg signing to goreleaser ([#911](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/911)) ([8ae3312](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/8ae3312ea09233323ac96d3d3ade755125ef1869))

### [0.28.3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.2...v0.28.3) (2022-03-10)


### BugFixes

* issue with ie-proxy ([#903](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/903)) ([e028bc8](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/e028bc8dde8bc60144f75170de09d4cf0b54c2e2))

### [0.28.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.1...v0.28.2) (2022-03-09)


### BugFixes

* Ran make deps to fix dependencies ([#899](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/899)) ([a65fcd6](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/a65fcd611e6c631e026ed0560ed9bd75b87708d2))

### [0.28.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.28.0...v0.28.1) (2022-03-09)


### BugFixes

* Release by updating go dependencies ([#894](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/894)) ([79ec370](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/79ec370e596356f1b4824e7b476fad76d15a210e))

## [0.28.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.27.0...v0.28.0) (2022-03-05)


### Features

* Implemented External OAuth Security Integration Resource ([#879](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/879)) ([83997a7](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/83997a742332f1376adfd31cf7e79c36c03760ff))


### BugFixes

* escape string escape_unenclosed_field ([#877](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/877)) ([6f5578f](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6f5578f55221f460f1dcc2fa48848cddea5ade20))

## [0.27.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.26.3...v0.27.0) (2022-03-02)


### Features

* Data source for list databases ([#861](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/861)) ([537428d](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/537428da16024707afab2b989f95f2fe2efc0e94))
* Expose GCP_PUBSUB_SERVICE_ACCOUNT attribute in notification integration ([#871](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/871)) ([9cb863c](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/9cb863cc1fb27f76030984917124bcbdef47dc7a))
* Support DIRECTORY option on stage create ([#872](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/872)) ([0ea9a1e](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/0ea9a1e0fb9617a2359ed1e1f60b572bd4df49a6))


### Misc

* Upgarde all dependencies to latest ([#878](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/878)) ([2f1c91a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/2f1c91a63859f8f9dc3075ab20aa1ded23c16179))

### [0.26.3](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.26.2...v0.26.3) (2022-02-08)


### BugFixes

* Remove keybase since moving to github actions ([#852](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/852)) ([6e14906](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/6e14906be91553c62b24e9ab0e8da7b12b37153f))

### [0.26.2](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.26.1...v0.26.2) (2022-02-07)


### BugFixes

* Update goreleaser env Dirty to false ([#850](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/850)) ([402f7e0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/402f7e0d0fb19d9cbe71f384883ebc3563dc82dc))

### [0.26.1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.26.0...v0.26.1) (2022-02-07)


### BugFixes

* Release tag ([#848](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/848)) ([610a85a](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/610a85a08c8c6c299aed423b14ecd9d115665a36))

## [0.26.0](https://github.com/Snowflake-Labs/terraform-provider-snowflake/compare/v0.25.36...v0.26.0) (2022-02-03)


### Features

* Add replication support ([#832](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/832)) ([f519cfc](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/f519cfc1fbefcda27da85b6a833834c0c9219a68))
* Release GH workflow ([#840](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/840)) ([c4b81e1](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/c4b81e1ec45749ed113143ec5a26ab0ad2fb5906))
* TitleLinter customized ([#842](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/842)) ([39c7e20](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/39c7e20108e6a8bb5f7cb98c8bd6a022d20f8f40))


### Misc

* Move titlelinter workflow ([#843](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/843)) ([be6c454](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/be6c4540f7a7bc25653a69f41deb2c533ae9a72e))


### BugFixes

* Allow multiple resources of the same object grant ([#824](https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/824)) ([7ac4d54](https://github.com/Snowflake-Labs/terraform-provider-snowflake/commit/7ac4d549c925d98f878cffed2447bbbb27379bd8))
