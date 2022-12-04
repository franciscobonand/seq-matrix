# Sequence Matrix Verifier

This application's purpose is to check if a sequence of strings (interpreted as an NxN matrix) given as input is valid or not.

## How it works

When executed, the application starts a server on `localhost:9001`, which provides the endpoint `/sequence`.
This endpoint expects a POST request with a JSON body in the given format:

```json
{
    "letters": ["DUHBHB", "DUBUHD", "UBUUHU", "BHBDHH", "DDDDUB", "UDBDUH"]
}
```

If the input is valid, it returns a JSON object containing the field `{ "is_valid": boolean }` to the user, and also stores this result in a MongoDB collection.

### Invalid inputs

If the input has one of the following characteristics it's considered invalid, and an error status code and message are returned:

- Input isn't a JSON object or doesn't contain the string array field `"letters"`.
- Some string in the `"letters"` array contains a character different than "B", "U", "D" or "H".
- The length of any item is different than the length of the array - **each string in the array represents a line in an NxN matrix, so these lengths must be the same**.
- The length is less than four. The reason for that is better explained in the next topic.

### Sequence validation

For a sequence to be considered valid, the NxN matrix that represents it must contain at least two sets of four letters repeated in a row, in any direction (vertical, horizontal or diagonal).

Example of a valid sequence:

```json
{
    "letters": ["DUHBHB", "DUBUHD", "UBUUHU", "BHBDHH", "DDDDUB", "UDBDUH"]
}
```

![Valid sequence example](images/valid_seq.jpg)

### Stats

Alongside with the `/sequence` endpoint, the server also provides another endpoint: `/stats`.
This endpoint expects a GET request with no body, and returns to the user some information about the previously analyzed sequences.

Response example:

```json
{
    "count_valid": 40,
    "count_invalid": 60,
    "ratio": 0.4
}
```

## How sequences are validated

As previously said, the string array which makes a sequence is interpreted as an NxN matrix. With this, the search for four repeated letters is made in four distinct directions of the matrix: rows, columns, primary diagonal and secondary diagonal:

![Validation directions](images/validation_dirs.jpg)

It's important to notice two things:

- As the validation doesn't require the values in the matrix to be modified, each direction is independent, so the validation is done concurrently/in parallel for better performance.
- When validating, it looks for four or more repeated letters. Because of this, the first and last three lines of the primary and secondary diagonals are not considered (as these lines' length is less than four).

> Known issue: For matrices 9x9 or bigger, two repetitons in the same line are only counted once.  
> Example: "DDDDHDDDD" isn't enought for the sequence to be considered valid.
