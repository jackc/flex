package flex_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/flex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShell(t *testing.T) {
	shell := flex.NewShell()
	shell.Register(&flex.Command{
		Name: "add",
		ParamsType: flex.NewType(func(tb flex.TypeBuilder) {
			tb.Field("a", flex.Require(), flex.Int32())
			tb.Field("b", flex.Require(), flex.Int32())
		}),
		ExecFunc: func(ctx context.Context, params *flex.Record) (map[string]interface{}, error) {
			a := params.Get("a").(int32)
			b := params.Get("b").(int32)
			sum := a + b

			return map[string]interface{}{"sum": sum}, nil
		},
	})

	response, err := shell.Exec(context.Background(), "add", map[string]interface{}{"a": 1, "b": 2})
	require.NoError(t, err)
	assert.EqualValues(t, 3, response["sum"])
}

func TestCommandExec(t *testing.T) {
	cmd := flex.Command{
		Name: "add",
		ParamsType: flex.NewType(func(tb flex.TypeBuilder) {
			tb.Field("a", flex.Require(), flex.Int32())
			tb.Field("b", flex.Require(), flex.Int32())
		}),
		ExecFunc: func(ctx context.Context, params *flex.Record) (map[string]interface{}, error) {
			a := params.Get("a").(int32)
			b := params.Get("b").(int32)
			sum := a + b

			return map[string]interface{}{"sum": sum}, nil
		},
	}

	response, err := cmd.Exec(context.Background(), map[string]interface{}{"a": 1, "b": 2})
	require.NoError(t, err)
	assert.EqualValues(t, 3, response["sum"])
}

func TestCommandExecParsesJSONIfOnlyExecJSONFuncAvailable(t *testing.T) {
	cmd := flex.Command{
		Name: "add",
		ParamsType: flex.NewType(func(tb flex.TypeBuilder) {
			tb.Field("a", flex.Require(), flex.Int32())
			tb.Field("b", flex.Require(), flex.Int32())
		}),
		ExecJSONFunc: func(ctx context.Context, params *flex.Record) ([]byte, error) {
			a := params.Get("a").(int32)
			b := params.Get("b").(int32)
			sum := a + b

			return []byte(fmt.Sprintf(`{"sum":%v}`, sum)), nil
		},
	}

	response, err := cmd.Exec(context.Background(), map[string]interface{}{"a": 1, "b": 2})
	require.NoError(t, err)
	assert.EqualValues(t, 3, response["sum"])
}

func TestCommandExecParsesJSONIfOnlyExecJSONFuncAvailableEmptyResponse(t *testing.T) {
	cmd := flex.Command{
		Name:       "add",
		ParamsType: flex.NewType(func(tb flex.TypeBuilder) {}),
		ExecJSONFunc: func(ctx context.Context, params *flex.Record) ([]byte, error) {
			return nil, nil
		},
	}

	response, err := cmd.Exec(context.Background(), map[string]interface{}{})
	require.NoError(t, err)
	assert.Nil(t, response)
}

func TestCommandExecWithNilParamsType(t *testing.T) {
	cmd := flex.Command{
		Name: "nop",
		ExecFunc: func(ctx context.Context, params *flex.Record) (map[string]interface{}, error) {
			return nil, nil
		},
	}

	response, err := cmd.Exec(context.Background(), map[string]interface{}{"a": 1, "b": 2})
	require.NoError(t, err)
	assert.Nil(t, response)
}

func TestCommandExecJSON(t *testing.T) {
	cmd := flex.Command{
		Name: "add",
		ParamsType: flex.NewType(func(tb flex.TypeBuilder) {
			tb.Field("a", flex.Require(), flex.Int32())
			tb.Field("b", flex.Require(), flex.Int32())
		}),
		ExecJSONFunc: func(ctx context.Context, params *flex.Record) ([]byte, error) {
			a := params.Get("a").(int32)
			b := params.Get("b").(int32)
			sum := a + b

			return []byte(fmt.Sprintf(`{"sum":%v}`, sum)), nil
		},
	}

	response, err := cmd.ExecJSON(context.Background(), map[string]interface{}{"a": 1, "b": 2})
	require.NoError(t, err)
	assert.Equal(t, []byte(`{"sum":3}`), response)
}

func TestCommandExecJSONMarshalsExecIfExecJSONUnavailable(t *testing.T) {
	cmd := flex.Command{
		Name: "add",
		ParamsType: flex.NewType(func(tb flex.TypeBuilder) {
			tb.Field("a", flex.Require(), flex.Int32())
			tb.Field("b", flex.Require(), flex.Int32())
		}),
		ExecFunc: func(ctx context.Context, params *flex.Record) (map[string]interface{}, error) {
			a := params.Get("a").(int32)
			b := params.Get("b").(int32)
			sum := a + b

			return map[string]interface{}{"sum": sum}, nil
		},
	}

	response, err := cmd.ExecJSON(context.Background(), map[string]interface{}{"a": 1, "b": 2})
	require.NoError(t, err)
	assert.Equal(t, []byte(`{"sum":3}`), response)
}

func TestCommandExecJSONMarshalsExecIfExecJSONUnavailableNilResponse(t *testing.T) {
	cmd := flex.Command{
		Name:       "add",
		ParamsType: flex.NewType(func(tb flex.TypeBuilder) {}),
		ExecFunc: func(ctx context.Context, params *flex.Record) (map[string]interface{}, error) {
			return nil, nil
		},
	}

	response, err := cmd.ExecJSON(context.Background(), map[string]interface{}{})
	require.NoError(t, err)
	assert.Nil(t, response)
}
