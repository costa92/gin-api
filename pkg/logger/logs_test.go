package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Logs(t *testing.T) {
	opts := Options{
		Level:            "test",
		Format:           "test",
		EnableCaller:     false,
		EnableColor:      true,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	errs := opts.Validate()
	fmt.Println(fmt.Sprintf("%s", errs))
	expected := `[unrecognized level: "test" not a valid log format: "test"]`
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))
}

func Test_WithName(t *testing.T) {
	defer Flush()
	log := WithName("test")
	log.Infow("Hello world!", "foo", "dd")
}

func Test_Infow(t *testing.T) {
	Infow("Hello world!", "foo", "dd")
}
