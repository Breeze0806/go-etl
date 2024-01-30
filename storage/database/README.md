# Guide for Developers of Database Storage

Database storage is a framework for database queries and SQL execution, serving as an abstraction for databases. Its underlying implementation relies on the `database/sql` interface of the Golang standard library.

## Introduction to Database Storage

Database storage facilitates the querying and execution of SQL through the encapsulation of the `database/sql` DB struct from the Golang standard library. The `db.go` file provides a rich set of methods, including not only the original methods of the `database/sql` DB such as `BeginTx`, `PingContext`, `QueryContext`, and `ExecContext` but also `FetchTable` and `FetchTableWithParam` for retrieving table structures, `FetchRecord` and `FetchRecordWithTx` for fetching records from a table, and `BatchExec`, `BatchExecWithTx`, and `BatchExecStmtWithTx` for executing write operations.

However, for different databases, the implementation of database storage can vary based on the specific database dialect. This document will introduce how to implement the database dialect interface.

## Introduction to the Database Dialect Interface

The prerequisite for implementing the dialect interface is that **the corresponding database driver can implement the `database/sql` interface of the Golang standard library**.

When implementing specifically, you can refer to the following directory structure. Here, MySQL is used as an example:

```go
storage--database--mysql----+--config.go        
                            |--doc.go
                            |--field.go
                            |--source.go
                            |--table.go
```

Using this approach, we have currently implemented support for MySQL, PostgreSQL, and DB2.

### Data Source Interface

```golang
// Dialect represents a database dialect.
type Dialect interface {
 Source(*BaseSource) (Source, error) // Data source
}

// Source represents a data source, including driver information, package information, configuration files, and connection information.
type Source interface {
 Config() *config.JSON   // Configuration information
 Key() string            // Generally, the connection information
 DriverName() string     // Driver name, used as the 1st parameter for sql.Open
 ConnectName() string    // Connection information, used as the 2nd parameter for sql.Open
 Table(*BaseTable) Table // Get the specific table structure interface
}
```

When implementing the `Source` interface, you can combine `BaseSource` to simplify the implementation. The `Table` method should return the specific table structure interface. Refer to the implementation in `source.go` of the MySQL package.

Additionally, the connection information relies on the configuration provided by `Config`. Currently, `Config` needs to be defined as follows to be compatible with the `dbms` package for implementing DataX plugins. Refer to the implementation in `config.go` of the MySQL package.

```go
type Config struct {
 URL      string `json:"url"`      // Database URL, including database address and other parameters
 Username string `json:"username"` // Username
 Password string `json:"password"` // Password
}
```

Furthermore, you need to use the `init` function to register the specific dialect:

```go
func init() {
 var d Dialect
 database.RegisterDialect(d.Name(), d)
}
```

### Table Structure Interface

```go
// Table represents a table structure.
type Table interface {
 fmt.Stringer

 Quoted() string   // Fully qualified quoted table name
 Instance() string // Instance name, e.g., database for MySQL
 Schema() string   // Schema name, e.g., username (schema name) for Oracle
 Name() string     // Table name, e.g., table for MySQL
 Fields() []Field  // Displays all columns
}

// FieldsFetcher is a supplementary interface for Table, used to specifically fetch all columns.
type FieldsFetcher interface {
 FetchFields(ctx context.Context, db *DB) error // Retrieves the specific columns
}

// FieldAdder is a supplementary interface for Table, used to add new columns to the table.
type FieldAdder interface {
 AddField(*BaseField) // Adds a specific column
}

// ExecParameter is a supplementary interface for Table, used to generate SQL statements for write operations.
type ExecParameter interface {
 ExecParam(string, *sql.TxOptions) (Parameter, bool)
}
```

When implementing the `Table` interface, you can combine `BaseTable` to simplify the implementation. The `Fields` method must return a collection of specific field interfaces for the corresponding database. Refer to the implementation in `table.go` of the MySQL package.

You can choose to implement either `FetchFields` or `FieldAdder`, but generally, `FieldAdder` is preferred. `ExecParameter` can be used to implement SQL statements for bulk inserts. For example, for MySQL, you can implement the `replace into` method for insertion. Currently, a universally applicable `insert` method is implemented by default, but for cases like using the `gora` driver for Oracle, the `insert` method may not be suitable.

```go
// Parameter represents an execution parameter with a table, transaction mode, and SQL statement.
type Parameter interface {
 Table() Table                                 // Table or view
 TxOptions() *sql.TxOptions                    // Transaction mode
 Query([]element.Record) (string, error)       // SQL prepare statement
 Agrs([]element.Record) ([]interface{}, error) // Prepare parameters
}
```

To implement the `replace into` method for insertion, you need to implement the `Parameter` interface. You can combine `BaseParam` to simplify the implementation. Refer to the implementation in `table.go` of the MySQL package.

### Field Interface

```go
// Field represents a database column.
type Field interface {
 fmt.Stringer

 Index() int                   // Index
 Name() string                 // Column name
 Quoted() string               // Quoted column name
 BindVar(int) string           // Placeholder symbol
 Select() string               // Select column name
 Type() FieldType              // Column type
 Scanner() Scanner             // Scanner, used to convert database data into a column
 Valuer(element.Column) Valuer // Valuer, used to convert a column into database data
}
```

When implementing the `Field` interface, you can combine `BaseField` to simplify the implementation. The `Type()` method must return the specific column type for the corresponding database. The `Scanner` must return the scanner for the corresponding database, and the `Valuer` must return the valuer for the corresponding database. Refer to the implementation in `field.go` of the MySQL package.

```go
// ColumnType represents a column type, abstracting `sql.ColumnType` and facilitating the implementation of corresponding functions.
type ColumnType interface {
 Name() string                                   // Column name
 ScanType() reflect.Type                         // Scan type
 Length() (length int64, ok bool)                // Length
 DecimalSize() (precision, scale int64, ok bool) // Precision
 Nullable() (nullable, ok bool)                  // Nullability
 DatabaseTypeName() string                       // Database type name
}

// FieldType represents a field type.
type FieldType interface {
 ColumnType

 IsSupportted() bool // Checks if the type is supported
}
```

When implementing the `FieldType` interface, you can combine `BaseFieldType` to simplify the implementation. `ColumnType` is essentially an abstraction of `sql.ColumnType`. Refer to the implementation in `field.go` of the MySQL package.

```go
// Scanner represents a column data scanner that converts database driver values into column data.
type Scanner interface {
 sql.Scanner

 Column() element.Column // Gets the column data
}
```

When implementing the `Scanner` interface, you can combine `BaseFieldType` to simplify the implementation. The role of the scanner is to convert the data read by the database driver into a single column of data. Refer to the implementation in `field.go` of the MySQL package.

```go
// Valuer represents a valuer that converts corresponding data into a value for the database driver.
type Valuer interface {
 driver.Valuer
}

// ValuerGoType is an optional functionality for the Field, used to determine the Go type for the corresponding driver value.
// It returns the corresponding value for the driver's value, facilitating the determination by GoValuer.
type ValuerGoType interface {
 GoType() GoType
}
```

When implementing the `Valuer` interface, you can combine `GoValuer` to simplify the implementation. To use `GoValuer`, you need to implement the `ValuerGoType` interface at the database layer's Field level. The role of the valuer is to convert a single column of data into the data type used for writing by the database driver. Refer to the implementation in `field.go` of the MySQL package.