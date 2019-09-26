# Failable

Perform failable operations.

```go
done, failed := failable.Run(func(fail failable.FailFunc) {
  err := expensiveOperation()
  if err != nil {
    fail(err) // Once fail is called the execution of the operation is stopped.
  }

  nextExpensiveOperation()
})

select {
case <-done:
  // Completed Successfully
case err <-failed:
  panic(err) // The same error passed to fail
}
```
