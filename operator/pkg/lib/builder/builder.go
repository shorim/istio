package builder

import (
	"encoding/json"

	istioOperator "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
)

type istioOperatorBuilder struct {
	istioOperator istioOperator.IstioOperator

	// Stored internally and returned on call to Get()
	mergeError error
}

// Creates new istioOperatorBuilder
func NewIstioOperatorBuilder() *istioOperatorBuilder {
	return &istioOperatorBuilder{}
}

func (b *istioOperatorBuilder) SetBaseOperator(base istioOperator.IstioOperator) *istioOperatorBuilder {
	b.istioOperator = base
	return b
}

// Mergeable exposes Merge function, types that implement this interface should be able to append their configuration to IstioOperator
type Mergeable interface {
	// merge should merge any new configuration with the istioOperator parameter, overwriting if needed
	Merge(istioOperator.IstioOperator) (istioOperator.IstioOperator, error)
}

// MergeWith executes merge from supplied Mergeable with the builder stored IstioOperator as parameter
func (b *istioOperatorBuilder) MergeWith(toMerge ...Mergeable) *istioOperatorBuilder {
	for _, merge := range toMerge {
		out, err := merge.Merge(b.istioOperator)
		if err != nil {
			b.mergeError = err
			return b
		}
		b.istioOperator = out
	}
	return b
}

// BuildString returns the built IstioOperator marshaled to JSON string
func (b *istioOperatorBuilder) BuildString() (string, error) {
	if b.mergeError != nil {
		return "", b.mergeError
	}

	s, err := json.Marshal(b.istioOperator)
	if err != nil {
		return "", err
	}

	return string(s), nil
}

// BuildJSONByteArray returns the built IstioOperator marshaled to JSON []byte
func (b *istioOperatorBuilder) BuildJSONByteArray() ([]byte, error) {
	if b.mergeError != nil {
		return nil, b.mergeError
	}

	return json.Marshal(b.istioOperator)
}

// Build returns the built IstioOperator
func (b *istioOperatorBuilder) Build() (istioOperator.IstioOperator, error) {
	if b.mergeError != nil {
		return istioOperator.IstioOperator{}, b.mergeError
	}

	return b.istioOperator, nil
}
