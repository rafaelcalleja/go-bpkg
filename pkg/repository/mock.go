package repository

type mockReleaseVersionFinder struct {
	LatestFn func(organization string, name string) (string, error)
}

func newMockReleaseVersionFinder() *mockReleaseVersionFinder {
	return new(mockReleaseVersionFinder)
}

func (m *mockReleaseVersionFinder) Latest(organization string, name string) (string, error) {
	return m.LatestFn(organization, name)
}
