Scrub Personal Identifying Information Problem
==================================================

#####Note: The input files contain unstructured JSON data

Application event tracking to external services often involves meta data about the user who took the action. When this event information is sent to a third party, we want to scrub all personally identifying data reported, without losing information about what fields were actually recorded.

In order to do this on different platforms, it is useful to implement a method called `scrub()` that will take a JSON object and an array of sensitive fields as an input and returns a modified JSON object that replaces alphanumeric characters with an asterisk ("\*") from values corresponding to keys matching the sensitive field names.

For example, if the sensitive fields we want to "scrub" are `name`, `phone`, and `email`, a JSON input of:

```
{
  "name": "Kelly Doe",
  "email": "kdoe@example.com",
  "id": 12324,
  "phone": "5551234567"
}
```

should return:

```
{
  "name": "***** ***",
  "email": "****@*******.***",
  "id": 12324,
  "phone": "**********"
}

```

## Requirements
You will need to implement a command line executable that takes two arguments: a text file with a list of sensitive fields and a JSON file of user data to scrub. Calling `./scrub sensitive_fields.txt input.json ` should output a scrubbed JSON version of `input.json` with the keys in `sensitive_fields.txt` "scrubbed".

Any valid JSON input should be able to be handled. Value types for sensitive keys should be handled as follows:
  - `String`: replace alphanumeric characters with "*"
  - `Number`: convert to string and replace alphanumeric characters with "*"
  - `Boolean`: replace entire value with "-"
  - `Array`: each value of the array should be evaluated as described by other field types
  - `Object`: if the key matches a sensitive field, all values of the nested object should be scrubbed as described by other field types. If the nested object does not correspond to a sensitive key, each key/value pair of the nested object should be evaluated as described by other field types
  - `null`: value should be unmodified

## Tests

A handful of example test "scrubs" are in the `/tests/` directory. Each subdirectory has an `input.json`, `sensitive_fields.txt`, and corresponding `output.json` that is expected.

---

## Sample Runs

### 1. Usage
```
$ ./scrub-pii   
Usage: scrub-pii <input json file> <sensitive fields file>
```
Command option: `-h`
```
$ ./scrub-pii -h
Usage of ./scrub-pii:
  -pretty
    	display pretty output; otherwise do: -pretty=false (default true)
```

### 2. Sample run on test `07_mixed_type_arrays`
![Scrub 07_mixed_type_arrays](https://github.com/developertogo/scrub-pii/blob/master/docs/assets/sample-pretty-test-07-run.jpg)

### 3. Sample run on test `06_nested_object` with no pretty output (i.e. -pretty=false)
![Scrub 06_nested_object](https://github.com/developertogo/scrub-pii/blob/master/docs/assets/sample-no-pretty-unit-test-run.jpg)

### 4. Sample run `all unit tests` run
![All unit tests run](https://github.com/developertogo/scrub-pii/blob/master/docs/assets/sample-unit-test-run.jpg)

### 5. Sample run `all unit tests` run with -v (`verbose` mode enabled)
![Verbosed all unit tests run](https://github.com/developertogo/scrub-pii/blob/master/docs/assets/sample-verbose-unit-test-run.jpg)