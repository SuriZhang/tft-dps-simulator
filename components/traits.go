package components

// Traits contains champion trait information
type Traits struct {
    List []string
}

// NewTraits creates a Traits component
func NewTraits(traitList []string) Traits {
    return Traits{
        List: traitList,
    }
}