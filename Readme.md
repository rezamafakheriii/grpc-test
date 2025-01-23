# gRPC Test Project

This project demonstrates a simple gRPC setup with two servers (order and payment) and a client. The flow is as follows:

1. The client calls the Order Server.
2. The Order Server internally calls the Payment Server.

The purpose of this project is to showcase error handling in gRPC services using a custom error package from go-lib and logging errors in interceptors.

## Project Structure

- **Order Server**: Handles order requests and calls the Payment Server.
- **Payment Server**: Handles payment requests.
- **Client**: Calls the Order Server to place an order.

## Purpose

The main goal of this project is to demonstrate:

- **Error Handling**: Using a custom error package (go-lib) and extending it with specific errors.
- **Logging**: Logging errors returned by servers in interceptors.
- **Recovery**: Recovering from panics in interceptors to prevent server shutdown.

## How to Run the Project

1. **Initialize the project:**

    ```bash
    go mod tidy
    ```

2. **Start the Payment Server:**

    ```bash
    go run payment/main.go
    ```

3. **Start the Order Server:**

    ```bash
    go run order/order.go
    ```

4. **Run the Client:**

    ```bash
    go run client/main.go
    ```

The client will call the Order Server, and you can observe the error flow in each service.

## How the Interceptor Works

The interceptor performs the following tasks:

### 1. Recover from Panic

The interceptor ensures that the server does not shut down in case of a panic. It recovers from the panic and logs the error.

```go
defer func() {
    if r := recover(); r != nil {
        if debugMode {
            slog.Error("Recovered from panic",
                slog.String("method", info.FullMethod),
                slog.Any("error", r),
                slog.String("stack_trace", string(debug.Stack())),
            )
        } else {
            slog.Error("Recovered from panic",
                slog.String("method", info.FullMethod),
                slog.Any("error", r),
            )
        }

        err = recoverFrom(serviceName, r)
    }
}()
```

#### Testing Panic Recovery

Comment out the above code and throw a panic or runtime error to see the difference in behavior when the recovery code is present versus when it is not.

### 2. Handle Errors

If the error is of type `appError` (custom error type), it is converted to a gRPC error. For other errors (e.g., errors from the Payment Server), the error is logged.

```go
if appErr, ok := err.(*appError); ok {
    slog.Error("app error",
        slog.String("method", info.FullMethod),
        slog.String("error", appErr.Error()),
        slog.Any("stack_trace", appErr.StackTrace()),
    )
} else {
    slog.Error("unknown error",
        slog.String("method", info.FullMethod),
        slog.String("error", err.Error()),
    )
}
```
## Why Use `appError` Stack Trace?
The `slog` stack trace shows where the log is called, not the origin of the error. Using `appError.StackTrace()` ensures we capture the actual origin of the error.


## Example Error Flow
1. Client calls the Order Server to place an order.
2. Order Server calls the Payment Server to process the payment.
3. If the Payment Server returns an error, the Order Server logs the error and returns a response to the client.
4. If a panic occurs in any server, the interceptor recovers from it and logs the error.


Feel free to modify and expand this `README.md` as needed!