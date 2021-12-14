package utils

// keep track of selections.
type SelectorWrapper struct {
	Selection map[int]struct{} // Keep track of which chapters have been selected by user.
	All       bool             // Keep track of whether user has selected All or not.
}

// HasSelections : Checks whether there are currently selections.
func (s *SelectorWrapper) HasSelections() bool {
	return len(s.Selection) != 0
}

// HasSelection : Checks whether the current row is selected.
func (s *SelectorWrapper) HasSelection(row int) bool {
	_, ok := s.Selection[row]
	return ok
}

// CopySelection : Returns a copy of the current Selection.
func (s *SelectorWrapper) CopySelection() map[int]struct{} {
	selection := map[int]struct{}{}
	for se := range s.Selection {
		selection[se] = struct{}{}
	}
	return selection
}

// AddSelection : Add a row to the Selection.
func (s *SelectorWrapper) AddSelection(row int) {
	s.Selection[row] = struct{}{}
}

// RemoveSelection : Remove a row from the Selection. No-op if row is not originally in Selection.
func (s *SelectorWrapper) RemoveSelection(row int) {
	delete(s.Selection, row)
}
