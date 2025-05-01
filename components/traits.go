package components

// Traits contains champion trait information
type Traits struct {
	list []string
}

// NewTraits creates a Traits component
func NewTraits(traitList []string) Traits {
	return Traits{
		list: traitList,
	}
}

// GetTraits returns the list of traits
func (t *Traits) GetTraits() []string {
	return t.list
}

// HasTrait checks if a specific trait is present
func (t *Traits) HasTrait(traitName string) bool {
	for _, t := range t.list {
		if t == traitName {
			return true
		}
	}
	return false
}

// AddTrait adds a trait to the list
func (t *Traits) AddTrait(traitName string) {
	for _, t := range t.list {
		if t == traitName {
			return // Trait already exists
		}
	}
	t.list = append(t.list, traitName)
}