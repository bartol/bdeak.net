# SQL

all snippets work in PostgreSQL

<!-- 
update
delete 
create table
foreign key
current date and time
current date
current time
get timestamp for one hour from now
COALESE
decimal column
table alias
column alias
count rows
ON DELETE CASCADE
ON UPDATE CASCADE
UNIQUE
multiple primary keys
extract year from date
extract month from date
partition table (by hash?)
explain statement
<> operator
IS NULL 
IS NOT NULL operators
IS DISTINCT FROM
LIMIT
get random row (with offset)
TSVECTOR for full text searching
-->

## Statements

### Select column from all rows

    SELECT <column-name> FROM <table-name>;

### Select all columns from all rows

    SELECT * FROM <table-name>;

### Select rows conditionaly

    SELECT * FROM <table-name>
    WHERE <column-name> = <column-value>;

more operators for where statement can be found below

### Select all columns from table alias

    SELECT <table-alias>.*
    FROM <table-name> AS <table-alias>;

### Concatenate multiple queries

    (SELECT * FROM <table-name>)
    UNION ALL
    (SELECT * FROM <table-name>)

### Concatenate multiple queries and return only distinct rows

    (SELECT * FROM <table-name>)
    UNION
    (SELECT * FROM <table-name>)

### Insert row

    INSERT INTO <table-name> (<column-name>, <column-name>, <column-name>) 
    VALUES (<column-value>, <column-value>, <column-value>);

### Insert multiple rows

    INSERT INTO <table-name> (<column-name>, <column-name>, <column-name>) 
    VALUES (<column-value>, <column-value>, <column-value>),
           (<column-value>, <column-value>, <column-value>),
           (<column-value>, <column-value>, <column-value>);

## Operators

## Functions

## Table joins

### Shorter syntax if join columns have same name

    SELECT *
    FROM <table-name>
    <JOIN-TYPE> JOIN <foreign-table-name> USING (<column-name>);

### Join one to many

each resource can have one category, multiple resources can have save category

    SELECT * 
    FROM <table-name> AS <table-alias>
    INNER JOIN <foreign-table-name> AS <foreign-table-alias>
    ON <table-alias>.<column-name> = <foreign-table-alias>.<foreign-column-name>;

## Column types

### Integer

    <column-name> INTEGER

### Auto-incrementing integer

    <column-name> SERIAL

### Text

    <column-name> TEXT

### Money

    <column-name> MONEY

possible value: `8.99`

### Decimal number

    <column-name> NUMERIC(<precision>, <scale>)

example parameters:

- `<precision>` = `9`
- `<scale>` = `2`

these can store decimal up to 9999999.99

### Date

    <column-name> DATE

possible value: `2020-11-30`

### Time

    <column-name> TIME

possible value: `11:24:02`

### Date and time

    <column-name> TIMESTAMP

possible value: `2020-11-30 11:24:02`

### Boolean

    <column-name> BOOLEAN

possible values: `'t'`, `'n'`

## Column constraints

### Primary key

    <table-name>_id SERIAL PRIMARY KEY

### Not null

    <column-name> <COLUMN-TYPE> NOT NULL

### Foreign key

    <column-name> INTEGER,
    FOREIGN KEY (<column-name>) REFERENCES <foreign-table-name>(<foreign-column-name>)

### Default value

    <column-name> <COLUMN-TYPE> DEFAULT <column-value>

## Table relationships

### Many to many

intersection table is needed

for example each account can have multiple products and each product can have multiple owners

    CREATE TABLE accounts_products (
        account_id INTEGER,
        product_id INTEGER,
        FOREIGN KEY (account_id) REFERENCES accounts(account_id),
        FOREIGN KEY (product_id) REFERENCES products(product_id)
    );

### Enum column (emulated)

for example resource status

    CREATE TABLE <table-name>_status (
        status TEXT PRIMARY KEY
    );

    CREATE TABLE <table-name> (
        -- ...
        status TEXT,
        FOREIGN KEY (status) REFERENCES <table-name>_status(status)
    );

## Null

from SQL Antipatterns book:

Suppose Stan is thirty years old, while Oliver's age is unknown. If I ask
you whether Stan is older than Oliver, your only possible answer is "I
don't know." If I ask you whether Stan is the same age as Oliver, your
answer is also "I don't know." If I ask you what is the sum of Stan's age
and Oliver's age, your answer is the same.

null in scalar expressions:

| Expression         | Expected | Actual |
| ------------------ | -------- | ------ |
| NULL = 0           | TRUE     | NULL   |
| NULL = 12345       | FALSE    | NULL   |
| NULL <> 12345      | TRUE     | NULL   |
| NULL + 12345       | 12345    | NULL   |
| NULL \|\| 'string' | 'string' | NULL   |
| NULL = NULL        | TRUE     | NULL   |
| NULL <> NULL       | FALSE    | NULL   |

null in boolean expressions:

| Expression        | Expected | Actual |
| ----------------- | -------- | ------ |
| NULL AND TRUE     | FALSE    | NULL   |
| NULL AND FALSE    | FALSE    | FALSE  |
| NULL OR FALSE     | FALSE    | NULL   |
| NULL OR TRUE      | TRUE     | TRUE   |
| NOT (NULL)        | TRUE     | NULL   |

## Comments

    -- comment

## Recommendations

- Use table relationships instead of storing values in comma separated text
- Don't let data spawn metadata
- Avoid float if you can
- Use null to signify a missing value for any data type
- Define `NOT NULL` constraint on columns whose missing value could break application
- Always spell out all the columns you need, instead of relying on wildcards or implicit column lists
- Generate different salt for each password and save it in the same table

## Note sources

- <https://pragprog.com/titles/bksqla/sql-antipatterns/>