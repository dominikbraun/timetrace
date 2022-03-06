# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.14.3] - 2022-03-06

### Fixed
* Fix output for hour format (#208)

## [0.14.2] - 2021-10-27

### Changed
* Make `tzdata` available in timetrace Docker image

## [0.14.1] - 2021-10-22

### Fixed
* Fixed panic when providing no argument for `timetrace start` (#199)

## [0.14.0] - 2021-10-02

### Added
* Add support for tags (#183)
* Add support for durations in decimal format (#180)
* Introduce `@ID` notation for selecting records (#116)

### Changed
* Print colliding records as table when a collision occurs (#174)

## [0.13.0] - 2021-08-10

### Added
* Add support for per-project configuration (#166)
* Add support for making projects billable by default (#166)
* Add support for deleting and reverting projects with their modules (#160)
* Add support for restoring records associated with a restored project (#160)
* Add `--non-billable` flag for `timetrace start` command (#172) 
* Add `--non-billable` flag for `timetrace report` command (#142)

### Changed
* Ask the user whether they want to delete all records when deleting a project (#160)

## [0.12.0] - 2021-07-30

### Changed
* `timetrace report`: Include project modules when filtering for a project (#143)

### Fixed
* Don't consider any backup records for calculations (#162)

### Other
* Refactorings and optimizations (#146, #164)
* Preparations for the larger and relevant 0.13.0 release

## [0.11.1] - 2021-07-05

### Changed
* Default to 'no' when asking the user for confirmation (#147)
* Ask the user for confirmation when deleting a project (#147, #155)
* Don't allow creation of records in the future (#152)

### Fixed
* Fix nil pointer dereference when checking record collisions (#150)
* Fix typos in documentation (#149)

## [0.11.0] - 2021-06-26

### Added
* Add `create record` command (#118)
* Add beta support for reports (#99)

## [0.10.0] - 2021-06-15

### Added
* Add `--output` flag for `timetrace status` (#129)
* Add `--format` flag for `timetrace status` (#114, #123, #124)
* Add Starship support
* Add Scoop support

### Changed
* Remove seconds from printed durations (#113)

## [0.9.0] - 2021-06-09

### Added
* Add the total tracked time for all listed records to `timetrace list records` (#106)

## [0.8.0] - 2021-06-06

### Added
* Introduce the `--revert` flag for `edit record`, `delete record`, `edit project` and `delete project` (#93)
* Add the overall break time for `timetrace status` (#100)

## [0.7.2] - 2021-06-04

### Fixed
* Fix critical error when starting tracking if there are no existing records (#103)

## [0.7.1] - 2021-06-02

### Fixed
* Fix unhandled error in `edit record` command if there are no tracked records (#96)

## [0.7.0] - 2021-05-30

### Added
* Add support for project modules

### Changed
* Consider project modules when filtering projects (#63)
* Display project modules when listing projects (#70)
* Require parent projects to exist when creating a module (#80)

## [0.6.1] - 2021-05-26

### Fixed
* Fix `timetrace stop` command (#86)

## [0.6.0] - 2021-05-26

### Added
* Add `timetrace edit record` command (#51)
* Add `latest` alias for `timetrace edit record` (#73)

### Changed
* Use an info output for `timetrace status` when there's no active tracking (#65)
* Always adhere to the `use12hours` setting for date- and time input and output (#67)

### Fixed
* Don't allow edting of incomplete records (#69)
* Don't allow re-creation of existing projects (#79)

## [0.5.0] - 2021-05-22

### Added
* Add record key to `list records` output (#44)
* Add `--project` filter to `list records` (#62)

### Changed
* Make command syntax and help sensible to `use12hours` in `config.yaml` (#44, #54)

### Fixed
* Fix Docker image labels (#35)

## [0.4.0] - 2021-05-20

### Added
* Add support for Bash autocompletion (#25)
* Add support for Snap (#31)
* Add filter for billable records (#33)

### Changed
* Use non-error output for `timetrace version` (#45)
* Colorize and stylize tables (#26, #48, #49)

## [0.3.1] - 2021-05-18

### Changed
* Introduce the `today` and `yesterday` aliases for `timetrace list records` (#32)
* Preview the target record and ask for confirmation when deleting a record (#37)

## [0.3.0] - 2021-05-17

### Added
* `timetrace delete record` command (#22)

### Changed
* Add project name to the output of `list records` (#24)

## [0.2.1] - 2021-05-16

### Changed
* Beautified output of time durations (#20)

### Fixed
* Fix `timetrace status` if there was no time tracking yet (#18)

## [0.2.0] - 2021-05-16

### Added
* `timetrace list records` command (#16)

## [0.1.4] - 2021-05-14

### Fixed
* Fix calculation of total tracked time

## [0.1.3] - 2021-05-14

### Fixed
* Fix determination of latest stored record

## [0.1.2] - 2021-05-14

### Fixed
* Support slashes and back-slashes in project keys (#9)

## [0.1.1] - 2021-05-14

### Fixed
* Fix status report for situations where time tracking is not active (#7)

## [0.1.0] - 2021-05-13

### The initial timetrace release.
