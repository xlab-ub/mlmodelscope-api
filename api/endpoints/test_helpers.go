package endpoints

type nullWriter struct{}

func (w *nullWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
