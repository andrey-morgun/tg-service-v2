package app

const (
	subjectSkeletonStatus = "subject.skeleton.status"
)

func (a *App) registerBrokerTopics() map[string]func([]byte) error {
	return map[string]func([]byte) error{
		subjectSkeletonStatus: a.statusBrokerHandler.CheckStatus,
	}
}
