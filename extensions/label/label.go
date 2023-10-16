package label

func BlankMetadata() *Metadata {
	return &Metadata{
		Labels: []*Label{},
	}
}

func (m *Metadata) MatchLabel(expectLabel *Label) bool {
	for _, l := range m.GetLabels() {
		if l.Value == expectLabel.Value && l.Key == expectLabel.Key {
			return true
		}
	}
	return false
}

func (m *Metadata) GetLabelByKey(key string) *Label {
	for _, l := range m.GetLabels() {
		if l.Key == key {
			return l
		}
	}
	return nil
}

func (m *Metadata) GetLabelsByKey(key string) []*Label {
	var labels []*Label
	for _, l := range m.GetLabels() {
		if l.Key == key {
			labels = append(labels, l)
		}
	}
	return labels
}

func (m *Metadata) RemoveLabel(expectLabel *Label) bool {
	var (
		i int
		l *Label
	)
	for i, l = range m.GetLabels() {
		if l.Key == expectLabel.Key && l.Value == expectLabel.Value {
			break
		}
	}
	m.Labels = append(m.Labels[i:], m.Labels[i+1:]...)
	return false
}

func (m *Metadata) AddLabel(newLabel *Label) error {
	if err := newLabel.Validate(); err != nil {
		return err
	}
	m.Labels = append(m.Labels, newLabel)
	return nil
}
