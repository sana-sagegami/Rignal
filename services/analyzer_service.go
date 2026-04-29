package services

import "auto-zen-backend/repositories"

type AnalyzerService struct {
    readinessRepo repositories.ReadinessRepository
    sleepRepo     repositories.SleepRepository
    summaryRepo   repositories.SummaryRepository
    apns          *apns.Sender
}

func NewAnalyzerService(
    rr repositories.ReadinessRepository,
    sr repositories.SleepRepository,
    smr repositories.SummaryRepository,
    apns *apns.Sender,
) *AnalyzerService {
    return &AnalyzerService{rr, sr, smr, apns}
}

func (s *AnalyzerService) RunDailyAnalysis(date time.Time) error {
    readiness, _ := s.readinessRepo.FindByDate(date)
    sleep, _     := s.sleepRepo.FindByDate(date)

    summary := s.analyze(readiness, sleep)

    if err := s.summaryRepo.Save(summary); err != nil {
        return err
    }

    // APNs で Swift に通知
    return s.apns.Push(summary)
}