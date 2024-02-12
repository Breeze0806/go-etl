# Developer Guide for Class Two-Dimensional Table File Storage

Class Two-Dimensional Table File Storage is a framework that abstracts the flow of class two-dimensional table files. This framework can support reading and writing various class two-dimensional table file formats.

## Input File Stream

```go
// Opener is an interface for an opener that can open an input stream
type Opener interface {
 Open(filename string) (stream InStream, err error) // Open an input stream with the filename
}

// InStream is an interface for an input stream
type InStream interface {
 Rows(conf *config.JSON) (rows Rows, err error) // Get a row reader
 Close() (err error)                            // Close the input stream
}

// Rows is an interface for a row reader
type Rows interface {
 Next() bool                                  // Get the next row, return false if there is no next row, true if there is
 Scan() (columns []element.Column, err error) // Scan each row's columns
 Error() error                                // Get the error of the next row
 Close() error                                // Close the row reader
}
```

The InStream input stream can obtain a row reader Rows by passing in a JSON configuration file, which converts a row of data into a record in Rows. For implementation details, refer to the csv package. Additionally, it is necessary to register the Opener.

```go
func init() {
 var opener Opener
 file.RegisterOpener("csv", &opener)
}
```

## Output File Stream

```go
// Creator is an interface for a creator that can create an output stream
type Creator interface {
 Create(filename string) (stream OutStream, err error) // Create an output stream with the filename
}

// OutStream is an interface for an output stream
type OutStream interface {
 Writer(conf *config.JSON) (writer StreamWriter, err error) // Create a writer
 Close() (err error)                                        // Close the output stream
}

// StreamWriter is an interface for an output stream writer
type StreamWriter interface {
 Write(record element.Record) (err error) // Write a record
 Flush() (err error)                      // Flush to the file
 Close() (err error)                      // Close the output stream writer
}
```

The OutStream output stream can obtain an output stream writer StreamWriter by passing in a JSON configuration file, which converts a record into a row of data in StreamWriter. For implementation details, refer to the csv package. Additionally, it is necessary to register the Creator.

```go
func init() {
 var creator Creator
 file.RegisterCreator("csv", &creator)
}
```