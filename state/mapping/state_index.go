package mapping

type (
	// StateIndex additional index of entity instance
	StateIndex struct {
		Name     string
		Uniq     bool
		Required bool
		Keyer    InstanceMultiKeyer // index can have multiple keys
	}

	// StateIndexDef additional index definition
	StateIndexDef struct {
		Name     string
		Fields   []string
		Required bool
		Multi    bool
		Keyer    InstanceMultiKeyer
	}
)
