package env

// DiffResult holds the comparison between two env files.
type DiffResult struct {
	Missing []string    // in source but not in target
	Extra   []string    // in target but not in source
	Changed []DiffEntry // different values
}

// DiffEntry represents a key with different values between files.
type DiffEntry struct {
	Key    string
	Source string
	Target string
}

// Compare finds differences between a source (example) and target (.env).
func Compare(source, target *EnvFile) *DiffResult {
	result := &DiffResult{}

	// Missing: in source but not in target
	for _, key := range source.Order {
		if _, exists := target.Vars[key]; !exists {
			result.Missing = append(result.Missing, key)
		}
	}

	// Extra: in target but not in source
	for _, key := range target.Order {
		if _, exists := source.Vars[key]; !exists {
			result.Extra = append(result.Extra, key)
		}
	}

	// Changed: in both but different values
	for _, key := range source.Order {
		if targetVal, exists := target.Vars[key]; exists {
			if source.Vars[key] != targetVal && source.Vars[key] != "" {
				result.Changed = append(result.Changed, DiffEntry{
					Key:    key,
					Source: source.Vars[key],
					Target: targetVal,
				})
			}
		}
	}

	return result
}

// HasDifferences returns true if there are any differences.
func (d *DiffResult) HasDifferences() bool {
	return len(d.Missing) > 0 || len(d.Extra) > 0 || len(d.Changed) > 0
}
