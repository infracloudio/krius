# Changelog

## [v0.0.1](https://github.com/infracloudio/krius/tree/v0.0.1) (2021-11-22)

[Full Changelog](https://github.com/infracloudio/krius/compare/241d8b7f8c9e4e6a43ee8663932e9a06a03d2751...v0.0.1)

**Implemented enhancements:**

- Support Ruler component in Thanos config [\#48](https://github.com/infracloudio/krius/issues/48)
- Support memcached/in-memory for caching query results in Query Frontend [\#45](https://github.com/infracloudio/krius/issues/45)

**Fixed bugs:**

- Bugs in generating a spec file [\#62](https://github.com/infracloudio/krius/issues/62)
- pre-flight checks should finish for all contexts before starting installation [\#33](https://github.com/infracloudio/krius/issues/33)

**Closed issues:**

- change command "krius uninstall" to "krius destroy"  [\#61](https://github.com/infracloudio/krius/issues/61)
- Initiate spec generation [\#30](https://github.com/infracloudio/krius/issues/30)
- The commands in README.md don't work.  [\#28](https://github.com/infracloudio/krius/issues/28)
- Add release process/script [\#27](https://github.com/infracloudio/krius/issues/27)
- Obj Store Configure Command [\#18](https://github.com/infracloudio/krius/issues/18)
- Helm need to be replaced with helm sdk [\#9](https://github.com/infracloudio/krius/issues/9)
- SingleCluster: Install prometheus-community/kube-prometheus-stack chart when `krius install prometheus` is given [\#3](https://github.com/infracloudio/krius/issues/3)
- Create a CLI structure require to install the components  [\#2](https://github.com/infracloudio/krius/issues/2)
- Create a CLI structure, require to install the components    [\#1](https://github.com/infracloudio/krius/issues/1)

**Merged pull requests:**

- Documenting [\#72](https://github.com/infracloudio/krius/pull/72) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Support targets and remote URL [\#71](https://github.com/infracloudio/krius/pull/71) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Added support for multiple obj storage/ upgrading obj storage [\#70](https://github.com/infracloudio/krius/pull/70) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Fix/change spec uninstall destroy [\#69](https://github.com/infracloudio/krius/pull/69) ([JESWINKNINAN](https://github.com/JESWINKNINAN))
- Logging framework [\#67](https://github.com/infracloudio/krius/pull/67) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Generate command fixes [\#66](https://github.com/infracloudio/krius/pull/66) ([YachikaRalhan](https://github.com/YachikaRalhan))
- fix: Bug fixes generate command [\#64](https://github.com/infracloudio/krius/pull/64) ([JESWINKNINAN](https://github.com/JESWINKNINAN))
- fixed the default thanos reciever port [\#60](https://github.com/infracloudio/krius/pull/60) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Move krius logo to center and increase size [\#59](https://github.com/infracloudio/krius/pull/59) ([sanketsudake](https://github.com/sanketsudake))
- Added default port for sidecar accessibility [\#58](https://github.com/infracloudio/krius/pull/58) ([YachikaRalhan](https://github.com/YachikaRalhan))
- No need for pre flight checks when uninstalling spec [\#57](https://github.com/infracloudio/krius/pull/57) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Feature/storage secret conflicts [\#56](https://github.com/infracloudio/krius/pull/56) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Supoort upgrading Thanos chart by install flag [\#55](https://github.com/infracloudio/krius/pull/55) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Support for Ruler component in Thanos [\#54](https://github.com/infracloudio/krius/pull/54) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Added svg logo in readme [\#53](https://github.com/infracloudio/krius/pull/53) ([YachikaRalhan](https://github.com/YachikaRalhan))
- feat: spec describe cluster with stack deployed status details [\#51](https://github.com/infracloudio/krius/pull/51) ([JESWINKNINAN](https://github.com/JESWINKNINAN))
- docs: Update readme with details of commands [\#50](https://github.com/infracloudio/krius/pull/50) ([JESWINKNINAN](https://github.com/JESWINKNINAN))
- Updated helm version to 3.6.1 and other dependencies for security purpose [\#49](https://github.com/infracloudio/krius/pull/49) ([YachikaRalhan](https://github.com/YachikaRalhan))
- added support for caching in Querier FE [\#47](https://github.com/infracloudio/krius/pull/47) ([YachikaRalhan](https://github.com/YachikaRalhan))
- feat: describe-cluster showing meta-details part from specfile [\#46](https://github.com/infracloudio/krius/pull/46) ([JESWINKNINAN](https://github.com/JESWINKNINAN))
- Uninstall spec implementation for multicluster [\#44](https://github.com/infracloudio/krius/pull/44) ([praddy26](https://github.com/praddy26))
- Added release process [\#43](https://github.com/infracloudio/krius/pull/43) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Added getting started details to readme such as generate a spec file,… [\#42](https://github.com/infracloudio/krius/pull/42) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Thanos Receiver implementation [\#38](https://github.com/infracloudio/krius/pull/38) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Spec generate command implemented [\#37](https://github.com/infracloudio/krius/pull/37) ([vaibhavp](https://github.com/vaibhavp))
- Thanos sidecar multi cluster setup [\#35](https://github.com/infracloudio/krius/pull/35) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Add golangci-lint checks for linting code [\#34](https://github.com/infracloudio/krius/pull/34) ([sanketsudake](https://github.com/sanketsudake))
- Prometheus setup with thanos sidecar [\#31](https://github.com/infracloudio/krius/pull/31) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Read yaml config file and unmarshal the clusters to respective structs [\#29](https://github.com/infracloudio/krius/pull/29) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Added Rule Schema for spec [\#26](https://github.com/infracloudio/krius/pull/26) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Add a release flag for passing in the release name [\#24](https://github.com/infracloudio/krius/pull/24) ([hellozee](https://github.com/hellozee))
- Configure obj and helm SDK changes [\#21](https://github.com/infracloudio/krius/pull/21) ([YachikaRalhan](https://github.com/YachikaRalhan))
- Added helm sdk to install helm chart [\#17](https://github.com/infracloudio/krius/pull/17) ([girishg4t](https://github.com/girishg4t))
- Added logic to install Prometheus chart [\#16](https://github.com/infracloudio/krius/pull/16) ([girishg4t](https://github.com/girishg4t))
- Created a basic project stucture and cli commands  [\#7](https://github.com/infracloudio/krius/pull/7) ([girishg4t](https://github.com/girishg4t))



