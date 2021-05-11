# :alarm_clock: timetrace

> timetrace is a simple CLI for tracking your working time.

![CLI screenshot](screenshot.png)

## Installation

### Homebrew

```
brew tap timetrace/timetrace
brew install timetrace
```

### Docker

```
docker container run -v ${HOME}:/data timetrace/timetrace
```

### Binary

Download the [latest release](https://github.com/timetrace/timetrace/releases)
and extract the binary into a directory like `/usr/local/bin` or
`C:\Program Files\timetrace`. Make sure the directory is in the `PATH` variable.

## Usage example

First, create a project you're working for:

```
timetrace create project make-coffee
```

Once the project is created, you're able to track work on that project.

```
timetrace start make-coffee
```

You can obtain your currently worked time using `timetrace status`. When you've
finished your work, stop tracking:

```
timetrace stop
```

You're also able to delete and edit projects and records (see below).

## Command reference

### Create a project

**Syntax:**

```
timetrace create project <KEY>
```

**Arguments:**

|Argument|Description|
|-|-|
|`KEY`|An unique project key.|

**Example:**

Create a project called `make-coffee`:

```
timetrace create project make-coffee
```

### Get a project

**Syntax:**

```
timetrace get project <KEY>
```

**Arguments:**

|Argument|Description|
|-|-|
|`KEY`|The project key.|

**Example:**

Display a project called `make-coffee`:

```
timetrace get project make-coffee
```

### Get a record

**Syntax:**

```
timetrace get record <YYYY-MM-DD-HH-MM>
```

**Arguments:**

|Argument|Description|
|-|-|
|`YYYY-MM-DD-HH-MM`|The start time of the desired record.|

**Example:**

At the moment, the start time has to be in 24h format, meaning 3PM = 15. Display
the record created on May 1st 2021, 3:00 PM:

```
timetrace get record 2021-05-01-15-00
```

### Edit a project

**Syntax:**

```
timetrace edit project <KEY>
```

**Arguments:**

|Argument|Description|
|-|-|
|`KEY`|The project key.|

**Example:**

Edit a project called `make-coffee`:

```
timetrace edit project make-coffee
```

### Delete a project

**Syntax:**

```
timetrace delete project <KEY>
```

**Arguments:**

|Argument|Description|
|-|-|
|`KEY`|The project key.|

**Example:**

Delete a project called `make-coffee`:

```
timetrace delete project make-coffee
```

### Start tracking

**Syntax:**

```
timetrace start <PROJECT KEY>
```

**Arguments:**

|Argument|Description|
|-|-|
|`PROJECT KEY`|The key of the project.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--billable`|`-b`|Mark the record as billable.|

**Example:**

Start working on a project called `make-coffee` and mark it as billable:

```
timetrace start --billable make-coffee
```

### Print current status

**Syntax:**

```
timetrace status
```

**Example:**

Print the current tracking status:

```
timetrace status
+-------------+--------------------+---------------+
|   PROJECT   | WORKED SINCE START | WORKED TODAY  |
+-------------+--------------------+---------------+
| make-coffee | 3h25m29.037343s    | 7h22m49.5749s |
+-------------+--------------------+---------------+
```

### Stop tracking

**Syntax:**

```
timetrace stop
```

**Example:**

Stop working on your current project:

```
timetrace stop
```

### Print version information

**Syntax:**

```
timetrace version
```

**Example:**

Print your installed timetrace version:

```
timetrace version
```

## Credits

This project uses [spf13/cobra](https://github.com/spf13/cobra),
[fatih/color](https://github.com/fatih/color),
[olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) and
[enescakir/emoji](https://github.com/enescakir/emoji).