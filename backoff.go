package backoff

import "time"

type Config[T any] struct {
	// Curve should be a function that returns an increasing value based on the
	// number of attempts and is used to determine how long in seconds to wait
	// before the next attempt.
	//
	// Eg. if Curve returns 10, the next attempt will be made in 10 seconds.
	Curve func(float64) float64
	// Func is the function that will be retried. It should return a value of
	// type *T and/or an error. If the error is not nil, the function will be
	// retried.
	Func func() (*T, error)
	// MaxAttempts is the maximum number of attempts to make before giving up.
	// If MaxAttempts is 0 the function will be retried indefinitely, and errors
	// will be logged but not returned.
	MaxAttempts int
	// If LogFailure is not nil, it will be called with the error returned by
	// Func each time it fails.
	LogFailure func(error)

	Result T
}

// Backoff will retry the function specified in the config until it returns a
// non-nil value or the maximum number of attempts is reached.
func Backoff[T any](conf Config[T]) (*T, []error) {
	if conf.Curve == nil || conf.Func == nil {
		return nil, []error{ErrInvalidConfig}
	}

	var errorChannel chan error
	result := make(chan *T, 1)
	errorList := []error{}

	if conf.MaxAttempts != 0 {
		errorChannel = make(chan error, conf.MaxAttempts)
	}

	go backoff[T](conf.Curve, conf.Func, conf.MaxAttempts, result, errorChannel, conf.LogFailure)

	res := <-result
	if res != nil {
		return res, nil
	}

	for len(errorChannel) > 0 {
		err := <-errorChannel
		errorList = append(errorList, err)
	}

	return nil, errorList
}

func backoff[T any](
	curve func(float64) float64,
	fn func() (*T, error),
	attempts int,
	result chan *T,
	errors chan error,
	logFailure func(error),
) {
	var res *T
	attempt := 0

	for attempts == 0 || attempt < attempts {
		<-time.After(time.Duration(curve(float64(attempt))) * time.Second)

		var err error
		res, err = fn()

		if err != nil {
			if logFailure != nil {
				logFailure(err)
			}
			if attempts != 0 && errors != nil {
				errors <- err
			}
		}
		if res != nil {
			// stop retrying
			break
		}
		attempt++
	}

	result <- res
}
